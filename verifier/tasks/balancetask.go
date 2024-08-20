package tasks

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/Cedro-Finance/CedroAVS/api/graphql"
	"github.com/Cedro-Finance/CedroAVS/core/tasks"
	"github.com/Cedro-Finance/CedroAVS/verifier/client/agents"
	"github.com/Cedro-Finance/CedroAVS/verifier/common"
	"github.com/Cedro-Finance/CedroAVS/verifier/opio"
	"github.com/Layr-Labs/eigensdk-go/logging"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

var MAX_REQID, _ = new(big.Int).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF", 16)

type OperatorDeviationTask struct {
	req     *agents.DeviationRequest
	retries uint64
}

type DeviationError struct {
	TaskAgent ethcommon.Address
	Error     error
}

type BalanceTaskProcessor struct {
	tasks  map[ethcommon.Address](*OperatorDeviationTask)
	taskMu sync.RWMutex

	signCh chan *tasks.SignedDeviationTaskRespRequest
	errCh  chan *DeviationError

	agClients map[string]common.IAgentPodClient
	logger    logging.Logger
	opio      *opio.OperatorIO

	reqCh chan agents.TaskRequest
}

var _ Taskprocessor = (*BalanceTaskProcessor)(nil)

func NewBalanceTask(opioClient *opio.OperatorIO, agClients map[string]common.IAgentPodClient, logger logging.Logger) *BalanceTaskProcessor {
	return &BalanceTaskProcessor{
		tasks:     make(map[ethcommon.Address]*OperatorDeviationTask),
		signCh:    make(chan *tasks.SignedDeviationTaskRespRequest, 10),
		agClients: agClients,
		logger:    logger,
		reqCh:     make(chan agents.TaskRequest),
		opio:      opioClient,
	}
}

func (r *BalanceTaskProcessor) Start(ctx context.Context, wg *sync.WaitGroup) {

	r.logger.Info("Deviation Task Processor Started")

	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case req := <-r.reqCh:
			r.logger.Info("Deviation Check request event received", "req", req)

			reqData := req.(*agents.DeviationRequest)

			newTask, err := r.CreateBalanceTask(reqData)

			if err != nil {
				r.logger.Error("Error while creating new task")
			}

			go r.processDeviationCheckRequest(ctx, newTask)

		case sign := <-r.signCh:
			r.logger.Info("Signature Received")

			go r.sendSignedTaskToAggregator(ctx, sign)

		case err := <-r.errCh:
			r.logger.Error("Errror occured", err)
		}
	}
}

func (r *BalanceTaskProcessor) GetRequestCh() chan agents.TaskRequest {
	return r.reqCh
}

func (r *BalanceTaskProcessor) GetSignatureCh() chan *tasks.SignedDeviationTaskRespRequest {
	return r.signCh
}

func (r *BalanceTaskProcessor) isOffchainTrigger(reqId *big.Int) bool {

	return reqId.Cmp(MAX_REQID) == 0
}

func (r *BalanceTaskProcessor) GetTaskByAgent(agent ethcommon.Address) *OperatorDeviationTask {
	r.taskMu.RLock()
	defer r.taskMu.RUnlock()

	res := r.tasks[agent]

	return res
}

func (r *BalanceTaskProcessor) CreateBalanceTask(req *agents.DeviationRequest) (*OperatorDeviationTask, error) {
	task := r.GetTaskByAgent(req.AgentAddress)

	if task == nil {
		r.taskMu.Lock()

		task = &OperatorDeviationTask{
			req:     req,
			retries: 0,
		}

		r.tasks[req.AgentAddress] = task
		r.taskMu.Unlock()
	}

	return task, nil
}

func (r *BalanceTaskProcessor) getAgentPodsForAgent(ctx context.Context, agentAddress ethcommon.Address) []common.IAgentPodClient {

	// agentRegisteredChains := []string{"ETH"}

	agentRegisteredChains, err := r.opio.GraphClient.GetAgentPodsChainLists(ctx, agentAddress.Hex())

	if err != nil {
		return []common.IAgentPodClient{}
	}

	agentList := make([]common.IAgentPodClient, len(agentRegisteredChains.ChainIds))

	for index, chain := range agentRegisteredChains.ChainIds {
		client := r.agClients[chain]

		agentPod, err := client.GetGraphClient().GetAgentSummaries(ctx, agentAddress.Hex())

		if err != nil {
			r.logger.Error("Error getting agentPods address from graph", "err", err)
			return nil
		}

		agentPod.ChainId = chain

		client.InitializeAgentPod(agentPod.AgentContract)

		agentList[index] = client
	}

	return agentList
}

func (b *BalanceTaskProcessor) processDeviationCheckRequest(ctx context.Context, task *OperatorDeviationTask) {
	//Now for every balance request event, there should be two jobs running in parallel
	//Routine 1: Fetch the balances of agents in all the chains to get the actual Balance
	//Routine 2: Fetch the total deposits & aprs of agents to get the ideal balance

	// deviationPercentage := big.NewInt(0)

	reqInfo := &tasks.BalanceTaskInfo{}
	var err error
	//Offchain Request
	if !b.isOffchainTrigger(task.req.ReqId) {
		reqInfo, err = b.opio.GRpcClient.CreateBalanceTaskToAggregator(task.req)

		if err != nil {
			b.logger.Error("Error while creating task to aggregator")
			b.errCh <- &DeviationError{
				TaskAgent: task.req.AgentAddress,
				Error:     fmt.Errorf("error while getting actual balance"),
			}
		}
	}

	idealBalCh := make(chan *big.Int)
	actualBalCh := make(chan *big.Int)

	go b.fetchAgentActualBalance(ctx, task, actualBalCh)
	go b.fetchAgentIdealBalance(ctx, task, idealBalCh)

	totalIdealBalance := <-idealBalCh
	totalActualBalance := <-actualBalCh

	if totalActualBalance == nil {
		b.errCh <- &DeviationError{
			TaskAgent: task.req.AgentAddress,
			Error:     fmt.Errorf("error while getting actual balance"),
		}

		return
	}

	if totalIdealBalance == nil {
		b.errCh <- &DeviationError{
			TaskAgent: task.req.AgentAddress,
			Error:     fmt.Errorf("error while getting ideal balance"),
		}

		return
	}

	deviation := calculateDeviationPercentage(totalActualBalance, totalIdealBalance)

	if deviation.Cmp(big.NewFloat(0)) > 0 {
		b.logger.Info("Deviation Found", "deviation", deviation.String())
		//Deviation Occured
		reqInfo, err = b.opio.GRpcClient.CreateBalanceTaskToAggregator(task.req)

		if err != nil {
			b.logger.Error("Error while creating task to aggregator")
			b.errCh <- &DeviationError{
				TaskAgent: task.req.AgentAddress,
				Error:     fmt.Errorf("error while getting actual balance"),
			}
		}

	}

	b.logger.Info("ProcessDeviationCheck",
		"RequestId", reqInfo.ReqId.String(),
		"idealBalance", totalIdealBalance.String(),
		"totalBalance", totalActualBalance.String(),
	)

	reqInfo.ActualBalance = totalActualBalance
	reqInfo.IdealBalance = totalIdealBalance
	reqInfo.ReferenceBlockNumber = task.req.ReferenceBlockNumber
	reqInfo.ReqId = task.req.ReqId
	reqInfo.AgentAddress = task.req.AgentAddress

	if (b.isOffchainTrigger(task.req.ReqId) && deviation.Cmp(big.NewFloat(0)) > 0) || (!b.isOffchainTrigger(task.req.ReqId)) {
		go b.signBalanceRespnse(reqInfo)
	}
}

func (r *BalanceTaskProcessor) fetchAgentActualBalance(ctx context.Context, task *OperatorDeviationTask, balCh chan<- *big.Int) {

	//Actual Balance is the total accumulated balance of the agent pod in all the chains for given agent address in event
	//reqEvent contains tht agent address, which will have agent pod contract in all the enrolled chains

	// agentAddr := req.AgentAddress
	agentPods := r.getAgentPodsForAgent(ctx, task.req.AgentAddress)

	if agentPods == nil {
		r.errCh <- &DeviationError{
			TaskAgent: task.req.AgentAddress,
			Error:     fmt.Errorf("error while getting agentPod contract address for given agent"),
		}
	}

	balanceCh := make([]chan *big.Int, len(agentPods))

	for index := range balanceCh {
		balanceCh[index] = make(chan *big.Int)
	}

	for index, agentPod := range agentPods {
		go agentPod.GetBalance(ctx, task.req, balanceCh[index])
	}

	totalBalance := big.NewInt(0)

	for _, ch := range balanceCh {
		data := <-ch
		if data == nil {
			balCh <- nil
			return
		}
		totalBalance = totalBalance.Add(totalBalance, data)
	}

	balCh <- totalBalance
}

func (r *BalanceTaskProcessor) fetchAgentIdealBalance(ctx context.Context, task *OperatorDeviationTask, balCh chan<- *big.Int) {

	agentRegisteredChains, err := r.opio.GraphClient.GetAgentPodsChainLists(ctx, task.req.AgentAddress.Hex())

	if err != nil {
		r.logger.Error("Error getting agents registered chains list", "err", err)
		balCh <- nil
		return
	}

	totalAmount := big.NewInt(0)

	for _, chain := range agentRegisteredChains.ChainIds {
		client := r.agClients[chain]

		deposits, err := client.GetGraphClient().GetUnclaimedDepositsLists(ctx, task.req.AgentAddress.Hex(), task.req.Timestamp)

		if err != nil {
			r.logger.Error("Error getting deposits address from graph", "err", err)
			balCh <- nil

			return
		}

		for _, deposit := range deposits {

			interest := CalculateInterest(&deposit, task.req.Timestamp)

			totalAmount = totalAmount.Add(totalAmount, interest)
			totalAmount = totalAmount.Add(totalAmount, deposit.PrincipalAmount)
		}
	}

	balCh <- totalAmount
}

func calculateDeviationPercentage(totalActualBalance, totalIdealBalance *big.Int) *big.Float {
	// Calculate the deviation: totalActualBalance - totalIdealBalance
	deviation := new(big.Int).Sub(totalActualBalance, totalIdealBalance)

	// Convert big.Int to big.Float for precise division
	deviationFloat := new(big.Float).SetInt(deviation)
	idealBalanceFloat := new(big.Float).SetInt(totalIdealBalance)

	// Calculate the deviation percentage: (deviation / totalIdealBalance) * 100
	deviationPercentage := new(big.Float).Quo(deviationFloat, idealBalanceFloat)
	deviationPercentage.Mul(deviationPercentage, big.NewFloat(100))

	return deviationPercentage
}

func CalculateInterest(depositInfo *graphql.DepositInfo, timestamp uint64) *big.Int {
	// Convert block timestamp to seconds
	blockTimestamp := time.Unix(int64(depositInfo.BlockTimestamp), 0)
	// Convert time interval to seconds
	timeInterval := big.NewInt(1).Exp(big.NewInt(10), big.NewInt(18), nil)
	timeInterval = timeInterval.Quo(timeInterval, depositInfo.TimeInterval)
	// Calculate maturity period
	maturityPeriod := blockTimestamp.Add(time.Duration(timeInterval.Int64()) * time.Second)
	// Determine whether to calculate interest up to maturity period or current timestamp
	var interestTime uint64
	if maturityPeriod.After(time.Unix(int64(timestamp), 0)) {
		interestTime = uint64(maturityPeriod.Unix())
	} else {
		interestTime = timestamp
	}
	// Calculate interest
	interest := big.NewInt(0)
	timeDiff := big.NewInt(int64(interestTime - depositInfo.BlockTimestamp))
	years := timeDiff.Int64() / (365 * 24 * 60 * 60)
	interest.Mul(depositInfo.PrincipalAmount, depositInfo.InterestRate)
	interest.Mul(interest, big.NewInt(int64(years)))
	return interest
}

func (b *BalanceTaskProcessor) signBalanceRespnse(balHeader *tasks.BalanceTaskInfo) {
	b.logger.Info("Sending Balance Header", "header", balHeader)

	hash, err := balHeader.GetHeaderHash()

	if err != nil {
		b.logger.Error("Error while hashing the messsage")
	}

	signature := b.opio.Signer.SignMessage(hash)

	b.logger.Info("Signed the message", signature)

	b.signCh <- &tasks.SignedDeviationTaskRespRequest{
		Balance:      *balHeader,
		BlsSignature: *signature,
	}
}

func (b *BalanceTaskProcessor) sendSignedTaskToAggregator(ctx context.Context, req *tasks.SignedDeviationTaskRespRequest) {
	maxRetry := 10
	retries := 0

RETRY_LOOP:
	for {
		response, err := b.opio.GRpcClient.SendSignedTaskResponse(ctx, req)

		if err != nil {
			b.logger.Error("Error while sending Signed Task response", "err", err)

			b.errCh <- &DeviationError{
				TaskAgent: req.Balance.AgentAddress,
				Error:     err,
			}

			if retries < maxRetry {
				b.logger.Info("Retrying sending signature to aggregator")
				retries++
				break RETRY_LOOP
			}

			return

		} else {
			b.taskMu.Lock()
			delete(b.tasks, req.Balance.AgentAddress)
			b.taskMu.Unlock()

			b.logger.Info("Successfully posted the signatures to the aggregator",
				"task", req.Balance.TaskIndex,
				"reqId", req.Balance.ReqId.String(),
				"agent", req.Balance.AgentAddress.Hex(),
				"response", response)

			return
		}
	}
}

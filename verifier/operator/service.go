package operator

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/Cedro-Finance/CedroAVS/api/graphql"
	// "github.com/Cedro-Finance/CedroAVS/core"

	// "github.com/Cedro-Finance/CedroAVS/verifier/client/avs"
	// "github.com/ethereum/go-ethereum/accounts/abi/bind"

	ethcommon "github.com/ethereum/go-ethereum/common"

	// "github.com/Cedro-Finance/CedroAVS/verifier/client/avs"
	"github.com/Cedro-Finance/CedroAVS/verifier/client/avs"
	"github.com/Cedro-Finance/CedroAVS/verifier/common"
)

func (r *Operator) getAgentPodsForAgent(ctx context.Context, agentAddress ethcommon.Address) []common.IAgentPodClient {

	// agentRegisteredChains := []string{"ETH"}

	agentRegisteredChains, err := r.graphClient.GetAgentPodsChainLists(ctx, agentAddress.Hex())

	if err != nil {
		return []common.IAgentPodClient{}
	}

	agentList := make([]common.IAgentPodClient, len(agentRegisteredChains.ChainIds))

	for index, chain := range agentRegisteredChains.ChainIds {
		client := r.chainAgents[chain]

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

func (r *Operator) ProcessDeviationCheckRequest(ctx context.Context, task *OperatorDeviationTask) {

	//Now for every balance request event, there should be two jobs running in parallel
	//Routine 1: Fetch the balances of agents in all the chains to get the actual Balance
	//Routine 2: Fetch the total deposits & aprs of agents to get the ideal balance

	idealBalCh := make(chan *big.Int)
	actualBalCh := make(chan *big.Int)

	go r.FetchAgentActualBalance(ctx, task, actualBalCh)
	go r.FetchAgentIdealBalance(ctx, task, idealBalCh)

	totalIdealBalance := <-idealBalCh
	totalActualBalance := <-actualBalCh

	if totalActualBalance == nil {
		r.errCh <- &DeviationError{
			TaskIndex: task.reqInfo.TaskIndex,
			Error:     fmt.Errorf("error while getting actual balance"),
		}

		return
	}

	if totalIdealBalance == nil {
		r.errCh <- &DeviationError{
			TaskIndex: task.reqInfo.TaskIndex,
			Error:     fmt.Errorf("error while getting ideal balance"),
		}

		return
	}

	r.logger.Info("ProcessDeviationCheck",
		"RequestId", task.reqInfo.ReqId.String(),
		"idealBalance", totalIdealBalance.String(),
		"totalBalance", totalActualBalance.String(),
	)

	task.reqInfo.ActualBalance = totalActualBalance
	task.reqInfo.IdealBalance = totalIdealBalance
	task.reqInfo.ReferenceBlockNumber = task.req.ReferenceBlockNumber
	task.reqInfo.ReqId = task.req.ReqId
	task.reqInfo.AgentAddress = task.req.AgentAddress

	go r.SignAndSendBalanceResponse(ctx, task.reqInfo)
}

func (r *Operator) FetchAgentActualBalance(ctx context.Context, task *OperatorDeviationTask, balCh chan<- *big.Int) {

	//Actual Balance is the total accumulated balance of the agent pod in all the chains for given agent address in event
	//reqEvent contains tht agent address, which will have agent pod contract in all the enrolled chains

	// agentAddr := req.AgentAddress
	agentPods := r.getAgentPodsForAgent(ctx, task.req.AgentAddress)

	if agentPods == nil {
		r.errCh <- &DeviationError{
			TaskIndex: task.reqInfo.TaskIndex,
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

func (r *Operator) FetchAgentIdealBalance(ctx context.Context, task *OperatorDeviationTask, balCh chan<- *big.Int) {

	agentRegisteredChains, err := r.graphClient.GetAgentPodsChainLists(ctx, task.req.AgentAddress.Hex())

	if err != nil {
		r.logger.Error("Error getting agents registered chains list", "err", err)
		balCh <- nil
		return
	}

	totalAmount := big.NewInt(0)

	for _, chain := range agentRegisteredChains.ChainIds {
		client := r.chainAgents[chain]

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

func (r *Operator) StartOffChainDeviationCheck(ctx context.Context) {

	r.logger.Info("Offchain trigger service started")

	// blockDetails, err := r.avs.Eth.BlockByNumber(ctx, big.NewInt(int64(1867349)))

	// if err != nil {
	// 	r.logger.Error("OffchainDeviationCheckTrigger errro geting block details", "err", err)

	// }

	// r.balanceReqCh <- &avs.Request{
	// 	ReqId:                big.NewInt(6),
	// 	ReferenceBlockNumber: 1867349		,
	// 	AgentAddress:         ethcommon.HexToAddress("0xdF1Cf915dB7A658dFDEb75aF3873B4e55F25Aa42"),
	// 	Timestamp:            blockDetails.Time(),
	// 	// IsOffChain:           true,
	// }

	dur := time.Duration(time.Minute * 1)

	ticker := time.NewTicker(dur)

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			r.logger.Infof("Offchain deviation trigger period")
		}

	RETRY_LOOP:
		for {
			agents, err := r.graphClient.GetEnrolledAgentsList(ctx)

			if err != nil {
				r.logger.Error("Error while getting all the enrolled agents in agent manager")
				break RETRY_LOOP
			}

			//Current Block Number

			blockNumber, err := r.avs.Eth.BlockNumber(ctx)

			if err != nil {
				r.logger.Error("OffchainDeviationCheckTrigger error while getting block number", "err", err)

				break RETRY_LOOP
			}

			blockDetails, err := r.avs.Eth.BlockByNumber(ctx, big.NewInt(int64(blockNumber)))

			if err != nil {
				r.logger.Error("OffchainDeviationCheckTrigger errro geting block details", "err", err)

				break RETRY_LOOP
			}

			for _, agent := range agents {
				maxUint256 := new(big.Int)
				maxUint256.SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF", 16)
				r.balanceReqCh <- &avs.Request{
					ReqId:                maxUint256,
					ReferenceBlockNumber: blockNumber,
					AgentAddress:         agent,
					Timestamp:            blockDetails.Time(),
					IsOffChain:           true,
				}
			}

			break
		}

	}

}

func (r *Operator) ProcessOffchainDeviationCheckRequest(ctx context.Context, task *OperatorDeviationTask) {

	//Now for every balance request event, there should be two jobs running in parallel
	//Routine 1: Fetch the balances of agents in all the chains to get the actual Balance
	//Routine 2: Fetch the total deposits & aprs of agents to get the ideal balance

	// Find the deviation threshold

	deviationPercentage, err := r.agentManager.GetDeviationPercentage(ctx, task.req.AgentAddress, task.req.ReferenceBlockNumber)

	if err != nil {
		r.logger.Error("Error while getting deviation percentage", "err", err)

	}

	idealBalCh := make(chan *big.Int)
	actualBalCh := make(chan *big.Int)

	go r.FetchAgentActualBalance(ctx, task, actualBalCh)
	go r.FetchAgentIdealBalance(ctx, task, idealBalCh)

	totalIdealBalance := <-idealBalCh
	totalActualBalance := <-actualBalCh

	if totalActualBalance == nil {
		r.errCh <- &DeviationError{
			TaskIndex: task.reqInfo.TaskIndex,
			Error:     fmt.Errorf("error while getting actual balance"),
		}

		return
	}

	if totalIdealBalance == nil {
		r.errCh <- &DeviationError{
			TaskIndex: task.reqInfo.TaskIndex,
			Error:     fmt.Errorf("error while getting ideal balance"),
		}

		return
	}

	r.logger.Info("ProcessDeviationCheck",
		"RequestId", task.reqInfo.ReqId.String(),
		"idealBalance", totalIdealBalance.String(),
		"totalBalance", totalActualBalance.String(),
	)

	task.reqInfo.ActualBalance = totalActualBalance
	task.reqInfo.IdealBalance = totalIdealBalance
	task.reqInfo.ReferenceBlockNumber = task.req.ReferenceBlockNumber
	task.reqInfo.ReqId = task.req.ReqId
	task.reqInfo.AgentAddress = task.req.AgentAddress

	devPercentFloat, _ := deviationPercentage.Float64()

	deviation := calculateDeviationPercentage(totalActualBalance, totalIdealBalance)

	if deviation.Cmp(big.NewFloat(devPercentFloat)) == 1 { // deviation is greater than threshold 10%

		r.logger.Info("deviation occured", "agent", task.req.AgentAddress.Hex(), "blockNumber", task.req.ReferenceBlockNumber)

		go r.SignAndSendBalanceResponse(ctx, task.reqInfo)
	} else {
		r.logger.Info("No deviation")
		go r.SignAndSendBalanceResponse(ctx, task.reqInfo)

	}

}

package operator

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/Cedro-Finance/CedroAVS/core"
	"github.com/Cedro-Finance/CedroAVS/core/logger"
	"github.com/Cedro-Finance/CedroAVS/verifier/client/avs"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/Cedro-Finance/CedroAVS/verifier/common"
)

func (r *Operator) getAgentPodsForAgent(ctx context.Context, agentAddress ethcommon.Address) []common.IAgentPodClient {

	agentRegisteredChains := []string{"ETH"}
	agentList := make([]common.IAgentPodClient, len(agentRegisteredChains))

	for index, chain := range agentRegisteredChains {
		client := r.chainAgents[chain]

		agentPod, err := client.GetGraphClient().GetAgentPodsForAgent(ctx, agentAddress.Hex())

		if err != nil {
			r.logger.Error("Error getting agentPods address from graph", "err", err)
			return nil
		}

		client.InitializeAgentPod(agentPod.AgentContract)

		agentList[index] = client
	}

	return agentList
}

func (r *Operator) ProcessDeviationCheckRequest(ctx context.Context, eventCh <-chan *avs.Request) {
	for {

		reqEvent := <-eventCh

		//Now for every balance request event, there should be two jobs running in parallel
		//Routine 1: Fetch the balances of agents in all the chains to get the actual Balance
		//Routine 2: Fetch the total deposits & aprs of agents to get the ideal balance

		reqInfo, err := r.gRpcClient.CreateBalanceTaskToAggregator(reqEvent)

		if err != nil {
			r.logger.Error("Create Task To Aggregator Failed. Trying again", "err", err)
			// continue
		}

		r.logger.Infof("Create Task successful", reqInfo)

		idealBalCh := make(chan *big.Int)
		actualBalCh := make(chan *big.Int)

		go r.FetchAgentActualBalance(ctx, reqEvent, actualBalCh)
		go r.FetchAgentIdealBalance(ctx, reqEvent, idealBalCh)

		totalIdealBalance := <-idealBalCh
		totalActualBalance := <-actualBalCh

		r.logger.Infof("ProcessDeviationCheck",
			"RequestId", reqEvent.ReqId.String(),
			"idealBalance", totalIdealBalance.String(),
			"totalBalance", totalActualBalance.String(),
		)

		reqInfo.ActualBalance = totalActualBalance
		reqInfo.IdealBalance = totalIdealBalance
		reqInfo.ReferenceBlockNumber = reqEvent.ReferenceBlockNumber
		reqInfo.ReqId = reqEvent.ReqId
		reqInfo.AgentAddress = reqEvent.AgentAddress

		// balResponse := &agents.AggregateBalanceResponse{
		// 	ReqId:                reqEvent.ReqId,
		// 	ActualBalance:        totalActualBalance,
		// 	IdealBalance:         totalIdealBalance,
		// 	ReferenceBlockNumber: reqEvent.ReferenceBlockNumber,
		// }

		go r.SignAndSendBalanceResponse(ctx, reqInfo)

	}
}

func (r *Operator) FetchAgentActualBalance(ctx context.Context, req *avs.Request, balCh chan<- *big.Int) {

	//Actual Balance is the total accumulated balance of the agent pod in all the chains for given agent address in event
	//reqEvent contains tht agent address, which will have agent pod contract in all the enrolled chains

	// agentAddr := req.AgentAddress
	agentPods := r.getAgentPodsForAgent(ctx, req.AgentAddress)

	if agentPods == nil {
		return
	}

	balanceCh := make([]chan *big.Int, len(agentPods))

	for index := range balanceCh {
		balanceCh[index] = make(chan *big.Int)
	}

	for index, agentPod := range agentPods {
		go agentPod.GetBalance(ctx, req, balanceCh[index])
	}

	totalBalance := big.NewInt(0)

	for _, ch := range balanceCh {
		data := <-ch
		totalBalance = totalBalance.Add(totalBalance, data)
	}

	balCh <- totalBalance
}

func (r *Operator) FetchAgentIdealBalance(ctx context.Context, req *avs.Request, balCh chan<- *big.Int) {
	// Call the graph data from here
	time.Sleep(time.Second * 5)
	balCh <- big.NewInt(5)
}

func (r *Operator) SignAndSendBalanceResponse(ctx context.Context, balHeader *core.BalanceTaskInfo) {

	r.logger.Info("Balance Heaader", "header", balHeader)

	hash, err := balHeader.GetBalanceHeaderHash()

	if err != nil {
		logger.GetLogger().Error("Error while hasing the message", "err", err)
		return
	}

	signature := r.blsKeyPair.SignMessage(hash)

	r.logger.Info("Signed the message", "sign", signature)

	request := &core.SignedTaskRespRequest{
		Balance:      *balHeader,
		BlsSignature: *signature,
		OperatorId:   r.operatorId,
	}

	err = r.gRpcClient.SendSignedTaskResponse(request)

	if err != nil {
		r.logger.Errorf("Error while sending Signed Task response")
	}
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

func (r *Operator) ProcessOffchainDeviationCheckRequest(ctx context.Context, eventCh <-chan *avs.Request) {
	for {
		reqEvent := <-eventCh

		//Now for every balance request event, there should be two jobs running in parallel
		//Routine 1: Fetch the balances of agents in all the chains to get the actual Balance
		//Routine 2: Fetch the total deposits & aprs of agents to get the ideal balance

		idealBalCh := make(chan *big.Int)
		actualBalCh := make(chan *big.Int)

		go r.FetchAgentActualBalance(ctx, reqEvent, actualBalCh)
		go r.FetchAgentIdealBalance(ctx, reqEvent, idealBalCh)

		totalIdealBalance := <-idealBalCh
		totalActualBalance := <-actualBalCh

		// Find the deviation threshold

		deviationPercentage, err := r.agentManager.Contract.DeviationPercentage(&bind.CallOpts{}, reqEvent.AgentAddress)

		if err != nil {
			r.logger.Error("Error while getting deviation percentage", "err", err)
		}

		fmt.Println(deviationPercentage)

		deviation := calculateDeviationPercentage(totalActualBalance, totalIdealBalance)

		if deviation.Cmp(big.NewFloat(10.0)) == 1 { // deviation is greater than threshold 10%
			// balResponse := &agents.AggregateBalanceResponse{
			// 	ReqId:                reqEvent.ReqId,
			// 	ActualBalance:        totalActualBalance,
			// 	IdealBalance:         totalIdealBalance,
			// 	ReferenceBlockNumber: reqEvent.ReferenceBlockNumber,
			// }

			// go r.SignAndSendBalanceResponse(ctx, balResponse)
			fmt.Println("")
		}

	}
}

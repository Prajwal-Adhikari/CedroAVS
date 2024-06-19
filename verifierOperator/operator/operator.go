package operator

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"sync"
	"time"

	"github.com/Cedro-Finance/CedroAVS/core"
	// "github.com/Cedro-Finance/CedroAVS/verifierOperator/client/agents"
	"github.com/Cedro-Finance/CedroAVS/verifierOperator/client/avs"
	"github.com/Cedro-Finance/CedroAVS/verifierOperator/common"
	"github.com/Cedro-Finance/CedroAVS/verifierOperator/config"
	"github.com/Cedro-Finance/CedroAVS/verifierOperator/keys"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/eth"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/wallet"
	"github.com/Layr-Labs/eigensdk-go/chainio/txmgr"
	sdkEcdsa "github.com/Layr-Labs/eigensdk-go/crypto/ecdsa"
	"github.com/Layr-Labs/eigensdk-go/signerv2"
	sdktypes "github.com/Layr-Labs/eigensdk-go/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	ethCommon "github.com/ethereum/go-ethereum/common"

	bls "github.com/Layr-Labs/eigensdk-go/crypto/bls"
	sdklogging "github.com/Layr-Labs/eigensdk-go/logging"

	"github.com/Cedro-Finance/CedroAVS/core/logger"
)

var RegisteredClients = map[string]func(config.ChainConfig, sdklogging.Logger) common.IAgentPodClient{}

type Operator struct {
	avs        *avs.AVSClient
	keystore   *keys.KeyStore
	blsKeyPair *bls.KeyPair

	agentManager *avs.AgentManager

	operatorId                sdktypes.OperatorId
	txMgr                     txmgr.TxManager
	logger                    sdklogging.Logger
	gRPCAggreagatorIpPortAddr string
	gRpcClient                *AggregatorGRpcClient

	chainAgents map[string]common.IAgentPodClient
}

type AgentPodListMap struct {
	Addr  string
	Chain string
}

func (r *Operator) mockAgentsForAgents(_ any) []common.IAgentPodClient {

	agentConfigList := []AgentPodListMap{
		{
			Addr:  "0x520D2e695c211AC600BDB41AF339D694c8b76425",
			Chain: "ethereum",
		},
	}

	agentList := make([]common.IAgentPodClient, len(agentConfigList))

	for index, agent := range agentConfigList {
		client := r.chainAgents[agent.Chain]
		client.InitializeAgentPod(ethCommon.HexToAddress(agent.Addr))

		agentList[index] = client
	}

	return agentList
}

func (r *Operator) FetchAgentActualBalance(ctx context.Context, req *avs.Request, balCh chan<- *big.Int) {

	//Actual Balance is the total accumulated balance of the agent pod in all the chains for given agent address in event
	//reqEvent contains tht agent address, which will have agent pod contract in all the enrolled chains

	// agentAddr := req.AgentAddress
	agentPods := r.mockAgentsForAgents(req.AgentAddress)

	fmt.Println(agentPods)

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

// func (r *Operator) ProcessOffchainDeviationCheckRequest(ctx context.Context, eventCh <-chan *avs.Request) {
// 	for {
// 		reqEvent := <-eventCh

// 		//Now for every balance request event, there should be two jobs running in parallel
// 		//Routine 1: Fetch the balances of agents in all the chains to get the actual Balance
// 		//Routine 2: Fetch the total deposits & aprs of agents to get the ideal balance

// 		idealBalCh := make(chan *big.Int)
// 		actualBalCh := make(chan *big.Int)

// 		go r.FetchAgentActualBalance(ctx, reqEvent, actualBalCh)
// 		go r.FetchAgentIdealBalance(ctx, reqEvent, idealBalCh)

// 		totalIdealBalance := <-idealBalCh
// 		totalActualBalance := <-actualBalCh

// 		// Find the deviation threshold

// 		deviation := calculateDeviationPercentage(totalActualBalance, totalIdealBalance)

// 		if deviation.Cmp(big.NewFloat(10.0)) == 1 { // deviation is greater than threshold 10%
// 			balResponse := &agents.AggregateBalanceResponse{
// 				ReqId:                reqEvent.ReqId,
// 				ActualBalance:        totalActualBalance,
// 				IdealBalance:         totalIdealBalance,
// 				ReferenceBlockNumber: reqEvent.ReferenceBlockNumber,
// 			}

// 			go r.SignAndSendBalanceResponse(ctx, balResponse)
// 		}

// 	}
// }

func (r *Operator) ProcessDeviationCheckRequest(ctx context.Context, eventCh <-chan *avs.Request) {
	for {

		reqEvent := <-eventCh

		//Now for every balance request event, there should be two jobs running in parallel
		//Routine 1: Fetch the balances of agents in all the chains to get the actual Balance
		//Routine 2: Fetch the total deposits & aprs of agents to get the ideal balance

		reqInfo, err := r.gRpcClient.CreateBalanceTaskToAggregator(reqEvent)

		if err != nil {
			fmt.Println(err)
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

func (r *Operator) SignAndSendBalanceResponse(ctx context.Context, balHeader *core.BalanceTaskInfo) {

	hash, err := balHeader.GetBalanceHeaderHash()

	r.logger.Info("BalancehEadr", "bal", balHeader)

	fmt.Println(hash)

	if err != nil {
		logger.GetLogger().Errorf("Error while hasing the message", "err", err)
		return
	}

	signature := r.blsKeyPair.SignMessage(hash)

	r.logger.Infof("Signed the message", "sign", signature)

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

func NewOperator(cfg config.Config) (*Operator, error) {
	var logLevel sdklogging.LogLevel

	if cfg.Production {
		logLevel = sdklogging.Production
	} else {
		logLevel = sdklogging.Development
	}

	logger, err := logger.NewZapLogger(logLevel)

	if err != nil {
		return nil, err
	}

	ethRpcClient, err := eth.NewClient(cfg.EthConfig.RpcUrl)

	if err != nil {
		panic(err)
	}

	//Operator Key Management

	operatorAddress, err := sdkEcdsa.GetAddressFromKeyStoreFile(cfg.KeyStore.Path)

	if err != nil {
		panic(err)
	}

	ecdsaKeypassword, ok := os.LookupEnv("OPERATOR_ECDSA_KEY_PASSWORD")

	if !ok {
		logger.Warnf("OPERATOR_ECDSA_KEY_PASSWORD env var not set")
		ecdsaKeypassword = cfg.KeyStore.Password
	}

	signerConfig := signerv2.Config{
		KeystorePath: cfg.KeyStore.Path,
		Password:     ecdsaKeypassword,
	}

	signerV2, _, err := signerv2.SignerFromConfig(signerConfig, big.NewInt(17000)) //Chain Id

	if err != nil {
		panic(err)
	}

	privateKey, err := sdkEcdsa.ReadKey(signerConfig.KeystorePath, signerConfig.Password)

	if err != nil {
		return nil, err
	}

	txSender, err := wallet.NewPrivateKeyWallet(ethRpcClient, signerV2, operatorAddress, logger)

	if err != nil {
		return nil, err
	}

	txMgr := txmgr.NewSimpleTxManager(txSender, ethRpcClient, logger, operatorAddress)

	avsClient := avs.NewAVSClient(&cfg, logger, txMgr)

	chainioConfig := clients.BuildAllConfig{
		EthHttpUrl:                 cfg.EthConfig.RpcUrl,
		RegistryCoordinatorAddr:    cfg.AVS.RegistryCoordinatorAddr.Hex(),
		OperatorStateRetrieverAddr: cfg.AVS.OperatorStateRetrieverAddr.Hex(),
		EthWsUrl:                   cfg.EthConfig.RpcUrl,
		AvsName:                    "CedroAVS",
		PromMetricsIpPortAddress:   "Cedro",
	}

	sdkClients, err := clients.BuildAll(chainioConfig, privateKey, logger)

	if err != nil {
		logger.Error("Error while bulding operator configs")

		fmt.Print(err)

		return nil, err
	}

	operatorId, err := sdkClients.AvsRegistryChainReader.GetOperatorId(&bind.CallOpts{}, operatorAddress)

	if err != nil {
		logger.Error("Cannot get operator id", "err", err)
		return nil, err
	}

	agentManager, err := avs.NewAgentManager(&cfg, logger, ethRpcClient)

	if err != nil {
		logger.Error("cannot get agent manager", err)
	}

	regAgents := make(map[string]common.IAgentPodClient)
	fmt.Println(RegisteredClients)

	for _, cfg := range cfg.Agents {
		regAgents[cfg.Name] = RegisteredClients[cfg.Name](cfg, logger)
	}

	gRpc, err := NewAggregatorGRpcClient(cfg, operatorId, operatorAddress, logger)

	if err != nil {
		logger.Errorf("Cannot initialize Aggregator Grpc client", "err", err)
		return nil, err
	}

	return &Operator{
		avs:                       avsClient,
		blsKeyPair:                cfg.BlsPrivateKey,
		keystore:                  &keys.KeyStore{},
		operatorId:                operatorId,
		txMgr:                     txMgr,
		logger:                    logger,
		gRPCAggreagatorIpPortAddr: cfg.GRpcAggregatorIpAddr,
		agentManager:              agentManager,
		chainAgents:               regAgents,
		gRpcClient:                gRpc,
	}, nil

}

func StartRelay(ctx context.Context, cfg *config.Config) {

	op, err := NewOperator(*cfg)

	if err != nil {
		panic("Couldn't start operator")
	}

	op.logger.Infof("Starting Operator")

	// First initialize Operator to the Aggregator

	op.logger.Infof("Init operator to aggregator.")
	err = op.gRpcClient.InitOperatorToAggregator()
	if err != nil {
		op.logger.Errorf("Init operator to aggregator failed: %v", err)
		panic(err)
	}

	op.logger.Infof("Init operator to aggregator succeed")

	balanceReqCh := make(chan *avs.Request) //Chanel to accept request from AVS contract
	// // deviationReqCh := make(chan *avs.Request) //Chanel to accept request from offline trigger

	var wg sync.WaitGroup

	wg.Add(1)

	go op.agentManager.PollEvents(ctx, balanceReqCh)

	go op.ProcessDeviationCheckRequest(ctx, balanceReqCh)

	wg.Wait()
}

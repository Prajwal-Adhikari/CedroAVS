package operator

import (
	"context"
	"math/big"
	"os"
	"sync"

	// "github.com/Cedro-Finance/CedroAVS/verifier/client/agents"
	"github.com/Cedro-Finance/CedroAVS/api/graphql"
	"github.com/Cedro-Finance/CedroAVS/core"
	"github.com/Cedro-Finance/CedroAVS/verifier/client/avs"
	"github.com/Cedro-Finance/CedroAVS/verifier/common"
	"github.com/Cedro-Finance/CedroAVS/verifier/config"
	"github.com/Cedro-Finance/CedroAVS/verifier/rpc"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/eth"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/wallet"
	"github.com/Layr-Labs/eigensdk-go/chainio/txmgr"
	sdkEcdsa "github.com/Layr-Labs/eigensdk-go/crypto/ecdsa"
	"github.com/Layr-Labs/eigensdk-go/signerv2"
	"github.com/Layr-Labs/eigensdk-go/types"
	sdktypes "github.com/Layr-Labs/eigensdk-go/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/Cedro-Finance/CedroAVS/core/logger"
	bls "github.com/Layr-Labs/eigensdk-go/crypto/bls"
	sdklogging "github.com/Layr-Labs/eigensdk-go/logging"
)

var RegisteredClients = map[string]func(config.ChainConfig, sdklogging.Logger) common.IAgentPodClient{}
var maxRetry = 5

type OperatorDeviationTask struct {
	req     *avs.Request
	reqInfo *core.BalanceTaskInfo
	retries uint64
}

type Operator struct {
	avs        *avs.AVSClient
	blsKeyPair *bls.KeyPair

	agentManager              *avs.AgentManager
	operatorId                sdktypes.OperatorId
	txMgr                     txmgr.TxManager
	logger                    sdklogging.Logger
	gRPCAggreagatorIpPortAddr string
	gRpcClient                *rpc.AggregatorGRpcClient
	chainAgents               map[string]common.IAgentPodClient
	graphClient               *graphql.GraphQLClient

	balanceTasks     map[sdktypes.TaskIndex](*OperatorDeviationTask)
	balanceTaskMutex sync.RWMutex

	balanceReqCh chan *avs.Request
	errCh        chan *DeviationError

	signatureCh chan *core.SignedTaskRespRequest
}

func (o *Operator) GetTaskByIndex(taskIndex sdktypes.TaskIndex) *OperatorDeviationTask {
	o.balanceTaskMutex.RLock()
	defer o.balanceTaskMutex.RUnlock()

	res := o.balanceTasks[taskIndex]

	return res
}

func (o *Operator) CreateOperatorTask(taskInfo *core.BalanceTaskInfo, req *avs.Request) (*OperatorDeviationTask, error) {
	task := o.GetTaskByIndex(taskInfo.TaskIndex)

	if task == nil {
		o.balanceTaskMutex.Lock()
		task = &OperatorDeviationTask{
			req:     req,
			retries: 0,
			reqInfo: taskInfo,
		}
		o.balanceTasks[taskInfo.TaskIndex] = task
		o.balanceTaskMutex.Unlock()
	}

	return task, nil
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
		return nil, err
	}

	//Operator Key Management
	operatorAddress, err := sdkEcdsa.GetAddressFromKeyStoreFile(cfg.KeyStore.Path)

	if err != nil {
		return nil, err
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
		return nil, err
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
		return nil, err
	}

	regAgents := make(map[string]common.IAgentPodClient)

	for _, cfg := range cfg.Agents {
		if RegisteredClients[cfg.ChainId] != nil {
			regAgents[cfg.ChainId] = RegisteredClients[cfg.ChainId](cfg, logger)
		}
	}

	gRpc, err := rpc.NewAggregatorGRpcClient(cfg, operatorId, operatorAddress, logger)

	if err != nil {
		logger.Errorf("Cannot initialize Aggregator Grpc client", "err", err)
		return nil, err
	}

	graphQLClient := graphql.NewGraphQLClient(cfg.GraphQLUrl, logger)

	return &Operator{
		avs:                       avsClient,
		blsKeyPair:                cfg.BlsPrivateKey,
		operatorId:                operatorId,
		txMgr:                     txMgr,
		logger:                    logger,
		gRPCAggreagatorIpPortAddr: cfg.GRpcAggregatorIpAddr,
		agentManager:              agentManager,
		chainAgents:               regAgents,
		gRpcClient:                gRpc,
		graphClient:               graphQLClient,
		balanceTasks:              make(map[uint32]*OperatorDeviationTask),

		balanceReqCh: make(chan *avs.Request),
		errCh:        make(chan *DeviationError),
		signatureCh:  make(chan *core.SignedTaskRespRequest, 10),
	}, nil

}

type DeviationError struct {
	TaskIndex types.TaskIndex
	Error     error
}

func Start(ctx context.Context, cfg *config.Config) {

	op, err := NewOperator(*cfg)

	if err != nil {
		op.logger.Error("Error while starting operator.", "err", err)
		return
	}

	op.logger.Info("Starting Operator", "operatorId", op.operatorId)

	// First initialize Operator to the Aggregator

	op.logger.Infof("Initializing operator to aggregator.")
	err = op.gRpcClient.InitOperatorToAggregator()
	if err != nil {
		op.logger.Errorf("Init operator to aggregator failed: %v", err)
		return
	}

	op.logger.Infof("Init operator to aggregator succeed")

	go op.agentManager.PollEvents(ctx, op.balanceReqCh, op.avs.AvsReader.AvsServiceBindings.ServiceManager)

	for {
		select {
		case <-ctx.Done():
			return
		case reqEvent := <-op.balanceReqCh:
			op.logger.Debug("Event Received", "req", reqEvent.ReqId, "agent", reqEvent.AgentAddress.Hex())

			// RETRY_LOOP:
			// 	for {
			reqInfo, err := op.gRpcClient.CreateBalanceTaskToAggregator(reqEvent)

			if err != nil {
				op.logger.Error("Create Task To Aggregator Failed. Trying again", "err", err)
				continue

				// break RETRY_LOOP
			}

			task, err := op.CreateOperatorTask(reqInfo, reqEvent)

			if err != nil {
				op.logger.Error("Create task in the ooerator side threw isseu")
				continue
			}

			go op.ProcessDeviationCheckRequest(ctx, task)
			// }

		case err := <-op.errCh:
			op.logger.Error("Error occured while performing the task", err.Error)

			task := op.GetTaskByIndex(err.TaskIndex)

			//TODO: Remove Hardcode value
			if task.retries < uint64(maxRetry) {
				op.logger.Info("Retrying the request again", "reqId", task.req.ReqId, "agent", task.req.AgentAddress.Hex())

				task.retries += 1

				op.balanceReqCh <- task.req
			}

		case signReq := <-op.signatureCh:
			op.logger.Info("Signature Receive")

			go op.SendSignedTaskToAggregator(ctx, signReq)

		}
	}
}

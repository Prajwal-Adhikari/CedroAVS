package operator

import (
	"context"
	"math/big"
	"os"
	"sync"

	// "github.com/Cedro-Finance/CedroAVS/verifier/client/agents"
	"github.com/Cedro-Finance/CedroAVS/api/graphql"
	avsclient "github.com/Cedro-Finance/CedroAVS/verifier/client/avs"
	"github.com/Cedro-Finance/CedroAVS/verifier/common"
	"github.com/Cedro-Finance/CedroAVS/verifier/opio"
	"github.com/Cedro-Finance/CedroAVS/verifier/rpc"
	"github.com/Cedro-Finance/CedroAVS/verifier/tasks"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/eth"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/wallet"
	"github.com/Layr-Labs/eigensdk-go/chainio/txmgr"
	sdkEcdsa "github.com/Layr-Labs/eigensdk-go/crypto/ecdsa"
	"github.com/Layr-Labs/eigensdk-go/signerv2"
	sdktypes "github.com/Layr-Labs/eigensdk-go/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	coreavs "github.com/Cedro-Finance/CedroAVS/core/avsClient"
	"github.com/Cedro-Finance/CedroAVS/core/config"
	"github.com/Cedro-Finance/CedroAVS/core/logger"
	bls "github.com/Layr-Labs/eigensdk-go/crypto/bls"
	sdklogging "github.com/Layr-Labs/eigensdk-go/logging"
)

var RegisteredClients = map[string]func(config.ChainConfig, sdklogging.Logger, coreavs.IAVSReader) common.IAgentPodClient{}

type Operator struct {
	avs        *avsclient.AVSClient
	blsKeyPair *bls.KeyPair

	operatorId                sdktypes.OperatorId
	txMgr                     txmgr.TxManager
	logger                    sdklogging.Logger
	gRPCAggreagatorIpPortAddr string
	gRpcClient                *rpc.AggregatorGRpcClient
	chainAgents               map[string]common.IAgentPodClient
	graphClient               *graphql.GraphQLClient

	deviateProcessor *tasks.BalanceTaskProcessor
	slashProcessor   *tasks.SlashTaskProcessor
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

	avsClient := avsclient.NewAVSClient(&cfg, logger, txMgr)

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

	regAgents := make(map[string]common.IAgentPodClient)

	for _, cfg := range cfg.Agents {
		if RegisteredClients[cfg.ChainId] != nil {
			regAgents[cfg.ChainId] = RegisteredClients[cfg.ChainId](cfg, logger, avsClient.AvsReader)
		}
	}

	gRpc, err := rpc.NewAggregatorGRpcClient(cfg, operatorId, operatorAddress, logger)

	if err != nil {
		logger.Errorf("Cannot initialize Aggregator Grpc client", "err", err)
		return nil, err
	}

	graphQLClient := graphql.NewGraphQLClient(cfg.GraphQLUrl, logger)

	opIoClient := opio.NewOpeartorIO(txMgr, cfg.BlsPrivateKey, gRpc, graphQLClient)

	deviationProcessor := tasks.NewBalanceTask(opIoClient, regAgents, logger)

	slashProcessor := tasks.NewSlashTask(opIoClient, logger)

	return &Operator{
		avs:                       avsClient,
		blsKeyPair:                cfg.BlsPrivateKey,
		operatorId:                operatorId,
		txMgr:                     txMgr,
		logger:                    logger,
		gRPCAggreagatorIpPortAddr: cfg.GRpcAggregatorIpAddr,

		chainAgents: regAgents,
		gRpcClient:  gRpc,
		graphClient: graphQLClient,

		deviateProcessor: deviationProcessor,
		slashProcessor:   slashProcessor,
	}, nil

}

type DeviationError struct {
	TaskIndex sdktypes.TaskIndex
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

	//Listen Ethereum Agent manager event
	for _, client := range op.chainAgents {
		go client.PollEvents(ctx, op.deviateProcessor.GetRequestCh(), op.slashProcessor.GetRequestCh())
	}

	//Start Task Processors

	var wg sync.WaitGroup

	wg.Add(2)

	go op.deviateProcessor.Start(ctx, &wg)
	go op.slashProcessor.Start(ctx, &wg)

	wg.Wait()
}

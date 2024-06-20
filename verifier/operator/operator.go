package operator

import (
	"context"
	"math/big"
	"os"
	"sync"

	// "github.com/Cedro-Finance/CedroAVS/verifier/client/agents"
	"github.com/Cedro-Finance/CedroAVS/api/graphql"
	"github.com/Cedro-Finance/CedroAVS/verifier/client/avs"
	"github.com/Cedro-Finance/CedroAVS/verifier/common"
	"github.com/Cedro-Finance/CedroAVS/verifier/config"
	"github.com/Cedro-Finance/CedroAVS/verifier/keys"
	"github.com/Cedro-Finance/CedroAVS/verifier/rpc"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/eth"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/wallet"
	"github.com/Layr-Labs/eigensdk-go/chainio/txmgr"
	sdkEcdsa "github.com/Layr-Labs/eigensdk-go/crypto/ecdsa"
	"github.com/Layr-Labs/eigensdk-go/signerv2"
	sdktypes "github.com/Layr-Labs/eigensdk-go/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/Cedro-Finance/CedroAVS/core/logger"
	bls "github.com/Layr-Labs/eigensdk-go/crypto/bls"
	sdklogging "github.com/Layr-Labs/eigensdk-go/logging"
)

var RegisteredClients = map[string]func(config.ChainConfig, sdklogging.Logger) common.IAgentPodClient{}

type Operator struct {
	avs        *avs.AVSClient
	keystore   *keys.KeyStore
	blsKeyPair *bls.KeyPair

	agentManager              *avs.AgentManager
	operatorId                sdktypes.OperatorId
	txMgr                     txmgr.TxManager
	logger                    sdklogging.Logger
	gRPCAggreagatorIpPortAddr string
	gRpcClient                *rpc.AggregatorGRpcClient
	chainAgents               map[string]common.IAgentPodClient

	graphClient *graphql.GraphQLClient
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
		keystore:                  &keys.KeyStore{},
		operatorId:                operatorId,
		txMgr:                     txMgr,
		logger:                    logger,
		gRPCAggreagatorIpPortAddr: cfg.GRpcAggregatorIpAddr,
		agentManager:              agentManager,
		chainAgents:               regAgents,
		gRpcClient:                gRpc,
		graphClient:               graphQLClient,
	}, nil

}

func Start(ctx context.Context, cfg *config.Config) {

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

	var wg sync.WaitGroup

	wg.Add(1)

	go op.agentManager.PollEvents(ctx, balanceReqCh, op.avs.AvsReader.AvsServiceBindings.ServiceManager)
	go op.ProcessDeviationCheckRequest(ctx, balanceReqCh)

	wg.Wait()
}

package aggregator

import (
	"context"
	"math/big"
	"os"
	"sync"

	"github.com/Cedro-Finance/CedroAVS/aggregator/rpc"
	"github.com/Cedro-Finance/CedroAVS/aggregator/taskprocessor"
	avs "github.com/Cedro-Finance/CedroAVS/core/avsClient"
	"github.com/Cedro-Finance/CedroAVS/core/tasks"

	// cedroservicemanager "github.com/Cedro-Finance/CedroAVS/core/bindings/CedroServiceManager"
	"github.com/Cedro-Finance/CedroAVS/core/config"
	"github.com/Cedro-Finance/CedroAVS/core/logger"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/eth"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/wallet"
	"github.com/Layr-Labs/eigensdk-go/chainio/txmgr"

	// "github.com/Layr-Labs/eigensdk-go/crypto/bls"
	sdkEcdsa "github.com/Layr-Labs/eigensdk-go/crypto/ecdsa"
	sdklogging "github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/Layr-Labs/eigensdk-go/signerv2"
)

type Aggregator struct {
	logger               sdklogging.Logger
	grpcServerIpPortAddr string
	gRpc                 *rpc.GRpcHandler
	avsWriter            avs.IAvsWriter

	taskprocessors map[string]*taskprocessor.TaskProcessor
}

func NewAggregator(cfg *config.Config) (*Aggregator, error) {
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

	// privateKey, err := sdkEcdsa.ReadKey(signerConfig.KeystorePath, signerConfig.Password)

	// if err != nil {
	// 	return nil, err
	// }

	txSender, err := wallet.NewPrivateKeyWallet(ethRpcClient, signerV2, operatorAddress, logger)

	if err != nil {
		return nil, err
	}

	txMgr := txmgr.NewSimpleTxManager(txSender, ethRpcClient, logger, operatorAddress)

	avsWriter, err := avs.NewAvsWriter(cfg, txMgr, ethRpcClient, logger)

	if err != nil {
		logger.Error("Error while setting avs writer", "err", err)
		return nil, err
	}

	tps := make(map[string]*taskprocessor.TaskProcessor)

	for _, tt := range tasks.TaskTypes {
		tp, err := taskprocessor.NewTaskProcessorService(tt, cfg, logger, avsWriter)

		if err != nil {
			panic(err)
		}

		tps[tt] = tp
	}

	aggGrpcHandler, err := rpc.NewAggregatorRpcHandler(logger, tps)

	if err != nil {
		logger.Errorf("Error while creating aggregator service", "err", err)
		return nil, err
	}

	grpc := rpc.NewGRpcHandler(logger, aggGrpcHandler)

	return &Aggregator{
		logger:               logger,
		gRpc:                 grpc,
		grpcServerIpPortAddr: cfg.GRpcAggregatorIpAddr,
		avsWriter:            avsWriter,
		taskprocessors:       tps,
	}, nil
}

func (a *Aggregator) Start(ctx context.Context) error {

	var wg sync.WaitGroup

	a.logger.Infof("Aggregator is starting")

	wg.Add(1)

	go func() {
		defer wg.Done()
		a.gRpc.Start(ctx, a.grpcServerIpPortAddr)
	}()

	wg.Add(len(a.taskprocessors))

	for _, tp := range a.taskprocessors {
		go func(tp *taskprocessor.TaskProcessor) {
			defer wg.Done()
			tp.Start(ctx)
		}(tp)
	}

	wg.Wait()

	return nil
}

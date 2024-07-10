package aggregator

import (
	"context"
	"math/big"
	"os"
	"sync"

	"github.com/Cedro-Finance/CedroAVS/aggregator/rpc"
	"github.com/Cedro-Finance/CedroAVS/aggregator/taskprocessor"
	avs "github.com/Cedro-Finance/CedroAVS/core/avsClient"

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

	avsWriter *avs.AvsWriter
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

	deviator, err := taskprocessor.NewDeviationProcessor(cfg, logger, avsWriter)

	if err != nil {
		logger.Errorf("Error while creating aggregator service", "err", err)
		return nil, err
	}

	slasher, err := taskprocessor.NewSlashProcessor(cfg, logger, avsWriter)

	if err != nil {
		logger.Errorf("Error while creating aggregator service", "err", err)
		return nil, err
	}

	aggGrpcHandler, err := rpc.NewAggregatorRpcHandler(logger, deviator, slasher)

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
	}, nil
}

func (a *Aggregator) Start(ctx context.Context, wg *sync.WaitGroup) error {
	defer func() {
		wg.Done()
	}()

	a.logger.Infof("Aggregator is starting")

	go a.gRpc.Start(ctx, a.grpcServerIpPortAddr)

	for {
		select {
		case <-ctx.Done():
			return nil

			// case blsAggServiceResp := <-a.service.GetResponseChannel():
			// 	a.logger.Info("Received response from blsAggregationService", "blsAggServiceResp", blsAggServiceResp)
			// a.sendAggregatedResponseToContract(blsAggServiceResp)
		}
	}

}

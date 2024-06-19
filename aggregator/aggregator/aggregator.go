package aggregator

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"sync"

	"github.com/Cedro-Finance/CedroAVS/aggregator/rpc"
	"github.com/Cedro-Finance/CedroAVS/core"
	avs "github.com/Cedro-Finance/CedroAVS/core/avsClient"
	"github.com/Cedro-Finance/CedroAVS/core/config"
	"github.com/Cedro-Finance/CedroAVS/core/logger"
	abi "github.com/Cedro-Finance/CedroAVS/verifier/client/avs/abi"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/eth"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/wallet"
	"github.com/Layr-Labs/eigensdk-go/chainio/txmgr"
	sdkEcdsa "github.com/Layr-Labs/eigensdk-go/crypto/ecdsa"
	sdklogging "github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/Layr-Labs/eigensdk-go/signerv2"
	sdktypes "github.com/Layr-Labs/eigensdk-go/types"
	"github.com/ethereum/go-ethereum/common"
)

type OperatorStatus struct {
	LastTime   int64               `json:"lastTime"`
	OperatorId sdktypes.OperatorId `json:"operatorId"`
}

type FinishedTaskStatus struct {
	Message          *core.BalanceTaskInfo
	TxHash           common.Hash
	BlockHash        common.Hash `json:"blockHash,omitempty"`
	BlockNumber      *big.Int    `json:"blockNumber,omitempty"`
	TransactionIndex uint        `json:"transactionIndex"`
}

type Aggregator struct {
	logger               sdklogging.Logger
	grpcServerIpPortAddr string
	gRpc                 *rpc.GRpcHandler
	service              *AggregatorService

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

	service, err := NewAggregatorService(cfg, logger)

	if err != nil {
		logger.Errorf("Error while creating aggregator service", "err", err)
		return nil, err
	}

	grpc := rpc.NewGRpcHandler(logger, service)

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

	return &Aggregator{
		logger:               logger,
		gRpc:                 grpc,
		grpcServerIpPortAddr: cfg.GRpcAggregatorIpAddr,
		service:              service,
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

		case blsAggServiceResp := <-a.service.GetResponseChannel():
			a.logger.Info("Received response from blsAggregationService", "blsAggServiceResp", blsAggServiceResp)
			a.sendAggregatedResponseToContract(blsAggServiceResp)
		}
	}

}

func (agg *Aggregator) sendAggregatedResponseToContract(blsAggServiceResp BLSAggregationServiceResponse) {
	// TODO: check if blsAggServiceResp contains an err
	if blsAggServiceResp.Err != nil {
		agg.logger.Error("BlsAggregationServiceResponse contains an error", "err", blsAggServiceResp.Err)
		// panicing to help with debugging (fail fast), but we shouldn't panic if we run this in production
		panic(blsAggServiceResp.Err)
	}
	nonSignerPubkeys := []abi.BN254G1Point{}
	for _, nonSignerPubkey := range blsAggServiceResp.NonSignersPubkeysG1 {
		nonSignerPubkeys = append(nonSignerPubkeys, core.ConvertToBN254G1Point(nonSignerPubkey))
	}
	quorumApks := []abi.BN254G1Point{}
	for _, quorumApk := range blsAggServiceResp.QuorumApksG1 {
		quorumApks = append(quorumApks, core.ConvertToBN254G1Point(quorumApk))
	}
	nonSignerStakesAndSignature := abi.IBLSSignatureCheckerNonSignerStakesAndSignature{
		NonSignerPubkeys:             nonSignerPubkeys,
		QuorumApks:                   quorumApks,
		ApkG2:                        core.ConvertToBN254G2Point(blsAggServiceResp.SignersApkG2),
		Sigma:                        core.ConvertToBN254G1Point(blsAggServiceResp.SignersAggSigG1.G1Point),
		NonSignerQuorumBitmapIndices: blsAggServiceResp.NonSignerQuorumBitmapIndices,
		QuorumApkIndices:             blsAggServiceResp.QuorumApkIndices,
		TotalStakeIndices:            blsAggServiceResp.TotalStakeIndices,
		NonSignerStakeIndices:        blsAggServiceResp.NonSignerStakeIndices,
	}

	agg.logger.Info("Threshold reached. Sending aggregated response onchain.",
		"taskIndex", blsAggServiceResp.TaskIndex,
	)

	fmt.Println(nonSignerStakesAndSignature)

	task := agg.service.GetTaskByIndex(blsAggServiceResp.TaskIndex)

	res, err := agg.avsWriter.SendConfirmationBatch(context.Background(), task, nonSignerStakesAndSignature)
	if err != nil {
		agg.logger.Error("Aggregator failed to respond to task", "err", err)
	}

	fmt.Println("RESPONSE", res)
}

package taskprocessor

import (
	// "context"

	"context"
	"fmt"

	blsaggregator "github.com/Cedro-Finance/CedroAVS/aggregator/bls"
	"github.com/Cedro-Finance/CedroAVS/api/graphql"
	"github.com/Cedro-Finance/CedroAVS/core"
	avs "github.com/Cedro-Finance/CedroAVS/core/avsClient"
	cedroservicemanager "github.com/Cedro-Finance/CedroAVS/core/bindings/CedroServiceManager"
	"github.com/Cedro-Finance/CedroAVS/core/config"
	sdkclients "github.com/Layr-Labs/eigensdk-go/chainio/clients"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/eth"
	"github.com/Layr-Labs/eigensdk-go/crypto/bls"
	sdkEcdsa "github.com/Layr-Labs/eigensdk-go/crypto/ecdsa"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/Layr-Labs/eigensdk-go/types"
)

type AggregatorService struct {
	logger                logging.Logger
	ethClient             eth.Client
	avsReader             *avs.AvsReader
	BlsAggregationService *blsaggregator.BLSAggregatorService

	cfg *config.Config

	graphClient *graphql.GraphQLClient
}

func NewAggregatorService(c *config.Config, logger logging.Logger) (*AggregatorService, error) {

	ethClient, err := eth.NewClient(c.EthConfig.RpcUrl)

	if err != nil {
		logger.Errorf("Error while creating Eth client", "err", err)
		return nil, err
	}

	graphClient := graphql.NewGraphQLClient(c.GraphQLUrl, logger)

	avsReader, err := avs.NewAvsReader(c.AVS, ethClient, logger, graphClient)

	if err != nil {
		logger.Errorf("Error while creating Avs reader", "err", err)
	}

	privateKey, err := sdkEcdsa.ReadKey(c.KeyStore.Path, c.KeyStore.Password)

	if err != nil {
		return nil, err
	}

	chainioConfig := sdkclients.BuildAllConfig{
		EthHttpUrl:                 c.EthConfig.RpcUrl,
		EthWsUrl:                   c.EthConfig.WsUrl,
		RegistryCoordinatorAddr:    c.AVS.RegistryCoordinatorAddr.String(),
		OperatorStateRetrieverAddr: c.AVS.OperatorStateRetrieverAddr.String(),
		AvsName:                    "CedroAVS",
		PromMetricsIpPortAddress:   "Cedro",
	}

	clients, err := sdkclients.BuildAll(chainioConfig, privateKey, logger)
	if err != nil {
		logger.Errorf("Cannot create sdk clients", "err", err)
		return nil, err
	}

	blsAggregationService := blsaggregator.NewBlsAggregatorService(avsReader, logger)

	return &AggregatorService{
		logger:                logger,
		avsReader:             avsReader,
		ethClient:             clients.EthHttpClient,
		BlsAggregationService: blsAggregationService,
		cfg:                   c,
		graphClient:           graphClient,
	}, nil

}

func (agg *AggregatorService) GetResponseChannel() <-chan blsaggregator.BLSAggregationServiceResponse {
	return agg.BlsAggregationService.GetResponseChannel()
}

type QuorumDetails struct {
	QuorumNums                 types.QuorumNums
	QuorumThresholdPercentages types.QuorumThresholdPercentages
}

func (agg *AggregatorService) GetQuorumDetails(ctx context.Context, refBlock uint64) (*QuorumDetails, error) {

	var err error

	agg.logger.Info("get from Contract", "referenceBlockNumber", refBlock)

	quorumNumbers, err := agg.avsReader.GetQuorumsByBlockNumber(context.Background(), uint32(refBlock))
	if err != nil {
		agg.logger.Error("GetQuorumCountByBlockNumber failed", "err", err)
		return nil, err
	}
	agg.logger.Info("get quorumNumbers ", "quorumNumbers", fmt.Sprintf("%v", quorumNumbers))

	quorumThresholdPercentages, err := agg.avsReader.GetQuorumThresholdPercentages(context.Background(), uint32(refBlock), quorumNumbers)
	if err != nil {
		agg.logger.Error("GetQuorumThresholdPercentages failed", "err", err)
		return nil, err
	}

	agg.logger.Info(
		"quorum datas",
		"numbers", fmt.Sprintf("%v", quorumNumbers.UnderlyingType()),
		"thresholdPercentages", fmt.Sprintf("%v", quorumThresholdPercentages.UnderlyingType()),
	)

	return &QuorumDetails{
		QuorumNums:                 quorumNumbers,
		QuorumThresholdPercentages: quorumThresholdPercentages,
	}, nil
}

func (agg *AggregatorService) ReceiveAggregatedSignatures(blsAggServiceResp blsaggregator.BLSAggregationServiceResponse) (cedroservicemanager.IBLSSignatureCheckerNonSignerStakesAndSignature, error) {
	// TODO: check if blsAggServiceResp contains an err
	if blsAggServiceResp.Err != nil {
		agg.logger.Error("BlsAggregationServiceResponse contains an error", "err", blsAggServiceResp.Err)
		// panicing to help with debugging (fail fast), but we shouldn't panic if we run this in production
		return cedroservicemanager.IBLSSignatureCheckerNonSignerStakesAndSignature{}, blsAggServiceResp.Err
	}
	nonSignerPubkeys := []cedroservicemanager.BN254G1Point{}
	for _, nonSignerPubkey := range blsAggServiceResp.NonSignersPubkeysG1 {
		nonSignerPubkeys = append(nonSignerPubkeys, core.ConvertToBN254G1Point(nonSignerPubkey))
	}
	quorumApks := []cedroservicemanager.BN254G1Point{}
	for _, quorumApk := range blsAggServiceResp.QuorumApksG1 {
		quorumApks = append(quorumApks, core.ConvertToBN254G1Point(quorumApk))
	}

	var aggSig *bls.Signature
	for _, sign := range blsAggServiceResp.SignersAggSigG1 {
		if aggSig == nil {
			aggSig = bls.NewZeroSignature().Add(sign)
		} else {
			aggSig.Add(sign)
		}
	}

	// Aggregate the aggregated public keys. We reuse the first aggregated public key as the accumulator
	var aggPubKey *bls.G2Point
	for _, pKey := range blsAggServiceResp.SignersApkG2 {
		if aggPubKey == nil {
			aggPubKey = bls.NewZeroG2Point().Add(pKey)
		} else {
			aggPubKey.Add(pKey)
		}
	}

	nonSignerStakesAndSignature := cedroservicemanager.IBLSSignatureCheckerNonSignerStakesAndSignature{
		NonSignerPubkeys:             nonSignerPubkeys,
		QuorumApks:                   quorumApks,
		ApkG2:                        core.ConvertToBN254G2Point(aggPubKey),
		Sigma:                        core.ConvertToBN254G1Point(aggSig.G1Point),
		NonSignerQuorumBitmapIndices: blsAggServiceResp.NonSignerQuorumBitmapIndices,
		QuorumApkIndices:             blsAggServiceResp.QuorumApkIndices,
		TotalStakeIndices:            blsAggServiceResp.TotalStakeIndices,
		NonSignerStakeIndices:        blsAggServiceResp.NonSignerStakeIndices,
	}

	return nonSignerStakesAndSignature, nil
}

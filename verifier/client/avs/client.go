package avs

import (
	"context"
	"time"

	// "fmt"
	"math/big"

	// "github.com/Cedro-Finance/CedroAVS/core/logger"
	"github.com/Cedro-Finance/CedroAVS/verifier/config"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/eth"
	"github.com/Layr-Labs/eigensdk-go/chainio/txmgr"
	"github.com/Layr-Labs/eigensdk-go/logging"
	ether "github.com/ethereum/go-ethereum"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const (
	defaultBlockTime = time.Duration(12) * time.Second

	defaultStartBlock                    = 1705750
	defaultHeightDifferenceForFilterLogs = 100
	defaultRequestId                     = 1

	defaultPollDuration = 5 * time.Second
)

func (avs *AVSClient) GetCurrentBlock(ctx context.Context) (uint64, error) {
	currentBlock, err := avs.Eth.BlockNumber(ctx)

	if err != nil {
		return 0, err
	}

	return currentBlock, nil
}

func (avs *AVSClient) GetFinalizedBlock(ctx context.Context) (uint64, error) {
	currentBlock, err := avs.GetCurrentBlock(ctx)

	if err != nil {
		avs.Logger.Error("Error while getting block information")

		return 0, err
	}

	return currentBlock - avs.BlockFinality, nil

}

func (avs *AVSClient) FilterLogs(
	ctx context.Context, fromHeight uint64, toHeight uint64,
	contractAddress ethCommon.Address, topics ethCommon.Hash) ([]types.Log, error) {

	logs, err := avs.Eth.FilterLogs(ctx, ether.FilterQuery{
		FromBlock: big.NewInt(int64(fromHeight)),
		ToBlock:   big.NewInt(int64(toHeight)),
		Addresses: []ethCommon.Address{contractAddress},
		Topics:    [][]ethCommon.Hash{{topics}},
	})
	if err != nil {
		return nil, err
	}
	return logs, err
}

type AVSClient struct {
	Eth              eth.Client
	StartBlock       uint64
	NextBlockHeight  uint64
	PollDuration     time.Duration
	CurrentRequestId uint64

	BlockTime     time.Duration
	BlockFinality uint64

	AvsReader *AvsReader
	AvsWriter *AvsWriter

	Logger logging.Logger
}

func NewAVSClient(cfg *config.Config, logger logging.Logger, txMgr txmgr.TxManager) *AVSClient {

	ethClient, err := eth.NewClient(cfg.EthConfig.RpcUrl)

	if err != nil {
		panic(err)
	}

	avsReader, err := NewAvsReader(&cfg.AVS, ethClient, logger)

	if err != nil {
		logger.Error("Error while Createing AVS Reader")
	}

	avsWriter, err := NewAvsWriter(&cfg.AVS, txMgr, ethClient, logger)

	if err != nil {
		logger.Error("Error while creating AVS Writer")
	}

	return &AVSClient{
		Eth:       ethClient,
		AvsReader: avsReader,
		AvsWriter: avsWriter,
		Logger:    logger,
	}
}

type Request struct {
	ReqId                *big.Int
	ReferenceBlockNumber uint64
	AgentAddress         ethCommon.Address
	Timestamp            uint64
}

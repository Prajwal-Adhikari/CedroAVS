package avs

import (
	"time"

	// "fmt"
	"math/big"

	// "github.com/Cedro-Finance/CedroAVS/core/logger"
	"github.com/Cedro-Finance/CedroAVS/verifier/config"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/eth"
	"github.com/Layr-Labs/eigensdk-go/chainio/txmgr"
	"github.com/Layr-Labs/eigensdk-go/logging"
	ethCommon "github.com/ethereum/go-ethereum/common"
)

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
		logger.Error("Error while Createing AVS Reader", err)
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

type TaskRequest interface {
	GetReqHash() string
}

type DeviationRequest struct {
	ReqId                *big.Int
	ReferenceBlockNumber uint64
	AgentAddress         ethCommon.Address
	Timestamp            uint64
}

func (d *DeviationRequest) GetReqHash() string {
	return d.AgentAddress.Hex()
}

type AgentSlashRequest struct {
	TxnId                *big.Int
	ReferenceBlockNumber uint64
	AgentAddress         ethCommon.Address
	InsufficientBalance  *big.Int
}

func (a *AgentSlashRequest) GetReqHash() string {
	return a.AgentAddress.Hex()
}

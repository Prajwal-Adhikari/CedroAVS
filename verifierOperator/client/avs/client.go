package avs

import (
	"context"
	"time"

	// "fmt"
	"math/big"

	// "github.com/Cedro-Finance/CedroAVS/core/logger"
	"github.com/Cedro-Finance/CedroAVS/verifierOperator/config"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/eth"
	"github.com/Layr-Labs/eigensdk-go/chainio/txmgr"
	"github.com/Layr-Labs/eigensdk-go/logging"
	ether "github.com/ethereum/go-ethereum"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const (
	defaultBlockTime = time.Duration(12) * time.Second

	defaultStartBlock                    = 1234
	defaultHeightDifferenceForFilterLogs = 100
	defaultRequestId                     = 1

	defaultPollDuration = defaultHeightDifferenceForFilterLogs * 2 / 3 * time.Second
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

// func (avs *AVSClient) NewAVSContract() *abi.AVS {
// 	contract, err := abi.NewAVS(avs.contract, avs.eth)

// 	if err != nil {
// 		logger.GetLogger().Error("Error while initializing contract")
// 		panic(err)
// 	}

// 	return contract
// }

// func (avs *AVSClient) GetRequestId() (*big.Int, error) {
// 	contract := avs.NewAVSContract()
// 	reqId, err := contract.ReqId(nil)

// 	if err != nil {
// 		return big.NewInt(0), err
// 	}

// 	return reqId, nil
// }

// func (avs *AVSClient) ParseGetBalanceLog(log types.Log) (*abi.AVSGetChainTVL, error) {
// 	parsed, err := avs.NewAVSContract().ParseGetChainTVL(log)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return parsed, nil
// }

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

func NewAVSClient(cfg *config.AVSConfig, logger logging.Logger, txMgr txmgr.TxManager) *AVSClient {

	ethClient, err := eth.NewClient(cfg.EthRpcUrl)

	if err != nil {
		panic(err)
	}

	avsReader, err := NewAvsReader(cfg, ethClient, logger)

	if err != nil {
		logger.Error("Error while Createing AVS Reader")
	}

	avsWriter, err := NewAvsWriter(cfg, txMgr, ethClient, logger)

	if err != nil {
		logger.Error("Error while creating AVS Writer")
	}

	return &AVSClient{
		Eth:              ethClient,
		StartBlock:       cfg.StartHeight,
		NextBlockHeight:  cfg.StartHeight,
		PollDuration:     time.Duration(time.Minute),
		CurrentRequestId: 0,
		BlockTime:        defaultBlockTime,
		BlockFinality:    cfg.BlockFinality,

		AvsReader: avsReader,
		AvsWriter: avsWriter,
		Logger:    logger,
	}
}

type Request struct {
	ReqId        *big.Int
	Height       *big.Int
	AgentAddress string
}

// polls the balance request events from the AVS contract on eth
func (c *AVSClient) PollEvents(ctx context.Context, ch chan<- *Request) {

	dur := c.PollDuration

	if dur == 0 {
		dur = defaultPollDuration
	}

	ticker := time.NewTicker(dur)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}

		// L1:
		// 	for {
		// 		finalizedBlock, err := c.GetFinalizedBlock(ctx)

		// 		if err != nil {
		// 			//handle error properly
		// 			break L1
		// 		}

		// 		if finalizedBlock < c.nextBlockHeight {
		// 			diff := c.nextBlockHeight - finalizedBlock
		// 			time.Sleep((time.Duration(diff) * c.blockTime))
		// 		}

		// 		currentReqId, err := c.GetRequestId()

		// 		if err != nil {
		// 			break L1
		// 		}

		// 		for startHeight := c.nextBlockHeight; startHeight <= finalizedBlock; startHeight += defaultHeightDifferenceForFilterLogs + 1 {
		// 			endHeight := startHeight + defaultHeightDifferenceForFilterLogs

		// 			if endHeight > finalizedBlock {
		// 				endHeight = finalizedBlock
		// 			}

		// 			logs, err := c.FilterLogs(ctx, startHeight, endHeight, c.RegistryCoordinatorAddr, ethCommon.Hash{}) //TODO: get topics from config

		// 			if err != nil {
		// 				//Implement Logger
		// 				logger.GetLogger().Error("")
		// 				break L1
		// 			}

		// 			for _, log := range logs {
		// 				parsed, err := c.ParseGetBalanceLog(log)

		// 				if err != nil {
		// 					break L1
		// 				}

		// 				if currentReqId == parsed.ReqId {
		// 					ch <- &Request{
		// 						ReqId:  parsed.ReqId,
		// 						Height: parsed.ReqId,
		// 					}

		// 					//TODO: Will there be condition when reqId 2 arrives before settling reqId 1?
		// 					c.currentRequestId += 1
		// 				}
		// 			}

		// 			c.nextBlockHeight = endHeight + 1
		// 		}
		// 	}

	}
}

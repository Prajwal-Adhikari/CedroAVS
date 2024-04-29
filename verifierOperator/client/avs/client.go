package avs

import (
	"context"
	"time"

	// "fmt"
	"math/big"

	"github.com/Cedro-Finance/CedroAVS/verifierOperator/client/avs/abi"
	"github.com/Cedro-Finance/CedroAVS/verifierOperator/config"
	"github.com/Cedro-Finance/CedroAVS/verifierOperator/logger"
	ether "github.com/ethereum/go-ethereum"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

const (
	defaultStartBlock                    = 1234
	defaultHeightDifferenceForFilterLogs = 100
	defaultRequestId                     = 1

	defaultPollDuration = defaultHeightDifferenceForFilterLogs * 2 / 3 * time.Second
)

func (avs *AVSClient) GetCurrentBlock(ctx context.Context) (uint64, error) {
	currentBlock, err := avs.eth.BlockNumber(ctx)

	if err != nil {
		return 0, err
	}

	return currentBlock, nil
}

func (avs *AVSClient) GetFinalizedBlock(ctx context.Context) (uint64, error) {
	currentBlock, err := avs.GetCurrentBlock(ctx)

	if err != nil {
		logger.GetLogger().Error("Error while getting block information")

		return 0, err
	}

	return currentBlock - avs.blockFinality, nil

}

func (avs *AVSClient) FilterLogs(
	ctx context.Context, fromHeight uint64, toHeight uint64,
	contractAddress ethCommon.Address, topics ethCommon.Hash) ([]types.Log, error) {

	logs, err := avs.eth.FilterLogs(ctx, ether.FilterQuery{
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

func (avs *AVSClient) NewAVSContract() *abi.AVS {
	contract, err := abi.NewAVS(avs.contract, avs.eth)

	if err != nil {
		logger.GetLogger().Error("Error while initializing contract")
		panic(err)
	}

	return contract
}

func (avs *AVSClient) GetRequestId() (*big.Int, error) {

	contract := avs.NewAVSContract()
	reqId, err := contract.ReqId(nil)

	if err != nil {
		return big.NewInt(0), err
	}

	return reqId, nil
}

func (avs *AVSClient) ParseGetBalanceLog(log types.Log) (*abi.AVSGetChainTVL, error) {
	parsed, err := avs.NewAVSContract().ParseGetChainTVL(log)

	if err != nil {
		return nil, err
	}

	return parsed, nil
}

type AVSClient struct {
	contract         ethCommon.Address
	eth              *ethclient.Client
	startBlock       uint64
	nextBlockHeight  uint64
	pollDuration     time.Duration
	currentRequestId uint64

	blockTime     time.Duration
	blockFinality uint64
}

func NewAVSClient(cfg *config.AVSConfig) *AVSClient {
	rpc, err := rpc.Dial(cfg.NodeUrl)
	if err != nil {
		logger.GetLogger().Error("failed to create rpc client for avs")
	}

	eth := ethclient.NewClient(rpc)

	if err != nil {
		panic(err)
	}

	return &AVSClient{
		contract:         ethCommon.HexToAddress(cfg.ContactAddress),
		eth:              eth,
		startBlock:       cfg.StartBlock,
		nextBlockHeight:  cfg.StartBlock,
		pollDuration:     time.Duration(time.Minute),
		currentRequestId: 0,

		blockTime:     cfg.BlockTime * time.Second,
		blockFinality: cfg.BlockFinality,
	}
}

type Request struct {
	reqId  *big.Int
	height *big.Int
}

func GetRequestId(req *Request) *big.Int {
	return req.reqId
}

// polls the balance request events from the AVS contract on eth
func (c *AVSClient) PollEvents(ctx context.Context, ch chan<- *Request) {

	dur := c.pollDuration

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

	L1:
		for {
			finalizedBlock, err := c.GetFinalizedBlock(ctx)

			if err != nil {
				//handle error properly
				break L1
			}

			if finalizedBlock < c.nextBlockHeight {
				diff := c.nextBlockHeight - finalizedBlock
				time.Sleep((time.Duration(diff) * c.blockTime))
				break L1
			}

			currentReqId, err := c.GetRequestId()

			if err != nil {
				break L1
			}

			for startHeight := c.nextBlockHeight; startHeight <= finalizedBlock; startHeight += defaultHeightDifferenceForFilterLogs + 1 {
				endHeight := startHeight + defaultHeightDifferenceForFilterLogs

				if endHeight > finalizedBlock {
					endHeight = finalizedBlock
				}

				logs, err := c.FilterLogs(ctx, startHeight, endHeight, c.contract, ethCommon.Hash{})

				if err != nil {
					//Implement Logger
					break L1
				}

				for _, log := range logs {
					parsed, err := c.ParseGetBalanceLog(log)

					if err != nil {
						break L1
					}

					if currentReqId == parsed.ReqId {
						ch <- &Request{
							reqId:  parsed.ReqId,
							height: parsed.ReqId,
						}

						//TODO: Will there be condition when reqId 2 arrives before settling reqId 1?
						c.currentRequestId += 1
					}
				}

				c.nextBlockHeight = endHeight + 1
			}
		}

	}
}

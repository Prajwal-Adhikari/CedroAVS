package avs

import (
	"context"
	"time"

	// "fmt"
	"math/big"

	"github.com/Cedro-Finance/CedroAVS/verifierOperator/client/avs/abi"
	"github.com/Cedro-Finance/CedroAVS/verifierOperator/config"
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

// type ethClientI interface {
// 	GetCurrentBlock(ctx context.Context) (uint64, error)
// 	FilterLogs(
// 		ctx context.Context, fromHeight uint64, toHeight uint64,
// 		contractAddress ethCommon.Address, topics ethCommon.Hash) ([]types.Log, error)
// }

type ethClient struct {
	eth *ethclient.Client
}

func NewEthClient(rpcUrl string) *ethClient {
	rpcClient, err := rpc.Dial(rpcUrl)

	if err != nil {
		panic(err)
	}

	return &ethClient{
		eth: ethclient.NewClient(rpcClient),
	}
}

func (eth *ethClient) GetCurrentBlock(ctx context.Context) (uint64, error) {
	currentBlock, err := eth.eth.BlockNumber(ctx)

	if err != nil {
		return 0, err
	}

	return currentBlock, nil
}

func (eth *ethClient) FilterLogs(
	ctx context.Context, fromHeight uint64, toHeight uint64,
	contractAddress ethCommon.Address, topics ethCommon.Hash) ([]types.Log, error) {

	logs, err := eth.eth.FilterLogs(ctx, ether.FilterQuery{
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

type AVSContract struct {
	avs *abi.AVS
}

func NewAVSContract(addr ethCommon.Address, ethClient *ethclient.Client) (*AVSContract, error) {
	avs, err := abi.NewAVS(addr, ethClient)

	if err != nil {
		return nil, err
	}

	return &AVSContract{
		avs: avs,
	}, err
}

func (avs *AVSContract) GetRequestId() (*big.Int, error)  {
	reqId, err := avs.avs.ReqId(nil)

	if (err != nil ) {
		return big.NewInt(0), err
	}

	return reqId, nil
}

func (avs *AVSContract) ParseGetBalanceLog(log types.Log) (*abi.AVSGetChainTVL, error) {
	parsed, err := avs.avs.ParseGetChainTVL(log)

	if err != nil {
		return nil, err
	}

	return parsed, nil
}

type AVSClient struct {
	contract        *AVSContract
	contractAddress ethCommon.Address
	ethClient       *ethClient
	startBlock      uint64
	nextBlockHeight uint64
	pollDuration    time.Duration
}

func NewAVSClient(cfg *config.AVSConfig) *AVSClient {
	ethClient := NewEthClient(cfg.NodeUrl)
	avsContract, err := NewAVSContract(ethCommon.HexToAddress(cfg.ContactAddress), ethClient.eth)

	if err != nil {
		panic(err)
	}

	return &AVSClient{
		contract:        avsContract,
		ethClient:       ethClient,
		startBlock:      cfg.StartBlock,
		nextBlockHeight: cfg.StartBlock,
		contractAddress: ethCommon.HexToAddress(cfg.ContactAddress),
		pollDuration:    time.Duration(time.Minute),
	}
}


type Request struct{}

// polls the balance request events from the AVS contract on eth
func (c *AVSClient) PollEvents(ctx context.Context, ch chan<- *big.Int) {

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
			currentBlock, err := c.ethClient.GetCurrentBlock(ctx)

			if err != nil {
				//handle error properly
				break L1
			}

			currentReqId, err := c.contract.GetRequestId()

			if err != nil {
				break L1
			}

			for startHeight := c.nextBlockHeight; startHeight <= currentBlock; startHeight += defaultHeightDifferenceForFilterLogs + 1 {
				endHeight := startHeight + defaultHeightDifferenceForFilterLogs

				if endHeight > currentBlock {
					endHeight = currentBlock
				}

				logs, err := c.ethClient.FilterLogs(ctx, startHeight, endHeight, c.contractAddress, ethCommon.Hash{})

				if err != nil {
					//Implement Logger
					break L1
				}

				for _, log := range logs {
					parsed, err := c.contract.ParseGetBalanceLog(log)

					if err != nil {
						break L1
					}

					if currentReqId == parsed.ReqId {
						ch <- parsed.ReqId
					}
				}

				c.nextBlockHeight = endHeight + 1
			}
		}

	}
}

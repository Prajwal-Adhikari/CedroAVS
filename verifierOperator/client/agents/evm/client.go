package evm

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/Cedro-Finance/CedroAVS/core/logger"
	"github.com/Cedro-Finance/CedroAVS/verifierOperator/client/agents/evm/abi"
	"github.com/Cedro-Finance/CedroAVS/verifierOperator/client/avs"
	"github.com/Cedro-Finance/CedroAVS/verifierOperator/common"
	"github.com/Cedro-Finance/CedroAVS/verifierOperator/config"
	"github.com/Layr-Labs/eigensdk-go/logging"
	ether "github.com/ethereum/go-ethereum"

	// ethAbi "github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

func NewClient(config config.ChainConfig, logger logging.Logger) common.IAgentPodClient {
	rpc, err := rpc.Dial(config.NodeUrl)

	if err != nil {
		logger.Error("failed to create rpc client")
		return nil
	}

	ethClient := ethclient.NewClient(rpc)

	return &Client{
		name:          config.Name,
		chain:         config.Chain,
		eth:           ethClient,
		blockFinality: config.BlockFinality,
		blockTime:     time.Duration(config.BlockTime) * time.Second,
		logger:        logger,
	}
}

type Client struct {
	name          string
	chain         string
	eth           *ethclient.Client
	blockFinality uint64
	blockTime     time.Duration //in second
	logger        logging.Logger
	agent         *abi.AgentPod
}

func (c *Client) GetCurrentBlock(ctx context.Context) (uint64, error) {
	return c.eth.BlockNumber(ctx)
}

func (c *Client) GetFinalizedBlock(ctx context.Context) (uint64, error) {
	currentBlock, err := c.GetCurrentBlock(ctx)

	if err != nil {
		logger.GetLogger().Error("Error while getting block information")

		return 0, err
	}

	return currentBlock - c.blockFinality, nil

}

func (c *Client) NonceAt(ctx context.Context, account ethCommon.Address, blockNumber *big.Int) (uint64, error) {
	return c.eth.NonceAt(ctx, account, blockNumber)
}

func (c *Client) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return c.eth.SuggestGasPrice(ctx)
}

func (c *Client) FilterLogs(
	ctx context.Context, fromHeight uint64, toHeight uint64,
	contractAddress ethCommon.Address, topics ethCommon.Hash) ([]types.Log, error) {

	logs, err := c.eth.FilterLogs(ctx, ether.FilterQuery{
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

func (c *Client) NewAgentContract(addr ethCommon.Address) *abi.AgentPod {
	agent, err := abi.NewAgentPod(addr, c.eth)

	if err != nil {
		fmt.Println("Contract not initialized")
		panic(err)
	}

	c.agent = agent

	return agent
}

func (c *Client) InitializeAgentPod(addr ethCommon.Address) {
	agent, err := abi.NewAgentPod(addr, c.eth)

	if err != nil {
		c.logger.Errorf("Error while initializing contract", "err", err)

		return
	}

	c.agent = agent
}

func (c *Client) GetBlockFromTimestamp(ctx context.Context, timestamp uint64) (uint64, error) {
	currentBlock, err := c.GetCurrentBlock(ctx)

	if err != nil {
		fmt.Println("Error getting block")

		return 0, err
	}

	block, err := c.eth.BlockByNumber(ctx, big.NewInt(int64(currentBlock)))

	if err != nil {
		fmt.Println("Error getting block info")
		return 0, err
	}

	timediff := block.Time() - timestamp

	blockDiff := timediff / uint64(c.blockTime)

	finalBlock := currentBlock + blockDiff

	return finalBlock, err
}

func (c *Client) GetBalance(ctx context.Context, reqEvent *avs.Request, balanceChan chan<- *big.Int) {
	if c.agent == nil {
		fmt.Print("Agent not initialized")
	}

	finalizedBlock, err := c.GetFinalizedBlock(ctx)

	if err != nil {
		logger.GetLogger().Error("Error while getting finalized block")
	}

	blockByTimestamp, err := c.GetBlockFromTimestamp(ctx, reqEvent.Timestamp)

	if err != nil {
		logger.GetLogger().Error("Error whole geting block from timestamp")
	}

	if blockByTimestamp > finalizedBlock {
		time.Sleep(time.Duration(blockByTimestamp-finalizedBlock) * c.blockTime)
	}

	balance, err := c.agent.GetBalanceOfAgent(&bind.CallOpts{
		BlockNumber: big.NewInt(int64(blockByTimestamp)),
	})

	if err != nil {
		fmt.Println("Errro while getting balance")
	}

	balanceChan <- balance
}

package evm

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/Cedro-Finance/CedroAVS/api/graphql"
	agentpod "github.com/Cedro-Finance/CedroAVS/core/bindings/AgentPod"
	"github.com/Cedro-Finance/CedroAVS/core/logger"
	"github.com/Cedro-Finance/CedroAVS/verifier/client/avs"
	"github.com/Cedro-Finance/CedroAVS/verifier/common"
	"github.com/Cedro-Finance/CedroAVS/verifier/config"
	"github.com/Layr-Labs/eigensdk-go/logging"

	// ethAbi "github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethCommon "github.com/ethereum/go-ethereum/common"
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

	gClient := graphql.NewGraphQLClient(config.GraphqlUrl, logger)

	return &Client{
		Name:          config.Name,
		Chain:         config.Chain,
		ChainId:       config.ChainId,
		Eth:           ethClient,
		BlockFinality: config.BlockFinality,
		BlockTime:     time.Duration(config.BlockTime) * time.Second,
		Logger:        logger,
		GraphClient:   gClient,
	}
}

type Client struct {
	Name          string
	Chain         string
	Eth           *ethclient.Client
	BlockFinality uint64
	BlockTime     time.Duration //in second
	Logger        logging.Logger
	Agent         *agentpod.AgentPod
	GraphClient   *graphql.GraphQLClient
	ChainId       string
}

func (c *Client) GetCurrentBlock(ctx context.Context) (uint64, error) {
	return c.Eth.BlockNumber(ctx)
}

func (c *Client) GetGraphClient() *graphql.GraphQLClient {
	return c.GraphClient
}

func (c *Client) GetChainId() string {
	return c.Chain
}

func (c *Client) GetFinalizedBlock(ctx context.Context) (uint64, error) {
	currentBlock, err := c.GetCurrentBlock(ctx)

	if err != nil {
		logger.GetLogger().Error("Error while getting block information")

		return 0, err
	}

	return currentBlock - c.BlockFinality, nil

}

func (c *Client) NewAgentContract(addr ethCommon.Address) *agentpod.AgentPod {
	agent, err := agentpod.NewAgentPod(addr, c.Eth)

	if err != nil {
		fmt.Println("Contract not initialized")
		panic(err)
	}

	c.Agent = agent

	return agent
}

func (c *Client) InitializeAgentPod(addr ethCommon.Address) {
	agent, err := agentpod.NewAgentPod(addr, c.Eth)

	if err != nil {
		c.Logger.Errorf("Error while initializing contract", "err", err)

		return
	}

	c.Agent = agent
}

func (c *Client) GetBlockFromTimestamp(ctx context.Context, timestamp uint64) (uint64, error) {
	currentBlock, err := c.GetCurrentBlock(ctx)

	if err != nil {
		fmt.Println("Error getting block")

		return 0, err
	}

	block, err := c.Eth.BlockByNumber(ctx, big.NewInt(int64(currentBlock)))

	if err != nil {
		fmt.Println("Error getting block info")
		return 0, err
	}

	timediff := block.Time() - timestamp

	blockDiff := timediff / uint64(c.BlockTime)

	finalBlock := currentBlock + blockDiff

	return finalBlock, err
}

func (c *Client) GetBalance(ctx context.Context, reqEvent *avs.Request, balanceChan chan<- *big.Int) {

	var balance *big.Int

	maxRetries := 10
	retryCount := 0

RETRY_LOOP:
	for {
		if retryCount <= maxRetries {
			finalizedBlock, err := c.GetFinalizedBlock(ctx)

			if err != nil {
				c.Logger.Error("Error while getting finalized block", "chain", c.ChainId, "err", err, "retryCount", retryCount)
				retryCount++
				break RETRY_LOOP
			}

			blockByTimestamp, err := c.GetBlockFromTimestamp(ctx, reqEvent.Timestamp)

			if err != nil {
				c.Logger.Error("Error while getting block by timestamp", "chain", c.ChainId, "err", err, "retryCount", retryCount)
				retryCount++
				return
			}

			if blockByTimestamp > finalizedBlock {
				time.Sleep(time.Duration(blockByTimestamp-finalizedBlock) * c.BlockTime)
			}

			balance, err = c.Agent.PrincipalBalanceOfAgentPod(&bind.CallOpts{
				BlockNumber: big.NewInt(int64(blockByTimestamp)),
			})

			if err != nil {
				c.Logger.Error("Error while getting balance", "chain", c.ChainId, "err", err, "retryCount", retryCount)
				retryCount++
				return
			}

			balanceChan <- balance
			return
		} else {
			c.Logger.Error("max retry threshold reached, error getting the balance from the agendpod client")

			balanceChan <- nil
			return
		}
	}

}

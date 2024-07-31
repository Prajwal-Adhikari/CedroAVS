package evm

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/Cedro-Finance/CedroAVS/api/graphql"
	avs "github.com/Cedro-Finance/CedroAVS/core/avsClient"
	agentmanager "github.com/Cedro-Finance/CedroAVS/core/bindings/AgentManager"
	agentpod "github.com/Cedro-Finance/CedroAVS/core/bindings/AgentPod"
	"github.com/Cedro-Finance/CedroAVS/core/config"
	"github.com/Cedro-Finance/CedroAVS/core/logger"
	"github.com/Cedro-Finance/CedroAVS/verifier/client/agents"
	"github.com/Cedro-Finance/CedroAVS/verifier/common"
	"github.com/Layr-Labs/eigensdk-go/logging"

	// ethAbi "github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

const (
	defaultStartBlock                    = 1705750
	defaultHeightDifferenceForFilterLogs = 100
	defaultRequestId                     = 1

	defaultPollDuration = 5 * time.Second
)

type IAgentManager interface {
	ParseDeviationCheck(log types.Log) (*agentmanager.AgentManagerDeviationCheck, error)
	ParseSlashAgent(log types.Log) (*agentmanager.AgentManagerSlashAgent, error)
}

func NewClient(config config.ChainConfig, logger logging.Logger, avsReader avs.IAVSReader) common.IAgentPodClient {
	rpc, err := rpc.Dial(config.NodeUrl)

	if err != nil {
		logger.Error("failed to create rpc client")
		return nil
	}

	ethClient := ethclient.NewClient(rpc)

	gClient := graphql.NewGraphQLClient(config.GraphqlUrl, logger)

	agentManager, err := agentmanager.NewAgentManager(ethCommon.HexToAddress(config.AgentManager), ethClient)

	if err != nil {
		panic(err)
	}

	return &Client{
		eth:             ethClient,
		blockTime:       time.Duration(config.BlockTime) * time.Second,
		logger:          logger,
		graphClient:     gClient,
		config:          config,
		avsReader:       avsReader,
		nextBlockHeight: config.AgentManagerStartHeight,
		agentManager:    agentManager,
	}
}

type Client struct {
	eth          common.EthClient
	blockTime    time.Duration //in second
	logger       logging.Logger
	agent        *agentpod.AgentPod
	graphClient  *graphql.GraphQLClient
	agentManager IAgentManager

	startHeight     uint64
	nextBlockHeight uint64

	config    config.ChainConfig
	avsReader avs.IAVSReader
}

func (c *Client) GetCurrentBlock(ctx context.Context) (uint64, error) {
	return c.eth.BlockNumber(ctx)
}

func (c *Client) GetGraphClient() *graphql.GraphQLClient {
	return c.graphClient
}

func (c *Client) GetChainId() string {
	return c.config.Chain
}

func (c *Client) GetFinalizedBlock(ctx context.Context) (uint64, error) {
	currentBlock, err := c.GetCurrentBlock(ctx)

	if err != nil {
		logger.GetLogger().Error("Error while getting block information")

		return 0, err
	}

	return currentBlock - c.config.BlockFinality, nil
}

func (c *Client) InitializeAgentPod(addr ethCommon.Address) {
	agent, err := agentpod.NewAgentPod(addr, c.eth)

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

	if block.Time() < timestamp {
		return currentBlock, nil
	}

	timediff := block.Time() - timestamp

	blockDiff := timediff / uint64(c.blockTime.Seconds())

	finalBlock := currentBlock - blockDiff

	return finalBlock, err
}

func (c *Client) GetBalance(ctx context.Context, reqEvent *agents.DeviationRequest, balanceChan chan<- *big.Int) {

	var balance *big.Int

	maxRetries := 10
	retryCount := 0

RETRY_LOOP:
	for {
		if retryCount <= maxRetries {
			finalizedBlock, err := c.GetFinalizedBlock(ctx)

			if err != nil {
				c.logger.Error("Error while getting finalized block", "chain", c.config.ChainId, "err", err, "retryCount", retryCount)
				retryCount++
				break RETRY_LOOP
			}

			blockByTimestamp, err := c.GetBlockFromTimestamp(ctx, reqEvent.Timestamp)

			if err != nil {
				c.logger.Error("Error while getting block by timestamp", "chain", c.config.ChainId, "err", err, "retryCount", retryCount)
				retryCount++
				return
			}

			if blockByTimestamp > finalizedBlock {
				time.Sleep(time.Duration(blockByTimestamp-finalizedBlock) * c.blockTime)
			}

			balance, err = c.agent.PrincipalBalanceOfAgentPod(&bind.CallOpts{
				BlockNumber: big.NewInt(int64(blockByTimestamp)),
			})

			if err != nil {
				c.logger.Error("Error while getting balance", "chain", c.config.ChainId, "err", err, "retryCount", retryCount)
				retryCount++
				return
			}

			balanceChan <- balance
			return
		} else {
			c.logger.Error("max retry threshold reached, error getting the balance from the agendpod client")

			balanceChan <- nil
			return
		}
	}

}

func (c *Client) getDeviationCheckEventSign() ethCommon.Hash {
	managerAbi, err := agentmanager.AgentManagerMetaData.GetAbi()

	if err != nil {
		return ethCommon.Hash{}
	}

	event := managerAbi.Events["DeviationCheck"]

	return event.ID
}

func (c *Client) getAgentSlashEventSign() ethCommon.Hash {
	managerAbi, err := agentmanager.AgentManagerMetaData.GetAbi()

	if err != nil {
		return ethCommon.Hash{}
	}

	event := managerAbi.Events["SlashAgent"]

	return event.ID
}

func (c *Client) processDeviationRequest(ctx context.Context, log types.Log, ch chan<- agents.TaskRequest) error {

	parsed, err := c.agentManager.ParseDeviationCheck(log)

	if err != nil {
		return err
	}

	blockInfo, err := c.eth.BlockByHash(ctx, log.BlockHash)

	if err != nil {
		c.logger.Error("Error while getting block inforamtion", "err", err)
		return err
	}

	c.logger.Info("Deviation Request Received",
		"RequestId", parsed.RequestID.String(),
		"Agent Address", parsed.Agent,
		"BlockNumber", log.BlockNumber,
	)

	lastReqId, err := c.avsReader.GetDeviationRequestId(ctx, parsed.Agent, uint32(log.BlockNumber))

	if err != nil {
		c.logger.Error("Error while fetching last requestId", "err", err)
		return err
	}

	if lastReqId.Cmp(parsed.RequestID) > -1 {
		c.logger.Warn("Request Id already processed", "reqId", parsed.RequestID.String(), "agent", parsed.Agent.Hex())
		return nil
	}

	ch <- &agents.DeviationRequest{
		ReqId:                parsed.RequestID,
		ReferenceBlockNumber: log.BlockNumber,
		AgentAddress:         parsed.Agent,
		Timestamp:            blockInfo.Time(),
	}

	return nil

}

func (c *Client) processAgentSlashRequest(ctx context.Context, log types.Log, ch chan<- agents.TaskRequest) error {

	parsed, err := c.agentManager.ParseSlashAgent(log)

	if err != nil {
		return err
	}

	c.logger.Info("Slash Agent Event Received", "req", parsed)

	ch <- &agents.AgentSlashRequest{
		TxnId:                parsed.TxnId,
		ReferenceBlockNumber: log.BlockNumber,
		AgentAddress:         parsed.Agent,
		InsufficientBalance:  parsed.AmountToSlash,
	}

	return nil
}

func (c *Client) PollEvents(ctx context.Context, ch chan<- agents.TaskRequest, slashCh chan<- agents.TaskRequest) {
	dur := time.Duration(0)

	if dur == 0 {
		dur = time.Duration(defaultPollDuration)
	}

	ticker := time.NewTicker(dur)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}

	POLL_LOOP:
		for {
			finalizedBlock, err := c.GetFinalizedBlock(ctx)

			if err != nil {
				c.logger.Error("Error while getting finalized block", "err", err)
				break POLL_LOOP
			}

			if finalizedBlock%100 == 0 {
				c.logger.Info("Polling Block", "blockk", finalizedBlock)
			}

			if finalizedBlock < c.nextBlockHeight {
				diff := c.nextBlockHeight - finalizedBlock

				time.Sleep(time.Duration(diff) * c.blockTime)
			}

			for startHeight := c.nextBlockHeight; startHeight <= finalizedBlock; startHeight += defaultHeightDifferenceForFilterLogs + 1 {

				endHeight := startHeight + defaultHeightDifferenceForFilterLogs

				if endHeight > finalizedBlock {
					endHeight = finalizedBlock
				}

				logs, err := c.eth.FilterLogs(ctx, ethereum.FilterQuery{
					FromBlock: big.NewInt(int64(startHeight)),
					ToBlock:   big.NewInt(int64(endHeight)),
					Topics:    [][]ethCommon.Hash{{c.getDeviationCheckEventSign(), c.getAgentSlashEventSign()}},
					Addresses: []ethCommon.Address{ethCommon.HexToAddress(c.config.AgentManager)},
				})

				if err != nil {
					c.logger.Error("Error while getting logs", "err", err)
					break POLL_LOOP
				}

				for _, log := range logs {
					switch log.Topics[0] {
					case c.getDeviationCheckEventSign():
						err := c.processDeviationRequest(ctx, log, ch)

						if err != nil {
							break POLL_LOOP
						}
						//handle deviation check
					case c.getAgentSlashEventSign():
						err := c.processAgentSlashRequest(ctx, log, slashCh)

						if err != nil {
							break POLL_LOOP
						}
					default:
						c.logger.Error("couldn't parse event log")
					}
				}

				c.nextBlockHeight = endHeight + 1

			}

		}
	}

}

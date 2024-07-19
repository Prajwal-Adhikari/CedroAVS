package evm

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/Cedro-Finance/CedroAVS/api/graphql"
	agentmanager "github.com/Cedro-Finance/CedroAVS/core/bindings/AgentManager"
	agentpod "github.com/Cedro-Finance/CedroAVS/core/bindings/AgentPod"
	"github.com/Cedro-Finance/CedroAVS/core/logger"
	"github.com/Cedro-Finance/CedroAVS/verifier/client/avs"
	"github.com/Cedro-Finance/CedroAVS/verifier/common"
	"github.com/Cedro-Finance/CedroAVS/verifier/config"
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
	defaultBlockTime = time.Duration(12) * time.Second

	defaultStartBlock                    = 1705750
	defaultHeightDifferenceForFilterLogs = 100
	defaultRequestId                     = 1

	defaultPollDuration = 5 * time.Second
)

func NewClient(config config.ChainConfig, logger logging.Logger, avsReader *avs.AvsReader) common.IAgentPodClient {
	rpc, err := rpc.Dial(config.NodeUrl)

	if err != nil {
		logger.Error("failed to create rpc client")
		return nil
	}

	ethClient := ethclient.NewClient(rpc)

	gClient := graphql.NewGraphQLClient(config.GraphqlUrl, logger)

	return &Client{
		Eth:             ethClient,
		BlockTime:       time.Duration(config.BlockTime) * time.Second,
		Logger:          logger,
		GraphClient:     gClient,
		AgentManager:    ethCommon.HexToAddress(config.AgentManager),
		config:          config,
		avsReader:       avsReader,
		NextBlockHeight: config.AgentManagerStartHeight,
	}
}

type Client struct {
	Eth          *ethclient.Client
	BlockTime    time.Duration //in second
	Logger       logging.Logger
	Agent        *agentpod.AgentPod
	GraphClient  *graphql.GraphQLClient
	AgentManager ethCommon.Address

	StartHeight     uint64
	NextBlockHeight uint64

	config config.ChainConfig

	avsReader *avs.AvsReader
}

func (c *Client) GetCurrentBlock(ctx context.Context) (uint64, error) {
	return c.Eth.BlockNumber(ctx)
}

func (c *Client) GetGraphClient() *graphql.GraphQLClient {
	return c.GraphClient
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

func (c *Client) GetBalance(ctx context.Context, reqEvent *avs.DeviationRequest, balanceChan chan<- *big.Int) {

	var balance *big.Int

	maxRetries := 10
	retryCount := 0

RETRY_LOOP:
	for {
		if retryCount <= maxRetries {
			finalizedBlock, err := c.GetFinalizedBlock(ctx)

			if err != nil {
				c.Logger.Error("Error while getting finalized block", "chain", c.config.ChainId, "err", err, "retryCount", retryCount)
				retryCount++
				break RETRY_LOOP
			}

			blockByTimestamp, err := c.GetBlockFromTimestamp(ctx, reqEvent.Timestamp)

			if err != nil {
				c.Logger.Error("Error while getting block by timestamp", "chain", c.config.ChainId, "err", err, "retryCount", retryCount)
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
				c.Logger.Error("Error while getting balance", "chain", c.config.ChainId, "err", err, "retryCount", retryCount)
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

func (c *Client) GetDeviationCheckEventSign() ethCommon.Hash {
	managerAbi, err := agentmanager.AgentManagerMetaData.GetAbi()

	if err != nil {
		return ethCommon.Hash{}
	}

	event := managerAbi.Events["DeviationCheck"]

	return event.ID
}

func (c *Client) GetAgentSlashEventSign() ethCommon.Hash {
	managerAbi, err := agentmanager.AgentManagerMetaData.GetAbi()

	if err != nil {
		return ethCommon.Hash{}
	}

	event := managerAbi.Events["SlashAgent"]

	return event.ID
}

func (c *Client) ProcessDeviationRequest(ctx context.Context, log types.Log, ch chan<- avs.TaskRequest) error {
	manager, err := agentmanager.NewAgentManager(c.AgentManager, c.Eth)

	if err != nil {
		return err
	}

	parsed, err := manager.ParseDeviationCheck(log)

	if err != nil {
		return err
	}

	blockInfo, err := c.Eth.BlockByHash(ctx, log.BlockHash)

	if err != nil {
		c.Logger.Error("Error while getting block inforamtion", "err", err)
		return err
	}

	c.Logger.Info("Deviation Request Received",
		"RequestId", parsed.RequestID.String(),
		"Agent Address", parsed.Agent,
		"BlockNumber", log.BlockNumber,
	)

	lastReqId, err := c.avsReader.GetDeviationRequestId(ctx, parsed.Agent, uint32(log.BlockNumber))

	if err != nil {
		c.Logger.Error("Error while fetching last requestId", "err", err)
		return err
	}

	if lastReqId.Cmp(parsed.RequestID) > -1 {
		c.Logger.Warn("Request Id already processed", "reqId", parsed.RequestID.String(), "agent", parsed.Agent.Hex())
		return nil
	}

	ch <- &avs.DeviationRequest{
		ReqId:                parsed.RequestID,
		ReferenceBlockNumber: log.BlockNumber,
		AgentAddress:         parsed.Agent,
		Timestamp:            blockInfo.Time(),
	}

	return nil

}

func (c *Client) ProcessAgentSlashRequest(ctx context.Context, log types.Log, ch chan<- avs.TaskRequest) error {

	manager, err := agentmanager.NewAgentManager(c.AgentManager, c.Eth)

	if err != nil {
		return err
	}

	parsed, err := manager.ParseSlashAgent(log)

	if err != nil {
		return err
	}

	c.Logger.Info("Slash Agent Event Received", "req", parsed)

	ch <- &avs.AgentSlashRequest{
		TxnId:                big.NewInt(10),
		ReferenceBlockNumber: 1957097,
		AgentAddress:         ethCommon.HexToAddress("0x17835C075FaBdd8c8fCC6Be06e924e3BAB0fCc3e"),
		InsufficientBalance:  big.NewInt(10000000000000000),
	}

	return nil
}

func (c *Client) PollEvents(ctx context.Context, ch chan<- avs.TaskRequest, slashCh chan<- avs.TaskRequest) {
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
				c.Logger.Error("Error while getting finalized block", "err", err)
				break POLL_LOOP
			}

			if finalizedBlock%100 == 0 {
				c.Logger.Info("Polling Block", "blockk", finalizedBlock)
			}

			if finalizedBlock < c.NextBlockHeight {
				diff := c.NextBlockHeight - finalizedBlock

				time.Sleep(time.Duration(diff) * c.BlockTime)
			}

			for startHeight := c.NextBlockHeight; startHeight <= finalizedBlock; startHeight += defaultHeightDifferenceForFilterLogs + 1 {

				endHeight := startHeight + defaultHeightDifferenceForFilterLogs

				if endHeight > finalizedBlock {
					endHeight = finalizedBlock
				}

				logs, err := c.Eth.FilterLogs(ctx, ethereum.FilterQuery{
					FromBlock: big.NewInt(int64(startHeight)),
					ToBlock:   big.NewInt(int64(endHeight)),
					Topics:    [][]ethCommon.Hash{{c.GetDeviationCheckEventSign(), c.GetAgentSlashEventSign()}},
					Addresses: []ethCommon.Address{c.AgentManager},
				})

				if err != nil {
					c.Logger.Error("Error while getting logs", "err", err)
					break POLL_LOOP
				}

				for _, log := range logs {
					switch log.Topics[0] {
					case c.GetDeviationCheckEventSign():
						err := c.ProcessDeviationRequest(ctx, log, ch)

						if err != nil {
							break POLL_LOOP
						}
						//handle deviation check
					case c.GetAgentSlashEventSign():
						err := c.ProcessAgentSlashRequest(ctx, log, slashCh)

						if err != nil {
							break POLL_LOOP
						}
					default:
						c.Logger.Error("couldn't parse event log")
					}
				}

				c.NextBlockHeight = endHeight + 1

			}

		}
	}

}

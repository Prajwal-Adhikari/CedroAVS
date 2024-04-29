package evm

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/Cedro-Finance/CedroAVS/verifierOperator/client/agents"
	"github.com/Cedro-Finance/CedroAVS/verifierOperator/client/agents/evm/abi"
	"github.com/Cedro-Finance/CedroAVS/verifierOperator/common"
	"github.com/Cedro-Finance/CedroAVS/verifierOperator/config"
	"github.com/Cedro-Finance/CedroAVS/verifierOperator/keys"
	"github.com/Cedro-Finance/CedroAVS/verifierOperator/logger"
	ether "github.com/ethereum/go-ethereum"

	// ethAbi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

const (
	defaultStartBlock                    = 1234
	defaultHeightDifferenceForFilterLogs = 100
	defaultRequestId                     = 1
	defaultPollDuration                  = defaultHeightDifferenceForFilterLogs * 2 / 3 * time.Second

	maxRetry      = 10
	retryInterval = 10 * time.Second
)

func NewClient(config config.ChainConfig) common.IClient {
	rpc, err := rpc.Dial(config.NodeUrl)

	if err != nil {
		logger.GetLogger().Error("failed to create rpc client")
		return nil
	}

	ethClient := ethclient.NewClient(rpc)

	return &Client{
		name:            config.Name,
		chain:           config.Chain,
		contractAddress: ethCommon.HexToAddress(config.ContactAddress),
		eth:             ethClient,
		startBlock:      100,
		nextBlockHeight: 100,
		pollDuration:    time.Second * 100,
		blockFinality:   config.BlockFinality,
		blockTime:       config.BlockTime * time.Second,
	}
}

type Client struct {
	name            string
	chain           string
	contractAddress ethCommon.Address

	eth             *ethclient.Client
	pollDuration    time.Duration
	startBlock      uint64
	nextBlockHeight uint64
	blockFinality   uint64
	blockTime       time.Duration //in second
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

func (c *Client) NewAgentContract() *abi.Agent {
	agent, err := abi.NewAgent(c.contractAddress, c.eth)

	if err != nil {
		fmt.Println("Contract not initialized")
		panic(err)
	}

	return agent
}

func (c *Client) Balance(ctx context.Context, balanceChan chan *agents.BalanceResponse) {
	// agent := c.NewAgentContract()

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
	pollLoop:
		for {
			finalizedBlock, err := c.GetFinalizedBlock(ctx)

			if err != nil {
				break pollLoop
			}

			if finalizedBlock < c.nextBlockHeight {
				diff := c.nextBlockHeight - finalizedBlock
				time.Sleep((time.Duration(diff) * c.blockTime))
				break pollLoop
			}

			for startHeight := c.nextBlockHeight; startHeight <= finalizedBlock; startHeight += defaultHeightDifferenceForFilterLogs + 1 {
				endHeight := startHeight + defaultHeightDifferenceForFilterLogs + 1

				if endHeight > finalizedBlock {
					endHeight = finalizedBlock
				}

				logs, err := c.FilterLogs(ctx, startHeight, endHeight, c.contractAddress, ethCommon.Hash{})

				if err != nil {
					break pollLoop
				}

				for _, log := range logs {
					parsed, err := c.NewAgentContract().ParseBalance(log)

					if err != nil {
						logger.GetLogger().Error("Error while parsing Balance event")
						break pollLoop
					}

					fmt.Print(parsed)
				}

				c.nextBlockHeight = endHeight + 1
			}
		}

	}
}

// snapshots the balace of lp in the respective chain
// sendtx method
func (c *Client) UpdateBalance(ctx context.Context, rxId *big.Int, address ethCommon.Address) *types.Transaction {

	ks := keys.LoadKeys(config.GetConfig().KeyStore)

	retries := 0

	var transaction *types.Transaction

	for retries <= maxRetry {
		time.Sleep(retryInterval)

		agent := c.NewAgentContract()
		gasPrice, err := c.SuggestGasPrice(ctx)

		if err != nil {
			logger.GetLogger().Error("Error getting gas price")
			continue
		}

		auth, err := bind.NewKeyedTransactorWithChainID(ks.PrivateKey, big.NewInt(10))

		if err != nil {
			logger.GetLogger().Error("Error getting private key")
			continue
		}

		nonce, err := c.eth.PendingNonceAt(ctx, crypto.PubkeyToAddress(*ks.PublicKey))

		if err != nil {
			logger.GetLogger().Error("Error getting Nonce")
			continue
		}

		// agAbi, _ := abi.AgentMetaData.GetAbi()

		// ethAbi.ABI.Pack(*agAbi, "", rxId);

		auth.GasPrice = gasPrice
		auth.Nonce = big.NewInt(int64(nonce))
		auth.GasLimit = uint64(600000) //TODO: Replace with actual value

		res, err := agent.UpdateBalance(auth, rxId, rxId)

		if err != nil {
			logger.GetLogger().Error("Error while transacting")
			retries++
			continue
		}

		transaction = res
	}

	// if (retries == maxRetry && transaction == nil) {
	// 	// operator no need to loop now.
	// }

	return transaction
}

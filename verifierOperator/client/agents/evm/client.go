package evm

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"time"

	"github.com/Cedro-Finance/CedroAVS/verifierOperator/client/agents/evm/abi"
	"github.com/Cedro-Finance/CedroAVS/verifierOperator/client/avs"
	"github.com/Cedro-Finance/CedroAVS/verifierOperator/common"
	"github.com/Cedro-Finance/CedroAVS/verifierOperator/config"
	"github.com/Cedro-Finance/CedroAVS/verifierOperator/logger"
	ether "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

const (
	defaultStartBlock                    = 1234
	defaultHeightDifferenceForFilterLogs = 100
	defaultRequestId                     = 1
	defaultPollDuration                  = defaultHeightDifferenceForFilterLogs * 2 / 3 * time.Second

	maxRetry = 10
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
		pollDuration: time.Second * 100,
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
}

func (c *Client) GetCurrentBlock(ctx context.Context) (uint64, error) {
	return c.eth.BlockNumber(ctx)
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

func (c *Client) Balance(balanceRequest *avs.Request, balanceChan chan<- *big.Int) *big.Int {
	// agent := c.NewAgentContract()

	return big.NewInt(10)

	// return agent.
}

// snapshots the balace of lp in the respective chain
// sendtx method
func (c *Client) UpdateBalance(rxId uint64, address ethCommon.Address) *types.Transaction{

	retries := 0

	var transaction *types.Transaction;

	for retries <= maxRetry {
		time.Sleep(retryInterval)
		
		agent := c.NewAgentContract()
		ctx := context.TODO()

		gasPrice, err := c.SuggestGasPrice(ctx)

		if err != nil {
			logger.GetLogger().Error("Error getting gas price")
			continue
		}

		auth, err := bind.NewKeyedTransactorWithChainID(&ecdsa.PrivateKey{}, big.NewInt(10))

		if err != nil {
			logger.GetLogger().Error("Error getting private key")
			continue
		}

		nonce, err := c.eth.PendingNonceAt(ctx, address)

		if err != nil {
			logger.GetLogger().Error("Error getting Nonce")
			continue
		}


		auth.GasPrice = gasPrice
		auth.Nonce = big.NewInt(int64(nonce))
		auth.GasLimit = uint64(600000) //TODO: Replace with actual value

		res, err := agent.UpdateBalance(auth, big.NewInt(int64(rxId)), big.NewInt(int64(rxId)))

		if err != nil {
			logger.GetLogger().Error("Error while transacting")
			retries++
			continue
		}

		transaction = res
	}
	
	return transaction
}

package eventlistener

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/Cedro-Finance/CedroAVS/agentEnroller/config"
	contractutils "github.com/Cedro-Finance/CedroAVS/agentEnroller/contract/utils"
	db "github.com/Cedro-Finance/CedroAVS/agentEnroller/database"
	"github.com/Cedro-Finance/CedroAVS/agentEnroller/eventhandler"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

func SyncToLatest() {
	appConfig := config.LoadConfig()
	RootChainData := appConfig.RootChainData
	FromBlock := appConfig.DeafultStartBlock

	// Retrieve the last confirmed block number from the database
	confirmedBlockNumber, err := db.ConfirmedBlockInstance().Get(RootChainData.ShortName)
	if err == nil {
		FromBlock = confirmedBlockNumber.BlockNumber
	}

	// Reduced one block from current synced state to resync any missed event on latest sync
	FromBlock = FromBlock - 1

	fmt.Printf("Syncing block events from block %v\n", FromBlock)

	client, err := contractutils.GetClient(RootChainData.NodeUrl)
	if err != nil {
		log.Fatal(err)
	}

	// Get the latest block number from chain
	latestBlock, err := client.BlockByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	latestBlockNumber := int64(latestBlock.NumberU64())
	fmt.Printf("Latest block: %v\n", latestBlockNumber)

	const blockRangeLimit = 50000

	for from := FromBlock; from <= latestBlockNumber; from += blockRangeLimit {
		to := from + blockRangeLimit - 1
		if to > latestBlockNumber {
			to = latestBlockNumber
		}

		fmt.Printf("Processing blocks from %v to %v\n", from, to)

		query := ethereum.FilterQuery{
			FromBlock: big.NewInt(int64(from)),
			ToBlock:   big.NewInt(int64(to)),
			Addresses: []common.Address{
				RootChainData.AgentManager,
			},
			Topics: [][]common.Hash{{common.HexToHash("0x2f10946b561f38d0050eb99e7eaa47b08ff6387d5ee19942bcfc754a7b530906")}},
		}

		logs, err := client.FilterLogs(context.Background(), query)
		if err != nil {
			log.Fatal(err)
		}

		contractAbi, err := abi.JSON(strings.NewReader(string(config.GetAgentManagerAbi())))
		if err != nil {
			log.Fatal(err)
		}

		for _, vLog := range logs {
			eventhandler.HandleEventReceived(contractAbi, vLog)
		}

		if err := db.ConfirmedBlockInstance().Update(RootChainData.ShortName, to); err != nil {
			log.Printf("Failed to update confirmed block: %v", err)
		}
	}

	fmt.Println("Syncing completed.")
}

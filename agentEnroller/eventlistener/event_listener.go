package eventlistener

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Cedro-Finance/CedroAVS/agentEnroller/config"
	"github.com/Cedro-Finance/CedroAVS/agentEnroller/eventhandler"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func Start() {
	contractABI := config.GetAgentManagerAbi()
	contractAddress := config.LoadConfig().RootChainData.AgentManager
	wssURL := config.LoadConfig().RootChainData.NodeWSURL
	client, err := ethclient.Dial(wssURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	// Parse the contract's ABI
	parsedABI, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		log.Fatalf("Failed to parse contract ABI: %v", err)
	}

	// Create a filter query to listen for specific events
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	// Subscribe to the logs
	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatalf("Failed to subscribe to logs: %v", err)
	}

	// Listen for the events
	fmt.Printf("Listening for events")
	for {
		select {
		case err := <-sub.Err():
			log.Fatalf("Error while receiving logs: %v", err)
		case vLog := <-logs:
			eventhandler.HandleEventReceived(parsedABI, vLog)
		}
	}
}

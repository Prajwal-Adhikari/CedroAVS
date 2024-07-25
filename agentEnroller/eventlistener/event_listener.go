package eventlistener

import (
	"context"
	"fmt"
	"log"
	"strings"

	eventhandler "github.com/Cedro-Finance/CedroAVS/agentEnroller/eventhandler"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Replace this with your contract's ABI
const contractABI = `[
	 {
      "anonymous": false,
      "inputs": [
        {
          "indexed": false,
          "internalType": "address",
          "name": "agent",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "bytes32[]",
          "name": "chainIDs",
          "type": "bytes32[]"
        }
      ],
      "name": "AgentPodChainList",
      "type": "event"
    }
]`

// TODO: Update to actual abi and contract adresses for agent pod
const contractAddress = "0x37071ECa44FDF16141623665511380a55cE0efEB"

func Start() {
	// TODO: Use rpc url from conifg
	client, err := ethclient.Dial("wss://ethereum-holesky-rpc.publicnode.com")
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
		Addresses: []common.Address{common.HexToAddress(contractAddress)},
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
			fmt.Println("Log received: ", vLog) // Handle the received log

			event := struct {
				Agent    common.Address
				ChainIDs []common.Hash
			}{}
			fmt.Printf("Event data: %+v\n", event)

			eventName := "AgentPodChainList"

			err := parsedABI.UnpackIntoInterface(&event, eventName, vLog.Data)
			if err != nil {
				log.Fatalf("Failed to unpack log data: %v", err)
			}
			event.Agent = common.HexToAddress(event.Agent.Hex())

			chainIDs := make([]string, len(event.ChainIDs))
			for i, chainID := range event.ChainIDs {
				chainIDs[i] = chainID.Hex()
			}

			eventhandler.EnrollAgentHandler(event.Agent, chainIDs)
		}
	}
}

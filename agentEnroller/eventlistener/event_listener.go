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
		"inputs": [],
		"stateMutability": "nonpayable",
		"type": "constructor"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"internalType": "address",
				"name": "oldOwner",
				"type": "address"
			},
			{
				"indexed": true,
				"internalType": "address",
				"name": "newOwner",
				"type": "address"
			}
		],
		"name": "OwnerSet",
		"type": "event"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "newOwner",
				"type": "address"
			}
		],
		"name": "changeOwner",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "getOwner",
		"outputs": [
			{
				"internalType": "address",
				"name": "",
				"type": "address"
			}
		],
		"stateMutability": "view",
		"type": "function"
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
				OldOwner common.Address
				NewOwner common.Address
			}{}
			fmt.Printf("Event data: %+v\n", event)

			// TODO: Update event name
			eventName := "OwnerSet"

			err := parsedABI.UnpackIntoInterface(&event, eventName, vLog.Data)
			if err != nil {
				log.Fatalf("Failed to unpack log data: %v", err)
			}
			event.OldOwner = common.HexToAddress(vLog.Topics[1].Hex())
			event.NewOwner = common.HexToAddress(vLog.Topics[2].Hex())

			eventhandler.EnrollAgentHandler(event)
		}
	}
}

package eventhandler

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Cedro-Finance/CedroAVS/agentEnroller/config"
	contractcaller "github.com/Cedro-Finance/CedroAVS/agentEnroller/contract/caller"
	contractdeployer "github.com/Cedro-Finance/CedroAVS/agentEnroller/contract/deployer"
	contractutils "github.com/Cedro-Finance/CedroAVS/agentEnroller/contract/utils"
	"github.com/Cedro-Finance/CedroAVS/agentEnroller/utils"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type EventMetadata struct {
	ContractAddress common.Address
	BlockNumber     uint64
	TxHash          common.Hash
	logIndex        uint
}

func enrollAgentHandler(owner common.Address, chainIDS []string) {
	fmt.Printf("Event data: %v %v\n", owner, chainIDS)

	var chainsToEnroll = chainIDS

	appConfig := config.LoadConfig()
	supportedChains := appConfig.Chains

	addresses := make(map[string]common.Address)

	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, chain := range chainsToEnroll {
		wg.Add(1)
		go func(chain string) {
			defer wg.Done()
			chainBytes, _ := hex.DecodeString(chain[2:])
			var chainBytes32 [32]byte
			copy(chainBytes32[:], chainBytes[:32])
			chainName := utils.Bytes32ToString(chainBytes32)

			chainData := contractutils.GetChainData(chainName, supportedChains)

			_, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			log.Println(chainData.NodeUrl)
			contractAddress, err := contractdeployer.DeployContract(owner, chainData)
			if err != nil {
				fmt.Printf("Failed to deploy contract on chain %s: %v\n", chain, err)
				return
			}

			mu.Lock()
			addresses[chain] = contractAddress
			mu.Unlock()
		}(chain)
	}

	wg.Wait()

	fmt.Printf("Deployed contract addresses: %+v\n", addresses)
	contractcaller.MakeManagerContractCalls(owner, addresses)
}

func parseEventData(vLog types.Log) EventMetadata {
	return EventMetadata{
		ContractAddress: vLog.Address,
		BlockNumber:     vLog.BlockNumber,
		TxHash:          vLog.TxHash,
		logIndex:        vLog.Index,
	}
}

func HandleEventReceived(parsedABI abi.ABI, vLog types.Log) {
	eventName := "AgentPodChainList"
	fmt.Println("Log received: ", vLog)

	event := struct {
		Agent    common.Address
		ChainIDs []common.Hash
	}{}
	fmt.Printf("Event data: %+v\n", event)

	eventData := parseEventData(vLog)

	fmt.Print(eventData)
	fmt.Println()

	err := parsedABI.UnpackIntoInterface(&event, eventName, vLog.Data)
	if err != nil {
		log.Fatalf("Failed to unpack log data: %v", err)
	}
	event.Agent = common.HexToAddress(event.Agent.Hex())

	chainIDs := make([]string, len(event.ChainIDs))
	for i, chainID := range event.ChainIDs {
		chainIDs[i] = chainID.Hex()
	}

	enrollAgentHandler(event.Agent, chainIDs)
}

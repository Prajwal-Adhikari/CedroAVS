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
	db "github.com/Cedro-Finance/CedroAVS/agentEnroller/database"
	"github.com/Cedro-Finance/CedroAVS/agentEnroller/utils"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type EventMetadata struct {
	ContractAddress common.Address
	BlockNumber     uint64
	TxHash          common.Hash
	LogIndex        uint
	TimeStamp       uint
}

type EventData struct {
	Agent    common.Address
	ChainIDs []common.Hash
}
type EventPayload struct {
	EventMetadata
	Data EventData
}

func enrollAgentHandler(owner common.Address, chainIDS []string) {
	fmt.Printf("Event data: %v %v\n", owner, chainIDS)

	var chainsToEnroll = chainIDS

	appConfig := config.LoadConfig()
	supportedChains := appConfig.Chains

	addresses := make(map[string]common.Address)

	var mu sync.Mutex
	var wg sync.WaitGroup

	var contractStatuses []db.ChainDeployment

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

			contractAddress, txHash, err := contractdeployer.DeployContract(owner, chainData)

			if err != nil {
				contractStatuses = append(contractStatuses, db.ChainDeployment{
					ManagerAddress:         chainData.AgentManager.String(),
					ContractAddress:        contractAddress.String(),
					Chain:                  chainName,
					TransactionHash:        txHash.String(),
					Status:                 "failed",
					Error:                  err.Error(),
					DeploymentNotification: []db.DeploymentNotification{},
				})

				fmt.Printf("Failed to deploy contract on chain %s: %v\n", chain, err)
				return
			}

			contractStatuses = append(contractStatuses, db.ChainDeployment{
				ManagerAddress:         chainData.AgentManager.String(),
				ContractAddress:        contractAddress.String(),
				Chain:                  chainName,
				TransactionHash:        txHash.String(),
				Status:                 "success",
				Error:                  "",
				DeploymentNotification: []db.DeploymentNotification{},
			})

			mu.Lock()
			addresses[chain] = contractAddress
			mu.Unlock()
		}(chain)
	}

	wg.Wait()

	db.ContractDeploymentInstance().Add(db.DBContractDeploymentData{
		Owner:  owner.String(),
		Chains: contractStatuses,
	})

	fmt.Printf("Deployed contract addresses: %+v\n", addresses)
	contractcaller.MakeManagerContractCalls(owner, addresses)
}

func parseEventData(vLog types.Log) EventMetadata {
	return EventMetadata{
		ContractAddress: vLog.Address,
		BlockNumber:     vLog.BlockNumber,
		TxHash:          vLog.TxHash,
		LogIndex:        vLog.Index,
	}
}

func HandleEventReceived(parsedABI abi.ABI, vLog types.Log) {
	eventName := "AgentPodChainList"
	fmt.Println("Log received: ", vLog.Index)

	event := EventData{}
	fmt.Printf("Event data: %+v\n", event)

	eventMetaData := parseEventData(vLog)

	fmt.Println(vLog.Data)

	err := parsedABI.UnpackIntoInterface(&event, eventName, vLog.Data)
	if err != nil {
		log.Fatalf("Failed to unpack log data: %v", err)
	}
	event.Agent = common.HexToAddress(event.Agent.Hex())

	chainIDs := make([]string, len(event.ChainIDs))
	for i, chainID := range event.ChainIDs {
		chainIDs[i] = chainID.Hex()
	}
	eventPayload := EventPayload{
		EventMetadata: eventMetaData,
		Data:          event,
	}

	_, err = db.EventServiceInstance().Get(db.EventGetFilter{
		TransactionHash: eventMetaData.TxHash.Hex(),
		LogIndex:        int(eventMetaData.LogIndex),
	})

	if err == mongo.ErrNoDocuments {
		db.EventServiceInstance().Add(db.DBEventData{
			Agent:           event.Agent.Hex(),
			ChainIDs:        chainIDs,
			Timestamp:       time.Now(),
			TransactionHash: eventMetaData.TxHash.Hex(),
			LogIndex:        int(eventMetaData.LogIndex),
		})

		fmt.Print(eventPayload)

		enrollAgentHandler(event.Agent, chainIDs)
	}
}

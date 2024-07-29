package contractcaller

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"

	"sync"

	"github.com/Cedro-Finance/CedroAVS/agentEnroller/config"
	contractutils "github.com/Cedro-Finance/CedroAVS/agentEnroller/contract/utils"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

func CallMapAgent(agent common.Address, chainRPC string, contractAddress common.Address, chainIDs []string, agentContracts []common.Address, privateKeyHex string) error {
	client, auth, err := contractutils.InitializeContractClient(chainRPC, privateKeyHex)
	if err != nil {
		return err
	}

	contractAbi := config.GetAgentManagerAbi()

	parsedAbi, err := abi.JSON(strings.NewReader(contractAbi))
	if err != nil {
		return fmt.Errorf("failed to parse ABI: %v", err)
	}

	chainIDBytes := make([][32]byte, len(chainIDs))
	for i, chainID := range chainIDs {
		var b [32]byte
		decodedChainID, err := hex.DecodeString(chainID[2:])
		if err != nil {
			return fmt.Errorf("failed to decode chain ID: %v", err)
		}
		copy(b[:], decodedChainID[:32])
		chainIDBytes[i] = b
	}

	// tx, err := bind.NewBoundContract(contractAddress, parsedAbi, client, client, client).Transact(auth, "mapAgent", agent, chainIDBytes, agentContracts)
	tx, err := bind.NewBoundContract(contractAddress, parsedAbi, client, client, client).Transact(auth, "mapAgent", agent, agentContracts[0])
	if err != nil {
		return fmt.Errorf("failed to transact: %v", err)
	}

	_, err = bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		return fmt.Errorf("failed to wait for transaction to be mined: %v", err)
	}

	fmt.Printf("Transaction complete: %s %s\n", chainRPC, tx.Hash())

	return nil
}

type ContractCallChainData struct {
	ChainName       string
	ContractAddress common.Address
}

func MakeManagerContractCalls(owner common.Address, addresses map[string]common.Address) {
	appConfig := config.LoadConfig()

	agentContracts := make([]common.Address, 0, len(addresses))
	chainIDs := make([]string, 0, len(addresses))

	for chainID, agentContractAddress := range addresses {
		agentContracts = append(agentContracts, agentContractAddress)
		chainIDs = append(chainIDs, chainID)
	}

	privateKeyHex := appConfig.WalletPrivateKey

	var wg sync.WaitGroup
	wg.Add(len(appConfig.Chains))

	for _, chainData := range appConfig.Chains {
		go func(chainData config.ChainConfig) {
			defer wg.Done()
			chainRPC := chainData.NodeUrl

			if chainRPC == "" {
				fmt.Printf("No rpc defined for chain: %s\n", chainData.Name)
				return
			}

			fmt.Println(chainData)

			err := CallMapAgent(owner, chainRPC, chainData.AgentManager, chainIDs, agentContracts, privateKeyHex)
			if err != nil {
				fmt.Printf("Failed to call mapAgent on chain %s: %v\n", chainData.Name, err)
			}
		}(chainData)
	}

	wg.Wait()

	fmt.Printf("Done Contract call")
}

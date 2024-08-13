package contractcaller

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"

	"sync"

	"github.com/Cedro-Finance/CedroAVS/agentEnroller/config"
	contractutils "github.com/Cedro-Finance/CedroAVS/agentEnroller/contract/utils"
	db "github.com/Cedro-Finance/CedroAVS/agentEnroller/database"
	"github.com/Cedro-Finance/CedroAVS/agentEnroller/utils"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

func CallMapAgent(agent common.Address, chainRPC string, contractAddress common.Address, chainIDs []string, agentContracts []common.Address, privateKeyHex string) (common.Hash, error) {
	client, auth, err := contractutils.InitializeContractClient(chainRPC, privateKeyHex)
	if err != nil {
		return common.Hash{}, err
	}

	contractAbi := config.GetAgentManagerAbi()

	parsedAbi, err := abi.JSON(strings.NewReader(contractAbi))
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to parse ABI: %v", err)
	}

	chainIDBytes := make([][32]byte, len(chainIDs))
	for i, chainID := range chainIDs {
		var b [32]byte
		decodedChainID, err := hex.DecodeString(chainID[2:])
		if err != nil {
			return common.Hash{}, fmt.Errorf("failed to decode chain ID: %v", err)
		}
		copy(b[:], decodedChainID[:32])
		chainIDBytes[i] = b
	}

	tx, err := bind.NewBoundContract(contractAddress, parsedAbi, client, client, client).Transact(auth, "mapAgent", agent, chainIDBytes, agentContracts)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to transact: %v", err)
	}

	_, err = bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		return tx.Hash(), fmt.Errorf("failed to wait for transaction to be mined: %v", err)
	}

	fmt.Printf("Transaction complete: %s %s\n", chainRPC, tx.Hash())

	return tx.Hash(), nil
}

type ContractCallChainData struct {
	ChainName       string
	ContractAddress common.Address
}

func MakeManagerContractCalls(owner common.Address, addresses map[string]common.Address) {
	appConfig := config.LoadConfig()

	agentContracts := make([]common.Address, 0, len(addresses))
	chainIDs := make([]string, 0, len(addresses))
	deployedChainNames := make([]string, 0, len(addresses))

	for chainID, agentContractAddress := range addresses {
		agentContracts = append(agentContracts, agentContractAddress)
		chainIDs = append(chainIDs, chainID)
		chainBytes, _ := hex.DecodeString(chainID[2:])
		var chainBytes32 [32]byte
		copy(chainBytes32[:], chainBytes[:32])
		chainName := utils.Bytes32ToString(chainBytes32)

		deployedChainNames = append(deployedChainNames, chainName)
	}

	privateKeyHex := appConfig.WalletPrivateKey
	var deploymentNotificationStatuses []db.DeploymentNotification

	var wg sync.WaitGroup
	wg.Add(len(appConfig.Chains))

	for _, chainData := range appConfig.Chains {
		go func(chainData config.ChainConfig) {
			defer wg.Done()
			chainRPC := chainData.NodeUrl

			chainDeploymentNotification := db.DeploymentNotification{
				Chain:          chainData.ShortName,
				ManagerAddress: chainData.AgentManager.String(),
			}

			if chainRPC == "" {
				chainDeploymentNotification.Status = "error"
				chainDeploymentNotification.Error = "No rpc defined for chain"

				fmt.Printf("[%v] No rpc defined for chain: %s\n", chainData.ShortName, chainData.Name)
				return
			}

			fmt.Println(chainData)

			txHash, err := CallMapAgent(owner, chainRPC, chainData.AgentManager, chainIDs, agentContracts, privateKeyHex)

			chainDeploymentNotification.TransactionHash = txHash.String()

			if err != nil {
				chainDeploymentNotification.Status = "error"
				chainDeploymentNotification.Error = err.Error()
				fmt.Printf("[%v] Failed to call mapAgent on chain %s: %v\n", chainData.ShortName, chainData.Name, err)
			} else {
				chainDeploymentNotification.Status = "success"
				chainDeploymentNotification.Error = ""
			}

			deploymentNotificationStatuses = append(deploymentNotificationStatuses, chainDeploymentNotification)
		}(chainData)
	}

	wg.Wait()

	db.ContractDeploymentInstance().BulkUpdateDeploymentNotification(owner, deploymentNotificationStatuses, deployedChainNames)

	fmt.Printf("Done Contract call")
}

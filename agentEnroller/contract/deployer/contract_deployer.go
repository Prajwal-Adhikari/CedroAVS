package contractdeployer

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Cedro-Finance/CedroAVS/agentEnroller/config"
	contractutils "github.com/Cedro-Finance/CedroAVS/agentEnroller/contract/utils"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func getAgentContractAddress(owner common.Address, managerAddress common.Address, client ethclient.Client) common.Address {
	contractAbi := config.GetAgentManagerAbi()

	parsedAbi, err := abi.JSON(strings.NewReader(contractAbi))
	if err != nil {
		fmt.Printf("Failed to pack method data: %v", err)
	}
	data, err := parsedAbi.Pack("agentContract", owner)
	if err != nil {
		fmt.Printf("Failed to pack method data: %v", err)
	}

	msg := ethereum.CallMsg{
		To:   &managerAddress,
		Data: data,
	}

	output, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		log.Fatalf("Failed to call contract method: %v", err)
	}

	var agentContractAddress common.Address
	err = parsedAbi.UnpackIntoInterface(&agentContractAddress, "agentContract", output)
	if err != nil {
		log.Fatalf("Failed to unpack the returned data: %v", err)
	}

	fmt.Printf("The agentContract for the address %s is: %s\n", owner.Hex(), agentContractAddress.Hex())

	return agentContractAddress
}

func DeployContract(owner common.Address, chainData config.ChainConfig) (common.Address, common.Hash, error) {
	appConfig := config.LoadConfig()

	client, auth, err := contractutils.InitializeContractClient(chainData.NodeUrl, appConfig.WalletPrivateKey)

	if err != nil {
		return common.Address{}, common.Hash{}, err
	}

	contractAbi, contractBytecode := config.GetAgentPodAbi()

	parsed, err := abi.JSON(strings.NewReader(contractAbi))
	if err != nil {
		return common.Address{}, common.Hash{}, fmt.Errorf("[%v] failed to parse abi: %w", chainData.ShortName, err)
	}

	agentPodAddress := getAgentContractAddress(owner, chainData.AgentManager, *client)

	if (agentPodAddress != common.Address{}) {
		return agentPodAddress, common.Hash{}, fmt.Errorf("agent already has contract deployed in chain")
	}

	// Deploy the contract
	address, tx, _, err := bind.DeployContract(auth, parsed, common.FromHex(contractBytecode), client, chainData.AgentManager, owner, chainData.Oracle)
	if err != nil {
		if tx != nil {
			return common.Address{}, tx.Hash(), fmt.Errorf("[%v] failed to deploy contract: %w", chainData.ShortName, err)
		}
		return common.Address{}, common.Hash{}, fmt.Errorf("[%v] failed to deploy contract: %w", chainData.ShortName, err)
	}

	fmt.Printf("[%v] Contract pending deploy: 0x%x\n", chainData.ShortName, address)
	fmt.Printf("[%v] Transaction waiting to be mined: 0x%x\n", chainData.ShortName, tx.Hash())

	// Wait for the transaction to be mined
	receipt, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		return common.Address{}, tx.Hash(), fmt.Errorf("failed to wait for the transaction to be mined: %w", err)
	}

	if receipt.Status == types.ReceiptStatusSuccessful {
		fmt.Printf("[%v] Contract deployed at address: 0x%x\n", chainData.ShortName, address)
		return address, tx.Hash(), nil
	}

	return common.Address{}, tx.Hash(), fmt.Errorf("[%v] contract deployment failed: 0x%x", chainData.ShortName, tx.Hash())
}

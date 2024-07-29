package contractdeployer

import (
	"context"
	"fmt"
	"strings"

	"github.com/Cedro-Finance/CedroAVS/agentEnroller/config"
	contractutils "github.com/Cedro-Finance/CedroAVS/agentEnroller/contract/utils"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func DeployContract(owner common.Address, chainData config.ChainConfig) (common.Address, error) {
	appConfig := config.LoadConfig()

	client, auth, err := contractutils.InitializeContractClient(chainData.NodeUrl, appConfig.WalletPrivateKey)

	if err != nil {
		return common.Address{}, err
	}

	contractAbi, contractBytecode := config.GetAgentPodAbi()

	parsed, err := abi.JSON(strings.NewReader(contractAbi))
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to parse abi: %w", err)
	}

	// Deploy the contract
	address, tx, _, err := bind.DeployContract(auth, parsed, common.FromHex(contractBytecode), client, chainData.AgentManager, owner, chainData.Oracle)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to deploy contract: %w", err)
	}

	fmt.Printf("Contract pending deploy: 0x%x\n", address)
	fmt.Printf("Transaction waiting to be mined: 0x%x\n", tx.Hash())

	// Wait for the transaction to be mined
	receipt, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to wait for the transaction to be mined: %w", err)
	}

	if receipt.Status == types.ReceiptStatusSuccessful {
		fmt.Printf("Contract deployed at address: 0x%x\n", address)
		return address, nil
	}

	return common.Address{}, fmt.Errorf("contract deployment failed: 0x%x", tx.Hash())
}

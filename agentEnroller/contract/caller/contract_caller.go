package contractcaller

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"sync"

	"github.com/Cedro-Finance/CedroAVS/agentEnroller/config"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// const contractAbi = `[
//
//	  {
//	    "inputs": [
//	      {
//	        "internalType": "address",
//	        "name": "_agent",
//	        "type": "address"
//	      },
//	      {
//	        "internalType": "bytes32[]",
//	        "name": "_chainIDs",
//	        "type": "bytes32[]"
//	      },
//	      {
//	        "internalType": "address[]",
//	        "name": "_agentContract",
//	        "type": "address[]"
//	      }
//	    ],
//	    "name": "mapAgent",
//	    "outputs": [],
//	    "stateMutability": "nonpayable",
//	    "type": "function"
//	  }
//	]`
const contractAbi = `[
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "_agent",
          "type": "address"
        },
        {
          "internalType": "address",
          "name": "_agentContract",
          "type": "address"
        }
      ],
      "name": "mapAgent",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    }
  ]`

func CallMapAgent(agent common.Address, chainRPC string, contractAddress common.Address, chainIDs []string, agentContracts []common.Address, privateKeyHex string) error {
	client, err := ethclient.Dial(chainRPC)
	if err != nil {
		return fmt.Errorf("failed to connect to Ethereum client: %v", err)
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return fmt.Errorf("failed to load private key: %v", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return errors.New("failed to cast public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return fmt.Errorf("failed to get nonce: %v", err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get gas price: %v", err)
	}

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get chain ID: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return fmt.Errorf("failed to create transactor: %v", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(3000000) // in wei
	auth.GasPrice = gasPrice

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

func GetChainData(chainName string, chains []config.ChainConfig) config.ChainConfig {
	var chainData config.ChainConfig

	for _, chain := range chains {
		if chain.ShortName == chainName {
			chainData = chain
			break
		}
	}

	return chainData
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

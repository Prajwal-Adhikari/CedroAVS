package contractutils

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/Cedro-Finance/CedroAVS/agentEnroller/config"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func InitializeContractClient(chainRPC string, privateKeyHex string) (*ethclient.Client, *bind.TransactOpts, error) {
	client, err := ethclient.Dial(chainRPC)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create transactor: %v", err)
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create transactor: %v", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, nil, fmt.Errorf("failed to create transactor: %v", err)
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create transactor: %v", err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create transactor: %v", err)
	}

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create transactor: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create transactor: %v", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      //in wei
	auth.GasLimit = uint64(3000000) //in wei
	auth.GasPrice = gasPrice
	return client, auth, nil
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

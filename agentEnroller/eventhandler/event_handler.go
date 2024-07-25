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
	"github.com/Cedro-Finance/CedroAVS/agentEnroller/utils"
	"github.com/ethereum/go-ethereum/common"
)

func EnrollAgentHandler(owner common.Address, chainIDS []string) {
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

			chainData := contractcaller.GetChainData(chainName, supportedChains)

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

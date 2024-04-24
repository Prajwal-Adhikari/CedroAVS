package main

import (
	"fmt"
	"sync"

	"github.com/Cedro-Finance/CedroAVS/verifierOperator/client/evm"
	"github.com/Cedro-Finance/CedroAVS/verifierOperator/config"
	"github.com/ethereum/go-ethereum/common"
)

func main() {
	configFile := config.LoadConfig("config.yaml")
	_ = configFile

	wg := sync.WaitGroup{}
	for _, chainConfig := range configFile.Chains {
		wg.Add(1)
		go func(cfg config.ChainConfig) {
			defer wg.Done()
			cl := evm.RegisteredChains[cfg.Name](cfg)
			cl.UpdateBalance(0, common.Address{})
			// wait for finality
			balance := cl.Balance()
			fmt.Println("balance is", balance)
		}(chainConfig)
	}
	wg.Wait()
}

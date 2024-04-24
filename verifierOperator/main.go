package main

import (
	"context"
	// "fmt"
	"math/big"
	"sync"

	"github.com/Cedro-Finance/CedroAVS/verifierOperator/client/avs"
	"github.com/Cedro-Finance/CedroAVS/verifierOperator/client/evm"
	"github.com/Cedro-Finance/CedroAVS/verifierOperator/config"
	"github.com/Cedro-Finance/CedroAVS/verifierOperator/logger"
	"github.com/ethereum/go-ethereum/common"
)

func main() {
	configFile := config.LoadConfig("config.yaml")
	_ = configFile

	avsClient := avs.NewAVSClient(&configFile.AVS)
	ctx := context.Background()
	reqCh := make(chan *big.Int)

	logger.InitLogging(configFile.Mode, configFile.LogConfig)
	logger.GetLogger().Info("Operator started")

	go avsClient.PollEvents(ctx, reqCh)

	for {
		select {
		case <-reqCh:
			wg := sync.WaitGroup{}
			for _, chainConfig := range configFile.Chains {
				wg.Add(1)
				go func(cfg config.ChainConfig) {
					defer wg.Done()
					cl := evm.RegisteredChains[cfg.Name](cfg)
					cl.UpdateBalance(0, common.Address{})
					// wait for finality
					// balance := cl.Balance()
					// fmt.Println("balance is", balance)
				}(chainConfig)
			}
			wg.Wait()
			return
		}
	}
}

package main

import (
	"fmt"
	"sync"

	"github.com/Cedro-Finance/CedroAVS/verifierOperator/client/evm"
	"github.com/Cedro-Finance/CedroAVS/verifierOperator/common"
	"github.com/Cedro-Finance/CedroAVS/verifierOperator/config"
)

func main() {
	configFile := config.LoadConfig("config.yaml")
	_ = configFile

	wg := sync.WaitGroup{}
	for _, v := range evm.RegisteredChains {
		wg.Add(1)
		go func(v func() common.IClient) {
			defer wg.Done()
			cl := v()
			balance := cl.Balance()
			fmt.Println("balance is", balance)
		}(v)
	}
	wg.Wait()
}

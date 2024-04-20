package main

import (
	"fmt"
	"sync"

	client "github.com/Cedro-Finance/CedroAVS/verifierOperator/client/evm"
	"github.com/Cedro-Finance/CedroAVS/verifierOperator/config"
)

func main() {
	configFile := config.LoadConfig("config.yaml")
	fmt.Println("config is ", configFile)
	wg := sync.WaitGroup{}
	for _, v := range client.RegisteredChains {
		wg.Add(1)
		go func(v func()) {
			defer wg.Done()
			v()
		}(v)
	}
	wg.Wait()
}

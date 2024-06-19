package main

import (
	"context"
	"fmt"

	// "fmt"
	"github.com/Cedro-Finance/CedroAVS/verifier/config"
	"github.com/Cedro-Finance/CedroAVS/verifier/operator"

	_ "github.com/Cedro-Finance/CedroAVS/verifier/client/agents/evm"
)

func main() {
	configFile, err := config.LoadConfig("config.yaml")

	if err != nil {
		fmt.Print(err)
		panic("Failed loading config")
	}
	_ = configFile

	ctx := context.Background()

	operator.StartRelay(ctx, configFile)
}

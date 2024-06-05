package main

import (
	"context"
	// "fmt"
	"github.com/Cedro-Finance/CedroAVS/verifierOperator/config"
	"github.com/Cedro-Finance/CedroAVS/verifierOperator/operator"
)

func main() {
	configFile, err := config.LoadConfig("config.yaml")

	if err != nil {
		panic("Failed loading config")
	}
	_ = configFile

	ctx := context.Background()

	operator.StartRelay(ctx, configFile)
}

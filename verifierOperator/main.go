package main

import (
	"context"
	"fmt"
	// "fmt"
	"github.com/Cedro-Finance/CedroAVS/verifierOperator/config"
	"github.com/Cedro-Finance/CedroAVS/verifierOperator/operator"
)

func main() {
	configFile, err := config.LoadConfig("verifierOperator/config.yaml")

	if err != nil {
		fmt.Print(err)
		panic("Failed loading config",)
	}
	_ = configFile

	ctx := context.Background()

	operator.StartRelay(ctx, configFile)
}

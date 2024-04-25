package main

import (
	"context"
	// "fmt"
	"github.com/Cedro-Finance/CedroAVS/verifierOperator/config"
	"github.com/Cedro-Finance/CedroAVS/verifierOperator/logger"
	"github.com/Cedro-Finance/CedroAVS/verifierOperator/relay"
)

func main() {
	configFile := config.LoadConfig("config.yaml")
	_ = configFile
	
	ctx := context.Background()

	logger.InitLogging(configFile.Mode, configFile.LogConfig)
	logger.GetLogger().Info("Operator started")

	relay.StartRelay(ctx, configFile)
}

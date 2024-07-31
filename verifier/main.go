package main

import (
	"context"
	"flag"
	"fmt"

	// "fmt"
	"github.com/Cedro-Finance/CedroAVS/core/config"
	"github.com/Cedro-Finance/CedroAVS/verifier/operator"

	_ "github.com/Cedro-Finance/CedroAVS/verifier/client/agents/evm"
)

var (
	configFile string
)

type FlagArgs struct {
	ConfigFile string
}

func init() {
	flag.StringVar(&configFile, "config", "", "config file")
}

func main() {

	flag.Parse()

	flagArgs := &FlagArgs{
		ConfigFile: configFile,
	}

	configFile, err := config.LoadConfig(flagArgs.ConfigFile)

	if err != nil {
		fmt.Print(err)
		panic("Failed loading config")
	}
	_ = configFile

	ctx := context.Background()

	operator.Start(ctx, configFile)
}

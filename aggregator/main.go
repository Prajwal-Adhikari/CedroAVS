package main

import (
	"context"
	"flag"
	"log"
	"sync"
	"github.com/Cedro-Finance/CedroAVS/aggregator/aggregator"
	"github.com/Cedro-Finance/CedroAVS/core/config"
)

// "context"

// env "github.com/Cedro-Finance/CedroAVS/aggregator/helpers"
// route "github.com/Cedro-Finance/CedroAVS/aggregator/routes"
// db "github.com/Cedro-Finance/CedroAVS/aggregator/services"
// "github.com/gin-gonic/gin"

type App struct{}

var (
	configFile string
)

func init() {
	flag.StringVar(&configFile, "config", "", "config file")
}

func main() {

	flag.Parse()

	flagArgs := &config.FlagArgs{
		ConfigFile: configFile,
	}

	config, err := config.LoadConfig(flagArgs.ConfigFile)

	if err != nil {
		log.Fatalf("Failed Loading Config %v", err)

	}

	_ = config

	aggr, err := aggregator.NewAggregator(config)

	if err != nil {
		panic("Counldn't start aggregator")
	}

	// ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	// defer stop()

	wg := &sync.WaitGroup{}
	wg.Add(1)

	err = aggr.Start(context.Background(), wg)
	if err != nil {
		log.Fatalln("Aggregator run failed", "err", err)
	}

	wg.Wait()
}

// Start Polling for the Database some reqId & if the threshold is reached, aggregate the signatueres & send it to the chain

package main

import (
	"context"
	"flag"
	"log"

	aggregator "github.com/Cedro-Finance/CedroAVS/aggregator/aggregation"
	"github.com/Cedro-Finance/CedroAVS/core/config"
)

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

	aggr, err := aggregator.NewAggregator(config)

	if err != nil {
		panic("Counldn't start aggregator")
	}

	// ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	// defer stop()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = aggr.Start(ctx)
	if err != nil {
		log.Fatalln("Aggregator run failed", "err", err)
	}

	<-ctx.Done()
}

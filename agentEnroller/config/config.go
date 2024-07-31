package config

import (
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"gopkg.in/yaml.v2"
)

type ChainConfig struct {
	Name          string         `yaml:"name"`
	Chain         string         `yaml:"chain"`
	ShortName     string         `yaml:"short_name"`
	ChainNameByte string         `yaml:"chain_name_byte"`
	NodeUrl       string         `yaml:"node_url"`
	NodeWSURL     string         `yaml:"node_wss_url"`
	AgentManager  common.Address `yaml:"agent_manager"`
	Oracle        common.Address `yaml:"oracle"`
}

type ConfigRaw struct {
	Chains        []ChainConfig `yaml:"chains"`
	RootChainName string        `yaml:"root_chain_name"`
}

type Config struct {
	Chains           []ChainConfig
	RootChainData    ChainConfig
	WalletPrivateKey string
}

var config *Config

func LoadConfig() *Config {
	if config != nil {
		return config
	}

	configData, err := os.ReadFile("config.yaml")

	if err != nil {
		log.Fatalf("Error loading config  %v", err)
	}

	var configRaw ConfigRaw

	err = yaml.Unmarshal(configData, &configRaw)

	if err != nil {
		log.Fatalf("Could not unmarshal config, %v", err)
	}

	var rootChainData ChainConfig

	for _, chain := range configRaw.Chains {
		if chain.ShortName == configRaw.RootChainName {
			rootChainData = chain
			break
		}
	}

	if (rootChainData == ChainConfig{}) {
		log.Fatalf("Root chain doesn't exist or is not defined")
	}

	config = &Config{
		Chains:           configRaw.Chains,
		RootChainData:    rootChainData,
		WalletPrivateKey: os.Getenv(("WALLET_PRIVATE_KEY")),
	}

	return config
}

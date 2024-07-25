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
	AgentManager  common.Address `yaml:"agent_manager"`
	Oracle        common.Address `yaml:"oracle"`
}

type ConfigRaw struct {
	Chains           []ChainConfig `yaml:"chains"`
	WalletPrivateKey string        `yaml:"wallet_private_key"`
}

type Config struct {
	Chains           []ChainConfig
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

	config = &Config{
		Chains:           configRaw.Chains,
		WalletPrivateKey: os.Getenv(("WALLET_PRIVATE_KEY")),
	}

	return config
}

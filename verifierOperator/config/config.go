package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	Development = "dev"
	Production  = "prod"
)

type LoggerConfig struct {
	Encoding   string `yaml:"encoding"`
	OutputDir  string `yaml:"output_dir"`
	OutputPath string `yaml:"-"` // calculated based on OutputDir
}

type AVSConfig struct {
	Chain          string `yaml:"chain"`
	ChainType      string `yaml:"chain_type"`
	ContactAddress string `yaml:"contract_address"`
	NodeUrl        string `yaml:"node_url"`
	StartBlock     uint64 `yaml:"start_block"`
	BlockFinality  uint64 `yaml:"block_finality"`
	BlockTime      time.Duration `yaml:"block_time"` //in second
}

type Config struct {
	Chains    []ChainConfig `yaml:"chains"`
	AVS       AVSConfig     `yaml:"avs"`
	LogConfig *LoggerConfig `yaml:"log"`
	Mode      string        `yaml:"mode"`
	KeyStore  *KeyStore     `yaml:"keystore"`
}

type ChainConfig struct {
	Name           string        `yaml:"name"`
	Chain          string        `yaml:"chain"`
	ContactAddress string        `yaml:"contract_address"`
	NodeUrl        string        `yaml:"node_url"`
	BlockFinality  uint64        `yaml:"block_finality"`
	BlockTime      time.Duration `yaml:"block_time"` //in second
}

var config *Config

type KeyStore struct {
	Path     string `yaml:"path"`
	Password string `yaml:"password"`
}

func GetConfig() *Config {
	return config
}

func LoadConfig(path string) *Config {
	file, err := os.ReadFile(path)
	if err != nil {
		panic("error in loading config")
	}

	config = new(Config)

	err = yaml.Unmarshal(file, config)

	if err != nil {
		panic("could not unmarshal config")
	}
	
	return config
}

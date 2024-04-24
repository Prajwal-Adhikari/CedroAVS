package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Chains   []ChainConfig `yaml:"chains"`
	KeyStore *KeyStore      `yaml:"keystore"`
}

type ChainConfig struct {
	Name           string `yaml:"name"`
	Chain          string `yaml:"chain"`
	ContactAddress string `yaml:"contract_address"`
}

type KeyStore struct {
	Path     string `yaml:"path"`
	Password string `yaml:"password"`
}

func LoadConfig(path string) *Config {
	file, err := os.ReadFile(path)
	if err != nil {
		panic("error in loading config")
	}

	config := &Config{}

	err = yaml.Unmarshal(file, config)
	if err != nil {
		panic("could not unmarshal config")
	}

	return config
}

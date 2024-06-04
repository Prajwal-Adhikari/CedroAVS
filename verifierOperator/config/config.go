package config

import (
	// "crypto/ecdsa"
	"fmt"
	"os"

	"github.com/Layr-Labs/eigensdk-go/chainio/clients/eth"
	// "github.com/Layr-Labs/eigensdk-go/chainio/txmgr"
	"github.com/Layr-Labs/eigensdk-go/crypto/bls"
	// sdklogging "github.com/Layr-Labs/eigensdk-go/logging"
	// "github.com/Layr-Labs/eigensdk-go/types"
	"github.com/ethereum/go-ethereum/common"
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

type AVSConfigRaw struct {
	RegistryCoordinatorAddr    string `yaml:"registryCoordinator"`
	OperatorStateRetrieverAddr string `yaml:"operatorStateRetriever"`
	EthRpcUrl                  string `yaml:"eth_rpc_url"`
	StartHeight                uint64 `yaml:"start_height"`
	BlockFinality              uint64 `yaml:"block_finality"`
}

type ChainConfig struct {
	Name           string `yaml:"name"`
	Chain          string `yaml:"chain"`
	ContactAddress string `yaml:"contract_address"`
	NodeUrl        string `yaml:"node_url"`
	BlockFinality  uint64 `yaml:"block_finality"`
	BlockTime      uint64 `yaml:"block_time"` //in second
}

type ConfigRaw struct {
	Chains                 []ChainConfig `yaml:"chains"`
	AVS                    AVSConfigRaw  `yaml:"avs"`
	LogConfig              *LoggerConfig `yaml:"log"`
	Mode                   string        `yaml:"mode"`
	KeyStore               *KeyStore     `yaml:"keystore"`
	BlsPrivateKeyStorePath string        `yaml:"bls_private_key_store_path"`
}

type AVSConfig struct {
	RegistryCoordinatorAddr    common.Address
	OperatorStateRetrieverAddr common.Address
	EthRpcUrl                  string
	StartHeight                uint64
	BlockFinality              uint64
}

type Config struct {
	BlsPrivateKey *bls.KeyPair
	EthClient     eth.Client
	Production    bool
	AVS           AVSConfig

	Agents []ChainConfig
}

var config *Config

type KeyStore struct {
	Path     string `yaml:"path"`
	Password string `yaml:"password"`
}

func GetConfig() *Config {
	return config
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		panic("error in loading config")
	}

	configRaw := new(ConfigRaw)

	err = yaml.Unmarshal(file, configRaw)

	if err != nil {
		fmt.Print(err)
		panic("could not unmarshal config")
	}

	blsKeypair, err := bls.NewKeyPairFromString(configRaw.BlsPrivateKeyStorePath)

	if err != nil {
		return nil, err
	}

	ethRpcUrl, ok := os.LookupEnv("ETH_RPC_URL")

	if ok && ethRpcUrl != "" {
		configRaw.AVS.EthRpcUrl = ethRpcUrl
	}

	avsRegistryCoordinatorAddress, rcOk := os.LookupEnv("AVS_REGISTRY_COORDINATOR_ADDRESS")
	operatorStateRetrieverAddress, osOk := os.LookupEnv("OPERATOR_STATE_RETRIEVER_ADDRESS")

	if rcOk && osOk && avsRegistryCoordinatorAddress != "" && operatorStateRetrieverAddress != "" {
		configRaw.AVS.OperatorStateRetrieverAddr = operatorStateRetrieverAddress
		configRaw.AVS.RegistryCoordinatorAddr = avsRegistryCoordinatorAddress
	}

	return &Config{
		AVS: AVSConfig{
			RegistryCoordinatorAddr:    common.HexToAddress(configRaw.AVS.RegistryCoordinatorAddr),
			OperatorStateRetrieverAddr: common.HexToAddress(configRaw.AVS.OperatorStateRetrieverAddr),
			EthRpcUrl:                  configRaw.AVS.EthRpcUrl,
			StartHeight:                configRaw.AVS.StartHeight,
			BlockFinality:              configRaw.AVS.BlockFinality,
		},
		BlsPrivateKey: blsKeypair,
		Agents:        configRaw.Chains,
	}, nil
}

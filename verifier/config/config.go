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

type EthConfigRaw struct {
	RpcUrl        string `yaml:"rpc_url"`
	StartHeight   uint64 `yaml:"start_height"`
	BlockFinality uint64 `yaml:"block_finality"`
	BlockTime     uint64 `yaml:"block_time"`
}

type AVSConfigRaw struct {
	RegistryCoordinatorAddr    string `yaml:"registryCoordinator"`
	OperatorStateRetrieverAddr string `yaml:"operatorStateRetriever"`
}

type ChainConfig struct {
	Name           string `yaml:"name"`
	Chain          string `yaml:"chain"`
	ChainId        string `yaml:"chain_id"`
	ContactAddress string `yaml:"contract_address"`
	NodeUrl        string `yaml:"node_url"`
	BlockFinality  uint64 `yaml:"block_finality"`
	BlockTime      uint64 `yaml:"block_time"` //in second
	GraphqlUrl     string `yaml:"graphql_url"`
}

type ConfigRaw struct {
	Chains                   []ChainConfig `yaml:"chains"`
	AVS                      AVSConfigRaw  `yaml:"avs"`
	EthConfig                EthConfigRaw  `yaml:"ethereum"`
	LogConfig                *LoggerConfig `yaml:"log"`
	Mode                     string        `yaml:"mode"`
	KeyStore                 *KeyStore     `yaml:"keystore"`
	GRpcAggregatorIPAddr     string        `yaml:"gRpc_aggregator_addr"`
	AgentManager             string        `yaml:"agent_manager"`
	AgentManagerStartHeight  uint64        `yaml:"agent_manager_start_height"`
	AgentManagerPollDuration uint64        `yaml:"agent_manager_poll_duration"`
	GraphAPIURL              string        `yaml:"graphql_url"`
}

type AVSConfig struct {
	RegistryCoordinatorAddr    common.Address
	OperatorStateRetrieverAddr common.Address
}

type EthConfig struct {
	RpcUrl        string
	StartHeight   uint64
	BlockFinality uint64
	BlockTime     uint64
	EthClient     eth.Client
}

type Config struct {
	BlsPrivateKey            *bls.KeyPair
	Production               bool
	AVS                      AVSConfig
	KeyStore                 *KeyStore
	GRpcAggregatorIpAddr     string
	Agents                   []ChainConfig
	AgentManagerAddr         common.Address
	EthConfig                EthConfig
	AgentManagerStartHeight  uint64
	AgentManagerPollDuration uint64
	GraphQLUrl               string
}

var config *Config

type KeyStore struct {
	Path     string `yaml:"path"`
	Password string `yaml:"password"`

	BlsPath     string `yaml:"bls_path"`
	BlsPassword string `yaml:"bls_password"`
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


	blsKeypair, err := bls.ReadPrivateKeyFromFile(configRaw.KeyStore.BlsPath, configRaw.KeyStore.BlsPassword)

	if err != nil {
		return nil, err
	}

	ethRpcUrl, ok := os.LookupEnv("ETH_RPC_URL")

	if ok && ethRpcUrl != "" {
		configRaw.EthConfig.RpcUrl = ethRpcUrl
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
		},
		BlsPrivateKey:            blsKeypair,
		Agents:                   configRaw.Chains,
		GRpcAggregatorIpAddr:     configRaw.GRpcAggregatorIPAddr,
		KeyStore:                 configRaw.KeyStore,
		AgentManagerAddr:         common.HexToAddress(configRaw.AgentManager),
		AgentManagerStartHeight:  configRaw.AgentManagerStartHeight,
		AgentManagerPollDuration: configRaw.AgentManagerPollDuration,
		GraphQLUrl:               configRaw.GraphAPIURL,
		EthConfig: EthConfig{
			RpcUrl:        configRaw.EthConfig.RpcUrl,
			BlockFinality: configRaw.EthConfig.BlockFinality,
			BlockTime:     configRaw.EthConfig.BlockTime,
		},
	}, nil
}

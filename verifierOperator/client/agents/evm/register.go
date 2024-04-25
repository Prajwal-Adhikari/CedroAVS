package evm

import (
	"github.com/Cedro-Finance/CedroAVS/verifierOperator/common"
	"github.com/Cedro-Finance/CedroAVS/verifierOperator/config"
)

var (
	RegisteredChains map[string]func(config.ChainConfig) common.IClient
)

func init() {
	RegisteredChains = make(map[string]func(config config.ChainConfig) common.IClient)
	RegisteredChains["ethereum"] = NewClient
	RegisteredChains["avalanche"] = NewClient
}

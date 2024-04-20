package evm

import (
	"github.com/Cedro-Finance/CedroAVS/verifierOperator/common"
)

var (
	RegisteredChains map[string]func() common.IClient
)

func init() {
	RegisteredChains = make(map[string]func() common.IClient)
	RegisteredChains["eth"] = NewClient
	RegisteredChains["avalanche"] = NewClient
}

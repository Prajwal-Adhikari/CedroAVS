package evm

import (
	"github.com/Cedro-Finance/CedroAVS/verifierOperator/relay"
)

func init() {
	relay.RegisteredClients["ethereum"] = NewClient
	relay.RegisteredClients["avalanche"] = NewClient
}

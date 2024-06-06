package evm

import "github.com/Cedro-Finance/CedroAVS/verifierOperator/operator"

func init() {
	operator.RegisteredClients["ethereum"] = NewClient
	operator.RegisteredClients["avalanche"] = NewClient
}

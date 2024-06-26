package evm

import (
	"github.com/Cedro-Finance/CedroAVS/verifier/operator"
)

func init() {
	operator.RegisteredClients["ETH"] = NewClient
	operator.RegisteredClients["BNB"] = NewClient
	operator.RegisteredClients["SEP"] = NewClient
	// operator.RegisteredClients["binance"] = NewClient
	// operator.RegisteredClients["sepolia"] = NewClient
}

package common

import (
	"context"
	"math/big"

	"github.com/Cedro-Finance/CedroAVS/verifier/client/avs"
	"github.com/ethereum/go-ethereum/common"
)

type IAgentPodClient interface {
	InitializeAgentPod(addr common.Address)
	GetBalance(context.Context, *avs.Request, chan<- *big.Int)
}

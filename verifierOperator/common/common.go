package common

import (
	"context"
	"math/big"

	"github.com/Cedro-Finance/CedroAVS/verifierOperator/client/agents"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type IClient interface {
	Balance(context.Context, chan *agents.BalanceResponse)
	UpdateBalance(context.Context, *big.Int, common.Address) *types.Transaction
}
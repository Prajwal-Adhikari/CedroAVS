package common

import (
	"math/big"

	"github.com/Cedro-Finance/CedroAVS/verifierOperator/client/avs"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type IClient interface {
	Balance(*avs.Request, chan<- *big.Int) *big.Int
	UpdateBalance(uint64, common.Address) *types.Transaction
}
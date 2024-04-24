package common

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type IClient interface {
	Balance() *big.Int
	UpdateBalance(uint64, common.Address)
}

package common

import (
	"math/big"
)

type IClient interface {
	Balance() *big.Int
}

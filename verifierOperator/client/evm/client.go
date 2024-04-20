package evm

import (
	"fmt"
	"math/big"

	"github.com/Cedro-Finance/CedroAVS/verifierOperator/common"
	ethCommon "github.com/ethereum/go-ethereum/common"
)

type Client struct {
	Chain    string
	Contract ethCommon.Address
}

func (c *Client) Balance() *big.Int {
	return big.NewInt(100)
}

func NewClient() common.IClient {
	fmt.Println("client is here")
	return &Client{}
}

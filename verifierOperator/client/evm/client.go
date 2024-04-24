package evm

import (
	"fmt"
	"math/big"

	"github.com/Cedro-Finance/CedroAVS/verifierOperator/client/avs"
	"github.com/Cedro-Finance/CedroAVS/verifierOperator/common"
	"github.com/Cedro-Finance/CedroAVS/verifierOperator/config"
	ethCommon "github.com/ethereum/go-ethereum/common"
)

type Client struct {
	name     string
	chain    string
	contract ethCommon.Address
}

func (c *Client) Balance(balanceRequest *avs.Request, balanceChan chan<-*big.Int) *big.Int {
	return ethCommon.Big0
}

// snapshots the balace of lp in the respective chain
// sendtx method
func (c *Client) UpdateBalance(rxId uint64, address ethCommon.Address) {}

func NewClient(config config.ChainConfig) common.IClient {
	fmt.Println("client is here")
	return &Client{
		name:     config.Name,
		chain:    config.Chain,
		contract: ethCommon.HexToAddress(config.ContactAddress),
	}
}

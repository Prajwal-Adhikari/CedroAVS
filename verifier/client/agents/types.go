package agents

import (
	"math/big"
	ethCommon "github.com/ethereum/go-ethereum/common"
)

type TaskRequest interface {
	GetReqHash() string
}

type DeviationRequest struct {
	ReqId                *big.Int
	ReferenceBlockNumber uint64
	AgentAddress         ethCommon.Address
	Timestamp            uint64
}

func (d *DeviationRequest) GetReqHash() string {
	return d.AgentAddress.Hex()
}

type AgentSlashRequest struct {
	TxnId                *big.Int
	ReferenceBlockNumber uint64
	AgentAddress         ethCommon.Address
	InsufficientBalance  *big.Int
}

func (a *AgentSlashRequest) GetReqHash() string {
	return a.AgentAddress.Hex()
}

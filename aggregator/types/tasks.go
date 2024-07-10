package types

import (
	"math/big"

	sdktypes "github.com/Layr-Labs/eigensdk-go/types"
	"github.com/ethereum/go-ethereum/common"
)

type ITaskInfo interface {}

type FinishedTaskStatus struct {
	Message          ITaskInfo
	TxHash           common.Hash
	BlockHash        common.Hash `json:"blockHash,omitempty"`
	BlockNumber      *big.Int    `json:"blockNumber,omitempty"`
	TransactionIndex uint        `json:"transactionIndex"`
}

type OperatorStatus struct {
	LastTime   int64               `json:"lastTime"`
	OperatorId sdktypes.OperatorId `json:"operatorId"`
}

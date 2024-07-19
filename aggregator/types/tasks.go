package types

import (
	"math/big"

	"github.com/Cedro-Finance/CedroAVS/core"
	sdktypes "github.com/Layr-Labs/eigensdk-go/types"
	"github.com/ethereum/go-ethereum/common"
)

type CreateTaskReq interface{}
type CreateTaskReqResp interface{}

type SignedTaskResRequest interface{}
type ITaskInfo interface {
	ToPbType() interface{}
	ToCedroServiceManagerHeader() interface{}

	EncodeSignHash() ([]byte, error)
	GetHeaderHash() (core.Bytes32, error)
}

type FinishedTaskStatus struct {
	Message          *core.BalanceTaskInfo
	TxHash           common.Hash
	BlockHash        common.Hash `json:"blockHash,omitempty"`
	BlockNumber      *big.Int    `json:"blockNumber,omitempty"`
	TransactionIndex uint        `json:"transactionIndex"`
}

type FinishedSlashTaskStatus struct{
	Message          *core.AgentSlashTaskInfo
	TxHash           common.Hash
	BlockHash        common.Hash `json:"blockHash,omitempty"`
	BlockNumber      *big.Int    `json:"blockNumber,omitempty"`
	TransactionIndex uint        `json:"transactionIndex"`
}

type OperatorStatus struct {
	LastTime   int64               `json:"lastTime"`
	OperatorId sdktypes.OperatorId `json:"operatorId"`
}

type ITaskProcessor interface {
	GetTaskByIndex(taskIndex sdktypes.TaskIndex) ITaskInfo
	NewIndex() sdktypes.TaskIndex
	GetTaskByRequestHash() ITaskInfo
	GetFinishedTaskByRequestHash() ITaskInfo
	SetFinishedTask(reqHash string, finished *FinishedTaskStatus)

	CreateTask(req CreateTaskReq) (CreateTaskReqResp, error)

	ProcessSignedTaskResponse(req SignedTaskResRequest) (*core.SignedTaskRespResponse, error)
}

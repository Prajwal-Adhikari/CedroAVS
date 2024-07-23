package types

import (
	"math/big"

	"github.com/Cedro-Finance/CedroAVS/core/tasks"
	sdktypes "github.com/Layr-Labs/eigensdk-go/types"
	"github.com/ethereum/go-ethereum/common"
)

type CreateTaskReq interface{}
type CreateTaskReqResp interface{}

type SignedTaskResRequest interface{}


type FinishedTaskStatus struct {
	Message          tasks.ITaskInfo
	TxHash           common.Hash
	BlockHash        common.Hash `json:"blockHash,omitempty"`
	BlockNumber      *big.Int    `json:"blockNumber,omitempty"`
	TransactionIndex uint        `json:"transactionIndex"`
}

type FinishedSlashTaskStatus struct{
	Message          *tasks.AgentSlashTaskInfo
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
	GetTaskByIndex(taskIndex sdktypes.TaskIndex) tasks.ITaskInfo
	NewIndex() sdktypes.TaskIndex
	GetTaskByRequestHash() tasks.ITaskInfo
	GetFinishedTaskByRequestHash() tasks.ITaskInfo
	SetFinishedTask(reqHash string, finished *FinishedTaskStatus)

	CreateTask(req CreateTaskReq) (CreateTaskReqResp, error)

	ProcessSignedTaskResponse(req SignedTaskResRequest) (*tasks.SignedTaskRespResponse, error)
}

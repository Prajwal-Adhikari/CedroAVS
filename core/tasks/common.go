package tasks

import (
	"fmt"
	"math/big"

	aggregator "github.com/Cedro-Finance/CedroAVS/api/grpc"
	"github.com/Layr-Labs/eigensdk-go/crypto/bls"
	sdktypes "github.com/Layr-Labs/eigensdk-go/types"
	"github.com/ethereum/go-ethereum/common"
)

const (
	AgentDeviationTask = "deviation"
	AgentSlashTask     = "agentslash"
)

var TaskTypes = [2]string{AgentDeviationTask, AgentSlashTask}

type SignedTaskRespResponse struct {
	Reply  bool
	TxHash [32]byte
}

func (r SignedTaskRespResponse) ToPbType() *aggregator.SignedTaskRespResponse {
	return &aggregator.SignedTaskRespResponse{
		Reply:  r.Reply,
		TxHash: r.TxHash[:],
	}
}

// The init operator request
type InitOperatorRequest struct {
	OperatorId                 sdktypes.OperatorId
	OperatorAddress            common.Address
	OperatorStateRetrieverAddr common.Address
	RegistryCoordinatorAddr    common.Address
}

func NewInitOperatorRequest(req *aggregator.InitOperatorRequest) (*InitOperatorRequest, error) {
	operatorId := req.GetOperatorId()
	if len(operatorId) != 32 {
		return nil, fmt.Errorf("operator ID len should be 32, got %d", len(operatorId))
	}

	if !common.IsHexAddress(req.GetOperatorAddress()) {
		return nil, fmt.Errorf("operatorAddress not a hex address")
	}
	operatorAddress := common.HexToAddress(req.OperatorAddress)

	if !common.IsHexAddress(req.GetOperatorStateRetrieverAddress()) {
		return nil, fmt.Errorf("operatorStateRetrieverAddr not a hex address")
	}
	operatorStateRetrieverAddr := common.HexToAddress(req.GetOperatorStateRetrieverAddress())

	if !common.IsHexAddress(req.GetRegistryCoordinatorAddress()) {
		return nil, fmt.Errorf("registryCoordinatorAddr not a hex address")
	}
	registryCoordinatorAddr := common.HexToAddress(req.GetRegistryCoordinatorAddress())

	res := &InitOperatorRequest{
		OperatorId:                 sdktypes.Bytes32(req.OperatorId),
		OperatorAddress:            operatorAddress,
		OperatorStateRetrieverAddr: operatorStateRetrieverAddr,
		RegistryCoordinatorAddr:    registryCoordinatorAddr,
	}

	copy(res.OperatorId[:], operatorId[:32])

	return res, nil
}

// The init operator response
type InitOperatorResponse struct {
	Ok  bool
	Res string
}

func (r InitOperatorResponse) ToPbType() *aggregator.InitOperatorResponse {
	return &aggregator.InitOperatorResponse{
		Ok:     r.Ok,
		Reason: r.Res,
	}
}

type ITaskInfo interface {
	// ToPbType() interface{}
	// ToCedroServiceManagerHeader() interface{}

	GetTaskIndex() sdktypes.TaskIndex
	EncodeSignHash() ([]byte, error)
	GetHeaderHash() (sdktypes.Bytes32, error)
	GetUniqueHash() string
}

type ICreateTaskRequest interface {
	GetHash() (string, error)
	GetAgent() common.Address
	GetIdentifier() *big.Int
	GetRefBlock() uint64
	GetTaskInfo(taskIndex sdktypes.TaskIndex, quorumNums sdktypes.QuorumNums, quorumThresholds sdktypes.QuorumThresholdPercentages) ITaskInfo
}

type ICreateTaskResponse struct {
	Info ITaskInfo
}

type CreateSignedTaskResRequest struct {
	Task         ITaskInfo
	BlsSignature bls.Signature
	OperatorId   sdktypes.OperatorId
}

type TaskInfo struct {
	QuorumNumbers              sdktypes.QuorumNums
	QuorumThresholdPercentages sdktypes.QuorumThresholdPercentages
	TaskIndex                  uint32
	ReferenceBlockNumber       uint64
}

package rpc

import (
	"github.com/Cedro-Finance/CedroAVS/aggregator/taskprocessor"
	"github.com/Cedro-Finance/CedroAVS/core/tasks"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

type AggregatorRpcHandler struct {
	logger logging.Logger

	taskprocessors map[string]*taskprocessor.TaskProcessor
}

func NewAggregatorRpcHandler(logger logging.Logger, tps map[string]*taskprocessor.TaskProcessor) (*AggregatorRpcHandler, error) {
	return &AggregatorRpcHandler{
		logger:         logger,
		taskprocessors: tps,
	}, nil
}

func (ag *AggregatorRpcHandler) InitOperator(req *tasks.InitOperatorRequest) (*tasks.InitOperatorResponse, error) {
	reply := &tasks.InitOperatorResponse{
		Ok: false,
	}

	// if ags.cfg.AVS.OperatorStateRetrieverAddr.Hex() != req.GetOperatorStateRetrieverAddr() {
	// 	return reply, nil
	// }

	reply.Ok = true

	return reply, nil
}

func (ag *AggregatorRpcHandler) CreateTask(task string, req tasks.ICreateTaskRequest) (tasks.ICreateTaskResponse, error) {
	res, err := ag.taskprocessors[task].CreateTask(req)

	return res, err
}

func (ag *AggregatorRpcHandler) ProcessSignedTask(task string, signed *tasks.CreateSignedTaskResRequest) (*tasks.SignedTaskRespResponse, error) {

	res, err := ag.taskprocessors[task].ProcessSignedTaskResponse(signed)

	return res, err
}

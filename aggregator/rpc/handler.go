package rpc

import (
	"fmt"

	"github.com/Cedro-Finance/CedroAVS/aggregator/taskprocessor"
	"github.com/Cedro-Finance/CedroAVS/core"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

type AggregatorRpcHandler struct {
	logger logging.Logger

	deviator *taskprocessor.DeviationProcess
	slasher  *taskprocessor.SlashProcess
}

func NewAggregatorRpcHandler(logger logging.Logger, deviator *taskprocessor.DeviationProcess, slasher *taskprocessor.SlashProcess) (*AggregatorRpcHandler, error) {
	return &AggregatorRpcHandler{
		logger:   logger,
		deviator: deviator,
		slasher:  slasher,
	}, nil
}

func (ag *AggregatorRpcHandler) InitOperator(req *core.InitOperatorRequest) (*core.InitOperatorResponse, error) {
	reply := &core.InitOperatorResponse{
		Ok: false,
	}

	// if ags.cfg.AVS.OperatorStateRetrieverAddr.Hex() != req.GetOperatorStateRetrieverAddr() {
	// 	return reply, nil
	// }

	reply.Ok = true

	return reply, nil
}

func (ag *AggregatorRpcHandler) CreateTask(task string, req interface{}) (interface{}, error) {
	switch task {
	case "deviation":
		devreq := req.(*core.CreateDeviationTaskRequest)
		task, err := ag.deviator.CreateTask(devreq)

		return task, err

	case "agentslash":
		sreq := req.(*core.CreateAgentSlashTaskRequest)
		task, err := ag.slasher.CreateTask(sreq)

		return task, err

	default:
		return nil, fmt.Errorf("invalid task type")
	}
}

func (ag *AggregatorRpcHandler) ProcessSignedTask(task string, signed interface{}) (interface{}, error) {
	switch task {
	case "deviation":
		devsign := signed.(*core.SignedDeviationTaskRespRequest)
		resp, err := ag.deviator.ProcessSignedTaskResponse(devsign)

		return resp, err

	case "agentslash":
		ssign := signed.(*core.SignedAgentSlashTaskRespRequest)
		resp, err := ag.slasher.ProcessSignedTaskResponse(ssign)

		return resp, err

	default:
		return nil, fmt.Errorf("invalid task type")
	}
}

package taskprocessor

import (
	"context"
	"fmt"
	"math/big"

	blsaggregator "github.com/Cedro-Finance/CedroAVS/aggregator/bls"
	"github.com/Cedro-Finance/CedroAVS/aggregator/types"
	avs "github.com/Cedro-Finance/CedroAVS/core/avsClient"
	"github.com/Cedro-Finance/CedroAVS/core/config"
	"github.com/Cedro-Finance/CedroAVS/core/tasks"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum/common"
)

type DeviationProcess struct {
	taskProcessor *TaskProcessor
	logger        logging.Logger

	avsWriter *avs.AvsWriter
}

func NewDeviationProcessor(c *config.Config, logger logging.Logger, avsWriter *avs.AvsWriter) (*DeviationProcess, error) {
	aggService, err := NewTaskProcessorService(c, logger)

	if err != nil {
		panic(err)
	}

	return &DeviationProcess{
		// tasks:         make(map[sdktypes.TaskIndex]*tasks.BalanceTaskInfo),
		// finishedTasks: make(map[string]*types.FinishedTaskStatus),
		taskProcessor: aggService,

		logger:    logger,
		avsWriter: avsWriter,
	}, nil
}

func (d *DeviationProcess) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case res := <-d.taskProcessor.blsAggrService.GetResponseChannel():
			d.SendAggregatedResponseToContract(res)
		}
	}
}

func (d *DeviationProcess) CreateTask(req *tasks.CreateDeviationTaskRequest) (*tasks.CreateDeviationTaskResponse, error) {
	d.logger.Info("Received Create Task", "reqId", req.ReqId, "Agent", req.AgentAddress)

	reqHash, err := req.GetHash()

	if err != nil {
		d.logger.Errorf("Error while getting request Hash", "err", err)
		return nil, err
	}

	finishedTask := d.taskProcessor.GetFinishedTaskByRequestHash(reqHash)

	if finishedTask != nil {
		d.logger.Error("Task already finished", "reqId", req.ReqId.String(), "agent", req.AgentAddress.Hex())

		return nil, err
	}

	task := d.taskProcessor.GetTaskByRequestHash(reqHash)

	taskres := task.(*tasks.BalanceTaskInfo)

	if task == nil {
		taskIndex := d.taskProcessor.NewIndex()
		task, err = d.sendNewTask(taskIndex, req.ReferenceBlockNumber, req.ReqId, req.AgentAddress)

		if err != nil {
			d.logger.Error("send new task failed", "err", err)
			return nil, err
		}
	}
	return &tasks.CreateDeviationTaskResponse{Info: *taskres}, nil
}

func (d *DeviationProcess) ProcessSignedTaskResponse(req *tasks.SignedDeviationTaskRespRequest) (*tasks.SignedTaskRespResponse, error) {
	d.logger.Info("New Signature Received from operator", "signature", req.BlsSignature.Serialize(), "operatorId", req.OperatorId, "balance", req.Balance)

	reqHash, err := tasks.GetRequestHash(req.Balance.ReqId, req.Balance.AgentAddress)

	if err != nil {
		d.logger.Errorf("Error while creating request hash", "err", fmt.Sprintf("%#f", err))

		return nil, err
	}

	taskIndex := req.Balance.TaskIndex

	d.taskProcessor.tasksMu.Lock()
	d.taskProcessor.tasks[taskIndex] = &req.Balance
	d.taskProcessor.tasksMu.Unlock()

	taskResponseDigest, err := req.Balance.GetHeaderHash()

	if err != nil {
		d.logger.Error("Error getting taskResponseDigest", "err", err)
		return nil, err
	}

	task := d.taskProcessor.GetTaskByRequestHash(reqHash)

	if task == nil {
		d.logger.Errorf("Task doesn't exist for given request hash", "err", fmt.Sprintf("%#f", err))
		return nil, err
	}

	err = d.taskProcessor.blsAggrService.ProcessNewSignature(
		context.Background(), taskIndex, taskResponseDigest.UnderlyingType(),
		&req.BlsSignature, req.OperatorId,
	)

	if err != nil {
		d.logger.Error("ProcessNewSignature error", "err", err)
	}

	return &tasks.SignedTaskRespResponse{}, nil
}

func (d *DeviationProcess) sendNewTask(taskIndex uint32, refBlock uint64, reqId *big.Int, agentAddr common.Address) (*tasks.BalanceTaskInfo, error) {

	quorumDetails, err := d.taskProcessor.GetQuorumDetails(context.Background(), refBlock)

	if err != nil {
		d.logger.Error("error while getting quorum details")
		return nil, err
	}

	newTask := &tasks.BalanceTaskInfo{
		QuorumNumbers:              quorumDetails.QuorumNums,
		QuorumThresholdPercentages: quorumDetails.QuorumThresholdPercentages,
		TaskIndex:                  taskIndex,
		ReferenceBlockNumber:       refBlock,
		ReqId:                      reqId,
		ActualBalance:              big.NewInt(0),
		IdealBalance:               big.NewInt(0),
		AgentAddress:               agentAddr,
	}

	err = d.taskProcessor.InitializeNewTask(context.Background(), taskIndex, newTask, *quorumDetails, refBlock)

	if err != nil {
		d.logger.Error("InitializeNewTask failed", "err", err)
		return nil, err
	}

	return newTask, nil
}

func (d *DeviationProcess) SendAggregatedResponseToContract(blsResp blsaggregator.BLSAggregationServiceResponse) {

	nonSignerStakes, err := d.taskProcessor.ReceiveAggregatedSignatures(blsResp)

	if err != nil {
		d.logger.Error("error while building non signer stake signatures")
		return
	}

	task := d.taskProcessor.GetTaskByIndex(blsResp.TaskIndex)

	taskres := task.(*tasks.BalanceTaskInfo)

	res, err := d.avsWriter.SendConfirmationBatch(context.Background(), taskres, nonSignerStakes)

	if err != nil {
		d.logger.Error("aggregator failed to respond to task", "err", err)
		return
	}

	if res != nil {
		reqHash, err := tasks.GetRequestHash(taskres.ReqId, taskres.AgentAddress)

		if err != nil {
			d.logger.Error("agg failed to find req hash")
		}

		d.taskProcessor.SetFinishedTask(reqHash, &types.FinishedTaskStatus{
			Message:          taskres,
			TxHash:           res.TxHash,
			BlockHash:        res.BlockHash,
			BlockNumber:      res.BlockNumber,
			TransactionIndex: res.TransactionIndex,
		})
	} else {
		d.logger.Error("the send confirm aler threw error", "err", res)
	}

}

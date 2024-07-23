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

type SlashProcess struct {
	taskProcessor *TaskProcessor

	logger logging.Logger

	avsWriter *avs.AvsWriter
}

func NewSlashProcessor(c *config.Config, logger logging.Logger, avsWriter *avs.AvsWriter) (*SlashProcess, error) {

	aggService, err := NewTaskProcessorService(c, logger)

	if err != nil {
		panic(err)
	}

	return &SlashProcess{
		taskProcessor: aggService,
		logger:        logger,
		avsWriter:     avsWriter,
	}, nil
}

func (d *SlashProcess) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case res := <-d.taskProcessor.blsAggrService.GetResponseChannel():
			d.SendAggregatedResponseToContract(res)
		}
	}
}

func (d *SlashProcess) CreateTask(req *tasks.CreateAgentSlashTaskRequest) (*tasks.CreateAgentSlashTaskResponse, error) {
	// d.tasksMu.Lock()
	// defer d.tasksMu.Unlock()

	d.logger.Info("Received Create Task", "reqId", req.TxnId, "Agent", req.AgentAddress)

	reqHash, err := req.GetHash()

	if err != nil {
		d.logger.Errorf("Error while getting request Hash", "err", err)
		return nil, err
	}

	finishedTask := d.taskProcessor.GetFinishedTaskByRequestHash(reqHash)

	if finishedTask != nil {
		d.logger.Error("Task already finished", "reqId", req.TxnId.String(), "agent", req.AgentAddress.Hex())

		return nil, err
	}

	task := d.taskProcessor.GetTaskByRequestHash(reqHash)

	taskres := task.(*tasks.AgentSlashTaskInfo)

	if task == nil {
		taskIndex := d.taskProcessor.NewIndex()
		task, err = d.sendNewTask(taskIndex, req.ReferenceBlockNumber, req.TxnId, req.AgentAddress)

		if err != nil {
			d.logger.Error("send new task failed", "err", err)
			return nil, err
		}
	}
	return &tasks.CreateAgentSlashTaskResponse{Info: *taskres}, nil
}

func (d *SlashProcess) ProcessSignedTaskResponse(req *tasks.SignedAgentSlashTaskRespRequest) (*tasks.SignedTaskRespResponse, error) {
	d.logger.Info("New Signature Received from operator", "signature", req.BlsSignature.Serialize(), "operatorId", req.OperatorId, "balance", req.Balance)

	reqHash, err := tasks.GetRequestHash(req.Balance.TxnId, req.Balance.AgentAddress)

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

func (d *SlashProcess) sendNewTask(taskIndex uint32, refBlock uint64, txnId *big.Int, agentAddr common.Address) (*tasks.AgentSlashTaskInfo, error) {

	quorumDetails, err := d.taskProcessor.GetQuorumDetails(context.Background(), refBlock)

	if err != nil {
		d.logger.Error("error while getting quorum details")
		return nil, err
	}

	newTask := &tasks.AgentSlashTaskInfo{
		QuorumNumbers:              quorumDetails.QuorumNums,
		QuorumThresholdPercentages: quorumDetails.QuorumThresholdPercentages,
		TaskIndex:                  taskIndex,
		ReferenceBlockNumber:       refBlock,
		AgentAddress:               agentAddr,
		InsufficientBalance:        big.NewInt(0),
		TxnId:                      txnId,
	}

	err = d.taskProcessor.InitializeNewTask(context.Background(), taskIndex, newTask, *quorumDetails, refBlock)

	if err != nil {
		d.logger.Error("InitializeNewTask failed", "err", err)
		return nil, err
	}

	return newTask, nil
}

func (d *SlashProcess) SendAggregatedResponseToContract(blsResp blsaggregator.BLSAggregationServiceResponse) {

	nonSignerStakes, err := d.taskProcessor.ReceiveAggregatedSignatures(blsResp)

	if err != nil {
		d.logger.Error("error while building non signer stake signatures")
		return
	}

	task := d.taskProcessor.GetTaskByIndex(blsResp.TaskIndex)

	taskres := task.(*tasks.AgentSlashTaskInfo)

	res, err := d.avsWriter.SendConfirmationAgentSlash(context.Background(), taskres, nonSignerStakes)

	if err != nil {
		d.logger.Error("aggregator failed to respond to task", "err", err)
		return
	}

	if res != nil {
		reqHash, err := tasks.GetRequestHash(big.NewInt(0), taskres.AgentAddress)

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

package taskprocessor

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	blsaggregator "github.com/Cedro-Finance/CedroAVS/aggregator/bls"
	"github.com/Cedro-Finance/CedroAVS/aggregator/types"
	"github.com/Cedro-Finance/CedroAVS/core"
	avs "github.com/Cedro-Finance/CedroAVS/core/avsClient"
	"github.com/Cedro-Finance/CedroAVS/core/config"
	"github.com/Layr-Labs/eigensdk-go/logging"
	sdktypes "github.com/Layr-Labs/eigensdk-go/types"
	"github.com/ethereum/go-ethereum/common"
)

type DeviationProcess struct {
	tasks   map[sdktypes.TaskIndex]*core.BalanceTaskInfo
	tasksMu sync.RWMutex

	finishedTasks   map[string]*types.FinishedTaskStatus
	finishedTasksMu sync.RWMutex

	nextTaskIndex   sdktypes.TaskIndex
	nextTaskIndexMu sync.RWMutex

	service *AggregatorService
	logger  logging.Logger

	avsWriter *avs.AvsWriter
}

func NewDeviationProcessor(c *config.Config, logger logging.Logger, avsWriter *avs.AvsWriter) (*DeviationProcess, error) {

	aggService, err := NewAggregatorService(c, logger)

	if err != nil {
		panic(err)
	}

	return &DeviationProcess{
		tasks:         make(map[sdktypes.TaskIndex]*core.BalanceTaskInfo),
		finishedTasks: make(map[string]*types.FinishedTaskStatus),
		service:       aggService,
		logger:        logger,
		avsWriter:     avsWriter,
	}, nil
}

func (d *DeviationProcess) GetTaskByIndex(taskIndex sdktypes.TaskIndex) *core.BalanceTaskInfo {

	d.tasksMu.RLock()
	defer d.tasksMu.Unlock()

	res := d.tasks[taskIndex]

	return res
}

func (d *DeviationProcess) NewIndex() sdktypes.TaskIndex {
	d.nextTaskIndexMu.Lock()
	defer d.nextTaskIndexMu.Unlock()

	res := d.nextTaskIndex
	d.nextTaskIndex += 1

	return res
}

func (d *DeviationProcess) GetTaskByRequestHash(reqHash string) *core.BalanceTaskInfo {
	d.tasksMu.RLock()
	defer d.tasksMu.Unlock()

	for _, task := range d.tasks {
		hash, err := core.GetRequestHash(task.ReqId, task.AgentAddress)

		if err != nil {
			return nil
		}

		if hash == reqHash {
			return task
		}
	}

	return nil
}

func (d *DeviationProcess) GetFinishedTaskByRequestHash(reqHash string) *types.FinishedTaskStatus {
	d.finishedTasksMu.RLock()
	defer d.finishedTasksMu.Unlock()

	return d.finishedTasks[reqHash]
}

func (d *DeviationProcess) SetFinishedTask(reqHash string, finished *types.FinishedTaskStatus) {
	d.finishedTasksMu.Lock()
	defer d.finishedTasksMu.Unlock()

	d.finishedTasks[reqHash] = finished
}

func (d *DeviationProcess) CreateTask(req *core.CreateDeviationTaskRequest) (*core.CreateDeviationTaskResponse, error) {
	d.logger.Infof("Received Create Task", "reqId", req.ReqId, "Agent", req.AgentAddress)

	reqHash, err := req.GetHash()

	if err != nil {
		d.logger.Errorf("Error while getting request Hash", "err", err)
		return nil, err
	}

	finishedTask := d.GetFinishedTaskByRequestHash(reqHash)

	if finishedTask != nil {
		d.logger.Error("Task already finished", "reqId", req.ReqId.String(), "agent", req.AgentAddress.Hex())

		return nil, err
	}

	task := d.GetTaskByRequestHash(reqHash)

	if task == nil {
		taskIndex := d.NewIndex()
		task, err = d.sendNewTask(reqHash, taskIndex, req.ReferenceBlockNumber, req.ReqId, req.AgentAddress)

		if err != nil {
			d.logger.Error("send new task failed", "err", err)
			return nil, err
		}
	}
	return &core.CreateDeviationTaskResponse{Info: *task}, nil
}

func (d *DeviationProcess) ProcessSignedTaskResponse(req *core.SignedDeviationTaskRespRequest) (*core.SignedTaskRespResponse, error) {
	d.logger.Info("New Signature Received from operator", "signature", req.BlsSignature.Serialize(), "operatorId", req.OperatorId, "balance", req.Balance)

	reqHash, err := core.GetRequestHash(req.Balance.ReqId, req.Balance.AgentAddress)

	if err != nil {
		d.logger.Errorf("Error while creating request hash", "err", fmt.Sprintf("%#f", err))

		return nil, err
	}

	taskIndex := req.Balance.TaskIndex

	d.tasksMu.Lock()
	d.tasks[taskIndex] = &req.Balance
	d.tasksMu.Unlock()

	taskResponseDigest, err := req.Balance.GetBalanceHeaderHash()

	if err != nil {
		d.logger.Error("Error getting taskResponseDigest", "err", err)
		return nil, err
	}

	task := d.GetTaskByRequestHash(reqHash)

	if task == nil {
		d.logger.Errorf("Task doesn't exist for given request hash", "err", fmt.Sprintf("%#f", err))
		return nil, err
	}

	err = d.service.BlsAggregationService.ProcessNewSignature(
		context.Background(), taskIndex, taskResponseDigest,
		&req.BlsSignature, req.OperatorId,
	)

	if err != nil {
		d.logger.Error("ProcessNewSignature error", "err", err)
	}

	return &core.SignedTaskRespResponse{}, nil
}

func (d *DeviationProcess) sendNewTask(reqHash string, taskIndex uint32, refBlock uint64, reqId *big.Int, agentAddr common.Address) (*core.BalanceTaskInfo, error) {

	quorumDetails, err := d.service.GetQuorumDetails(context.Background(), refBlock)

	if err != nil {
		d.logger.Error("error while getting quorum details")
		return nil, err
	}

	newTask := &core.BalanceTaskInfo{
		QuorumNumbers:              quorumDetails.QuorumNums,
		QuorumThresholdPercentages: quorumDetails.QuorumThresholdPercentages,
		TaskIndex:                  taskIndex,
		ReferenceBlockNumber:       refBlock,
		ReqId:                      reqId,
		ActualBalance:              big.NewInt(0),
		IdealBalance:               big.NewInt(0),
		AgentAddress:               agentAddr,
	}

	d.tasksMu.Lock()
	d.tasks[taskIndex] = newTask
	d.tasksMu.Unlock()

	taskTimeToExpiry := time.Duration(10) * time.Minute

	d.logger.Infof("InitializeNewTask %v %v", taskIndex, taskTimeToExpiry)

	err = d.service.BlsAggregationService.InitializeNewTask(
		taskIndex,
		uint32(newTask.ReferenceBlockNumber),
		newTask.QuorumNumbers,
		newTask.QuorumThresholdPercentages,
		taskTimeToExpiry,
	)
	if err != nil {
		d.logger.Error("InitializeNewTask failed", "err", err)
		return nil, err
	}

	return newTask, nil
}

func (d *DeviationProcess) SendAggregatedResponseToContract(blsResp blsaggregator.BLSAggregationServiceResponse) {

	nonSignerStakes, err := d.service.ReceiveAggregatedSignatures(blsResp)

	if err != nil {
		d.logger.Error("error while building non signer stake signatures")
		return
	}

	task := d.GetTaskByIndex(blsResp.TaskIndex)

	res, err := d.avsWriter.SendConfirmationBatch(context.Background(), task, nonSignerStakes)

	if err != nil {
		d.logger.Error("aggregator failed to respond to task", "err", err)
		return
	}

	if res != nil {
		reqHash, err := core.GetRequestHash(task.ReqId, task.AgentAddress)

		if err != nil {
			d.logger.Error("agg failed to find req hash")
		}

		d.SetFinishedTask(reqHash, &types.FinishedTaskStatus{
			Message:          task,
			TxHash:           res.TxHash,
			BlockHash:        res.BlockHash,
			BlockNumber:      res.BlockNumber,
			TransactionIndex: res.TransactionIndex,
		})
	} else {
		d.logger.Error("the send confirm aler threw error", "err", res)
	}

}

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

type SlashProcess struct {
	tasks   map[sdktypes.TaskIndex]*core.AgentSlashTaskInfo
	tasksMu sync.RWMutex

	finishedTasks   map[string]*types.FinishedSlashTaskStatus
	finishedTasksMu sync.RWMutex

	nextTaskIndex   sdktypes.TaskIndex
	nextTaskIndexMu sync.RWMutex

	service *AggregatorService

	logger logging.Logger

	avsWriter *avs.AvsWriter
}

func NewSlashProcessor(c *config.Config, logger logging.Logger, avsWriter *avs.AvsWriter) (*SlashProcess, error) {

	aggService, err := NewAggregatorService(c, logger)

	if err != nil {
		panic(err)
	}

	return &SlashProcess{
		tasks:         make(map[sdktypes.TaskIndex]*core.AgentSlashTaskInfo),
		finishedTasks: make(map[string]*types.FinishedSlashTaskStatus),
		service:       aggService,
		logger:        logger,
		avsWriter:     avsWriter,
		nextTaskIndex: 0,
	}, nil
}

func (d *SlashProcess) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case res := <-d.service.GetResponseChannel():
			d.SendAggregatedResponseToContract(res)
		}
	}
}


func (d *SlashProcess) GetTaskByIndex(taskIndex sdktypes.TaskIndex) *core.AgentSlashTaskInfo {

	d.tasksMu.RLock()
	defer d.tasksMu.RUnlock()

	res := d.tasks[taskIndex]

	return res
}

func (d *SlashProcess) NewIndex() sdktypes.TaskIndex {
	d.nextTaskIndexMu.Lock()
	defer d.nextTaskIndexMu.Unlock()

	res := d.nextTaskIndex
	d.nextTaskIndex += 1

	return res
}

func (d *SlashProcess) GetTaskByRequestHash(reqHash string) *core.AgentSlashTaskInfo {
	d.tasksMu.RLock()
	defer d.tasksMu.RUnlock()

	for _, task := range d.tasks {
		hash, err := core.GetRequestHash(task.TxnId, task.AgentAddress)

		if err != nil {
			return nil
		}

		if hash == reqHash {
			return task
		}
	}

	return nil
}

func (d *SlashProcess) GetFinishedTaskByRequestHash(hash string) *types.FinishedSlashTaskStatus {
	d.finishedTasksMu.RLock()
	defer d.finishedTasksMu.RUnlock()

	return d.finishedTasks[hash]
}

func (d *SlashProcess) SetFinishedTask(reqHash string, finished *types.FinishedSlashTaskStatus) {
	d.finishedTasksMu.Lock()
	defer d.finishedTasksMu.Unlock()

	d.finishedTasks[reqHash] = finished
}

func (d *SlashProcess) CreateTask(req *core.CreateAgentSlashTaskRequest) (*core.CreateAgentSlashTaskResponse, error) {
	// d.tasksMu.Lock()
	// defer d.tasksMu.Unlock()

	d.logger.Info("Received Create Task", "reqId", req.TxnId, "Agent", req.AgentAddress)

	reqHash, err := req.GetHash()

	if err != nil {
		d.logger.Errorf("Error while getting request Hash", "err", err)
		return nil, err
	}

	finishedTask := d.GetFinishedTaskByRequestHash(reqHash)

	if finishedTask != nil {
		d.logger.Error("Task already finished", "reqId", req.TxnId.String(), "agent", req.AgentAddress.Hex())

		return nil, err
	}

	task := d.GetTaskByRequestHash(reqHash)

	if task == nil {
		taskIndex := d.NewIndex()
		task, err = d.sendNewTask(reqHash, taskIndex, req.ReferenceBlockNumber, req.TxnId, req.AgentAddress)

		if err != nil {
			d.logger.Error("send new task failed", "err", err)
			return nil, err
		}
	}
	return &core.CreateAgentSlashTaskResponse{Info: *task}, nil
}

func (d *SlashProcess) ProcessSignedTaskResponse(req *core.SignedAgentSlashTaskRespRequest) (*core.SignedTaskRespResponse, error) {
	d.logger.Info("New Signature Received from operator", "signature", req.BlsSignature.Serialize(), "operatorId", req.OperatorId, "balance", req.Balance)

	reqHash, err := core.GetRequestHash(req.Balance.TxnId, req.Balance.AgentAddress)

	if err != nil {
		d.logger.Errorf("Error while creating request hash", "err", fmt.Sprintf("%#f", err))

		return nil, err
	}

	taskIndex := req.Balance.TaskIndex

	d.tasksMu.Lock()
	d.tasks[taskIndex] = &req.Balance
	d.tasksMu.Unlock()

	taskResponseDigest, err := req.Balance.GetHeaderHash()

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
		context.Background(), taskIndex, taskResponseDigest.UnderlyingType(),
		&req.BlsSignature, req.OperatorId,
	)

	if err != nil {
		d.logger.Error("ProcessNewSignature error", "err", err)
	}

	return &core.SignedTaskRespResponse{}, nil
}

func (d *SlashProcess) sendNewTask(reqHash string, taskIndex uint32, refBlock uint64, txnId *big.Int, agentAddr common.Address) (*core.AgentSlashTaskInfo, error) {

	quorumDetails, err := d.service.GetQuorumDetails(context.Background(), refBlock)

	if err != nil {
		d.logger.Error("error while getting quorum details")
		return nil, err
	}

	newTask := &core.AgentSlashTaskInfo{
		QuorumNumbers:              quorumDetails.QuorumNums,
		QuorumThresholdPercentages: quorumDetails.QuorumThresholdPercentages,
		TaskIndex:                  taskIndex,
		ReferenceBlockNumber:       refBlock,
		AgentAddress:               agentAddr,
		InsufficientBalance:        big.NewInt(0),
		TxnId:                      txnId,
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

func (d *SlashProcess) SendAggregatedResponseToContract(blsResp blsaggregator.BLSAggregationServiceResponse) {

	nonSignerStakes, err := d.service.ReceiveAggregatedSignatures(blsResp)

	if err != nil {
		d.logger.Error("error while building non signer stake signatures")
		return
	}

	task := d.GetTaskByIndex(blsResp.TaskIndex)

	res, err := d.avsWriter.SendConfirmationAgentSlash(context.Background(), task, nonSignerStakes)

	if err != nil {
		d.logger.Error("aggregator failed to respond to task", "err", err)
		return
	}

	if res != nil {
		reqHash, err := core.GetRequestHash(big.NewInt(0), task.AgentAddress)

		if err != nil {
			d.logger.Error("agg failed to find req hash")
		}

		d.SetFinishedTask(reqHash, &types.FinishedSlashTaskStatus{
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

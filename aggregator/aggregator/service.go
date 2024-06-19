package aggregator

import (
	// "context"

	"bytes"
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/Cedro-Finance/CedroAVS/aggregator/rpc"
	"github.com/Cedro-Finance/CedroAVS/core"
	avs "github.com/Cedro-Finance/CedroAVS/core/avsClient"
	"github.com/Cedro-Finance/CedroAVS/core/config"
	sdkclients "github.com/Layr-Labs/eigensdk-go/chainio/clients"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/eth"
	sdkEcdsa "github.com/Layr-Labs/eigensdk-go/crypto/ecdsa"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/Layr-Labs/eigensdk-go/types"
	"github.com/ethereum/go-ethereum/common"
	// "github.com/Layr-Labs/eigensdk-go/services/operatorsinfo"
)

type AggregatorService struct {
	logger                logging.Logger
	ethClient             eth.Client
	avsReader             *avs.AvsReader
	blsAggregationService *BLSAggregatorService

	cfg *config.Config

	tasks   map[types.TaskIndex]*core.BalanceTaskInfo
	tasksMu sync.RWMutex

	finishedTasks   map[[32]byte]*FinishedTaskStatus
	finishedTasksMu sync.RWMutex

	nextTaskIndex   types.TaskIndex
	nextTaskIndexMu sync.RWMutex
}

func NewAggregatorService(c *config.Config, logger logging.Logger) (*AggregatorService, error) {

	ethClient, err := eth.NewClient(c.EthConfig.RpcUrl)

	if err != nil {
		logger.Errorf("Error while creating Eth client", "err", err)
		return nil, err
	}

	avsReader, err := avs.NewAvsReader(c.AVS, ethClient, logger)

	if err != nil {
		logger.Errorf("Error while creating Avs reader", "err", err)
	}

	privateKey, err := sdkEcdsa.ReadKey(c.KeyStore.Path, c.KeyStore.Password)

	if err != nil {
		return nil, err
	}

	fmt.Println(c.EthConfig.RpcUrl, c.EthConfig.WsUrl)

	chainioConfig := sdkclients.BuildAllConfig{
		EthHttpUrl:                 c.EthConfig.RpcUrl,
		EthWsUrl:                   c.EthConfig.WsUrl,
		RegistryCoordinatorAddr:    c.AVS.RegistryCoordinatorAddr.String(),
		OperatorStateRetrieverAddr: c.AVS.OperatorStateRetrieverAddr.String(),
		AvsName:                    "CedroAVS",
		PromMetricsIpPortAddress:   "Cedro",
	}

	clients, err := sdkclients.BuildAll(chainioConfig, privateKey, logger)
	if err != nil {
		logger.Errorf("Cannot create sdk clients", "err", err)
		return nil, err
	}

	// operatorsinfoService := operatorsinfo.NewOperatorsInfoServiceInMemory(context.Background(), clients.AvsRegistryChainSubscriber, clients.AvsRegistryChainReader, logger)
	// avsRegistryService := avsregistry.NewAvsRegistryServiceChainCaller(avsReader, operatorsinfoService, logger)
	blsAggregationService := NewBlsAggregatorService(avsReader, logger)

	return &AggregatorService{
		logger:                logger,
		avsReader:             avsReader,
		ethClient:             clients.EthHttpClient,
		blsAggregationService: blsAggregationService,
		tasks:                 make(map[types.TaskIndex]*core.BalanceTaskInfo),
		finishedTasks:         make(map[[32]byte]*FinishedTaskStatus),
		cfg:                   c,
	}, nil

}

var _ rpc.AggregatorRpcHandler = (*AggregatorService)(nil)

func (ags *AggregatorService) GetTaskByRequestHash(reqHash []byte) *core.BalanceTaskInfo {
	ags.tasksMu.RLock()
	defer ags.tasksMu.RUnlock()

	for _, task := range ags.tasks {
		hash, err := core.GetRequestHash(task.ReqId, task.AgentAddress)
		if err != nil {
			return nil
		}
		if bytes.Equal(hash, reqHash) {
			return task
		}
	}

	return nil
}

func (agg *AggregatorService) GetTaskByIndex(taskIndex types.TaskIndex) *core.BalanceTaskInfo {
	agg.tasksMu.RLock()
	defer agg.tasksMu.RUnlock()

	res := agg.tasks[taskIndex]

	return res
}

func (agg *AggregatorService) newIndex() types.TaskIndex {
	agg.nextTaskIndexMu.Lock()
	defer agg.nextTaskIndexMu.Unlock()

	res := agg.nextTaskIndex
	agg.nextTaskIndex += 1

	return res
}

// Called by Operator
func (ags *AggregatorService) InitOperator(req *core.InitOperatorRequest) (*core.InitOperatorResponse, error) {
	reply := &core.InitOperatorResponse{
		Ok: false,
	}

	// if ags.cfg.AVS.OperatorStateRetrieverAddr.Hex() != req.GetOperatorStateRetrieverAddr() {
	// 	return reply, nil
	// }

	reply.Ok = true

	return reply, nil
}

func (ags *AggregatorService) CreateTask(req *core.CreateTaskRequest) (*core.CreateTaskResponse, error) {
	ags.logger.Infof("Received Create Task", "reqId", req.ReqId, "Agent", req.AgentAddress)

	reqHash, err := core.GetRequestHash(req.ReqId, req.AgentAddress)

	if err != nil {
		ags.logger.Errorf("Error while getting request Hash", "err", err)
		return nil, err
	}

	task := ags.GetTaskByRequestHash(reqHash)

	if task == nil {
		taskIndex := ags.newIndex()
		task, err = ags.sendNewTask(core.Bytes32(reqHash), taskIndex, req.ReferenceBlockNumber, req.ReqId, req.AgentAddress)

		if err != nil {
			ags.logger.Error("send new task failed", "err", err)
			return nil, err
		}
	}
	return &core.CreateTaskResponse{Info: *task}, nil
}

// func convertToFixedByteArray(b []byte) ([32]byte, error) {
// 	if len(b) != 32 {
// 		return [32]byte{}, fmt.Errorf("byte slice is not 32 bytes long")
// 	}
// 	var array [32]byte
// 	copy(array[:], b)
// 	return array, nil
// }

func (ags *AggregatorService) ProcessSignedTaskResponse(req *core.SignedTaskRespRequest) (*core.SignedTaskRespResponse, error) {
	ags.logger.Info("New Signature Received from operator", "signature", req.BlsSignature.Serialize(), "operatorId", req.OperatorId, "balance", req.Balance)

	reqHash, err := core.GetRequestHash(req.Balance.ReqId, req.Balance.AgentAddress)

	if err != nil {
		ags.logger.Errorf("Error while creating request hash", "err", fmt.Sprintf("%#f", err))

		return nil, err
	}

	taskIndex := req.Balance.TaskIndex

	ags.tasksMu.Lock()
	ags.tasks[taskIndex] = &req.Balance
	ags.tasksMu.Unlock()

	taskResponseDigest, err := req.Balance.GetBalanceHeaderHash()

	if err != nil {
		ags.logger.Error("Error getting taskResponseDigest", "err", err)
		return nil, err
	}

	task := ags.GetTaskByRequestHash(reqHash)

	if task == nil {
		ags.logger.Errorf("Task doesn't exist for given request hash", "err", fmt.Sprintf("%#f", err))
		return nil, err
	}

	if err != nil {
		ags.logger.Error("Error while typecasting", "err", err)
		return nil, err
	}

	err = ags.blsAggregationService.ProcessNewSignature(
		context.Background(), taskIndex, taskResponseDigest,
		&req.BlsSignature, req.OperatorId,
	)

	if err != nil {
		ags.logger.Error("ProcessNewSignature error", "err", err)
	}

	fmt.Println("Task", task)

	return &core.SignedTaskRespResponse{}, nil
}

func (agg *AggregatorService) GetResponseChannel() <-chan BLSAggregationServiceResponse {
	return agg.blsAggregationService.GetResponseChannel()
}

func (agg *AggregatorService) sendNewTask(reqHash core.Bytes32, taskIndex uint32, refBlock uint64, reqId *big.Int, agentAddr common.Address) (*core.BalanceTaskInfo, error) {
	agg.logger.Info("Aggregator sending new task", "alert", reqHash, "task", taskIndex)

	var err error

	agg.logger.Info("get from Contract", "referenceBlockNumber", refBlock)

	quorumNumbers, err := agg.avsReader.GetQuorumsByBlockNumber(context.Background(), uint32(refBlock))
	if err != nil {
		agg.logger.Error("GetQuorumCountByBlockNumber failed", "err", err)
		return nil, err
	}
	agg.logger.Info("get quorumNumbers ", "quorumNumbers", fmt.Sprintf("%v", quorumNumbers))

	// if len(quorumNumbers) < len(agg.cfg.QuorumNums) {
	// 	agg.logger.Error("the cfg quorum numbers is larger to the layer1, it will commit failed")
	// 	return nil, fmt.Errorf("the quorum numbers is larger to the layer1 %v, expected %v", agg.cfg.QuorumNums, quorumNumbers)
	// }

	// // just use config value
	// quorumNumbers = agg.cfg.QuorumNums

	quorumThresholdPercentages, err := agg.avsReader.GetQuorumThresholdPercentages(context.Background(), uint32(refBlock), quorumNumbers)
	if err != nil {
		agg.logger.Error("GetQuorumThresholdPercentages failed", "err", err)
		return nil, err
	}

	agg.logger.Info(
		"quorum datas",
		"numbers", fmt.Sprintf("%v", quorumNumbers.UnderlyingType()),
		"thresholdPercentages", fmt.Sprintf("%v", quorumThresholdPercentages.UnderlyingType()),
	)

	newAlertTask := &core.BalanceTaskInfo{
		QuorumNumbers:              quorumNumbers,
		QuorumThresholdPercentages: quorumThresholdPercentages,
		TaskIndex:                  taskIndex,
		ReferenceBlockNumber:       refBlock,
		ReqId:                      reqId,
		ActualBalance:              big.NewInt(0),
		IdealBalance:               big.NewInt(0),
		AgentAddress:               agentAddr,
	}

	agg.tasksMu.Lock()
	agg.tasks[taskIndex] = newAlertTask
	agg.tasksMu.Unlock()

	// TODO(samlaf): we use seconds for now, but we should ideally pass a blocknumber to the blsAggregationService
	// and it should monitor the chain and only expire the task aggregation once the chain has reached that block number.
	taskTimeToExpiry := time.Duration(10) * time.Minute

	agg.logger.Infof("InitializeNewTask %v %v", taskIndex, taskTimeToExpiry)

	err = agg.blsAggregationService.InitializeNewTask(
		taskIndex,
		uint32(newAlertTask.ReferenceBlockNumber),
		newAlertTask.QuorumNumbers,
		newAlertTask.QuorumThresholdPercentages,
		taskTimeToExpiry,
	)
	if err != nil {
		agg.logger.Error("InitializeNewTask failed", "err", err)
		return nil, err
	}
	return newAlertTask, nil
}

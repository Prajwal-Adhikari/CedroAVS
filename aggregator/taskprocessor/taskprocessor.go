package taskprocessor

import (
	"context"
	"fmt"
	"sync"
	"time"

	blsaggregator "github.com/Cedro-Finance/CedroAVS/aggregator/bls"
	"github.com/Cedro-Finance/CedroAVS/aggregator/types"
	"github.com/Cedro-Finance/CedroAVS/api/graphql"
	"github.com/Cedro-Finance/CedroAVS/core"
	avs "github.com/Cedro-Finance/CedroAVS/core/avsClient"
	cedroservicemanager "github.com/Cedro-Finance/CedroAVS/core/bindings/CedroServiceManager"
	"github.com/Cedro-Finance/CedroAVS/core/config"
	"github.com/Cedro-Finance/CedroAVS/core/tasks"
	sdkclients "github.com/Layr-Labs/eigensdk-go/chainio/clients"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/eth"
	"github.com/Layr-Labs/eigensdk-go/crypto/bls"
	sdkEcdsa "github.com/Layr-Labs/eigensdk-go/crypto/ecdsa"
	"github.com/Layr-Labs/eigensdk-go/logging"
	sdktypes "github.com/Layr-Labs/eigensdk-go/types"
)

type ITaskProcessor interface {
	Start(ctx context.Context)

	CreateTask(req tasks.ICreateTaskRequest) (tasks.ICreateTaskResponse, error)

	ProcessSignedTaskResponse(req *tasks.CreateSignedTaskResRequest) (*tasks.SignedTaskRespResponse, error)
}

type TaskProcessor struct {
	taskType string

	tasks   map[sdktypes.TaskIndex]tasks.ITaskInfo
	tasksMu sync.RWMutex

	finishedTasks   map[string]*types.FinishedTaskStatus
	finishedTasksMu sync.RWMutex

	nextTaskIndex   sdktypes.TaskIndex
	nextTaskIndexMu sync.RWMutex

	//Aggregator service

	logger    logging.Logger
	ethClient eth.Client
	avsReader avs.IAVSReader

	blsAggrService *blsaggregator.BLSAggregatorService

	cfg         *config.Config
	graphClient *graphql.GraphQLClient

	avsWriter avs.IAvsWriter
}

func NewTaskProcessorService(taskType string, c *config.Config, logger logging.Logger, avsWriter avs.IAvsWriter) (*TaskProcessor, error) {

	ethClient, err := eth.NewClient(c.EthConfig.RpcUrl)

	if err != nil {
		logger.Error("Error while creating eth client", "err", err)
		return nil, err
	}

	graphClinet := graphql.NewGraphQLClient(c.GraphQLUrl, logger)

	avsReader, err := avs.NewAvsReader(c.AVS, ethClient, logger, graphClinet)

	if err != nil {
		logger.Error("error while creating avs reader", "err", err)
		return nil, err
	}

	privateKey, err := sdkEcdsa.ReadKey(c.KeyStore.Path, c.KeyStore.Password)

	if err != nil {
		return nil, err
	}

	chainioConfig := sdkclients.BuildAllConfig{
		EthHttpUrl:                 c.EthConfig.RpcUrl,
		EthWsUrl:                   c.EthConfig.RpcUrl,
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

	blsAggregationService := blsaggregator.NewBlsAggregatorService(avsReader, logger)

	return &TaskProcessor{
		tasks:         make(map[uint32]tasks.ITaskInfo),
		finishedTasks: make(map[string]*types.FinishedTaskStatus),

		logger:         logger,
		avsReader:      avsReader,
		ethClient:      clients.EthHttpClient,
		blsAggrService: blsAggregationService,
		cfg:            c,
		graphClient:    graphClinet,
		avsWriter:      avsWriter,
		taskType:       taskType,
	}, nil

}

func (t *TaskProcessor) Start(ctx context.Context) {
	t.logger.Infof("%s taskprocessor started", t.taskType)
	for {
		select {
		case <-ctx.Done():
			return
		case res := <-t.blsAggrService.GetResponseChannel():
			fmt.Println("Response signatre", res)
			// t.SendAggregatedResponseToContract(res)
		}
	}
}

func (t *TaskProcessor) CreateTask(req tasks.ICreateTaskRequest) (tasks.ICreateTaskResponse, error) {
	t.logger.Info("Received Create Task", "reqId", req.GetIdentifier().String(), "Agent", req.GetAgent().Hex(), "task", t.taskType)

	reqHash, err := req.GetHash()

	if err != nil {
		t.logger.Error("error getting request hash", "err", err)
		return tasks.ICreateTaskResponse{}, err

	}

	finishedTask := t.GetFinishedTaskByRequestHash(reqHash)

	if finishedTask != nil {
		t.logger.Error("Task already finished", "reqId", req.GetIdentifier().String(), "agent", req.GetAgent().Hex())

		return tasks.ICreateTaskResponse{}, fmt.Errorf("task already done")
	}

	task := t.GetTaskByRequestHash(reqHash)

	if task == nil {
		taskIndex := t.NewIndex()

		quorumDetails, err := t.GetQuorumDetails(context.TODO(), req.GetRefBlock())

		if err != nil {
			t.logger.Error("error while getting quorum details")
			return tasks.ICreateTaskResponse{}, nil
		}

		task = req.GetTaskInfo(taskIndex, quorumDetails.QuorumNums, quorumDetails.QuorumThresholdPercentages)

		err = t.InitializeNewTask(context.Background(), taskIndex, task, quorumDetails, req.GetRefBlock())

		if err != nil {
			t.logger.Error("InitializeNewTask failed", "err", err)
			return tasks.ICreateTaskResponse{}, err
		}

	}

	return tasks.ICreateTaskResponse{Info: task}, nil

}

func (t *TaskProcessor) ProcessSignedTaskResponse(req *tasks.CreateSignedTaskResRequest) (*tasks.SignedTaskRespResponse, error) {
	t.logger.Info("New Signature Received from operator for task", "tasktype", t.taskType, "signature", req.BlsSignature.Serialize(), "task", req.Task)

	//TODO
	reqHash := req.Task.GetUniqueHash()

	// if err != nil {
	// 	t.logger.Info("error while getting header hash", "err", err)
	// 	return nil, err
	// }

	taskIndex := req.Task.GetTaskIndex()

	//TODO: Check how mutex is performing here?
	// t.tasksMu.Lock()
	// t.tasks[taskIndex] = req.Task
	// t.tasksMu.Unlock()

	taskResponseDigest, err := req.Task.GetHeaderHash()

	if err != nil {
		t.logger.Error("errro getting taskresponse digest", "err", err)
		return nil, err
	}

	task := t.GetTaskByRequestHash(reqHash)

	if task == nil {
		t.logger.Error("task doesn't exist for given request hash", "err", err)

		return nil, fmt.Errorf("task doesn't exist for given request hash")
	}

	err = t.blsAggrService.ProcessNewSignature(context.TODO(), taskIndex, taskResponseDigest.UnderlyingType(), &req.BlsSignature, req.OperatorId)
	if err != nil {
		t.logger.Error("ProcessNewSignature error", "err", err)
	}

	return &tasks.SignedTaskRespResponse{}, err

}

func (t *TaskProcessor) GetTaskByIndex(taskIndex sdktypes.TaskIndex) tasks.ITaskInfo {
	t.tasksMu.RLock()
	defer t.tasksMu.RUnlock()

	res := t.tasks[taskIndex]

	return res
}

func (t *TaskProcessor) NewIndex() sdktypes.TaskIndex {
	t.nextTaskIndexMu.Lock()
	defer t.nextTaskIndexMu.Unlock()

	res := t.nextTaskIndex
	t.nextTaskIndex += 1

	return res
}

func (d *TaskProcessor) GetTaskByRequestHash(reqHash string) tasks.ITaskInfo {

	d.tasksMu.RLock()
	defer d.tasksMu.RUnlock()

	for _, task := range d.tasks {
		hash := task.GetUniqueHash()

		// if err != nil {
		// 	return nil
		// }

		if hash == reqHash {
			return task
		}
	}

	return nil
}

func (t *TaskProcessor) GetFinishedTaskByRequestHash(reqHash string) *types.FinishedTaskStatus {
	t.finishedTasksMu.RLock()
	defer t.finishedTasksMu.RUnlock()

	return t.finishedTasks[reqHash]
}

func (t *TaskProcessor) SetFinishedTask(reqHash string, finished *types.FinishedTaskStatus) {
	t.finishedTasksMu.Lock()
	defer t.finishedTasksMu.Unlock()

	t.finishedTasks[reqHash] = finished
}

type QuorumDetails struct {
	QuorumNums                 sdktypes.QuorumNums
	QuorumThresholdPercentages sdktypes.QuorumThresholdPercentages
}

func (agg *TaskProcessor) GetQuorumDetails(ctx context.Context, refBlock uint64) (*QuorumDetails, error) {

	var err error

	agg.logger.Info("get from Contract", "referenceBlockNumber", refBlock)

	quorumNumbers, err := agg.avsReader.GetQuorumsByBlockNumber(context.Background(), uint32(refBlock))
	if err != nil {
		agg.logger.Error("GetQuorumCountByBlockNumber failed", "err", err)
		return nil, err
	}
	agg.logger.Info("get quorumNumbers ", "quorumNumbers", fmt.Sprintf("%v", quorumNumbers))

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

	return &QuorumDetails{
		QuorumNums:                 quorumNumbers,
		QuorumThresholdPercentages: quorumThresholdPercentages,
	}, nil
}

func (agg *TaskProcessor) InitializeNewTask(ctx context.Context, taskIndex sdktypes.TaskIndex, task tasks.ITaskInfo, quorumDetails *QuorumDetails, refBlock uint64) error {

	agg.tasksMu.Lock()
	agg.tasks[taskIndex] = task
	agg.tasksMu.Unlock()

	//TODO: Define properly
	taskTimeToExpiry := time.Duration(10) * time.Minute

	err := agg.blsAggrService.InitializeNewTask(
		taskIndex,
		uint32(refBlock),
		quorumDetails.QuorumNums,
		quorumDetails.QuorumThresholdPercentages,
		taskTimeToExpiry,
	)

	return err
}

func (agg *TaskProcessor) ReceiveAggregatedSignatures(blsAggServiceResp blsaggregator.BLSAggregationServiceResponse) (cedroservicemanager.IBLSSignatureCheckerNonSignerStakesAndSignature, error) {
	// TODO: check if blsAggServiceResp contains an err
	if blsAggServiceResp.Err != nil {
		agg.logger.Error("BlsAggregationServiceResponse contains an error", "err", blsAggServiceResp.Err)
		// panicing to help with debugging (fail fast), but we shouldn't panic if we run this in production
		return cedroservicemanager.IBLSSignatureCheckerNonSignerStakesAndSignature{}, blsAggServiceResp.Err
	}
	nonSignerPubkeys := []cedroservicemanager.BN254G1Point{}
	for _, nonSignerPubkey := range blsAggServiceResp.NonSignersPubkeysG1 {
		nonSignerPubkeys = append(nonSignerPubkeys, core.ConvertToBN254G1Point(nonSignerPubkey))
	}
	quorumApks := []cedroservicemanager.BN254G1Point{}
	for _, quorumApk := range blsAggServiceResp.QuorumApksG1 {
		quorumApks = append(quorumApks, core.ConvertToBN254G1Point(quorumApk))
	}

	var aggSig *bls.Signature
	for _, sign := range blsAggServiceResp.SignersAggSigG1 {
		if aggSig == nil {
			aggSig = bls.NewZeroSignature().Add(sign)
		} else {
			aggSig.Add(sign)
		}
	}

	// Aggregate the aggregated public keys. We reuse the first aggregated public key as the accumulator
	var aggPubKey *bls.G2Point
	for _, pKey := range blsAggServiceResp.SignersApkG2 {
		if aggPubKey == nil {
			aggPubKey = bls.NewZeroG2Point().Add(pKey)
		} else {
			aggPubKey.Add(pKey)
		}
	}

	nonSignerStakesAndSignature := cedroservicemanager.IBLSSignatureCheckerNonSignerStakesAndSignature{
		NonSignerPubkeys:             nonSignerPubkeys,
		QuorumApks:                   quorumApks,
		ApkG2:                        core.ConvertToBN254G2Point(aggPubKey),
		Sigma:                        core.ConvertToBN254G1Point(aggSig.G1Point),
		NonSignerQuorumBitmapIndices: blsAggServiceResp.NonSignerQuorumBitmapIndices,
		QuorumApkIndices:             blsAggServiceResp.QuorumApkIndices,
		TotalStakeIndices:            blsAggServiceResp.TotalStakeIndices,
		NonSignerStakeIndices:        blsAggServiceResp.NonSignerStakeIndices,
	}

	return nonSignerStakesAndSignature, nil
}

func (t *TaskProcessor) SendAggregatedResponseToContract(blsResp blsaggregator.BLSAggregationServiceResponse) {

	nonSignerStakes, err := t.ReceiveAggregatedSignatures(blsResp)

	if err != nil {
		t.logger.Error("error while building non signer stake signatures")
		return
	}

	task := t.GetTaskByIndex(blsResp.TaskIndex)

	res, err := t.avsWriter.SendConfirmationBatch(context.Background(), task, nonSignerStakes)

	if err != nil {
		t.logger.Error("aggregator failed to respond to task", "err", err)
		return
	}

	if res != nil {
		//delete from unfinisehed task

		reqHash := task.GetUniqueHash()

		// if err != nil {
		// 	t.logger.Error("agg failed to find req hash")
		// }

		t.SetFinishedTask(reqHash, &types.FinishedTaskStatus{
			Message:          task,
			TxHash:           res.TxHash,
			BlockHash:        res.BlockHash,
			BlockNumber:      res.BlockNumber,
			TransactionIndex: res.TransactionIndex,
		})
	} else {
		t.logger.Error("the send confirm aler threw error", "err", res)
	}

}

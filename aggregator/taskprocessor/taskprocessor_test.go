package taskprocessor

import (
	"math/big"
	"sync"
	"testing"

	blsaggregator "github.com/Cedro-Finance/CedroAVS/aggregator/bls"
	"github.com/Cedro-Finance/CedroAVS/aggregator/types"
	"github.com/Cedro-Finance/CedroAVS/api/graphql"
	avs "github.com/Cedro-Finance/CedroAVS/core/avsClient"
	"github.com/Cedro-Finance/CedroAVS/core/config"
	"github.com/Cedro-Finance/CedroAVS/core/tasks"
	"github.com/Layr-Labs/eigensdk-go/crypto/bls"
	"github.com/Layr-Labs/eigensdk-go/logging"
	sdktypes "github.com/Layr-Labs/eigensdk-go/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockCreateTaskReq struct {
	Req int
}

func (m *mockCreateTaskReq) GetHash() (string, error) {
	return "", nil
}

func (m *mockCreateTaskReq) GetAgent() common.Address {
	return common.Address{}
}

func (m *mockCreateTaskReq) GetIdentifier() *big.Int {
	return big.NewInt(0)
}

func (m *mockCreateTaskReq) GetRefBlock() uint64 {
	return uint64(1)
}

func (m *mockCreateTaskReq) GetTaskInfo(taskIndex sdktypes.TaskIndex, quorumNums sdktypes.QuorumNums, quorumThresholds sdktypes.QuorumThresholdPercentages) tasks.ITaskInfo {
	return &mockTaskInfo{
		Value: 100,
	}
}

type mockTaskInfo struct {
	Value int
}

func (m *mockTaskInfo) EncodeSignHash() ([]byte, error) {
	return nil, nil
}

func (m *mockTaskInfo) GetHeaderHash() (sdktypes.Bytes32, error) {
	return sdktypes.Bytes32{}, nil
}

func (m *mockTaskInfo) GetTaskIndex() sdktypes.TaskIndex {
	return uint32(1)
}

func (m *mockTaskInfo) GetUniqueHash() string {
	return ""
}

func TestTaskProcessor(t *testing.T) {
	// mockConfig := &config.Config{
	// 	BlsPrivateKey: &bls.KeyPair{},
	// 	Production:    false,
	// 	AVS: config.AVSConfig{
	// 		RegistryCoordinatorAddr:    common.HexToAddress("0xeAAb9DA7B87Ee17dA9BB40Fcfa4e15E67d58c02C"),
	// 		OperatorStateRetrieverAddr: common.HexToAddress("0xbc5DB5e02F8ff5C2dBB69784e318348d3d85E2e7"),
	// 	},
	// 	KeyStore: &config.KeyStore{
	// 		Path:        "",
	// 		Password:    "",
	// 		BlsPath:     "",
	// 		BlsPassword: "",
	// 	},
	// 	GRpcAggregatorIpAddr: "",
	// 	Agents:               []config.ChainConfig{},
	// 	AgentManagerAddr:     common.HexToAddress("0xD0621425d41cb00989d44D84FE10eF35FB0D879B"),
	// 	EthConfig: config.EthConfig{
	// 		RpcUrl:        "https://ethereum-holesky-rpc.publicnode.com",
	// 		WsUrl:         "https://ethereum-holesky-rpc.publicnode.com",
	// 		StartHeight:   0,
	// 		BlockFinality: 0,
	// 		BlockTime:     0,
	// 		EthClient:     nil,
	// 	},
	// 	AgentManagerStartHeight:  0,
	// 	AgentManagerPollDuration: 0,
	// 	GraphQLUrl:               "https://api.studio.thegraph.com/query/41099/cedro-avs/v0.0.4x",
	// }

	logger := logging.NewNoopLogger()

	testOp1 := sdktypes.TestOperator{
		OperatorId:     sdktypes.OperatorId{0},
		StakePerQuorum: map[sdktypes.QuorumNum]sdktypes.StakeAmount{0: big.NewInt(100), 1: big.NewInt(200)},
		BlsKeypair:     newBlsKeyPairPanics("0x1"),
	}

	blockNum := uint32(1)
	blsAggServ := blsaggregator.NewBlsAggregatorService(avs.NewFakeAvsRegistryService(blockNum, []sdktypes.TestOperator{testOp1}), logger)

	taskprocessor := &TaskProcessor{
		taskType:        "",
		tasks:           map[uint32]tasks.ITaskInfo{},
		tasksMu:         sync.RWMutex{},
		finishedTasks:   map[string]*types.FinishedTaskStatus{},
		finishedTasksMu: sync.RWMutex{},
		nextTaskIndex:   0,
		nextTaskIndexMu: sync.RWMutex{},
		logger:          logger,
		ethClient:       nil,
		avsReader:       avs.NewFakeAvsRegistryService(blockNum, []sdktypes.TestOperator{testOp1}),
		blsAggrService:  blsAggServ,
		cfg:             &config.Config{},
		graphClient:     &graphql.GraphQLClient{},
		avsWriter:       &avs.AvsWriter{},
	}

	quorumNumbers := sdktypes.QuorumNums{0}
	quorumThresholdPercentages := []sdktypes.QuorumThresholdPercentage{100}

	t.Run("create task test", func(t *testing.T) {

		taskIndex := uint32(1)
		reqTask := &mockCreateTaskReq{
			Req: 1,
		}

		res, err := taskprocessor.CreateTask(reqTask)

		assert.Nil(t, err)

		wantedRes := tasks.ICreateTaskResponse{
			Info: reqTask.GetTaskInfo(taskIndex, quorumNumbers, quorumThresholdPercentages),
		}

		assert.Equal(t, wantedRes, res)

		delete(taskprocessor.tasks, 0)
	})

	t.Run("process new signature- without create task", func(t *testing.T) {
		taskResponse := &mockTaskInfo{123}

		taskResponseDigest, err := taskResponse.GetHeaderHash()

		require.Nil(t, err)

		signedReq := &tasks.CreateSignedTaskResRequest{
			Task:         taskResponse,
			BlsSignature: *testOp1.BlsKeypair.SignMessage(taskResponseDigest),
			OperatorId:   testOp1.OperatorId,
		}

		_, err = taskprocessor.ProcessSignedTaskResponse(signedReq)
		require.EqualError(t, err, "task doesn't exist for given request hash")
	})

	t.Run("process new signature", func(t *testing.T) {
		reqTask := &mockCreateTaskReq{
			Req: 1,
		}

		res, err := taskprocessor.CreateTask(reqTask)

		assert.Nil(t, err)

		taskResponseDigest, err := res.Info.GetHeaderHash()

		require.Nil(t, err)

		signedReq := &tasks.CreateSignedTaskResRequest{
			Task:         res.Info,
			BlsSignature: *testOp1.BlsKeypair.SignMessage(taskResponseDigest),
			OperatorId:   testOp1.OperatorId,
		}

		_, err = taskprocessor.ProcessSignedTaskResponse(signedReq)
		// require.EqualError(t, err, "task doesn't exist for given request hash")

		require.Nil(t, err)
	})

}

func newBlsKeyPairPanics(hexKey string) *bls.KeyPair {
	keypair, err := bls.NewKeyPairFromString(hexKey)
	if err != nil {
		panic(err)
	}
	return keypair
}


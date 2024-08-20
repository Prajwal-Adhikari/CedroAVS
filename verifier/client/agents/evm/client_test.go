package evm

import (
	"context"

	"math/big"
	"testing"
	"time"

	avs "github.com/Cedro-Finance/CedroAVS/core/avsClient"
	agentmanager "github.com/Cedro-Finance/CedroAVS/core/bindings/AgentManager"
	"github.com/Cedro-Finance/CedroAVS/core/config"
	"github.com/Cedro-Finance/CedroAVS/verifier/client/agents"
	mockagentmanager "github.com/Cedro-Finance/CedroAVS/verifier/client/agents/evm/mocks"
	"github.com/Cedro-Finance/CedroAVS/verifier/common/mocks"
	"github.com/Layr-Labs/eigensdk-go/logging"

	// "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetFinalizedBlock(t *testing.T) {
	ethCl := new(mocks.EthClient)

	currentBlock := uint64(1000)
	blockFinality := 10

	ethCl.On("BlockNumber", mock.Anything).Return(currentBlock, nil)
	ethCl.On("BlockByNumber", mock.Anything, currentBlock).Return(types.Block{}, nil)

	cl := &Client{
		eth:       ethCl,
		blockTime: time.Duration(10) * time.Second,
		logger:    logging.NewNoopLogger(),
		config: config.ChainConfig{
			BlockFinality: uint64(blockFinality),
		},
	}

	curBlock, err := cl.GetCurrentBlock(context.Background())
	require.Nil(t, err)

	finalizedBlock, err := cl.GetFinalizedBlock(context.Background())

	require.Nil(t, err)
	require.Equal(t, currentBlock, curBlock)

	require.Equal(t, currentBlock-uint64(blockFinality), finalizedBlock)
}

func TestInitializeAgentPod(t *testing.T) {

	ethCl := new(mocks.EthClient)

	cl := &Client{
		eth:             ethCl,
		blockTime:       time.Duration(10) * time.Second,
		logger:          logging.NewNoopLogger(),
		config:          config.ChainConfig{},
		avsReader:       &avs.AvsReader{},
		nextBlockHeight: uint64(0),
	}

	require.Nil(t, cl.agent)

	agentAddr := ethCommon.HexToAddress("0x01234")

	cl.InitializeAgentPod(agentAddr)

	assert.NotNil(t, cl.agent)
}

func TestGetBlockFromTimestamp(t *testing.T) {
	ethCl := new(mocks.EthClient)

	currentBlock := uint64(1000)
	currTime := uint64(1527211625)
	blockFinality := 10

	ethCl.On("BlockNumber", mock.Anything).Return(currentBlock, nil)
	ethCl.On("BlockByNumber", mock.Anything, big.NewInt(int64(currentBlock))).Return(types.NewBlock(&types.Header{
		Time: currTime,
	}, nil, nil, nil), nil)

	// blockTime := 1

	startHeight := 1000

	cl := &Client{
		eth:       ethCl,
		blockTime: time.Duration(10) * time.Second,
		logger:    logging.NewNoopLogger(),
		config: config.ChainConfig{
			BlockFinality: uint64(blockFinality),
		},
		nextBlockHeight: uint64(startHeight),
	}

	blockNum, err := cl.GetBlockFromTimestamp(context.Background(), currTime-5*cl.config.BlockFinality)

	require.Nil(t, err)

	require.Equal(t, currentBlock-5, blockNum)

}

func TestPollEvents(t *testing.T) {
	currentBlock := uint64(1000)
	currTime := uint64(1527211625)
	blockFinality := 10
	startHeight := uint64(100)

	deviationEventHash := ethCommon.HexToHash("0xf1f6ae32e492ed94fefcd957830c206af14cae19c30c54be0814d27e8cdb5ac2")
	slashEventHash := ethCommon.HexToHash("0xc01e2e0defd55e48b29c6b1c0abaa6e53cf62dd50601aadb27e0b2706d3da2ea")

	agentManagerAdd := ethCommon.HexToAddress("0x0123")

	mockLog := types.Log{
		Address:     agentManagerAdd,
		Topics:      []ethCommon.Hash{deviationEventHash},
		BlockHash:   ethCommon.HexToHash("hash"),
		BlockNumber: 105,
	}

	mockLog2 := types.Log{
		Address:     agentManagerAdd,
		Topics:      []ethCommon.Hash{deviationEventHash},
		BlockHash:   ethCommon.HexToHash("hash2"),
		BlockNumber: 205,
	}

	mockSlashLog := types.Log{
		Address:   agentManagerAdd,
		Topics:    []ethCommon.Hash{slashEventHash},
		BlockHash: ethCommon.HexToHash("hash"),
	}

	manCl := new(mockagentmanager.IAgentManager)

	manCl.On("ParseDeviationCheck", mockLog).Return(&agentmanager.AgentManagerDeviationCheck{
		Agent:     agentManagerAdd,
		RequestID: big.NewInt(1),
	}, nil)

	manCl.On("ParseDeviationCheck", mockLog2).Return(&agentmanager.AgentManagerDeviationCheck{
		Agent:     agentManagerAdd,
		RequestID: big.NewInt(2),
	}, nil)

	manCl.On("ParseSlashAgent", mockSlashLog).Return(&agentmanager.AgentManagerSlashAgent{
		Agent:         agentManagerAdd,
		TxnId:         big.NewInt(1),
		AmountToSlash: big.NewInt(1000),
	}, nil)

	ethCl := new(mocks.EthClient)
	ethCl.On("BlockNumber", mock.Anything).Return(currentBlock, nil)
	ethCl.On("BlockByNumber", mock.Anything, big.NewInt(int64(currentBlock))).Return(types.NewBlock(&types.Header{
		Time: currTime,
	}, nil, nil, nil), nil)

	for index := startHeight; index <= currentBlock; index += defaultHeightDifferenceForFilterLogs + 1 {
		endHeight := index + defaultHeightDifferenceForFilterLogs

		if endHeight < 201 {
			ethCl.On("FilterLogs", mock.Anything, ethereum.FilterQuery{
				FromBlock: big.NewInt(int64(index)),
				ToBlock:   big.NewInt(int64(endHeight)),
				Addresses: []ethCommon.Address{agentManagerAdd},
				Topics:    [][]ethCommon.Hash{{deviationEventHash, slashEventHash}},
			}).Return([]types.Log{mockLog, mockSlashLog}, nil)
		} else if endHeight < 303 {
			ethCl.On("FilterLogs", mock.Anything, ethereum.FilterQuery{
				FromBlock: big.NewInt(int64(index)),
				ToBlock:   big.NewInt(int64(endHeight)),
				Addresses: []ethCommon.Address{agentManagerAdd},
				Topics:    [][]ethCommon.Hash{{deviationEventHash, slashEventHash}},
			}).Return([]types.Log{mockLog2}, nil)
		} else {
			ethCl.On("FilterLogs", mock.Anything, mock.Anything).Return([]types.Log{}, nil)
		}

	}

	ethCl.On("BlockByHash", mock.Anything, mock.Anything).Return(types.NewBlock(&types.Header{
		Time: currTime,
	}, nil, nil, nil), nil)

	// blockTime := 1

	avsReader := avs.NewFakeAvsRegistryService(uint32(startHeight), nil)

	cl := &Client{
		eth:       ethCl,
		blockTime: time.Duration(10) * time.Second,
		logger:    logging.NewNoopLogger(),
		config: config.ChainConfig{
			BlockFinality: uint64(blockFinality),
			AgentManager:  agentManagerAdd.Hex(),
		},
		startHeight:     uint64(startHeight),
		nextBlockHeight: uint64(startHeight),
		agentManager:    manCl,
		avsReader:       avsReader,
	}

	ch := make(chan agents.TaskRequest)

	slashCh := make(chan agents.TaskRequest)

	cl.startHeight = startHeight
	cl.nextBlockHeight = startHeight

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	require.Equal(t, startHeight, cl.nextBlockHeight)

	go cl.PollEvents(ctx, ch, slashCh)

	reqEvent := <-ch

	require.NotNil(t, reqEvent)

	require.Equal(t, big.NewInt(1), reqEvent.(*agents.DeviationRequest).ReqId)
	require.Equal(t, agentManagerAdd, reqEvent.(*agents.DeviationRequest).AgentAddress)

	slashEvent := <-slashCh

	require.NotNil(t, slashEvent)

	require.Equal(t, big.NewInt(1), slashEvent.(*agents.AgentSlashRequest).TxnId)
	require.Equal(t, agentManagerAdd, slashEvent.(*agents.AgentSlashRequest).AgentAddress)

	// require.Equal(t, uint64(100), cl.nextBlockHeight)

	reqEvent = <-ch
	require.Equal(t, big.NewInt(2), reqEvent.(*agents.DeviationRequest).ReqId)
	// require.Equal(t, uint64(302), cl.nextBlockHeight)

}

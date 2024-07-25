package blsaggregator

import (
	"context"
	"math/big"
	"testing"
	"time"

	avs "github.com/Cedro-Finance/CedroAVS/core/avsClient"
	"github.com/Layr-Labs/eigensdk-go/crypto/bls"
	"github.com/Layr-Labs/eigensdk-go/logging"

	// "github.com/Layr-Labs/eigensdk-go/services/avsregistry"
	"github.com/Layr-Labs/eigensdk-go/types"
	"github.com/stretchr/testify/require"
)

func newBlsKeyPairPanics(hexKey string) *bls.KeyPair {
	keypair, err := bls.NewKeyPairFromString(hexKey)
	if err != nil {
		panic(err)
	}
	return keypair
}

type mockTaskInfo struct {
	Value int
}

func (m *mockTaskInfo) EncodeSignHash() ([]byte, error) {
	return nil, nil
}

func (m *mockTaskInfo) GetHeaderHash() (types.Bytes32, error) {
	return types.Bytes32{}, nil
}

func TestBlsAgg(t *testing.T) {
	taskTimeToExpiry := 500 * time.Second

	t.Run("1 quorum 1 operator 1 correct signatuere", func(t *testing.T) {
		testOp1 := types.TestOperator{
			OperatorId:     types.OperatorId{0},
			StakePerQuorum: map[types.QuorumNum]types.StakeAmount{0: big.NewInt(100), 1: big.NewInt(200)},
			BlsKeypair:     newBlsKeyPairPanics("0x1"),
		}

		blockNum := uint32(1)
		taskIndex := types.TaskIndex(0)
		quorumNumbers := types.QuorumNums{0}
		quorumThresholdPercentages := []types.QuorumThresholdPercentage{100}
		taskResponse := &mockTaskInfo{123}

		taskResponseDigest, err := taskResponse.GetHeaderHash()

		require.Nil(t, err)

		blsSig := testOp1.BlsKeypair.SignMessage(taskResponseDigest)

		fakeAvsreader := avs.NewFakeAvsRegistryService(blockNum, []types.TestOperator{testOp1})

		noopLogger := logging.NewNoopLogger()
		blsAggServ := NewBlsAggregatorService(fakeAvsreader, noopLogger)

		err = blsAggServ.InitializeNewTask(taskIndex, blockNum, quorumNumbers, quorumThresholdPercentages, taskTimeToExpiry)
		require.Nil(t, err)
		err = blsAggServ.ProcessNewSignature(context.Background(), taskIndex, taskResponseDigest.UnderlyingType(), blsSig, testOp1.OperatorId)
		require.Nil(t, err)

		aggSign := make(map[types.QuorumNum]*bls.Signature)
		aggPubKeys := make(map[types.QuorumNum]*bls.G2Point, len(quorumNumbers))

		aggSign[0] = testOp1.BlsKeypair.SignMessage(taskResponseDigest)
		aggPubKeys[0] = testOp1.BlsKeypair.GetPubKeyG2()

		wantAggregationServiceResponse := BLSAggregationServiceResponse{
			Err:                          nil,
			TaskIndex:                    taskIndex,
			TaskResponse:                 taskResponseDigest,
			TaskResponseDigest:           taskResponseDigest.UnderlyingType(),
			NonSignersPubkeysG1:          []*bls.G1Point{},
			QuorumApksG1:                 []*bls.G1Point{testOp1.BlsKeypair.GetPubKeyG1()},
			SignersApkG2:                 aggPubKeys,
			SignersAggSigG1:              aggSign,
			NonSignerQuorumBitmapIndices: nil,
			QuorumApkIndices:             nil,
			TotalStakeIndices:            nil,
			NonSignerStakeIndices:        nil,
		}
		gotAggregationServiceResponse := <-blsAggServ.aggregatedResponsesC
		require.Equal(t, wantAggregationServiceResponse, gotAggregationServiceResponse)
		require.EqualValues(t, taskIndex, gotAggregationServiceResponse.TaskIndex)

	})

	t.Run("1 quorum 3 operator 3 correct signatuere", func(t *testing.T) {
		testOp1 := types.TestOperator{
			OperatorId:     types.OperatorId{1},
			StakePerQuorum: map[types.QuorumNum]types.StakeAmount{0: big.NewInt(300), 1: big.NewInt(200)},
			BlsKeypair:     newBlsKeyPairPanics("0x1"),
		}
		testOp2 := types.TestOperator{
			OperatorId:     types.OperatorId{2},
			StakePerQuorum: map[types.QuorumNum]types.StakeAmount{0: big.NewInt(100), 1: big.NewInt(200)},
			BlsKeypair:     newBlsKeyPairPanics("0x2"),
		}
		testOp3 := types.TestOperator{
			OperatorId:     types.OperatorId{3},
			StakePerQuorum: map[types.QuorumNum]types.StakeAmount{0: big.NewInt(100), 1: big.NewInt(200)},
			BlsKeypair:     newBlsKeyPairPanics("0x3"),
		}

		blockNum := uint32(1)
		taskIndex := types.TaskIndex(0)
		quorumNumbers := types.QuorumNums{0}
		quorumThresholdPercentages := []types.QuorumThresholdPercentage{100}
		taskResponse := &mockTaskInfo{123}

		taskResponseDigest, err := taskResponse.GetHeaderHash()

		require.Nil(t, err)

		fakeAvsreader := avs.NewFakeAvsRegistryService(blockNum, []types.TestOperator{testOp1, testOp2, testOp3})

		noopLogger := logging.NewNoopLogger()
		blsAggServ := NewBlsAggregatorService(fakeAvsreader, noopLogger)

		err = blsAggServ.InitializeNewTask(taskIndex, blockNum, quorumNumbers, quorumThresholdPercentages, taskTimeToExpiry)
		require.Nil(t, err)

		blsSig1 := testOp1.BlsKeypair.SignMessage(taskResponseDigest)
		err = blsAggServ.ProcessNewSignature(context.Background(), taskIndex, taskResponseDigest.UnderlyingType(), blsSig1, testOp1.OperatorId)
		require.Nil(t, err)
		

		blsSig2 := testOp2.BlsKeypair.SignMessage(taskResponseDigest)
		err = blsAggServ.ProcessNewSignature(context.Background(), taskIndex, taskResponseDigest.UnderlyingType(), blsSig2, testOp2.OperatorId)
		require.Nil(t, err)

		blsSig3 := testOp3.BlsKeypair.SignMessage(taskResponseDigest)
		err = blsAggServ.ProcessNewSignature(context.Background(), taskIndex, taskResponseDigest.UnderlyingType(), blsSig3, testOp3.OperatorId)
		require.Nil(t, err)

		aggSign := make(map[types.QuorumNum]*bls.Signature)
		aggPubKeys := make(map[types.QuorumNum]*bls.G2Point, len(quorumNumbers))

		aggSign[0] = blsSig1.Add(blsSig2).Add(blsSig3)
		aggPubKeys[0] = testOp1.BlsKeypair.GetPubKeyG2().Add(testOp2.BlsKeypair.GetPubKeyG2()).Add(testOp3.BlsKeypair.GetPubKeyG2())

		//Total quorumApksG1
		qApk := testOp1.BlsKeypair.GetPubKeyG1().Add(testOp2.BlsKeypair.GetPubKeyG1()).Add(testOp3.BlsKeypair.GetPubKeyG1())

		wantAggregationServiceResponse := BLSAggregationServiceResponse{
			Err:                          nil,
			TaskIndex:                    taskIndex,
			TaskResponse:                 taskResponseDigest,
			TaskResponseDigest:           taskResponseDigest.UnderlyingType(),
			NonSignersPubkeysG1:          []*bls.G1Point{},
			QuorumApksG1:                 []*bls.G1Point{qApk},
			SignersApkG2:                 aggPubKeys,
			SignersAggSigG1:              aggSign,
			NonSignerQuorumBitmapIndices: nil,
			QuorumApkIndices:             nil,
			TotalStakeIndices:            nil,
			NonSignerStakeIndices:        nil,
		}

		gotAggregationServiceResponse := <-blsAggServ.aggregatedResponsesC
		require.Equal(t, wantAggregationServiceResponse, gotAggregationServiceResponse)
		require.EqualValues(t, taskIndex, gotAggregationServiceResponse.TaskIndex)
	})

	t.Run("1 quorum 3 operator 2 correct signatuere", func(t *testing.T) {
		testOp1 := types.TestOperator{
			OperatorId:     types.OperatorId{1},
			StakePerQuorum: map[types.QuorumNum]types.StakeAmount{0: big.NewInt(250), 1: big.NewInt(200)},
			BlsKeypair:     newBlsKeyPairPanics("0x1"),
		}
		testOp2 := types.TestOperator{
			OperatorId:     types.OperatorId{2},
			StakePerQuorum: map[types.QuorumNum]types.StakeAmount{0: big.NewInt(150), 1: big.NewInt(200)},
			BlsKeypair:     newBlsKeyPairPanics("0x2"),
		}
		testOp3 := types.TestOperator{
			OperatorId:     types.OperatorId{3},
			StakePerQuorum: map[types.QuorumNum]types.StakeAmount{0: big.NewInt(100), 1: big.NewInt(200)},
			BlsKeypair:     newBlsKeyPairPanics("0x3"),
		}

		blockNum := uint32(1)
		taskIndex := types.TaskIndex(0)
		quorumNumbers := types.QuorumNums{0}
		quorumThresholdPercentages := []types.QuorumThresholdPercentage{60}
		taskResponse := &mockTaskInfo{123}

		taskResponseDigest, err := taskResponse.GetHeaderHash()

		require.Nil(t, err)

		fakeAvsreader := avs.NewFakeAvsRegistryService(blockNum, []types.TestOperator{testOp1, testOp2, testOp3})

		noopLogger := logging.NewNoopLogger()
		blsAggServ := NewBlsAggregatorService(fakeAvsreader, noopLogger)

		err = blsAggServ.InitializeNewTask(taskIndex, blockNum, quorumNumbers, quorumThresholdPercentages, taskTimeToExpiry)
		require.Nil(t, err)

		blsSig1 := testOp1.BlsKeypair.SignMessage(taskResponseDigest)
		err = blsAggServ.ProcessNewSignature(context.Background(), taskIndex, taskResponseDigest.UnderlyingType(), blsSig1, testOp1.OperatorId)
		require.Nil(t, err)

		blsSig2 := testOp2.BlsKeypair.SignMessage(taskResponseDigest)
		err = blsAggServ.ProcessNewSignature(context.Background(), taskIndex, taskResponseDigest.UnderlyingType(), blsSig2, testOp2.OperatorId)
		require.Nil(t, err)

		aggSign := make(map[types.QuorumNum]*bls.Signature)
		aggPubKeys := make(map[types.QuorumNum]*bls.G2Point, len(quorumNumbers))

		aggSign[0] = blsSig1.Add(blsSig2)
		aggPubKeys[0] = testOp1.BlsKeypair.GetPubKeyG2().Add(testOp2.BlsKeypair.GetPubKeyG2())

		//Total quorumApksG1
		qApk := testOp1.BlsKeypair.GetPubKeyG1().Add(testOp2.BlsKeypair.GetPubKeyG1()).Add(testOp3.BlsKeypair.GetPubKeyG1())

		wantAggregationServiceResponse := BLSAggregationServiceResponse{
			Err:                          nil,
			TaskIndex:                    taskIndex,
			TaskResponse:                 taskResponseDigest,
			TaskResponseDigest:           taskResponseDigest.UnderlyingType(),
			NonSignersPubkeysG1:          []*bls.G1Point{testOp3.BlsKeypair.GetPubKeyG1()},
			QuorumApksG1:                 []*bls.G1Point{qApk},
			SignersApkG2:                 aggPubKeys,
			SignersAggSigG1:              aggSign,
			NonSignerQuorumBitmapIndices: nil,
			QuorumApkIndices:             nil,
			TotalStakeIndices:            nil,
			NonSignerStakeIndices:        nil,
		}

		gotAggregationServiceResponse := <-blsAggServ.aggregatedResponsesC
		require.Equal(t, wantAggregationServiceResponse, gotAggregationServiceResponse)
		require.EqualValues(t, taskIndex, gotAggregationServiceResponse.TaskIndex)
	})

	t.Run("2 quorum 2 operator 2 correct signatuere", func(t *testing.T) {
		testOp1 := types.TestOperator{
			OperatorId:     types.OperatorId{1},
			StakePerQuorum: map[types.QuorumNum]types.StakeAmount{0: big.NewInt(250), 1: big.NewInt(200)},
			BlsKeypair:     newBlsKeyPairPanics("0x1"),
		}
		testOp2 := types.TestOperator{
			OperatorId:     types.OperatorId{2},
			StakePerQuorum: map[types.QuorumNum]types.StakeAmount{0: big.NewInt(150), 1: big.NewInt(200)},
			BlsKeypair:     newBlsKeyPairPanics("0x2"),
		}

		blockNum := uint32(1)
		taskIndex := types.TaskIndex(0)
		quorumNumbers := types.QuorumNums{0, 1}
		quorumThresholdPercentages := []types.QuorumThresholdPercentage{100, 100}
		taskResponse := &mockTaskInfo{123}

		taskResponseDigest, err := taskResponse.GetHeaderHash()

		require.Nil(t, err)

		fakeAvsreader := avs.NewFakeAvsRegistryService(blockNum, []types.TestOperator{testOp1, testOp2})

		noopLogger := logging.NewNoopLogger()
		blsAggServ := NewBlsAggregatorService(fakeAvsreader, noopLogger)

		err = blsAggServ.InitializeNewTask(taskIndex, blockNum, quorumNumbers, quorumThresholdPercentages, taskTimeToExpiry)
		require.Nil(t, err)

		blsSig1 := testOp1.BlsKeypair.SignMessage(taskResponseDigest)
		err = blsAggServ.ProcessNewSignature(context.Background(), taskIndex, taskResponseDigest.UnderlyingType(), blsSig1, testOp1.OperatorId)
		require.Nil(t, err)

		blsSig2 := testOp2.BlsKeypair.SignMessage(taskResponseDigest)
		err = blsAggServ.ProcessNewSignature(context.Background(), taskIndex, taskResponseDigest.UnderlyingType(), blsSig2, testOp2.OperatorId)
		require.Nil(t, err)

		aggSign := make(map[types.QuorumNum]*bls.Signature)
		aggPubKeys := make(map[types.QuorumNum]*bls.G2Point, len(quorumNumbers))

		aggSign[0] = bls.NewZeroSignature().Add(blsSig1).Add(blsSig2)
		aggPubKeys[0] = bls.NewZeroG2Point().Add(testOp1.BlsKeypair.GetPubKeyG2()).Add(testOp2.BlsKeypair.GetPubKeyG2())

		aggSign[1] = bls.NewZeroSignature().Add(blsSig1).Add(blsSig2)
		aggPubKeys[1] = bls.NewZeroG2Point().Add(testOp1.BlsKeypair.GetPubKeyG2()).Add(testOp2.BlsKeypair.GetPubKeyG2())

		//Total quorumApksG1

		wantAggregationServiceResponse := BLSAggregationServiceResponse{
			Err:                 nil,
			TaskIndex:           taskIndex,
			TaskResponse:        taskResponseDigest,
			TaskResponseDigest:  taskResponseDigest.UnderlyingType(),
			NonSignersPubkeysG1: []*bls.G1Point{},
			QuorumApksG1: []*bls.G1Point{
				bls.NewZeroG1Point().Add(testOp1.BlsKeypair.GetPubKeyG1()).Add(testOp2.BlsKeypair.GetPubKeyG1()),
				bls.NewZeroG1Point().Add(testOp1.BlsKeypair.GetPubKeyG1()).Add(testOp2.BlsKeypair.GetPubKeyG1()),
			},
			SignersApkG2:                 aggPubKeys,
			SignersAggSigG1:              aggSign,
			NonSignerQuorumBitmapIndices: nil,
			QuorumApkIndices:             nil,
			TotalStakeIndices:            nil,
			NonSignerStakeIndices:        nil,
		}

		gotAggregationServiceResponse := <-blsAggServ.aggregatedResponsesC
		require.Equal(t, wantAggregationServiceResponse, gotAggregationServiceResponse)
		require.EqualValues(t, taskIndex, gotAggregationServiceResponse.TaskIndex)
	})

	t.Run("1 quorum 1 operator 0 correct signatuere - task expired", func(t *testing.T) {
		testOp1 := types.TestOperator{
			OperatorId:     types.OperatorId{1},
			StakePerQuorum: map[types.QuorumNum]types.StakeAmount{0: big.NewInt(250), 1: big.NewInt(200)},
			BlsKeypair:     newBlsKeyPairPanics("0x1"),
		}

		blockNum := uint32(1)
		taskIndex := types.TaskIndex(0)
		quorumNumbers := types.QuorumNums{0}
		quorumThresholdPercentages := []types.QuorumThresholdPercentage{100}

		fakeAvsreader := avs.NewFakeAvsRegistryService(blockNum, []types.TestOperator{testOp1})

		noopLogger := logging.NewNoopLogger()
		blsAggServ := NewBlsAggregatorService(fakeAvsreader, noopLogger)

		err := blsAggServ.InitializeNewTask(taskIndex, blockNum, quorumNumbers, quorumThresholdPercentages, taskTimeToExpiry)
		require.Nil(t, err)

		aggSign := make(map[types.QuorumNum]*bls.Signature)
		aggPubKeys := make(map[types.QuorumNum]*bls.G2Point, len(quorumNumbers))

		aggSign[0] = bls.NewZeroSignature()
		aggPubKeys[0] = bls.NewZeroG2Point()

		//Total quorumApksG1
		wantAggregationServiceResponse := BLSAggregationServiceResponse{
			Err:       TaskExpiredErrorFn(taskIndex),
			TaskIndex: taskIndex,
		}

		gotAggregationServiceResponse := <-blsAggServ.aggregatedResponsesC
		require.Equal(t, wantAggregationServiceResponse, gotAggregationServiceResponse)
		require.EqualValues(t, taskIndex, gotAggregationServiceResponse.TaskIndex)
	})

	t.Run("1 quorum 1 operator 0 correct signatuere - task expired", func(t *testing.T) {
		testOp1 := types.TestOperator{
			OperatorId:     types.OperatorId{1},
			StakePerQuorum: map[types.QuorumNum]types.StakeAmount{0: big.NewInt(250), 1: big.NewInt(200)},
			BlsKeypair:     newBlsKeyPairPanics("0x1"),
		}

		blockNum := uint32(1)
		taskIndex := types.TaskIndex(0)
		quorumNumbers := types.QuorumNums{0}
		quorumThresholdPercentages := []types.QuorumThresholdPercentage{100}

		fakeAvsreader := avs.NewFakeAvsRegistryService(blockNum, []types.TestOperator{testOp1})

		noopLogger := logging.NewNoopLogger()
		blsAggServ := NewBlsAggregatorService(fakeAvsreader, noopLogger)

		err := blsAggServ.InitializeNewTask(taskIndex, blockNum, quorumNumbers, quorumThresholdPercentages, taskTimeToExpiry)
		require.Nil(t, err)

		aggSign := make(map[types.QuorumNum]*bls.Signature)
		aggPubKeys := make(map[types.QuorumNum]*bls.G2Point, len(quorumNumbers))

		aggSign[0] = bls.NewZeroSignature()
		aggPubKeys[0] = bls.NewZeroG2Point()

		//Total quorumApksG1
		wantAggregationServiceResponse := BLSAggregationServiceResponse{
			Err:       TaskExpiredErrorFn(taskIndex),
			TaskIndex: taskIndex,
		}

		gotAggregationServiceResponse := <-blsAggServ.aggregatedResponsesC
		require.Equal(t, wantAggregationServiceResponse, gotAggregationServiceResponse)
		require.EqualValues(t, taskIndex, gotAggregationServiceResponse.TaskIndex)
	})

	t.Run("1 quorum 2 operator 1 correct signatuere - task expired", func(t *testing.T) {
		testOp1 := types.TestOperator{
			OperatorId:     types.OperatorId{1},
			StakePerQuorum: map[types.QuorumNum]types.StakeAmount{0: big.NewInt(250), 1: big.NewInt(200)},
			BlsKeypair:     newBlsKeyPairPanics("0x1"),
		}

		testOp2 := types.TestOperator{
			OperatorId:     types.OperatorId{2},
			StakePerQuorum: map[types.QuorumNum]types.StakeAmount{0: big.NewInt(250), 1: big.NewInt(200)},
			BlsKeypair:     newBlsKeyPairPanics("0x2"),
		}

		blockNum := uint32(1)
		taskIndex := types.TaskIndex(0)
		quorumNumbers := types.QuorumNums{0}
		quorumThresholdPercentages := []types.QuorumThresholdPercentage{100}
		taskResponse := &mockTaskInfo{123}

		taskResponseDigest, err := taskResponse.GetHeaderHash()

		require.Nil(t, err)
		fakeAvsreader := avs.NewFakeAvsRegistryService(blockNum, []types.TestOperator{testOp1, testOp2})

		noopLogger := logging.NewNoopLogger()
		blsAggServ := NewBlsAggregatorService(fakeAvsreader, noopLogger)

		err = blsAggServ.InitializeNewTask(taskIndex, blockNum, quorumNumbers, quorumThresholdPercentages, taskTimeToExpiry)
		require.Nil(t, err)

		blsSig1 := testOp1.BlsKeypair.SignMessage(taskResponseDigest)
		err = blsAggServ.ProcessNewSignature(context.Background(), taskIndex, taskResponseDigest.UnderlyingType(), blsSig1, testOp1.OperatorId)
		require.Nil(t, err)

		aggSign := make(map[types.QuorumNum]*bls.Signature)
		aggPubKeys := make(map[types.QuorumNum]*bls.G2Point, len(quorumNumbers))

		aggSign[0] = bls.NewZeroSignature()
		aggPubKeys[0] = bls.NewZeroG2Point()

		//Total quorumApksG1
		wantAggregationServiceResponse := BLSAggregationServiceResponse{
			Err:       TaskExpiredErrorFn(taskIndex),
			TaskIndex: taskIndex,
		}

		gotAggregationServiceResponse := <-blsAggServ.aggregatedResponsesC
		require.Equal(t, wantAggregationServiceResponse, gotAggregationServiceResponse)
		require.EqualValues(t, taskIndex, gotAggregationServiceResponse.TaskIndex)
	})
}

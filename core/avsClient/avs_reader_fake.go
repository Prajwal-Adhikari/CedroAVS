package avs

import (
	"context"
	"errors"
	"math/big"

	opstateretriever "github.com/Layr-Labs/eigensdk-go/contracts/bindings/OperatorStateRetriever"
	"github.com/Layr-Labs/eigensdk-go/crypto/bls"
	"github.com/Layr-Labs/eigensdk-go/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type FakeAvsRegistryService struct {
	operators map[types.BlockNum]map[types.OperatorId]types.OperatorAvsState
}

// GetOperatorAddrsInQuorumsAtCurrentBlock implements IAVSReader.
func (f *FakeAvsRegistryService) GetOperatorAddrsInQuorumsAtCurrentBlock(opts *bind.CallOpts, quorumNumbers types.QuorumNums) ([][]common.Address, error) {
	panic("unimplemented")
}

// GetOperatorFromId implements IAVSReader.
func (f *FakeAvsRegistryService) GetOperatorFromId(opts *bind.CallOpts, operatorId types.Bytes32) (common.Address, error) {
	panic("unimplemented")
}

// GetOperatorId implements IAVSReader.
func (f *FakeAvsRegistryService) GetOperatorId(opts *bind.CallOpts, operatorAddress common.Address) ([32]byte, error) {
	panic("unimplemented")
}

// GetOperatorStakeInQuorumsOfOperatorAtCurrentBlock implements IAVSReader.
func (f *FakeAvsRegistryService) GetOperatorStakeInQuorumsOfOperatorAtCurrentBlock(opts *bind.CallOpts, operatorId types.Bytes32) (map[types.QuorumNum]*big.Int, error) {
	panic("unimplemented")
}

// GetOperatorsStakeInQuorumsAtBlock implements IAVSReader.
func (f *FakeAvsRegistryService) GetOperatorsStakeInQuorumsAtBlock(opts *bind.CallOpts, quorumNumbers types.QuorumNums, blockNumber uint32) ([][]opstateretriever.OperatorStateRetrieverOperator, error) {
	panic("unimplemented")
}

// GetOperatorsStakeInQuorumsAtCurrentBlock implements IAVSReader.
func (f *FakeAvsRegistryService) GetOperatorsStakeInQuorumsAtCurrentBlock(opts *bind.CallOpts, quorumNumbers types.QuorumNums) ([][]opstateretriever.OperatorStateRetrieverOperator, error) {
	panic("unimplemented")
}

// GetOperatorsStakeInQuorumsOfOperatorAtBlock implements IAVSReader.
func (f *FakeAvsRegistryService) GetOperatorsStakeInQuorumsOfOperatorAtBlock(opts *bind.CallOpts, operatorId types.Bytes32, blockNumber uint32) (types.QuorumNums, [][]opstateretriever.OperatorStateRetrieverOperator, error) {
	panic("unimplemented")
}

// GetOperatorsStakeInQuorumsOfOperatorAtCurrentBlock implements IAVSReader.
func (f *FakeAvsRegistryService) GetOperatorsStakeInQuorumsOfOperatorAtCurrentBlock(opts *bind.CallOpts, operatorId types.Bytes32) (types.QuorumNums, [][]opstateretriever.OperatorStateRetrieverOperator, error) {
	panic("unimplemented")
}

// GetQuorumCount implements IAVSReader.
func (f *FakeAvsRegistryService) GetQuorumCount(opts *bind.CallOpts) (uint8, error) {
	panic("unimplemented")
}

// GetQuorumThresholdPercentages implements IAVSReader.
func (f *FakeAvsRegistryService) GetQuorumThresholdPercentages(ctx context.Context, blockNumber uint32, quorums types.QuorumNums) (types.QuorumThresholdPercentages, error) {
	return []types.QuorumThresholdPercentage{100}, nil
}

// GetQuorumsByBlockNumber implements IAVSReader.
func (f *FakeAvsRegistryService) GetQuorumsByBlockNumber(ctx context.Context, blockNumber uint32) (types.QuorumNums, error) {
	return []types.QuorumNum{0}, nil
}

// IsOperatorRegistered implements IAVSReader.
func (f *FakeAvsRegistryService) IsOperatorRegistered(opts *bind.CallOpts, operatorAddress common.Address) (bool, error) {
	panic("unimplemented")
}

// QueryExistingRegisteredOperatorPubKeys implements IAVSReader.
func (f *FakeAvsRegistryService) QueryExistingRegisteredOperatorPubKeys(ctx context.Context, startBlock *big.Int, stopBlock *big.Int) ([]common.Address, []types.OperatorPubkeys, error) {
	panic("unimplemented")
}

// QueryExistingRegisteredOperatorSockets implements IAVSReader.
func (f *FakeAvsRegistryService) QueryExistingRegisteredOperatorSockets(ctx context.Context, startBlock *big.Int, stopBlock *big.Int) (map[types.Bytes32]types.Socket, error) {
	panic("unimplemented")
}

func NewFakeAvsRegistryService(blockNum types.BlockNum, operators []types.TestOperator) *FakeAvsRegistryService {
	fakeAvsRegistryService := &FakeAvsRegistryService{
		operators: map[types.BlockNum]map[types.OperatorId]types.OperatorAvsState{
			blockNum: {},
		},
	}
	for _, operator := range operators {
		fakeAvsRegistryService.operators[blockNum][operator.OperatorId] = types.OperatorAvsState{
			OperatorId: operator.OperatorId,
			OperatorInfo: types.OperatorInfo{
				Pubkeys: types.OperatorPubkeys{G1Pubkey: operator.BlsKeypair.GetPubKeyG1(), G2Pubkey: operator.BlsKeypair.GetPubKeyG2()},
				Socket:  "localhost:8080",
			},
			StakePerQuorum: operator.StakePerQuorum,
			BlockNumber:    blockNum,
		}
	}
	return fakeAvsRegistryService
}

var _ IAVSReader = (*FakeAvsRegistryService)(nil)

func (f *FakeAvsRegistryService) GetOperatorsAvsStateAtBlock(ctx context.Context, quorumNumbers types.QuorumNums, blockNumber types.BlockNum) (map[types.OperatorId]types.OperatorAvsState, error) {
	operatorsAvsState, ok := f.operators[blockNumber]
	if !ok {
		return nil, errors.New("block number not found")
	}
	return operatorsAvsState, nil
}

func (f *FakeAvsRegistryService) GetQuorumsAvsStateAtBlock(ctx context.Context, quorumNumbers types.QuorumNums, blockNumber types.BlockNum) (map[types.QuorumNum]types.QuorumAvsState, error) {
	operatorsAvsState, ok := f.operators[blockNumber]
	if !ok {
		return nil, errors.New("block number not found")
	}
	quorumsAvsState := make(map[types.QuorumNum]types.QuorumAvsState)
	for _, quorumNum := range quorumNumbers {
		aggPubkeyG1 := bls.NewG1Point(big.NewInt(0), big.NewInt(0))
		totalStake := big.NewInt(0)
		for _, operator := range operatorsAvsState {
			// only include operators that have a stake in this quorum
			if stake, ok := operator.StakePerQuorum[quorumNum]; ok {
				aggPubkeyG1.Add(operator.OperatorInfo.Pubkeys.G1Pubkey)
				totalStake.Add(totalStake, stake)
			}
		}
		quorumsAvsState[quorumNum] = types.QuorumAvsState{
			QuorumNumber: quorumNum,
			AggPubkeyG1:  aggPubkeyG1,
			TotalStake:   totalStake,
			BlockNumber:  blockNumber,
		}
	}
	return quorumsAvsState, nil
}

func (f *FakeAvsRegistryService) GetCheckSignaturesIndices(
	opts *bind.CallOpts, referenceBlockNumber types.BlockNum,
	quorumNumbers types.QuorumNums, nonSignerOperatorIds []types.OperatorId,
) (opstateretriever.OperatorStateRetrieverCheckSignaturesIndices, error) {
	return opstateretriever.OperatorStateRetrieverCheckSignaturesIndices{}, nil
}

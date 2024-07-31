package avs

import (
	"context"
	"math/big"

	"github.com/Cedro-Finance/CedroAVS/api/graphql"
	"github.com/Cedro-Finance/CedroAVS/core/config"
	sdkavsregistry "github.com/Layr-Labs/eigensdk-go/chainio/clients/avsregistry"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/eth"
	"github.com/Layr-Labs/eigensdk-go/crypto/bls"
	"github.com/Layr-Labs/eigensdk-go/logging"
	sdktypes "github.com/Layr-Labs/eigensdk-go/types"
	"github.com/Layr-Labs/eigensdk-go/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type IAVSReader interface {
	sdkavsregistry.AvsRegistryReader

	//GetQuorumsByBlockNumber
	GetQuorumsByBlockNumber(ctx context.Context, blockNumber uint32) (sdktypes.QuorumNums, error)

	//GetQuorumThresholdPercentages
	GetQuorumThresholdPercentages(ctx context.Context, blockNumber uint32, quorums sdktypes.QuorumNums) (sdktypes.QuorumThresholdPercentages, error)

	//
	GetOperatorsAvsStateAtBlock(ctx context.Context, quorumNumbers sdktypes.QuorumNums, blockNumber sdktypes.BlockNum) (map[sdktypes.OperatorId]sdktypes.OperatorAvsState, error)

	GetQuorumsAvsStateAtBlock(ctx context.Context, quorumNumbers sdktypes.QuorumNums, blockNumber sdktypes.BlockNum) (map[sdktypes.QuorumNum]sdktypes.QuorumAvsState, error)

	//Get Latest Request ID
	GetDeviationRequestId(ctx context.Context, agent common.Address, blockNumber uint32) (*big.Int, error)

	//Get Latest Request ID
	GetAgentSlashRequestId(ctx context.Context, agent common.Address, blockNumber uint32) (*big.Int, error)
}

type AvsReader struct {
	sdkavsregistry.AvsRegistryReader
	AvsServiceBindings *AVSManagerBinding
	graphClient        *graphql.GraphQLClient
}

func NewAvsReader(c config.AVSConfig, ethClient eth.Client, logger logging.Logger, graphClient *graphql.GraphQLClient) (IAVSReader, error) {
	regrAddr := c.RegistryCoordinatorAddr
	osAddr := c.OperatorStateRetrieverAddr

	avsManagerBindings, err := NewAvsManagerBindings(regrAddr, osAddr, ethClient)

	if err != nil {
		return nil, err
	}

	avsRegistryReader, err := sdkavsregistry.BuildAvsRegistryChainReader(regrAddr, osAddr, ethClient, logger)

	if err != nil {
		return nil, err
	}

	return &AvsReader{
		AvsServiceBindings: avsManagerBindings,
		AvsRegistryReader:  avsRegistryReader,
		graphClient:        graphClient,
	}, nil
}

func (r *AvsReader) GetQuorumsByBlockNumber(ctx context.Context, blockNumber uint32) (sdktypes.QuorumNums, error) {
	quorumCount, err := r.AvsRegistryReader.GetQuorumCount(&bind.CallOpts{
		Context:     ctx,
		BlockNumber: big.NewInt(int64(blockNumber)),
	})
	if err != nil {
		return nil, err
	}

	res := make([]sdktypes.QuorumNum, 0, quorumCount)
	for i := uint8(0); i < quorumCount; i++ {
		res = append(res, sdktypes.QuorumNum(i))
	}

	return res, nil
}

func (r *AvsReader) GetQuorumThresholdPercentages(ctx context.Context, blockNumber uint32, quorums sdktypes.QuorumNums) (sdktypes.QuorumThresholdPercentages, error) {
	quorumThresholdPercentage, err := r.AvsServiceBindings.ServiceManager.QuorumThresholdPercentage(&bind.CallOpts{
		Context:     ctx,
		BlockNumber: big.NewInt(int64(blockNumber)),
	})

	if err != nil {
		return nil, err
	}

	res := make([]sdktypes.QuorumThresholdPercentage, 0, len(quorums))
	for i := 0; i < len(quorums); i++ {
		res = append(res, sdktypes.QuorumThresholdPercentage(uint8(quorumThresholdPercentage)))
	}

	return res, nil
}

func (ar *AvsReader) GetOperatorsAvsStateAtBlock(ctx context.Context, quorumNumbers sdktypes.QuorumNums, blockNumber sdktypes.BlockNum) (map[sdktypes.OperatorId]sdktypes.OperatorAvsState, error) {
	operatorsAvsState := make(map[sdktypes.OperatorId]sdktypes.OperatorAvsState)
	// Get operator state for each quorum by querying BLSOperatorStateRetriever (this call is why this service implementation is called ChainCaller)
	operatorsStakesInQuorums, err := ar.AvsRegistryReader.GetOperatorsStakeInQuorumsAtBlock(&bind.CallOpts{Context: ctx}, quorumNumbers, blockNumber)
	if err != nil {
		return nil, utils.WrapError("Failed to get operator state", err)
	}
	numquorums := len(quorumNumbers)
	if len(operatorsStakesInQuorums) != numquorums {
		// ar.logger.Fatal("Number of quorums returned from GetOperatorsStakeInQuorumsAtBlock does not match number of quorums requested. Probably pointing to old contract or wrong implementation.", "service", "AvsRegistryServiceChainCaller")
	}

	for quorumIdx, quorumNum := range quorumNumbers {
		for _, operator := range operatorsStakesInQuorums[quorumIdx] {
			info, err := ar.graphClient.GetOperatorInfo(ctx, operator.OperatorId)
			if err != nil {
				return nil, utils.WrapError("Failed to find pubkeys for operator while building operatorsAvsState", err)
			}
			if operatorAvsState, ok := operatorsAvsState[operator.OperatorId]; ok {
				operatorAvsState.StakePerQuorum[quorumNum] = operator.Stake
			} else {
				stakePerQuorum := make(map[sdktypes.QuorumNum]sdktypes.StakeAmount)
				stakePerQuorum[quorumNum] = operator.Stake
				operatorsAvsState[operator.OperatorId] = sdktypes.OperatorAvsState{
					OperatorId:     operator.OperatorId,
					OperatorInfo:   info,
					StakePerQuorum: stakePerQuorum,
					BlockNumber:    blockNumber,
				}
				operatorsAvsState[operator.OperatorId].StakePerQuorum[quorumNum] = operator.Stake
			}
		}
	}

	return operatorsAvsState, nil
}

func (ar *AvsReader) GetQuorumsAvsStateAtBlock(ctx context.Context, quorumNumbers sdktypes.QuorumNums, blockNumber sdktypes.BlockNum) (map[sdktypes.QuorumNum]sdktypes.QuorumAvsState, error) {
	operatorsAvsState, err := ar.GetOperatorsAvsStateAtBlock(ctx, quorumNumbers, blockNumber)
	if err != nil {
		return nil, utils.WrapError("Failed to get quorum state", err)
	}
	quorumsAvsState := make(map[sdktypes.QuorumNum]sdktypes.QuorumAvsState)
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
		quorumsAvsState[quorumNum] = sdktypes.QuorumAvsState{
			QuorumNumber: quorumNum,
			AggPubkeyG1:  aggPubkeyG1,
			TotalStake:   totalStake,
			BlockNumber:  blockNumber,
		}
	}
	return quorumsAvsState, nil
}

func (r *AvsReader) GetDeviationRequestId(ctx context.Context, agent common.Address, blockNumber uint32) (*big.Int, error) {
	reqId, err := r.AvsServiceBindings.ServiceManager.AgentsDeviationRequestId(&bind.CallOpts{
		Context:     ctx,
		BlockNumber: big.NewInt(int64(blockNumber)),
	}, agent)

	if err != nil {
		return nil, err
	}

	return reqId, nil
}

func (r *AvsReader) GetAgentSlashRequestId(ctx context.Context, agent common.Address, blockNumber uint32) (*big.Int, error) {
	reqId, err := r.AvsServiceBindings.ServiceManager.AgentsFraudRequestId(&bind.CallOpts{
		Context:     ctx,
		BlockNumber: big.NewInt(int64(blockNumber)),
	}, agent)

	if err != nil {
		return nil, err
	}

	return reqId, nil
}

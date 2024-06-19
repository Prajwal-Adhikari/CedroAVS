package avs

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"

	"github.com/Cedro-Finance/CedroAVS/core/config"
	sdkavsregistry "github.com/Layr-Labs/eigensdk-go/chainio/clients/avsregistry"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/eth"
	"github.com/Layr-Labs/eigensdk-go/crypto/bls"
	"github.com/Layr-Labs/eigensdk-go/logging"
	sdktypes "github.com/Layr-Labs/eigensdk-go/types"
	"github.com/Layr-Labs/eigensdk-go/utils"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/shurcooL/graphql"
)

type IAVSReader interface {
	sdkavsregistry.AvsRegistryReader

	//GetQuorumsByBlockNumber
	GetQuorumsByBlockNumber(ctx context.Context, blockNumber uint32) (sdktypes.QuorumNums, error)

	//GetQuorumThresholdPercentages
	GetQuorumThresholdPercentages(ctx context.Context, blockNumber uint32, quorums sdktypes.QuorumNums) (sdktypes.QuorumThresholdPercentages, error)
}

type AvsReader struct {
	sdkavsregistry.AvsRegistryReader
	AvsServiceBindings *AVSManagerBinding
}

// Define the GraphQL query structure
type OperatorQuery struct {
	Operator struct {
		Operator   graphql.String
		PubkeyG1_X graphql.String   `graphql:"pubkeyG1_X"`
		PubkeyG1_Y graphql.String   `graphql:"pubkeyG1_Y"`
		PubkeyG2_X []graphql.String `graphql:"pubkeyG2_X"`
		PubkeyG2_Y []graphql.String `graphql:"pubkeyG2_Y"`
		ID         graphql.String   `graphql:"id"`
	} `graphql:"operator(id: $id)"`
}

func NewAvsReader(c config.AVSConfig, ethClient eth.Client, logger logging.Logger) (*AvsReader, error) {
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
			info, err := ar.getOperatorInfo(ctx, operator.OperatorId)
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

func (ar *AvsReader) getOperatorInfo(ctx context.Context, opId sdktypes.OperatorId) (sdktypes.OperatorInfo, error) {

	fmt.Println("I am at graphql client")
	fmt.Println(opId)

	operatorId := hex.EncodeToString(opId[:])

	fmt.Println(operatorId)
	fmt.Println(fmt.Sprintf("0x%s", operatorId))

	// Create a new GraphQL client with the specified endpoint
	client := graphql.NewClient("https://api.studio.thegraph.com/query/41099/cedro-avs/v0.0.3h", nil)

	// Define the query variables
	vars := map[string]interface{}{
		"id": graphql.String(fmt.Sprintf("0x%s", operatorId)),
	}

	// Create an instance of the query struct to hold the response data
	var query OperatorQuery

	// Execute the query
	err := client.Query(ctx, &query, vars)
	if err != nil {
		log.Fatalf("failed to execute query: %v", err)
	}

	opInfo, err := convertIndexedOperatorInfoGqlToOperatorInfo(&query)

	if err != nil {
		log.Fatal("Errro converting gql result", err)
	}

	return *opInfo, nil
}

func convertIndexedOperatorInfoGqlToOperatorInfo(operator *OperatorQuery) (*sdktypes.OperatorInfo, error) {

	pubkeyG1 := new(bn254.G1Affine)
	_, err := pubkeyG1.X.SetString(string(operator.Operator.PubkeyG1_X))
	if err != nil {
		return nil, err
	}
	_, err = pubkeyG1.Y.SetString(string(operator.Operator.PubkeyG1_Y))
	if err != nil {
		return nil, err
	}

	pubkeyG2 := new(bn254.G2Affine)
	_, err = pubkeyG2.X.A1.SetString(string(operator.Operator.PubkeyG2_X[0]))
	if err != nil {
		return nil, err
	}
	_, err = pubkeyG2.X.A0.SetString(string(operator.Operator.PubkeyG2_X[1]))
	if err != nil {
		return nil, err
	}
	_, err = pubkeyG2.Y.A1.SetString(string(operator.Operator.PubkeyG2_Y[0]))
	if err != nil {
		return nil, err
	}
	_, err = pubkeyG2.Y.A0.SetString(string(operator.Operator.PubkeyG2_Y[1]))
	if err != nil {
		return nil, err
	}

	return &sdktypes.OperatorInfo{
		Socket: sdktypes.Socket(""),
		Pubkeys: sdktypes.OperatorPubkeys{
			G1Pubkey: &bls.G1Point{G1Affine: pubkeyG1},
			G2Pubkey: &bls.G2Point{G2Affine: pubkeyG2},
		},
	}, nil

}

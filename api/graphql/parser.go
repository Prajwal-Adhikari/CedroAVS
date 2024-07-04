package graphql

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"

	"github.com/Layr-Labs/eigensdk-go/crypto/bls"
	sdktypes "github.com/Layr-Labs/eigensdk-go/types"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/ethereum/go-ethereum/common"
)

func ConvertIndexedOperatorInfoGqlToOperatorInfo(operator *OperatorQuery) (*sdktypes.OperatorInfo, error) {
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

// type AgentPodsInfo struct {
// 	Agent         common.Address
// 	AgentContract common.Address
// 	ChainId       string
// }

// func ConvertAgentPodsInfo(query *AgentPodMappingQuery) (AgentPodsInfo, error) {
// 	chainId := string(query.AgentPods.ChainId)
// 	chainIdBytes, err := hex.DecodeString(chainId[2:])

// 	if err != nil {
// 		log.Println("Error")
// 	}

// 	chainName := string(chainIdBytes)

// 	podInfo := AgentPodsInfo{
// 		Agent:         common.HexToAddress(string(query.AgentPods.Agent)),
// 		AgentContract: common.HexToAddress(string(query.AgentPods.AgentContract)),
// 		ChainId:       strings.Trim(chainName, "\x00"),
// 	}

// 	return podInfo, nil

// }

type AgentPodChainListInfo struct {
	Agent    common.Address
	ChainIds []string
}

func ConvertAgentPodsChains(query *AgentPodChainListsQuery) (AgentPodChainListInfo, error) {

	if len(query.AgendPodChains) == 0 {
		return AgentPodChainListInfo{}, fmt.Errorf("no chains found for given agent")
	}

	chainIds := query.AgendPodChains[0].ChainIDs
	parsedChainIds := make([]string, 0)

	for _, cId := range chainIds {
		chainIdBytes, err := hex.DecodeString(string(cId)[2:])

		if err != nil {
			log.Println("Error")
		}

		chainName := strings.Trim(string(chainIdBytes), "\x00")

		parsedChainIds = append(parsedChainIds, chainName)

	}

	podInfo := AgentPodChainListInfo{
		Agent:    common.HexToAddress(string(query.AgendPodChains[0].Agent)),
		ChainIds: parsedChainIds,
	}

	return podInfo, nil
}

type EnrolledAgent struct {
	Agent common.Address
}

func ConvertEnrolledAgent(query *EnrolledAgentsListQuery) ([]common.Address, error) {

	agents := make([]common.Address, 0)

	for _, agent := range query.AgendPodChains {
		agents = append(agents, common.HexToAddress(string(agent.Agent)))
	}

	return agents, nil
}

type DepositInfo struct {
	IsClaimed       bool
	InterestRate    *big.Int
	TxnID           *big.Int
	PrincipalAmount *big.Int
	BlockTimestamp  uint64
	TimeInterval    *big.Int
	AgentID         common.Address
	BlockNumber     uint64
}

func ConvertDepositListInfo(deposits []DeposisInfoQuery) ([]DepositInfo, error) {

	depositList := make([]DepositInfo, 0)

	for _, deposit := range deposits {

		interestRate, ok := new(big.Int).SetString(string(deposit.InterestRate), 10)

		if !ok {
			return []DepositInfo{}, nil
		}

		isClaimed := bool(deposit.IsClaimed)
		txId, ok := new(big.Int).SetString(string(deposit.TxnID), 10)

		if !ok {
			return []DepositInfo{}, nil
		}

		principalAmount, ok := new(big.Int).SetString(string(deposit.PrincipalAmount), 10)

		if !ok {
			return []DepositInfo{}, nil
		}

		timeInterval, ok := new(big.Int).SetString(string(deposit.TimeInterval), 10)

		if !ok {
			return []DepositInfo{}, nil
		}

		blockTimestamp, err := strconv.ParseUint(string(deposit.BlockTimestamp), 10, 64)

		if err != nil {
			return []DepositInfo{}, err
		}

		blockNumber, err := strconv.ParseUint(string(deposit.BlockNumber), 10, 64)

		if err != nil {
			return []DepositInfo{}, err
		}

		dep := DepositInfo{
			IsClaimed:       isClaimed,
			InterestRate:    interestRate,
			PrincipalAmount: principalAmount,
			BlockTimestamp:  blockTimestamp,
			AgentID:         common.HexToAddress(string(deposit.Agent)),
			TxnID:           txId,
			TimeInterval:    timeInterval,
			BlockNumber:     blockNumber,
		}

		depositList = append(depositList, dep)
	}

	return depositList, nil
}

type AgentSummary struct {
	AgentContract     common.Address
	Agent             common.Address
	ChainId           string
	TotalDepositCount int
	TotalClaimCount   int
}

func ConvertAgentSummaryInfo(query *AgentSummaryQuery) (AgentSummary, error) {
	agentAddress := common.HexToAddress(string(query.AgentSummary.Agent))
	agentPod := common.HexToAddress(string(query.AgentSummary.AgentContract))

	chainId := string(query.AgentSummary.ChainId)

	chainIdBytes, err := hex.DecodeString(chainId[2:])

	if err != nil {
		log.Println("Error")
	}

	chainName := string(chainIdBytes)

	depositCount, err := strconv.Atoi(string(query.AgentSummary.TotalDepositCount))

	if err != nil {
		return AgentSummary{}, err
	}

	claimCount, err := strconv.Atoi(string(query.AgentSummary.TotalClaimCount))

	if err != nil {
		return AgentSummary{}, err
	}

	return AgentSummary{
		AgentContract:     agentPod,
		Agent:             agentAddress,
		TotalDepositCount: depositCount,
		TotalClaimCount:   claimCount,
		ChainId:           strings.Trim(chainName, "\x00"),
	}, nil
}

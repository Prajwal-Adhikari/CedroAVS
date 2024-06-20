package graphql

import (
	"encoding/hex"
	"fmt"
	"log"
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

type AgentPodsInfo struct {
	Agent         common.Address
	AgentContract common.Address
	ChainId       string
}

func ConvertAgentPodsInfo(query *AgentPodMappingQuery) (AgentPodsInfo, error) {
	fmt.Println(query)
	chainId := string(query.AgentPods.ChainId)
	chainIdBytes, err := hex.DecodeString(chainId[2:])

	if err != nil {
		log.Println("Error")
	}

	chainName := string(chainIdBytes)

	podInfo := AgentPodsInfo{
		Agent:         common.HexToAddress(string(query.AgentPods.Agent)),
		AgentContract: common.HexToAddress(string(query.AgentPods.AgentContract)),
		ChainId:       strings.Trim(chainName, "\x00"),
	}

	return podInfo, nil

}

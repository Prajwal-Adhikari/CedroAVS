package tasks

import (
	"encoding/hex"
	"fmt"
	"math/big"

	aggregator "github.com/Cedro-Finance/CedroAVS/api/grpc"
	"github.com/Cedro-Finance/CedroAVS/core"
	cedroservicemanager "github.com/Cedro-Finance/CedroAVS/core/bindings/CedroServiceManager"
	"github.com/Layr-Labs/eigensdk-go/crypto/bls"
	sdktypes "github.com/Layr-Labs/eigensdk-go/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/crypto/sha3"
)

type BalanceTaskInfo struct {
	QuorumNumbers              sdktypes.QuorumNums
	QuorumThresholdPercentages sdktypes.QuorumThresholdPercentages
	TaskIndex                  uint32
	ReferenceBlockNumber       uint64

	ReqId         *big.Int
	ActualBalance *big.Int
	IdealBalance  *big.Int
	AgentAddress  common.Address
}

var _ ITaskInfo = (*BalanceTaskInfo)(nil)

func NewBalanceTaskInfo(req *aggregator.BalanceTaskInfo) (*BalanceTaskInfo, error) {

	actualBalance := new(big.Int)
	idealBalance := new(big.Int)
	reqId := new(big.Int)

	actualBalance = actualBalance.SetBytes(req.GetActualBalance())
	idealBalance = idealBalance.SetBytes(req.GetIdealBalance())
	reqId = reqId.SetBytes(req.GetReqId())

	res := &BalanceTaskInfo{
		QuorumNumbers:              core.ConvertQuorumNumbersFromBytes(req.GetQuorumNumbers()),
		QuorumThresholdPercentages: core.ConvertQuorumThresholdPercentagesFromBytes(req.GetQuorumThresholdPercentages()),
		TaskIndex:                  req.GetTaskIndex(),
		ReferenceBlockNumber:       req.GetReferenceBlockNumber(),
		ActualBalance:              actualBalance,
		IdealBalance:               idealBalance,
		ReqId:                      reqId,
		AgentAddress:               common.HexToAddress(req.GetAgentAddress()),
	}

	return res, nil
}

func (r *BalanceTaskInfo) ToPbType() *aggregator.BalanceTaskInfo {
	return &aggregator.BalanceTaskInfo{
		QuorumNumbers:              r.QuorumNumbers.UnderlyingType(),
		QuorumThresholdPercentages: r.QuorumThresholdPercentages.UnderlyingType(),
		TaskIndex:                  r.TaskIndex,
		ReferenceBlockNumber:       r.ReferenceBlockNumber,
		ActualBalance:              r.ActualBalance.Bytes(),
		IdealBalance:               r.IdealBalance.Bytes(),
		ReqId:                      r.ReqId.Bytes(),
		AgentAddress:               r.AgentAddress.Hex(),
	}
}

func (r *BalanceTaskInfo) ToCedroServiceManagerHeader() cedroservicemanager.CedroServiceDataDeviationBalanceHeader {
	return cedroservicemanager.CedroServiceDataDeviationBalanceHeader{
		AgentAddress:               r.AgentAddress,
		ReferenceBlockNumber:       uint32(r.ReferenceBlockNumber),
		QuorumNumbers:              r.QuorumNumbers.UnderlyingType(),
		QuorumThresholdPercentages: r.QuorumThresholdPercentages.UnderlyingType(),
		RequestId:                  r.ReqId,
		ActualAmount:               r.ActualBalance,
		IdealAmount:                r.IdealBalance,
	}
}

func (r *BalanceTaskInfo) EncodeSignHash() ([]byte, error) {

	tp, err := abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{
			Name: "agentAddress",
			Type: "address",
		},
		{
			Name: "referenceBlockNumber",
			Type: "uint32",
		},
		{
			Name: "requestId",
			Type: "uint256",
		},
		{
			Name: "actualAmount",
			Type: "uint256",
		},
		{
			Name: "idealAmount",
			Type: "uint256",
		},
	})

	if err != nil {
		return nil, err
	}

	s := struct {
		AgentAddress         common.Address
		ReferenceBlockNumber uint32
		RequestId            *big.Int
		ActualAmount         *big.Int
		IdealAmount          *big.Int
	}{
		AgentAddress:         r.AgentAddress,
		ReferenceBlockNumber: uint32(r.ReferenceBlockNumber),
		ActualAmount:         r.ActualBalance,
		IdealAmount:          r.IdealBalance,
		RequestId:            r.ReqId,
	}

	arguments := abi.Arguments{
		{
			Type: tp,
		},
	}

	bytes, err := arguments.Pack(s)

	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func (r *BalanceTaskInfo) GetHeaderHash() (sdktypes.Bytes32, error) {
	headerByte, err := r.EncodeSignHash()
	if err != nil {
		return sdktypes.Bytes32{}, err
	}

	var headerHash sdktypes.Bytes32
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(headerByte)
	copy(headerHash[:], hasher.Sum(nil)[:32])

	return headerHash.UnderlyingType(), nil
}

func (r *BalanceTaskInfo) GetTaskIndex() sdktypes.TaskIndex {
	return r.TaskIndex
}

func (r *BalanceTaskInfo) GetUniqueHash() string {
	return r.AgentAddress.Hex() + r.ReqId.String()
}

type CreateDeviationTaskRequest struct {
	ReqId                *big.Int
	ReferenceBlockNumber uint64
	AgentAddress         common.Address
}

func NewCreateDeviationTaskRequest(req *aggregator.CreateDeviationTaskRequest) (*CreateDeviationTaskRequest, error) {

	z := new(big.Int)

	reqId := z.SetBytes(req.GetReqId())

	res := &CreateDeviationTaskRequest{
		ReqId:                reqId,
		ReferenceBlockNumber: req.GetReferenceBlockNumber(),
		AgentAddress:         common.HexToAddress(req.GetAgentAddress()),
	}

	return res, nil
}

func (t *CreateDeviationTaskRequest) GetHash() (string, error) {
	hash := t.AgentAddress.Hex() + t.ReqId.String()

	return hash, nil
}

func (t *CreateDeviationTaskRequest) GetAgent() common.Address {
	return t.AgentAddress
}
func (t *CreateDeviationTaskRequest) GetIdentifier() *big.Int {
	return t.ReqId
}
func (t *CreateDeviationTaskRequest) GetRefBlock() uint64 {
	return t.ReferenceBlockNumber
}

func (t *CreateDeviationTaskRequest) GetTaskInfo(taskIndex sdktypes.TaskIndex, quorumNums sdktypes.QuorumNums, quorumThresholds sdktypes.QuorumThresholdPercentages) ITaskInfo {
	return &BalanceTaskInfo{
		QuorumNumbers:              quorumNums,
		QuorumThresholdPercentages: quorumThresholds,
		TaskIndex:                  taskIndex,
		ReferenceBlockNumber:       t.ReferenceBlockNumber,
		ReqId:                      t.ReqId,
		ActualBalance:              big.NewInt(0),
		IdealBalance:               big.NewInt(0),
		AgentAddress:               t.AgentAddress,
	}
}

type CreateDeviationTaskResponse struct {
	Info BalanceTaskInfo
}

func (r CreateDeviationTaskResponse) ToPbType() *aggregator.CreateDeviationTaskResponse {
	return &aggregator.CreateDeviationTaskResponse{
		Info: r.Info.ToPbType(),
	}
}

type SignedDeviationTaskRespRequest struct {
	Balance      BalanceTaskInfo
	BlsSignature bls.Signature
	OperatorId   sdktypes.OperatorId
}

func NewSignedDeviationTaskRespRequest(req *aggregator.SignedDeviationTaskRespRequest) (*CreateSignedTaskResRequest, error) {
	operatorId := req.GetOperatorId()

	if len(operatorId) != 32 {
		return nil, fmt.Errorf("operator ID len should be 32, got %d", len(operatorId))
	}

	balHead, err := NewBalanceTaskInfo(req.GetBalanceHeader())

	if err != nil {
		return nil, fmt.Errorf("new alert task info failed: %v", err.Error())
	}

	signRaw := req.GetOperatorRequestSignature()
	if len(signRaw) != 64 {
		return nil, fmt.Errorf("operatorRequestSignature len should be 64")
	}

	g1Point := bls.NewZeroG1Point().Deserialize(signRaw)
	sign := bls.Signature{G1Point: g1Point}

	res := &CreateSignedTaskResRequest{
		Task:         balHead,
		BlsSignature: sign,
	}

	copy(res.OperatorId[:], operatorId[:32])

	return res, nil
}

func GetRequestHash(reqId *big.Int, agentAddr common.Address) (string, error) {
	reqType, err := abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{
			Name: "reqId",
			Type: "uint256",
		},
		{
			Name: "agentAddress",
			Type: "address",
		},
	})

	if err != nil {
		return "", err
	}

	s := struct {
		ReqId        *big.Int
		AgentAddress common.Address
	}{
		ReqId:        reqId,
		AgentAddress: agentAddr,
	}

	arguments := abi.Arguments{
		{
			Type: reqType,
		},
	}

	bytes, err := arguments.Pack(s)

	if err != nil {
		return "", err
	}

	hasher := sha3.NewLegacyKeccak256()
	var headerHash [32]byte
	hasher.Write(bytes)
	copy(headerHash[:], hasher.Sum(nil)[:32])

	return hex.EncodeToString(headerHash[:]), nil
}

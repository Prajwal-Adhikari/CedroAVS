package core

import (
	"encoding/hex"
	"fmt"
	"math/big"

	// "fmt"
	"log/slog"

	// aggregator "github.com/Cedro-Finance/CedroAVS/api/grpc"
	aggregator "github.com/Cedro-Finance/CedroAVS/api/grpc"
	cedroservicemanager "github.com/Cedro-Finance/CedroAVS/core/bindings/CedroServiceManager"
	"github.com/Layr-Labs/eigensdk-go/crypto/bls"
	sdktypes "github.com/Layr-Labs/eigensdk-go/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/crypto/sha3"
)

type Bytes32 [32]byte

func (m Bytes32) LogValue() slog.Value {
	return slog.StringValue(m.String())
}

func (m Bytes32) String() string {
	return hex.EncodeToString(m[:])
}

func (m *Bytes32) UnderlyingType() [32]byte {
	return *m
}

type BalanceTaskInfo struct {
	ReqId                      *big.Int
	QuorumNumbers              sdktypes.QuorumNums
	QuorumThresholdPercentages sdktypes.QuorumThresholdPercentages
	TaskIndex                  uint32
	ActualBalance              *big.Int
	IdealBalance               *big.Int
	ReferenceBlockNumber       uint64
	AgentAddress               common.Address
}

func NewBalanceTaskInfo(req *aggregator.BalanceTaskInfo) (*BalanceTaskInfo, error) {

	actualBalance := new(big.Int)
	idealBalance := new(big.Int)
	reqId := new(big.Int)

	actualBalance = actualBalance.SetBytes(req.GetActualBalance())
	idealBalance = idealBalance.SetBytes(req.GetIdealBalance())
	reqId = reqId.SetBytes(req.GetReqId())

	res := &BalanceTaskInfo{
		QuorumNumbers:              ConvertQuorumNumbersFromBytes(req.GetQuorumNumbers()),
		QuorumThresholdPercentages: ConvertQuorumThresholdPercentagesFromBytes(req.GetQuorumThresholdPercentages()),
		TaskIndex:                  req.GetTaskIndex(),
		ReferenceBlockNumber:       req.GetReferenceBlockNumber(),
		ActualBalance:              actualBalance,
		IdealBalance:               idealBalance,
		ReqId:                      reqId,
		AgentAddress:               common.HexToAddress(req.GetAgentAddress()),
	}

	return res, nil
}

func (r BalanceTaskInfo) ToPbType() *aggregator.BalanceTaskInfo {

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

func (r BalanceTaskInfo) ToCedroServiceManagerBalanceHeader() cedroservicemanager.DATADeviationBalanceHeader {
	return cedroservicemanager.DATADeviationBalanceHeader{
		AgentAddress:               r.AgentAddress,
		QuorumNumbers:              r.QuorumNumbers.UnderlyingType(),
		QuorumThresholdPercentages: r.QuorumThresholdPercentages.UnderlyingType(),
		ReferenceBlockNumber:       uint32(r.ReferenceBlockNumber),
		ActualAmount:               r.ActualBalance,
		IdealAmount:                r.IdealBalance,
		RequestId:                  r.ReqId,
	}
}

func (r *BalanceTaskInfo) EncodeSignHash() ([]byte, error) {

	tp, err := abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{
			Name: "agentAddress",
			Type: "address",
		},
		{
			Name: "actualAmount",
			Type: "uint256",
		},
		{
			Name: "idealAmount",
			Type: "uint256",
		},
		{
			Name: "requestId",
			Type: "uint256",
		},
		{
			Name: "referenceBlockNumber",
			Type: "uint32",
		},
	})

	if err != nil {
		return nil, err
	}

	s := struct {
		AgentAddress         common.Address
		RequestId            *big.Int
		ActualAmount         *big.Int
		IdealAmount          *big.Int
		ReferenceBlockNumber uint32
	}{
		AgentAddress:         r.AgentAddress,
		ActualAmount:         r.ActualBalance,
		IdealAmount:          r.IdealBalance,
		RequestId:            r.ReqId,
		ReferenceBlockNumber: uint32(r.ReferenceBlockNumber),
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

func (r BalanceTaskInfo) GetBalanceHeaderHash() ([32]byte, error) {
	headerByte, err := r.EncodeSignHash()
	if err != nil {
		return [32]byte{}, err
	}

	var headerHash [32]byte
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(headerByte)
	copy(headerHash[:], hasher.Sum(nil)[:32])

	return headerHash, nil
}

// The init operator request
type InitOperatorRequest struct {
	OperatorId                 sdktypes.OperatorId
	OperatorAddress            common.Address
	OperatorStateRetrieverAddr common.Address
	RegistryCoordinatorAddr    common.Address
}

func NewInitOperatorRequest(req *aggregator.InitOperatorRequest) (*InitOperatorRequest, error) {
	operatorId := req.GetOperatorId()
	if len(operatorId) != 32 {
		return nil, fmt.Errorf("operator ID len should be 32, got %d", len(operatorId))
	}

	if !common.IsHexAddress(req.GetOperatorAddress()) {
		return nil, fmt.Errorf("operatorAddress not a hex address")
	}
	operatorAddress := common.HexToAddress(req.OperatorAddress)

	if !common.IsHexAddress(req.GetOperatorStateRetrieverAddr()) {
		return nil, fmt.Errorf("operatorStateRetrieverAddr not a hex address")
	}
	operatorStateRetrieverAddr := common.HexToAddress(req.GetOperatorStateRetrieverAddr())

	if !common.IsHexAddress(req.GetRegistryCoordinatorAddr()) {
		return nil, fmt.Errorf("registryCoordinatorAddr not a hex address")
	}
	registryCoordinatorAddr := common.HexToAddress(req.GetRegistryCoordinatorAddr())

	res := &InitOperatorRequest{
		OperatorId:                 sdktypes.Bytes32(req.OperatorId),
		OperatorAddress:            operatorAddress,
		OperatorStateRetrieverAddr: operatorStateRetrieverAddr,
		RegistryCoordinatorAddr:    registryCoordinatorAddr,
	}

	copy(res.OperatorId[:], operatorId[:32])

	return res, nil
}

// The init operator response
type InitOperatorResponse struct {
	Ok  bool
	Res string
}

func (r InitOperatorResponse) ToPbType() *aggregator.InitOperatorResponse {
	return &aggregator.InitOperatorResponse{
		Ok:     r.Ok,
		Reason: r.Res,
	}
}

type CreateTaskRequest struct {
	ReqId                *big.Int
	ReferenceBlockNumber uint64
	AgentAddress         common.Address
}

func NewCreateTaskRequest(req *aggregator.CreateTaskRequest) (*CreateTaskRequest, error) {

	z := new(big.Int)

	reqId := z.SetBytes(req.GetReqId())

	res := &CreateTaskRequest{
		ReqId:                reqId,
		ReferenceBlockNumber: req.GetReferenceBlockNumber(),
		AgentAddress:         common.HexToAddress(req.GetAgentAddress()),
	}

	return res, nil
}

type CreateTaskResponse struct {
	Info BalanceTaskInfo
}

func (r CreateTaskResponse) ToPbType() *aggregator.CreateTaskResponse {
	return &aggregator.CreateTaskResponse{
		Info: r.Info.ToPbType(),
	}
}

type SignedTaskRespRequest struct {
	Balance      BalanceTaskInfo
	BlsSignature bls.Signature
	OperatorId   sdktypes.OperatorId
}

func NewSignedTaskRespRequest(req *aggregator.SignedTaskRespRequest) (*SignedTaskRespRequest, error) {
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

	res := &SignedTaskRespRequest{
		Balance:      *balHead,
		BlsSignature: sign,
	}

	copy(res.OperatorId[:], operatorId[:32])

	return res, nil
}

type SignedTaskRespResponse struct {
	Reply  bool
	TxHash [32]byte
}

func (r SignedTaskRespResponse) ToPbType() *aggregator.SignedTaskRespResponse {
	return &aggregator.SignedTaskRespResponse{
		Reply:  r.Reply,
		TxHash: r.TxHash[:],
	}
}

func GetRequestHash(reqId *big.Int, agentAddr common.Address) ([]byte, error) {
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
		return nil, err
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
		return nil, err
	}

	hasher := sha3.NewLegacyKeccak256()
	var headerHash [32]byte
	hasher.Write(bytes)
	copy(headerHash[:], hasher.Sum(nil)[:32])

	return headerHash[:], nil
}

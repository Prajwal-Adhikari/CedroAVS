package core

import (
	"encoding/hex"
	"fmt"
	"math/big"

	// "fmt"

	// aggregator "github.com/Cedro-Finance/CedroAVS/api/grpc"
	aggregator "github.com/Cedro-Finance/CedroAVS/api/grpc"
	cedroservicemanager "github.com/Cedro-Finance/CedroAVS/core/bindings/CedroServiceManager"
	"github.com/Layr-Labs/eigensdk-go/crypto/bls"
	sdktypes "github.com/Layr-Labs/eigensdk-go/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/crypto/sha3"
)

// type TaskInfo struct {
// 	QuorumNumbers              sdktypes.QuorumNums
// 	QuorumThresholdPercentages sdktypes.QuorumThresholdPercentages
// 	TaskIndex                  uint32
// 	ReferenceBlockNumber       uint64
// }

const (
	AgentDeviationTask = "deviation"
	AgentSlashTask     = "agentslash"
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

func (r BalanceTaskInfo) ToCedroServiceManagerHeader() cedroservicemanager.CedroServiceDataDeviationBalanceHeader {
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

func (r BalanceTaskInfo) GetHeaderHash() (Bytes32, error) {
	headerByte, err := r.EncodeSignHash()
	if err != nil {
		return Bytes32{}, err
	}

	var headerHash Bytes32
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(headerByte)
	copy(headerHash[:], hasher.Sum(nil)[:32])

	return headerHash.UnderlyingType(), nil
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

	if !common.IsHexAddress(req.GetOperatorStateRetrieverAddress()) {
		return nil, fmt.Errorf("operatorStateRetrieverAddr not a hex address")
	}
	operatorStateRetrieverAddr := common.HexToAddress(req.GetOperatorStateRetrieverAddress())

	if !common.IsHexAddress(req.GetRegistryCoordinatorAddress()) {
		return nil, fmt.Errorf("registryCoordinatorAddr not a hex address")
	}
	registryCoordinatorAddr := common.HexToAddress(req.GetRegistryCoordinatorAddress())

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
	hash, err := GetRequestHash(t.ReqId, t.AgentAddress)

	return hash, err
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

func NewSignedDeviationTaskRespRequest(req *aggregator.SignedDeviationTaskRespRequest) (*SignedDeviationTaskRespRequest, error) {
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

	res := &SignedDeviationTaskRespRequest{
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

//Agent Slash Task Datas

type AgentSlashTaskInfo struct {
	QuorumNumbers              sdktypes.QuorumNums
	QuorumThresholdPercentages sdktypes.QuorumThresholdPercentages
	TaskIndex                  uint32
	InsufficientBalance        *big.Int
	ReferenceBlockNumber       uint64
	AgentAddress               common.Address
	TxnId                      *big.Int
}

func NewAgentSlashTaskInfo(req *aggregator.AgentSlashTaskInfo) (*AgentSlashTaskInfo, error) {

	insufBalance := new(big.Int)
	txnId := new(big.Int)

	insufBalance = insufBalance.SetBytes(req.GetInsufficientBalance())
	txnId = txnId.SetBytes(req.GetTxnId())

	res := &AgentSlashTaskInfo{
		QuorumNumbers:              ConvertQuorumNumbersFromBytes(req.GetQuorumNumbers()),
		QuorumThresholdPercentages: ConvertQuorumThresholdPercentagesFromBytes(req.GetQuorumThresholdPercentages()),
		TaskIndex:                  req.GetTaskIndex(),
		ReferenceBlockNumber:       req.GetReferenceBlockNumber(),
		AgentAddress:               common.HexToAddress(req.GetAgentAddress()),
		InsufficientBalance:        insufBalance,
		TxnId:                      txnId,
	}

	return res, nil
}

func (r AgentSlashTaskInfo) ToPbType() *aggregator.AgentSlashTaskInfo {

	return &aggregator.AgentSlashTaskInfo{
		QuorumNumbers:              r.QuorumNumbers.UnderlyingType(),
		QuorumThresholdPercentages: r.QuorumThresholdPercentages.UnderlyingType(),
		TaskIndex:                  r.TaskIndex,
		ReferenceBlockNumber:       r.ReferenceBlockNumber,
		AgentAddress:               r.AgentAddress.Hex(),
		InsufficientBalance:        r.InsufficientBalance.Bytes(),
		TxnId:                      r.TxnId.Bytes(),
	}
}

func (r AgentSlashTaskInfo) ToCedroServiceManagerHeader() cedroservicemanager.CedroServiceDataFraudBalanceHeader {
	return cedroservicemanager.CedroServiceDataFraudBalanceHeader{
		AgentAddress:               r.AgentAddress,
		ReferenceBlockNumber:       uint32(r.ReferenceBlockNumber),
		QuorumNumbers:              r.QuorumNumbers.UnderlyingType(),
		QuorumThresholdPercentages: r.QuorumThresholdPercentages.UnderlyingType(),
		TransactionId:              r.TxnId,
		InsufficientAmount:         r.InsufficientBalance,

	}
}

func (r *AgentSlashTaskInfo) EncodeSignHash() ([]byte, error) {

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
			Name: "transactionId",
			Type: "uint256",
		},
		{
			Name: "insufficientAmount",
			Type: "uint256",
		},
	})

	if err != nil {
		return nil, err
	}

	s := struct {
		AgentAddress         common.Address
		ReferenceBlockNumber uint32
		TransactionId        *big.Int
		InsufficientAmount   *big.Int
	}{
		AgentAddress:         r.AgentAddress,
		ReferenceBlockNumber: uint32(r.ReferenceBlockNumber),
		TransactionId:        r.TxnId,
		InsufficientAmount:   r.InsufficientBalance,
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

func (r AgentSlashTaskInfo) GetHeaderHash() (Bytes32, error) {
	headerByte, err := r.EncodeSignHash()
	if err != nil {
		return Bytes32{}, err
	}

	var headerHash Bytes32
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(headerByte)
	copy(headerHash[:], hasher.Sum(nil)[:32])

	return headerHash.UnderlyingType(), nil
}

type CreateAgentSlashTaskRequest struct {
	TxnId                *big.Int
	ReferenceBlockNumber uint64
	AgentAddress         common.Address
}

func NewCreateAgentSlashTaskRequest(req *aggregator.CreateAgentSlashTaskRequest) (*CreateAgentSlashTaskRequest, error) {

	z := new(big.Int)

	reqId := z.SetBytes(req.GetTxnId())
	res := &CreateAgentSlashTaskRequest{
		ReferenceBlockNumber: req.GetReferenceBlockNumber(),
		AgentAddress:         common.HexToAddress(req.GetAgentAddress()),
		TxnId:                reqId,
	}

	return res, nil
}

func (t *CreateAgentSlashTaskRequest) GetHash() (string, error) {
	hash, err := GetRequestHash(t.TxnId, t.AgentAddress)

	return hash, err
}

type CreateAgentSlashTaskResponse struct {
	Info AgentSlashTaskInfo
}

func (r CreateAgentSlashTaskResponse) ToPbType() *aggregator.CreateAgentSlashTaskResponse {
	return &aggregator.CreateAgentSlashTaskResponse{
		Info: r.Info.ToPbType(),
	}
}

type SignedAgentSlashTaskRespRequest struct {
	Balance      AgentSlashTaskInfo
	BlsSignature bls.Signature
	OperatorId   sdktypes.OperatorId
}

func NewSignedAgentSlashTaskRespRequest(req *aggregator.SignedAgentSlashTaskRespRequest) (*SignedAgentSlashTaskRespRequest, error) {
	operatorId := req.GetOperatorId()

	if len(operatorId) != 32 {
		return nil, fmt.Errorf("operator ID len should be 32, got %d", len(operatorId))
	}

	balHead, err := NewAgentSlashTaskInfo(req.GetFraudHeader())

	if err != nil {
		return nil, fmt.Errorf("new alert task info failed: %v", err.Error())
	}

	signRaw := req.GetOperatorRequestSignature()
	if len(signRaw) != 64 {
		return nil, fmt.Errorf("operatorRequestSignature len should be 64")
	}

	g1Point := bls.NewZeroG1Point().Deserialize(signRaw)
	sign := bls.Signature{G1Point: g1Point}

	res := &SignedAgentSlashTaskRespRequest{
		Balance:      *balHead,
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

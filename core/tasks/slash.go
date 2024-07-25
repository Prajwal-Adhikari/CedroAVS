package tasks

import (
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
		QuorumNumbers:              core.ConvertQuorumNumbersFromBytes(req.GetQuorumNumbers()),
		QuorumThresholdPercentages: core.ConvertQuorumThresholdPercentagesFromBytes(req.GetQuorumThresholdPercentages()),
		TaskIndex:                  req.GetTaskIndex(),
		ReferenceBlockNumber:       req.GetReferenceBlockNumber(),
		AgentAddress:               common.HexToAddress(req.GetAgentAddress()),
		InsufficientBalance:        insufBalance,
		TxnId:                      txnId,
	}

	return res, nil
}

func (r *AgentSlashTaskInfo) ToPbType() *aggregator.AgentSlashTaskInfo {

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

func (r *AgentSlashTaskInfo) ToCedroServiceManagerHeader() cedroservicemanager.CedroServiceDataFraudBalanceHeader {
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

func (r *AgentSlashTaskInfo) GetHeaderHash() (sdktypes.Bytes32, error) {
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

func (r *AgentSlashTaskInfo) GetTaskIndex() sdktypes.TaskIndex {
	return r.TaskIndex
}

func (r *AgentSlashTaskInfo) GetUniqueHash() string {
	return r.AgentAddress.Hex() + r.TxnId.String()
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
	hash := t.AgentAddress.Hex() + t.TxnId.String()

	return hash, nil
}

func (t *CreateAgentSlashTaskRequest) GetAgent() common.Address {
	return t.AgentAddress
}
func (t *CreateAgentSlashTaskRequest) GetIdentifier() *big.Int {
	return t.TxnId
}
func (t *CreateAgentSlashTaskRequest) GetRefBlock() uint64 {
	return t.ReferenceBlockNumber
}

func (t *CreateAgentSlashTaskRequest) GetTaskInfo(taskIndex sdktypes.TaskIndex, quorumNums sdktypes.QuorumNums, quorumThresholds sdktypes.QuorumThresholdPercentages) ITaskInfo {
	return &AgentSlashTaskInfo{
		QuorumNumbers:              quorumNums,
		QuorumThresholdPercentages: quorumThresholds,
		TaskIndex:                  taskIndex,
		ReferenceBlockNumber:       t.ReferenceBlockNumber,
		TxnId:                      t.TxnId,
		InsufficientBalance:        big.NewInt(0),
		AgentAddress:               t.AgentAddress,
	}
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

func NewSignedAgentSlashTaskRespRequest(req *aggregator.SignedAgentSlashTaskRespRequest) (*CreateSignedTaskResRequest, error) {
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

	res := &CreateSignedTaskResRequest{
		Task:         balHead,
		BlsSignature: sign,
	}

	copy(res.OperatorId[:], operatorId[:32])

	return res, nil
}

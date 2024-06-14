package core

import (
	"encoding/hex"
	"fmt"
	"log/slog"

	aggregator "github.com/Cedro-Finance/CedroAVS/api/grpc"
	sdktypes "github.com/Layr-Labs/eigensdk-go/types"
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

type TaskIndex uint32

type BalanceTaskInfo struct {
	RequestId  uint64
	QuorumNumbers              sdktypes.QuorumNums
	QuorumThresholdPercentages sdktypes.QuorumThresholdPercentages
	TaskIndex                  TaskIndex
	ReferenceBlockNumber       uint64
}


func NewBalanceTaskInfo(req *aggregator.BalanceTaskInfo) (*BalanceTaskInfo, error) {
	alertHash := req.GetBalanceHash()
	if len(alertHash) != 32 {
		return nil, fmt.Errorf("balanceHash len should be 32")
	}

	res := &BalanceTaskInfo{
		QuorumNumbers:              ConvertQuorumNumbersFromBytes(req.GetQuorumNumbers()),
		QuorumThresholdPercentages: ConvertQuorumThresholdPercentagesFromBytes(req.GetQuorumThresholdPercentages()),
		TaskIndex:                  TaskIndex(req.GetTaskIndex()),
		ReferenceBlockNumber:       req.GetReferenceBlockNumber(),
		RequestId: 100,
	}

	// copy(res.RequestId[:], alertHash[:32])

	return res, nil
}
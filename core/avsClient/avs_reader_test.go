package avs

import (
	"context"

	sdktypes "github.com/Layr-Labs/eigensdk-go/types"
)

type MockAVSReader struct {
	getQuorum          func() (sdktypes.QuorumNums, error)
	getQuorumThreshold func() (sdktypes.QuorumThresholdPercentages, error)
}

func (m *MockAVSReader) GetQuorumsByBlockNumber(ctx context.Context, blockNumber uint32) (sdktypes.QuorumNums, error) {
	return m.getQuorum()
}

func (m *MockAVSReader) GetQuorumThresholdPercentages(ctx context.Context, blockNumber uint32, quorums sdktypes.QuorumNums) (sdktypes.QuorumThresholdPercentages, error) {
	return m.getQuorumThreshold()
}

// func TestSomething(t *testing.T) {
// 	mckAvsReader := &MockAVSReader{
// 		getQuorum: func() (sdktypes.QuorumNums, error) {
// 			return []sdktypes.QuorumNum{1,2,3}, nil 	
// 		},
// 		getQuorumThreshold: func() (sdktypes.QuorumThresholdPercentages, error) {
// 			return nil, nil 
// 		},
// 	}
// 	taskProcessor := 
// }

package common

import (
	"context"
	"math/big"

	"github.com/Cedro-Finance/CedroAVS/api/graphql"
	"github.com/Cedro-Finance/CedroAVS/verifier/client/agents"
	"github.com/ethereum/go-ethereum/common"
)

type IAgentPodClient interface {
	InitializeAgentPod(addr common.Address)
	GetBalance(context.Context, *agents.DeviationRequest, chan<- *big.Int)
	GetGraphClient() *graphql.GraphQLClient
	PollEvents(context.Context, chan<- agents.TaskRequest, chan<- agents.TaskRequest)

	GetCurrentBlock(context.Context) (uint64, error)
	GetChainId() string

	GetFinalizedBlock(context.Context) (uint64, error)

	GetBlockFromTimestamp(context.Context, uint64) (uint64, error)

	
}

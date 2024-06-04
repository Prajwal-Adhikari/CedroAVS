package common

import (
	"context"

	"github.com/Cedro-Finance/CedroAVS/verifierOperator/client/agents"
	"github.com/Cedro-Finance/CedroAVS/verifierOperator/client/avs"
)

type IClient interface {
	GetBalance(context.Context, *avs.Request, chan<- *agents.BalanceResponse)
}


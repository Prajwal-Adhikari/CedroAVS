package tasks

import (
	"context"
	"sync"

	"github.com/Cedro-Finance/CedroAVS/verifier/client/agents"
)

type Taskprocessor interface {
	Start(ctx context.Context, wg *sync.WaitGroup)
	GetRequestCh() chan agents.TaskRequest
}

type TaskProcessor struct {
}

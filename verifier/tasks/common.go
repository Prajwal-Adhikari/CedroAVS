package tasks

import (
	"context"
	"sync"

	"github.com/Cedro-Finance/CedroAVS/verifier/client/avs"
)

type Taskprocessor interface {
	Start(ctx context.Context, wg *sync.WaitGroup)
	GetRequestCh() chan avs.TaskRequest
	// GetType() string
}

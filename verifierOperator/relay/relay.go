package relay

import (
	"context"
	"fmt"
	"math/big"

	"github.com/Cedro-Finance/CedroAVS/verifierOperator/client/avs"
	"github.com/Cedro-Finance/CedroAVS/verifierOperator/common"
)

var RegisteredClients map[string]common.IClient

type Relay struct {
	avs *avs.AVSClient
}

func (r *Relay) EigenEventPoller(eventCh chan<- *avs.Request) {
	request := r.avs.PollEvents()
	eventCh <- request
}

func (r *Relay) GetChainBalances(eventCh <-chan *avs.Request) {
	requestEvent := <-eventCh
	balanceCh := make(chan *big.Int, len(RegisteredClients))
	for _, v := range RegisteredClients {
		go v.Balance(requestEvent, balanceCh) // TODO: filter balance related event
	}
	ctx := context.TODO() // remove later
	for {
		select {
		case <-ctx.Done():
			return
		case balance := <-balanceCh:
			fmt.Println("balance is", balance)

		}
	}
}

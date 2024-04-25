package relay

import (
	"context"
	"fmt"
	"math/big"

	"github.com/Cedro-Finance/CedroAVS/verifierOperator/client/agents"
	"github.com/Cedro-Finance/CedroAVS/verifierOperator/client/avs"
	"github.com/Cedro-Finance/CedroAVS/verifierOperator/common"
	"github.com/Cedro-Finance/CedroAVS/verifierOperator/config"
	ethCommon "github.com/ethereum/go-ethereum/common"
)

var RegisteredClients map[string] common.IClient

type Relay struct {
	avs *avs.AVSClient

}


func (r *Relay) GetChainBalances(eventCh <-chan *big.Int) {
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

func (r *Relay) UpdateChainBalances(reqCh <-chan *avs.Request, balanceCh chan *big.Int) {
	reqEvent := <-reqCh

	for _, v := range RegisteredClients {
		go v.UpdateBalance(reqEvent.reqId, ethCommon.Address{})
	}
}

func StartRelay(ctx context.Context, cfg *config.Config) {
	avsClient := avs.NewAVSClient(&cfg.AVS);
	balanceReqCh := make(chan *avs.Request) //Chanel to accept request from AVS contract
	balanceCh := make(chan *agents.BalanceResponse)  //Channel to accept balance from agent contract



	r := Relay {
		avs: avsClient,

	}
	go r.avs.PollEvents(ctx, balanceReqCh)
	go r.UpdateChainBalances(balanceReqCh, balanceCh)
	// go r.GetChainBalances(balanceReqCh, balanceCh)
}
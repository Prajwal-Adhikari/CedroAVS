package relay

import (
	"context"
	"fmt"

	"github.com/Cedro-Finance/CedroAVS/verifierOperator/client/agents"
	"github.com/Cedro-Finance/CedroAVS/verifierOperator/client/avs"
	"github.com/Cedro-Finance/CedroAVS/verifierOperator/common"
	"github.com/Cedro-Finance/CedroAVS/verifierOperator/config"
	"github.com/Cedro-Finance/CedroAVS/verifierOperator/keys"
	ethCommon "github.com/ethereum/go-ethereum/common"
)

var RegisteredClients map[string]func(config.ChainConfig) common.IClient

type Relay struct {
	avs *avs.AVSClient
	agents []common.IClient
	keystore *keys.KeyStore

}


func (r *Relay) GetChainBalances(ctx context.Context, eventCh <-chan *avs.Request, balanceCh chan *agents.BalanceResponse) {

	for _, ag := range r.agents {
		go ag.Balance(ctx, balanceCh)
	}
	for {
		select {
		case <-ctx.Done():
			return
		case balance := <-balanceCh:
			//TODO: handle whether this response is same as the request or we don't need to?
			fmt.Println("balance is", balance)
		}
	}
}

func (r *Relay) UpdateChainBalances(ctx context.Context, reqCh <-chan *avs.Request) {
	for {
		reqEvent := <-reqCh
		
		for _, ag := range r.agents {
			go ag.UpdateBalance(ctx, avs.GetRequestId(reqEvent), ethCommon.Address{})
		}
	}
}

func StartRelay(ctx context.Context, cfg *config.Config) {
	avsClient := avs.NewAVSClient(&cfg.AVS);
	balanceReqCh := make(chan *avs.Request) //Chanel to accept request from AVS contract
	balanceCh := make(chan *agents.BalanceResponse, len(RegisteredClients))


	regAgents := make([]common.IClient, len(cfg.Chains))
	for _, cfg := range cfg.Chains {
		regAgents = append(regAgents, RegisteredClients[cfg.Name](cfg))
	}

	keystore := keys.LoadKeys(cfg.KeyStore)
	
	r := Relay {
		avs: avsClient,
		agents: regAgents,
		keystore: keystore,
	}



	go r.avs.PollEvents(ctx, balanceReqCh)
	go r.UpdateChainBalances(ctx, balanceReqCh)
	go r.GetChainBalances(ctx, balanceReqCh, balanceCh)
}
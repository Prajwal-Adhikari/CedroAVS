package eventlistener

import (
	"context"
	"log"
	"math/big"
	"strings"

	"github.com/Cedro-Finance/CedroAVS/agentEnroller/config"
	contractutils "github.com/Cedro-Finance/CedroAVS/agentEnroller/contract/utils"
	"github.com/Cedro-Finance/CedroAVS/agentEnroller/eventhandler"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

func SyncToLatest(startBlock int64) {
	RootChainData := config.LoadConfig().RootChainData

	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(startBlock),
		Addresses: []common.Address{
			RootChainData.AgentManager,
		},
	}

	client, _ := contractutils.GetClient(RootChainData.NodeUrl)

	logs, err := client.FilterLogs(context.Background(), query)

	if err != nil {
		log.Fatal(err)
	}

	contractAbi, err := abi.JSON(strings.NewReader(string(config.GetAgentManagerAbi())))
	if err != nil {
		log.Fatal(err)
	}

	for _, vLog := range logs {
		eventhandler.HandleEventReceived(contractAbi, vLog)
	}
}

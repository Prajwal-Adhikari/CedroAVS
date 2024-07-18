package eventhandler

import (
	"fmt"

	contractdeployer "github.com/Cedro-Finance/CedroAVS/agentEnroller/contract/deployer"
	"github.com/ethereum/go-ethereum/common"
)

type EnrollAgentEventData struct {
	OldOwner common.Address
	NewOwner common.Address
}

func EnrollAgentHandler(eventData EnrollAgentEventData) {
	fmt.Printf("Event data: %+v\n", eventData)
	// TODO: Read params from eventdata and add go routine for deploying to different chains
	// TODO: RPC urls to be mapped from configs
	contractdeployer.DeployContract(eventData.NewOwner, "https://ethereum-holesky-rpc.publicnode.com")
}

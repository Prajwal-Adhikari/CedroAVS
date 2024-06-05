package avs

import (
	// "github.com/Cedro-Finance/CedroAVS/verifierOperator/config"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/avsregistry"
	// "github.com/Layr-Labs/eigensdk-go/chainio/clients/eth"
	"github.com/Layr-Labs/eigensdk-go/chainio/txmgr"
	// "github.com/Layr-Labs/eigensdk-go/logging"
	// "github.com/ethereum/go-ethereum/common"
)

type IAvsWriter interface {
	avsregistry.AvsRegistryWriter
}

type AvsWriter struct {
	avsregistry.AvsRegistryWriter
	AvsContractBinding *AVSManagerBinding

	TxMgr txmgr.TxManager
}


// var _ IAvsWriter = (*AvsWriter)(nil)

// func NewAvsWritedr(cfg *config.AVSConfig, ethClient eth.Client) (*AvsWriter, error) {
// 	regrAddr := common.HexToAddress(cfg.RegistryCoordinatorAddr)
// 	osAddr := common.HexToAddress(cfg.OperatorStateRetrieverAddr)

// 	avsManagerBindings, err := NewAvsManagerBindings(regrAddr, osAddr, ethClient)

// 	if err != nil {
// 		return nil, err
// 	}

// 	avsRegistryWriter, err := avsregistry.BuildAvsRegistryChainWriter(regrAddr, osAddr, logging.NewNoopLogger(), ethClient, )
// }
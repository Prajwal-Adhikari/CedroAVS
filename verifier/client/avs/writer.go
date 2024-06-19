package avs

import (
	// "github.com/Cedro-Finance/CedroAVS/verifier/config"
	"github.com/Cedro-Finance/CedroAVS/verifier/config"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/avsregistry"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/eth"
	"github.com/Layr-Labs/eigensdk-go/logging"

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
	AvsContractBindings *AVSManagerBinding

	TxMgr  txmgr.TxManager
	logger logging.Logger
}

var _ IAvsWriter = (*AvsWriter)(nil)

func NewAvsWriter(c *config.AVSConfig, txMgr txmgr.TxManager, ethClient eth.Client, logger logging.Logger) (*AvsWriter, error) {
	regrAddr := c.RegistryCoordinatorAddr
	osAddr := c.OperatorStateRetrieverAddr

	avsManagerBindings, err := NewAvsManagerBindings(regrAddr, osAddr, ethClient)

	if err != nil {
		logger.Error("Failed to create contract bindings", "err", err)
		return nil, err
	}

	avsRegistryWriter, err := avsregistry.BuildAvsRegistryChainWriter(regrAddr, osAddr, logger, ethClient, txMgr)
	if err != nil {
		return nil, err
	}

	return &AvsWriter{
		AvsRegistryWriter:   avsRegistryWriter,
		AvsContractBindings: avsManagerBindings,
		logger:              logger,
		TxMgr:               txMgr,
	}, nil

}

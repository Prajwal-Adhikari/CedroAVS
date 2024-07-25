package avs

import (
	// "github.com/Cedro-Finance/CedroAVS/verifier/config"
	"context"
	"reflect"

	cedroservicemanager "github.com/Cedro-Finance/CedroAVS/core/bindings/CedroServiceManager"
	"github.com/Cedro-Finance/CedroAVS/core/config"
	"github.com/Cedro-Finance/CedroAVS/core/tasks"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/avsregistry"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/eth"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum/core/types"

	// "github.com/Layr-Labs/eigensdk-go/chainio/clients/eth"
	"github.com/Layr-Labs/eigensdk-go/chainio/txmgr"
	// "github.com/Layr-Labs/eigensdk-go/logging"
	// "github.com/ethereum/go-ethereum/common"
)

type IAvsWriter interface {
	avsregistry.AvsRegistryWriter

	SendConfirmationBatch(ctx context.Context, balanceHeader tasks.ITaskInfo, nonSignerStakeSignature cedroservicemanager.IBLSSignatureCheckerNonSignerStakesAndSignature) (*types.Receipt, error)
}

type AvsWriter struct {
	avsregistry.AvsRegistryWriter
	AvsContractBindings *AVSManagerBinding

	TxMgr  txmgr.TxManager
	logger logging.Logger
}

var _ IAvsWriter = (*AvsWriter)(nil)

func NewAvsWriter(c *config.Config, txMgr txmgr.TxManager, ethClient eth.Client, logger logging.Logger) (IAvsWriter, error) {
	regrAddr := c.AVS.RegistryCoordinatorAddr
	osAddr := c.AVS.OperatorStateRetrieverAddr

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

func (w *AvsWriter) SendConfirmationBatch(ctx context.Context, balanceHeader tasks.ITaskInfo, nonSignerStakeSignature cedroservicemanager.IBLSSignatureCheckerNonSignerStakesAndSignature) (*types.Receipt, error) {
	txOpts, err := w.TxMgr.GetNoSendTxOpts()
	if err != nil {
		w.logger.Errorf("Error getting tx opts")
		return nil, err
	}

	w.logger.Info("Submitting Batch COnfirmation", "bal", balanceHeader)

	var tx *types.Transaction

	switch reflect.TypeOf(balanceHeader) {
	case reflect.TypeOf(&tasks.BalanceTaskInfo{}):

		balHeader := balanceHeader.(*tasks.BalanceTaskInfo)

		parsedHeader := balHeader.ToCedroServiceManagerHeader()

		tx, err = w.AvsContractBindings.ServiceManager.ConfirmAgentDeviation(txOpts, parsedHeader, nonSignerStakeSignature)

	case reflect.TypeOf(&tasks.AgentSlashTaskInfo{}):

		balHeader := balanceHeader.(*tasks.AgentSlashTaskInfo)

		parsedHeader := balHeader.ToCedroServiceManagerHeader()

		tx, err = w.AvsContractBindings.ServiceManager.ConfirmAgentFraud(txOpts, parsedHeader, nonSignerStakeSignature)

	}

	if err != nil {
		w.logger.Error("Error submitting SubmitTaskResponse tx while calling respondToTask", "err", err)
		return nil, err
	}
	receipt, err := w.TxMgr.Send(ctx, tx)
	if err != nil {
		w.logger.Errorf("Error submitting CreateNewTask tx")
		return nil, err
	}
	return receipt, nil
}

func (w *AvsWriter) SendConfirmationAgentSlash(ctx context.Context, slashHeader *tasks.AgentSlashTaskInfo, nonSignerStakeSignature cedroservicemanager.IBLSSignatureCheckerNonSignerStakesAndSignature) (*types.Receipt, error) {
	txOpts, err := w.TxMgr.GetNoSendTxOpts()

	if err != nil {
		w.logger.Error("Error getting tx opts")
		return nil, err
	}

	parsedHeader := slashHeader.ToCedroServiceManagerHeader()

	w.logger.Info("Submitting Batch Confirmation for Agent Slash", "agent", slashHeader.AgentAddress)

	tx, err := w.AvsContractBindings.ServiceManager.ConfirmAgentFraud(txOpts, parsedHeader, nonSignerStakeSignature)
	if err != nil {
		w.logger.Error("Error submitting SubmitTaskResponse tx while calling respondToTask", "err", err)
		return nil, err
	}
	receipt, err := w.TxMgr.Send(ctx, tx)
	if err != nil {
		w.logger.Errorf("Error submitting agentslash task tx")
		return nil, err
	}
	return receipt, nil
}

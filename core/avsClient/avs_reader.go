package avs

import (
	"context"
	"math/big"

	"github.com/Cedro-Finance/CedroAVS/core/config"
	sdkavsregistry "github.com/Layr-Labs/eigensdk-go/chainio/clients/avsregistry"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/eth"
	"github.com/Layr-Labs/eigensdk-go/logging"
	sdktypes "github.com/Layr-Labs/eigensdk-go/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

type IAVSReader interface {
	sdkavsregistry.AvsRegistryReader

	//GetQuorumsByBlockNumber
	GetQuorumsByBlockNumber(ctx context.Context, blockNumber uint32) (sdktypes.QuorumNums, error)

	//GetQuorumThresholdPercentages
	GetQuorumThresholdPercentages(ctx context.Context, blockNumber uint32, quorums sdktypes.QuorumNums) (sdktypes.QuorumThresholdPercentages, error)
}

type AvsReader struct {
	sdkavsregistry.AvsRegistryReader
	AvsServiceBindings *AVSManagerBinding
}

func NewAvsReader(c config.AVSConfig, ethClient eth.Client, logger logging.Logger) (*AvsReader, error) {
	regrAddr := c.RegistryCoordinatorAddr
	osAddr := c.OperatorStateRetrieverAddr

	avsManagerBindings, err := NewAvsManagerBindings(regrAddr, osAddr, ethClient)

	if err != nil {
		return nil, err
	}

	avsRegistryReader, err := sdkavsregistry.BuildAvsRegistryChainReader(regrAddr, osAddr, ethClient, logger)

	if err != nil {
		return nil, err
	}

	return &AvsReader{
		AvsServiceBindings: avsManagerBindings,
		AvsRegistryReader:  avsRegistryReader,
	}, nil
}

func (r *AvsReader) GetQuorumsByBlockNumber(ctx context.Context, blockNumber uint32) (sdktypes.QuorumNums, error) {
	quorumCount, err := r.AvsRegistryReader.GetQuorumCount(&bind.CallOpts{
		Context:     ctx,
		BlockNumber: big.NewInt(int64(blockNumber)),
	})
	if err != nil {
		return nil, err
	}

	res := make([]sdktypes.QuorumNum, 0, quorumCount)
	for i := uint8(0); i < quorumCount; i++ {
		res = append(res, sdktypes.QuorumNum(i))
	}

	return res, nil
}

func (r *AvsReader) GetQuorumThresholdPercentages(ctx context.Context, blockNumber uint32, quorums sdktypes.QuorumNums) (sdktypes.QuorumThresholdPercentages, error) {
	quorumThresholdPercentage, err := r.AvsServiceBindings.ServiceManager.QuorumThresholdPercentage(&bind.CallOpts{
		Context:     ctx,
		BlockNumber: big.NewInt(int64(blockNumber)),
	})

	if err != nil {
		return nil, err
	}

	res := make([]sdktypes.QuorumThresholdPercentage, 0, len(quorums))
	for i := 0; i < len(quorums); i++ {
		res = append(res, sdktypes.QuorumThresholdPercentage(uint8(quorumThresholdPercentage)))
	}

	return res, nil
}

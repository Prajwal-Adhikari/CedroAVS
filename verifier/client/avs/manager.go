package avs

import (
	"context"
	"fmt"
	"math/big"
	"time"

	abi "github.com/Cedro-Finance/CedroAVS/verifier/client/avs/abi"
	"github.com/Cedro-Finance/CedroAVS/verifier/config"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/eth"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
)

type AgentManager struct {
	Eth     eth.Client
	Address common.Address

	PollDuration    time.Duration
	StartBlock      uint64
	NextBlockHeight uint64
	BlockFinality   uint64
	BlockTime       time.Duration

	Logger logging.Logger

	Contract *abi.AgentManager
}

func NewAgentManager(cfg *config.Config, logger logging.Logger, ethClient eth.Client) (*AgentManager, error) {

	contract, err := abi.NewAgentManager(cfg.AgentManagerAddr, ethClient)

	fmt.Println(cfg.AgentManagerAddr)

	if err != nil {
		return nil, err
	}

	return &AgentManager{
		Eth:             ethClient,
		Address:         cfg.AgentManagerAddr,
		PollDuration:    time.Duration(cfg.AgentManagerPollDuration) * time.Second,
		StartBlock:      cfg.AgentManagerStartHeight,
		NextBlockHeight: cfg.AgentManagerStartHeight,
		BlockFinality:   cfg.EthConfig.BlockFinality,
		BlockTime:       time.Duration(cfg.EthConfig.BlockTime) * time.Second,
		Logger:          logger,
		Contract:        contract,
	}, nil

}

func (ag *AgentManager) PollEvents(ctx context.Context, ch chan<- *Request) {
	dur := ag.PollDuration

	if dur == 0 {
		dur = defaultPollDuration
	}

	ticker := time.NewTicker(dur)
	defer ticker.Stop()

	for {

		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
		}

	L1:

		for {
			currentBlock, err := ag.Eth.BlockNumber(ctx)

			fmt.Println(currentBlock)

			if err != nil {
				ag.Logger.Error("Error while getting current block", "err", err)

				break L1
			}

			finalizedBlock := currentBlock - ag.BlockFinality

			if finalizedBlock < ag.NextBlockHeight {
				diff := ag.NextBlockHeight - finalizedBlock

				time.Sleep((time.Duration(diff) * ag.BlockTime))
			}

			for startHeight := ag.NextBlockHeight; startHeight <= finalizedBlock; startHeight += defaultHeightDifferenceForFilterLogs + 1 {
				endHeight := startHeight + defaultHeightDifferenceForFilterLogs

				if endHeight > finalizedBlock {
					endHeight = finalizedBlock
				}

				logs, err := ag.Eth.FilterLogs(ctx, ethereum.FilterQuery{
					FromBlock: big.NewInt(int64(startHeight)),
					ToBlock:   big.NewInt(int64(endHeight)),
					Topics:    [][]common.Hash{{common.HexToHash("0xa08d0dfd9bd009969096b75b05bd74056e062358eac1eb43af41815db2bd1b77")}},
					Addresses: []common.Address{ag.Address},
				})

				if err != nil {
					ag.Logger.Error("Error while getting logs", err)
					break L1
				}

				for _, log := range logs {
					parsed, err := ag.Contract.ParseDeviationCheck(log)

					if err != nil {
						break L1
					}

					blockInfo, err := ag.Eth.BlockByHash(ctx, log.BlockHash)

					if err != nil {
						break L1
					}

					ag.Logger.Infof("Deviation Request Received",
						"RequestId", parsed.RequestID.String(),
						"Agent Address", parsed.Agent,
						"BlockNumber", log.BlockNumber,
					)

					ch <- &Request{
						ReqId:                parsed.RequestID,
						ReferenceBlockNumber: log.BlockNumber,
						AgentAddress:         parsed.Agent,
						Timestamp:            blockInfo.Time(),
					}
				}

				ag.NextBlockHeight = endHeight + 1
			}
		}

	}

}

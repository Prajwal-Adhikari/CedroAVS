package avs

import (
	"context"
	"fmt"
	"math/big"
	"time"

	agentmanager "github.com/Cedro-Finance/CedroAVS/core/bindings/AgentManager"
	cedroservicemanager "github.com/Cedro-Finance/CedroAVS/core/bindings/CedroServiceManager"
	"github.com/Cedro-Finance/CedroAVS/verifier/config"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/eth"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

const (
	defaultBlockTime = time.Duration(12) * time.Second

	defaultStartBlock                    = 1705750
	defaultHeightDifferenceForFilterLogs = 100
	defaultRequestId                     = 1

	defaultPollDuration = 5 * time.Second
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

	Contract *agentmanager.AgentManager
}

func NewAgentManager(cfg *config.Config, logger logging.Logger, ethClient eth.Client) (*AgentManager, error) {

	contract, err := agentmanager.NewAgentManager(cfg.AgentManagerAddr, ethClient)

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

func (ag *AgentManager) PollEvents(ctx context.Context, ch chan<- *Request, serviceManager *cedroservicemanager.CedroServiceManager) {
	dur := ag.PollDuration

	fmt.Println(dur)

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

			fmt.Println(currentBlock, ag.NextBlockHeight)

			if err != nil {
				ag.Logger.Error("Error while getting current block", "err", err)

				break L1
			}

			if currentBlock%100 == 0 {
				ag.Logger.Info("Polling Block", "blockk", currentBlock)
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
					Topics:    [][]common.Hash{{common.HexToHash("0xf1f6ae32e492ed94fefcd957830c206af14cae19c30c54be0814d27e8cdb5ac2")}},
					Addresses: []common.Address{ag.Address},
				})

				if err != nil {
					ag.Logger.Error("Error while getting logs", err)
					break L1
				}

				for _, log := range logs {
					parsed, err := ag.Contract.ParseDeviationCheck(log)

					if err != nil {
						ag.Logger.Error("Error while parsing deviation event", "err", err)
						break L1
					}

					blockInfo, err := ag.Eth.BlockByHash(ctx, log.BlockHash)

					if err != nil {
						ag.Logger.Error("Error while getting block inforamtion", "err", err)
						break L1
					}

					ag.Logger.Info("Deviation Request Received",
						"RequestId", parsed.RequestID.String(),
						"Agent Address", parsed.Agent,
						"BlockNumber", log.BlockNumber,
					)

					lastReqId, err := serviceManager.AgentsDeviationRequestId(&bind.CallOpts{}, parsed.Agent)

					if err != nil {
						ag.Logger.Error("Error while fetching last requestId", "err", err)
						break L1
					}

					if lastReqId.Cmp(parsed.RequestID) > -1 {
						ag.Logger.Warn("Request Id already processed", "reqId", parsed.RequestID.String(), "agent", parsed.Agent.Hex())
						break
					}

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

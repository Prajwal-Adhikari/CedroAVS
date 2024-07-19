package operator

import (
	// "context"

	// "math/big"
	// "time"
	// "github.com/Cedro-Finance/CedroAVS/core"

	// "github.com/Cedro-Finance/CedroAVS/verifier/client/avs"
	// "github.com/ethereum/go-ethereum/accounts/abi/bind"



	// "github.com/Cedro-Finance/CedroAVS/verifier/client/avs"
	// "github.com/Cedro-Finance/CedroAVS/verifier/client/avs"

)




// func (r *Operator) StartOffChainDeviationCheck(ctx context.Context) {

// 	r.logger.Info("Offchain trigger service started")

// 	// blockDetails, err := r.avs.Eth.BlockByNumber(ctx, big.NewInt(int64(1867349)))

// 	// if err != nil {
// 	// 	r.logger.Error("OffchainDeviationCheckTrigger errro geting block details", "err", err)

// 	// }

// 	// r.balanceReqCh <- &avs.Request{
// 	// 	ReqId:                big.NewInt(6),
// 	// 	ReferenceBlockNumber: 1867349		,
// 	// 	AgentAddress:         ethcommon.HexToAddress("0xdF1Cf915dB7A658dFDEb75aF3873B4e55F25Aa42"),
// 	// 	Timestamp:            blockDetails.Time(),
// 	// 	// IsOffChain:           true,
// 	// }

// 	dur := time.Duration(time.Minute * 1)

// 	ticker := time.NewTicker(dur)

// 	for {
// 		select {
// 		case <-ctx.Done():
// 			return

// 		case <-ticker.C:
// 			r.logger.Infof("Offchain deviation trigger period")
// 		}

// 	RETRY_LOOP:
// 		for {
// 			agents, err := r.graphClient.GetEnrolledAgentsList(ctx)

// 			if err != nil {
// 				r.logger.Error("Error while getting all the enrolled agents in agent manager")
// 				break RETRY_LOOP
// 			}

// 			//Current Block Number

// 			blockNumber, err := r.avs.Eth.BlockNumber(ctx)

// 			if err != nil {
// 				r.logger.Error("OffchainDeviationCheckTrigger error while getting block number", "err", err)

// 				break RETRY_LOOP
// 			}

// 			blockDetails, err := r.avs.Eth.BlockByNumber(ctx, big.NewInt(int64(blockNumber)))

// 			if err != nil {
// 				r.logger.Error("OffchainDeviationCheckTrigger errro geting block details", "err", err)

// 				break RETRY_LOOP
// 			}

// 			for _, agent := range agents {
// 				maxUint256 := new(big.Int)
// 				maxUint256.SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF", 16)
// 				r.balanceReqCh <- &avs.DeviationRequest{
// 					ReqId:                maxUint256,
// 					ReferenceBlockNumber: blockNumber,
// 					AgentAddress:         agent,
// 					Timestamp:            blockDetails.Time(),
// 					IsOffChain:           true,
// 				}
// 			}

// 			break
// 		}

// 	}

// }



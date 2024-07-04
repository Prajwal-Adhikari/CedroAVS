package operator

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/Cedro-Finance/CedroAVS/core"
	// "github.com/Cedro-Finance/CedroAVS/verifier/client/avs"
	// "github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func (r *Operator) SignAndSendBalanceResponse(ctx context.Context, balHeader *core.BalanceTaskInfo) {

	r.logger.Info("Balance Heaader", "header", balHeader)

	hash, err := balHeader.GetBalanceHeaderHash()

	fmt.Println("==========HASH===============", hex.EncodeToString(hash[:]))

	if err != nil {
		r.logger.Error("Error while hasing the message", "err", err)

		r.errCh <- &DeviationError{
			TaskIndex: balHeader.TaskIndex,
			Error:     err,
		}
		return
	}

	signature := r.blsKeyPair.SignMessage(hash)

	r.logger.Info("Signed the message", "sign", signature)

	request := &core.SignedTaskRespRequest{
		Balance:      *balHeader,
		BlsSignature: *signature,
		OperatorId:   r.operatorId,
	}

	r.signatureCh <- request

}

func (r *Operator) SendSignedTaskToAggregator(ctx context.Context, req *core.SignedTaskRespRequest) {

	maxRetry := 10
	retries := 0

RETRY_LOOP:

	for {
		response, err := r.gRpcClient.SendSignedTaskResponse(req)

		if err != nil {
			r.logger.Error("Error while sending Signed Task response", "err", err)

			r.errCh <- &DeviationError{
				TaskIndex: req.Balance.TaskIndex,
				Error:     err,
			}

			if retries < maxRetry {
				r.logger.Info("Retrying sending signature to aggregator")
				retries++
				break RETRY_LOOP
			}

			return

		} else {
			r.balanceTaskMutex.Lock()
			delete(r.balanceTasks, req.Balance.TaskIndex)
			r.balanceTaskMutex.Unlock()

			r.logger.Info("Successfully posted the signatures to the aggregator",
				"task", req.Balance.TaskIndex,
				"reqId", req.Balance.ReqId.String(),
				"agent", req.Balance.AgentAddress.Hex(),
				"response", response)

			return
		}

	}

}

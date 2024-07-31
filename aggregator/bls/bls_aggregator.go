package blsaggregator

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"sync"
	"time"

	avs "github.com/Cedro-Finance/CedroAVS/core/avsClient"
	"github.com/Layr-Labs/eigensdk-go/crypto/bls"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/Layr-Labs/eigensdk-go/types"
	"github.com/Layr-Labs/eigensdk-go/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

var (
	// TODO: refactor these errors to use a custom struct with taskIndex field instead of wrapping taskIndex in the error string directly.
	//       see https://go.dev/blog/go1.13-errors
	TaskInitializationErrorFn = func(err error, taskIndex types.TaskIndex) error {
		return fmt.Errorf("Failed to initialize task %d: %w", taskIndex, err)
	}
	TaskAlreadyInitializedErrorFn = func(taskIndex types.TaskIndex) error {
		return fmt.Errorf("task %d already initialized", taskIndex)
	}
	TaskExpiredErrorFn = func(taskIndex types.TaskIndex) error {
		return fmt.Errorf("task %d expired", taskIndex)
	}
	TaskNotFoundErrorFn = func(taskIndex types.TaskIndex) error {
		return fmt.Errorf("task %d not initialized or already completed", taskIndex)
	}
	OperatorNotPartOfTaskQuorumErrorFn = func(operatorId types.OperatorId, taskIndex types.TaskIndex) error {
		return fmt.Errorf("operator %x not part of task %d's quorum", operatorId, taskIndex)
	}
	HashFunctionError = func(err error) error {
		return fmt.Errorf("failed to hash task response: %w", err)
	}
	SignatureVerificationError = func(err error) error {
		return fmt.Errorf("failed to verify signature: %w", err)
	}
	IncorrectSignatureError = errors.New("signature verification failed Incorrect Signature")
)

type BLSAggregationServiceResponse struct {
	Err                error                    // if Err is not nil, the other fields are not valid
	TaskIndex          types.TaskIndex          // unique identifier of the task
	TaskResponse       interface{}              // the task response that was signed
	TaskResponseDigest types.TaskResponseDigest // digest of the task response that was signed
	// The below 8 fields are the data needed to build the IBLSSignatureChecker.NonSignerStakesAndSignature struct
	// users of this service will need to build the struct themselves by converting the bls points
	// into the BN254.G1/G2Point structs that the IBLSSignatureChecker expects
	// given that those are different for each AVS service manager that individually inherits BLSSignatureChecker
	NonSignersPubkeysG1          []*bls.G1Point
	QuorumApksG1                 []*bls.G1Point
	SignersApkG2                 map[types.QuorumNum]*bls.G2Point
	SignersAggSigG1              map[types.QuorumNum]*bls.Signature
	NonSignerQuorumBitmapIndices []uint32
	QuorumApkIndices             []uint32
	TotalStakeIndices            []uint32
	NonSignerStakeIndices        [][]uint32
}

type BLSAggregatorService struct {
	logger    logging.Logger
	avsReader avs.IAVSReader

	aggregatedResponsesC chan BLSAggregationServiceResponse

	signedTaskRespsCs map[types.TaskIndex]chan types.SignedTaskResponseDigest

	taskChansMutex sync.RWMutex
}

func NewBlsAggregatorService(avsReader avs.IAVSReader, logger logging.Logger) *BLSAggregatorService {
	return &BLSAggregatorService{
		avsReader:            avsReader,
		logger:               logger,
		aggregatedResponsesC: make(chan BLSAggregationServiceResponse),
		signedTaskRespsCs:    make(map[types.TaskIndex]chan types.SignedTaskResponseDigest),
		taskChansMutex:       sync.RWMutex{},
	}
}

func (b *BLSAggregatorService) GetResponseChannel() <-chan BLSAggregationServiceResponse {
	return b.aggregatedResponsesC
}

func (b *BLSAggregatorService) InitializeNewTask(taskIndex types.TaskIndex,
	taskCreatedBlock uint32,
	quorumNumbers types.QuorumNums,
	quorumThresholdPercentages types.QuorumThresholdPercentages,
	timeToExpiry time.Duration) error {
	b.logger.Debug("AggregatorService initializing new task", "taskIndex", taskIndex, "taskCreatedBlock", taskCreatedBlock, "quorumNumbers", quorumNumbers, "quorumThresholdPercentages", quorumThresholdPercentages, "timeToExpiry", timeToExpiry)

	if _, taskExists := b.signedTaskRespsCs[taskIndex]; taskExists {
		return fmt.Errorf("task %d already initialized", taskIndex)
	}

	signedTaskRespsC := make(chan types.SignedTaskResponseDigest)
	b.taskChansMutex.Lock()
	b.signedTaskRespsCs[taskIndex] = signedTaskRespsC
	b.taskChansMutex.Unlock()

	go b.singleTaskAggregatorGoroutineFunc(taskIndex, taskCreatedBlock, quorumNumbers, quorumThresholdPercentages, timeToExpiry, signedTaskRespsC)

	return nil

}

func (b *BLSAggregatorService) ProcessNewSignature(ctx context.Context,
	taskIndex types.TaskIndex,
	taskResponse types.Bytes32,
	blsSignature *bls.Signature,
	operatorId types.OperatorId) error {

	b.logger.Info("Signature Received for Processing", "opId", operatorId, "taskIndex", taskIndex, "taskResponse", taskResponse)

	b.taskChansMutex.Lock()
	taskC, taskInitialized := b.signedTaskRespsCs[taskIndex]
	b.taskChansMutex.Unlock()
	if !taskInitialized {
		b.logger.Info("Task not initilized")
		return TaskNotFoundErrorFn(taskIndex)
	}

	signatureVerificationErrorC := make(chan error)
	// send the task to the goroutine processing this task
	// and return the error (if any) returned by the signature verification routine

	select {
	// we need to send this as part of select because if the goroutine is processing another SignedTaskResponseDigest
	// and cannot receive this one, we want the context to be able to cancel the request
	case <-ctx.Done():
		return ctx.Err()
	case taskC <- types.SignedTaskResponseDigest{
		TaskResponseDigest:          taskResponse,
		BlsSignature:                blsSignature,
		OperatorId:                  operatorId,
		SignatureVerificationErrorC: signatureVerificationErrorC,
	}:
		// note that we need to wait synchronously here for this response because we want to
		// send back an informative error message to the operator who sent his signature to the aggregator
		err := <-signatureVerificationErrorC

		return err
	}
}

func (b *BLSAggregatorService) singleTaskAggregatorGoroutineFunc(
	taskIndex types.TaskIndex,
	taskCreatedBlock uint32,
	quorumNumbers types.QuorumNums,
	quorumThresholdPercentages []types.QuorumThresholdPercentage,
	timeToExpiry time.Duration,
	signedTaskRespsC <-chan types.SignedTaskResponseDigest,
) {
	defer b.closeTaskGoroutine(taskIndex)

	quorumThresholdPercentagesMap := make(map[types.QuorumNum]types.QuorumThresholdPercentage)
	for i, quorumNumber := range quorumNumbers {
		quorumThresholdPercentagesMap[quorumNumber] = quorumThresholdPercentages[i]
	}

	operatorsAvsStateDict, err := b.avsReader.GetOperatorsAvsStateAtBlock(context.Background(), quorumNumbers, taskCreatedBlock)
	if err != nil {
		b.aggregatedResponsesC <- BLSAggregationServiceResponse{
			Err:       TaskInitializationErrorFn(fmt.Errorf("aggregatorService failed to get operators state from avs registry at blockNum %d: %w", taskCreatedBlock, err), taskIndex),
			TaskIndex: taskIndex,
		}
		return
	}

	quorumsAvsStakeDict, err := b.avsReader.GetQuorumsAvsStateAtBlock(context.Background(), quorumNumbers, taskCreatedBlock)
	if err != nil {
		b.aggregatedResponsesC <- BLSAggregationServiceResponse{
			Err:       TaskInitializationErrorFn(fmt.Errorf("aggregator failed to get quorums state from avs registry: %w", err), taskIndex),
			TaskIndex: taskIndex,
		}
		return
	}

	totalStakePerQuorum := make(map[types.QuorumNum]*big.Int)
	for quorumNum, quorumAvsState := range quorumsAvsStakeDict {
		totalStakePerQuorum[quorumNum] = quorumAvsState.TotalStake
	}

	quorumApksG1 := []*bls.G1Point{}
	for _, quorumNumber := range quorumNumbers {
		quorumApksG1 = append(quorumApksG1, quorumsAvsStakeDict[quorumNumber].AggPubkeyG1)
	}

	// TODO instead of taking a TTE, we should take a block as input
	// and monitor the chain and only close the task goroutine when that block is reached
	taskExpiredTimer := time.NewTimer(timeToExpiry)

	// aggregatedOperatorsDict := map[types.TaskResponseDigest]aggregatedOperators{}

	aggSigs := make(map[types.QuorumNum]*bls.Signature, len(quorumNumbers))
	aggPubKeys := make(map[types.QuorumNum]*bls.G2Point, len(quorumNumbers))
	signerMap := make(map[types.OperatorId]bool)

	stakeSigned := make(map[types.QuorumNum]*big.Int, len(quorumNumbers))
	for _, quorumID := range quorumNumbers {
		stakeSigned[quorumID] = big.NewInt(0)
		aggSigs[quorumID] = bls.NewZeroSignature()
		aggPubKeys[quorumID] = bls.NewZeroG2Point()
	}

	for {
		select {
		case signedTaskResponseDigest := <-signedTaskRespsC:
			b.logger.Debug("Task goroutine received new signed task response digest", "taskIndex", taskIndex, "signedTaskResponseDigest", signedTaskResponseDigest)

			err := b.verifySignature(taskIndex, signedTaskResponseDigest, operatorsAvsStateDict)
			signedTaskResponseDigest.SignatureVerificationErrorC <- err
			if err != nil {
				continue
			}

			// compute the taskResponseDigest using the hash function
			taskResponseDigest := signedTaskResponseDigest.TaskResponseDigest

			for _, quorumNum := range quorumNumbers {

				stake, ok := operatorsAvsStateDict[signedTaskResponseDigest.OperatorId].StakePerQuorum[quorumNum]

				if !ok || stake.Cmp(big.NewInt(0)) == 0 {
					continue
				}

				signerMap[signedTaskResponseDigest.OperatorId] = true

				stakeSigned[quorumNum].Add(stakeSigned[quorumNum], stake)

				aggSigs[quorumNum].Add(signedTaskResponseDigest.BlsSignature)
				aggPubKeys[quorumNum].Add(operatorsAvsStateDict[signedTaskResponseDigest.OperatorId].OperatorInfo.Pubkeys.G2Pubkey)
			}

			if checkIfStakeThresholdsMet(b.logger, stakeSigned, totalStakePerQuorum, quorumThresholdPercentagesMap) {
				nonSignersOperatorIds := []types.OperatorId{}
				for operatorId := range operatorsAvsStateDict {
					if _, operatorSigned := signerMap[operatorId]; !operatorSigned {
						nonSignersOperatorIds = append(nonSignersOperatorIds, operatorId)
					}
				}

				// the contract requires a sorted nonSignersOperatorIds
				sort.SliceStable(nonSignersOperatorIds, func(i, j int) bool {
					iOprInt := new(big.Int).SetBytes(nonSignersOperatorIds[i][:])
					jOprInt := new(big.Int).SetBytes(nonSignersOperatorIds[j][:])
					return iOprInt.Cmp(jOprInt) == -1
				})

				nonSignersG1Pubkeys := []*bls.G1Point{}
				for _, operatorId := range nonSignersOperatorIds {
					operator := operatorsAvsStateDict[operatorId]
					nonSignersG1Pubkeys = append(nonSignersG1Pubkeys, operator.OperatorInfo.Pubkeys.G1Pubkey)
				}

				quorumAggPubKeys := make(map[types.QuorumNum]*bls.G1Point, len(quorumNumbers))

				for _, quorumNum := range quorumNumbers {

					quorumAggKey := quorumsAvsStakeDict[quorumNum].AggPubkeyG1

					quorumAggPubKeys[quorumNum] = quorumAggKey

					signerAggkey := bls.NewZeroG1Point().Add(quorumAggKey)

					for opInd, nsk := range nonSignersG1Pubkeys {

						stake, ok := operatorsAvsStateDict[nonSignersOperatorIds[opInd]].StakePerQuorum[quorumNum]

						if !ok || stake.Cmp(big.NewInt(0)) == 0 {
							continue
						}

						signerAggkey = signerAggkey.Sub(nsk)
					}

					ok, err := signerAggkey.VerifyEquivalence(aggPubKeys[quorumNum])

					if !ok || err != nil {
						b.aggregatedResponsesC <- BLSAggregationServiceResponse{
							Err:       utils.WrapError(errors.New("failed to get verify quorum signature equivalent"), err),
							TaskIndex: taskIndex,
						}
						continue
					}

					ok, err = aggSigs[quorumNum].Verify(aggPubKeys[quorumNum], taskResponseDigest)
					if !ok || err != nil {
						b.aggregatedResponsesC <- BLSAggregationServiceResponse{
							Err:       utils.WrapError(errors.New("failed to get verify quorum signature"), err),
							TaskIndex: taskIndex,
						}
						return
					}

				}

				indices, err := b.avsReader.GetCheckSignaturesIndices(&bind.CallOpts{}, taskCreatedBlock, quorumNumbers, nonSignersOperatorIds)
				if err != nil {
					b.aggregatedResponsesC <- BLSAggregationServiceResponse{
						Err:       utils.WrapError(errors.New("failed to get check signatures indices"), err),
						TaskIndex: taskIndex,
					}
					return
				}

				blsAggregationServiceResponse := BLSAggregationServiceResponse{
					Err:                          nil,
					TaskIndex:                    taskIndex,
					TaskResponse:                 signedTaskResponseDigest.TaskResponseDigest,
					TaskResponseDigest:           taskResponseDigest,
					NonSignersPubkeysG1:          nonSignersG1Pubkeys,
					QuorumApksG1:                 quorumApksG1,
					SignersApkG2:                 aggPubKeys,
					SignersAggSigG1:              aggSigs,
					NonSignerQuorumBitmapIndices: indices.NonSignerQuorumBitmapIndices,
					QuorumApkIndices:             indices.QuorumApkIndices,
					TotalStakeIndices:            indices.TotalStakeIndices,
					NonSignerStakeIndices:        indices.NonSignerStakeIndices,
				}

				fmt.Println(blsAggregationServiceResponse)
				// b.aggregatedResponsesC <- blsAggregationServiceResponse
				b.closeTaskGoroutine(taskIndex)
				taskExpiredTimer.Stop()
				return
			}
		case <-taskExpiredTimer.C:
			b.aggregatedResponsesC <- BLSAggregationServiceResponse{
				Err:       TaskExpiredErrorFn(taskIndex),
				TaskIndex: taskIndex,
			}
			return
		}
	}

}

func (b *BLSAggregatorService) closeTaskGoroutine(taskIndex types.TaskIndex) {
	b.taskChansMutex.Lock()
	delete(b.signedTaskRespsCs, taskIndex)
	b.taskChansMutex.Unlock()

	b.logger.Info("Closed single goroutine", "taskIndex", taskIndex)
}

func (a *BLSAggregatorService) verifySignature(
	taskIndex types.TaskIndex,
	signedTaskResponseDigest types.SignedTaskResponseDigest,
	operatorsAvsStateDict map[types.OperatorId]types.OperatorAvsState,
) error {
	_, ok := operatorsAvsStateDict[signedTaskResponseDigest.OperatorId]
	if !ok {
		a.logger.Warnf("Operator %#v not found. Skipping message.", signedTaskResponseDigest.OperatorId)
		return OperatorNotPartOfTaskQuorumErrorFn(signedTaskResponseDigest.OperatorId, taskIndex)
	}

	taskResponseDigest := signedTaskResponseDigest.TaskResponseDigest

	// verify that the msg actually came from the correct operator
	operatorG2Pubkey := operatorsAvsStateDict[signedTaskResponseDigest.OperatorId].OperatorInfo.Pubkeys.G2Pubkey
	if operatorG2Pubkey == nil {
		a.logger.Error("Operator G2 pubkey not found", "operatorId", signedTaskResponseDigest.OperatorId, "taskId", taskIndex)
		return fmt.Errorf("taskId %d: Operator G2 pubkey not found (operatorId: %x)", taskIndex, signedTaskResponseDigest.OperatorId)
	}
	a.logger.Debug("Verifying signed task response digest signature",
		"operatorG2Pubkey", operatorG2Pubkey,
		"taskResponseDigest", taskResponseDigest,
		"blsSignature", signedTaskResponseDigest.BlsSignature,
	)

	// if the operator signs a digest that is not the digest of the TaskResponse submitted in ProcessNewTask
	// then the signature will not be verified
	signatureVerified, err := signedTaskResponseDigest.BlsSignature.Verify(operatorG2Pubkey, taskResponseDigest)
	if err != nil {
		return SignatureVerificationError(err)
	}
	if !signatureVerified {
		return IncorrectSignatureError
	}
	return nil
}

func checkIfStakeThresholdsMet(
	logger logging.Logger,
	signedStakePerQuorum map[types.QuorumNum]*big.Int,
	totalStakePerQuorum map[types.QuorumNum]*big.Int,
	quorumThresholdPercentagesMap map[types.QuorumNum]types.QuorumThresholdPercentage,
) bool {
	for quorumNum, quorumThresholdPercentage := range quorumThresholdPercentagesMap {
		signedStakeByQuorum, ok := signedStakePerQuorum[quorumNum]
		if !ok {
			// signedStakePerQuorum not contain the quorum,
			// this case means signedStakePerQuorum has not signed for each quorum.
			// even the total stake for this quorum is zero.
			return false
		}

		totalStakeByQuorum, ok := totalStakePerQuorum[quorumNum]
		if !ok {
			// Note this case should not happen.
			// The `totalStakePerQuorum` is got from the contract, so if we not found the
			// totalStakeByQuorum, that means the code have a bug.
			logger.Errorf("TotalStake not found for quorum %d.", quorumNum)
			return false
		}

		// we check that signedStake >= totalStake * quorumThresholdPercentage / 100
		// to be exact (and do like the contracts), we actually check that
		// signedStake * 100 >= totalStake * quorumThresholdPercentage
		signedStake := big.NewInt(0).Mul(signedStakeByQuorum, big.NewInt(100))
		thresholdStake := big.NewInt(0).Mul(totalStakeByQuorum, big.NewInt(int64(quorumThresholdPercentage)))
		if signedStake.Cmp(thresholdStake) < 0 {
			return false
		}
	}
	return true
}

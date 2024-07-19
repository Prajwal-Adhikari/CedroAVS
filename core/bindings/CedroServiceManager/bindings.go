// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package cedroservicemanager

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// BN254G1Point is an auto generated low-level Go binding around an user-defined struct.
type BN254G1Point struct {
	X *big.Int
	Y *big.Int
}

// BN254G2Point is an auto generated low-level Go binding around an user-defined struct.
type BN254G2Point struct {
	X [2]*big.Int
	Y [2]*big.Int
}

// CedroServiceDataDeviationBalanceHeader is an auto generated low-level Go binding around an user-defined struct.
type CedroServiceDataDeviationBalanceHeader struct {
	AgentAddress               common.Address
	ReferenceBlockNumber       uint32
	QuorumNumbers              []byte
	QuorumThresholdPercentages []byte
	RequestId                  *big.Int
	ActualAmount               *big.Int
	IdealAmount                *big.Int
}

// CedroServiceDataFraudBalanceHeader is an auto generated low-level Go binding around an user-defined struct.
type CedroServiceDataFraudBalanceHeader struct {
	AgentAddress               common.Address
	ReferenceBlockNumber       uint32
	QuorumNumbers              []byte
	QuorumThresholdPercentages []byte
	TransactionId              *big.Int
	InsufficientAmount         *big.Int
}

// CedroServiceDataOperatorChallengeEvidence is an auto generated low-level Go binding around an user-defined struct.
type CedroServiceDataOperatorChallengeEvidence struct {
	MessageHash [32]byte
}

// IBLSSignatureCheckerNonSignerStakesAndSignature is an auto generated low-level Go binding around an user-defined struct.
type IBLSSignatureCheckerNonSignerStakesAndSignature struct {
	NonSignerQuorumBitmapIndices []uint32
	NonSignerPubkeys             []BN254G1Point
	QuorumApks                   []BN254G1Point
	ApkG2                        BN254G2Point
	Sigma                        BN254G1Point
	QuorumApkIndices             []uint32
	TotalStakeIndices            []uint32
	NonSignerStakeIndices        [][]uint32
}

// IBLSSignatureCheckerQuorumStakeTotals is an auto generated low-level Go binding around an user-defined struct.
type IBLSSignatureCheckerQuorumStakeTotals struct {
	SignedStakeForQuorum []*big.Int
	TotalStakeForQuorum  []*big.Int
}

// IRewardsCoordinatorRewardsSubmission is an auto generated low-level Go binding around an user-defined struct.
type IRewardsCoordinatorRewardsSubmission struct {
	StrategiesAndMultipliers []IRewardsCoordinatorStrategyAndMultiplier
	Token                    common.Address
	Amount                   *big.Int
	StartTimestamp           uint32
	Duration                 uint32
}

// IRewardsCoordinatorStrategyAndMultiplier is an auto generated low-level Go binding around an user-defined struct.
type IRewardsCoordinatorStrategyAndMultiplier struct {
	Strategy   common.Address
	Multiplier *big.Int
}

// ISignatureUtilsSignatureWithSaltAndExpiry is an auto generated low-level Go binding around an user-defined struct.
type ISignatureUtilsSignatureWithSaltAndExpiry struct {
	Signature []byte
	Salt      [32]byte
	Expiry    *big.Int
}

// CedroServiceManagerMetaData contains all meta data concerning the CedroServiceManager contract.
var CedroServiceManagerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"ContractTAG\",\"type\":\"string\"}],\"name\":\"InsufficientThreshold\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"ContractTAG\",\"type\":\"string\"}],\"name\":\"InsufficientThresholdPercentages\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"ContractTAG\",\"type\":\"string\"}],\"name\":\"InvalidAmount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"ContractTAG\",\"type\":\"string\"}],\"name\":\"InvalidQuorumParam\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"ContractTAG\",\"type\":\"string\"}],\"name\":\"InvalidQuorumThresholdPercentage\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"ContractTAG\",\"type\":\"string\"}],\"name\":\"InvalidReferenceBlockNum\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"ContractTAG\",\"type\":\"string\"}],\"name\":\"InvalidRequestId\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"ContractTAG\",\"type\":\"string\"}],\"name\":\"InvalidSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"contractTAG\",\"type\":\"string\"}],\"name\":\"MessageHashAlreadyAdded\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"contractTAG\",\"type\":\"string\"}],\"name\":\"NotInAllowlist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"contractTAG\",\"type\":\"string\"}],\"name\":\"OperatorRegistrationPaused\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"agentAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"SlashingAmount\",\"type\":\"uint256\"}],\"name\":\"AgentBalanceConfirmed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"agentAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"SlashingAmount\",\"type\":\"uint256\"}],\"name\":\"AgentFraudConfirmed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"OperatorDeRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"OperatorRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"thresholdpercentage\",\"type\":\"uint256\"}],\"name\":\"QuorumThresholdPercentageChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"prevRewardsInitiator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newRewardsInitiator\",\"type\":\"address\"}],\"name\":\"RewardsInitiatorUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"value\",\"type\":\"bool\"}],\"name\":\"StaleStakesForbiddenUpdate\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"THRESHOLD_DENOMINATOR\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_avsDirectory\",\"outputs\":[{\"internalType\":\"contractIAVSDirectory\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_registryCoordinator\",\"outputs\":[{\"internalType\":\"contractIRegistryCoordinator\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_rewardsCoordinator\",\"outputs\":[{\"internalType\":\"contractIRewardsCoordinator\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"operators\",\"type\":\"address[]\"}],\"name\":\"addOperatorsToAllowlist\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"agentManager\",\"outputs\":[{\"internalType\":\"contractIAgentManager\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"agentsDeviationRequestId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"agentsFraudRequestId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"avsDirectory\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"blsApkRegistry\",\"outputs\":[{\"internalType\":\"contractIBLSApkRegistry\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"messageHash\",\"type\":\"bytes32\"}],\"internalType\":\"structCedroServiceData.OperatorChallengeEvidence\",\"name\":\"evidence\",\"type\":\"tuple\"}],\"name\":\"challengeOperators\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"msgHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"quorumNumbers\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"referenceBlockNumber\",\"type\":\"uint32\"},{\"components\":[{\"internalType\":\"uint32[]\",\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structBN254.G1Point[]\",\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structBN254.G1Point[]\",\"name\":\"quorumApks\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256[2]\",\"name\":\"X\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"Y\",\"type\":\"uint256[2]\"}],\"internalType\":\"structBN254.G2Point\",\"name\":\"apkG2\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structBN254.G1Point\",\"name\":\"sigma\",\"type\":\"tuple\"},{\"internalType\":\"uint32[]\",\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\"},{\"internalType\":\"uint32[]\",\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\"},{\"internalType\":\"uint32[][]\",\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\"}],\"internalType\":\"structIBLSSignatureChecker.NonSignerStakesAndSignature\",\"name\":\"params\",\"type\":\"tuple\"}],\"name\":\"checkSignatures\",\"outputs\":[{\"components\":[{\"internalType\":\"uint96[]\",\"name\":\"signedStakeForQuorum\",\"type\":\"uint96[]\"},{\"internalType\":\"uint96[]\",\"name\":\"totalStakeForQuorum\",\"type\":\"uint96[]\"}],\"internalType\":\"structIBLSSignatureChecker.QuorumStakeTotals\",\"name\":\"\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"agentAddress\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"referenceBlockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"quorumNumbers\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"quorumThresholdPercentages\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actualAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"idealAmount\",\"type\":\"uint256\"}],\"internalType\":\"structCedroServiceData.DeviationBalanceHeader\",\"name\":\"deviationBalanceHeader\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint32[]\",\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structBN254.G1Point[]\",\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structBN254.G1Point[]\",\"name\":\"quorumApks\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256[2]\",\"name\":\"X\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"Y\",\"type\":\"uint256[2]\"}],\"internalType\":\"structBN254.G2Point\",\"name\":\"apkG2\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structBN254.G1Point\",\"name\":\"sigma\",\"type\":\"tuple\"},{\"internalType\":\"uint32[]\",\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\"},{\"internalType\":\"uint32[]\",\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\"},{\"internalType\":\"uint32[][]\",\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\"}],\"internalType\":\"structIBLSSignatureChecker.NonSignerStakesAndSignature\",\"name\":\"nonSignerStakesAndSignature\",\"type\":\"tuple\"}],\"name\":\"confirmAgentDeviation\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"agentAddress\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"referenceBlockNumber\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"quorumNumbers\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"quorumThresholdPercentages\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"transactionId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"insufficientAmount\",\"type\":\"uint256\"}],\"internalType\":\"structCedroServiceData.FraudBalanceHeader\",\"name\":\"fraudBalanceHeader\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint32[]\",\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structBN254.G1Point[]\",\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structBN254.G1Point[]\",\"name\":\"quorumApks\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256[2]\",\"name\":\"X\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"Y\",\"type\":\"uint256[2]\"}],\"internalType\":\"structBN254.G2Point\",\"name\":\"apkG2\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structBN254.G1Point\",\"name\":\"sigma\",\"type\":\"tuple\"},{\"internalType\":\"uint32[]\",\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\"},{\"internalType\":\"uint32[]\",\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\"},{\"internalType\":\"uint32[][]\",\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\"}],\"internalType\":\"structIBLSSignatureChecker.NonSignerStakesAndSignature\",\"name\":\"nonSignerStakesAndSignature\",\"type\":\"tuple\"}],\"name\":\"confirmAgentFraud\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"messageHash\",\"type\":\"bytes32\"}],\"name\":\"contains\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"contractIStrategy\",\"name\":\"strategy\",\"type\":\"address\"},{\"internalType\":\"uint96\",\"name\":\"multiplier\",\"type\":\"uint96\"}],\"internalType\":\"structIRewardsCoordinator.StrategyAndMultiplier[]\",\"name\":\"strategiesAndMultipliers\",\"type\":\"tuple[]\"},{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"startTimestamp\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"duration\",\"type\":\"uint32\"}],\"internalType\":\"structIRewardsCoordinator.RewardsSubmission[]\",\"name\":\"rewardsSubmissions\",\"type\":\"tuple[]\"}],\"name\":\"createAVSRewardsSubmission\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"delegation\",\"outputs\":[{\"internalType\":\"contractIDelegationManager\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"deregisterOperatorFromAVS\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"getOperatorRestakedStrategies\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRestakeableStrategies\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"__roleManagerAddress\",\"type\":\"address\"},{\"internalType\":\"contractIAVSDirectory\",\"name\":\"__avsDirectory\",\"type\":\"address\"},{\"internalType\":\"contractIRewardsCoordinator\",\"name\":\"__rewardsCoordinator\",\"type\":\"address\"},{\"internalType\":\"contractIStakeRegistry\",\"name\":\"__stakeRegistry\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"referenceBlockNumber\",\"type\":\"uint32\"}],\"name\":\"latestServeUntilBlock\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"operatorAllowlist\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"quorumThresholdPercentage\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"salt\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"expiry\",\"type\":\"uint256\"}],\"internalType\":\"structISignatureUtils.SignatureWithSaltAndExpiry\",\"name\":\"operatorSignature\",\"type\":\"tuple\"}],\"name\":\"registerOperatorToAVS\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"registryCoordinator\",\"outputs\":[{\"internalType\":\"contractIRegistryCoordinator\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rewardsInitiator\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"roleManager\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"roleManagerAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIAgentManager\",\"name\":\"_agentManagerAddress\",\"type\":\"address\"}],\"name\":\"setAgentManagerAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"status\",\"type\":\"bool\"}],\"name\":\"setPauseStatus\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"thresholdPercentage\",\"type\":\"uint8\"}],\"name\":\"setQuorumThresholdPercentage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIRegistryCoordinator\",\"name\":\"__registryCoordinator\",\"type\":\"address\"}],\"name\":\"setRegistryCoordinatorAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newRewardsInitiator\",\"type\":\"address\"}],\"name\":\"setRewardsInitiator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_roleManagerAddress\",\"type\":\"address\"}],\"name\":\"setRoleManager\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"value\",\"type\":\"bool\"}],\"name\":\"setStaleStakesForbidden\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"stakeRegistry\",\"outputs\":[{\"internalType\":\"contractIStakeRegistry\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"staleStakesForbidden\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"msgHash\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structBN254.G1Point\",\"name\":\"apk\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256[2]\",\"name\":\"X\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"Y\",\"type\":\"uint256[2]\"}],\"internalType\":\"structBN254.G2Point\",\"name\":\"apkG2\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structBN254.G1Point\",\"name\":\"sigma\",\"type\":\"tuple\"}],\"name\":\"trySignatureAndApkVerification\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"pairingSuccessful\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"siganatureIsValid\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_metadataURI\",\"type\":\"string\"}],\"name\":\"updateAVSMetadataURI\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// CedroServiceManagerABI is the input ABI used to generate the binding from.
// Deprecated: Use CedroServiceManagerMetaData.ABI instead.
var CedroServiceManagerABI = CedroServiceManagerMetaData.ABI

// CedroServiceManager is an auto generated Go binding around an Ethereum contract.
type CedroServiceManager struct {
	CedroServiceManagerCaller     // Read-only binding to the contract
	CedroServiceManagerTransactor // Write-only binding to the contract
	CedroServiceManagerFilterer   // Log filterer for contract events
}

// CedroServiceManagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type CedroServiceManagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CedroServiceManagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CedroServiceManagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CedroServiceManagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CedroServiceManagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CedroServiceManagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CedroServiceManagerSession struct {
	Contract     *CedroServiceManager // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// CedroServiceManagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CedroServiceManagerCallerSession struct {
	Contract *CedroServiceManagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// CedroServiceManagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CedroServiceManagerTransactorSession struct {
	Contract     *CedroServiceManagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// CedroServiceManagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type CedroServiceManagerRaw struct {
	Contract *CedroServiceManager // Generic contract binding to access the raw methods on
}

// CedroServiceManagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CedroServiceManagerCallerRaw struct {
	Contract *CedroServiceManagerCaller // Generic read-only contract binding to access the raw methods on
}

// CedroServiceManagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CedroServiceManagerTransactorRaw struct {
	Contract *CedroServiceManagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCedroServiceManager creates a new instance of CedroServiceManager, bound to a specific deployed contract.
func NewCedroServiceManager(address common.Address, backend bind.ContractBackend) (*CedroServiceManager, error) {
	contract, err := bindCedroServiceManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CedroServiceManager{CedroServiceManagerCaller: CedroServiceManagerCaller{contract: contract}, CedroServiceManagerTransactor: CedroServiceManagerTransactor{contract: contract}, CedroServiceManagerFilterer: CedroServiceManagerFilterer{contract: contract}}, nil
}

// NewCedroServiceManagerCaller creates a new read-only instance of CedroServiceManager, bound to a specific deployed contract.
func NewCedroServiceManagerCaller(address common.Address, caller bind.ContractCaller) (*CedroServiceManagerCaller, error) {
	contract, err := bindCedroServiceManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CedroServiceManagerCaller{contract: contract}, nil
}

// NewCedroServiceManagerTransactor creates a new write-only instance of CedroServiceManager, bound to a specific deployed contract.
func NewCedroServiceManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*CedroServiceManagerTransactor, error) {
	contract, err := bindCedroServiceManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CedroServiceManagerTransactor{contract: contract}, nil
}

// NewCedroServiceManagerFilterer creates a new log filterer instance of CedroServiceManager, bound to a specific deployed contract.
func NewCedroServiceManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*CedroServiceManagerFilterer, error) {
	contract, err := bindCedroServiceManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CedroServiceManagerFilterer{contract: contract}, nil
}

// bindCedroServiceManager binds a generic wrapper to an already deployed contract.
func bindCedroServiceManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CedroServiceManagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CedroServiceManager *CedroServiceManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CedroServiceManager.Contract.CedroServiceManagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CedroServiceManager *CedroServiceManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CedroServiceManager.Contract.CedroServiceManagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CedroServiceManager *CedroServiceManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CedroServiceManager.Contract.CedroServiceManagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CedroServiceManager *CedroServiceManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CedroServiceManager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CedroServiceManager *CedroServiceManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CedroServiceManager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CedroServiceManager *CedroServiceManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CedroServiceManager.Contract.contract.Transact(opts, method, params...)
}

// THRESHOLDDENOMINATOR is a free data retrieval call binding the contract method 0xef024458.
//
// Solidity: function THRESHOLD_DENOMINATOR() view returns(uint256)
func (_CedroServiceManager *CedroServiceManagerCaller) THRESHOLDDENOMINATOR(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CedroServiceManager.contract.Call(opts, &out, "THRESHOLD_DENOMINATOR")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// THRESHOLDDENOMINATOR is a free data retrieval call binding the contract method 0xef024458.
//
// Solidity: function THRESHOLD_DENOMINATOR() view returns(uint256)
func (_CedroServiceManager *CedroServiceManagerSession) THRESHOLDDENOMINATOR() (*big.Int, error) {
	return _CedroServiceManager.Contract.THRESHOLDDENOMINATOR(&_CedroServiceManager.CallOpts)
}

// THRESHOLDDENOMINATOR is a free data retrieval call binding the contract method 0xef024458.
//
// Solidity: function THRESHOLD_DENOMINATOR() view returns(uint256)
func (_CedroServiceManager *CedroServiceManagerCallerSession) THRESHOLDDENOMINATOR() (*big.Int, error) {
	return _CedroServiceManager.Contract.THRESHOLDDENOMINATOR(&_CedroServiceManager.CallOpts)
}

// AvsDirectory is a free data retrieval call binding the contract method 0x91a83d61.
//
// Solidity: function _avsDirectory() view returns(address)
func (_CedroServiceManager *CedroServiceManagerCaller) AvsDirectory(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CedroServiceManager.contract.Call(opts, &out, "_avsDirectory")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AvsDirectory is a free data retrieval call binding the contract method 0x91a83d61.
//
// Solidity: function _avsDirectory() view returns(address)
func (_CedroServiceManager *CedroServiceManagerSession) AvsDirectory() (common.Address, error) {
	return _CedroServiceManager.Contract.AvsDirectory(&_CedroServiceManager.CallOpts)
}

// AvsDirectory is a free data retrieval call binding the contract method 0x91a83d61.
//
// Solidity: function _avsDirectory() view returns(address)
func (_CedroServiceManager *CedroServiceManagerCallerSession) AvsDirectory() (common.Address, error) {
	return _CedroServiceManager.Contract.AvsDirectory(&_CedroServiceManager.CallOpts)
}

// CedroRegistryCoordinator is a free data retrieval call binding the contract method 0x59b7b232.
//
// Solidity: function _registryCoordinator() view returns(address)
func (_CedroServiceManager *CedroServiceManagerCaller) CedroRegistryCoordinator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CedroServiceManager.contract.Call(opts, &out, "_registryCoordinator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// CedroRegistryCoordinator is a free data retrieval call binding the contract method 0x59b7b232.
//
// Solidity: function _registryCoordinator() view returns(address)
func (_CedroServiceManager *CedroServiceManagerSession) CedroRegistryCoordinator() (common.Address, error) {
	return _CedroServiceManager.Contract.CedroRegistryCoordinator(&_CedroServiceManager.CallOpts)
}

// CedroRegistryCoordinator is a free data retrieval call binding the contract method 0x59b7b232.
//
// Solidity: function _registryCoordinator() view returns(address)
func (_CedroServiceManager *CedroServiceManagerCallerSession) CedroRegistryCoordinator() (common.Address, error) {
	return _CedroServiceManager.Contract.CedroRegistryCoordinator(&_CedroServiceManager.CallOpts)
}

// RewardsCoordinator is a free data retrieval call binding the contract method 0x026c7b09.
//
// Solidity: function _rewardsCoordinator() view returns(address)
func (_CedroServiceManager *CedroServiceManagerCaller) RewardsCoordinator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CedroServiceManager.contract.Call(opts, &out, "_rewardsCoordinator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RewardsCoordinator is a free data retrieval call binding the contract method 0x026c7b09.
//
// Solidity: function _rewardsCoordinator() view returns(address)
func (_CedroServiceManager *CedroServiceManagerSession) RewardsCoordinator() (common.Address, error) {
	return _CedroServiceManager.Contract.RewardsCoordinator(&_CedroServiceManager.CallOpts)
}

// RewardsCoordinator is a free data retrieval call binding the contract method 0x026c7b09.
//
// Solidity: function _rewardsCoordinator() view returns(address)
func (_CedroServiceManager *CedroServiceManagerCallerSession) RewardsCoordinator() (common.Address, error) {
	return _CedroServiceManager.Contract.RewardsCoordinator(&_CedroServiceManager.CallOpts)
}

// AgentManager is a free data retrieval call binding the contract method 0x7622f78d.
//
// Solidity: function agentManager() view returns(address)
func (_CedroServiceManager *CedroServiceManagerCaller) AgentManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CedroServiceManager.contract.Call(opts, &out, "agentManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AgentManager is a free data retrieval call binding the contract method 0x7622f78d.
//
// Solidity: function agentManager() view returns(address)
func (_CedroServiceManager *CedroServiceManagerSession) AgentManager() (common.Address, error) {
	return _CedroServiceManager.Contract.AgentManager(&_CedroServiceManager.CallOpts)
}

// AgentManager is a free data retrieval call binding the contract method 0x7622f78d.
//
// Solidity: function agentManager() view returns(address)
func (_CedroServiceManager *CedroServiceManagerCallerSession) AgentManager() (common.Address, error) {
	return _CedroServiceManager.Contract.AgentManager(&_CedroServiceManager.CallOpts)
}

// AgentsDeviationRequestId is a free data retrieval call binding the contract method 0xd51b1779.
//
// Solidity: function agentsDeviationRequestId(address ) view returns(uint256)
func (_CedroServiceManager *CedroServiceManagerCaller) AgentsDeviationRequestId(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _CedroServiceManager.contract.Call(opts, &out, "agentsDeviationRequestId", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AgentsDeviationRequestId is a free data retrieval call binding the contract method 0xd51b1779.
//
// Solidity: function agentsDeviationRequestId(address ) view returns(uint256)
func (_CedroServiceManager *CedroServiceManagerSession) AgentsDeviationRequestId(arg0 common.Address) (*big.Int, error) {
	return _CedroServiceManager.Contract.AgentsDeviationRequestId(&_CedroServiceManager.CallOpts, arg0)
}

// AgentsDeviationRequestId is a free data retrieval call binding the contract method 0xd51b1779.
//
// Solidity: function agentsDeviationRequestId(address ) view returns(uint256)
func (_CedroServiceManager *CedroServiceManagerCallerSession) AgentsDeviationRequestId(arg0 common.Address) (*big.Int, error) {
	return _CedroServiceManager.Contract.AgentsDeviationRequestId(&_CedroServiceManager.CallOpts, arg0)
}

// AgentsFraudRequestId is a free data retrieval call binding the contract method 0x113ea8ea.
//
// Solidity: function agentsFraudRequestId(address ) view returns(uint256)
func (_CedroServiceManager *CedroServiceManagerCaller) AgentsFraudRequestId(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _CedroServiceManager.contract.Call(opts, &out, "agentsFraudRequestId", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AgentsFraudRequestId is a free data retrieval call binding the contract method 0x113ea8ea.
//
// Solidity: function agentsFraudRequestId(address ) view returns(uint256)
func (_CedroServiceManager *CedroServiceManagerSession) AgentsFraudRequestId(arg0 common.Address) (*big.Int, error) {
	return _CedroServiceManager.Contract.AgentsFraudRequestId(&_CedroServiceManager.CallOpts, arg0)
}

// AgentsFraudRequestId is a free data retrieval call binding the contract method 0x113ea8ea.
//
// Solidity: function agentsFraudRequestId(address ) view returns(uint256)
func (_CedroServiceManager *CedroServiceManagerCallerSession) AgentsFraudRequestId(arg0 common.Address) (*big.Int, error) {
	return _CedroServiceManager.Contract.AgentsFraudRequestId(&_CedroServiceManager.CallOpts, arg0)
}

// CedroAvsDirectory is a free data retrieval call binding the contract method 0x6b3aa72e.
//
// Solidity: function avsDirectory() view returns(address)
func (_CedroServiceManager *CedroServiceManagerCaller) CedroAvsDirectory(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CedroServiceManager.contract.Call(opts, &out, "avsDirectory")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// CedroAvsDirectory is a free data retrieval call binding the contract method 0x6b3aa72e.
//
// Solidity: function avsDirectory() view returns(address)
func (_CedroServiceManager *CedroServiceManagerSession) CedroAvsDirectory() (common.Address, error) {
	return _CedroServiceManager.Contract.CedroAvsDirectory(&_CedroServiceManager.CallOpts)
}

// CedroAvsDirectory is a free data retrieval call binding the contract method 0x6b3aa72e.
//
// Solidity: function avsDirectory() view returns(address)
func (_CedroServiceManager *CedroServiceManagerCallerSession) CedroAvsDirectory() (common.Address, error) {
	return _CedroServiceManager.Contract.CedroAvsDirectory(&_CedroServiceManager.CallOpts)
}

// BlsApkRegistry is a free data retrieval call binding the contract method 0x5df45946.
//
// Solidity: function blsApkRegistry() view returns(address)
func (_CedroServiceManager *CedroServiceManagerCaller) BlsApkRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CedroServiceManager.contract.Call(opts, &out, "blsApkRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BlsApkRegistry is a free data retrieval call binding the contract method 0x5df45946.
//
// Solidity: function blsApkRegistry() view returns(address)
func (_CedroServiceManager *CedroServiceManagerSession) BlsApkRegistry() (common.Address, error) {
	return _CedroServiceManager.Contract.BlsApkRegistry(&_CedroServiceManager.CallOpts)
}

// BlsApkRegistry is a free data retrieval call binding the contract method 0x5df45946.
//
// Solidity: function blsApkRegistry() view returns(address)
func (_CedroServiceManager *CedroServiceManagerCallerSession) BlsApkRegistry() (common.Address, error) {
	return _CedroServiceManager.Contract.BlsApkRegistry(&_CedroServiceManager.CallOpts)
}

// CheckSignatures is a free data retrieval call binding the contract method 0x6efb4636.
//
// Solidity: function checkSignatures(bytes32 msgHash, bytes quorumNumbers, uint32 referenceBlockNumber, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) params) view returns((uint96[],uint96[]), bytes32)
func (_CedroServiceManager *CedroServiceManagerCaller) CheckSignatures(opts *bind.CallOpts, msgHash [32]byte, quorumNumbers []byte, referenceBlockNumber uint32, params IBLSSignatureCheckerNonSignerStakesAndSignature) (IBLSSignatureCheckerQuorumStakeTotals, [32]byte, error) {
	var out []interface{}
	err := _CedroServiceManager.contract.Call(opts, &out, "checkSignatures", msgHash, quorumNumbers, referenceBlockNumber, params)

	if err != nil {
		return *new(IBLSSignatureCheckerQuorumStakeTotals), *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new(IBLSSignatureCheckerQuorumStakeTotals)).(*IBLSSignatureCheckerQuorumStakeTotals)
	out1 := *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)

	return out0, out1, err

}

// CheckSignatures is a free data retrieval call binding the contract method 0x6efb4636.
//
// Solidity: function checkSignatures(bytes32 msgHash, bytes quorumNumbers, uint32 referenceBlockNumber, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) params) view returns((uint96[],uint96[]), bytes32)
func (_CedroServiceManager *CedroServiceManagerSession) CheckSignatures(msgHash [32]byte, quorumNumbers []byte, referenceBlockNumber uint32, params IBLSSignatureCheckerNonSignerStakesAndSignature) (IBLSSignatureCheckerQuorumStakeTotals, [32]byte, error) {
	return _CedroServiceManager.Contract.CheckSignatures(&_CedroServiceManager.CallOpts, msgHash, quorumNumbers, referenceBlockNumber, params)
}

// CheckSignatures is a free data retrieval call binding the contract method 0x6efb4636.
//
// Solidity: function checkSignatures(bytes32 msgHash, bytes quorumNumbers, uint32 referenceBlockNumber, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) params) view returns((uint96[],uint96[]), bytes32)
func (_CedroServiceManager *CedroServiceManagerCallerSession) CheckSignatures(msgHash [32]byte, quorumNumbers []byte, referenceBlockNumber uint32, params IBLSSignatureCheckerNonSignerStakesAndSignature) (IBLSSignatureCheckerQuorumStakeTotals, [32]byte, error) {
	return _CedroServiceManager.Contract.CheckSignatures(&_CedroServiceManager.CallOpts, msgHash, quorumNumbers, referenceBlockNumber, params)
}

// Contains is a free data retrieval call binding the contract method 0x1d1a696d.
//
// Solidity: function contains(bytes32 messageHash) view returns(bool)
func (_CedroServiceManager *CedroServiceManagerCaller) Contains(opts *bind.CallOpts, messageHash [32]byte) (bool, error) {
	var out []interface{}
	err := _CedroServiceManager.contract.Call(opts, &out, "contains", messageHash)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Contains is a free data retrieval call binding the contract method 0x1d1a696d.
//
// Solidity: function contains(bytes32 messageHash) view returns(bool)
func (_CedroServiceManager *CedroServiceManagerSession) Contains(messageHash [32]byte) (bool, error) {
	return _CedroServiceManager.Contract.Contains(&_CedroServiceManager.CallOpts, messageHash)
}

// Contains is a free data retrieval call binding the contract method 0x1d1a696d.
//
// Solidity: function contains(bytes32 messageHash) view returns(bool)
func (_CedroServiceManager *CedroServiceManagerCallerSession) Contains(messageHash [32]byte) (bool, error) {
	return _CedroServiceManager.Contract.Contains(&_CedroServiceManager.CallOpts, messageHash)
}

// Delegation is a free data retrieval call binding the contract method 0xdf5cf723.
//
// Solidity: function delegation() view returns(address)
func (_CedroServiceManager *CedroServiceManagerCaller) Delegation(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CedroServiceManager.contract.Call(opts, &out, "delegation")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Delegation is a free data retrieval call binding the contract method 0xdf5cf723.
//
// Solidity: function delegation() view returns(address)
func (_CedroServiceManager *CedroServiceManagerSession) Delegation() (common.Address, error) {
	return _CedroServiceManager.Contract.Delegation(&_CedroServiceManager.CallOpts)
}

// Delegation is a free data retrieval call binding the contract method 0xdf5cf723.
//
// Solidity: function delegation() view returns(address)
func (_CedroServiceManager *CedroServiceManagerCallerSession) Delegation() (common.Address, error) {
	return _CedroServiceManager.Contract.Delegation(&_CedroServiceManager.CallOpts)
}

// GetOperatorRestakedStrategies is a free data retrieval call binding the contract method 0x33cfb7b7.
//
// Solidity: function getOperatorRestakedStrategies(address operator) view returns(address[])
func (_CedroServiceManager *CedroServiceManagerCaller) GetOperatorRestakedStrategies(opts *bind.CallOpts, operator common.Address) ([]common.Address, error) {
	var out []interface{}
	err := _CedroServiceManager.contract.Call(opts, &out, "getOperatorRestakedStrategies", operator)

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetOperatorRestakedStrategies is a free data retrieval call binding the contract method 0x33cfb7b7.
//
// Solidity: function getOperatorRestakedStrategies(address operator) view returns(address[])
func (_CedroServiceManager *CedroServiceManagerSession) GetOperatorRestakedStrategies(operator common.Address) ([]common.Address, error) {
	return _CedroServiceManager.Contract.GetOperatorRestakedStrategies(&_CedroServiceManager.CallOpts, operator)
}

// GetOperatorRestakedStrategies is a free data retrieval call binding the contract method 0x33cfb7b7.
//
// Solidity: function getOperatorRestakedStrategies(address operator) view returns(address[])
func (_CedroServiceManager *CedroServiceManagerCallerSession) GetOperatorRestakedStrategies(operator common.Address) ([]common.Address, error) {
	return _CedroServiceManager.Contract.GetOperatorRestakedStrategies(&_CedroServiceManager.CallOpts, operator)
}

// GetRestakeableStrategies is a free data retrieval call binding the contract method 0xe481af9d.
//
// Solidity: function getRestakeableStrategies() view returns(address[])
func (_CedroServiceManager *CedroServiceManagerCaller) GetRestakeableStrategies(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _CedroServiceManager.contract.Call(opts, &out, "getRestakeableStrategies")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetRestakeableStrategies is a free data retrieval call binding the contract method 0xe481af9d.
//
// Solidity: function getRestakeableStrategies() view returns(address[])
func (_CedroServiceManager *CedroServiceManagerSession) GetRestakeableStrategies() ([]common.Address, error) {
	return _CedroServiceManager.Contract.GetRestakeableStrategies(&_CedroServiceManager.CallOpts)
}

// GetRestakeableStrategies is a free data retrieval call binding the contract method 0xe481af9d.
//
// Solidity: function getRestakeableStrategies() view returns(address[])
func (_CedroServiceManager *CedroServiceManagerCallerSession) GetRestakeableStrategies() ([]common.Address, error) {
	return _CedroServiceManager.Contract.GetRestakeableStrategies(&_CedroServiceManager.CallOpts)
}

// LatestServeUntilBlock is a free data retrieval call binding the contract method 0xeaefd27d.
//
// Solidity: function latestServeUntilBlock(uint32 referenceBlockNumber) view returns(uint32)
func (_CedroServiceManager *CedroServiceManagerCaller) LatestServeUntilBlock(opts *bind.CallOpts, referenceBlockNumber uint32) (uint32, error) {
	var out []interface{}
	err := _CedroServiceManager.contract.Call(opts, &out, "latestServeUntilBlock", referenceBlockNumber)

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// LatestServeUntilBlock is a free data retrieval call binding the contract method 0xeaefd27d.
//
// Solidity: function latestServeUntilBlock(uint32 referenceBlockNumber) view returns(uint32)
func (_CedroServiceManager *CedroServiceManagerSession) LatestServeUntilBlock(referenceBlockNumber uint32) (uint32, error) {
	return _CedroServiceManager.Contract.LatestServeUntilBlock(&_CedroServiceManager.CallOpts, referenceBlockNumber)
}

// LatestServeUntilBlock is a free data retrieval call binding the contract method 0xeaefd27d.
//
// Solidity: function latestServeUntilBlock(uint32 referenceBlockNumber) view returns(uint32)
func (_CedroServiceManager *CedroServiceManagerCallerSession) LatestServeUntilBlock(referenceBlockNumber uint32) (uint32, error) {
	return _CedroServiceManager.Contract.LatestServeUntilBlock(&_CedroServiceManager.CallOpts, referenceBlockNumber)
}

// OperatorAllowlist is a free data retrieval call binding the contract method 0xaa99c067.
//
// Solidity: function operatorAllowlist(address ) view returns(bool)
func (_CedroServiceManager *CedroServiceManagerCaller) OperatorAllowlist(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _CedroServiceManager.contract.Call(opts, &out, "operatorAllowlist", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// OperatorAllowlist is a free data retrieval call binding the contract method 0xaa99c067.
//
// Solidity: function operatorAllowlist(address ) view returns(bool)
func (_CedroServiceManager *CedroServiceManagerSession) OperatorAllowlist(arg0 common.Address) (bool, error) {
	return _CedroServiceManager.Contract.OperatorAllowlist(&_CedroServiceManager.CallOpts, arg0)
}

// OperatorAllowlist is a free data retrieval call binding the contract method 0xaa99c067.
//
// Solidity: function operatorAllowlist(address ) view returns(bool)
func (_CedroServiceManager *CedroServiceManagerCallerSession) OperatorAllowlist(arg0 common.Address) (bool, error) {
	return _CedroServiceManager.Contract.OperatorAllowlist(&_CedroServiceManager.CallOpts, arg0)
}

// QuorumThresholdPercentage is a free data retrieval call binding the contract method 0x4deabc21.
//
// Solidity: function quorumThresholdPercentage() view returns(uint8)
func (_CedroServiceManager *CedroServiceManagerCaller) QuorumThresholdPercentage(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _CedroServiceManager.contract.Call(opts, &out, "quorumThresholdPercentage")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// QuorumThresholdPercentage is a free data retrieval call binding the contract method 0x4deabc21.
//
// Solidity: function quorumThresholdPercentage() view returns(uint8)
func (_CedroServiceManager *CedroServiceManagerSession) QuorumThresholdPercentage() (uint8, error) {
	return _CedroServiceManager.Contract.QuorumThresholdPercentage(&_CedroServiceManager.CallOpts)
}

// QuorumThresholdPercentage is a free data retrieval call binding the contract method 0x4deabc21.
//
// Solidity: function quorumThresholdPercentage() view returns(uint8)
func (_CedroServiceManager *CedroServiceManagerCallerSession) QuorumThresholdPercentage() (uint8, error) {
	return _CedroServiceManager.Contract.QuorumThresholdPercentage(&_CedroServiceManager.CallOpts)
}

// RegistryCoordinator is a free data retrieval call binding the contract method 0x6d14a987.
//
// Solidity: function registryCoordinator() view returns(address)
func (_CedroServiceManager *CedroServiceManagerCaller) RegistryCoordinator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CedroServiceManager.contract.Call(opts, &out, "registryCoordinator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RegistryCoordinator is a free data retrieval call binding the contract method 0x6d14a987.
//
// Solidity: function registryCoordinator() view returns(address)
func (_CedroServiceManager *CedroServiceManagerSession) RegistryCoordinator() (common.Address, error) {
	return _CedroServiceManager.Contract.RegistryCoordinator(&_CedroServiceManager.CallOpts)
}

// RegistryCoordinator is a free data retrieval call binding the contract method 0x6d14a987.
//
// Solidity: function registryCoordinator() view returns(address)
func (_CedroServiceManager *CedroServiceManagerCallerSession) RegistryCoordinator() (common.Address, error) {
	return _CedroServiceManager.Contract.RegistryCoordinator(&_CedroServiceManager.CallOpts)
}

// RewardsInitiator is a free data retrieval call binding the contract method 0xfc299dee.
//
// Solidity: function rewardsInitiator() view returns(address)
func (_CedroServiceManager *CedroServiceManagerCaller) RewardsInitiator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CedroServiceManager.contract.Call(opts, &out, "rewardsInitiator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RewardsInitiator is a free data retrieval call binding the contract method 0xfc299dee.
//
// Solidity: function rewardsInitiator() view returns(address)
func (_CedroServiceManager *CedroServiceManagerSession) RewardsInitiator() (common.Address, error) {
	return _CedroServiceManager.Contract.RewardsInitiator(&_CedroServiceManager.CallOpts)
}

// RewardsInitiator is a free data retrieval call binding the contract method 0xfc299dee.
//
// Solidity: function rewardsInitiator() view returns(address)
func (_CedroServiceManager *CedroServiceManagerCallerSession) RewardsInitiator() (common.Address, error) {
	return _CedroServiceManager.Contract.RewardsInitiator(&_CedroServiceManager.CallOpts)
}

// RoleManager is a free data retrieval call binding the contract method 0x00435da5.
//
// Solidity: function roleManager() view returns(address)
func (_CedroServiceManager *CedroServiceManagerCaller) RoleManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CedroServiceManager.contract.Call(opts, &out, "roleManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RoleManager is a free data retrieval call binding the contract method 0x00435da5.
//
// Solidity: function roleManager() view returns(address)
func (_CedroServiceManager *CedroServiceManagerSession) RoleManager() (common.Address, error) {
	return _CedroServiceManager.Contract.RoleManager(&_CedroServiceManager.CallOpts)
}

// RoleManager is a free data retrieval call binding the contract method 0x00435da5.
//
// Solidity: function roleManager() view returns(address)
func (_CedroServiceManager *CedroServiceManagerCallerSession) RoleManager() (common.Address, error) {
	return _CedroServiceManager.Contract.RoleManager(&_CedroServiceManager.CallOpts)
}

// RoleManagerAddress is a free data retrieval call binding the contract method 0xd33b8d23.
//
// Solidity: function roleManagerAddress() view returns(address)
func (_CedroServiceManager *CedroServiceManagerCaller) RoleManagerAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CedroServiceManager.contract.Call(opts, &out, "roleManagerAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RoleManagerAddress is a free data retrieval call binding the contract method 0xd33b8d23.
//
// Solidity: function roleManagerAddress() view returns(address)
func (_CedroServiceManager *CedroServiceManagerSession) RoleManagerAddress() (common.Address, error) {
	return _CedroServiceManager.Contract.RoleManagerAddress(&_CedroServiceManager.CallOpts)
}

// RoleManagerAddress is a free data retrieval call binding the contract method 0xd33b8d23.
//
// Solidity: function roleManagerAddress() view returns(address)
func (_CedroServiceManager *CedroServiceManagerCallerSession) RoleManagerAddress() (common.Address, error) {
	return _CedroServiceManager.Contract.RoleManagerAddress(&_CedroServiceManager.CallOpts)
}

// StakeRegistry is a free data retrieval call binding the contract method 0x68304835.
//
// Solidity: function stakeRegistry() view returns(address)
func (_CedroServiceManager *CedroServiceManagerCaller) StakeRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CedroServiceManager.contract.Call(opts, &out, "stakeRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// StakeRegistry is a free data retrieval call binding the contract method 0x68304835.
//
// Solidity: function stakeRegistry() view returns(address)
func (_CedroServiceManager *CedroServiceManagerSession) StakeRegistry() (common.Address, error) {
	return _CedroServiceManager.Contract.StakeRegistry(&_CedroServiceManager.CallOpts)
}

// StakeRegistry is a free data retrieval call binding the contract method 0x68304835.
//
// Solidity: function stakeRegistry() view returns(address)
func (_CedroServiceManager *CedroServiceManagerCallerSession) StakeRegistry() (common.Address, error) {
	return _CedroServiceManager.Contract.StakeRegistry(&_CedroServiceManager.CallOpts)
}

// StaleStakesForbidden is a free data retrieval call binding the contract method 0xb98d0908.
//
// Solidity: function staleStakesForbidden() view returns(bool)
func (_CedroServiceManager *CedroServiceManagerCaller) StaleStakesForbidden(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _CedroServiceManager.contract.Call(opts, &out, "staleStakesForbidden")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// StaleStakesForbidden is a free data retrieval call binding the contract method 0xb98d0908.
//
// Solidity: function staleStakesForbidden() view returns(bool)
func (_CedroServiceManager *CedroServiceManagerSession) StaleStakesForbidden() (bool, error) {
	return _CedroServiceManager.Contract.StaleStakesForbidden(&_CedroServiceManager.CallOpts)
}

// StaleStakesForbidden is a free data retrieval call binding the contract method 0xb98d0908.
//
// Solidity: function staleStakesForbidden() view returns(bool)
func (_CedroServiceManager *CedroServiceManagerCallerSession) StaleStakesForbidden() (bool, error) {
	return _CedroServiceManager.Contract.StaleStakesForbidden(&_CedroServiceManager.CallOpts)
}

// TrySignatureAndApkVerification is a free data retrieval call binding the contract method 0x171f1d5b.
//
// Solidity: function trySignatureAndApkVerification(bytes32 msgHash, (uint256,uint256) apk, (uint256[2],uint256[2]) apkG2, (uint256,uint256) sigma) view returns(bool pairingSuccessful, bool siganatureIsValid)
func (_CedroServiceManager *CedroServiceManagerCaller) TrySignatureAndApkVerification(opts *bind.CallOpts, msgHash [32]byte, apk BN254G1Point, apkG2 BN254G2Point, sigma BN254G1Point) (struct {
	PairingSuccessful bool
	SiganatureIsValid bool
}, error) {
	var out []interface{}
	err := _CedroServiceManager.contract.Call(opts, &out, "trySignatureAndApkVerification", msgHash, apk, apkG2, sigma)

	outstruct := new(struct {
		PairingSuccessful bool
		SiganatureIsValid bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.PairingSuccessful = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.SiganatureIsValid = *abi.ConvertType(out[1], new(bool)).(*bool)

	return *outstruct, err

}

// TrySignatureAndApkVerification is a free data retrieval call binding the contract method 0x171f1d5b.
//
// Solidity: function trySignatureAndApkVerification(bytes32 msgHash, (uint256,uint256) apk, (uint256[2],uint256[2]) apkG2, (uint256,uint256) sigma) view returns(bool pairingSuccessful, bool siganatureIsValid)
func (_CedroServiceManager *CedroServiceManagerSession) TrySignatureAndApkVerification(msgHash [32]byte, apk BN254G1Point, apkG2 BN254G2Point, sigma BN254G1Point) (struct {
	PairingSuccessful bool
	SiganatureIsValid bool
}, error) {
	return _CedroServiceManager.Contract.TrySignatureAndApkVerification(&_CedroServiceManager.CallOpts, msgHash, apk, apkG2, sigma)
}

// TrySignatureAndApkVerification is a free data retrieval call binding the contract method 0x171f1d5b.
//
// Solidity: function trySignatureAndApkVerification(bytes32 msgHash, (uint256,uint256) apk, (uint256[2],uint256[2]) apkG2, (uint256,uint256) sigma) view returns(bool pairingSuccessful, bool siganatureIsValid)
func (_CedroServiceManager *CedroServiceManagerCallerSession) TrySignatureAndApkVerification(msgHash [32]byte, apk BN254G1Point, apkG2 BN254G2Point, sigma BN254G1Point) (struct {
	PairingSuccessful bool
	SiganatureIsValid bool
}, error) {
	return _CedroServiceManager.Contract.TrySignatureAndApkVerification(&_CedroServiceManager.CallOpts, msgHash, apk, apkG2, sigma)
}

// AddOperatorsToAllowlist is a paid mutator transaction binding the contract method 0x435a6f28.
//
// Solidity: function addOperatorsToAllowlist(address[] operators) returns()
func (_CedroServiceManager *CedroServiceManagerTransactor) AddOperatorsToAllowlist(opts *bind.TransactOpts, operators []common.Address) (*types.Transaction, error) {
	return _CedroServiceManager.contract.Transact(opts, "addOperatorsToAllowlist", operators)
}

// AddOperatorsToAllowlist is a paid mutator transaction binding the contract method 0x435a6f28.
//
// Solidity: function addOperatorsToAllowlist(address[] operators) returns()
func (_CedroServiceManager *CedroServiceManagerSession) AddOperatorsToAllowlist(operators []common.Address) (*types.Transaction, error) {
	return _CedroServiceManager.Contract.AddOperatorsToAllowlist(&_CedroServiceManager.TransactOpts, operators)
}

// AddOperatorsToAllowlist is a paid mutator transaction binding the contract method 0x435a6f28.
//
// Solidity: function addOperatorsToAllowlist(address[] operators) returns()
func (_CedroServiceManager *CedroServiceManagerTransactorSession) AddOperatorsToAllowlist(operators []common.Address) (*types.Transaction, error) {
	return _CedroServiceManager.Contract.AddOperatorsToAllowlist(&_CedroServiceManager.TransactOpts, operators)
}

// ChallengeOperators is a paid mutator transaction binding the contract method 0x248b3979.
//
// Solidity: function challengeOperators((bytes32) evidence) returns()
func (_CedroServiceManager *CedroServiceManagerTransactor) ChallengeOperators(opts *bind.TransactOpts, evidence CedroServiceDataOperatorChallengeEvidence) (*types.Transaction, error) {
	return _CedroServiceManager.contract.Transact(opts, "challengeOperators", evidence)
}

// ChallengeOperators is a paid mutator transaction binding the contract method 0x248b3979.
//
// Solidity: function challengeOperators((bytes32) evidence) returns()
func (_CedroServiceManager *CedroServiceManagerSession) ChallengeOperators(evidence CedroServiceDataOperatorChallengeEvidence) (*types.Transaction, error) {
	return _CedroServiceManager.Contract.ChallengeOperators(&_CedroServiceManager.TransactOpts, evidence)
}

// ChallengeOperators is a paid mutator transaction binding the contract method 0x248b3979.
//
// Solidity: function challengeOperators((bytes32) evidence) returns()
func (_CedroServiceManager *CedroServiceManagerTransactorSession) ChallengeOperators(evidence CedroServiceDataOperatorChallengeEvidence) (*types.Transaction, error) {
	return _CedroServiceManager.Contract.ChallengeOperators(&_CedroServiceManager.TransactOpts, evidence)
}

// ConfirmAgentDeviation is a paid mutator transaction binding the contract method 0x1008ef92.
//
// Solidity: function confirmAgentDeviation((address,uint32,bytes,bytes,uint256,uint256,uint256) deviationBalanceHeader, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature) returns()
func (_CedroServiceManager *CedroServiceManagerTransactor) ConfirmAgentDeviation(opts *bind.TransactOpts, deviationBalanceHeader CedroServiceDataDeviationBalanceHeader, nonSignerStakesAndSignature IBLSSignatureCheckerNonSignerStakesAndSignature) (*types.Transaction, error) {
	return _CedroServiceManager.contract.Transact(opts, "confirmAgentDeviation", deviationBalanceHeader, nonSignerStakesAndSignature)
}

// ConfirmAgentDeviation is a paid mutator transaction binding the contract method 0x1008ef92.
//
// Solidity: function confirmAgentDeviation((address,uint32,bytes,bytes,uint256,uint256,uint256) deviationBalanceHeader, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature) returns()
func (_CedroServiceManager *CedroServiceManagerSession) ConfirmAgentDeviation(deviationBalanceHeader CedroServiceDataDeviationBalanceHeader, nonSignerStakesAndSignature IBLSSignatureCheckerNonSignerStakesAndSignature) (*types.Transaction, error) {
	return _CedroServiceManager.Contract.ConfirmAgentDeviation(&_CedroServiceManager.TransactOpts, deviationBalanceHeader, nonSignerStakesAndSignature)
}

// ConfirmAgentDeviation is a paid mutator transaction binding the contract method 0x1008ef92.
//
// Solidity: function confirmAgentDeviation((address,uint32,bytes,bytes,uint256,uint256,uint256) deviationBalanceHeader, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature) returns()
func (_CedroServiceManager *CedroServiceManagerTransactorSession) ConfirmAgentDeviation(deviationBalanceHeader CedroServiceDataDeviationBalanceHeader, nonSignerStakesAndSignature IBLSSignatureCheckerNonSignerStakesAndSignature) (*types.Transaction, error) {
	return _CedroServiceManager.Contract.ConfirmAgentDeviation(&_CedroServiceManager.TransactOpts, deviationBalanceHeader, nonSignerStakesAndSignature)
}

// ConfirmAgentFraud is a paid mutator transaction binding the contract method 0x2c1e244d.
//
// Solidity: function confirmAgentFraud((address,uint32,bytes,bytes,uint256,uint256) fraudBalanceHeader, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature) returns()
func (_CedroServiceManager *CedroServiceManagerTransactor) ConfirmAgentFraud(opts *bind.TransactOpts, fraudBalanceHeader CedroServiceDataFraudBalanceHeader, nonSignerStakesAndSignature IBLSSignatureCheckerNonSignerStakesAndSignature) (*types.Transaction, error) {
	return _CedroServiceManager.contract.Transact(opts, "confirmAgentFraud", fraudBalanceHeader, nonSignerStakesAndSignature)
}

// ConfirmAgentFraud is a paid mutator transaction binding the contract method 0x2c1e244d.
//
// Solidity: function confirmAgentFraud((address,uint32,bytes,bytes,uint256,uint256) fraudBalanceHeader, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature) returns()
func (_CedroServiceManager *CedroServiceManagerSession) ConfirmAgentFraud(fraudBalanceHeader CedroServiceDataFraudBalanceHeader, nonSignerStakesAndSignature IBLSSignatureCheckerNonSignerStakesAndSignature) (*types.Transaction, error) {
	return _CedroServiceManager.Contract.ConfirmAgentFraud(&_CedroServiceManager.TransactOpts, fraudBalanceHeader, nonSignerStakesAndSignature)
}

// ConfirmAgentFraud is a paid mutator transaction binding the contract method 0x2c1e244d.
//
// Solidity: function confirmAgentFraud((address,uint32,bytes,bytes,uint256,uint256) fraudBalanceHeader, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature) returns()
func (_CedroServiceManager *CedroServiceManagerTransactorSession) ConfirmAgentFraud(fraudBalanceHeader CedroServiceDataFraudBalanceHeader, nonSignerStakesAndSignature IBLSSignatureCheckerNonSignerStakesAndSignature) (*types.Transaction, error) {
	return _CedroServiceManager.Contract.ConfirmAgentFraud(&_CedroServiceManager.TransactOpts, fraudBalanceHeader, nonSignerStakesAndSignature)
}

// CreateAVSRewardsSubmission is a paid mutator transaction binding the contract method 0xfce36c7d.
//
// Solidity: function createAVSRewardsSubmission(((address,uint96)[],address,uint256,uint32,uint32)[] rewardsSubmissions) returns()
func (_CedroServiceManager *CedroServiceManagerTransactor) CreateAVSRewardsSubmission(opts *bind.TransactOpts, rewardsSubmissions []IRewardsCoordinatorRewardsSubmission) (*types.Transaction, error) {
	return _CedroServiceManager.contract.Transact(opts, "createAVSRewardsSubmission", rewardsSubmissions)
}

// CreateAVSRewardsSubmission is a paid mutator transaction binding the contract method 0xfce36c7d.
//
// Solidity: function createAVSRewardsSubmission(((address,uint96)[],address,uint256,uint32,uint32)[] rewardsSubmissions) returns()
func (_CedroServiceManager *CedroServiceManagerSession) CreateAVSRewardsSubmission(rewardsSubmissions []IRewardsCoordinatorRewardsSubmission) (*types.Transaction, error) {
	return _CedroServiceManager.Contract.CreateAVSRewardsSubmission(&_CedroServiceManager.TransactOpts, rewardsSubmissions)
}

// CreateAVSRewardsSubmission is a paid mutator transaction binding the contract method 0xfce36c7d.
//
// Solidity: function createAVSRewardsSubmission(((address,uint96)[],address,uint256,uint32,uint32)[] rewardsSubmissions) returns()
func (_CedroServiceManager *CedroServiceManagerTransactorSession) CreateAVSRewardsSubmission(rewardsSubmissions []IRewardsCoordinatorRewardsSubmission) (*types.Transaction, error) {
	return _CedroServiceManager.Contract.CreateAVSRewardsSubmission(&_CedroServiceManager.TransactOpts, rewardsSubmissions)
}

// DeregisterOperatorFromAVS is a paid mutator transaction binding the contract method 0xa364f4da.
//
// Solidity: function deregisterOperatorFromAVS(address operator) returns()
func (_CedroServiceManager *CedroServiceManagerTransactor) DeregisterOperatorFromAVS(opts *bind.TransactOpts, operator common.Address) (*types.Transaction, error) {
	return _CedroServiceManager.contract.Transact(opts, "deregisterOperatorFromAVS", operator)
}

// DeregisterOperatorFromAVS is a paid mutator transaction binding the contract method 0xa364f4da.
//
// Solidity: function deregisterOperatorFromAVS(address operator) returns()
func (_CedroServiceManager *CedroServiceManagerSession) DeregisterOperatorFromAVS(operator common.Address) (*types.Transaction, error) {
	return _CedroServiceManager.Contract.DeregisterOperatorFromAVS(&_CedroServiceManager.TransactOpts, operator)
}

// DeregisterOperatorFromAVS is a paid mutator transaction binding the contract method 0xa364f4da.
//
// Solidity: function deregisterOperatorFromAVS(address operator) returns()
func (_CedroServiceManager *CedroServiceManagerTransactorSession) DeregisterOperatorFromAVS(operator common.Address) (*types.Transaction, error) {
	return _CedroServiceManager.Contract.DeregisterOperatorFromAVS(&_CedroServiceManager.TransactOpts, operator)
}

// Initialize is a paid mutator transaction binding the contract method 0xf8c8765e.
//
// Solidity: function initialize(address __roleManagerAddress, address __avsDirectory, address __rewardsCoordinator, address __stakeRegistry) returns()
func (_CedroServiceManager *CedroServiceManagerTransactor) Initialize(opts *bind.TransactOpts, __roleManagerAddress common.Address, __avsDirectory common.Address, __rewardsCoordinator common.Address, __stakeRegistry common.Address) (*types.Transaction, error) {
	return _CedroServiceManager.contract.Transact(opts, "initialize", __roleManagerAddress, __avsDirectory, __rewardsCoordinator, __stakeRegistry)
}

// Initialize is a paid mutator transaction binding the contract method 0xf8c8765e.
//
// Solidity: function initialize(address __roleManagerAddress, address __avsDirectory, address __rewardsCoordinator, address __stakeRegistry) returns()
func (_CedroServiceManager *CedroServiceManagerSession) Initialize(__roleManagerAddress common.Address, __avsDirectory common.Address, __rewardsCoordinator common.Address, __stakeRegistry common.Address) (*types.Transaction, error) {
	return _CedroServiceManager.Contract.Initialize(&_CedroServiceManager.TransactOpts, __roleManagerAddress, __avsDirectory, __rewardsCoordinator, __stakeRegistry)
}

// Initialize is a paid mutator transaction binding the contract method 0xf8c8765e.
//
// Solidity: function initialize(address __roleManagerAddress, address __avsDirectory, address __rewardsCoordinator, address __stakeRegistry) returns()
func (_CedroServiceManager *CedroServiceManagerTransactorSession) Initialize(__roleManagerAddress common.Address, __avsDirectory common.Address, __rewardsCoordinator common.Address, __stakeRegistry common.Address) (*types.Transaction, error) {
	return _CedroServiceManager.Contract.Initialize(&_CedroServiceManager.TransactOpts, __roleManagerAddress, __avsDirectory, __rewardsCoordinator, __stakeRegistry)
}

// RegisterOperatorToAVS is a paid mutator transaction binding the contract method 0x9926ee7d.
//
// Solidity: function registerOperatorToAVS(address operator, (bytes,bytes32,uint256) operatorSignature) returns()
func (_CedroServiceManager *CedroServiceManagerTransactor) RegisterOperatorToAVS(opts *bind.TransactOpts, operator common.Address, operatorSignature ISignatureUtilsSignatureWithSaltAndExpiry) (*types.Transaction, error) {
	return _CedroServiceManager.contract.Transact(opts, "registerOperatorToAVS", operator, operatorSignature)
}

// RegisterOperatorToAVS is a paid mutator transaction binding the contract method 0x9926ee7d.
//
// Solidity: function registerOperatorToAVS(address operator, (bytes,bytes32,uint256) operatorSignature) returns()
func (_CedroServiceManager *CedroServiceManagerSession) RegisterOperatorToAVS(operator common.Address, operatorSignature ISignatureUtilsSignatureWithSaltAndExpiry) (*types.Transaction, error) {
	return _CedroServiceManager.Contract.RegisterOperatorToAVS(&_CedroServiceManager.TransactOpts, operator, operatorSignature)
}

// RegisterOperatorToAVS is a paid mutator transaction binding the contract method 0x9926ee7d.
//
// Solidity: function registerOperatorToAVS(address operator, (bytes,bytes32,uint256) operatorSignature) returns()
func (_CedroServiceManager *CedroServiceManagerTransactorSession) RegisterOperatorToAVS(operator common.Address, operatorSignature ISignatureUtilsSignatureWithSaltAndExpiry) (*types.Transaction, error) {
	return _CedroServiceManager.Contract.RegisterOperatorToAVS(&_CedroServiceManager.TransactOpts, operator, operatorSignature)
}

// SetAgentManagerAddress is a paid mutator transaction binding the contract method 0xda35a805.
//
// Solidity: function setAgentManagerAddress(address _agentManagerAddress) returns()
func (_CedroServiceManager *CedroServiceManagerTransactor) SetAgentManagerAddress(opts *bind.TransactOpts, _agentManagerAddress common.Address) (*types.Transaction, error) {
	return _CedroServiceManager.contract.Transact(opts, "setAgentManagerAddress", _agentManagerAddress)
}

// SetAgentManagerAddress is a paid mutator transaction binding the contract method 0xda35a805.
//
// Solidity: function setAgentManagerAddress(address _agentManagerAddress) returns()
func (_CedroServiceManager *CedroServiceManagerSession) SetAgentManagerAddress(_agentManagerAddress common.Address) (*types.Transaction, error) {
	return _CedroServiceManager.Contract.SetAgentManagerAddress(&_CedroServiceManager.TransactOpts, _agentManagerAddress)
}

// SetAgentManagerAddress is a paid mutator transaction binding the contract method 0xda35a805.
//
// Solidity: function setAgentManagerAddress(address _agentManagerAddress) returns()
func (_CedroServiceManager *CedroServiceManagerTransactorSession) SetAgentManagerAddress(_agentManagerAddress common.Address) (*types.Transaction, error) {
	return _CedroServiceManager.Contract.SetAgentManagerAddress(&_CedroServiceManager.TransactOpts, _agentManagerAddress)
}

// SetPauseStatus is a paid mutator transaction binding the contract method 0xc38bb537.
//
// Solidity: function setPauseStatus(bool status) returns()
func (_CedroServiceManager *CedroServiceManagerTransactor) SetPauseStatus(opts *bind.TransactOpts, status bool) (*types.Transaction, error) {
	return _CedroServiceManager.contract.Transact(opts, "setPauseStatus", status)
}

// SetPauseStatus is a paid mutator transaction binding the contract method 0xc38bb537.
//
// Solidity: function setPauseStatus(bool status) returns()
func (_CedroServiceManager *CedroServiceManagerSession) SetPauseStatus(status bool) (*types.Transaction, error) {
	return _CedroServiceManager.Contract.SetPauseStatus(&_CedroServiceManager.TransactOpts, status)
}

// SetPauseStatus is a paid mutator transaction binding the contract method 0xc38bb537.
//
// Solidity: function setPauseStatus(bool status) returns()
func (_CedroServiceManager *CedroServiceManagerTransactorSession) SetPauseStatus(status bool) (*types.Transaction, error) {
	return _CedroServiceManager.Contract.SetPauseStatus(&_CedroServiceManager.TransactOpts, status)
}

// SetQuorumThresholdPercentage is a paid mutator transaction binding the contract method 0x3fe9f9a1.
//
// Solidity: function setQuorumThresholdPercentage(uint8 thresholdPercentage) returns()
func (_CedroServiceManager *CedroServiceManagerTransactor) SetQuorumThresholdPercentage(opts *bind.TransactOpts, thresholdPercentage uint8) (*types.Transaction, error) {
	return _CedroServiceManager.contract.Transact(opts, "setQuorumThresholdPercentage", thresholdPercentage)
}

// SetQuorumThresholdPercentage is a paid mutator transaction binding the contract method 0x3fe9f9a1.
//
// Solidity: function setQuorumThresholdPercentage(uint8 thresholdPercentage) returns()
func (_CedroServiceManager *CedroServiceManagerSession) SetQuorumThresholdPercentage(thresholdPercentage uint8) (*types.Transaction, error) {
	return _CedroServiceManager.Contract.SetQuorumThresholdPercentage(&_CedroServiceManager.TransactOpts, thresholdPercentage)
}

// SetQuorumThresholdPercentage is a paid mutator transaction binding the contract method 0x3fe9f9a1.
//
// Solidity: function setQuorumThresholdPercentage(uint8 thresholdPercentage) returns()
func (_CedroServiceManager *CedroServiceManagerTransactorSession) SetQuorumThresholdPercentage(thresholdPercentage uint8) (*types.Transaction, error) {
	return _CedroServiceManager.Contract.SetQuorumThresholdPercentage(&_CedroServiceManager.TransactOpts, thresholdPercentage)
}

// SetRegistryCoordinatorAddress is a paid mutator transaction binding the contract method 0x947cd465.
//
// Solidity: function setRegistryCoordinatorAddress(address __registryCoordinator) returns()
func (_CedroServiceManager *CedroServiceManagerTransactor) SetRegistryCoordinatorAddress(opts *bind.TransactOpts, __registryCoordinator common.Address) (*types.Transaction, error) {
	return _CedroServiceManager.contract.Transact(opts, "setRegistryCoordinatorAddress", __registryCoordinator)
}

// SetRegistryCoordinatorAddress is a paid mutator transaction binding the contract method 0x947cd465.
//
// Solidity: function setRegistryCoordinatorAddress(address __registryCoordinator) returns()
func (_CedroServiceManager *CedroServiceManagerSession) SetRegistryCoordinatorAddress(__registryCoordinator common.Address) (*types.Transaction, error) {
	return _CedroServiceManager.Contract.SetRegistryCoordinatorAddress(&_CedroServiceManager.TransactOpts, __registryCoordinator)
}

// SetRegistryCoordinatorAddress is a paid mutator transaction binding the contract method 0x947cd465.
//
// Solidity: function setRegistryCoordinatorAddress(address __registryCoordinator) returns()
func (_CedroServiceManager *CedroServiceManagerTransactorSession) SetRegistryCoordinatorAddress(__registryCoordinator common.Address) (*types.Transaction, error) {
	return _CedroServiceManager.Contract.SetRegistryCoordinatorAddress(&_CedroServiceManager.TransactOpts, __registryCoordinator)
}

// SetRewardsInitiator is a paid mutator transaction binding the contract method 0x3bc28c8c.
//
// Solidity: function setRewardsInitiator(address newRewardsInitiator) returns()
func (_CedroServiceManager *CedroServiceManagerTransactor) SetRewardsInitiator(opts *bind.TransactOpts, newRewardsInitiator common.Address) (*types.Transaction, error) {
	return _CedroServiceManager.contract.Transact(opts, "setRewardsInitiator", newRewardsInitiator)
}

// SetRewardsInitiator is a paid mutator transaction binding the contract method 0x3bc28c8c.
//
// Solidity: function setRewardsInitiator(address newRewardsInitiator) returns()
func (_CedroServiceManager *CedroServiceManagerSession) SetRewardsInitiator(newRewardsInitiator common.Address) (*types.Transaction, error) {
	return _CedroServiceManager.Contract.SetRewardsInitiator(&_CedroServiceManager.TransactOpts, newRewardsInitiator)
}

// SetRewardsInitiator is a paid mutator transaction binding the contract method 0x3bc28c8c.
//
// Solidity: function setRewardsInitiator(address newRewardsInitiator) returns()
func (_CedroServiceManager *CedroServiceManagerTransactorSession) SetRewardsInitiator(newRewardsInitiator common.Address) (*types.Transaction, error) {
	return _CedroServiceManager.Contract.SetRewardsInitiator(&_CedroServiceManager.TransactOpts, newRewardsInitiator)
}

// SetRoleManager is a paid mutator transaction binding the contract method 0xf1d588c5.
//
// Solidity: function setRoleManager(address _roleManagerAddress) returns()
func (_CedroServiceManager *CedroServiceManagerTransactor) SetRoleManager(opts *bind.TransactOpts, _roleManagerAddress common.Address) (*types.Transaction, error) {
	return _CedroServiceManager.contract.Transact(opts, "setRoleManager", _roleManagerAddress)
}

// SetRoleManager is a paid mutator transaction binding the contract method 0xf1d588c5.
//
// Solidity: function setRoleManager(address _roleManagerAddress) returns()
func (_CedroServiceManager *CedroServiceManagerSession) SetRoleManager(_roleManagerAddress common.Address) (*types.Transaction, error) {
	return _CedroServiceManager.Contract.SetRoleManager(&_CedroServiceManager.TransactOpts, _roleManagerAddress)
}

// SetRoleManager is a paid mutator transaction binding the contract method 0xf1d588c5.
//
// Solidity: function setRoleManager(address _roleManagerAddress) returns()
func (_CedroServiceManager *CedroServiceManagerTransactorSession) SetRoleManager(_roleManagerAddress common.Address) (*types.Transaction, error) {
	return _CedroServiceManager.Contract.SetRoleManager(&_CedroServiceManager.TransactOpts, _roleManagerAddress)
}

// SetStaleStakesForbidden is a paid mutator transaction binding the contract method 0x416c7e5e.
//
// Solidity: function setStaleStakesForbidden(bool value) returns()
func (_CedroServiceManager *CedroServiceManagerTransactor) SetStaleStakesForbidden(opts *bind.TransactOpts, value bool) (*types.Transaction, error) {
	return _CedroServiceManager.contract.Transact(opts, "setStaleStakesForbidden", value)
}

// SetStaleStakesForbidden is a paid mutator transaction binding the contract method 0x416c7e5e.
//
// Solidity: function setStaleStakesForbidden(bool value) returns()
func (_CedroServiceManager *CedroServiceManagerSession) SetStaleStakesForbidden(value bool) (*types.Transaction, error) {
	return _CedroServiceManager.Contract.SetStaleStakesForbidden(&_CedroServiceManager.TransactOpts, value)
}

// SetStaleStakesForbidden is a paid mutator transaction binding the contract method 0x416c7e5e.
//
// Solidity: function setStaleStakesForbidden(bool value) returns()
func (_CedroServiceManager *CedroServiceManagerTransactorSession) SetStaleStakesForbidden(value bool) (*types.Transaction, error) {
	return _CedroServiceManager.Contract.SetStaleStakesForbidden(&_CedroServiceManager.TransactOpts, value)
}

// UpdateAVSMetadataURI is a paid mutator transaction binding the contract method 0xa98fb355.
//
// Solidity: function updateAVSMetadataURI(string _metadataURI) returns()
func (_CedroServiceManager *CedroServiceManagerTransactor) UpdateAVSMetadataURI(opts *bind.TransactOpts, _metadataURI string) (*types.Transaction, error) {
	return _CedroServiceManager.contract.Transact(opts, "updateAVSMetadataURI", _metadataURI)
}

// UpdateAVSMetadataURI is a paid mutator transaction binding the contract method 0xa98fb355.
//
// Solidity: function updateAVSMetadataURI(string _metadataURI) returns()
func (_CedroServiceManager *CedroServiceManagerSession) UpdateAVSMetadataURI(_metadataURI string) (*types.Transaction, error) {
	return _CedroServiceManager.Contract.UpdateAVSMetadataURI(&_CedroServiceManager.TransactOpts, _metadataURI)
}

// UpdateAVSMetadataURI is a paid mutator transaction binding the contract method 0xa98fb355.
//
// Solidity: function updateAVSMetadataURI(string _metadataURI) returns()
func (_CedroServiceManager *CedroServiceManagerTransactorSession) UpdateAVSMetadataURI(_metadataURI string) (*types.Transaction, error) {
	return _CedroServiceManager.Contract.UpdateAVSMetadataURI(&_CedroServiceManager.TransactOpts, _metadataURI)
}

// CedroServiceManagerAgentBalanceConfirmedIterator is returned from FilterAgentBalanceConfirmed and is used to iterate over the raw logs and unpacked data for AgentBalanceConfirmed events raised by the CedroServiceManager contract.
type CedroServiceManagerAgentBalanceConfirmedIterator struct {
	Event *CedroServiceManagerAgentBalanceConfirmed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CedroServiceManagerAgentBalanceConfirmedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CedroServiceManagerAgentBalanceConfirmed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CedroServiceManagerAgentBalanceConfirmed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CedroServiceManagerAgentBalanceConfirmedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CedroServiceManagerAgentBalanceConfirmedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CedroServiceManagerAgentBalanceConfirmed represents a AgentBalanceConfirmed event raised by the CedroServiceManager contract.
type CedroServiceManagerAgentBalanceConfirmed struct {
	AgentAddress   common.Address
	SlashingAmount *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterAgentBalanceConfirmed is a free log retrieval operation binding the contract event 0x135c0f63d57042d9134244f2f8134005a2b3525c9943e69ce69c3ae1c3d611f8.
//
// Solidity: event AgentBalanceConfirmed(address agentAddress, uint256 SlashingAmount)
func (_CedroServiceManager *CedroServiceManagerFilterer) FilterAgentBalanceConfirmed(opts *bind.FilterOpts) (*CedroServiceManagerAgentBalanceConfirmedIterator, error) {

	logs, sub, err := _CedroServiceManager.contract.FilterLogs(opts, "AgentBalanceConfirmed")
	if err != nil {
		return nil, err
	}
	return &CedroServiceManagerAgentBalanceConfirmedIterator{contract: _CedroServiceManager.contract, event: "AgentBalanceConfirmed", logs: logs, sub: sub}, nil
}

// WatchAgentBalanceConfirmed is a free log subscription operation binding the contract event 0x135c0f63d57042d9134244f2f8134005a2b3525c9943e69ce69c3ae1c3d611f8.
//
// Solidity: event AgentBalanceConfirmed(address agentAddress, uint256 SlashingAmount)
func (_CedroServiceManager *CedroServiceManagerFilterer) WatchAgentBalanceConfirmed(opts *bind.WatchOpts, sink chan<- *CedroServiceManagerAgentBalanceConfirmed) (event.Subscription, error) {

	logs, sub, err := _CedroServiceManager.contract.WatchLogs(opts, "AgentBalanceConfirmed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CedroServiceManagerAgentBalanceConfirmed)
				if err := _CedroServiceManager.contract.UnpackLog(event, "AgentBalanceConfirmed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAgentBalanceConfirmed is a log parse operation binding the contract event 0x135c0f63d57042d9134244f2f8134005a2b3525c9943e69ce69c3ae1c3d611f8.
//
// Solidity: event AgentBalanceConfirmed(address agentAddress, uint256 SlashingAmount)
func (_CedroServiceManager *CedroServiceManagerFilterer) ParseAgentBalanceConfirmed(log types.Log) (*CedroServiceManagerAgentBalanceConfirmed, error) {
	event := new(CedroServiceManagerAgentBalanceConfirmed)
	if err := _CedroServiceManager.contract.UnpackLog(event, "AgentBalanceConfirmed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CedroServiceManagerAgentFraudConfirmedIterator is returned from FilterAgentFraudConfirmed and is used to iterate over the raw logs and unpacked data for AgentFraudConfirmed events raised by the CedroServiceManager contract.
type CedroServiceManagerAgentFraudConfirmedIterator struct {
	Event *CedroServiceManagerAgentFraudConfirmed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CedroServiceManagerAgentFraudConfirmedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CedroServiceManagerAgentFraudConfirmed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CedroServiceManagerAgentFraudConfirmed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CedroServiceManagerAgentFraudConfirmedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CedroServiceManagerAgentFraudConfirmedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CedroServiceManagerAgentFraudConfirmed represents a AgentFraudConfirmed event raised by the CedroServiceManager contract.
type CedroServiceManagerAgentFraudConfirmed struct {
	AgentAddress   common.Address
	SlashingAmount *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterAgentFraudConfirmed is a free log retrieval operation binding the contract event 0x5cfe7f821ee446d640013d716cfadf320f8648b26d6e7acaa7d32129a619c121.
//
// Solidity: event AgentFraudConfirmed(address agentAddress, uint256 SlashingAmount)
func (_CedroServiceManager *CedroServiceManagerFilterer) FilterAgentFraudConfirmed(opts *bind.FilterOpts) (*CedroServiceManagerAgentFraudConfirmedIterator, error) {

	logs, sub, err := _CedroServiceManager.contract.FilterLogs(opts, "AgentFraudConfirmed")
	if err != nil {
		return nil, err
	}
	return &CedroServiceManagerAgentFraudConfirmedIterator{contract: _CedroServiceManager.contract, event: "AgentFraudConfirmed", logs: logs, sub: sub}, nil
}

// WatchAgentFraudConfirmed is a free log subscription operation binding the contract event 0x5cfe7f821ee446d640013d716cfadf320f8648b26d6e7acaa7d32129a619c121.
//
// Solidity: event AgentFraudConfirmed(address agentAddress, uint256 SlashingAmount)
func (_CedroServiceManager *CedroServiceManagerFilterer) WatchAgentFraudConfirmed(opts *bind.WatchOpts, sink chan<- *CedroServiceManagerAgentFraudConfirmed) (event.Subscription, error) {

	logs, sub, err := _CedroServiceManager.contract.WatchLogs(opts, "AgentFraudConfirmed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CedroServiceManagerAgentFraudConfirmed)
				if err := _CedroServiceManager.contract.UnpackLog(event, "AgentFraudConfirmed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAgentFraudConfirmed is a log parse operation binding the contract event 0x5cfe7f821ee446d640013d716cfadf320f8648b26d6e7acaa7d32129a619c121.
//
// Solidity: event AgentFraudConfirmed(address agentAddress, uint256 SlashingAmount)
func (_CedroServiceManager *CedroServiceManagerFilterer) ParseAgentFraudConfirmed(log types.Log) (*CedroServiceManagerAgentFraudConfirmed, error) {
	event := new(CedroServiceManagerAgentFraudConfirmed)
	if err := _CedroServiceManager.contract.UnpackLog(event, "AgentFraudConfirmed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CedroServiceManagerInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the CedroServiceManager contract.
type CedroServiceManagerInitializedIterator struct {
	Event *CedroServiceManagerInitialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CedroServiceManagerInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CedroServiceManagerInitialized)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CedroServiceManagerInitialized)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CedroServiceManagerInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CedroServiceManagerInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CedroServiceManagerInitialized represents a Initialized event raised by the CedroServiceManager contract.
type CedroServiceManagerInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_CedroServiceManager *CedroServiceManagerFilterer) FilterInitialized(opts *bind.FilterOpts) (*CedroServiceManagerInitializedIterator, error) {

	logs, sub, err := _CedroServiceManager.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &CedroServiceManagerInitializedIterator{contract: _CedroServiceManager.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_CedroServiceManager *CedroServiceManagerFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *CedroServiceManagerInitialized) (event.Subscription, error) {

	logs, sub, err := _CedroServiceManager.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CedroServiceManagerInitialized)
				if err := _CedroServiceManager.contract.UnpackLog(event, "Initialized", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInitialized is a log parse operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_CedroServiceManager *CedroServiceManagerFilterer) ParseInitialized(log types.Log) (*CedroServiceManagerInitialized, error) {
	event := new(CedroServiceManagerInitialized)
	if err := _CedroServiceManager.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CedroServiceManagerOperatorDeRegisteredIterator is returned from FilterOperatorDeRegistered and is used to iterate over the raw logs and unpacked data for OperatorDeRegistered events raised by the CedroServiceManager contract.
type CedroServiceManagerOperatorDeRegisteredIterator struct {
	Event *CedroServiceManagerOperatorDeRegistered // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CedroServiceManagerOperatorDeRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CedroServiceManagerOperatorDeRegistered)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CedroServiceManagerOperatorDeRegistered)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CedroServiceManagerOperatorDeRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CedroServiceManagerOperatorDeRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CedroServiceManagerOperatorDeRegistered represents a OperatorDeRegistered event raised by the CedroServiceManager contract.
type CedroServiceManagerOperatorDeRegistered struct {
	Operator common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterOperatorDeRegistered is a free log retrieval operation binding the contract event 0xb2c38c6252ee2d17f80059fb47a790e20f7bd75e7ba577685375e5484f412d73.
//
// Solidity: event OperatorDeRegistered(address indexed operator)
func (_CedroServiceManager *CedroServiceManagerFilterer) FilterOperatorDeRegistered(opts *bind.FilterOpts, operator []common.Address) (*CedroServiceManagerOperatorDeRegisteredIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _CedroServiceManager.contract.FilterLogs(opts, "OperatorDeRegistered", operatorRule)
	if err != nil {
		return nil, err
	}
	return &CedroServiceManagerOperatorDeRegisteredIterator{contract: _CedroServiceManager.contract, event: "OperatorDeRegistered", logs: logs, sub: sub}, nil
}

// WatchOperatorDeRegistered is a free log subscription operation binding the contract event 0xb2c38c6252ee2d17f80059fb47a790e20f7bd75e7ba577685375e5484f412d73.
//
// Solidity: event OperatorDeRegistered(address indexed operator)
func (_CedroServiceManager *CedroServiceManagerFilterer) WatchOperatorDeRegistered(opts *bind.WatchOpts, sink chan<- *CedroServiceManagerOperatorDeRegistered, operator []common.Address) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _CedroServiceManager.contract.WatchLogs(opts, "OperatorDeRegistered", operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CedroServiceManagerOperatorDeRegistered)
				if err := _CedroServiceManager.contract.UnpackLog(event, "OperatorDeRegistered", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOperatorDeRegistered is a log parse operation binding the contract event 0xb2c38c6252ee2d17f80059fb47a790e20f7bd75e7ba577685375e5484f412d73.
//
// Solidity: event OperatorDeRegistered(address indexed operator)
func (_CedroServiceManager *CedroServiceManagerFilterer) ParseOperatorDeRegistered(log types.Log) (*CedroServiceManagerOperatorDeRegistered, error) {
	event := new(CedroServiceManagerOperatorDeRegistered)
	if err := _CedroServiceManager.contract.UnpackLog(event, "OperatorDeRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CedroServiceManagerOperatorRegisteredIterator is returned from FilterOperatorRegistered and is used to iterate over the raw logs and unpacked data for OperatorRegistered events raised by the CedroServiceManager contract.
type CedroServiceManagerOperatorRegisteredIterator struct {
	Event *CedroServiceManagerOperatorRegistered // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CedroServiceManagerOperatorRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CedroServiceManagerOperatorRegistered)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CedroServiceManagerOperatorRegistered)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CedroServiceManagerOperatorRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CedroServiceManagerOperatorRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CedroServiceManagerOperatorRegistered represents a OperatorRegistered event raised by the CedroServiceManager contract.
type CedroServiceManagerOperatorRegistered struct {
	Operator common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterOperatorRegistered is a free log retrieval operation binding the contract event 0x4d0eb1f4bac8744fd2be119845e23b3befc88094b42bcda1204c65694a00f9e5.
//
// Solidity: event OperatorRegistered(address indexed operator)
func (_CedroServiceManager *CedroServiceManagerFilterer) FilterOperatorRegistered(opts *bind.FilterOpts, operator []common.Address) (*CedroServiceManagerOperatorRegisteredIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _CedroServiceManager.contract.FilterLogs(opts, "OperatorRegistered", operatorRule)
	if err != nil {
		return nil, err
	}
	return &CedroServiceManagerOperatorRegisteredIterator{contract: _CedroServiceManager.contract, event: "OperatorRegistered", logs: logs, sub: sub}, nil
}

// WatchOperatorRegistered is a free log subscription operation binding the contract event 0x4d0eb1f4bac8744fd2be119845e23b3befc88094b42bcda1204c65694a00f9e5.
//
// Solidity: event OperatorRegistered(address indexed operator)
func (_CedroServiceManager *CedroServiceManagerFilterer) WatchOperatorRegistered(opts *bind.WatchOpts, sink chan<- *CedroServiceManagerOperatorRegistered, operator []common.Address) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _CedroServiceManager.contract.WatchLogs(opts, "OperatorRegistered", operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CedroServiceManagerOperatorRegistered)
				if err := _CedroServiceManager.contract.UnpackLog(event, "OperatorRegistered", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOperatorRegistered is a log parse operation binding the contract event 0x4d0eb1f4bac8744fd2be119845e23b3befc88094b42bcda1204c65694a00f9e5.
//
// Solidity: event OperatorRegistered(address indexed operator)
func (_CedroServiceManager *CedroServiceManagerFilterer) ParseOperatorRegistered(log types.Log) (*CedroServiceManagerOperatorRegistered, error) {
	event := new(CedroServiceManagerOperatorRegistered)
	if err := _CedroServiceManager.contract.UnpackLog(event, "OperatorRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CedroServiceManagerQuorumThresholdPercentageChangedIterator is returned from FilterQuorumThresholdPercentageChanged and is used to iterate over the raw logs and unpacked data for QuorumThresholdPercentageChanged events raised by the CedroServiceManager contract.
type CedroServiceManagerQuorumThresholdPercentageChangedIterator struct {
	Event *CedroServiceManagerQuorumThresholdPercentageChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CedroServiceManagerQuorumThresholdPercentageChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CedroServiceManagerQuorumThresholdPercentageChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CedroServiceManagerQuorumThresholdPercentageChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CedroServiceManagerQuorumThresholdPercentageChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CedroServiceManagerQuorumThresholdPercentageChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CedroServiceManagerQuorumThresholdPercentageChanged represents a QuorumThresholdPercentageChanged event raised by the CedroServiceManager contract.
type CedroServiceManagerQuorumThresholdPercentageChanged struct {
	Thresholdpercentage *big.Int
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterQuorumThresholdPercentageChanged is a free log retrieval operation binding the contract event 0x16f5a711b43272565a899d1435ebd000562a82f51d7ab2d9f0cf963096442a1e.
//
// Solidity: event QuorumThresholdPercentageChanged(uint256 thresholdpercentage)
func (_CedroServiceManager *CedroServiceManagerFilterer) FilterQuorumThresholdPercentageChanged(opts *bind.FilterOpts) (*CedroServiceManagerQuorumThresholdPercentageChangedIterator, error) {

	logs, sub, err := _CedroServiceManager.contract.FilterLogs(opts, "QuorumThresholdPercentageChanged")
	if err != nil {
		return nil, err
	}
	return &CedroServiceManagerQuorumThresholdPercentageChangedIterator{contract: _CedroServiceManager.contract, event: "QuorumThresholdPercentageChanged", logs: logs, sub: sub}, nil
}

// WatchQuorumThresholdPercentageChanged is a free log subscription operation binding the contract event 0x16f5a711b43272565a899d1435ebd000562a82f51d7ab2d9f0cf963096442a1e.
//
// Solidity: event QuorumThresholdPercentageChanged(uint256 thresholdpercentage)
func (_CedroServiceManager *CedroServiceManagerFilterer) WatchQuorumThresholdPercentageChanged(opts *bind.WatchOpts, sink chan<- *CedroServiceManagerQuorumThresholdPercentageChanged) (event.Subscription, error) {

	logs, sub, err := _CedroServiceManager.contract.WatchLogs(opts, "QuorumThresholdPercentageChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CedroServiceManagerQuorumThresholdPercentageChanged)
				if err := _CedroServiceManager.contract.UnpackLog(event, "QuorumThresholdPercentageChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseQuorumThresholdPercentageChanged is a log parse operation binding the contract event 0x16f5a711b43272565a899d1435ebd000562a82f51d7ab2d9f0cf963096442a1e.
//
// Solidity: event QuorumThresholdPercentageChanged(uint256 thresholdpercentage)
func (_CedroServiceManager *CedroServiceManagerFilterer) ParseQuorumThresholdPercentageChanged(log types.Log) (*CedroServiceManagerQuorumThresholdPercentageChanged, error) {
	event := new(CedroServiceManagerQuorumThresholdPercentageChanged)
	if err := _CedroServiceManager.contract.UnpackLog(event, "QuorumThresholdPercentageChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CedroServiceManagerRewardsInitiatorUpdatedIterator is returned from FilterRewardsInitiatorUpdated and is used to iterate over the raw logs and unpacked data for RewardsInitiatorUpdated events raised by the CedroServiceManager contract.
type CedroServiceManagerRewardsInitiatorUpdatedIterator struct {
	Event *CedroServiceManagerRewardsInitiatorUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CedroServiceManagerRewardsInitiatorUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CedroServiceManagerRewardsInitiatorUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CedroServiceManagerRewardsInitiatorUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CedroServiceManagerRewardsInitiatorUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CedroServiceManagerRewardsInitiatorUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CedroServiceManagerRewardsInitiatorUpdated represents a RewardsInitiatorUpdated event raised by the CedroServiceManager contract.
type CedroServiceManagerRewardsInitiatorUpdated struct {
	PrevRewardsInitiator common.Address
	NewRewardsInitiator  common.Address
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterRewardsInitiatorUpdated is a free log retrieval operation binding the contract event 0xe11cddf1816a43318ca175bbc52cd0185436e9cbead7c83acc54a73e461717e3.
//
// Solidity: event RewardsInitiatorUpdated(address prevRewardsInitiator, address newRewardsInitiator)
func (_CedroServiceManager *CedroServiceManagerFilterer) FilterRewardsInitiatorUpdated(opts *bind.FilterOpts) (*CedroServiceManagerRewardsInitiatorUpdatedIterator, error) {

	logs, sub, err := _CedroServiceManager.contract.FilterLogs(opts, "RewardsInitiatorUpdated")
	if err != nil {
		return nil, err
	}
	return &CedroServiceManagerRewardsInitiatorUpdatedIterator{contract: _CedroServiceManager.contract, event: "RewardsInitiatorUpdated", logs: logs, sub: sub}, nil
}

// WatchRewardsInitiatorUpdated is a free log subscription operation binding the contract event 0xe11cddf1816a43318ca175bbc52cd0185436e9cbead7c83acc54a73e461717e3.
//
// Solidity: event RewardsInitiatorUpdated(address prevRewardsInitiator, address newRewardsInitiator)
func (_CedroServiceManager *CedroServiceManagerFilterer) WatchRewardsInitiatorUpdated(opts *bind.WatchOpts, sink chan<- *CedroServiceManagerRewardsInitiatorUpdated) (event.Subscription, error) {

	logs, sub, err := _CedroServiceManager.contract.WatchLogs(opts, "RewardsInitiatorUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CedroServiceManagerRewardsInitiatorUpdated)
				if err := _CedroServiceManager.contract.UnpackLog(event, "RewardsInitiatorUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRewardsInitiatorUpdated is a log parse operation binding the contract event 0xe11cddf1816a43318ca175bbc52cd0185436e9cbead7c83acc54a73e461717e3.
//
// Solidity: event RewardsInitiatorUpdated(address prevRewardsInitiator, address newRewardsInitiator)
func (_CedroServiceManager *CedroServiceManagerFilterer) ParseRewardsInitiatorUpdated(log types.Log) (*CedroServiceManagerRewardsInitiatorUpdated, error) {
	event := new(CedroServiceManagerRewardsInitiatorUpdated)
	if err := _CedroServiceManager.contract.UnpackLog(event, "RewardsInitiatorUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CedroServiceManagerStaleStakesForbiddenUpdateIterator is returned from FilterStaleStakesForbiddenUpdate and is used to iterate over the raw logs and unpacked data for StaleStakesForbiddenUpdate events raised by the CedroServiceManager contract.
type CedroServiceManagerStaleStakesForbiddenUpdateIterator struct {
	Event *CedroServiceManagerStaleStakesForbiddenUpdate // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CedroServiceManagerStaleStakesForbiddenUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CedroServiceManagerStaleStakesForbiddenUpdate)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CedroServiceManagerStaleStakesForbiddenUpdate)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CedroServiceManagerStaleStakesForbiddenUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CedroServiceManagerStaleStakesForbiddenUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CedroServiceManagerStaleStakesForbiddenUpdate represents a StaleStakesForbiddenUpdate event raised by the CedroServiceManager contract.
type CedroServiceManagerStaleStakesForbiddenUpdate struct {
	Value bool
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterStaleStakesForbiddenUpdate is a free log retrieval operation binding the contract event 0x40e4ed880a29e0f6ddce307457fb75cddf4feef7d3ecb0301bfdf4976a0e2dfc.
//
// Solidity: event StaleStakesForbiddenUpdate(bool value)
func (_CedroServiceManager *CedroServiceManagerFilterer) FilterStaleStakesForbiddenUpdate(opts *bind.FilterOpts) (*CedroServiceManagerStaleStakesForbiddenUpdateIterator, error) {

	logs, sub, err := _CedroServiceManager.contract.FilterLogs(opts, "StaleStakesForbiddenUpdate")
	if err != nil {
		return nil, err
	}
	return &CedroServiceManagerStaleStakesForbiddenUpdateIterator{contract: _CedroServiceManager.contract, event: "StaleStakesForbiddenUpdate", logs: logs, sub: sub}, nil
}

// WatchStaleStakesForbiddenUpdate is a free log subscription operation binding the contract event 0x40e4ed880a29e0f6ddce307457fb75cddf4feef7d3ecb0301bfdf4976a0e2dfc.
//
// Solidity: event StaleStakesForbiddenUpdate(bool value)
func (_CedroServiceManager *CedroServiceManagerFilterer) WatchStaleStakesForbiddenUpdate(opts *bind.WatchOpts, sink chan<- *CedroServiceManagerStaleStakesForbiddenUpdate) (event.Subscription, error) {

	logs, sub, err := _CedroServiceManager.contract.WatchLogs(opts, "StaleStakesForbiddenUpdate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CedroServiceManagerStaleStakesForbiddenUpdate)
				if err := _CedroServiceManager.contract.UnpackLog(event, "StaleStakesForbiddenUpdate", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseStaleStakesForbiddenUpdate is a log parse operation binding the contract event 0x40e4ed880a29e0f6ddce307457fb75cddf4feef7d3ecb0301bfdf4976a0e2dfc.
//
// Solidity: event StaleStakesForbiddenUpdate(bool value)
func (_CedroServiceManager *CedroServiceManagerFilterer) ParseStaleStakesForbiddenUpdate(log types.Log) (*CedroServiceManagerStaleStakesForbiddenUpdate, error) {
	event := new(CedroServiceManagerStaleStakesForbiddenUpdate)
	if err := _CedroServiceManager.contract.UnpackLog(event, "StaleStakesForbiddenUpdate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

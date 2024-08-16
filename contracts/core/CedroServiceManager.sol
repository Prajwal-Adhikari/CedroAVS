// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.24;

import {EnumerableSet} from "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";
import {ServiceManagerBase, IAVSDirectory, IRewardsCoordinator, IRegistryCoordinator, ISignatureUtils, IStakeRegistry} from "./eigenContracts/ServiceManagerBase.sol";
import {BLSSignatureChecker} from "./eigenContracts/BLSSignatureChecker.sol";
import {IServiceManager} from "./eigenContracts/interfaces/IServiceManager.sol";
import {IAgentManager} from "../operator/interfaces/IAgentManager.sol";
import {IRoleManager} from "./interfaces/IRoleManager.sol";
import {CedroServiceManagerStorage} from "./CedroServiceManagerStorage.sol";
import {CedroServiceData} from "./libraries/CedroServiceData.sol";

string constant TAG = "CedroServiceManager";

contract CedroServiceManager is
    CedroServiceManagerStorage,
    ServiceManagerBase,
    BLSSignatureChecker
{
    using EnumerableSet for EnumerableSet.Bytes32Set;

    IAgentManager public agentManager;

    error OperatorRegistrationPaused(string contractTAG);
    error MessageHashAlreadyAdded(string contractTAG);
    error InvalidSender(string ContractTAG);
    error InvalidQuorumParam(string ContractTAG);
    error InvalidReferenceBlockNum(string ContractTAG);
    error InvalidRequestId(string ContractTAG);
    error InsufficientThreshold(string ContractTAG);
    error InvalidAmount(string ContractTAG);
    error InvalidQuorumThresholdPercentage(string ContractTAG);
    error InsufficientThresholdPercentages(string ContractTAG);

    event OperatorRegistered(address indexed operator);
    event OperatorDeRegistered(address indexed operator);
    event QuorumThresholdPercentageChanged(uint256 thresholdpercentage);
    event AgentBalanceConfirmed(address agentAddress, uint256 SlashingAmount);
    event AgentFraudConfirmed(address agentAddress, uint256 SlashingAmount);

    /**
     * @dev Modifier that checks if system is paused or not.
     **/
    modifier Ispaused() {
        require(!pausedStatus, "AVS System is paused");
        _;
    }

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    function initialize(
        address __roleManagerAddress,
        IAVSDirectory __avsDirectory,
        IRewardsCoordinator __rewardsCoordinator,
        IStakeRegistry __stakeRegistry
    ) public initializer {
        _BLSSignatureChecker_init(__roleManagerAddress);
        __ServiceManagerBase_init(
            __avsDirectory,
            __rewardsCoordinator,
            __stakeRegistry,
            __roleManagerAddress
        );
    }

    //////////////////////////////////////////////////////////////////////////////
    //                          Operatoion Verifications                          //
    //////////////////////////////////////////////////////////////////////////////

    /**
    //  * @notice Confirm the AI Agent Blance accross the chains and check if it is deviated
    //  * @param  deviationBalanceHeader -
    //  * @param nonSignerStakesAndSignature -
    //  **/
    function confirmAgentDeviation(
        CedroServiceData.DeviationBalanceHeader calldata deviationBalanceHeader,
        NonSignerStakesAndSignature memory nonSignerStakesAndSignature
    ) external onlyRoles(IRoleManager(roleManagerAddress).VERIFIER_ADMIN()) {
        if (tx.origin != msg.sender) revert InvalidSender(TAG);

        if (deviationBalanceHeader.referenceBlockNumber >= block.number)
            revert InvalidReferenceBlockNum(TAG);

        if (
            deviationBalanceHeader.requestId != type(uint256).max &&
            (deviationBalanceHeader.requestId <=
                agentsDeviationRequestId[deviationBalanceHeader.agentAddress] ||
                deviationBalanceHeader.requestId == 0)
        ) revert InvalidRequestId(TAG);

        if (
            deviationBalanceHeader.quorumNumbers.length !=
            deviationBalanceHeader.quorumThresholdPercentages.length
        ) revert InvalidQuorumParam(TAG);

        // calculate reducedBalancehHeaderHash which nodes signed
        bytes32 reducedBalanceHeaderHash = hashDeviationBalanceHeader(
            deviationBalanceHeader
        );

        (
            QuorumStakeTotals memory quorumStakeTotals,
            bytes32 signatoryRecordHash
        ) = checkSignatures(
                reducedBalanceHeaderHash,
                deviationBalanceHeader.quorumNumbers, // use list of uint8s instead of uint256 bitmap to not iterate 256 times
                deviationBalanceHeader.referenceBlockNumber,
                nonSignerStakesAndSignature
            );

        for (
            uint256 i = 0;
            i < deviationBalanceHeader.quorumThresholdPercentages.length;
            ++i
        ) {
            uint8 currentQuorumThresholdPercentages = uint8(
                deviationBalanceHeader.quorumThresholdPercentages[i]
            );
            if (currentQuorumThresholdPercentages > 100) {
                revert InvalidQuorumThresholdPercentage(TAG);
            }
            if (currentQuorumThresholdPercentages < quorumThresholdPercentage) {
                revert InsufficientThresholdPercentages(TAG);
            }
            if (
                quorumStakeTotals.signedStakeForQuorum[i] *
                    THRESHOLD_DENOMINATOR <
                quorumStakeTotals.totalStakeForQuorum[i] *
                    currentQuorumThresholdPercentages
            ) {
                revert InsufficientThreshold(TAG);
            }
        }

        // store deviation message hash
        bool success = _resolvedMessageHashes.add(reducedBalanceHeaderHash);
        if (!success) {
            revert MessageHashAlreadyAdded(TAG);
        }

        // calcualte the amount to be slashed
        uint256 slashingAmount = calculateSlashingAmount(
            deviationBalanceHeader.idealAmount,
            deviationBalanceHeader.actualAmount,
            agentManager.deviationPercentage(
                deviationBalanceHeader.agentAddress
            )
        );

        agentManager.slashAgentForDeviation(
            deviationBalanceHeader.agentAddress,
            slashingAmount
        );

        //calls to agent manager when check is not initiated offchain, request Id is sent max to denote it is offchain call
        if (deviationBalanceHeader.requestId != type(uint256).max) {
            agentManager.receivedResponse(deviationBalanceHeader.agentAddress);
            agentsDeviationRequestId[
                deviationBalanceHeader.agentAddress
            ] = deviationBalanceHeader.requestId;
        }

        emit AgentBalanceConfirmed(
            deviationBalanceHeader.agentAddress,
            slashingAmount
        );
    }

    /**
     * @notice Challenge Agent performance to process bond slashing
     * @param bankRunBalanceHeader -
     * @param nonSignerStakesAndSignature
     **/
    function confirmAgentBankRun(
        CedroServiceData.BankRunBalanceHeader calldata bankRunBalanceHeader,
        NonSignerStakesAndSignature memory nonSignerStakesAndSignature
    ) external onlyRoles(IRoleManager(roleManagerAddress).VERIFIER_ADMIN()) {
        if (tx.origin != msg.sender) revert InvalidSender(TAG);

        if (bankRunBalanceHeader.referenceBlockNumber >= block.number)
            revert InvalidReferenceBlockNum(TAG);

        if (bankRunBalanceHeader.insufficientAmount <= 0)
            revert InvalidAmount(TAG);

        if (
            bankRunBalanceHeader.quorumNumbers.length !=
            bankRunBalanceHeader.quorumThresholdPercentages.length
        ) revert InvalidQuorumParam(TAG);

        if (
            bankRunBalanceHeader.transactionId <=
            agentsDeviationRequestId[bankRunBalanceHeader.agentAddress] ||
            bankRunBalanceHeader.transactionId == 0
        ) revert InvalidRequestId(TAG);

        // calculate reducedFraudBalancehHeaderHash which nodes signed
        bytes32 reducedFraudBalanceHeaderHash = hashFraudBalanceHeader(
            bankRunBalanceHeader
        );

        // check the signature
        (
            QuorumStakeTotals memory quorumStakeTotals,
            bytes32 signatoryRecordHash
        ) = checkSignatures(
                reducedFraudBalanceHeaderHash,
                bankRunBalanceHeader.quorumNumbers, // use list of uint8s instead of uint256 bitmap to not iterate 256 times
                bankRunBalanceHeader.referenceBlockNumber,
                nonSignerStakesAndSignature
            );

        for (
            uint256 i = 0;
            i < bankRunBalanceHeader.quorumThresholdPercentages.length;
            ++i
        ) {
            uint8 currentQuorumThresholdPercentages = uint8(
                bankRunBalanceHeader.quorumThresholdPercentages[i]
            );
            if (currentQuorumThresholdPercentages > 100) {
                revert InvalidQuorumThresholdPercentage(TAG);
            }
            if (currentQuorumThresholdPercentages < quorumThresholdPercentage) {
                revert InsufficientThresholdPercentages(TAG);
            }
            if (
                quorumStakeTotals.signedStakeForQuorum[i] *
                    THRESHOLD_DENOMINATOR <
                quorumStakeTotals.totalStakeForQuorum[i] *
                    currentQuorumThresholdPercentages
            ) {
                revert InsufficientThreshold(TAG);
            }
        }

        // store deviation message hash
        bool success = _resolvedMessageHashes.add(
            reducedFraudBalanceHeaderHash
        );

        if (!success) {
            revert MessageHashAlreadyAdded(TAG);
        }

        agentManager.slashAgentForBankRun(
            bankRunBalanceHeader.agentAddress,
            bankRunBalanceHeader.insufficientAmount,
            bankRunBalanceHeader.transactionId,
            bankRunBalanceHeader.chainId,
            bankRunBalanceHeader.agentContract
        );

        emit AgentFraudConfirmed(
            bankRunBalanceHeader.agentAddress,
            bankRunBalanceHeader.insufficientAmount
        );
    }

    /**
     * @notice Challenge Operator performance to process slashing
     * @param evidence - evidence of the fraud.
     **/
    function challengeOperators(
        CedroServiceData.OperatorChallengeEvidence calldata evidence
    ) external {
        // :TO DO
    }

    //////////////////////////////////////////////////////////////////////////////
    //                              Contract Functions                             //
    //////////////////////////////////////////////////////////////////////////////

    /**
     * @notice Whitelists operators
     * This function can only be accessed by account with role AVS_ADMIN.
     * @param operators - list of operators to whitelist
     **/

    /// @notice hash the deviation header
    function hashDeviationBalanceHeader(
        CedroServiceData.DeviationBalanceHeader calldata deviationBalanceHeader
    ) internal pure returns (bytes32) {
        return
            keccak256(
                abi.encode(reduceDeviationBalanceHeader(deviationBalanceHeader))
            );
    }

    /// @notice hash the fraud header
    function hashFraudBalanceHeader(
        CedroServiceData.BankRunBalanceHeader calldata deviationBalanceHeader
    ) internal pure returns (bytes32) {
        return
            keccak256(
                abi.encode(reduceFraudBalanceHeader(deviationBalanceHeader))
            );
    }

    function setRegistryCoordinatorAddress(
        IRegistryCoordinator __registryCoordinator
    ) external onlyRoles(IRoleManager(roleManagerAddress).AVS_ADMIN()) {
        registryCoordinator = __registryCoordinator;
        _registryCoordinator = __registryCoordinator;
        stakeRegistry = registryCoordinator.stakeRegistry();
        blsApkRegistry = registryCoordinator.blsApkRegistry();
        delegation = stakeRegistry.delegation();
    }

    function setQuorumThresholdPercentage(
        uint8 thresholdPercentage
    ) external onlyRoles(IRoleManager(roleManagerAddress).AVS_ADMIN()) {
        if (thresholdPercentage > 100) {
            revert InvalidQuorumThresholdPercentage(TAG);
        }
        quorumThresholdPercentage = thresholdPercentage;
        emit QuorumThresholdPercentageChanged(thresholdPercentage);
    }

    function setAgentManagerAddress(
        IAgentManager _agentManagerAddress
    ) external onlyRoles(IRoleManager(roleManagerAddress).AVS_ADMIN()) {
        agentManager = _agentManagerAddress;
    }

    function setPauseStatus(
        bool status
    ) external onlyRole(IRoleManager(roleManagerAddress).AVS_ADMIN()) {
        pausedStatus = status;
    }

    /**
     * @notice converts a balance header to a reduced balance header
     * @param deviationBalanceHeader the balance header to convert
     */

    function reduceDeviationBalanceHeader(
        CedroServiceData.DeviationBalanceHeader calldata deviationBalanceHeader
    ) internal pure returns (CedroServiceData.ReducedDeviationHeader memory) {
        return
            CedroServiceData.ReducedDeviationHeader({
                agentAddress: deviationBalanceHeader.agentAddress,
                referenceBlockNumber: deviationBalanceHeader
                    .referenceBlockNumber,
                requestId: deviationBalanceHeader.requestId,
                actualAmount: deviationBalanceHeader.actualAmount,
                idealAmount: deviationBalanceHeader.idealAmount
            });
    }

    /**
     * @notice Calculates the amount that needs to be slashed of the agent
     * @dev using Basis Points as measurement units - 1% = 100 basis points
     * @param idealAmount the ideal balance that ai agent should have
     * @param actualAmount the actual balance possesed by ai agent
     * @param deviationPercentage the accepted deviation percentage of an ai agent
     */
    function calculateSlashingAmount(
        uint256 idealAmount,
        uint256 actualAmount,
        uint256 deviationPercentage
    ) internal pure returns (uint256) {
        return ((((1e18 - deviationPercentage)) * idealAmount - actualAmount) /
            1e18);
    }

    /// @notice Checks if messageHash exists
    function contains(bytes32 messageHash) external view returns (bool) {
        return _resolvedMessageHashes.contains(messageHash);
    }

    /**
     * @notice converts a balance header to a reduced balance header
     * @param bankRunBalanceHeader the balance header to convert
     */

    function reduceFraudBalanceHeader(
        CedroServiceData.BankRunBalanceHeader calldata bankRunBalanceHeader
    ) internal pure returns (CedroServiceData.ReducedBankRunHeader memory) {
        return
            CedroServiceData.ReducedBankRunHeader({
                agentAddress: bankRunBalanceHeader.agentAddress,
                agentContract: bankRunBalanceHeader.agentContract,
                referenceBlockNumber: bankRunBalanceHeader.referenceBlockNumber,
                chainId: bankRunBalanceHeader.chainId,
                transactionId: bankRunBalanceHeader.transactionId,
                insufficientAmount: bankRunBalanceHeader.insufficientAmount
            });
    }

    /// @notice Given a reference block number, returns the block until which operators must serve.
    function latestServeUntilBlock(
        uint32 referenceBlockNumber
    ) external view returns (uint32) {}
}

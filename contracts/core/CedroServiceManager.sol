// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.24;

import {EnumerableSet} from "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";
import {ReentrancyGuardUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";
import {ECDSA} from "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import {Pausable} from "./eigenContracts/Pausable.sol";
import {IPauserRegistry} from "./eigenContracts/interfaces/IPauserRegistry.sol";
import {ServiceManagerBase, IAVSDirectory, IRewardsCoordinator, IRegistryCoordinator, ISignatureUtils, IStakeRegistry} from "./eigenContracts/ServiceManagerBase.sol";
import {BLSSignatureChecker} from "./eigenContracts/BLSSignatureChecker.sol";
import {IServiceManager} from "./eigenContracts/interfaces/IServiceManager.sol";

import {IRoleManager} from "./interfaces/IRoleManager.sol";
import {CedroServiceManagerStorage} from "./CedroServiceManagerStorage.sol";
import {DATA} from "./libraries/Data.sol";

string constant TAG = "CedroServiceManager";

contract CedroServiceManager is
    CedroServiceManagerStorage,
    ServiceManagerBase,
    BLSSignatureChecker
{
    using EnumerableSet for EnumerableSet.Bytes32Set;

    //status of the operator registration
    bool public OperatorRegistrationStatus;
    //address of the rolemanager contract

    uint256 public roleManagerAddressOne;

    error OperatorRegistrationPaused(string contractTAG);
    error MessageHashAlreadyAdded(string contractTAG);
    error InvalidSender(string ContractTAG);
    error InvalidQuorumParam(string ContractTAG);
    error InvalidReferenceBlockNum(string ContractTAG);
    error InvalidRequestId(string ContractTAG);
    error InvalidQuorumThresholdPercentage(string ContractTAG);
    error InsufficientThresholdPercentages(string ContractTAG);
    error InsufficientThreshold(string ContractTAG);

    event OperatorRegistered(address indexed operator);
    event OperatorDeRegistered(address indexed operator);
    event AgentBalanceConfirmed();
    event AgentFraudConfirmed();

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
        address _roleManager,
        IAVSDirectory __avsDirectory,
        IRewardsCoordinator ___paymentCoordinator,
        IStakeRegistry __stakeRegistry
    ) public initializer {
        _BLSSignatureChecker_init(_roleManager);
        __ServiceManagerBase_init(
            __avsDirectory,
            ___paymentCoordinator,
            __stakeRegistry,
            _roleManager
        );
    }


    //////////////////////////////////////////////////////////////////////////////
    //                          Operatoion Verifications                          //
    //////////////////////////////////////////////////////////////////////////////

    // /**
    //  * @notice Confirm the AI Agent Blance accross the chains
    //  * @param  balanceHeader -
    //  * @param nonSignerStakesAndSignature -
    //  **/
    function confirmAgentDeviation(
        DATA.DeviationBalanceHeader calldata balanceHeader,
        NonSignerStakesAndSignature memory nonSignerStakesAndSignature
    ) external {
        if (tx.origin != msg.sender) revert InvalidSender(TAG);

        if (balanceHeader.referenceBlockNumber >= block.number) {
            revert InvalidReferenceBlockNum(TAG);
        }

        if (
            balanceHeader.requestId <=
            agentDeviationRequestId[balanceHeader.agentOperator]
        ) {
            revert InvalidRequestId(TAG);
        }

        if (
            balanceHeader.quorumNumbers.length !=
            balanceHeader.quorumThresholdPercentages.length
        ) revert InvalidQuorumParam(TAG);

        // calculate reducedBalancehHeaderHash which nodes signed
        bytes32 reducedBalanceHeaderHash = hashDeviationBalanceHeader(
            balanceHeader
        );

        (
            QuorumStakeTotals memory quorumStakeTotals,
            bytes32 signatoryRecordHash
        ) = checkSignatures(
                reducedBalanceHeaderHash,
                balanceHeader.quorumNumbers, // use list of uint8s instead of uint256 bitmap to not iterate 256 times
                balanceHeader.referenceBlockNumber,
                nonSignerStakesAndSignature
            );

        for (
            uint256 i = 0;
            i < balanceHeader.quorumThresholdPercentages.length;
            ++i
        ) {
            uint8 currentQuorumThresholdPercentages = uint8(
                balanceHeader.quorumThresholdPercentages[i]
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

        emit AgentBalanceConfirmed();

        agentDeviationRequestId[balanceHeader.agentOperator] = balanceHeader
            .requestId;
    }

    /**
     * @notice Challenge Agent performance to process bond slashing
     * @param fraudBalanceHeader - evidence of fraud.
     **/
    function confirmAgentFraud(
        DATA.FraudBalanceHeader calldata fraudBalanceHeader,
        NonSignerStakesAndSignature memory nonSignerStakesAndSignature
    ) external {
        if (tx.origin != msg.sender) revert InvalidSender(TAG);

        if (fraudBalanceHeader.referenceBlockNumber >= block.number) {
            revert InvalidReferenceBlockNum(TAG);
        }

        if (
            fraudBalanceHeader.requestId <=
            agentFraudRequestId[fraudBalanceHeader.agentOperator]
        ) {
            revert InvalidRequestId(TAG);
        }

        if (
            fraudBalanceHeader.quorumNumbers.length ==
            fraudBalanceHeader.quorumThresholdPercentages.length
        ) revert InvalidQuorumParam(TAG);

        // calculate reducedBalancehHeaderHash which nodes signed
        bytes32 reducedBalanceHeaderHash = hashFraudBalanceHeader(
            fraudBalanceHeader
        );

        // check the signature
        (
            QuorumStakeTotals memory quorumStakeTotals,
            bytes32 signatoryRecordHash
        ) = checkSignatures(
                reducedBalanceHeaderHash,
                fraudBalanceHeader.quorumNumbers, // use list of uint8s instead of uint256 bitmap to not iterate 256 times
                fraudBalanceHeader.referenceBlockNumber,
                nonSignerStakesAndSignature
            );

        for (
            uint256 i = 0;
            i < fraudBalanceHeader.quorumThresholdPercentages.length;
            ++i
        ) {
            uint8 currentQuorumThresholdPercentages = uint8(
                fraudBalanceHeader.quorumThresholdPercentages[i]
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

        uint32 requestIdMemory = agentDeviationRequestId[
            fraudBalanceHeader.agentOperator
        ];

        // store deviation message hash
        bool success = _resolvedMessageHashes.add(reducedBalanceHeaderHash);
        if (!success) {
            revert MessageHashAlreadyAdded(TAG);
        }

        agentDeviationRequestId[fraudBalanceHeader.agentOperator] =
            requestIdMemory +
            1;

        emit AgentFraudConfirmed();
    }

    /**
     * @notice Challenge Operator performance to process slashing
     * @param evidence - evidence of the fraud.
     **/
    function challengeOperators(
        DATA.OperatorChallengeEvidence calldata evidence
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
        DATA.DeviationBalanceHeader calldata balanceHeader
    ) internal pure returns (bytes32) {
        return
            keccak256(abi.encode(reduceDeviationBalanceHeader(balanceHeader)));
    }

    /// @notice hash the fraud header
    function hashFraudBalanceHeader(
        DATA.FraudBalanceHeader calldata balanceHeader
    ) internal pure returns (bytes32) {
        return keccak256(abi.encode(reduceFraudBalanceHeader(balanceHeader)));
    }

    function setPauseStatus(
        bool status
    ) external onlyRole(IRoleManager(roleManager).AVS_ADMIN()) {
        pausedStatus = status;
    }

    /**
     * @notice converts a balance header to a reduced balance header
     * @param balanceHeader the balance header to convert
     */

    function reduceDeviationBalanceHeader(
        DATA.DeviationBalanceHeader calldata balanceHeader
    ) internal pure returns (DATA.ReducedDeviationHeader memory) {
        return
            DATA.ReducedDeviationHeader({
                actualAmount: balanceHeader.actualAmount,
                idealAmount: balanceHeader.idealAmount,
                referenceRequestId: balanceHeader.requestId,
                referenceBlockNumber: balanceHeader.referenceBlockNumber
            });
    }

    /// @notice Checks if messageHash exists
    function contains(bytes32 messageHash) external view returns (bool) {
        return _resolvedMessageHashes.contains(messageHash);
    }

    /**
     * @notice converts a balance header to a reduced balance header
     * @param balanceHeader the balance header to convert
     */

    function reduceFraudBalanceHeader(
        DATA.FraudBalanceHeader calldata balanceHeader
    ) internal pure returns (DATA.ReducedFraudHeader memory) {
        return
            DATA.ReducedFraudHeader({
                insufficientAmount: balanceHeader.insufficientAmount,
                referenceBlockNumber: balanceHeader.referenceBlockNumber
            });
    }

    /// @notice Given a reference block number, returns the block until which operators must serve.
    function latestServeUntilBlock(
        uint32 referenceBlockNumber
    ) external view returns (uint32) {}
}

// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.24;
import {BN254} from "../eigenContracts/libraries/BN254.sol";

library CedroServiceData {
    struct DeviationBalanceHeader {
        address agentAddress;
        uint32 referenceBlockNumber;
        bytes quorumNumbers;
        bytes quorumThresholdPercentages;
        uint256 requestId;
        uint256 actualAmount;
        uint256 idealAmount;
    }

    struct FraudBalanceHeader {
        address agentAddress;
        uint32 referenceBlockNumber;
        bytes quorumNumbers;
        bytes quorumThresholdPercentages;
        uint256 transactionId;
        uint256 insufficientAmount;
    }

    struct ReducedDeviationHeader {
        address agentAddress;
        uint32 referenceBlockNumber;
        uint256 requestId;
        uint256 actualAmount;
        uint256 idealAmount;
       
    }

    struct ReducedFraudHeader {
        address agentAddress;
        uint32 referenceBlockNumber;
        uint256 transactionId;
        uint256 insufficientAmount;
    }

    struct OperatorChallengeEvidence {
        bytes32 messageHash;
        // :TO DO
    }
}
// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.24;
import {BN254} from "../eigenContracts/libraries/BN254.sol";

library DATA {
    struct DeviationBalanceHeader {
        address agentOperator;
        bytes quorumNumbers;
        bytes quorumThresholdPercentages;
        uint32 referenceBlockNumber;
        uint32 requestId;
        uint256 actualAmount;
        uint256 idealAmount;
    }
    struct FraudBalanceHeader {
        bytes quorumNumbers;
        bytes quorumThresholdPercentages;
        uint32 requestId;
        uint32 referenceBlockNumber;
        uint256 insufficientAmount;
    }

    struct ReducedDeviationHeader {
        uint256 actualAmount;
        uint256 idealAmount;
        uint32 referenceRequestId;
        uint32 referenceBlockNumber;
    }
    struct ReducedFraudHeader {
        uint256 insufficientAmount;
        uint32 referenceBlockNumber;
    }

    struct AgentChallengeEvidence {
        bytes32 messageHash;
    }

}

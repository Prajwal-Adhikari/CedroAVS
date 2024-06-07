// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.24;
import {BN254} from "../eigenContracts/libraries/BN254.sol";

library DATA {
    struct DeviationBalanceHeader {
        address agentOperator;
        bytes quorumNumbers;
        bytes quorumThresholdPercentages;
        uint32 referenceBlockNumber;
        uint256 actualAmount;
        uint256 idealAmount;
    }
    struct FraudBalanceHeader {
        bytes quorumNumbers;
        bytes quorumThresholdPercentages;
        uint32 referenceBlockNumber;
        uint256 insufficientAmount;
    }

    struct ReducedDeviationHeader {
        uint256 actualAmount;
        uint256 idealAmount;
        uint32 referenceBlockNumber;
    }
    struct ReducedFraudHeader {
        uint256 insufficientAmount;
        uint32 referenceBlockNumber;
    }

    struct AgentChallengeEvidence {
        bytes32 messageHash;
    }

    struct NonSignerStakesAndSignature {
        uint32[] nonSignerQuorumBitmapIndices; // is the indices of all nonsigner quorum bitmaps
        BN254.G1Point[] nonSignerPubkeys; // is the G1 pubkeys of all nonsigners
        BN254.G1Point[] quorumApks; // is the aggregate G1 pubkey of each quorum
        BN254.G2Point apkG2; // is the aggregate G2 pubkey of all signers
        BN254.G1Point sigma; // is the aggregate G1 signature of all signers
        uint32[] quorumApkIndices; // is the indices of each quorum aggregate pubkey
        uint32[] totalStakeIndices; // is the indices of each quorums total stake
        uint32[][] nonSignerStakeIndices; // is the indices of each non signers stake within a quorum
    }
    struct NonSignerStakesAndSignatureDemo {
        uint32 quorumApkIndices; // is the indices of each quorum aggregate pubkey
        uint32 totalStakeIndices; // is the indices of each quorums total stake
    }

    struct AgentData {
        uint256 Balance;
        uint256 lastCheckedTime;
        uint256 averageInterestRate;
    }
}

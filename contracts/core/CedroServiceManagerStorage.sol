// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.24;

import {EnumerableSet} from "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";


contract CedroServiceManagerStorage{

    // CONSTANTS
    uint256 public constant THRESHOLD_DENOMINATOR = 100;

    /// @notice Minimal quorum threshold percentage
    uint8 public quorumThresholdPercentage;

    /// @notice The current batchId
    uint32 public requestId;

    /// @notice Resolved deviation message hashes, prevent aggregator from replay any resolved message hash
    EnumerableSet.Bytes32Set internal _resolvedMessageHashes;

    // storage gap for upgradeability
    // slither-disable-next-line shadowing-state
    uint256[47] private __GAP;
}

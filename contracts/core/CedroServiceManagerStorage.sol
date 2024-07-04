// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.24;

import {EnumerableSet} from "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";

contract CedroServiceManagerStorage {
    // CONSTANTS
    uint256 public constant THRESHOLD_DENOMINATOR = 100;

    /// @notice Minimal quorum threshold percentage
    uint8 public quorumThresholdPercentage;

    ///@notice Pause status of the system
    bool pausedStatus;

    /// @notice The current requestid of the  agent deviation check, specific to agent operator
    mapping(address => uint256) public agentsDeviationRequestId;

    /// @notice The current requestid of the  agent fraud check, specific to agent operator
    mapping(address => uint256) public agentsFraudRequestId;

    /// @notice Resolved deviation message hashes, prevent aggregator from replay any resolved message hash
    EnumerableSet.Bytes32Set internal _resolvedMessageHashes;

    // storage gap for upgradeability
    // slither-disable-next-line shadowing-state
    uint256[47] private __GAP;
}

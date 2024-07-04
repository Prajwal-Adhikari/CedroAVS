// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.12;

import {ContextUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/ContextUpgradeable.sol";

import {IDelegationManager} from "./interfaces/IDelegationManager.sol";
import {IStrategyManager, IStrategy} from "./interfaces/IStrategyManager.sol";
import {IRegistryCoordinator} from "./interfaces/IRegistryCoordinator.sol";
import {IRoleManager} from "../interfaces/IRoleManager.sol";
import {IStakeRegistry} from "./interfaces/IStakeRegistry.sol";

/**
 * @title Storage variables for the `StakeRegistry` contract.
 * @author Layr Labs, Inc.
 * @notice This storage contract is separate from the logic to simplify the upgrade process.
 */
abstract contract StakeRegistryStorage is IStakeRegistry, ContextUpgradeable {
    /// @notice Constant used as a divisor in calculating weights.
    uint256 public constant WEIGHTING_DIVISOR = 1e18;
    /// @notice Maximum length of dynamic arrays in the `strategyParams` mapping.
    uint8 public constant MAX_WEIGHING_FUNCTION_LENGTH = 32;
    /// @notice Constant used as a divisor in dealing with BIPS amounts.
    uint256 internal constant MAX_BIPS = 10000;

    /// @notice The address of the Delegation contract for EigenLayer.
    IDelegationManager public delegation;

    /// @notice the coordinator contract that this registry is associated with
    address public registryCoordinator;

    //address of the rolemanager contract
    address public roleManagerAddress;

    /// @notice In order to register for a quorum i, an operator must have at least `minimumStakeForQuorum[i]`
    /// evaluated by this contract's 'VoteWeigher' logic.
    mapping(uint8 => uint96) public minimumStakeForQuorum;

    /// @notice History of the total stakes for each quorum
    mapping(uint8 => StakeUpdate[]) internal _totalStakeHistory;

    /// @notice mapping from operator's operatorId to the history of their stake updates
    mapping(bytes32 => mapping(uint8 => StakeUpdate[]))
        internal operatorStakeHistory;

    /**
     * @notice mapping from quorum number to the list of strategies considered and their
     * corresponding multipliers for that specific quorum
     */
    mapping(uint8 => StrategyParams[]) public strategyParams;
    mapping(uint8 => IStrategy[]) public strategiesPerQuorum;

    /**
     * @dev Modifier that checks if caller has a specific role.
     **/
    modifier onlyRole(bytes32 role) {
        IRoleManager(roleManagerAddress).checkRole(role, _msgSender());
        _;
    }

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        // disable initializers so that the implementation contract cannot be initialized
        _disableInitializers();
    }

    function _stakeRegistryStorage_init(
        address _roleManagerAddress,
        IDelegationManager _delegationManager
    ) internal initializer {
        delegation = _delegationManager;
        roleManagerAddress = _roleManagerAddress;
    }

    function setRegisteryCoordinator(
        IRegistryCoordinator _registryCoordinator
    ) external onlyRole(IRoleManager(roleManagerAddress).AVS_ADMIN()) {
        registryCoordinator = address(_registryCoordinator);
    }

    // storage gap for upgradeability
    // slither-disable-next-line shadowing-state
    uint256[45] private __GAP;
}

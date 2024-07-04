// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.12;

import {ContextUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/ContextUpgradeable.sol";

import {IRegistryCoordinator} from "./interfaces/IRegistryCoordinator.sol";
import {IRoleManager} from "../interfaces/IRoleManager.sol";
import {IIndexRegistry} from "./interfaces/IIndexRegistry.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";

/**
 * @title Storage variables for the `IndexRegistry` contract.
 * @author Layr Labs, Inc.
 * @notice This storage contract is separate from the logic to simplify the upgrade process.
 */
abstract contract IndexRegistryStorage is ContextUpgradeable, IIndexRegistry {
    /// @notice The value that is returned when an operator does not exist at an index at a certain block
    bytes32 public constant OPERATOR_DOES_NOT_EXIST_ID = bytes32(0);

    /// @notice The RegistryCoordinator contract for this middleware
    address public registryCoordinator;

    //@notice roleManager Contract
    address public roleManagerAddress;

    /// @notice maps quorumNumber => operator id => current operatorIndex
    /// NOTE: This mapping is NOT updated when an operator is deregistered,
    /// so it's possible that an index retrieved from this mapping is inaccurate.
    /// If you're querying for an operator that might be deregistered, ALWAYS
    /// check this index against the latest `_operatorIndexHistory` entry
    mapping(uint8 => mapping(bytes32 => uint32)) public currentOperatorIndex;
    /// @notice maps quorumNumber => operatorIndex => historical operator ids at that index
    mapping(uint8 => mapping(uint32 => OperatorUpdate[]))
        internal _operatorIndexHistory;
    /// @notice maps quorumNumber => historical number of unique registered operators
    mapping(uint8 => QuorumUpdate[]) internal _operatorCountHistory;

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

    function _indexRegistryStorage_init(
        address _roleManagerAddress
    ) internal initializer {
        roleManagerAddress = _roleManagerAddress;
    }

    function setRegisteryCoordinator(
        IRegistryCoordinator _registryCoordinator
    ) external onlyRole(IRoleManager(roleManagerAddress).AVS_ADMIN()) {
        registryCoordinator = address(_registryCoordinator);
    }

    /**
     * @notice Configures role manager address
     * This function can only be accessed by the account with role DEFAULT_ADMIN_ROLE.
     * @param _roleManagerAddress - address of the roleManager contract.
     **/
    function setRoleManager(
        address _roleManagerAddress
    ) external onlyRole(IRoleManager(roleManagerAddress).AVS_ADMIN()) {
        roleManagerAddress = _roleManagerAddress;
    }
 
    // storage gap for upgradeability
    uint256[47] private __GAP;
}

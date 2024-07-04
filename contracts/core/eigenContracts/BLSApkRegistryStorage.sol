// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.12;

import {ContextUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/ContextUpgradeable.sol";

import {IBLSApkRegistry} from "./interfaces/IBLSApkRegistry.sol";
import {IRegistryCoordinator} from "./interfaces/IRegistryCoordinator.sol";
import {IRoleManager} from "../interfaces/IRoleManager.sol";
import {BN254} from "./libraries/BN254.sol";

abstract contract BLSApkRegistryStorage is IBLSApkRegistry, ContextUpgradeable {
    /// @notice the hash of the zero pubkey aka BN254.G1Point(0,0)
    bytes32 internal constant ZERO_PK_HASH =
        hex"ad3228b676f7d3cd4284a5443f17f1962b36e491b30a40b2405849e597ba5fb5";

    /// @notice the registry coordinator contract
    address public registryCoordinator;

    //@notice roleManager Contract
    address public roleManagerAddress;

    // storage for individual pubkeys
    /// @notice maps operator address to pubkey hash
    mapping(address => bytes32) public operatorToPubkeyHash;
    /// @notice maps pubkey hash to operator address
    mapping(bytes32 => address) public pubkeyHashToOperator;
    /// @notice maps operator address to pubkeyG1
    mapping(address => BN254.G1Point) public operatorToPubkey;

    // storage for aggregate pubkeys (APKs)
    /// @notice maps quorumNumber => historical aggregate pubkey updates
    mapping(uint8 => ApkUpdate[]) public apkHistory;
    /// @notice maps quorumNumber => current aggregate pubkey of quorum
    mapping(uint8 => BN254.G1Point) public currentApk;

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

    function _BLSApkRegistryStorage_init(
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
    uint256[45] private __GAP;
}

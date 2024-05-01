// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.24;

import {ContextUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/ContextUpgradeable.sol";
import {ReentrancyGuardUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";
import {ECDSA} from "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";

import {Pausable} from "./eigenContracts/Pausable.sol";
import {ServiceManagerBase} from "./eigenContracts/ServiceManagerBase.sol";
import {IAVSDirectory} from "./eigenContracts/interfaces/IAVSDirectory.sol";
import {IStakeRegistry} from "./eigenContracts/interfaces/IStakeRegistry.sol";
import {ISignatureUtils} from "./eigenContracts/interfaces/ISignatureUtils.sol";
import {IServiceManager} from "./eigenContracts/interfaces/IServiceManager.sol";
import {IPauserRegistry} from "./eigenContracts/interfaces/IPauserRegistry.sol";
import {IPaymentCoordinator} from "./eigenContracts/interfaces/IPaymentCoordinator.sol";
import {IRegistryCoordinator} from "./eigenContracts/interfaces/IRegistryCoordinator.sol";
import {ECDSA} from "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";

import {IRoleManager} from "./interfaces/IRoleManager.sol";
import {CedroServiceManagerStorage} from "./CedroServiceManagerStorage.sol";
import {IRoleManager} from "./interfaces/IRoleManager.sol";

string constant TAG = "CedroServiceManager";

contract CedroServiceManager is
    ContextUpgradeable,
    CedroServiceManagerStorage,
    ServiceManagerBase
{
    //status of the operator registration
    bool OperatorRegistrationStatus;
    //address of the rolemanager contract
    address public roleManager;

    mapping(address => bool) public operatorAllowlist;

    error OperatorRegistrationPaused(string contractTAG);
    error NotInAllowlist(string contractTAG);

    event OperatorRegistered(address indexed operator);
    event OperatorDeRegistered(address indexed operator);

    /**
     * @dev Modifier that checks if operator is whitelisted
     **/
    modifier onlyWhitelistedOperator() {
        require(operatorAllowlist[msg.sender], "Operator is not whitelisted");
        _;
    }

    /**
     * @dev Modifier that checks if caller has a specific role.
     **/
    modifier onlyRole(bytes32 role) {
        IRoleManager(roleManager).checkRole(role, _msgSender());
        _;
    }

    constructor(
        IAVSDirectory __avsDirectory,
        IPaymentCoordinator ___paymentCoordinator,
        IRegistryCoordinator __registryCoordinator,
        IStakeRegistry __stakeRegistry
    )
        ServiceManagerBase(
            __avsDirectory,
            ___paymentCoordinator,
            __registryCoordinator,
            __stakeRegistry
        )
    {
        _disableInitializers();
    }

    /**
     * @notice Whitelists operators
     * This function can only be accessed by account with role AVS_ADMIN.
     * @param operators - list of operators to whitelist
     **/
    function addOperatorsToAllowlist(
        address[] calldata operators
    ) external onlyRole(IRoleManager(roleManager).AVS_ADMIN()) {
        for (uint256 i; i < operators.length; i++) {
            operatorAllowlist[operators[i]] = true;
        }
    }

    /**
     * @notice Removes operators from whitelist
     * This function can only be accessed by account with role AVS_ADMIN.
     * @param operators - list of operators to remove from whitelist
     **/
    function removeOperatorsFromAllowlist(
        address[] calldata operators
    ) external onlyRole(IRoleManager(roleManager).AVS_ADMIN()) {
        for (uint256 i; i < operators.length; i++) {
            delete operatorAllowlist[operators[i]];
        }
    }

     /**
     * @notice Configures role manager address
     * This function can only be accessed by the account with role DEFAULT_ADMIN_ROLE.
     * @param roleManagerAddress - address of the roleManager contract.
     **/
    function setRoleManager(
        address roleManagerAddress
    ) external onlyRole(IRoleManager(roleManager).AVS_ADMIN()) {
        roleManager = roleManagerAddress;
    }


    //////////////////////////////////////////////////////////////////////////////
    //                          Operator Registration                           //
    //////////////////////////////////////////////////////////////////////////////

    /**
     * @notice Register an operator with the AVS. Forwards call to EigenLayer' AVSDirectory.
     * @param operatorSignature The signature, salt, and expiry of the operator's signature.
     */
    function registerOperatorToAVS(
        address operator,
        ISignatureUtils.SignatureWithSaltAndExpiry memory operatorSignature
    ) public override(ServiceManagerBase) onlyRole(IRoleManager(roleManager).REGISTRYCOORDINATOR_ADMIN()) {
         if (! operatorAllowlist[operator]) {
            revert NotInAllowlist(TAG);
        }
        _avsDirectory.registerOperatorToAVS(operator, operatorSignature);
        emit OperatorRegistered(operator);
    }

    /**
     * @notice Deregister an operator from the AVS. Forwards a call to EigenLayer's AVSDirectory.
     */
    function deregisterOperatorFromAVS(address operator) public override(ServiceManagerBase) onlyRole(IRoleManager(roleManager).REGISTRYCOORDINATOR_ADMIN()) {
        _avsDirectory.deregisterOperatorFromAVS(operator);
        emit OperatorDeRegistered(operator);
    }

    
}

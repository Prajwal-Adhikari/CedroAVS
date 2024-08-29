// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.24;

import "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/utils/ContextUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {IRoleManager} from "./interfaces/IRoleManager.sol";
import "hardhat/console.sol";

string constant TAG = "Pool";

contract Pool is Initializable, ContextUpgradeable, ReentrancyGuardUpgradeable {
    using SafeERC20 for IERC20;

    address public owner;
    address public agentPod;
    //prz:
     address public roleManager;

    error IncorrectRequestId(
        string contractTAG,
        uint256 sentId,
        uint256 requiredId
    );
    error InsufficientValue(
        string contractTAG,
        uint256 valueSent,
        uint256 valueToBeGreater
    );
    error ArgumentsMismatch(string contractTAG);
    error AssetAlreadyExists(bytes32 asset);
    error AssetNotApprovedByAgentPod(bytes32 asset);

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    function initialize(address _owner, address _agentPod) public initializer {
        __ReentrancyGuard_init();
        owner = _owner;
        agentPod = _agentPod;
    }
        /**
     * @dev Modifier that checks if caller has a specific role.
     **/

    modifier onlyRole(bytes32 role) {
        IRoleManager(roleManager).checkRole(role, _msgSender());
        _;
    }

    /**
     * @notice Configures role manager address
     * This function can only be accessed by the account with role DEFAULT_ADMIN_ROLE.
     * @param roleManagerAddress - address of the roleManager contract.
     **/
    // function setRoleManager(
    //     address roleManagerAddress
    // ) external onlyRole(IRoleManager(roleManager).DEFAULT_ADMIN_ROLE()) {
    //     roleManager = roleManagerAddress;
    // }
function setRoleManager(
    address roleManagerAddress
) external {
    roleManager = roleManagerAddress;
}

    /**
     * @notice Configures corresponding agentPod
     * This function can be accessed by the role MANAGER
     * @param _agentPod - address of the AgentPod
     **/

     //prz: Added role AVS_ADMIN as MANAGER 
    function setAgentPod(address _agentPod) public onlyRole(IRoleManager(roleManager).AVS_ADMIN()) {
        agentPod = _agentPod;
    }

    // To be called only by agent pod
    function trasferAssetToDepositor(
        address _user,
        address[] memory _token,
        uint256[] memory _amount
    ) public {
        if (_token.length != _amount.length) revert ArgumentsMismatch(TAG);

        for (uint256 i; i < _token.length; ) {
            IERC20(_token[i]).safeTransfer(_user, _amount[i]);
            unchecked {
                i++;
            }
        }
    }

//prz:
    // To be called only by owner
    function trasferAssetToAgentPod(
        address[] memory _token,
        uint256[] memory _amount
    ) public onlyRole(IRoleManager(roleManager).DEFAULT_ADMIN_ROLE()){
        if (_token.length != _amount.length) revert ArgumentsMismatch(TAG);

        for (uint256 i; i < _token.length; ) {
            IERC20(_token[i]).safeTransfer(agentPod, _amount[i]);
            unchecked {
                i++;
            }
        }
    }
}

// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.24;

import "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/utils/ContextUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";

string constant TAG = "Pool";

contract Pool is Initializable, ContextUpgradeable, ReentrancyGuardUpgradeable {
    using SafeERC20 for IERC20;

    address public owner;
    address public agentPod;

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
     * @notice Configures corresponding agentPod
     * This function can be accessed by the role MANAGER
     * @param _agentPod - address of the AgentPod
     **/
    function setAgentPod(address _agentPod) public {
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

    // To be called only by owner
    function trasferAssetToAgentPod(
        address[] memory _token,
        uint256[] memory _amount
    ) public {
        if (_token.length != _amount.length) revert ArgumentsMismatch(TAG);

        for (uint256 i; i < _token.length; ) {
            IERC20(_token[i]).safeTransfer(agentPod, _amount[i]);
            unchecked {
                i++;
            }
        }
    }
}

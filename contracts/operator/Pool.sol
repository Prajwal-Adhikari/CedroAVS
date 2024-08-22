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
<<<<<<< HEAD

    // To be called for LIDO claim action
    function claimWithdrawals(address _withdrawContractAddress) internal {
        ILidoWithdrawalQueue withdrawContract = ILidoWithdrawalQueue(
            _withdrawContractAddress
        );

        uint256[] memory requestIds = withdrawContract.getWithdrawalRequests(
            address(this)
        );

        require(requestIds.length > 0, "Withdrawal requests ids not found");

        WithdrawalRequestStatus[] memory reqStatusArray = withdrawContract
            .getWithdrawalStatus(requestIds);

        require(requestIds.length > 0, "Withdrawal requests status not found");

        // Count the number of matured requests
        uint count = 0;
        for (uint i = 0; i < reqStatusArray.length; i++) {
            if (reqStatusArray[i].isFinalized && !reqStatusArray[i].isClaimed) {
                count++;
            }
        }

        // Revert if there are no matured withdrawals
        require(count > 0, "No claimable withdrawals yet");

        // store matured requests in an array
        uint256[] memory maturedRequestIds = new uint256[](count);

        uint index = 0;
        for (uint i = 0; i < reqStatusArray.length; i++) {
            if (reqStatusArray[i].isFinalized && !reqStatusArray[i].isClaimed) {
                maturedRequestIds[index] = requestIds[i];
                index++;
            }
        }

        require(
            requestIds.length > 0,
            "Matured Withdrawal requests status not found"
        );

        for (uint i = 0; i < maturedRequestIds.length; i++) {
            withdrawContract.claimWithdrawal(maturedRequestIds[i]);
        }
    }

    function stringToBytes32(
        string memory source
    ) public pure returns (bytes32 result) {
        bytes memory tempEmptyStringTest = bytes(source);
        if (tempEmptyStringTest.length == 0) {
            return 0x0;
        }

        assembly {
            result := mload(add(source, 32))
        }
    }
=======
>>>>>>> 235aad504e9f6b12c999cdc39415aba9c5cd7c2b
}

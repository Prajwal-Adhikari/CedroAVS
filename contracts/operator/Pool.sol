// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.24;

import "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/utils/ContextUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "../actions/interfaces/IAave.sol";
import "../actions/lido/interfaces/ILido.sol";
import "../actions/lido/interfaces/ILidoWithdrawalQueue.sol";

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

    /**
     * @notice Allows agent to perform DeFi action that requires direct interaction as msg.sender
     * This function can be accessed by the owner
     * @param _payload - fixed bytes pattern to be provided for a DeFi action
     * @param _protocol - address of the DeFi protocol
     **/
    function takeAction(bytes memory _payload, address _protocol) public {
        bytes32 id;
        uint8 methodName;
        address asset;
        uint256 amount;
        assembly {
            id := mload(add(_payload, 34))
            methodName := mload(add(_payload, 35))
            asset := mload(add(_payload, 107))
            amount := mload(add(_payload, 139))
        }

        if (id == stringToBytes32("AAVE")) {
            if (methodName == 1) {
                uint256 interestRateMode;
                uint16 referralCode;
                assembly {
                    interestRateMode := mload(add(_payload, 171))
                    referralCode := mload(add(_payload, 173))
                }
                IAave(_protocol).borrow(
                    asset,
                    amount,
                    interestRateMode,
                    referralCode,
                    _msgSender()
                );
            } else if (methodName == 2) {
                uint256 interestRateMode;
                assembly {
                    interestRateMode := mload(add(_payload, 171))
                }
                IAave(_protocol).repayWithATokens(
                    asset,
                    amount,
                    interestRateMode
                );
            }
        }

        if (id == stringToBytes32("LIDO")) {
            // 1: Stake ETH
            if (methodName == 1) {
                // This method stakes our eth and returns equivalent amount of stEth.
                // Referral address is set to Zero address.
                ILido(_protocol).submit{value: amount}(address(0));
            }
            // 2: Request for withdrawal
            else if (methodName == 2) {
                uint256[] memory _amount = new uint256[](1);
                _amount[0] = amount;

                // TODO: How to get the address of Lido contract here ? Currently using zero address as placeholder.
                ILido(address(0)).approve(_protocol, amount);
                ILidoWithdrawalQueue(_protocol).requestWithdrawals(
                    _amount,
                    address(this)
                );
            }
            //  3: Claim withdrawals
            else if (methodName == 3) {
                // If there is no any matured withdrawal request, the transaction will revert. 
                claimWithdrawals(_protocol);
            }
        }
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
}

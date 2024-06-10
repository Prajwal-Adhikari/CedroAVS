// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.24;

import "@openzeppelin/contracts/utils/math/SafeMath.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "@openzeppelin/contracts/utils/Context.sol";
import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "./interfaces/IOracle.sol";
import "./interfaces/IAgentManager.sol";
import "./interfaces/IAgentPod.sol";
import "./libraries/Data.sol";
import "../actions/interfaces/IAave.sol";

string constant TAG = "Pool";

contract Pool is Context, ReentrancyGuard {
    using ECDSA for bytes32;
    using SafeMath for uint256;
    using SafeERC20 for IERC20;

    address public owner;
    address public agentPod;
    address public oppurtunity;

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

    constructor(address _owner, address _agentPod) {
        owner = _owner;
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

    function stringToBytes32(string memory source)
        public
        pure
        returns (bytes32 result)
    {
        bytes memory tempEmptyStringTest = bytes(source);
        if (tempEmptyStringTest.length == 0) {
            return 0x0;
        }

        assembly {
            result := mload(add(source, 32))
        }
    }
}

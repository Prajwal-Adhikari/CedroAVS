// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.24;

import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/utils/ContextUpgradeable.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "./interfaces/IAave.sol";

string constant TAG = "AaveInterface";

/**
 * @title AaveInterface contract
 * @notice Used to iterface call for Aave methods
 **/
contract AaveInterface is Initializable, ContextUpgradeable {
    mapping(address => address) public rToken;
    // Method Types
    uint8 internal constant SUPPLY = 1;
    uint8 internal constant WITHDRAW = 2;
    uint8 internal constant REPAY = 3;

    error UnsupportedFunctionType(string contractTAG, uint8 functionType);

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    /**
     * @notice Configures respective asset to aTokens
     * This function can only be accessed by the account with role MANAGER.
     * @param _asset - address of the token
     * @param _aToken -  aToken address of the asset
     **/
    function setRtoken(address _asset, address _aToken) public {
        rToken[_aToken] = _asset;
    }

    /**
     * @notice Initilize the contract.
     **/
    function initialize() public initializer {}

    /**
     * @notice Performs external call to respective Aave methods
     * This function can only be accessed by the account with role MANAGER.
     * @param _payload - payload containing params information
     * @param _user -  address onBehalf of which txn is done (agent address)
     * @param _protocol - address of the protocol (Aave address)
     **/
    function extractPayload(
        bytes memory _payload,
        address _user,
        address _protocol
    ) public {
        uint8 methodName;
        assembly {
            methodName := mload(add(_payload, 35))
        }
        if (methodName == SUPPLY) {
            address asset;
            uint256 amount;
            uint16 referralCode;

            assembly {
                asset := mload(add(_payload, 55))
                amount := mload(add(_payload, 87))
                referralCode := mload(add(_payload, 89))
            }

            IERC20(asset).approve(_protocol, amount);

            IAave(_protocol).supply(asset, amount, _user, referralCode);
        } else if (methodName == WITHDRAW) {
            address aToken;
            uint256 amount;

            assembly {
                aToken := mload(add(_payload, 55))
                amount := mload(add(_payload, 87))
            }
            IERC20(aToken).approve(_protocol, amount);
            IAave(_protocol).withdraw(rToken[aToken], amount, _user);
        } else if (methodName == REPAY) {
            address asset;
            uint256 amount;
            uint256 interestRateMode;

            assembly {
                asset := mload(add(_payload, 55))
                amount := mload(add(_payload, 87))
                interestRateMode := mload(add(_payload, 119))
            }
            IERC20(asset).approve(_protocol, amount);
            IAave(_protocol).repay(asset, amount, interestRateMode, _user);
        } else {
            revert UnsupportedFunctionType(TAG, methodName);
        }
    }
}

// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.24;

import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/utils/ContextUpgradeable.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "./interfaces/IUniswap.sol";

string constant TAG = "UniswapInterface";

/**
 * @title UniswapInterface contract
 * @notice Used to iterface call for UniswapInterface methods
 **/
contract UniswapInterface is Initializable, ContextUpgradeable {
    mapping(address => address) public rToken;
    // Method Types
    uint8 internal constant SWAP = 1;
  
    error UnsupportedFunctionType(string contractTAG, uint8 functionType);

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
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
        if (methodName == SWAP) {
            address assetIn;
            address assetOut;
            uint256 amountIn;
            uint256 amountOutMin;
            uint256 deadline;

            address[] memory path = new address[](2);

            assembly {
                assetIn := mload(add(_payload, 55))
                assetOut := mload(add(_payload, 107))
                amountIn := mload(add(_payload, 87))
                amountOutMin := mload(add(_payload, 139))
                deadline := mload(add(_payload, 171))
            }
            path[0] = assetIn;
            path[1] = assetOut;

            IERC20(assetIn).approve(_protocol, amountIn);

            IUniswap(_protocol).swapExactTokensForTokens(amountIn, amountOutMin, path, _user, deadline);
        } else {
            revert UnsupportedFunctionType(TAG, methodName);
        }
    }
}

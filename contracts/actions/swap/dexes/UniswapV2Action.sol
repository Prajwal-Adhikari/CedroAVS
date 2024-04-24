// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.24;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/utils/ContextUpgradeable.sol";

import "../../interfaces/IRoleManager.sol";
import "../interfaces/IUniSwapV2Router.sol";

string constant TAG = "UniSwapV2Action";

contract UniSwapV2Action is ReentrancyGuardUpgradeable, ContextUpgradeable {
    // Request types
    uint8 internal constant TYPE_SWAPEXACTTOKENSFORTOKENS = 1;
    uint8 internal constant TYPE_SWAPTOKENSFOREXACTTOKENS = 2;

    address public routerAddress;
    address public roleManager;
    // pause swap
    bool public swapPause;

    //key:contractAddress
    // Example: contractAddress = keccak256("UNISWAPV2") for Uniswap V2
    mapping(address => bool) public allowedAsset;

    error DEXNotSupported(string contractTAG);
    error AssetNotSupported(string contractTAG);
    error InvalidTokenAmount(string contractTAG);
    error UnsupportedFunctionType(string contractTAG, uint8 swapType);

    /**
     * @dev Modifier that checks if caller has a specific role.
     **/
    modifier onlyRole(bytes32 role) {
        IRoleManager(roleManager).checkRole(role, _msgSender());
        _;
    }

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    /**
     * @notice Initilize the contract.
     * @param roleManagerAddress - Address of the roleManager contract.
     **/
    function initialize(address roleManagerAddress) public initializer {
        roleManager = roleManagerAddress;
    }

    /**
     * @notice Configures role manager address
     * This function can only be accessed by the account with role DEFAULT_ADMIN_ROLE.
     * @param roleManagerAddress - address of the roleManager contract.
     **/
    function setRoleManager(
        address roleManagerAddress
    ) external onlyRole(IRoleManager(roleManager).ACTIONS_ADMIN()) {
        roleManager = roleManagerAddress;
    }

    /**
     * @notice Configures router address of the dex
     * This function can only be accessed by the account with role ACTIONS_ADMIN.
     * @param _routerAddress - dex router address.
     **/
    function setDexRouter(
        address _routerAddress
    ) external onlyRole(IRoleManager(roleManager).ACTIONS_ADMIN()) {
        routerAddress = _routerAddress;
    }

    /**
     * @notice Configures the allowed asset to swap
     * This function can only be accessed by the account with role ACTIONS_ADMIN.
     * @param contractAddress - allowed asset address.
     **/
    function setAllowedAssets(
        address contractAddress,
        bool status
    ) external onlyRole(IRoleManager(roleManager).ACTIONS_ADMIN()) {
        allowedAsset[contractAddress] = status;
    }

    /**
     * @notice Configures swap status
     * This function can only be accessed by account with role ACTIONS_ADMIN.
     * @param status - status of the swap.
     **/
    function setSwapStatus(
        bool status
    ) external onlyRole(IRoleManager(roleManager).ACTIONS_ADMIN()) {
        swapPause = status;
    }

    /**
     * @notice performs the swap operation
     * This function can only be accessed by the account with role SWAP_ACTION().
     * @param payload - contains information for the swap
     **/
    function swap(
        bytes memory payload
    ) external nonReentrant onlyRole(IRoleManager(roleManager).SWAP_ACTION()) {
        address srcToken;
        address dstToken;
        address recipient;
        uint8 swapType;
        uint256 amountIn;
        uint256 amountOut;
        uint256 deadline;

        assembly {
            srcToken := mload(add(payload, 52))
            dstToken := mload(add(payload, 72))
            recipient := mload(add(payload, 92))
            swapType := mload(add(payload,93))
            amountIn := mload(add(payload, 125))
            amountOut := mload(add(payload, 157))
            deadline := mload(add(payload, 189))
        }

        if (!allowedAsset[srcToken]) revert AssetNotSupported(TAG);

        if (!allowedAsset[dstToken]) revert AssetNotSupported(TAG);

        if (amountIn == 0) revert InvalidTokenAmount(TAG);

        IERC20(srcToken).transferFrom(_msgSender(), address(this), amountIn);

        IERC20(srcToken).approve(routerAddress, amountIn);

        if (swapType == TYPE_SWAPEXACTTOKENSFORTOKENS) {
            IUniSwapV2Router(routerAddress).swapExactTokensForTokens(
                amountIn,
                amountOut,
                getPathForTokenToToken(srcToken, dstToken),
                recipient,
                deadline
            );
        } else if (swapType == TYPE_SWAPTOKENSFOREXACTTOKENS) {
            IUniSwapV2Router(routerAddress).swapTokensForExactTokens(
                amountOut,
                amountIn,
                getPathForTokenToToken(srcToken, dstToken),
                recipient,
                deadline
            );
        } else {
            revert UnsupportedFunctionType(TAG, swapType);
        }
    }

    /**
     * @notice Constructs the required parameter for the swap
     * @param _fromToken - token to send
     * @param _toToken - token to receive
     **/
    function getPathForTokenToToken(
        address _fromToken,
        address _toToken
    ) private pure returns (address[] memory) {
        address[] memory path = new address[](2);
        path[0] = _fromToken;
        path[1] = _toToken;
        return path;
    }
}

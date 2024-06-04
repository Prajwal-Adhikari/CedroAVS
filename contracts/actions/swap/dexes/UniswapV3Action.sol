// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.24;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/utils/ContextUpgradeable.sol";

import "../../interfaces/IRoleManager.sol";
import "../interfaces/IUniSwapV3Router.sol";

string constant TAG = "UniSwapV3Action";

contract UniSwapV3Action is ReentrancyGuardUpgradeable, ContextUpgradeable {
    // Request types
    uint8 internal constant TYPE_EXACTINPUTSINGLE = 1;
    uint8 internal constant TYPE_EXACTINPUT = 2;
    uint8 internal constant TYPE_EXACTOUTPUTSINGLE = 3;
    uint8 internal constant TYPE_EXACTOUTPUT = 4;

    // pause swap
    bool public swapPause;
    address public routerAddress;
    address public roleManager;

    //key:contractAddress
    // Example: contractAddress =  token address of USDC
    mapping(address => bool) public allowedAsset;

    event ReceiveToken(address user, bytes32 id, uint256 tokenAmount);

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
        address middleWareToken;
        address dstToken;
        uint8 swapType;
        uint256 amountIn;
        uint256 amountOut;
        uint256 deadline;
        uint24 fee;
        address recipient;
        uint160 sqrtPriceLimitX96;

        assembly {
            srcToken := mload(add(payload, 52))
            middleWareToken := mload(add(payload, 72))
            dstToken := mload(add(payload, 92))
            swapType := mload(add(payload, 93))
            amountIn := mload(add(payload, 125))
            amountOut := mload(add(payload, 157))
            deadline := mload(add(payload, 189))
            fee := mload(add(payload, 192))
            recipient := mload(add(payload, 212))
            sqrtPriceLimitX96 := mload(add(payload, 232))
        }

        if (!allowedAsset[srcToken]) revert AssetNotSupported(TAG);

        if (!allowedAsset[dstToken]) revert AssetNotSupported(TAG);

        IERC20(srcToken).approve(routerAddress, amountIn);

        if (swapType == TYPE_EXACTINPUTSINGLE) {
            IUniSwapV3Router.ExactInputSingleParams
                memory params = IUniSwapV3Router.ExactInputSingleParams(
                    srcToken,
                    dstToken,
                    fee,
                    recipient,
                    deadline,
                    amountIn,
                    amountOut,
                    sqrtPriceLimitX96
                );

            IUniSwapV3Router(routerAddress).exactInputSingle(params);
        } else if (swapType == TYPE_EXACTINPUT) {
            IUniSwapV3Router.ExactInputParams memory params = IUniSwapV3Router
                .ExactInputParams(
                    getPathForTokenToToken(
                        srcToken,
                        middleWareToken,
                        dstToken,
                        fee
                    ),
                    recipient,
                    deadline,
                    amountIn,
                    amountOut
                );

            IUniSwapV3Router(routerAddress).exactInput(params);
        } else if (swapType == TYPE_EXACTOUTPUT) {
            IUniSwapV3Router.ExactOutputParams memory params = IUniSwapV3Router
                .ExactOutputParams(
                    getPathForTokenToToken(
                        srcToken,
                        middleWareToken,
                        dstToken,
                        fee
                    ),
                    recipient,
                    deadline,
                    amountOut,
                    amountIn
                );

            IUniSwapV3Router(routerAddress).exactOutput(params);
        } else if (swapType == TYPE_EXACTOUTPUTSINGLE) {
            IUniSwapV3Router.ExactOutputSingleParams
                memory params = IUniSwapV3Router.ExactOutputSingleParams(
                    srcToken,
                    dstToken,
                    fee,
                    recipient,
                    deadline,
                    amountOut,
                    amountIn,
                    sqrtPriceLimitX96
                );

            IUniSwapV3Router(routerAddress).exactOutputSingle(params);
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
        address middleWareToken,
        address _toToken,
        uint24 fee
    ) private pure returns (bytes memory path) {
        path = abi.encodePacked(
            _fromToken,
            fee,
            middleWareToken,
            fee,
            _toToken
        );
    }
}

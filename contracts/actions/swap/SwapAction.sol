// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.24;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/utils/ContextUpgradeable.sol";

import "../interfaces/IRoleManager.sol";
import "./interfaces/ISwapRouter.sol";

string constant TAG = "SwapAction";

/**
 * @title SwapAction
 * @notice Swap Action is deployed on all the chains Cedro supports.
 * @notice Operator interact with swap action to swap the funds.
 **/

contract SwapActions is ContextUpgradeable, ReentrancyGuardUpgradeable {
    //address of the rolemanager contract
    address public roleManager;
    // swap pause for all dex
    bool public swapPause;

    //Information required for the dex
    struct DexInfo {
        bool enable;
        uint128 payloadLength;
        address dexContractAddress;
    }

    //key: dexName
    // Example: dexName = DexInfo 
    mapping(bytes32 => DexInfo) public dexes;

    event Swap(address user, bytes32 dexName);

    error InvalidDEX(string contractTAG);
    error SwapActionPaused(string contractTAG);
    error DexSwapPaused(string contractTAG);
    error DEXNotSupported(string contractTAG);
    error BadPayload(
        string contractTAG,
        uint256 currentLength,
        uint256 requiredLength
    );

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
     * @notice Configures the dex allowed
     * This function can only be accessed by account with role ACTIONS_ADMIN.
     * @param dexName - name of the dex
     * @param enable - status of the dex
     * @param payloadLength - length of the payload allowed
     * @param dexContractAddress - address of the dex contract deployed
     **/
    function setDex(
        bytes32 dexName,
        bool enable,
        uint128 payloadLength,
        address dexContractAddress
    ) external onlyRole(IRoleManager(roleManager).ACTIONS_ADMIN()) {
        DexInfo storage dex = dexes[dexName];
        dex.enable = enable;
        dex.dexContractAddress = dexContractAddress;
        dex.payloadLength = payloadLength;
    }

    /**
     * @notice Performs a swap operation on a specified decentralized exchange (DEX)
     * @dev  Different dex requires different payload
     * This function can only be accessed by account with role OPERATOR_CORE.
     * @param payload - contains information for the swap
     **/

    function performSwap(
        bytes memory payload
    )
        external
        nonReentrant
        onlyRole(IRoleManager(roleManager).OPERATOR_CORE())
    {
        if (swapPause) revert SwapActionPaused(TAG);

        bytes32 dexName;
        address srcToken;
        uint256 amountIn;

        assembly {
            dexName := mload(add(payload, 32))
            srcToken := mload(add(payload, 52))
            amountIn := mload(add(payload, 124))
        }

        address dexContract = dexes[dexName].dexContractAddress;

        if (dexContract == address(0)) revert DEXNotSupported(TAG);

        if (dexes[dexName].enable) revert DexSwapPaused(TAG);

        if (payload.length != dexes[dexName].payloadLength)
            revert BadPayload(
                TAG,
                payload.length,
                dexes[dexName].payloadLength
            );

        IERC20(srcToken).transferFrom(_msgSender(), dexContract, amountIn);

        ISwapRouter(dexContract).swap(payload);

        emit Swap(_msgSender(), dexName);
    }
}

// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.24;

// @title swap interface for different dex

interface ISwapRouter {
    /**
     * @notice performs the swap.
     * This function can only be accessed by account with role SwapAction.
     * @param payload - information containing token amount and receipient.
     **/
    function swap(bytes memory payload) external;

    /**
     * @notice Configures role manager address
     * This function can only be accessed by the account with role ACTION_ADMIN.
     * @param roleManagerAddress - address of the roleManager contract.
     **/
    function setRoleManager(address roleManagerAddress) external;

    /**
     * @notice Configures dex router address
     * This function can only be accessed by the account with role DEFAULT_ADMIN.
     * @param routerAddress - dex router address.
     **/
    function setDexRouter(address routerAddress) external;

    /**
     * @notice Configures the allowed asset to swap
     * This function can only be accessed by the account with role ACTIONS_ADMIN.
     * @param contractAddress - allowed asset address.
     **/
    function setAllowedAssets(address contractAddress, bool status) external;
}

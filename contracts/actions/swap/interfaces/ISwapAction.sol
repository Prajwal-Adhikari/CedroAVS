// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.24;

// @title swap interface of different dex

interface ISwapActions {
    /**
     * @notice Configures role manager address
     * This function can only be accessed by the account with role ACTION_ADMIN.
     * @param roleManagerAddress - address of the roleManager contract.
     **/
    function setRoleManager(address roleManagerAddress) external;

    /**
     * @notice Configures swap status
     * This function can only be accessed by account with role ACTIONS_ADMIN.
     * @param status - status of the swap.
     **/
    function setSwapStatus(bool status) external;

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
    ) external;

    /**
     * @notice Swaps the required tokens
     * This function can only be accessed by account with role Operator Core.
     * @param payload - information containing swap data.
     **/
    function performSwap(bytes memory payload) external;
}

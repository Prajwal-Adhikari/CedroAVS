// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.24;

interface IPriceFeed {
    /**
     * @notice Configures role manager address
     * This function can only be accessed by the account with role DEFAULT_ADMIN_ROLE.
     * @param roleManagerAddress - address of the roleManager contract.
     **/
    function setRoleManager(address roleManagerAddress) external;

    /**
     * @notice Proposes the prices of assets
     * This function can only be accessed by the account with role PROPOSE_PRICE.
     * @param _symbols - list of pool ids of the assets
     * @param _prices - prices of the asset
     **/
    function proposePrices(
        bytes32[] memory _symbols,
        uint256[] memory _prices
    ) external;

    /** @notice Validates the prices of assets
     * This function can only be accessed by the account with role VALIDATE_PRICE.
     * @param _symbols - list of pool ids of the assets
     * @param _prices - prices of the asset
     **/
    function validatePrices(
        bytes32[] memory _symbols,
        uint256[] memory _prices
    ) external;

    /**
     * @notice Returns price of assets
     * @param _symbols - list of the asset's pool id
     **/
    function getPrices(
        bytes32[] memory _symbols
    ) external view returns (uint256[] memory results);

    /**
     * @notice Returns price of an asset
     * @param _symbol - pool id of the asset
     **/
    function getPrice(bytes32 _symbol) external view returns (uint256 result);
}

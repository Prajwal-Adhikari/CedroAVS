// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.24;

interface IAIAgent {

    /**
     * @notice Approves list of asset from Operator contract
     * This function can only be accessed by _owner
     * @param _assets - list of ids of the asset
     **/
    function approveAsset(bytes32[] memory _assets) external;

    /**
     * @notice Removes asset from the approved list
     * This function can only be accessed by the account with role AIAGENT.
     * @param _asset - id of the asset
     **/
    function delistAsset(bytes32 _asset) external;

    /**
     * @notice Accept deposits from user
     * This function can only be accessed by the account with role Agent.
     * @param _asset - id of the asset
     * @param _amount - amount of the asset
     * @param _timeInterval - time for which user wants to deposit
     * @param _interestRate - rate provided by AI agent
     * @param _timeThreshold - block time till which signature is valid
     * @param _signature - msg signed by the AI agent
     **/
    function depositAsset(
        bytes32 _asset,
        uint256 _amount,
        uint256 _timeInterval,
        uint256 _interestRate,
        uint256 _timeThreshold,
        bytes memory _signature
    ) external payable;

    /**
     * @notice Emits event with assets data and total liquidity in USD
     * This function can only be accessed by the account with role SNAPSHOT_OPERATOR
     * @param _requestId - request number
     **/
    function snapshotAssetData(uint256 _requestId) external;

    /**
     * @notice Returns balance of the asset holded by the contract
     * @param _token - address of the asset
     **/
    function getAssetBalance(address _token)
        external
        view
        returns (uint256 assetBalance);
}

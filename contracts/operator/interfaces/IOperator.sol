// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.24;

interface IOperator {

    /**
     * @notice Enrolls agent to the operator
     * This function can only be accessed by the account with role AI_AGENT.
     **/
    function enrollAgent() external;

    /**
     * @notice Removes agent from the operator
     * This function can only be accessed by the account with role OPERATOR.
     * @param _agent - address of the agent
     **/
    function delistAgent(address _agent) external;

    /**
     * @notice Approves list of asset from Operator Manager contract
     * This function can only be accessed by the account with role OPERATOR.
     * @param _assets - list of ids of the asset
     **/
    function approveAsset(bytes32[] memory _assets) external;

    /**
     * @notice Removes asset from the approved list
     * This function can only be accessed by the account with role OPERATOR.
     * @param _asset - id of the asset
     **/
    function delistAsset(bytes32 _asset) external;

    /**
     * @notice Updates bond amount by agent
     * This function can only be accessed by the account with role Agent.
     **/
    function updateBondAmount() external payable;

    /**
     * @notice Returns the status of the asset
     * @param asset - id of the asset.
     **/
    function isAssetAccepted(bytes32 asset) external view returns (bool);
}

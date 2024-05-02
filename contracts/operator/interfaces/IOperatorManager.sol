// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.24;
import "../libraries/Data.sol";

interface IOperatorManager {
    /**
     * @notice Returns the pool info of the asset
     * @param _assetId - id of the asset.
     **/
    function pools(
        bytes32 _assetId
    ) external view returns (Data.PoolInfo memory);

    function mapAgent(
        address _agentContract,
        address _operatorContract
    ) external;
}

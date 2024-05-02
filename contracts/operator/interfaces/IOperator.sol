// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.24;
import "../libraries/Data.sol";

interface IOperator {
    function isAssetAccepted(bytes32 asset) external view returns (bool);
}

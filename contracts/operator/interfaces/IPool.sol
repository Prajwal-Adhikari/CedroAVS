// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.24;


interface IPool {

    function trasferAssetToDepositor(
        address _user,
        address[] memory _token,
        uint256[] memory _amount
    ) external ;

    function takeAction(bytes memory payload, address protocol) external;
}
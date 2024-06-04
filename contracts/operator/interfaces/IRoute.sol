// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.24;


interface IRoute {
    function send(
        bytes memory payload,
        bytes32[] memory chainIDs,
        address sender
    ) external payable;

   
}

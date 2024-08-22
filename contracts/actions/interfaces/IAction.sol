// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.24;

interface IAction {


    function extractPayloads(
        bytes memory _payload
    ) external view returns (bytes memory, uint256, address);


}

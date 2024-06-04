// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.24;

interface IAction {
    function transferPayload(bytes memory payload, address to) external;

    function extractPayload(
        bytes memory payload,
        address user,
        address oppurtunity
    ) external;

    function protocol(bytes32 id) external view returns (address);

    function protocolInterface(bytes32 id) external view returns (address);
}

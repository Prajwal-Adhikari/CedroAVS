// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.24;

interface IAdapter {
    function enrollAgents(
        address agentManager,
        address agent,
        address oracle,
        bytes32[] memory chainIDs,
        uint256[] memory fees
    ) external payable;

    function enrollAgent(
        address agentManager,
        address agent,
        address oracle
    ) external payable;
}

// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.24;

interface ICompound {
    function supply(address asset, uint amount) external;

    function withdraw(address asset, uint amount) external;

    function claim(address comet, address src, bool shouldAccrue) external;

    function deposit() external payable;
}

// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.24;

interface ILido {
   function approve(address _spender, uint256 _amount) external returns (bool);
   function submit(address _referral) external payable returns (uint256);
}
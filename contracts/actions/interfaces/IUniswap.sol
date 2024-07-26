// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.24;

// @title swap interface for Uniswap V2 and Sushiswap

interface IUniswap {
    /// @notice Function for swapping tokens via Uniswap V2 and Sushiswap
    function swapExactTokensForTokens(
        uint amountIn,
        uint amountOutMin,
        address[] calldata path,
        address to,
        uint deadline
    ) external returns (uint[] memory amounts);

    function swapTokensForExactTokens(
        uint amountOut,
        uint amountInMax,
        address[] calldata path,
        address to,
        uint deadline
    ) external returns (uint[] memory amounts);
}


// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.24;

import { MessagingFee, OFTReceipt, SendParam } from "@layerzerolabs/lz-evm-oapp-v2/contracts/oft/interfaces/IOFT.sol";

interface IStargateInterface {
    function prepareTakeTaxi(
        address _stargate,
        uint32 _dstEid,
        uint256 _amount,
        address _receiver
    ) external view returns (uint256 valueToSend, SendParam memory sendParam, MessagingFee memory messagingFee);
    

}

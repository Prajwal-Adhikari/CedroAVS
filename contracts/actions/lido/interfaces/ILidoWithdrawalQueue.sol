// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.24;
import "../common/Structs.sol";

interface ILidoWithdrawalQueue {
   
    function requestWithdrawals(uint256[] calldata _amounts, address _owner)
        external 
        returns (uint256[] memory requestIds);

    function getWithdrawalRequests(
        address _owner
    ) external view returns (uint256[] memory requestsIds);

    function getWithdrawalStatus(
        uint256[] calldata _requestIds
    ) external view returns (WithdrawalRequestStatus[] memory statuses);

    function findCheckpointHints(
        uint256[] calldata _requestIds,
        uint256 _firstIndex,
        uint256 _lastIndex
    ) external view returns (uint256[] memory hintIds);

    function getLastCheckpointIndex() external view returns (uint256);

    function claimWithdrawals(
        uint256[] calldata _requestIds,
        uint256[] calldata _hints
    ) external;


    function claimWithdrawal(uint256 _requestId) external;

}

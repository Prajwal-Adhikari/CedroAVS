// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.24;

library Data {

    struct PoolInfo {
        bool enable;
        address token;
    }

    struct SnapshotInfo {
        bool inProgress;
        uint256 snapshotRequestTime;
        uint256 requestID;
    }

    struct UserDepositData {
        uint256 txnID;
        bytes32 asset;
        uint256 amount;
        uint256 principalAmount;
        uint256 timeInterval;
        uint256 interestRate;
        uint256 createdAt;
        uint256 claimedAt;
        bool claimed;
    }
   
    struct AssetsData {
        uint256 totalBalance;
        bytes32 assetId;
        address assetAddress;
    }


  
    
}

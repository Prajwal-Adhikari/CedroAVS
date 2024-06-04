// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.24;

library Data {
    enum STATUS {
        CREATED,
        MATURED,
        CLAIMED
    }

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
        address user;
        bytes32 asset;
        uint256 amount;
        uint256 principalAmount;
        uint256 timeInterval;
        uint256 interestRate;
        uint256 createdAt;
        uint256 claimedAt;
        bool claimed;
        STATUS status;
    }

    struct AssetsData {
        uint256 balance;
        bytes32 assetId;
        address assetAddress;
    }
}

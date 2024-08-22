// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;
import "./interfaces/ILidoWithdrawalQueue.sol";

contract LidoInteraction {
    address lidoContract;
    address withdrawQueueContract;

    function setLidoContract(address _address) external {
        lidoContract = _address;
    }

    function setWithdrawQueueContract(address _address) external {
        withdrawQueueContract = _address;
    }

    function extractPayload(
        bytes memory rawPayload
    )
        public
        view
        returns (
            bytes memory actionPayload,
            uint256 returnValue,
            address _protocol
        )
    {
        uint8 actionType;
        uint8 receiver;
        bytes32 id;
        uint8 methodName;
        address asset;
        uint256 amount;
        assembly {
            actionType := mload(add(rawPayload, 1))
            receiver := mload(add(rawPayload, 2))
            id := mload(add(rawPayload, 34))
            methodName := mload(add(rawPayload, 35))
            asset := mload(add(rawPayload, 55))
            amount := mload(add(rawPayload, 87))
        }

        // Deposit eth
        if (methodName == 1) {
            actionPayload = abi.encodeWithSignature(
                "submit(address)",
                address(0)
            );
            returnValue = amount;
            _protocol = lidoContract;
        }

        // Lido token approval for spender
        if (methodName == 2) {
            address _spender;

            assembly {
                _spender := mload(add(rawPayload, 107))
            }

            actionPayload = abi.encodeWithSignature(
                "approve(address,uint256)",
                _spender,
                amount
            );
            _protocol = lidoContract;
        }

        // Request withdrawal
        if (methodName == 3) {
            address _owner;

            assembly {
                _owner := mload(add(rawPayload, 107))
            }

            uint256[] memory _amount = new uint256[](1);
            _amount[0] = amount;

            actionPayload = abi.encodeWithSignature(
                "requestWithdrawals(uint256[],address)",
                _amount,
                _owner
            );
            _protocol = withdrawQueueContract;
        }

        // Claim withdrawal
        if (methodName == 4) {
            (
                uint256[] memory maturedRequestIds,
                uint256[] memory hintIds
            ) = getMaturedIdsAndHintIds();

            actionPayload = abi.encodeWithSignature(
                "claimWithdrawals(uint256[],uint256[])",
                maturedRequestIds,
                hintIds
            );
            _protocol = withdrawQueueContract;
        }
    }

    function getMaturedIdsAndHintIds()
        internal
        view
        returns (uint256[] memory maturedRequestIds, uint256[] memory hintIds)
    {
        ILidoWithdrawalQueue withdrawContract = ILidoWithdrawalQueue(
            withdrawQueueContract
        );

        uint256[] memory requestIds = withdrawContract.getWithdrawalRequests(
            msg.sender
        );

        require(
            requestIds.length > 0,
            string(
                abi.encodePacked("Withdrawal requests not found ", requestIds)
            )
        );

        WithdrawalRequestStatus[] memory reqStatusArray = withdrawContract
            .getWithdrawalStatus(requestIds);

        require(
            requestIds.length > 0,
            string(
                abi.encodePacked(
                    "Withdrawal requests status not found ",
                    reqStatusArray.length
                )
            )
        );

        // Count the number of matured requests
        uint count = 0;
        for (uint i = 0; i < reqStatusArray.length; i++) {
            if (reqStatusArray[i].isFinalized && !reqStatusArray[i].isClaimed) {
                count++;
            }
        }

        // Revert if there are no matured withdrawals
        require(count > 0, "No claimable withdrawals yet");

        // store matured requests in an array
        maturedRequestIds = new uint256[](count);

        uint index = 0;
        for (uint i = 0; i < reqStatusArray.length; i++) {
            if (reqStatusArray[i].isFinalized && !reqStatusArray[i].isClaimed) {
                maturedRequestIds[index] = requestIds[i];
                index++;
            }
        }

        require(
            requestIds.length > 0,
            string(
                abi.encodePacked(
                    "Matured Withdrawal requests status not found ",
                    maturedRequestIds
                )
            )
        );

        uint256 lastCheckpointIndex = withdrawContract.getLastCheckpointIndex();
        hintIds = withdrawContract.findCheckpointHints(
            maturedRequestIds,
            1,
            lastCheckpointIndex
        );
    }
}

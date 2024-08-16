// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.24;

import "@openzeppelin/contracts-upgradeable/utils/ContextUpgradeable.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "../interfaces/ICompound.sol";

string constant TAG = "CompoundInteraction";

/**
 * @title Compound Interaction contract
 * @notice Used to interact with the compound
 **/
contract CompoundInteraction is ContextUpgradeable {
    address rewardContract;
    address wethContract;
    address nativeToken;
    mapping(address => address) addressToInteract;

    error InvalidAsset(string contractTAG);

    function extractPayload(
        bytes memory _payload
    )
        public
        view
        returns (
            bytes memory bytesData,
            uint256 nativeAssetValue,
            address protocolAddress
        )
    {
        uint8 methodToCall;
        assembly {
            methodToCall := mload(add(_payload, 85))
        }

        if (methodToCall == 1) {
            return extractSupplyData(_payload);
        } else if (methodToCall == 2) {
            return extractWithdrawData(_payload);
        } else if (methodToCall == 3) {
            return extractRewardData(_payload);
        } else {
            return extractDepositData(_payload);
        }
    }

    function setRewardContract(address rewardConractAddress) external {
        rewardContract = rewardConractAddress;
    }

    function setNativeToken(address nativeTokenAddress) external {
        nativeToken = nativeTokenAddress;
    }

    function setCometAddress(
        address assetAddress,
        address cometAddress
    ) external {
        addressToInteract[assetAddress] = cometAddress;
    }

    function setWethContract(address wethContractAddress) external {
        wethContract = wethContractAddress;
    }

    function extractSupplyData(
        bytes memory _payload
    )
        internal
        view
        returns (
            bytes memory bytesData,
            uint256 nativeAssetValue,
            address protocolAddress
        )
    {
        address asset;
        uint amount;

        assembly {
            asset := mload(add(_payload, 52))
            amount := mload(add(_payload, 84))
        }
        protocolAddress = addressToInteract[asset];

        bytesData = abi.encodeWithSelector(
            ICompound.supply.selector,
            asset,
            amount
        );
        nativeAssetValue = 0;
    }

    function extractWithdrawData(
        bytes memory _payload
    )
        internal
        view
        returns (
            bytes memory bytesData,
            uint256 nativeAssetValue,
            address protocolAddress
        )
    {
        address asset;
        uint amount;

        assembly {
            asset := mload(add(_payload, 52))
            amount := mload(add(_payload, 84))
        }
        bytesData = abi.encodeWithSelector(
            ICompound.withdraw.selector,
            asset,
            amount
        );

        nativeAssetValue = 0;
        protocolAddress = addressToInteract[asset];
    }

    function extractRewardData(
        bytes memory _payload
    )
        internal
        view
        returns (
            bytes memory bytesData,
            uint256 nativeAssetValue,
            address protocolAddress
        )
    {
        address comet;
        address src;
        bool shouldAccrue;

        assembly {
            comet := mload(add(_payload, 105))
            src := mload(add(_payload, 125))
            shouldAccrue := mload(add(_payload, 126))
        }
        bytesData = abi.encodeWithSelector(
            ICompound.claim.selector,
            comet,
            src,
            shouldAccrue
        );
        nativeAssetValue = 0;
        protocolAddress = rewardContract;
    }

    function extractDepositData(
        bytes memory _payload
    )
        internal
        view
        returns (
            bytes memory bytesData,
            uint256 nativeAssetValue,
            address protocolAddress
        )
    {
        address asset;
        uint256 amount;
        assembly {
            asset := mload(add(_payload, 52))
            amount := mload(add(_payload, 84))
        }

        if (asset != nativeToken) revert InvalidAsset(TAG);

        bytesData = abi.encodeWithSelector(ICompound.deposit.selector);

        nativeAssetValue = amount;

        protocolAddress = wethContract;
    }
}

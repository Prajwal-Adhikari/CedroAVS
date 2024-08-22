// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.24;

import {ContextUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/ContextUpgradeable.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";

import {IGmx} from "./interfaces/IGmx.sol";
import {IRoleManager} from "../interfaces/IRoleManager.sol";

string constant TAG = "GmxInteraction";

/**
 * @title AaveInterface contract
 * @notice Used to iterface call for Aave methods
 **/
contract GMXInteraction is ContextUpgradeable {
    address nativeToken;
    address orderVault;
    address depositVault;
    address gmxProtocol;
    address roleManagerAddress;

    mapping(address => bool) allowedAsset;

    error UnsupportedFunctionType(string contractTAG, uint8 functionType);

    modifier onlyRoles(bytes32 role) {
        IRoleManager(roleManagerAddress).checkRole(role, _msgSender());
        _;
    }

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    /**
     * @notice Initilize the contract.
     **/
    function initialize(
        address _nativeTokenAddrees,
        address _orderVaultAddress,
        address _depositVaultAddress,
        address _gmxProtocolAddress,
        address _roleManagerAddress
    ) public initializer {
        nativeToken = _nativeTokenAddrees;
        orderVault = _orderVaultAddress;
        depositVault = _depositVaultAddress;
        gmxProtocol = _gmxProtocolAddress;
        roleManagerAddress = _roleManagerAddress;
    }

    function extractPayload(
        bytes memory _payload
    )
        public
        view
        returns (
            bytes memory extractedPayload,
            uint256 nativeAssetAmountValue,
            address protocolAddress
        )
    {
        uint8 methodToCall;
        assembly {
            methodToCall := mload(add(_payload, 85))
        }

        protocolAddress = gmxProtocol;

        if (methodToCall == 1) {
            bytes[] memory extractedData;
            (extractedData, nativeAssetAmountValue) = extractOrderArray(
                _payload
            );
            extractedPayload = abi.encodeWithSelector(
                IGmx.multicall.selector,
                extractedData
            );
        } else if (methodToCall == 2) {
            bytes[] memory extractedData;
            (extractedData, nativeAssetAmountValue) = extractUpdateOrderArray(
                _payload
            );
            extractedPayload = abi.encodeWithSelector(
                IGmx.multicall.selector,
                extractedData
            );
        } else if (methodToCall == 3) {
            bytes[] memory extractedData;
            (extractedData, nativeAssetAmountValue) = extractCancelOrderArray(
                _payload
            );
            extractedPayload = abi.encodeWithSelector(
                IGmx.multicall.selector,
                extractedData
            );
        } else if (methodToCall == 4) {
            bytes[] memory extractedData;
            (extractedData, nativeAssetAmountValue) = extractDepositArray(
                _payload
            );
            extractedPayload = abi.encodeWithSelector(
                IGmx.multicall.selector,
                extractedData
            );
        } else if (methodToCall == 5) {
            bytes[] memory extractedData;
            (extractedData, nativeAssetAmountValue) = extractCancelDepositArray(
                _payload
            );
            extractedPayload = abi.encodeWithSelector(
                IGmx.multicall.selector,
                extractedData
            );
        } else if (methodToCall == 6) {
            bytes[] memory extractedData;
            (extractedData, nativeAssetAmountValue) = extractWithdrawalArray(
                _payload
            );
            extractedPayload = abi.encodeWithSelector(
                IGmx.multicall.selector,
                extractedData
            );
        } else if (methodToCall == 7) {
            bytes[] memory extractedData;
            (
                extractedData,
                nativeAssetAmountValue
            ) = extractCancelWithdrawalArray(_payload);
            extractedPayload = abi.encodeWithSelector(
                IGmx.multicall.selector,
                extractedData
            );
        } else if (methodToCall == 8) {
            bytes[] memory extractedData;
            (extractedData, nativeAssetAmountValue) = extractCreateShiftArray(
                _payload
            );
            extractedPayload = abi.encodeWithSelector(
                IGmx.multicall.selector,
                extractedData
            );
        } else if (methodToCall == 9) {
            bytes[] memory extractedData;
            (extractedData, nativeAssetAmountValue) = extractCancelShiftArray(
                _payload
            );
            extractedPayload = abi.encodeWithSelector(
                IGmx.multicall.selector,
                extractedData
            );
        } else if (methodToCall == 10) {
            bytes[] memory extractedData;
            (
                extractedData,
                nativeAssetAmountValue
            ) = extractClaimAffiliateRewardsArray(_payload);
            extractedPayload = abi.encodeWithSelector(
                IGmx.multicall.selector,
                extractedData
            );
        } else if (methodToCall == 11) {
            bytes[] memory extractedData;
            (
                extractedData,
                nativeAssetAmountValue
            ) = extractClaimCollateralArray(_payload);
            extractedPayload = abi.encodeWithSelector(
                IGmx.multicall.selector,
                extractedData
            );
        } else {
            bytes[] memory extractedData;
            (
                extractedData,
                nativeAssetAmountValue
            ) = extractClaimFundingFeesArray(_payload);
            extractedPayload = abi.encodeWithSelector(
                IGmx.multicall.selector,
                extractedData
            );
        }
    }

    function extractOrderFirstFunction(
        bytes memory _payload
    ) public view returns (bytes memory) {
        uint256 executingFees;

        assembly {
            executingFees := mload(add(_payload, 193))
        }

        return
            abi.encodeWithSelector(
                IGmx.sendWnt.selector,
                orderVault,
                executingFees
            );
    }

    function extractOrderSecondFunction(
        bytes memory _payload
    ) public view returns (bytes memory) {
        address tokenToSend;
        uint256 sendingTokenAmount;

        assembly {
            tokenToSend := mload(add(_payload, 393))
        }
        sendingTokenAmount = 1;

        return
            abi.encodeWithSelector(
                IGmx.sendTokens.selector,
                tokenToSend,
                orderVault,
                sendingTokenAmount
            );
    }

    function extractOrderArray(
        bytes memory _payload
    )
        public
        view
        returns (bytes[] memory dataArray, uint256 nativeAssetAmount)
    {
        address asset;
        uint256 amount;
        uint256 executionFee;
        assembly {
            asset := mload(add(_payload, 393))
            amount := mload(add(_payload, 65))
            executionFee := mload(add(_payload, 65))
        }

        bytes memory firstFunction = extractOrderFirstFunction(_payload);
        bytes memory orderData = extractOrderParams(_payload);

        if (asset == nativeToken) {
            dataArray = new bytes[](2);
            dataArray[0] = firstFunction;
            dataArray[1] = orderData;
        } else {
            dataArray = new bytes[](3);
            bytes memory secondFunction = extractOrderSecondFunction(_payload);
            dataArray[0] = firstFunction;
            dataArray[1] = secondFunction;
            dataArray[2] = orderData;
        }

        if (asset == nativeToken) {
            nativeAssetAmount = amount + executionFee;
        } else {
            nativeAssetAmount = executionFee;
        }
    }

    function extractUpdateOrderArray(
        bytes memory _payload
    )
        public
        pure
        returns (bytes[] memory dataArray, uint256 nativeAssetAmount)
    {
        bytes32 key;
        uint256 sizeDeltaUsd;
        uint256 acceptablePrice;
        uint256 triggerPrice;
        uint256 minOutputAmount;
        bool autoCancel;

        assembly {
            key := mload(add(_payload, 117))
            sizeDeltaUsd := mload(add(_payload, 149))
            acceptablePrice := mload(add(_payload, 181))
            triggerPrice := mload(add(_payload, 213))
            minOutputAmount := mload(add(_payload, 245))
            autoCancel := mload(add(_payload, 246))
        }

        dataArray = new bytes[](1);
        bytes memory updateOrderBytes = abi.encodeWithSelector(
            IGmx.updateOrder.selector,
            key,
            sizeDeltaUsd,
            acceptablePrice,
            minOutputAmount,
            autoCancel
        );
        dataArray[0] = updateOrderBytes;

        nativeAssetAmount = 0;
    }

    function extractCancelOrderArray(
        bytes memory _payload
    )
        public
        pure
        returns (bytes[] memory dataArray, uint256 nativeAssetAmount)
    {
        bytes32 key;
        assembly {
            key := mload(add(_payload, 117))
        }

        dataArray = new bytes[](1);
        bytes memory cancelOrderBytes = abi.encodeWithSelector(
            IGmx.cancelOrder.selector,
            key
        );
        dataArray[0] = cancelOrderBytes;
        nativeAssetAmount = 0;
    }

    function extractDepositArray(
        bytes memory _payload
    )
        public
        view
        returns (bytes[] memory dataArray, uint256 nativeAssetAmount)
    {
        address asset;
        uint256 amount;
        uint256 executingFee;

        assembly {
            asset := mload(add(_payload, 52))
            amount := mload(add(_payload, 84))
            executingFee := mload(add(_payload, 84))
        }

        bytes memory firstFunction = extractOrderFirstFunction(_payload);
        bytes memory depositData = extractDepositParams(_payload);

        if (asset == nativeToken) {
            dataArray = new bytes[](2);
            dataArray[0] = firstFunction;
            dataArray[1] = depositData;
        } else {
            dataArray = new bytes[](3);
            bytes memory secondFunction = extractOrderSecondFunction(_payload);
            dataArray[0] = firstFunction;
            dataArray[1] = secondFunction;
            dataArray[2] = depositData;
        }

        if (asset == nativeToken) {
            nativeAssetAmount = amount + executingFee;
        } else {
            nativeAssetAmount = executingFee;
        }
    }

    function extractCancelDepositArray(
        bytes memory _payload
    )
        public
        pure
        returns (bytes[] memory dataArray, uint256 nativeAssetAmount)
    {
        bytes32 key;
        assembly {
            key := mload(add(_payload, 393))
        }

        dataArray = new bytes[](1);
        bytes memory cancelOrderBytes = abi.encodeWithSelector(
            IGmx.cancelDeposit.selector,
            key
        );
        dataArray[0] = cancelOrderBytes;
        nativeAssetAmount = 0;
    }

    function extractWithdrawalArray(
        bytes memory _payload
    )
        public
        view
        returns (bytes[] memory dataArray, uint256 nativeAssetAmount)
    {
        address asset;
        uint256 executionFee;

        assembly {
            asset := mload(add(_payload, 393))
            executionFee := mload(add(_payload, 393))
        }

        bytes memory firstFunction = extractOrderFirstFunction(_payload);
        bytes memory withdrawData = extractWithdrawalParams(_payload);

        dataArray = new bytes[](2);
        dataArray[0] = firstFunction;
        dataArray[1] = withdrawData;

        nativeAssetAmount = executionFee;
    }

    function extractCancelWithdrawalArray(
        bytes memory _payload
    )
        public
        pure
        returns (bytes[] memory dataArray, uint256 nativeAssetAmount)
    {
        bytes32 key;
        uint256 executingFee;
        assembly {
            key := mload(add(_payload, 393))
        }
        dataArray = new bytes[](1);
        bytes memory cancelOrderBytes = abi.encodeWithSelector(
            IGmx.cancelWithdrawal.selector,
            key
        );
        dataArray[0] = cancelOrderBytes;
        nativeAssetAmount = executingFee;
    }

    function extractCreateShiftArray(
        bytes memory _payload
    )
        public
        view
        returns (bytes[] memory dataArray, uint256 nativeAssetAmount)
    {
        uint256 executionFee;
        assembly {
            executionFee := mload(add(_payload, 65))
        }

        bytes memory firstFunction = extractOrderFirstFunction(_payload);
        bytes memory shiftData = extractShiftParams(_payload);

        dataArray = new bytes[](2);
        dataArray[0] = firstFunction;
        dataArray[1] = shiftData;

        nativeAssetAmount = executionFee;
    }

    function extractCancelShiftArray(
        bytes memory _payload
    )
        public
        pure
        returns (bytes[] memory dataArray, uint256 nativeAssetAmount)
    {
        bytes32 key;
        assembly {
            key := mload(add(_payload, 117))
        }
        dataArray = new bytes[](1);
        bytes memory cancelShiftBytes = abi.encodeWithSelector(
            IGmx.cancelShift.selector,
            key
        );
        dataArray[0] = cancelShiftBytes;
        nativeAssetAmount = 0;
    }

    function extractClaimAffiliateRewardsArray(
        bytes memory _payload
    )
        public
        pure
        returns (bytes[] memory dataArray, uint256 nativeAssetAmount)
    {
        uint8 marketsLength;
        uint8 tokensLength;
        address receiver;
        address[] memory markets;
        address[] memory tokens;

        assembly {
            marketsLength := mload(add(_payload, 86))
            tokensLength := mload(add(_payload, 87))
            receiver := mload(add(_payload, 107))
        }

        markets = new address[](marketsLength);
        tokens = new address[](tokensLength);

        // Pointer to read addresses from the payload
        uint256 pointer = 105; // Start reading addresses after receiver

        for (uint256 i = 0; i < marketsLength; i++) {
            assembly {
                mstore(
                    add(markets, mul(i, 0x20)),
                    mload(add(_payload, pointer))
                )
            }
            pointer += 20;
        }

        for (uint256 i = 0; i < tokensLength; i++) {
            assembly {
                mstore(add(tokens, mul(i, 0x20)), mload(add(_payload, pointer)))
            }
        }

        bytes memory claimAffiliateRewardsData = abi.encodeWithSelector(
            IGmx.claimAffiliateRewards.selector,
            markets,
            tokens,
            receiver
        );

        dataArray = new bytes[](1);
        dataArray[0] = claimAffiliateRewardsData;
        nativeAssetAmount = 0;
    }

    function extractClaimCollateralArray(
        bytes memory _payload
    )
        public
        pure
        returns (bytes[] memory dataArray, uint256 nativeAssetAmount)
    {
        uint8 marketsLength;
        uint8 tokensLength;
        uint8 timeKeysLength;
        address receiver;
        address[] memory markets;
        address[] memory tokens;
        uint256[] memory timeKeys;

        assembly {
            marketsLength := mload(add(_payload, 86))
            tokensLength := mload(add(_payload, 87))
            timeKeysLength := mload(add(_payload, 88))
            receiver := mload(add(_payload, 108))
        }

        markets = new address[](marketsLength);
        tokens = new address[](tokensLength);
        timeKeys = new uint256[](timeKeysLength);

        // Pointer to read addresses from the payload
        uint256 pointer = 108; // Start reading addresses after receiver

        for (uint256 i = 0; i < marketsLength; i++) {
            assembly {
                mstore(
                    add(markets, mul(i, 0x20)),
                    mload(add(_payload, pointer))
                )
            }
            pointer += 20;
        }

        for (uint256 i = 0; i < tokensLength; i++) {
            assembly {
                mstore(add(tokens, mul(i, 0x20)), mload(add(_payload, pointer)))
            }
            pointer += 20;
        }

        for (uint256 i = 0; i < tokensLength; i++) {
            assembly {
                mstore(
                    add(timeKeys, mul(i, 0x32)),
                    mload(add(_payload, pointer))
                )
            }
        }

        bytes memory claimCollateralData = abi.encodeWithSelector(
            IGmx.claimCollateral.selector,
            markets,
            tokens,
            timeKeys,
            receiver
        );
        dataArray = new bytes[](1);
        dataArray[0] = claimCollateralData;
        nativeAssetAmount = 0;
    }

    function extractClaimFundingFeesArray(
        bytes memory _payload
    )
        public
        pure
        returns (bytes[] memory dataArray, uint256 nativeAssetAmount)
    {
        uint8 marketsLength;
        uint8 tokensLength;
        address receiver;
        address[] memory markets;
        address[] memory tokens;

        assembly {
            marketsLength := mload(add(_payload, 86))
            tokensLength := mload(add(_payload, 87))
            receiver := mload(add(_payload, 107))
        }

        markets = new address[](marketsLength);
        tokens = new address[](tokensLength);

        // Pointer to read addresses from the payload
        uint256 pointer = 107; // Start reading addresses after receiver

        for (uint256 i = 0; i < marketsLength; i++) {
            assembly {
                mstore(
                    add(markets, mul(i, 0x20)),
                    mload(add(_payload, pointer))
                )
            }
            pointer += 20;
        }

        for (uint256 i = 0; i < tokensLength; i++) {
            assembly {
                mstore(add(tokens, mul(i, 0x20)), mload(add(_payload, pointer)))
            }
        }

        bytes memory claimCollateralData = abi.encodeWithSelector(
            IGmx.claimFundingFees.selector,
            markets,
            tokens,
            receiver
        );
        dataArray = new bytes[](1);
        dataArray[0] = claimCollateralData;
        nativeAssetAmount = 0;
    }

    function extractOrderParams(
        bytes memory _payload
    ) internal pure returns (bytes memory) {
        uint8 orderType;
        uint8 decreasePositionSwapType;
        bool isLong;
        bool shouldUnwrapNativeToken;
        bytes32 referralCode;

        assembly {
            // Load other parameters
            orderType := mload(add(_payload, 258))
            decreasePositionSwapType := mload(add(_payload, 259))
            isLong := mload(add(_payload, 260))
            shouldUnwrapNativeToken := mload(add(_payload, 261))
            referralCode := mload(add(_payload, 293))
        }

        IGmx.CreateOrderParamsAddresses
            memory orderParamsAddresses = extractOrderParamsAddresses(_payload);
        IGmx.CreateOrderParamsNumbers
            memory orderParamsNumbers = extractOrderParamsNumbers(_payload);

        IGmx.CreateOrderParams memory params = IGmx.CreateOrderParams({
            addresses: orderParamsAddresses,
            numbers: orderParamsNumbers,
            orderType: IGmx.OrderType(orderType),
            decreasePositionSwapType: IGmx.DecreasePositionSwapType(
                decreasePositionSwapType
            ),
            isLong: isLong,
            shouldUnwrapNativeToken: shouldUnwrapNativeToken,
            referralCode: referralCode
        });

        return abi.encodeWithSelector(IGmx.createOrder.selector, params);
    }

    function extractOrderParamsAddresses(
        bytes memory _payload
    ) internal pure returns (IGmx.CreateOrderParamsAddresses memory) {
        address receiver;
        address callbackContract;
        address uiFeeReceiver;
        address market;
        address initialCollateralToken;
        address[] memory swapPath;

        assembly {
            // Load addresses
            receiver := mload(add(_payload, 313))
            callbackContract := mload(add(_payload, 333))
            uiFeeReceiver := mload(add(_payload, 353))
            market := mload(add(_payload, 373))
            initialCollateralToken := mload(add(_payload, 393))
        }
        uint256 swapPathLength = (_payload.length - 393) / 32;

        swapPath = new address[](swapPathLength);
        for (uint256 i = 0; i < swapPathLength; i++) {
            assembly {
                mstore(
                    add(add(swapPath, 32), mul(i, 32)),
                    mload(add(_payload, add(425, mul(i, 32))))
                )
            }
        }

        return
            IGmx.CreateOrderParamsAddresses({
                receiver: receiver,
                callbackContract: callbackContract,
                uiFeeReceiver: uiFeeReceiver,
                market: market,
                initialCollateralToken: initialCollateralToken,
                swapPath: swapPath
            });
    }

    // Function to extract numeric parameters from the encoded bytes
    function extractOrderParamsNumbers(
        bytes memory _payload
    ) internal pure returns (IGmx.CreateOrderParamsNumbers memory) {
        // Check the length of the encoded data

        uint256 sizeDeltaUsd;
        uint256 initialCollateralDeltaAmount;
        uint256 triggerPrice;
        uint256 acceptablePrice;
        uint256 executionFee;
        uint256 callbackGasLimit;
        uint256 minOutputAmount;

        assembly {
            // Load numeric parameters
            sizeDeltaUsd := mload(add(_payload, 65))
            initialCollateralDeltaAmount := mload(add(_payload, 97))
            triggerPrice := mload(add(_payload, 129))
            acceptablePrice := mload(add(_payload, 161))
            executionFee := mload(add(_payload, 193))
            callbackGasLimit := mload(add(_payload, 225))
            minOutputAmount := mload(add(_payload, 257))
        }

        return
            IGmx.CreateOrderParamsNumbers({
                sizeDeltaUsd: sizeDeltaUsd,
                initialCollateralDeltaAmount: initialCollateralDeltaAmount,
                triggerPrice: triggerPrice,
                acceptablePrice: acceptablePrice,
                executionFee: executionFee,
                callbackGasLimit: callbackGasLimit,
                minOutputAmount: minOutputAmount
            });
    }

    function extractDepositParams(
        bytes memory _payload
    ) internal pure returns (bytes memory) {
        bool shouldUnwrapNativeToken;
        uint256 minMarketTokens;
        uint256 executionFee;
        uint256 callbackGasLimit;
        address receiver;
        address callbackContract;
        address uiFeeReceiver;
        address market;
        address initialLongToken;
        address initialShortToken;
        address[] memory longTokenSwapPath;
        address[] memory shortTokenSwapPath;

        assembly {
            // Load numeric parameters
            shouldUnwrapNativeToken := mload(add(_payload, 86))
            minMarketTokens := mload(add(_payload, 118))
            executionFee := mload(add(_payload, 150))
            callbackGasLimit := mload(add(_payload, 182))
            receiver := mload(add(_payload, 202))
            callbackContract := mload(add(_payload, 222))
            uiFeeReceiver := mload(add(_payload, 242))
            market := mload(add(_payload, 262))
            initialLongToken := mload(add(_payload, 282))
            initialShortToken := mload(add(_payload, 302))
        }

        (longTokenSwapPath, shortTokenSwapPath) = getDepositPath(_payload);

        IGmx.CreateDepositParams memory depositParams = IGmx
            .CreateDepositParams({
                receiver: receiver,
                callbackContract: callbackContract,
                uiFeeReceiver: uiFeeReceiver,
                market: market,
                initialLongToken: initialLongToken,
                initialShortToken: initialShortToken,
                longTokenSwapPath: longTokenSwapPath,
                shortTokenSwapPath: shortTokenSwapPath,
                minMarketTokens: minMarketTokens,
                shouldUnwrapNativeToken: shouldUnwrapNativeToken,
                executionFee: executionFee,
                callbackGasLimit: callbackGasLimit
            });

        return
            abi.encodeWithSelector(IGmx.createDeposit.selector, depositParams);
    }

    function getDepositPath(
        bytes memory _payload
    )
        internal
        pure
        returns (
            address[] memory longTokenSwapPath,
            address[] memory shortTokenSwapPath
        )
    {
        uint8 longTokenSwapPathLength;
        uint8 shortTokenSwapPathLength;
        assembly {
            longTokenSwapPathLength := mload(add(_payload, 303))
            shortTokenSwapPathLength := mload(add(_payload, 304))
        }

        // Pointer to read addresses from the payload
        uint256 pointer = 304; // Start reading addresses after receiver

        for (uint256 i = 0; i < longTokenSwapPathLength; i++) {
            assembly {
                mstore(
                    add(longTokenSwapPath, mul(i, 0x20)),
                    mload(add(_payload, pointer))
                )
            }
            pointer += 20;
        }

        for (uint256 i = 0; i < shortTokenSwapPathLength; i++) {
            assembly {
                mstore(
                    add(shortTokenSwapPath, mul(i, 0x20)),
                    mload(add(_payload, pointer))
                )
            }
        }
    }

    function extractWithdrawalParams(
        bytes memory _payload
    ) internal pure returns (bytes memory) {
        address receiver;
        address callbackContract;
        address uiFeeReceiver;
        address market;
        uint256 minLongTokenAmount;
        uint256 minShortTokenAmount;
        bool shouldUnwrapNativeToken;
        uint256 executionFee;
        uint256 callbackGasLimit;
        address[] memory longTokenSwapPath;
        address[] memory shortTokenSwapPath;

        assembly {
            receiver := mload(add(_payload, 105))
            callbackContract := mload(add(_payload, 125))
            uiFeeReceiver := mload(add(_payload, 145))
            market := mload(add(_payload, 165))
            minLongTokenAmount := mload(add(_payload, 197))
            minShortTokenAmount := mload(add(_payload, 229))
            shouldUnwrapNativeToken := mload(add(_payload, 230))
            executionFee := mload(add(_payload, 262))
            callbackGasLimit := mload(add(_payload, 294))
        }

        (longTokenSwapPath, shortTokenSwapPath) = getWithdrawPath(_payload);

        IGmx.CreateWithdrawalParams memory WithdrawParams = IGmx
            .CreateWithdrawalParams({
                receiver: receiver,
                callbackContract: callbackContract,
                uiFeeReceiver: uiFeeReceiver,
                market: market,
                longTokenSwapPath: longTokenSwapPath,
                shortTokenSwapPath: shortTokenSwapPath,
                minLongTokenAmount: minLongTokenAmount,
                minShortTokenAmount: minShortTokenAmount,
                shouldUnwrapNativeToken: shouldUnwrapNativeToken,
                executionFee: executionFee,
                callbackGasLimit: callbackGasLimit
            });

        return
            abi.encodeWithSelector(
                IGmx.createWithdrawal.selector,
                WithdrawParams
            );
    }

    function getWithdrawPath(
        bytes memory _payload
    )
        internal
        pure
        returns (
            address[] memory longTokenSwapPath,
            address[] memory shortTokenSwapPath
        )
    {
        uint8 longTokenSwapPathLength;
        uint8 shortTokenSwapPathLength;

        assembly {
            longTokenSwapPathLength := mload(add(_payload, 295))
            shortTokenSwapPathLength := mload(add(_payload, 296))
        }

        // Pointer to read addresses from the payload
        uint256 pointer = 296; // Start reading addresses after receiver

        for (uint256 i = 0; i < longTokenSwapPathLength; i++) {
            assembly {
                mstore(
                    add(longTokenSwapPath, mul(i, 0x20)),
                    mload(add(_payload, pointer))
                )
            }
            pointer += 20;
        }

        for (uint256 i = 0; i < shortTokenSwapPathLength; i++) {
            assembly {
                mstore(
                    add(shortTokenSwapPath, mul(i, 0x20)),
                    mload(add(_payload, pointer))
                )
            }
        }
    }

    function extractShiftParams(
        bytes memory _payload
    ) internal pure returns (bytes memory) {
        address receiver;
        address callbackContract;
        address uiFeeReceiver;
        address fromMarket;
        address toMarket;
        uint256 minMarketTokens;
        uint256 executionFee;
        uint256 callbackGasLimit;
        assembly {
            receiver := mload(add(_payload, 105))
            callbackContract := mload(add(_payload, 125))
            uiFeeReceiver := mload(add(_payload, 145))
            fromMarket := mload(add(_payload, 165))
            toMarket := mload(add(_payload, 185))
            minMarketTokens := mload(add(_payload, 217))
            executionFee := mload(add(_payload, 249))
            callbackGasLimit := mload(add(_payload, 281))
        }

        IGmx.CreateShiftParams memory createShiftParams = IGmx
            .CreateShiftParams({
                receiver: receiver,
                callbackContract: callbackContract,
                uiFeeReceiver: uiFeeReceiver,
                fromMarket: fromMarket,
                toMarket: toMarket,
                minMarketTokens: minMarketTokens,
                executionFee: executionFee,
                callbackGasLimit: callbackGasLimit
            });

        return
            abi.encodeWithSelector(
                IGmx.createShift.selector,
                createShiftParams
            );
    }

    function setNativeToken(address _nativeToken) internal {
        nativeToken = _nativeToken;
    }

    function setOrderVault(
        address _orderVault
    ) external onlyRoles(IRoleManager(roleManagerAddress).ACTIONS_ADMIN()) {
        orderVault = _orderVault;
    }

    function setDepositVault(
        address _depositVault
    ) external onlyRoles(IRoleManager(roleManagerAddress).ACTIONS_ADMIN()) {
        depositVault = _depositVault;
    }

    function setGmxProtocolVault(
        address _gmxProtocolVault
    ) external onlyRoles(IRoleManager(roleManagerAddress).ACTIONS_ADMIN()) {
        gmxProtocol = _gmxProtocolVault;
    }

    function setRoleManagerAddress(
        address _roleManagerAddress
    ) external onlyRoles(IRoleManager(roleManagerAddress).ACTIONS_ADMIN()) {
        roleManagerAddress = _roleManagerAddress;
    }

    function stringToBytes32(
        string memory source
    ) public pure returns (bytes32 result) {
        bytes memory tempEmptyStringTest = bytes(source);
        if (tempEmptyStringTest.length == 0) {
            return 0x0;
        }

        assembly {
            result := mload(add(source, 32))
        }
    }
}

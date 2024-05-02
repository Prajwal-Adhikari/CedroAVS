// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.24;

import "@openzeppelin/contracts/utils/math/SafeMath.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "@openzeppelin/contracts/utils/Context.sol";
import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "./interfaces/IOracle.sol";
import "./interfaces/IOperator.sol";
import "./interfaces/IOperatorManager.sol";
import "./libraries/Data.sol";

string constant TAG = "AIAgent";

contract AIAgent is Context, ReentrancyGuard {
    using ECDSA for bytes32;
    using SafeMath for uint256;
    using SafeERC20 for IERC20;

    uint256 public constant ROOT_DECIMALS = 18;

    uint256 public balanceRequestId;
    uint256 public confirmationThreshold;
    mapping(uint256 => uint256) public confirmatonCount;
    mapping(address => mapping(uint256 => bool)) public signed;
    address public oracle;
    address public owner;
    mapping(uint256 => Data.UserDepositData) public userDeposits;
    address public operator;
    address public operatorManager;
    uint256 public txnIdCounter;
    mapping(bytes32 => bool) public isAssetAccepted;
    bytes32[] public acceptedAssetList;

    event getAssetsData(uint256 balacnceRequestId, Data.AssetsData[] assetData, uint256 totalLiquidityUSD);
    error IncorrectRequestId(
        string contractTAG,
        uint256 sentId,
        uint256 requiredId
    );

    error InsufficientValue(
        string contractTAG,
        uint256 valueSent,
        uint256 valueToBeGreater
    );

    error ArgumentsMismatch(string contractTAG);
    error AssetNotEnabledByOperator(bytes32 asset);
    error AssetAlreadyExists(bytes32 asset);
    error AssetNotEnabledByAgent(bytes32 asset);
    error SignatureCorrupted(bytes signature);
    error AlreadySigned(address operator, uint256 requestID);

    constructor(
        address _operatorManager,
        address _operator,
        address _owner,
        address _oracle
    ) {
        operatorManager = _operatorManager;
        operator = _operator;
        owner = _owner;
        oracle = _oracle;
    }

    function approveAsset(bytes32[] memory _assets) external {
        uint256 len = _assets.length;
        for (uint256 i; i < len; i++) {
            if (!IOperator(operator).isAssetAccepted(_assets[i]))
                revert AssetNotEnabledByOperator(_assets[i]);
            if (isAssetAccepted[_assets[i]])
                revert AssetAlreadyExists(_assets[i]);

            acceptedAssetList.push(_assets[i]);
            isAssetAccepted[_assets[i]] = true;
        }
    }

    function delistAsset(bytes32 _asset) external {
        uint256 len = acceptedAssetList.length;
        for (uint256 i; i < len; i++) {
            if (acceptedAssetList[i] == _asset) {
                acceptedAssetList[i] = acceptedAssetList[len - 1];
                acceptedAssetList.pop();
                isAssetAccepted[_asset] = false;
                break;
            }
        }
    }

    /**
     * @notice Accept deposits from user
     * This function can only be accessed by the account with role Agent.
     * @param _asset - id of the asset
     * @param _amount - amount of the asset
     * @param _timeInterval - time for which user wants to deposit
     * @param _interestRate - rate provided by AI agent
     **/
    function depositAsset(
        bytes32 _asset,
        uint256 _amount,
        uint256 _timeInterval,
        uint256 _interestRate,
        uint256 _timeThreshold,
        bytes memory signature
    ) external payable {
        if (!isAssetAccepted[_asset]) revert AssetNotEnabledByAgent(_asset);

        if (
            owner !=
            isDepositorValid(
                _asset,
                _amount,
                _timeInterval,
                _interestRate,
                _timeThreshold,
                signature
            )
        ) revert SignatureCorrupted(signature);

        Data.UserDepositData storage userData = userDeposits[txnIdCounter];
        userData.txnID = txnIdCounter;
        userData.asset = _asset;
        userData.amount = _amount;
        userData.principalAmount = _amount
            .mul(IOracle(oracle).getSinglePrice(_asset))
            .div(1e18);
        userData.timeInterval = _timeInterval;
        userData.interestRate = _interestRate;
        userData.createdAt = block.timestamp;
        Data.PoolInfo memory pool = IOperatorManager(operatorManager).pools(
            _asset
        );
        _pay(pool.token, _amount);
        txnIdCounter++;
    }

    function _pay(address token, uint256 amount) internal {
        if (token == address(0)) {
            if (msg.value < amount)
                revert InsufficientValue(TAG, msg.value, amount - 1);
        } else {
            IERC20(token).safeTransferFrom(_msgSender(), address(this), amount);
        }
    }

    function isDepositorValid(
        bytes32 _asset,
        uint256 _amount,
        uint256 _timeInterval,
        uint256 _interestRate,
        uint256 _timeThreshold,
        bytes memory signature
    ) internal view returns (address) {
        // Build the hash and check the sig
        // We only accept sigs from the system
        bytes32 msgHash = keccak256(
            abi.encodePacked(
                _msgSender(),
                _asset,
                _amount,
                _timeInterval,
                _interestRate,
                _timeThreshold
            )
        );
        bytes32 signedHash = keccak256(
            abi.encodePacked("\x19Ethereum Signed Message:\n32", msgHash)
        );
        return signedHash.recover(signature);
    }

    //TODO: to be called by whilelisted operator
    function getBalanceData(uint256 requestId)
        external
        returns (Data.AssetsData[] memory assetData)
    {   
        uint256 totalLiquidityUSD;
        uint256 len = acceptedAssetList.length;
        assetData = new Data.AssetsData[](len);
        confirmatonCount[requestId]++;
        if (signed[_msgSender()][requestId])
            revert AlreadySigned(_msgSender(), requestId);
        if (confirmatonCount[requestId] < confirmationThreshold)
            return assetData;
        if (requestId != balanceRequestId)
            revert IncorrectRequestId(TAG, requestId, balanceRequestId);

        for (uint256 i; i < len; i++) {
            assetData[i].assetId = acceptedAssetList[i];

            Data.PoolInfo memory pool = IOperatorManager(operatorManager).pools(
                acceptedAssetList[i]
            );
            assetData[i].assetAddress = pool.token;
            if (pool.token == address(0)) {
                assetData[i].balance = address(this).balance;
            } else {
                assetData[i].balance = IERC20(pool.token).balanceOf(address(this));
            }
            totalLiquidityUSD += _buildAmount(assetData[i].balance, pool.token)
                .mul(IOracle(oracle).getSinglePrice(acceptedAssetList[i]))
                .div(1e18);
        }

        signed[_msgSender()][requestId] = true;
        balanceRequestId++;

        emit getAssetsData(balanceRequestId, assetData, totalLiquidityUSD);
    }

    function getAssetBalance(bytes32 assetId)
        public
        view
        returns (uint256 assetBalance)
    {
        Data.PoolInfo memory pool = IOperatorManager(operatorManager).pools(
            assetId
        );
        assetBalance = IERC20(pool.token).balanceOf(address(this));
    }

    function _buildAmount(
        uint256 _localAmount,
        address _token
    ) internal view returns (uint256 amount) {
        uint8 decimals = _getTokenDecimal(_token);

        if (decimals > ROOT_DECIMALS) {
            amount = _localAmount / (10 ** (decimals - ROOT_DECIMALS));
        } else {
            amount = _localAmount * (10 ** (ROOT_DECIMALS - decimals));
        }
    }

    function _getTokenDecimal(
        address _token
    ) internal view returns (uint8 decimals) {
        decimals = 18;
        if (_token != address(0)) {
            (, bytes memory data) = _token.staticcall(
                abi.encodeWithSignature("decimals()")
            );
            decimals = abi.decode(data, (uint8));
        }
    }

}

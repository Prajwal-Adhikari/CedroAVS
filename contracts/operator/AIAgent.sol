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

    address public oracle;
    address public owner;
    mapping(uint256 => Data.UserDepositData) public userDeposits; // txnID => userDeposits
    address public operator;
    address public operatorManager;
    uint256 public txnIdCounter;
    mapping(bytes32 => bool) public isAssetAccepted; // assetID => isAccepted
    bytes32[] public acceptedAssetList;

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
    error StaleMessage(bytes signature);

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

    /**
     * @notice Approves list of asset from Operator contract
     * This function can only be accessed by _owner
     * @param _assets - list of ids of the asset
     **/
    function approveAsset(bytes32[] memory _assets) external {
        uint256 len = _assets.length;
        for (uint256 i; i < len; ) {
            if (!IOperator(operator).isAssetAccepted(_assets[i]))
                revert AssetNotEnabledByOperator(_assets[i]);
            if (isAssetAccepted[_assets[i]])
                revert AssetAlreadyExists(_assets[i]);

            acceptedAssetList.push(_assets[i]);
            isAssetAccepted[_assets[i]] = true;
            unchecked {
                i++;
            }
        }
    }

    /**
     * @notice Removes asset from the approved list
     * This function can only be accessed by the account with role AIAGENT.
     * @param _asset - id of the asset
     **/
    function delistAsset(bytes32 _asset) external {
        uint256 len = acceptedAssetList.length;
        for (uint256 i; i < len; ) {
            if (acceptedAssetList[i] == _asset) {
                acceptedAssetList[i] = acceptedAssetList[len - 1];
                acceptedAssetList.pop();
                isAssetAccepted[_asset] = false;
                break;
            }
            unchecked {
                i++;
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
     * @param _timeThreshold - block time till which signature is valid
     * @param _signature - msg signed by the AI agent
     **/
    function depositAsset(
        bytes32 _asset,
        uint256 _amount,
        uint256 _timeInterval,
        uint256 _interestRate,
        uint256 _timeThreshold,
        bytes memory _signature
    ) external payable {
        if (!isAssetAccepted[_asset]) revert AssetNotEnabledByAgent(_asset);
        if (block.timestamp > _timeThreshold) revert StaleMessage(_signature);
        if (
            owner !=
            isDepositorValid(
                _asset,
                _amount,
                _timeInterval,
                _interestRate,
                _timeThreshold,
                _signature
            )
        ) revert SignatureCorrupted(_signature);

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
        ++txnIdCounter;
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
        bytes memory _signature
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
        return signedHash.recover(_signature);
    }

    /**
     * @notice Returns balance of the asset holded by the contract
     * @param _token - address of the asset
     **/
    function getAssetBalance(address _token)
        public
        view
        returns (uint256 assetBalance)
    {
        if (_token == address(0)) {
            assetBalance = address(this).balance;
        } else {
            assetBalance = IERC20(_token).balanceOf(address(this));
        }
    }

    function _buildAmount(uint256 _localAmount, address _token)
        internal
        view
        returns (uint256 amount)
    {
        uint8 decimals = _getTokenDecimal(_token);

        if (decimals > ROOT_DECIMALS) {
            amount = _localAmount / (10**(decimals - ROOT_DECIMALS));
        } else {
            amount = _localAmount * (10**(ROOT_DECIMALS - decimals));
        }
    }

    function _getTokenDecimal(address _token)
        internal
        view
        returns (uint8 decimals)
    {
        decimals = 18;
        if (_token != address(0)) {
            (, bytes memory data) = _token.staticcall(
                abi.encodeWithSignature("decimals()")
            );
            decimals = abi.decode(data, (uint8));
        }
    }
}

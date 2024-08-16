// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.24;

import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "@openzeppelin/contracts/utils/Context.sol";
import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {IStargate} from "@stargatefinance/stg-evm-v2/src/interfaces/IStargate.sol";
import {MessagingFee, OFTReceipt, SendParam} from "@layerzerolabs/lz-evm-oapp-v2/contracts/oft/interfaces/IOFT.sol";
import "./interfaces/IOracle.sol";
import "./interfaces/IPool.sol";
import "./interfaces/IStargateInterface.sol";
import "./interfaces/IStargateRouter.sol";
import "./interfaces/IStargateEthVault.sol";
import "./interfaces/IAgentManager.sol";
import "../actions/interfaces/IAction.sol";
import "./libraries/Data.sol";

string constant TAG = "AgentPod";

/**
 * @title AgentPod contract
 * @notice Interface provided to each agent where it can interact with users
 **/
contract AgentPod is Context, ReentrancyGuard {
    using ECDSA for bytes32;
    using SafeERC20 for IERC20;

    address public oracle;
    address public owner;
    address public operator;
    address public agentManager;
    address public pool;

    uint256 public constant ROOT_DECIMALS = 18;
    uint256 public constant SECONDS_IN_A_YEAR = 365 * 86400;

    uint256 public txnIdCounter;
    uint256 public totalReward;
    uint256 public rewardAmount;
    bytes32 public chainID;

    mapping(uint256 => Data.UserDepositData) public userDeposits; // txnID => userDeposits
    mapping(address => uint256[]) public userTxnIDs;
    mapping(bytes32 => bool) public isAssetAccepted; // assetID => isAccepted
    mapping(bytes32 => bool) public isAssetDebt;
    mapping(bytes => bool) public signatureUsed; // signature => isUsed
    mapping(address => uint256) public reward;

    bytes32[] public acceptedAssetList;
    bytes32[] public debtAssetList;
    mapping(uint256 => bool) public slashEmitted;
    event TransferAssetToPool(
        uint256 txnId,
        address depositor,
        bytes32[] assets,
        uint256[] amounts,
        Data.STATUS status
    );

    event ProtocolResponse(bytes32 id,bytes payload, bytes response);

    error DepositNotMatured(uint256 _txnID, address _user);
    error ActionReverted(bytes32 id, bytes payload);
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
    error IncorrectRatio(string contractTAG);
    error AssetNotEnabledByAgentManager(bytes32 asset);
    error AssetAlreadyExists(bytes32 asset);
    error AssetNotEnabledByAgent(bytes32 asset);
    error SignatureCorrupted(bytes signature);
    error StaleMessage(bytes signature);
    error SignatureAlreadyUsed(bytes signature);
    error AssetAlreadyClaimed(uint256 _txnID);
    error InsufficientQuantity();

    constructor(
        address _agentManager,
        address _owner,
        address _oracle
    ) {
        agentManager = _agentManager;
        owner = _owner;
        oracle = _oracle;
    }

    receive() external payable {}

    /**
     * @notice Fetch oracle from AgentManger and sets it in the contract
     * This function can only be accessed by the account with role MANAGER.
     **/
    function setOracle() external {
        oracle = IAgentManager(agentManager).oracle();
    }

    /**
     * @notice Configures a pool contract where agent can transfer the matured funds
     * This function can only be accessed by the account with role MANAGER.
     * @param _pool - address of the pool
     **/
    function setPool(address _pool) external {
        pool = _pool;
    }

    /**
     * @notice Configures reward
     * This function can only be accessed by the account with role MANAGER.
     * @param _reward - amount in USD to be given as reward
     **/
    function setReward(uint256 _reward) external {
        rewardAmount = _reward;
    }

    /**
     * @notice Approves list of asset from AgentManager contract
     * This function can only be accessed by the owner
     * @param _assets - list of ids of the asset
     **/
    function approveAsset(bytes32[] memory _assets) external {
        uint256 len = _assets.length;
        for (uint256 i; i < len; ) {
            Data.PoolInfo memory poolInfo = IAgentManager(agentManager).pools(
                _assets[i]
            );
            if (!poolInfo.enable)
                revert AssetNotEnabledByAgentManager(_assets[i]);
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
     * This function can only be accessed by the owner.
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
     * @notice Configures list of debt assets
     * This function can only be accessed by account with role MANAGER.
     * @param _assets - ids of the debt assets
     **/
    function addDebtAsset(bytes32[] memory _assets) external {
        uint256 len = _assets.length;
        for (uint256 i; i < len; ) {
            if (isAssetDebt[_assets[i]]) revert AssetAlreadyExists(_assets[i]);

            debtAssetList.push(_assets[i]);
            isAssetDebt[_assets[i]] = true;
            unchecked {
                i++;
            }
        }
    }

    /**
     * @notice Removes debt asset from the list
     * This function can only be accessed by account with role MANAGER.
     * @param _asset - id of the debt asset
     **/
    function removeDebtAsset(bytes32 _asset) external {
        uint256 len = debtAssetList.length;
        for (uint256 i; i < len; ) {
            if (debtAssetList[i] == _asset) {
                debtAssetList[i] = debtAssetList[len - 1];
                debtAssetList.pop();
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
        if (signatureUsed[_signature]) revert SignatureAlreadyUsed(_signature);
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
        userData.user = _msgSender();
        userData.asset = _asset;
        userData.amount = _amount;

        userData.principalAmount =
            (_amount * (IOracle(oracle).getSinglePrice(_asset))) /
            (1e18);

        userData.timeInterval = _timeInterval;
        userData.interestRate = _interestRate;
        userData.createdAt = block.timestamp;
        userData.status = Data.STATUS.CREATED;

        userTxnIDs[userData.user].push(txnIdCounter);

        Data.PoolInfo memory poolInfo = IAgentManager(agentManager).pools(
            _asset
        );
        ++txnIdCounter;
        signatureUsed[_signature] = true;
        _pay(poolInfo.token, _amount);
        IAgentManager(agentManager).emitDeposit(
            owner,
            userData.txnID,
            userData.user,
            _asset,
            _amount,
            userData.principalAmount,
            _timeInterval,
            _interestRate
        );
    }

    /**
     * @notice  Transfers assets from AgentPod to Pool contract
     * This function can only be accessed by the owner.
     * @param _assets - id of the assets
     * @param _amounts - amount of the assets
     **/
    function transferAssetToPool(
        bytes32[] memory _assets,
        uint256[] memory _amounts
    ) public {
        if (_assets.length != _amounts.length) revert ArgumentsMismatch(TAG);
        for (uint256 i; i < _assets.length; ) {
            Data.PoolInfo memory poolInfo = IAgentManager(agentManager).pools(
                _assets[i]
            );
            IERC20(poolInfo.token).safeTransfer(pool, _amounts[i]);
            unchecked {
                i++;
            }
        }
    }

    /**
     * @notice Faciliates claim to users
     * This function can be accessed by the user with a valid txnId
     * @param _txnId - key of deposit
     * @param _assets - id of the assets to claim
     * @param _ratio - proportion in which assets are wished to received. Sum to be 1.
     **/
    function claim(
        uint256 _txnId,
        bytes32[] memory _assets,
        uint256[] memory _ratio
    ) public {
        if (_assets.length != _ratio.length) revert ArgumentsMismatch(TAG);
        Data.UserDepositData storage userData = userDeposits[_txnId];
        if (userData.status == Data.STATUS.CLAIMED)
            revert AssetAlreadyClaimed(_txnId);
        uint256 ratioTotal;
        uint256[] memory amounts = new uint256[](_assets.length);
        address[] memory tokens = new address[](_assets.length);
        for (uint256 i; i < _ratio.length; ) {
            ratioTotal += _ratio[i];
            unchecked {
                i++;
            }
        }
        if (ratioTotal != 1e18) revert IncorrectRatio(TAG);

        if (
            block.timestamp <
            (userData.createdAt + (userData.timeInterval / (1e18)))
        ) revert DepositNotMatured(_txnId, userData.user);

        uint256 amountToClaim = (userData.principalAmount *
            (((userData.timeInterval / SECONDS_IN_A_YEAR) *
                (userData.interestRate)) /
                (1e18) +
                (1e18))) / (1e18);

        uint256 totalBalance = poolBalanceOfAgentPod() +
            principalBalanceOfAgentPod();

        if (amountToClaim > totalBalance) {
            if (!slashEmitted[_txnId]) {
                IAgentManager(agentManager).emitSlash(
                    owner,
                    amountToClaim - (totalBalance),
                    chainID,
                    _txnId
                );
                slashEmitted[_txnId] = true;
            }
            return;
        }

        if (amountToClaim > poolBalanceOfAgentPod()) {
            for (uint256 i; i < _assets.length; i++) {
                uint256 deviation = (amountToClaim) - (poolBalanceOfAgentPod());
                amounts[i] =
                    (deviation * (_ratio[i])) /
                    (IOracle(oracle).getSinglePrice(_assets[i]));
            }

            transferAssetToPool(_assets, amounts);
            reward[_msgSender()] += rewardAmount;
            totalReward += rewardAmount;
            userData.status = Data.STATUS.MATURED;
            emit TransferAssetToPool(
                _txnId,
                userData.user,
                _assets,
                amounts,
                userData.status
            );
            return;
        }

        for (uint256 i; i < _assets.length; ) {
            amounts[i] =
                ((amountToClaim + reward[userData.user]) * (_ratio[i])) /
                (IOracle(oracle).getSinglePrice(_assets[i]));
            Data.PoolInfo memory poolInfo = IAgentManager(agentManager).pools(
                _assets[i]
            );
            tokens[i] = poolInfo.token;
            unchecked {
                i++;
            }
        }
        reward[userData.user] = 0;
        IPool(pool).trasferAssetToDepositor(userData.user, tokens, amounts);
        userData.status = Data.STATUS.CLAIMED;
        IAgentManager(agentManager).emitClaim(
            owner,
            _txnId,
            userData.user,
            _assets,
            amounts,
            userData.status
        );
    }

    /**
     * @notice Allows agent to perform DeFi action
     * This function can be accessed by the owner
     * @param _payload - fixed bytes pattern to be provided for a DeFi action
     **/
    function takeAction(bytes memory _payload) public {
        bytes32 id;
        address asset;
        uint256 amount;
        bool success;
        bytes memory data;

        assembly {
            id := mload(add(_payload, 32))
            asset := mload(add(_payload, 52))
            amount := mload(add(_payload, 84))
        }

        if (amount != 0 && asset != address(0))
            IERC20(asset).approve(
                IAgentManager(agentManager).protocol(id),
                amount
            );

        (bytes memory payload, uint256 nativeValue, address contractAddress) = IAction(
            IAgentManager(agentManager).payloadFetchFrom(id)
        ).extractPayloads(_payload);
        if (nativeValue > 0) {
            (success, data) = (contractAddress).call{
                value: nativeValue
            }(payload);
        } else {
            (success, data) = (contractAddress).call(
                payload
            );
        }
        if (!success) {
            revert ActionReverted(id, _payload);
        } else {
            emit ProtocolResponse(id, _payload, data);
        }
    }

    /**
     * @notice Returns balance of agent in Pool contract
     **/
    function debtBalanceOfAgentPod() public view returns (uint256 amountInUSD) {
        uint256 len = debtAssetList.length;
        for (uint256 i; i < len; ) {
            Data.PoolInfo memory poolInfo = IAgentManager(agentManager).pools(
                debtAssetList[i]
            );
            amountInUSD +=
                (IERC20(poolInfo.token).balanceOf(pool) *
                    (IOracle(oracle).getSinglePrice(debtAssetList[i]))) /
                (1e18);
            unchecked {
                i++;
            }
        }
    }

    /**
     * @notice Returns balance of agent in Pool contract
     **/
    function poolBalanceOfAgentPod() public view returns (uint256 amountInUSD) {
        uint256 len = acceptedAssetList.length;
        for (uint256 i; i < len; ) {
            Data.PoolInfo memory poolInfo = IAgentManager(agentManager).pools(
                acceptedAssetList[i]
            );
            amountInUSD +=
                (IERC20(poolInfo.token).balanceOf(pool) *
                    (IOracle(oracle).getSinglePrice(acceptedAssetList[i]))) /
                (1e18);
            unchecked {
                i++;
            }
        }
    }

    /**
     * @notice Returns balance of agent in AgentPod contract
     **/
    function principalBalanceOfAgentPod()
        public
        view
        returns (uint256 amountInUSD)
    {
        uint256 len = acceptedAssetList.length;
        for (uint256 i; i < len; ) {
            Data.PoolInfo memory poolInfo = IAgentManager(agentManager).pools(
                acceptedAssetList[i]
            );
            amountInUSD +=
                (IERC20(poolInfo.token).balanceOf(address(this)) *
                    (IOracle(oracle).getSinglePrice(acceptedAssetList[i]))) /
                (1e18);
            unchecked {
                i++;
            }
        }
    }

    /**
     * @notice Returns totalBalance of the agent
     **/
    function getBalanceOfAgent() public view returns (uint256) {
        if (
            principalBalanceOfAgentPod() + poolBalanceOfAgentPod() <
            debtBalanceOfAgentPod()
        ) return 0;
        return
            principalBalanceOfAgentPod() +
            poolBalanceOfAgentPod() -
            debtBalanceOfAgentPod();
    }

    function _pay(address _token, uint256 _amount) internal {
        if (_token == address(0)) {
            if (msg.value < _amount)
                revert InsufficientValue(TAG, msg.value, _amount - 1);
        } else {
            IERC20(_token).safeTransferFrom(
                _msgSender(),
                address(this),
                _amount
            );
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
                IAgentManager(agentManager).chainID(),
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

    /**
     * @notice Returns list of assets accepted for deposits from users
     **/
    function getAcceptedAssetList() public view returns (bytes32[] memory) {
        bytes32[] memory result = new bytes32[](acceptedAssetList.length);
        for (uint256 i = 0; i < acceptedAssetList.length; i++) {
            result[i] = acceptedAssetList[i];
        }
        return result;
    }

    /**
     * @notice Returns list of debt assets
     **/
    function getDebtAssetList() public view returns (bytes32[] memory) {
        bytes32[] memory result = new bytes32[](debtAssetList.length);
        for (uint256 i = 0; i < debtAssetList.length; i++) {
            result[i] = debtAssetList[i];
        }
        return result;
    }

    /**
     * @notice Returns interest accumulated till date for a particular deposit
     **/
    function getInterestAccumulated(uint256 _txnId)
        external
        view
        returns (uint256)
    {
        Data.UserDepositData storage userData = userDeposits[_txnId];
        uint256 timeDiff = (block.timestamp - userData.createdAt) * (1e18);
        if (userData.timeInterval < timeDiff) timeDiff = userData.timeInterval;
        uint256 amountToClaim = (userData.principalAmount *
            (((timeDiff / SECONDS_IN_A_YEAR) * (userData.interestRate)) /
                (1e18))) / (1e18);
        return amountToClaim;
    }

    /**
     * @notice Allows swap of token to agentPod contracts of other chains
     **/
    function sendToken(
        address _stargate,
        bytes32 chainId,
        uint256 _amount
    ) external {
        SendParam memory sendParam;
        MessagingFee memory messagingFee;
        uint256 valueToSend;

        IERC20(IStargate(_stargate).token()).approve(_stargate, _amount);

        (valueToSend, sendParam, messagingFee) = IStargateInterface(
            IAgentManager(agentManager).stargateIntegration()
        ).prepareTakeTaxi(
                _stargate,
                IAgentManager(agentManager).mapIDs(chainId),
                _amount,
                IAgentManager(agentManager).agentPodContracts(owner, chainId)
            );

        IStargate(_stargate).sendToken{value: valueToSend}(
            sendParam,
            messagingFee,
            address(this)
        );
    }

    /**
     * @notice Swaps liquidity across chain
     * This function can only be accessed by account with role SWAP.
     * @param qty - amount to be swapped
     * @param id - id of the pool
     * @param dstChainId - chain id of destination
     * @param srcPoolId - pool id of source token in Stargate
     * @param dstPoolId - pool id of destination token in Stargate
     * @param amountMin - The minimum amount to receive.
     **/
    function swapRequest(
        uint256 qty,
        bytes32 id,
        bytes32 dstChainId,
        uint16 srcPoolId,
        uint16 dstPoolId,
        uint256 amountMin,
        uint256 valueToSend
    ) external {
        if (valueToSend == 0) revert InsufficientValue(TAG, valueToSend, 0);
        if (qty == 0) revert InsufficientQuantity();

        Data.PoolInfo memory poolInfo = IAgentManager(agentManager).pools(id);
        if (!poolInfo.enable) revert AssetNotEnabledByAgentManager(id);

        if (poolInfo.token == address(0)) {
            if (valueToSend <= qty)
                revert InsufficientValue(TAG, valueToSend, qty);
            // wrap the ETH into WETH
            IStargateEthVault(IAgentManager(agentManager).stargateEthVault())
                .deposit{value: qty}();
            IStargateEthVault(IAgentManager(agentManager).stargateEthVault())
                .approve(IAgentManager(agentManager).stargateRouter(), qty);
        } else {
            IERC20(poolInfo.token).approve(
                IAgentManager(agentManager).stargateRouter(),
                qty
            );
        }

        bytes memory destination = abi.encodePacked(
                IAgentManager(agentManager).agentPodContracts(owner, dstChainId)
            );

        IStargateRouter(IAgentManager(agentManager).stargateRouter()).swap{
            value: valueToSend
        }(
            uint16(IAgentManager(agentManager).mapIDsV1(dstChainId)), // the destination chain id
            srcPoolId, // the source Stargate poolId
            dstPoolId, // the destination Stargate poolId
            payable(address(this)), // refund adddress. if msg.sender pays too much gas, return extra eth
            qty, // total tokens to send to destination chain
            amountMin, // min amount allowed out
            IStargateRouter.lzTxObj(0, 0, "0x"), // default lzTxObj
            destination, // destination address, the sgReceive() implementer
            bytes("") // bytes payload
        );
    }
}

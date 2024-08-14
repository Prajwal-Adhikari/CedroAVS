// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.24;

import "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/utils/ContextUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "./interfaces/IOracle.sol";
import "./libraries/Data.sol";

string constant TAG = "AgentManager";

/**
 * @title AgentManager contract
 * @notice Used to enroll and manage agents in the protocol
 **/
contract AgentManager is
    Initializable,
    ContextUpgradeable,
    ReentrancyGuardUpgradeable
{
    address public oracle;
    address public stargateIntegration;
    address public stargateRouter;
    address public stargateEthVault;
    uint256 public slashingTimeThreshold;
    bytes32 public chainID;

    address[] public agents;
    mapping(address => bool) public isAgent;
    mapping(address => address) public agentContract;
    mapping(bytes32 => Data.PoolInfo) public pools;
    mapping(address => uint256) public requestID;
    mapping(address => mapping(uint256 => uint256)) public requestTimestamp;
    mapping(address => mapping(uint256 => bool)) public requestFulfilled;
    mapping(address => uint256) public bondAmount;
    mapping(address => uint256) public deviationPercentage;
    uint256 public bondThreshold;
    mapping(address => mapping(bytes32 => address)) public agentPodContracts;
    mapping(bytes32 => uint32) public mapIDs;
    mapping(bytes32 => uint16) public mapIDsV1;
    mapping(bytes32 => address) public payloadFetchFrom;
    mapping(bytes32 => address) public protocol;

    event DeviationCheck(address agent, uint256 requestID);
    event BondAmount(address agent, uint256 bondAmount);
    event AgentPodCreated(
        address agent,
        bytes32[] chainID,
        address[] agentContract
    );
    event AgentPodChainList(address agent, bytes32[] chainIDs);
    event Deposit(
        address agent,
        uint256 txnID,
        address depositor,
        bytes32 asset,
        uint256 amount,
        uint256 principalAmount,
        uint256 timeInterval,
        uint256 interestRate
    );

    event Claim(
        address agent,
        uint256 txnId,
        address depositor,
        bytes32[] assets,
        uint256[] amounts,
        Data.STATUS status
    );

    event SlashAgent(
        address agent,
        uint256 amountToSlash,
        bytes32 chainID,
        uint256 txnId
    );

    error InsufficientFees();
    error ArgumentsMismatch(string contractTAG);
    error SlashingTimeLimitExceeded(address agent, uint256 requestID);
    error PreviousDeviationInProgress(address agent, uint256 requestID);
    error NotAnAgent(address agent);
    error AIAgentAlreadyExists(address agent);
    error BondAmountLessThanThreshold(address agent, uint256 bondAmount);
    error DeviationAlreadySet(address agent, uint256 deviationPercentage);
    error TransferEthFailed();

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    function initialize() public initializer {
        __ReentrancyGuard_init();
    }

    /**
     * @notice Configures oracle
     * This function can only be accessed by the account with role MANAGER.
     * @param _oracle - address of the oracle
     **/
    function setOracle(address _oracle) external {
        oracle = _oracle;
    }

    /**
     * @notice Configures chain id
     * This function can only be accessed by the account with role MANAGER.
     * @param _chainID - protocol specific key for the chain
     **/
    function setChainID(bytes32 _chainID) external {
        chainID = _chainID;
    }

    /**
     * @notice Maps chain ids of the protocol to LayerZero V2 chain ids
     * This function can only be accessed by the account with role MANAGER.
     * @param _chainIDs - protocol specific keys for the chains
     * @param _chainIds - LayerZero V2 chain ids
     **/
    function mapChainID(bytes32[] memory _chainIDs, uint32[] memory _chainIds)
        external
    {
        if (_chainIDs.length != _chainIds.length) revert ArgumentsMismatch(TAG);
        for (uint256 i; i < _chainIDs.length; ) {
            mapIDs[_chainIDs[i]] = _chainIds[i];
            unchecked {
                i++;
            }
        }
    }

    /**
     * @notice Maps chain ids of the protocol to LayerZero V1 chain ids
     * This function can only be accessed by the account with role MANAGER.
     * @param _chainIDs - protocol specific keys for the chains
     * @param _chainIds - LayerZero V1 chain ids
     **/
    function mapChainIDV1(bytes32[] memory _chainIDs, uint16[] memory _chainIds)
        external
    {
        if (_chainIDs.length != _chainIds.length) revert ArgumentsMismatch(TAG);
        for (uint256 i; i < _chainIDs.length; ) {
            mapIDsV1[_chainIDs[i]] = _chainIds[i];
            unchecked {
                i++;
            }
        }
    }

    /**
     * @notice Configures chain id
     * This function can only be accessed by the account with role MANAGER.
     * @param _chainID - protocol specific key for the chain
     * @param _fetchAddress - address of interaction contract to fetch payload
     **/
    function setPayloadFetchAddress(bytes32 _chainID, address _fetchAddress) external {
        payloadFetchFrom[_chainID] = _fetchAddress;
    }

    /**
     * @notice Configures bond threshold
     * This function can only be accessed by the account with role AGENT.
     * @param _bondThreshold - minimum amount of bond to be deposited
     **/
    function setBondThreshold(uint256 _bondThreshold) external {
        bondThreshold = _bondThreshold;
    }

    /**
     * @notice Configures deviation percentage
     * This function can only be accessed by the account with role AGENT.
     * @param _deviationPercentage - deviation that is allowed for agent before slashing
     **/
    function setDeviationPercentage(uint256 _deviationPercentage) external {
        if (deviationPercentage[_msgSender()] != 0)
            revert DeviationAlreadySet(
                _msgSender(),
                deviationPercentage[_msgSender()]
            );
        deviationPercentage[_msgSender()] = _deviationPercentage;
    }

    /**
     * @notice Configures bond amount for agent
     * This function can only be accessed by the account with role AGENT.
     **/
    function updateBondAmount() external payable {
        bondAmount[_msgSender()] += msg.value;
        if (bondAmount[_msgSender()] < bondThreshold)
            revert BondAmountLessThanThreshold(
                _msgSender(),
                bondAmount[_msgSender()]
            );
        emit BondAmount(_msgSender(), bondAmount[_msgSender()]);
    }

    /**
     * @notice Enrolls agent
     * This function can only be accessed by the account with role Agent.
     * @param _chainIDs - ids of the chain where agent wants to create AgentPod contract
     **/
    function enrollAgent(bytes32[] memory _chainIDs) external payable {
        if (isAgent[_msgSender()]) revert AIAgentAlreadyExists(_msgSender());
        agents.push(_msgSender());
        isAgent[_msgSender()] = true;
        emit AgentPodChainList(_msgSender(), _chainIDs);
    }

    /**
     * @notice Configures agents to their respective AgentPod contracts
     * This function can only be accessed by the account with role Manager.
     * @param _agent - address of the agent.
     * @param _chainIDs - chain keys
     * @param _agentPods - address of AgentPod contracts
     **/
    function mapAgent(
        address _agent,
        bytes32[] memory _chainIDs,
        address[] memory _agentPods
    ) external {
        if (_chainIDs.length != _agentPods.length)
            revert ArgumentsMismatch(TAG);
        for (uint256 i; i < _chainIDs.length; ) {
            agentPodContracts[_agent][_chainIDs[i]] = _agentPods[i];
            unchecked {
                i++;
            }
        }

        emit AgentPodCreated(_agent, _chainIDs, _agentPods);
    }

    /**
     * @notice Removes agent from the agent list
     * This function can only be accessed by the account with role MANAGER.
     * @param _agent - address of the agent
     **/
    function delistAgent(address _agent) external {
        uint256 len = agents.length;
        for (uint256 i; i < len; ) {
            if (agents[i] == _agent) {
                agents[i] = agents[len - 1];
                agents.pop();
                isAgent[_agent] = false;
                break;
            }

            unchecked {
                i++;
            }
        }
    }

    /**
     * @notice Configures pool for an asset
     * This function can only be accessed by account with role MANAGER.
     * @param _id - pool id for the asset.
     * @param _token - address of the asset.
     * @param _enable - status of the pool.
     **/
    function setPools(
        bytes32[] memory _id,
        address[] memory _token,
        bool[] memory _enable
    ) external {
        if (_token.length != _id.length || _token.length != _enable.length)
            revert ArgumentsMismatch(TAG);
        for (uint256 i; i < _token.length; ) {
            Data.PoolInfo storage pool = pools[_id[i]];
            pool.token = _token[i];
            pool.enable = _enable[i];

            unchecked {
                i++;
            }
        }
    }

    /**
     * @notice Slash an agent in case of major deviation in holdings
     * @param _agent - address of the agent.
     * @param _slashAmount - amount to be slashed in USD
     **/
    function slashAgentForDeviation(address _agent, uint256 _slashAmount)
        external
    {
        if (
            block.timestamp >
            requestTimestamp[_agent][requestID[_agent]] +
                (slashingTimeThreshold)
        ) revert SlashingTimeLimitExceeded(_agent, requestID[_agent]);

        _slashAmount =
            (_slashAmount * (1e18)) /
            IOracle(oracle).getSinglePrice(chainID);

        if (_slashAmount > bondAmount[_agent])
            _slashAmount = bondAmount[_agent];
        bondAmount[_agent] -= _slashAmount;
        (bool success, ) = (payable(agentPodContracts[_agent][chainID])).call{
            value: _slashAmount
        }("");

        if (!success) revert TransferEthFailed();
    }

    /**
     * @notice Slash an agent in case of bank run
     * @param _agent - address of the agent.
     * @param _slashAmount - amount to be slashed in USD
     * @param _txnID - id of the txn for which slash is initiated
     * @param _chainID - id of the chain from which slash event is generated
     * @param _agentContract - address to which the fund is to be sent
     **/
    function slashAgentForBankRun(
        address _agent,
        uint256 _slashAmount,
        uint256 _txnID,
        bytes32 _chainID,
        address _agentContract
    ) external {
        // swap interface
    }

    /**
     * @notice Invoked after receiving responses from operator for deviation check
     * This function can only be accessed by CedroServiceManager
     **/
    function receivedResponse(address _agent) external {
        requestFulfilled[_agent][requestID[_agent]] = true;
    }

    /**
     * @notice Invokes operator to check for deviation of an agent
     * @param _agent - address of the agent
     **/
    function invokeDeviationCheck(address _agent) external {
        if (
            !(requestFulfilled[_agent][requestID[_agent]] ||
                requestID[_agent] == 0)
        ) revert PreviousDeviationInProgress(_agent, requestID[_agent]);
        ++requestID[_agent];
        requestTimestamp[_agent][requestID[_agent]] = block.timestamp;
        emit DeviationCheck(_agent, requestID[_agent]);
    }

    function emitDeposit(
        address _agent,
        uint256 _txnID,
        address _depositor,
        bytes32 _asset,
        uint256 _amount,
        uint256 _principalAmount,
        uint256 _timeInterval,
        uint256 _interestRate
    ) external {
        emit Deposit(
            _agent,
            _txnID,
            _depositor,
            _asset,
            _amount,
            _principalAmount,
            _timeInterval,
            _interestRate
        );
    }

    function emitClaim(
        address _agent,
        uint256 _txnId,
        address _depositor,
        bytes32[] memory _assets,
        uint256[] memory _amounts,
        Data.STATUS _status
    ) external {
        emit Claim(_agent, _txnId, _depositor, _assets, _amounts, _status);
    }

    function emitSlash(
        address _agent,
        uint256 _amountToSlash,
        bytes32 _chainID,
        uint256 _txnId
    ) external {
        emit SlashAgent(_agent, _amountToSlash, _chainID, _txnId);
    }
}
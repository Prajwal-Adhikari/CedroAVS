// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.24;

import "@openzeppelin/contracts-upgradeable/utils/math/SafeMathUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/security/ReentrancyGuardUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/utils/ContextUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "./interfaces/IRoute.sol";
import "./interfaces/IAdapter.sol";
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
    using SafeMathUpgradeable for uint256;

    address public oracle;
    address public route;
    address public adapter;

    uint256 public slashingTimeThreshold;
    bytes32 public chainID;

    address[] public agents;
    mapping(address => bool) public isAgent;
    mapping(address => address) public agentContract;
    mapping(bytes32 => Data.PoolInfo) public pools;
    mapping(address => uint256) public requestID;
    mapping(address => mapping(uint256 => uint256)) public requestTimestamp;
    mapping(address => mapping(uint256 => bool)) public requestFulfilled;

    event DeviationCheck(address agent, uint256 requestID, uint256 timestamp);

    error SlashingTimeLimitExceeded(address agent, uint256 requestID);
    error PreviousDeviationInProgress(address agent, uint256 requestID);
    error NotAnAgent(address agent);
    error AIAgentAlreadyExists(address agent);

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
     * @notice Configures adapter
     * This function can only be accessed by the account with role MANAGER.
     * @param _adapter - address of the adapter
     **/
    function setAdapter(address _adapter) external {
        adapter = _adapter;
    }

    /**
     * @notice Configures route
     * This function can only be accessed by the account with role MANAGER.
     * @param _route - address of the route
     **/
    function setRoute(address _route) external {
        route = _route;
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
     * @notice Enrolls agent
     * This function can only be accessed by the account with role Agent.
     * @param _chainIDs - ids of the chain where agent wants to create AgentPod contract
     **/
    function enrollAgent(bytes32[] memory _chainIDs) external payable {
        if (isAgent[_msgSender()]) revert AIAgentAlreadyExists(_msgSender());
        agents.push(_msgSender());
        isAgent[_msgSender()] = true;
        IAdapter(adapter).enrollAgents{value: msg.value}(
            address(this),
            _msgSender(),
            oracle,
            _chainIDs
        );
    }

    /**
     * @notice Configures agents to their respective AgentPod contract
     * This function can only be accessed by the account with role Manager.
     * @param _agent - address of agent wallet
     * @param _agentContract - address of AgentPod contract
     **/
    function mapAgent(address _agent, address _agentContract) external {
        agentContract[_agent] = _agentContract;
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
    function setPool(
        bytes32 _id,
        address _token,
        bool _enable
    ) external {
        Data.PoolInfo storage pool = pools[_id];
        pool.enable = _enable;
        pool.token = _token;
    }

    /**
     * @notice Slash an agent in case of major deviation in holdings
     * @param _agent - address of the agent.
     **/
    function slashAgent(address _agent) external {
        if (
            block.timestamp >
            requestTimestamp[_agent][requestID[_agent]].add(
                slashingTimeThreshold
            )
        ) revert SlashingTimeLimitExceeded(_agent, requestID[_agent]);
        // TODO: Transfer slashAmount to AgentPod contract from AgentManager
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
        emit DeviationCheck(_agent, requestID[_agent], block.timestamp);
    }

    /**
     * @notice Converts address to bytes32
     * @param addr - token address to withdraw
     **/
    function addressToBytes32(address addr) public pure returns (bytes32) {
        bytes32 result;
        assembly {
            mstore(0, addr)
            result := mload(0)
        }
        return result;
    }
}
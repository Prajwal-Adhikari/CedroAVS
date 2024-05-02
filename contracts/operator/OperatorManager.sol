// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.24;

import "@openzeppelin/contracts-upgradeable/utils/math/SafeMathUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/security/ReentrancyGuardUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/utils/ContextUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";

import "./libraries/Data.sol";
import {Operator} from "./Operator.sol";

string constant TAG = "OperatorManager";

contract OperatorManager is
    Initializable,
    ContextUpgradeable,
    ReentrancyGuardUpgradeable
{
    using SafeMathUpgradeable for uint256;

    address[] public operators;
    mapping(address => bool) public isOperator;
    mapping(address => Operator) public operatorContract;
    mapping(bytes32 => Data.PoolInfo) public pools;
    mapping(address => Data.SnapshotInfo) public snapshots;
    mapping(address => address) public agentToOperator;

    address public oracle;
    uint256 public slashingTimeThreshold;

    event UpdateAgentBalance(
        address agent,
        uint256 requestID,
        uint256 timestamp
    );

    error PreviousSnapshotInProgress(address agent, uint256 requestID);
    error SlashingTimeLimitExceeded(address agent, uint256 requestID);
    error NotAnAgent(address agent);
    error OperatorAlreadyExists(address operator);

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    function initialize() public initializer {
        __ReentrancyGuard_init();
    }

    function setOracle(address _oracle) external {
        oracle = _oracle;
    }

    /**
     * @notice Enrolls operator by service manager
     * This function can only be accessed by the account with role OPERATOR.
     **/
    function enrollOperator() external {
        if (isOperator[msg.sender]) revert OperatorAlreadyExists(msg.sender);
        operators.push(msg.sender);
        isOperator[msg.sender] = true;
        Operator OperatorInstance = new Operator(
            address(this),
            msg.sender,
            oracle
        );
        operatorContract[msg.sender] = OperatorInstance;
    }

    /**
     * @notice Removes agent from the operator
     * This function can only be accessed by the account with role MANAGER.
     * @param _operator - address of the agent
     **/
    function delistOperator(address _operator) external {
        uint256 len = operators.length;
        for (uint256 i; i < len; i++) {
            if (operators[i] == _operator) {
                operators[i] = operators[len - 1];
                operators.pop();
                isOperator[_operator] = false;
                break;
            }
        }
    }

    /**
     * @notice Configures pool for an asset
     * This function can only be accessed by account with role MANAGER.
     * @param id - pool id for the asset.
     * @param token - address of the asset.
     * @param enable - status of the pool
     **/
    function setPool(
        bytes32 id,
        address token,
        bool enable
    ) external {
        Data.PoolInfo storage pool = pools[id];
        pool.enable = enable;
        pool.token = token;
    }

    // To be called by only operator contract
    function mapAgent(address _agent, address _operatorContract) external {
        agentToOperator[_agent] = _operatorContract;
    }

    /**
     * @notice Triggers a snapshot update event for the agent
     * This function can only be accessed by any user.
     **/
    function updateSnapshot(address _agent) external payable {
        // TODO: take some bond amount to trigger update event or whitelist callers
        Data.SnapshotInfo storage snapshot = snapshots[_agent];
        if (snapshot.inProgress) {
            revert PreviousSnapshotInProgress(_agent, snapshot.requestID);
        }
        snapshot.inProgress = true;
        snapshot.snapshotRequestTime = block.timestamp;
        snapshot.requestID++;

        emit UpdateAgentBalance(
            _agent,
            snapshot.requestID,
            snapshot.snapshotRequestTime
        );
    }

    // TODO: function claimSlashBond(address slasher) external

    /**
     * @notice Slash an agent in case of major deviation in holdings
     * This function can only be accessed by Operator.
     **/
    function slashAgent(address _agent, uint256 _slashAmount) external {
        Data.SnapshotInfo memory snapshot = snapshots[_agent];
        if (
            block.timestamp >
            snapshot.snapshotRequestTime.add(slashingTimeThreshold)
        ) revert SlashingTimeLimitExceeded(_agent, snapshot.requestID);
        // TODO: Transfer slashAmount to AIAgent contract from operator
    }
}

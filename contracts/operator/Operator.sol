// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.24;

import "@openzeppelin/contracts/utils/math/SafeMath.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "./libraries/Data.sol";
import "./interfaces/IServiceManager.sol";
import {AIAgent} from "./AIAgent.sol";

string constant TAG = "Operator";

contract Operator {
    using SafeMath for uint256;
    address[] public agents;
    mapping(address => bool) public isAgent;
    mapping(address => AIAgent) public agentContract;
    mapping(address => uint256) public agentBond;
    mapping(bytes32 => bool) public isAssetAccepted;
    bytes32[] public acceptedAssetList;
    address public oracle;
    address public owner;
    address public serviceManager;

    error NotAnAgent(address agent);
    error AgentAlreadyExists(address agent);
    error AssetNotEnabledByServiceManager(bytes32 asset);
    error AssetAlreadyExists(bytes32 asset);

    constructor(address _serviceManager, address _owner, address _oracle) {
        serviceManager = _serviceManager;
        owner = _owner;
        oracle = _oracle;
    }

    function setOracle(address _oracle) external {
        oracle = _oracle;
    }

    /**
     * @notice Enrolls agent to the operator
     * This function can only be accessed by the account with role AI_AGENT.
     **/
    function enrollAgent() external {
        if (isAgent[msg.sender]) revert AgentAlreadyExists(msg.sender);
        agents.push(msg.sender);
        isAgent[msg.sender] = true;
        AIAgent AIAgentInstance = new AIAgent(serviceManager, address(this), msg.sender, oracle);
        agentContract[msg.sender] = AIAgentInstance;
        IServiceManager(serviceManager).mapAgent(msg.sender,address(this));
    }

    /**
     * @notice Removes agent from the operator
     * This function can only be accessed by the account with role OPERATOR.
     * @param _agent - address of the agent
     **/
    function delistAgent(address _agent) external {
        uint256 len = agents.length;
        for (uint256 i; i < len; i++) {
            if (agents[i] == _agent) {
                agents[i] = agents[len - 1];
                agents.pop();
                isAgent[_agent] = false;
                break;
            }
        }
    }

    function approveAsset(bytes32[] memory _assets) external {
        uint256 len = _assets.length;
        for (uint256 i; i < len; i++) {

        Data.PoolInfo memory pool = IServiceManager(serviceManager).pools(_assets[i]);
        if (! pool.enable)
            revert AssetNotEnabledByServiceManager(_assets[i]);
        if(isAssetAccepted[_assets[i]])
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
     * @notice Updates bond amount by agent
     * This function can only be accessed by the account with role Agent.
     **/
    function updateBondAmount() external payable {
        if (!isAgent[msg.sender]) revert NotAnAgent(msg.sender);
        agentBond[msg.sender] += msg.value;
    }


}

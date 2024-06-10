// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.24;
import "../libraries/Data.sol";
import {AgentPod} from ".././AgentPod.sol";

interface IAgentManager {
    /**
     * @notice Configures oracle
     * This function can only be accessed by the account with role MANAGER.
     * @param _oracle - address of the oracle
     **/
    function setOracle(address _oracle) external;

    /**
     * @notice Configures chain id
     * This function can only be accessed by the account with role MANAGER.
     * @param _chainID - protocol specific key for the chain
     **/
    function setChainID(bytes32 _chainID) external;

    /**
     * @notice Enrolls operator by service manager
     * This function can only be accessed by the account with role OPERATOR.
     **/
    function enrollOperator() external;

    /**
     * @notice Removes agent from the operator
     * This function can only be accessed by the account with role MANAGER.
     * @param _operator - address of the agent
     **/
    function delistOperator(address _operator) external;

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
    ) external;

    /**
     * @notice Returns the pool info of the asset
     * @param _assetId - id of the asset.
     **/
    function pools(bytes32 _assetId)
        external
        view
        returns (Data.PoolInfo memory);

    /**
     * @notice Returns the address of the oracle
     **/
    function oracle() external view returns (address);

    /**
     * @notice Returns the protocol specific key of the chain
     **/
    function chainID() external view returns (bytes32);

    /**
     * @notice Maps the agent to its corresponding operator
     * This function can only be accessed by account with role OPERATOR.
     * @param _agent - address of the agent.
     * @param _agentPod - address of the agentPod contract.
     **/
    function mapAgent(address _agent, address _agentPod) external;

    /**
     * @notice Triggers a snapshot update event for the agent.
     * This function can only be accessed by any user.
     * @param _agent - address of the agent.
     **/
    function updateSnapshot(address _agent) external;

    /**
     * @notice Slash an agent in case of major deviation in holdings
     * This function can only be accessed by Operator.
     * @param _agent - address of the agent.
     * @param _agent - amount to be slashed
     **/
    function slashAgent(address _agent, uint256 _slashAmount) external;
}

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

    function deviationPercentage(address _agent)
        external
        view
        returns (uint256);

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
    ) external;

    function emitDeposit(
        address agent,
        uint256 txnID,
        address depositor,
        bytes32 asset,
        uint256 amount,
        uint256 principalAmount,
        uint256 timeInterval,
        uint256 interestRate
    ) external;

    function emitClaim(
        address agent,
        uint256 txnId,
        address depositor,
        bytes32[] memory assets,
        uint256[] memory amounts,
        Data.STATUS status
    ) external;

    function emitSlash(
        address _agent,
        uint256 _amountToSlash,
        bytes32 _chainID,
        uint256 _txnId
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
    ) external;

    function mapIDs(bytes32 _chainID) external view returns (uint32);

    function agentPodContracts(address _agent, bytes32 _chainID)
        external
        view
        returns (address);

    function stargateIntegration() external view returns (address);

    function stargateRouter() external view returns (address);

    function stargateEthVault() external view returns (address);

    /**
     * @notice Invoked after receiving responses from operator for deviation check
     * This function can only be accessed by CedroServiceManager
     **/
    function receivedResponse(address _agent) external;

    /**
     * @notice Slash an agent in case of major deviation in holdings
     * This function can only be accessed by Operator.
     * @param _agent - address of the agent.
     * @param _agent - amount to be slashed
     **/
    function slashAgentForDeviation(address _agent, uint256 _slashAmount)
        external;

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
    ) external;
}

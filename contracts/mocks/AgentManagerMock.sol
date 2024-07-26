// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.24;

import "@openzeppelin/contracts-upgradeable/utils/ContextUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";

import "hardhat/console.sol";

string constant TAG = "AgentManager";

/**
 * @title AgentManager contract
 * @notice Used to enroll and manage agents in the protocol
 **/
contract AgentManagerMock is Initializable, ContextUpgradeable {
    mapping(address => uint256) deviationPercentages;

    function setDeviationPercentage(
        address agentAddress,
        uint256 _deviationPercentage
    ) external {
        deviationPercentages[agentAddress] = _deviationPercentage;
    }

    function slashAgent(
        address agentAddress,
        uint256 slashingAmount
    ) external pure {
        console.log(" The slashed agent", agentAddress, slashingAmount);
    }

    function slashAgentForBankRun(address _agent, uint256 _slashAmount, uint256 _txnId, bytes32 _chainID, address _agentContract) external pure {
        console.log(" The slashed agent", _agent, _slashAmount,_txnId);
    }

    function receivedResponse(address agentAddress) external pure {
        console.log("Agent Responsed", agentAddress);
    }

    function deviationPercentage(
        address agentAddress
    ) external view returns (uint256) {
        return deviationPercentages[agentAddress];
    }
}

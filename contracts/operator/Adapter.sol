// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.24;

import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "@openzeppelin/contracts/utils/Context.sol";
import "../operator/interfaces/IRoute.sol";
import "../operator/interfaces/IAgentManager.sol";
import {AgentPod} from "./AgentPod.sol";

string constant TAG = "Adapter";

contract Adapter is Context, ReentrancyGuard {
    address public route;

    /**
     * @notice Configures route
     * This function can only be accessed by the account with role MANAGER.
     * @param _route - address of the route
     **/
    function setRoute(address _route) external {
        route = _route;
    }

    function enrollAgents(
        address agentManager,
        address agent,
        address oracle,
        bytes32[] memory chainIDs
    ) public payable {
        AgentPod AIInstance = new AgentPod(agentManager, agent, oracle);
        IAgentManager(agentManager).mapAgent(agent, address(AIInstance));

        if (chainIDs.length != 0) {
            bytes memory payload = abi.encodePacked(addressToBytes32(agent));
            IRoute(route).send{value: msg.value}(payload, chainIDs, agent);
        }
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

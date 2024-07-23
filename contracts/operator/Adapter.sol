// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.24;

import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/utils/ContextUpgradeable.sol";
import "../operator/interfaces/IRoute.sol";
import "../operator/interfaces/IAgentManager.sol";

import {AgentPod} from "./AgentPod.sol";

string constant TAG = "Adapter";

contract Adapter is Initializable, ContextUpgradeable {
    address public route;

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    function initialize() public initializer {}

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
        bytes32[] memory chainIDs,
        uint256[] memory fees
    ) public payable {
        AgentPod AIInstance = new AgentPod(agentManager, agent, oracle);
        IAgentManager(agentManager).mapAgent(agent, address(AIInstance));

        bytes memory payload = abi.encodePacked(addressToBytes32(agent));
        IRoute(route).send{value: msg.value}(payload, chainIDs, fees, agent);
    }

    function enrollAgent(
        address agentManager,
        address agent,
        address oracle
    ) public payable {
        AgentPod AIInstance = new AgentPod(agentManager, agent, oracle);
        IAgentManager(agentManager).mapAgent(agent, address(AIInstance));
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

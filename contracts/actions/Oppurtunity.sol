// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.24;

import "@chainlink/contracts/src/v0.8/interfaces/AggregatorV3Interface.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/utils/ContextUpgradeable.sol";
import "./interfaces/IAction.sol";

string constant TAG = "Oppurtunity";

interface IERC20 {
    function decimals() external view returns (uint8);
}

/**
 * @title Oppurtunity contract
 * @notice Used to transfer payload to respective defi oppurtunities
 **/
contract Oppurtunity is Initializable, ContextUpgradeable {
    mapping(bytes32 => address) public protocolInterface;
    mapping(bytes32 => address) public protocol;

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    /**
     * @notice Initilize the contract.
     **/
    function initialize() public initializer {}

    function setProtocol(bytes32 id, address defiAddress) public {
        protocol[id] = defiAddress;
    }

    function setProtocolInterface(bytes32 id, address defiInterface) public {
        protocolInterface[id] = defiInterface;
    }

    function transferPayload(bytes memory _payload, address _to) public {
        bytes32 id;
        assembly {
            id := mload(add(_payload, 34))
        }
        IAction(protocolInterface[id]).extractPayload(
            _payload,
            _to,
            protocol[id]
        );
    }
}

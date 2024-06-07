// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.24;
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {LzApp} from "../layerZero/lzApp/LzApp.sol";
import "../interfaces/IAgentPod.sol";
import "../interfaces/IAgentManager.sol";
import "../interfaces/IAdapter.sol";
import {Adapter} from "../Adapter.sol";

string constant TAG = "LzRoute";

/**
 * Title LzRoute
 * Uses LayerZero functions for sending and receiving messages
 **/

contract LzRoute is Initializable, LzApp {
    address public agentManager;
    address public oracle;
    address public adapter;

    mapping(bytes32 => uint16) public chainID;
    mapping(uint16 => bytes32) public chainKey;

    uint256 public gasUsage;
    bytes32 public route;

    error InsufficientGas(string contractTAG);

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    function initialize(
        address _lzEndpoint,
        uint16 _localChainId,
        uint256 _gasUsage
    ) public initializer {
        __LzApp_init(_lzEndpoint, _localChainId);
        if (_gasUsage < 200000) revert InsufficientGas(TAG);
        gasUsage = _gasUsage;
    }

    function setChainKeyID(bytes32 chainId, uint16 id) external {
        chainID[chainId] = id;
        chainKey[id] = chainId;
    }

    function setRoute(bytes32 routeKey) external {
        route = routeKey;
    }

    function setAdapter(address _adapter) external {
        adapter = _adapter;
    }

    /**
     * @dev Receive a cross chain message sent from the branch of chain with chain id `_srcChainId`.
     *
     **/
    function _blockingLzReceive(
        uint16 _srcChainId,
        bytes memory,
        uint64,
        bytes memory _payload
    ) internal override {
        bytes32 id;
        assembly {
            id := mload(add(_payload, 32))
        }
        address agent = _bytesToAddress(id);
        bytes32[] memory chains;
        IAdapter(adapter).enrollAgents(
            agentManager,
            agent,
            agentManager,
            chains
        );
    }

    function send(
        bytes memory payload,
        bytes32[] memory chainIDs,
        address sender
    ) external payable {
        uint256 len = chainIDs.length;
        uint256 value = msg.value; //TO DO: add respective value to diff chains
        for (uint256 i; i < len; ) {
            _lzSend(
                chainID[chainIDs[i]],
                payload,
                payable(sender),
                address(this),
                abi.encodePacked(uint16(1), gasUsage),
                value
            );
            unchecked {
                i++;
            }
        }
    }

    function estimateFees(bytes32 chainId, uint256 airdropAmount)
        public
        view
        returns (uint256)
    {
        bytes memory payload = abi.encode("", chainID[chainId], msg.sender, 0);
        uint256 fees;
        if (airdropAmount > 0) {
            (fees, ) = lzEndpoint.estimateFees(
                chainID[chainId],
                address(this),
                payload,
                false,
                abi.encodePacked(uint16(2), gasUsage, airdropAmount, address(0))
            );
        } else {
            (fees, ) = lzEndpoint.estimateFees(
                chainID[chainId],
                address(this),
                payload,
                false,
                abi.encodePacked(uint16(1), gasUsage)
            );
        }

        return fees;
    }

    /**
     * @notice Converts bytes32 to address
     * @param bys - The bytes32 data
     **/
    function _bytesToAddress(bytes32 bys) private pure returns (address addr) {
        return address(uint160(uint256(bys)));
    }
}

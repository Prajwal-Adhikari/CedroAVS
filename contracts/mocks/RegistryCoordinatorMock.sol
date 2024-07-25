// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.12;

import "@openzeppelin/contracts-upgradeable/utils/ContextUpgradeable.sol";

import {RegistryCoordinator} from "../../contracts/core/eigenContracts/RegistryCoordinator.sol";
import {BN254} from "../../contracts/core/eigenContracts/libraries/BN254.sol";
import {IPauserRegistry} from "../../contracts/core/eigenContracts/interfaces/IPauserRegistry.sol";
import {ISignatureUtils} from "../../contracts/core/eigenContracts/interfaces/ISignatureUtils.sol";
import {IBLSApkRegistry} from "../../contracts/core/eigenContracts/interfaces/IBLSApkRegistry.sol";
import {IStakeRegistry} from "../../contracts/core/eigenContracts/interfaces/IStakeRegistry.sol";
import {IIndexRegistry} from "../../contracts/core/eigenContracts/interfaces/IIndexRegistry.sol";
import {IServiceManager} from "../../contracts/core/eigenContracts/interfaces/IServiceManager.sol";
import {IRegistryCoordinator} from "../../contracts/core/eigenContracts/interfaces/IRegistryCoordinator.sol";

import "hardhat/console.sol";

// wrapper around the RegistryCoordinator contract that exposes required functionality, for unit testing
contract RegistryCoordinatorMock is ContextUpgradeable, RegistryCoordinator {
    using BN254 for BN254.G1Point;

    constructor() {}

    // constructor(address ownerAddress) BLSApkRegistry(ownerAddress) {
    // }

    function RegistryCoordinatorMockInit(
        address ownerAddress,
        address _churnApprover,
        address _ejector,
        IServiceManager _serviceManager,
        IStakeRegistry _stakeRegistry,
        IBLSApkRegistry _blsApkRegistry,
        IIndexRegistry _indexRegistry
    ) external initializer {
        _RegistryCoordinatorStorage_init(
            ownerAddress,
            _serviceManager,
            _stakeRegistry,
            _blsApkRegistry,
            _indexRegistry
        );
        __EIP712_init("AVSRegistryCoordinator", "v0.0.1");
        _setChurnApprover(_churnApprover);
        _setEjector(_ejector);
    }

    function registerOperatorMock(
        address operatorAddress,
        bytes calldata quorumNumbers,
        string calldata socket,
        IBLSApkRegistry.PubkeyRegistrationParams calldata params,
        SignatureWithSaltAndExpiry memory operatorSignature
    ) external  {
        bytes32 operatorId = _getOrCreateOperatorId(operatorAddress, params);
        // Register the operator in each of the registry contracts and update the operator's
        // quorum bitmap and registration status
        uint32[] memory numOperatorsPerQuorum = _registerOperator({
            operator: operatorAddress,
            operatorId: operatorId,
            quorumNumbers: quorumNumbers,
            operatorSignature: operatorSignature
        }).numOperatorsPerQuorum;

        // For each quorum, validate that the new operator count does not exceed the maximum
        // (If it does, an operator needs to be replaced -- see `registerOperatorWithChurn`)
        for (uint256 i = 0; i < quorumNumbers.length; i++) {
            uint8 quorumNumber = uint8(quorumNumbers[i]);

            require(
                numOperatorsPerQuorum[i] <=
                    _quorumParams[quorumNumber].maxOperatorCount,
                "RegistryCoordinator.registerOperator: operator count exceeds maximum"
            );
        }
    }
}

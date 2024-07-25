// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.12;

import "@openzeppelin/contracts-upgradeable/utils/ContextUpgradeable.sol";

import {BLSApkRegistry} from "../../contracts/core/eigenContracts/BLSApkRegistry.sol";
import {BN254} from "../../contracts/core/eigenContracts/libraries/BN254.sol";


// wrapper around the BLSApkRegistry contract that exposes internal functionality, for unit testing _other functionality_.
contract BLSApkRegistryMock is ContextUpgradeable, BLSApkRegistry {
    using BN254 for BN254.G1Point;

    constructor() {}

    // constructor(address ownerAddress) BLSApkRegistry(ownerAddress) {
    // }

    function BLSApkRegistryMockinitialize(address ownerAddress) external  initializer {
        _BLSApkRegistryStorage_init(ownerAddress);
    }

    function setBLSPublicKey(
        address account,
        BN254.G1Point memory pk
    ) external {
        bytes32 pubkeyHash = BN254.hashG1Point(pk);
        // store updates
        operatorToPubkeyHash[account] = pubkeyHash;
        pubkeyHashToOperator[pubkeyHash] = account;
        operatorToPubkey[account] = pk;
    }
}

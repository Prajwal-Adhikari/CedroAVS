// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.12;

import {BLSApkRegistry} from "../../contracts/core/eigenContracts/BLSApkRegistry.sol";
import {BN254} from "../../contracts/core/eigenContracts/libraries/BN254.sol";

// wrapper around the BLSApkRegistry contract that exposes internal functionality, for unit testing _other functionality_.
contract BLSApkRegistryMock is BLSApkRegistry {
    using BN254 for BN254.G1Point;

    constructor(address OwnerAddress) BLSApkRegistry(OwnerAddress) {}

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

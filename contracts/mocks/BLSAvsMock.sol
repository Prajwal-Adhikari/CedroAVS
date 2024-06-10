// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.12;

import {BLSSignatureChecker} from "../../contracts/core/eigenContracts/BLSSignatureChecker.sol";
import {BLSApkRegistry} from "../../contracts/core/eigenContracts/BLSApkRegistry.sol";
import {BN254} from "../../contracts/core/eigenContracts/libraries/BN254.sol";
import {IRegistryCoordinator} from "../../contracts/core/eigenContracts/interfaces/IRegistryCoordinator.sol";
import {IDelegationManager} from "../../contracts/core/eigenContracts/interfaces/IDelegationManager.sol";
import {IStrategy} from "../../contracts/core/eigenContracts/interfaces/IStrategy.sol";
import {IStakeRegistry} from "../../contracts/core/eigenContracts/interfaces/IStakeRegistry.sol";
import {BitmapUtils} from "../../contracts/core/eigenContracts/libraries/BitmapUtils.sol";
import {DATA} from "../../contracts/core/libraries/Data.sol";
import {OperatorStateRetriever} from "../../contracts/core/eigenContracts/OperatorStateRetriever.sol";


contract BLSAvsMock is OperatorStateRetriever {
    using BN254 for BN254.G1Point;
    OperatorStateRetriever public operatorStateRetriever;
    address stakeRegistry;
    address delegationManagerAddress;

    bytes32 msgHash =
        0x16917735a95f9a155028c2190c39ba6c7be2ee45ccf802ed9817616b7bff76c2;
    uint256 aggSignerPrivKey = 69;
    /// @notice StakeRegistry, Constant used as a divisor in calculating weights.
    uint256 public constant WEIGHTING_DIVISOR = 1e18;
    BN254.G2Point aggSignerApkG2;
    BN254.G1Point sigmaDetail;
    uint256 quorumBitmap = 1;
    bytes quorumNumbersInput = BitmapUtils.bitmapToBytesArray(quorumBitmap);

    constructor(OperatorStateRetriever _operatorStateRetriever) {
        operatorStateRetriever = _operatorStateRetriever;
        //  _setAggregatePublicKeysAndSignature();
    }

    function _generateSignerAndNonSignerPrivateKeys(
        uint256 pseudoRandomNumber,
        uint256 numSigners,
        uint256 numNonSigners
    ) internal view returns (uint256[] memory, uint256[] memory) {
        uint256[] memory signerPrivateKeys = new uint256[](numSigners);
        // generate numSigners numbers that add up to aggSignerPrivKey mod BN254.FR_MODULUS
        uint256 sum = 0;
        for (uint i = 0; i < numSigners - 1; i++) {
            signerPrivateKeys[i] =
                uint256(
                    keccak256(
                        abi.encodePacked(
                            "signerPrivateKey",
                            pseudoRandomNumber,
                            i
                        )
                    )
                ) %
                BN254.FR_MODULUS;
            sum = addmod(sum, signerPrivateKeys[i], BN254.FR_MODULUS);
        }
        // signer private keys need to add to aggSignerPrivKey
        signerPrivateKeys[numSigners - 1] = addmod(
            aggSignerPrivKey,
            BN254.FR_MODULUS - (sum % BN254.FR_MODULUS),
            BN254.FR_MODULUS
        );

        uint256[] memory nonSignerPrivateKeys = new uint256[](numNonSigners);
        for (uint i = 0; i < numNonSigners; i++) {
            nonSignerPrivateKeys[i] =
                uint256(
                    keccak256(
                        abi.encodePacked(
                            "nonSignerPrivateKey",
                            pseudoRandomNumber,
                            i
                        )
                    )
                ) %
                BN254.FR_MODULUS;
        }

        // Sort nonSignerPrivateKeys in order of ascending pubkeyHash
        // Uses insertion sort to sort array in place
        for (uint i = 1; i < nonSignerPrivateKeys.length; i++) {
            uint privateKey = nonSignerPrivateKeys[i];
            bytes32 pubkeyHash = _toPubkeyHash(privateKey);
            uint j = i;

            // Move elements of nonSignerPrivateKeys[0..i-1] that are greater than the current key
            // to one position ahead of their current position
            while (
                j > 0 && _toPubkeyHash(nonSignerPrivateKeys[j - 1]) > pubkeyHash
            ) {
                nonSignerPrivateKeys[j] = nonSignerPrivateKeys[j - 1];
                j--;
            }
            nonSignerPrivateKeys[j] = privateKey;
        }

        return (signerPrivateKeys, nonSignerPrivateKeys);
    }

    function getDeviationHeader()
        public
        view
        returns (DATA.DeviationBalanceHeader memory)
    {
        bytes memory quorumThresholdPercentages = new bytes(1);
        quorumThresholdPercentages[0] = bytes1(uint8(65));
        return
            DATA.DeviationBalanceHeader({
                agentOperator: 0xDA29BB71669f46F2a779b4b62f03644A84eE3479,
                quorumNumbers: quorumNumbersInput,
                quorumThresholdPercentages: quorumThresholdPercentages,
                referenceBlockNumber: 25,
                actualAmount: 1350,
                idealAmount: 1500
            });
    }

    function _setAggregatePublicKeysAndSignature() internal {
        // aggSignerPrivKey*g2
        aggSignerApkG2.X[
                1
            ] = 19101821850089705274637533855249918363070101489527618151493230256975900223847;
        aggSignerApkG2.X[
                0
            ] = 5334410886741819556325359147377682006012228123419628681352847439302316235957;
        aggSignerApkG2.Y[
                1
            ] = 354176189041917478648604979334478067325821134838555150300539079146482658331;
        aggSignerApkG2.Y[
                0
            ] = 4185483097059047421902184823581361466320657066600218863748375739772335928910;

        sigmaDetail = BN254.hashToG1(msgHash).scalar_mul(aggSignerPrivKey);
    }

    function getNonSignerStakeAndSignature(
        IRegistryCoordinator registryCoordinator,
        uint32 referenceBlockNumber, 
        bytes calldata quorumNumbers, 
        bytes32[] calldata nonSignerOperatorIds
    ) public view  returns (DATA.NonSignerStakesAndSignature memory) {
       
        OperatorStateRetriever.CheckSignaturesIndices
            memory checkSignaturesIndices = operatorStateRetriever
                .getCheckSignaturesIndices(
                    registryCoordinator,
                    referenceBlockNumber,
                    quorumNumbers,
                    nonSignerOperatorIds
                );

        uint256[] memory privateKeys = new uint256[](1);
        privateKeys[0] = 253;
        BN254.G1Point[] memory nonSignerPubkeys = new BN254.G1Point[](
            privateKeys.length
        );

        for (uint256 i = 0; i < privateKeys.length; i++) {
            nonSignerPubkeys[i] = BN254.generatorG1().scalar_mul(
                privateKeys[i]
            );
        }
        BN254.G1Point[] memory quorumApksData = new BN254.G1Point[](
            quorumNumbers.length
        );

        //  return
        //     DATA.NonSignerStakesAndSignatureDemo({
        //         quorumApkIndices: 32,
        //         totalStakeIndices: 42
        //     });

        return
            DATA.NonSignerStakesAndSignature({
                nonSignerQuorumBitmapIndices: checkSignaturesIndices
                    .nonSignerQuorumBitmapIndices,
                nonSignerPubkeys: nonSignerPubkeys,
                quorumApks: quorumApksData,
                apkG2: aggSignerApkG2,
                sigma: sigmaDetail,
                quorumApkIndices: checkSignaturesIndices.quorumApkIndices,
                totalStakeIndices: checkSignaturesIndices.totalStakeIndices,
                nonSignerStakeIndices: checkSignaturesIndices
                    .nonSignerStakeIndices
            });
    }

    function signMessage(
        uint256 privateKey,
        address signer,
        address registryCoordinatorAddress
    ) external view returns (BN254.G1Point memory) {
        BN254.G1Point memory messageHash = IRegistryCoordinator(
            registryCoordinatorAddress
        ).pubkeyRegistrationMessageHash(signer);
        return BN254.scalar_mul(messageHash, privateKey);
    }

    function getG1PubKey(
        uint256 privateKey
    ) external view returns (BN254.G1Point memory) {
        return BN254.generatorG1().scalar_mul(privateKey);
    }

    // function getG2PubKey() external view returns(BN254.G2Point memory){
    //     BN254.G2Point memory messageHash = BN254.generatorG2();
    //      console.log("inside messageHash");
    //     return BN254.generatorG2();
    // }

    function _toPubkeyHash(uint privKey) internal view returns (bytes32) {
        return BN254.generatorG1().scalar_mul(privKey).hashG1Point();
    }

    function converBitmapToBytesArray(
        uint256 quorumBitmaps
    ) public pure returns (bytes memory /*bytesArray*/) {
        return BitmapUtils.bitmapToBytesArray(quorumBitmaps);
    }
}

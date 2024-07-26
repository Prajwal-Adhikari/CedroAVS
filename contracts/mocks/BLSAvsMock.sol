// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.12;

import {BLSSignatureChecker} from "../../contracts/core/eigenContracts/BLSSignatureChecker.sol";
import {BLSApkRegistry} from "../../contracts/core/eigenContracts/BLSApkRegistry.sol";
import {BN254} from "../../contracts/core/eigenContracts/libraries/BN254.sol";
import {IRegistryCoordinator} from "../../contracts/core/eigenContracts/interfaces/IRegistryCoordinator.sol";
import {IDelegationManager} from "../../contracts/core/eigenContracts/interfaces/IDelegationManager.sol";
import {ISignatureUtils} from "../../contracts/core/eigenContracts/interfaces/ISignatureUtils.sol";
import {IStrategy} from "../../contracts/core/eigenContracts/interfaces/IStrategy.sol";
import {IIndexRegistry} from "../../contracts/core/eigenContracts/interfaces/IIndexRegistry.sol";
import {IBLSSignatureChecker} from "../../contracts/core/eigenContracts/interfaces/IBLSSignatureChecker.sol";
import {BitmapUtils} from "../../contracts/core/eigenContracts/libraries/BitmapUtils.sol";
import {CedroServiceData} from "../../contracts/core/libraries/CedroServiceData.sol";
import {OperatorStateRetriever} from "../../contracts/core/eigenContracts/OperatorStateRetriever.sol";
import {IBLSApkRegistry} from "../../contracts/core/eigenContracts/interfaces/IBLSApkRegistry.sol";
import {BLSApkRegistryMock} from "./ BLSApkRegistryMock.sol";
import {DelegationMock} from "./DelegationMock.sol";
import {StakeRegistry} from "../../contracts/core/eigenContracts/StakeRegistry.sol";
import {RegistryCoordinator} from "../../contracts/core/eigenContracts/RegistryCoordinator.sol";
import {CedroServiceManager} from "../../contracts/core/CedroServiceManager.sol";

import "hardhat/console.sol";

contract BLSAvsMock is OperatorStateRetriever {
    using BN254 for BN254.G1Point;
    OperatorStateRetriever public operatorStateRetriever;
    StakeRegistry stakeRegistry;
    DelegationMock delegationManagerAddress;
    RegistryCoordinator registryCoordinator;
    CedroServiceManager cedroServiceManager;
    IIndexRegistry indexRegistry;

    bytes32 msgHash =
        0xd6f4db7833690dfab0ec49d35dfe04e1ee9c162e5d3316de21f394a96a652938;
    bytes32 fraudMsgHash =
        0xb1d2ac1967f0fc912756edfea132ebc4ad9fbe9a3f5b34c78c1889bc87c1fba3;
    uint256 aggSignerPrivKey =
        15436643226594533498284352790494521296965667247196462338996131163212826097054;
    /// @notice StakeRegistry, Constant used as a divisor in calculating weights.
    uint256 public constant WEIGHTING_DIVISOR = 1e18;
    BN254.G2Point aggSignerApkG2;
    BN254.G1Point sigmaDetail;
    BN254.G1Point quorumApksData;

    // uint256 quorumBitmap = 1;
    // bytes quorumNumbersInput = BitmapUtils.bitmapToBytesArray(quorumBitmap);
    uint96 defaultStake = 1 ether;
    address defaultOperator =
        address(uint160(uint256(keccak256("defaultOperator"))));
    bytes32 defaultOperatorId;
    BN254.G1Point internal defaultPubKey =
        BN254.G1Point(
            18260007818883133054078754218619977578772505796600400998181738095793040006897,
            3432351341799135763167709827653955074218841517684851694584291831827675065899
        );
    BN254.G1Point sigma;
    uint256 MAX_QUORUM_BITMAP = type(uint192).max;
    BLSApkRegistryMock blsApkRegistryMock;
    DelegationMock delegationMock;
    string defaultSocket = "69.69.69.69:420";
    IBLSApkRegistry.PubkeyRegistrationParams pubkeyRegistrationParams;

    constructor(
        CedroServiceManager _cedroServiceManager,
        OperatorStateRetriever _operatorStateRetriever,
        RegistryCoordinator _registryCoordinator,
        BLSApkRegistryMock _blsApkRegistry,
        DelegationMock _delegationMock,
        StakeRegistry _stakeRegistry
    ) {
        operatorStateRetriever = _operatorStateRetriever;
        registryCoordinator = _registryCoordinator;
        blsApkRegistryMock = _blsApkRegistry;
        delegationMock = _delegationMock;
        stakeRegistry = _stakeRegistry;
        cedroServiceManager = _cedroServiceManager;
        indexRegistry = registryCoordinator.indexRegistry();
        _setAggregatePublicKeysAndSignature();
    }

    function getDeviationHeader()
        public
        pure
        returns (CedroServiceData.DeviationBalanceHeader memory)
    {
        bytes memory quorumThresholdPercentages = new bytes(1);
        quorumThresholdPercentages[0] = bytes1(uint8(0));
        // quorumThresholdPercentages[1] = bytes1(uint8(0));

        bytes memory quorumNumbersInput = new bytes(1);
        quorumNumbersInput[0] = bytes1(uint8(0));
        // quorumNumbersInput[1] = bytes1(uint8(1));
        return
            CedroServiceData.DeviationBalanceHeader({
                agentAddress: 0xdF1Cf915dB7A658dFDEb75aF3873B4e55F25Aa42,
                referenceBlockNumber: 1860891,
                quorumNumbers: quorumNumbersInput,
                quorumThresholdPercentages: quorumThresholdPercentages,
                requestId: 3,
                actualAmount: 1920000000000000000000,
                idealAmount: 3110400000000000001920000000000000000000
            });
    }

    function getFraudHeader()
        public
        pure
        returns (CedroServiceData.BankRunBalanceHeader memory)
    {
        bytes memory quorumThresholdPercentages = new bytes(1);
        quorumThresholdPercentages[0] = bytes1(uint8(0));
        // quorumThresholdPercentages[1] = bytes1(uint8(0));

        bytes memory quorumNumbersInput = new bytes(1);
        quorumNumbersInput[0] = bytes1(uint8(0));
        // quorumNumbersInput[1] = bytes1(uint8(1));
        return
            CedroServiceData.BankRunBalanceHeader({
                agentAddress: 0xdF1Cf915dB7A658dFDEb75aF3873B4e55F25Aa42,
                agentContract: 0xdF1Cf915dB7A658dFDEb75aF3873B4e55F25Aa42,
                referenceBlockNumber: 1860891,
                chainId: "0x424e42",
                quorumNumbers: quorumNumbersInput,
                quorumThresholdPercentages: quorumThresholdPercentages,
                transactionId: 3,
                insufficientAmount: 1920000000000000000000
            });
    }

    function _setAggregatePublicKeysAndSignature() internal {
        aggSignerApkG2.X[
                0
            ] = 1230730797123182311850807131235073178385821783983372276639854957091263955384;
        aggSignerApkG2.X[
                1
            ] = 7034359866883434113182571131287568261842779761851733188885430766201745859762;
        aggSignerApkG2.Y[
                0
            ] = 17453839935845982565550211636126433446255050520010422383769263780951237157828;
        aggSignerApkG2.Y[
                1
            ] = 16271686786459442426557467568799518089585926782040698082337113377336362878238;

        sigmaDetail = BN254.hashToG1(fraudMsgHash).scalar_mul(aggSignerPrivKey);
        sigma
            .X = 14700674594185938633638077079219260099284573239895592630666105257857049367678;
        sigma
            .Y = 6079443873574917784354775434183978606910175512518269837094952080875748687466;
    }

    function getNonsignerAddress(
        uint256 operatorNumber
    ) internal pure returns (BN254.G1Point memory nonSignerPubKey) {
        if (operatorNumber == 1) {
            nonSignerPubKey = BN254.G1Point(0, 0);

            nonSignerPubKey
                .X = 14110567285676345389244725775337526580401753248023284461771486167266462355103;
            nonSignerPubKey
                .Y = 9720413144746361199461732946893043367171487505859650901952835867397436923362;
        }
        if (operatorNumber == 2) {
            nonSignerPubKey = BN254.G1Point(0, 0);

            nonSignerPubKey
                .X = 11909330550897052025764217059883091939734032688455609568393118274846029799265;
            nonSignerPubKey
                .Y = 18867945654008209292625458880315629237005264805214925261372079909097324926955;
        }
    }

    function getDeviationNonSignerStakeAndSignature(
        uint32 referenceBlockNumber,
        uint256 quorumBitmaps,
        address[] memory nonSigners
    )
        public
        view
        returns (BLSSignatureChecker.NonSignerStakesAndSignature memory)
    {
        bytes32[] memory nonSignerOperatorIds = new bytes32[](
            nonSigners.length
        );
        for (uint256 i; i < nonSigners.length; i++) {
            nonSignerOperatorIds[i] = registryCoordinator.getOperatorId(
                nonSigners[i]
            );
        }
        //  bytes32[] memory nonSignerOperatorIds = new bytes32[](nonSignerOperatorIds.length);
        bytes memory quorumNumbers = BitmapUtils.bitmapToBytesArray(
            quorumBitmaps
        );
        BLSSignatureChecker.NonSignerStakesAndSignature
            memory nonSignerStakesAndSignature;
        OperatorStateRetriever.CheckSignaturesIndices
            memory checkSignaturesIndices = operatorStateRetriever
                .getCheckSignaturesIndices(
                    registryCoordinator,
                    referenceBlockNumber,
                    quorumNumbers,
                    nonSignerOperatorIds
                );

        nonSignerStakesAndSignature.quorumApks = new BN254.G1Point[](
            quorumNumbers.length
        );
        nonSignerStakesAndSignature.nonSignerPubkeys = new BN254.G1Point[](
            nonSigners.length
        );

        console.log("Quorum length :", quorumNumbers.length);
        for (
            uint256 i = 0;
            i < nonSignerStakesAndSignature.nonSignerPubkeys.length;
            i++
        ) {
            (uint256 pubkeysX, uint256 pubkeysY) = blsApkRegistryMock
                .operatorToPubkey(nonSigners[i]);
            nonSignerStakesAndSignature.nonSignerPubkeys[i].X = pubkeysX;
            nonSignerStakesAndSignature.nonSignerPubkeys[i].Y = pubkeysY;
        }

        for (
            uint256 i = 0;
            i < nonSignerStakesAndSignature.quorumApks.length;
            i++
        ) {
            BN254.G1Point[] memory pubKeysFinal = new BN254.G1Point[](
                quorumNumbers.length
            );

            bytes32[] memory OperatorIds = indexRegistry
                .getOperatorListAtBlockNumber(uint8(i), referenceBlockNumber);
            for (uint8 j; j < OperatorIds.length; j++) {
                BN254.G1Point[] memory pubKeysData = new BN254.G1Point[](
                    OperatorIds.length
                );
                (uint256 pubkeysX, uint256 pubkeysY) = (
                    blsApkRegistryMock.operatorToPubkey(
                        blsApkRegistryMock.getOperatorFromPubkeyHash(
                            OperatorIds[j]
                        )
                    )
                );
                pubKeysData[j].X = pubkeysX;
                pubKeysData[j].Y = pubkeysY;
                pubKeysFinal[i] = pubKeysFinal[i].plus(pubKeysData[j]);
            }
            nonSignerStakesAndSignature.quorumApks[
                    i
                ] = nonSignerStakesAndSignature.quorumApks[i].plus(
                pubKeysFinal[i]
            );
            console.log(
                "Pubkeys quorumapk X:",
                nonSignerStakesAndSignature.quorumApks[i].X
            );
            console.log(
                "Pubkeys quorumapk Y:",
                nonSignerStakesAndSignature.quorumApks[i].Y
            );
        }

        nonSignerStakesAndSignature
            .nonSignerQuorumBitmapIndices = checkSignaturesIndices
            .nonSignerQuorumBitmapIndices;
        nonSignerStakesAndSignature.apkG2 = aggSignerApkG2;
        nonSignerStakesAndSignature.sigma = sigma;
        nonSignerStakesAndSignature.quorumApkIndices = checkSignaturesIndices
            .quorumApkIndices;
        nonSignerStakesAndSignature.totalStakeIndices = checkSignaturesIndices
            .totalStakeIndices;
        nonSignerStakesAndSignature
            .nonSignerStakeIndices = checkSignaturesIndices
            .nonSignerStakeIndices;

        return nonSignerStakesAndSignature;
    }

    function getFraudNonSignerStakeAndSignature(
        uint32 referenceBlockNumber,
        uint256 quorumBitmaps,
        address[] memory nonSigners
    )
        public
        view
        returns (BLSSignatureChecker.NonSignerStakesAndSignature memory)
    {
        bytes32[] memory nonSignerOperatorIds = new bytes32[](
            nonSigners.length
        );
        for (uint256 i; i < nonSigners.length; i++) {
            nonSignerOperatorIds[i] = registryCoordinator.getOperatorId(
                nonSigners[i]
            );
        }
        //  bytes32[] memory nonSignerOperatorIds = new bytes32[](nonSignerOperatorIds.length);
        bytes memory quorumNumbers = BitmapUtils.bitmapToBytesArray(
            quorumBitmaps
        );
        BLSSignatureChecker.NonSignerStakesAndSignature
            memory nonSignerStakesAndSignature;
        OperatorStateRetriever.CheckSignaturesIndices
            memory checkSignaturesIndices = operatorStateRetriever
                .getCheckSignaturesIndices(
                    registryCoordinator,
                    referenceBlockNumber,
                    quorumNumbers,
                    nonSignerOperatorIds
                );

        nonSignerStakesAndSignature.quorumApks = new BN254.G1Point[](
            quorumNumbers.length
        );
        nonSignerStakesAndSignature.nonSignerPubkeys = new BN254.G1Point[](
            nonSigners.length
        );

        console.log("Quorum length :", quorumNumbers.length);
        for (
            uint256 i = 0;
            i < nonSignerStakesAndSignature.nonSignerPubkeys.length;
            i++
        ) {
            (uint256 pubkeysX, uint256 pubkeysY) = blsApkRegistryMock
                .operatorToPubkey(nonSigners[i]);
            nonSignerStakesAndSignature.nonSignerPubkeys[i].X = pubkeysX;
            nonSignerStakesAndSignature.nonSignerPubkeys[i].Y = pubkeysY;
        }

        for (
            uint256 i = 0;
            i < nonSignerStakesAndSignature.quorumApks.length;
            i++
        ) {
            BN254.G1Point[] memory pubKeysFinal = new BN254.G1Point[](
                quorumNumbers.length
            );

            bytes32[] memory OperatorIds = indexRegistry
                .getOperatorListAtBlockNumber(uint8(i), referenceBlockNumber);
            for (uint8 j; j < OperatorIds.length; j++) {
                BN254.G1Point[] memory pubKeysData = new BN254.G1Point[](
                    OperatorIds.length
                );
                (uint256 pubkeysX, uint256 pubkeysY) = (
                    blsApkRegistryMock.operatorToPubkey(
                        blsApkRegistryMock.getOperatorFromPubkeyHash(
                            OperatorIds[j]
                        )
                    )
                );
                pubKeysData[j].X = pubkeysX;
                pubKeysData[j].Y = pubkeysY;
                pubKeysFinal[i] = pubKeysFinal[i].plus(pubKeysData[j]);
            }
            nonSignerStakesAndSignature.quorumApks[
                    i
                ] = nonSignerStakesAndSignature.quorumApks[i].plus(
                pubKeysFinal[i]
            );
            console.log(
                "Pubkeys quorumapk X:",
                nonSignerStakesAndSignature.quorumApks[i].X
            );
            console.log(
                "Pubkeys quorumapk Y:",
                nonSignerStakesAndSignature.quorumApks[i].Y
            );
        }

        nonSignerStakesAndSignature
            .nonSignerQuorumBitmapIndices = checkSignaturesIndices
            .nonSignerQuorumBitmapIndices;
        nonSignerStakesAndSignature.apkG2 = aggSignerApkG2;
        nonSignerStakesAndSignature.sigma = sigmaDetail;
        nonSignerStakesAndSignature.quorumApkIndices = checkSignaturesIndices
            .quorumApkIndices;
        nonSignerStakesAndSignature.totalStakeIndices = checkSignaturesIndices
            .totalStakeIndices;
        nonSignerStakesAndSignature
            .nonSignerStakeIndices = checkSignaturesIndices
            .nonSignerStakeIndices;

        return nonSignerStakesAndSignature;
    }

    function testAgentDeviation() external pure {}

    function _incrementAddress(
        address start,
        uint256 inc
    ) internal pure returns (address) {
        return address(uint160(uint256(uint160(start) + inc)));
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

    /**
     * @dev Set the operator weight for a given quorum. Note we have to do this by setting delegationMock operatorShares
     * Given each quorum must have at least one strategy, we set operatorShares for this strategy to this weight
     * Returns actual weight calculated set for operator shares in DelegationMock since multiplier and WEIGHTING_DIVISOR calculations
     * can give small rounding errors.
     */
    function _setOperatorWeight(
        address operator,
        uint8 quorumNumber,
        uint96 weight
    ) internal returns (uint96) {
        // Set StakeRegistry operator weight by setting DelegationManager operator shares
        (IStrategy strategy, uint96 multiplier) = stakeRegistry.strategyParams(
            quorumNumber,
            0
        );

        uint256 actualWeight = ((uint256(weight) * WEIGHTING_DIVISOR) /
            uint256(multiplier));
        delegationMock.setOperatorShares(operator, strategy, actualWeight);
        return uint96(actualWeight);
    }

    function converBitmapToBytesArray(
        uint256 quorumBitmaps
    ) public pure returns (bytes memory /*bytesArray*/) {
        return BitmapUtils.bitmapToBytesArray(quorumBitmaps);
    }

    function testDeviation(
        uint32 referenceBlockNumber,
        uint256 quorumBitmaps,
        address[] memory nonSigners
    ) external {
        //   function testDeviation() external view {

        bytes32 messagehash = 0xa952620b763ac36ee313124f0e665018fc281b1ba479c19cadb3c3e5f8fee122;
        BN254.G1Point memory sigmaValue = BN254
            .hashToG1(messagehash)
            .scalar_mul(aggSignerPrivKey);
        sigmaValue = sigmaValue.plus(sigmaValue);
        console.log("Sigma Value X ", sigmaValue.X);
        console.log("Sigma value Y : ", sigmaValue.Y);
        BN254.G1Point[] memory quorumApksValue = new BN254.G1Point[](2);
        quorumApksValue[0]
            .X = 14110567285676345389244725775337526580401753248023284461771486167266462355103;
        quorumApksValue[0]
            .Y = 9720413144746361199461732946893043367171487505859650901952835867397436923362;
        quorumApksValue[1]
            .X = 14110567285676345389244725775337526580401753248023284461771486167266462355103;
        quorumApksValue[1]
            .Y = 9720413144746361199461732946893043367171487505859650901952835867397436923362;

        BN254.G1Point memory apk = BN254.G1Point(0, 0);
        // for(uint256 i=0;i<2;i++){
        //      apk = apk.plus(quorumApksValue[i]);
        //      }
        console.log("Console.log Apk", apk.X);
        console.log("Console.log Apk", apk.Y);

        apk
            .X = 14110567285676345389244725775337526580401753248023284461771486167266462355103;
        apk
            .Y = 9720413144746361199461732946893043367171487505859650901952835867397436923362;

        (bool pairingSuccessful, bool siganatureIsValid) = cedroServiceManager
            .trySignatureAndApkVerification(
                messagehash,
                apk,
                aggSignerApkG2,
                sigma
            );
        console.log("Pairing successful", pairingSuccessful);
        console.log("Signature is valid", siganatureIsValid);
        cedroServiceManager.confirmAgentDeviation(
            getDeviationHeader(),
            getDeviationNonSignerStakeAndSignature(
                referenceBlockNumber,
                quorumBitmaps,
                nonSigners
            )
        );
    }

    function testFraud(
        uint32 referenceBlockNumber,
        uint256 quorumBitmaps,
        address[] memory nonSigners
    ) external {
        bytes32 messagehash = 0x456cefdadb1fda7899f355e0ad99a17429f8fd97b2d29f8cd811c8801f4b1d4f;
        BN254.G1Point memory sigmaValue = BN254
            .hashToG1(messagehash)
            .scalar_mul(aggSignerPrivKey);
        sigmaValue = sigmaValue.plus(sigmaValue);
        console.log("Sigma Value X ", sigmaValue.X);
        console.log("Sigma value Y : ", sigmaValue.Y);
        BN254.G1Point[] memory quorumApksValue = new BN254.G1Point[](2);
        quorumApksValue[0]
            .X = 14110567285676345389244725775337526580401753248023284461771486167266462355103;
        quorumApksValue[0]
            .Y = 9720413144746361199461732946893043367171487505859650901952835867397436923362;
        quorumApksValue[1]
            .X = 14110567285676345389244725775337526580401753248023284461771486167266462355103;
        quorumApksValue[1]
            .Y = 9720413144746361199461732946893043367171487505859650901952835867397436923362;

        BN254.G1Point memory apk = BN254.G1Point(0, 0);
        // for(uint256 i=0;i<2;i++){
        //      apk = apk.plus(quorumApksValue[i]);
        //      }
        console.log("Console.log Apk", apk.X);
        console.log("Console.log Apk", apk.Y);

        apk
            .X = 14110567285676345389244725775337526580401753248023284461771486167266462355103;
        apk
            .Y = 9720413144746361199461732946893043367171487505859650901952835867397436923362;

        (bool pairingSuccessful, bool siganatureIsValid) = cedroServiceManager
            .trySignatureAndApkVerification(
                messagehash,
                apk,
                aggSignerApkG2,
                sigmaDetail
            );
        console.log("Pairing successful", pairingSuccessful);
        console.log("Signature is valid", siganatureIsValid);
        cedroServiceManager.confirmAgentBankRun(
            getFraudHeader(),
            getFraudNonSignerStakeAndSignature(
                referenceBlockNumber,
                quorumBitmaps,
                nonSigners
            )
        );
    }
}

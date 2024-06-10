//Testing of connection with the Eigen Layer

const { assert, expect } = require("chai");
const chai = require("chai");
const chaiAlmost = require("chai-almost");

chai.use(chaiAlmost());
const { ethers, upgrades, config } = require("hardhat");
const { anyValue } = require("@nomicfoundation/hardhat-chai-matchers/withArgs");

//contracts
let cedroServiceManger,
  avsDirectoryMock,
  blsAvsMock,
  operatorStateRetriever,
  delegationMock,
  paymentCoordinatorMock,
  registryCoordinator,
  stakeRegistry,
  serviceManagerBase,
  indexRegistry,
  blsApkRegistryMock,
  pauserRegistry,
  roleManagerCore;

async function uintToAddress(uintValue) {
  // Convert the uint to a BigNumber
  const bigNumber = ethers.BigNumber.from(uintValue);

  // Convert the BigNumber to a hexadecimal string
  const hexString = bigNumber.toHexString();

  // Ensure the string starts with '0x'
  const addressLikeString = "0x" + hexString;

  return addressLikeString;
}

async function getBlsprivKey() {
  // Determine the order r (for BLS12-381 curve)
  const r =
    "0x73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001";

  // Generate a random 32-byte hex string
  const randomBytes = ethers.randomBytes(32);

  // Convert the random bytes to a BigNumber
  // const randomBigNumber = ethers.BigInt(ethers.toUtf8String(randomBytes));

  // Take the random BigNumber modulo r to get a scalar in the finite field Fr
  const randomScalar = randomBigNumber.mod(ethers.BigNumber.from(r));
}

describe("Cedro Service Manager Tests", function () {
  beforeEach(async function () {
    [admin, operator1, operator2, operator3] = await ethers.getSigners();

    const CedroServiceManager = await ethers.getContractFactory(
      "CedroServiceManager"
    );
    const RegistryCoordinator = await ethers.getContractFactory(
      "RegistryCoordinator"
    );

    const AVSDirectoryMock = await ethers.getContractFactory(
      "AVSDirectoryMock"
    );
    const PaymentCoordinatorMock = await ethers.getContractFactory(
      "PaymentCoordinatorMock"
    );
    const BLSAvsMock = await ethers.getContractFactory("BLSAvsMock");
    const DelegationMock = await ethers.getContractFactory("DelegationMock");
    const BlsApkRegistryMock = await ethers.getContractFactory(
      "BLSApkRegistryMock"
    );
    const IndexRegistry = await ethers.getContractFactory("IndexRegistry");
    const StakeRegistry = await ethers.getContractFactory("StakeRegistry");
    const PauserRegistry = await ethers.getContractFactory("PauserRegistry");
    const OperatorStateRetriever = await ethers.getContractFactory(
      "OperatorStateRetriever"
    );
    const RoleManagerCore = await ethers.getContractFactory("RoleManagerCore");

    roleManagerCore = await upgrades.deployProxy(RoleManagerCore, [
      admin.address,
    ]);
    paymentCoordinatorMock = await PaymentCoordinatorMock.deploy();
    delegationMock = await DelegationMock.deploy();
    avsDirectoryMock = await AVSDirectoryMock.deploy(delegationMock.target);
    pauserRegistry = await PauserRegistry.deploy(admin.address, admin.address);
    blsApkRegistryMock = await BlsApkRegistryMock.deploy(admin.address);
    indexRegistry = await IndexRegistry.deploy(admin.address);
    stakeRegistry = await StakeRegistry.deploy(
      admin.address,
      delegationMock.target
    );
    operatorStateRetriever = await OperatorStateRetriever.deploy();
    blsAvsMock = await BLSAvsMock.deploy(operatorStateRetriever);
    cedroServiceManger = await upgrades.deployProxy(CedroServiceManager, [
      roleManagerCore.target,
      avsDirectoryMock.target,
      paymentCoordinatorMock.target,
      stakeRegistry.target,
    ]);

    registryCoordinator = await upgrades.deployProxy(RegistryCoordinator, [
      admin.address,
      admin.address,
      cedroServiceManger.target,
      stakeRegistry.target,
      blsApkRegistryMock.target,
      indexRegistry.target,
    ]);

    blsApkRegistryMock.setRegisterCoordinator(registryCoordinator.target);
    indexRegistry.setRegisterCoordinator(registryCoordinator.target);
    stakeRegistry.setRegisterCoordinator(registryCoordinator.target);
    cedroServiceManger.setRegistryCoordinatorAddress(
      registryCoordinator.target
    );
    cedroServiceManger.setRegistryCoordinator(registryCoordinator.target);

    const blockNumber = await ethers.provider.getBlockNumber();
    block = await ethers.provider.getBlock(blockNumber);
  });

  describe("1. Testing CedroServiceManager Functionality ", () => {
    it("Testing Operator Registration in AVS", async () => {
      privateKey =
        "0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d";
      wallet = new ethers.Wallet(privateKey);
      signingKey = new ethers.SigningKey(wallet.privateKey);
      let salt = ethers.keccak256(
        ethers.solidityPackedKeccak256(["string"], ["default"])
      );
      let expiry = block.timestamp + 100;
      let digestHash =
        await avsDirectoryMock.calculateOperatorAVSRegistrationDigestHash(
          operator1.address,
          cedroServiceManger.target,
          salt,
          expiry
        );
      signature = signingKey.sign(digestHash);
      signatureConcat = ethers.Signature.from(signature).serialized;
      await delegationMock.setIsOperator(operator1.address, true);
      await cedroServiceManger.addOperatorsToAllowlist([operator1.address]);
      await cedroServiceManger.registerOperatorToAVS(operator1, [
        signatureConcat,
        salt,
        expiry,
      ]);
    });

    it("Creating Quorum in AVS", async () => {
      const numQuorumsToAdd = 2; // Example value

      let minimumStakeForQuorum = new Array(numQuorumsToAdd);

      let operatorSetParams = new Array(numQuorumsToAdd);

      for (let i = 0; i < minimumStakeForQuorum.length; i++) {
        minimumStakeForQuorum[i] = BigInt(i + 1);
      }

      for (let i = 0; i < operatorSetParams.length; i++) {
        operatorSetParams[i] = [10, 15000, 150];
      }

      for (let i = 0; i < numQuorumsToAdd.length; i++) {
        minimumStakeForQuorum[i] = BigInt(i + 1);
      }

      const WEIGHTING_DIVISOR = 1e18;
      const quorumStrategiesConsideredAndMultipliers = [];
      for (let i = 0; i < numQuorumsToAdd; i++) {
        const bigNumber = BigInt(i);

        const hexString = bigNumber.toString();

        const addressLikeString = "0x" + hexString;

        strategyId = operator1.address;

        // Here, you would typically call a method on the contract to add a strategy
        // Since we're mocking, we just push the strategy details to our array
        quorumStrategiesConsideredAndMultipliers[i] = [
          [strategyId, WEIGHTING_DIVISOR],
        ];

        for (let i = 0; i < operatorSetParams.length; i++) {
          registryCoordinator.createQuorum(
            operatorSetParams[i],
            minimumStakeForQuorum[i],
            quorumStrategiesConsideredAndMultipliers[i]
          );
        }
      }
    });

    it("Register operator and enroll in the quorum", async () => {
      //Create Quorum
      await delegationMock.setIsOperator(operator1.address, true);
      await cedroServiceManger.addOperatorsToAllowlist([operator1.address]);

      const numQuorumsToAdd = 2; // Example value

      let minimumStakeForQuorum = new Array(numQuorumsToAdd);

      let operatorSetParams = new Array(numQuorumsToAdd);

      for (let i = 0; i < minimumStakeForQuorum.length; i++) {
        minimumStakeForQuorum[i] = BigInt(i + 1);
      }

      for (let i = 0; i < operatorSetParams.length; i++) {
        operatorSetParams[i] = [10, 15000, 150];
      }

      const WEIGHTING_DIVISOR = 1e12;
      const quorumStrategiesConsideredAndMultipliers = [];
      for (let i = 0; i < numQuorumsToAdd; i++) {
        strategyId = operator1.address;

        // Here, you would typically call a method on the contract to add a strategy
        // Since we're mocking, we just push the strategy details to our array
        quorumStrategiesConsideredAndMultipliers[i] = [
          [strategyId, WEIGHTING_DIVISOR],
        ];
      }

      for (let i = 0; i < operatorSetParams.length; i++) {
        await registryCoordinator.createQuorum(
          operatorSetParams[i],
          minimumStakeForQuorum[i],
          quorumStrategiesConsideredAndMultipliers[i]
        );
      }

      privateKey =
        "0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d";
      wallet = new ethers.Wallet(privateKey);
      signingKey = new ethers.SigningKey(wallet.privateKey);
      let salt = ethers.keccak256(
        ethers.solidityPackedKeccak256(["string"], ["default"])
      );
      let expiry = block.timestamp + 100;
      let digestHash =
        await avsDirectoryMock.calculateOperatorAVSRegistrationDigestHash(
          operator1.address,
          cedroServiceManger.target,
          salt,
          expiry
        );
      signature = signingKey.sign(digestHash);
      signatureConcat = ethers.Signature.from(signature).serialized;
      let quorumNumbers = await blsAvsMock.converBitmapToBytesArray(1);

      blsPrivateKey = 853;

      let pubkeyRegistrationSignature = await blsAvsMock.signMessage(
        blsPrivateKey,
        operator1.address,
        registryCoordinator.target
      );
      // let pubkeyRegistrationSignature =[13540890549107538365317826005663225429149650894201009774295406557850371632651n,2966724547965158980192270741744292034179150792486555967152513253931665577541n]
      // let pubkeyG1 = await blsAvsMock.getG1PubKey(blsPrivateKey);
      let pubkeyG1x =
        21008703310028844141620162130587521182514488709533704622337026710714836074802n;
      let pubkeyG1y =
        1145283673551783492487423178395020830991754869805176606412560778073024405948n;
      let pubkeyG2x = [
        12302775405227555806769733102102076488786351091672743244822361848998753983454n,
        12302775405227555806769733102102076488786351091672743244822361848998753983454n,
      ];
      let pubkeyG2y = [
        12302775405227555806769733102102076488786351091672743244822361848998753983454n,
        12302775405227555806769733102102076488786351091672743244822361848998753983454n,
      ];
      // let pubkeyG2 = await blsAvsMock.getG2PubKey();

      pubkeyRegistrationParams = [
        [pubkeyG1x, pubkeyG1y],
        [pubkeyG1x, pubkeyG1y],
        [pubkeyG2x, pubkeyG2y],
      ];
      let operatorSignatureParams = [signatureConcat, salt, expiry];
      await blsApkRegistryMock.setBLSPublicKey(operator1, [
        pubkeyG1x,
        pubkeyG1y,
      ]);

      let weight = ethers.parseEther("1.0");
      defaultStake = 1;
      let WEIGHTING_DIVISORS = BigInt(1e18);

      //   // _setOperatorWeight(operator, quorumNumbers[i], defaultStake);
      let strategyDetails = await stakeRegistry.strategyParams(0, 0);
      let actualWeight = (weight * WEIGHTING_DIVISORS) / WEIGHTING_DIVISORS;
      delegationMock.setOperatorShares(
        operator1.address,
        strategyDetails[0],
        actualWeight
      );

      await registryCoordinator.connect(operator1).registerOperator(
        quorumNumbers,
        [
          [pubkeyG1x, pubkeyG1y],
          [pubkeyG1x, pubkeyG1y],
          [pubkeyG2x, pubkeyG2y],
        ],
        operatorSignatureParams
      );
    });

    it.only("Test the signature verification flow", async () => {
      await delegationMock.setIsOperator(operator1.address, true);
      await cedroServiceManger.addOperatorsToAllowlist([operator1.address]);

      const numQuorumsToAdd = 2; // Example value

      let minimumStakeForQuorum = new Array(numQuorumsToAdd);

      let operatorSetParams = new Array(numQuorumsToAdd);

      for (let i = 0; i < minimumStakeForQuorum.length; i++) {
        minimumStakeForQuorum[i] = BigInt(i + 1);
      }

      for (let i = 0; i < operatorSetParams.length; i++) {
        operatorSetParams[i] = [10, 15000, 150];
      }

      const WEIGHTING_DIVISOR = 1e12;
      const quorumStrategiesConsideredAndMultipliers = [];
      for (let i = 0; i < numQuorumsToAdd; i++) {
        strategyId = operator1.address;

        // Here, you would typically call a method on the contract to add a strategy
        // Since we're mocking, we just push the strategy details to our array
        quorumStrategiesConsideredAndMultipliers[i] = [
          [strategyId, WEIGHTING_DIVISOR],
        ];
      }

      for (let i = 0; i < operatorSetParams.length; i++) {
        await registryCoordinator.createQuorum(
          operatorSetParams[i],
          minimumStakeForQuorum[i],
          quorumStrategiesConsideredAndMultipliers[i]
        );
      }

      privateKey =
        "0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d";
      wallet = new ethers.Wallet(privateKey);
      signingKey = new ethers.SigningKey(wallet.privateKey);
      let salt = ethers.keccak256(
        ethers.solidityPackedKeccak256(["string"], ["default"])
      );
      let expiry = block.timestamp + 100;
      let digestHash =
        await avsDirectoryMock.calculateOperatorAVSRegistrationDigestHash(
          operator1.address,
          cedroServiceManger.target,
          salt,
          expiry
        );
      signature = signingKey.sign(digestHash);
      signatureConcat = ethers.Signature.from(signature).serialized;
      let quorumNumbers = await blsAvsMock.converBitmapToBytesArray(1);
      blsPrivateKey = 853;

      let pubkeyRegistrationSignature = await blsAvsMock.signMessage(
        blsPrivateKey,
        operator1.address,
        registryCoordinator.target
      );
      // let pubkeyRegistrationSignature =[13540890549107538365317826005663225429149650894201009774295406557850371632651n,2966724547965158980192270741744292034179150792486555967152513253931665577541n]
      // let pubkeyG1 = await blsAvsMock.getG1PubKey(blsPrivateKey);
      let pubkeyG1x =
        21008703310028844141620162130587521182514488709533704622337026710714836074802n;
      let pubkeyG1y =
        1145283673551783492487423178395020830991754869805176606412560778073024405948n;
      let pubkeyG2x = [
        12302775405227555806769733102102076488786351091672743244822361848998753983454n,
        12302775405227555806769733102102076488786351091672743244822361848998753983454n,
      ];
      let pubkeyG2y = [
        12302775405227555806769733102102076488786351091672743244822361848998753983454n,
        12302775405227555806769733102102076488786351091672743244822361848998753983454n,
      ];
      // let pubkeyG2 = await blsAvsMock.getG2PubKey();

      pubkeyRegistrationParams = [
        [pubkeyG1x, pubkeyG1y],
        [pubkeyG1x, pubkeyG1y],
        [pubkeyG2x, pubkeyG2y],
      ];
      let operatorSignatureParams = [signatureConcat, salt, expiry];
      await blsApkRegistryMock.setBLSPublicKey(operator1, [
        pubkeyG1x,
        pubkeyG1y,
      ]);

      let weight = ethers.parseEther("1.0");
      defaultStake = 1;
      let WEIGHTING_DIVISORS = BigInt(1e18);

      //   // _setOperatorWeight(operator, quorumNumbers[i], defaultStake);
      let strategyDetails = await stakeRegistry.strategyParams(0, 0);
      let actualWeight = (weight * WEIGHTING_DIVISORS) / WEIGHTING_DIVISORS;
      delegationMock.setOperatorShares(
        operator1.address,
        strategyDetails[0],
        actualWeight
      );

      await registryCoordinator.connect(operator1).registerOperator(
        quorumNumbers,
        [
          [pubkeyG1x, pubkeyG1y],
          [pubkeyG1x, pubkeyG1y],
          [pubkeyG2x, pubkeyG2y],
        ],
        operatorSignatureParams
      );
      /* ==> De register Operator <== */
      // await registryCoordinator.connect(operator1).deregisterOperator(quorumNumbers);

      /* ==> confirming signature Check <== */
      let nonRandomNumber = 111;
      let numNonSigners = 0;
      let quorumBitmap = 1;
      let referenceBlockNumbers = 123;
      let nonsigenrOperatorIds = await blsApkRegistryMock.getOperatorId(operator1.address);
      // console.log(bls);
      let deviationHeader = await blsAvsMock.getDeviationHeader();
      let deviationHeaderData = [
        deviationHeader[0],
        deviationHeader[1],
        deviationHeader[2],
        deviationHeader[3],
        deviationHeader[4],
        deviationHeader[5],
      ];
      // let nonsignRequiredData = registryCoordinator.target;
      let NonSignerStakesAndSignature =
        await blsAvsMock.getNonSignerStakeAndSignature(
          registryCoordinator.target,
          referenceBlockNumbers,
          quorumNumbers,
          [ nonsigenrOperatorIds]
        );
      // console.log("non signer Data :", NonSignerStakesAndSignature);
      console.log(NonSignerStakesAndSignature[0][0].X);
      let nonSignerQuorumBitmapIndicesData = [NonSignerStakesAndSignature[0][0].X,NonSignerStakesAndSignature[0][0].y];
      // let nonSignerPubkeysData = NonSignerStakesAndSignature[1]
      // let NonSignerStakeAndSignatureData = [
      //   nonSignerQuorumBitmapIndicesData,
      // ];

      // console.log(await blsAvsMock.getDeviationHeader());
      // console.log(await blsAvsMock.getNonSignerStakeAndSignature(registryCoordinator.target, referenceBlockNumber, quorumNumbers, bls));

      // await cedroServiceManger.confirmAgentDeviation((await blsAvsMock.getDeviationHeader(), await blsAvsMock.getNonSignerStakeAndSignature(registryCoordinator.target,referenceBlockNumber,quorumNumbers,bls)));
      // await cedroServiceManger.confirmAgentDeviation([[NonSignerStakesAndSignature[0][0]],]);
    });
  });
});

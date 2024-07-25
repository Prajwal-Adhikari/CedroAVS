//Testing of connection with the Eigen Layer

const { assert, expect } = require("chai");
const chai = require("chai");
const chaiAlmost = require("chai-almost");

chai.use(chaiAlmost());
const { ethers, upgrades, config } = require("hardhat");
const { anyValue } = require("@nomicfoundation/hardhat-chai-matchers/withArgs");

const { mine } = require("@nomicfoundation/hardhat-network-helpers");

//contracts
let cedroServiceManager,
  avsDirectoryMock,
  blsAvsMock,
  operatorStateRetriever,
  delegationMock,
  rewardsCoordinatorMock,
  registryCoordinator,
  stakeRegistry,
  indexRegistry,
  blsApkRegistryMock,
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

describe("Cedro Service Manager Tests", function () {
  beforeEach(async function () {
    [admin, operator1, operator2, operator3] = await ethers.getSigners();

    const RoleManagerCore = await ethers.getContractFactory("RoleManagerCore");
    const CedroServiceManager = await ethers.getContractFactory(
      "CedroServiceManager"
    );
    const RegistryCoordinator = await ethers.getContractFactory(
      "RegistryCoordinator"
    );
    const AVSDirectoryMock = await ethers.getContractFactory(
      "AVSDirectoryMock"
    );
    const BLSAvsMock = await ethers.getContractFactory("BLSAvsMock");
    const DelegationMock = await ethers.getContractFactory("DelegationMock");
    const BlsApkRegistryMock = await ethers.getContractFactory(
      "BLSApkRegistryMock"
    );
    const RewardsCoordinatorMock = await ethers.getContractFactory(
      "RewardsCoordinatorMock"
    );
    const IndexRegistry = await ethers.getContractFactory("IndexRegistry");
    const StakeRegistry = await ethers.getContractFactory("StakeRegistry");
    // const PauserRegistry = await ethers.getContractFactory("PauserRegistry");
    const OperatorStateRetriever = await ethers.getContractFactory(
      "OperatorStateRetriever"
    );

    roleManagerCore = await upgrades.deployProxy(RoleManagerCore, [
      admin.address,
    ]);

    console.log("RoleManager core:", roleManagerCore.target);

    // await grantRoles();
    await roleManagerCore.grantRole(
      await roleManagerCore.AVS_ADMIN(),
      admin.address
    );
    await roleManagerCore.grantRole(
      await roleManagerCore.VERIFIER_ADMIN(),
      admin.address
    );

    delegationMock = await DelegationMock.deploy();
    avsDirectoryMock = await AVSDirectoryMock.deploy(delegationMock.target);
    rewardsCoordinatorMock = await RewardsCoordinatorMock.deploy();
    // pauserRegistry = await PauserRegistry.deploy(admin.address, admin.address);
    blsApkRegistryMock = await upgrades.deployProxy(
      BlsApkRegistryMock,
      [roleManagerCore.target],
      {
        initialize: "BLSApkRegistryMockinitialize",
      }
    );

    console.log("Bls apk",blsApkRegistryMock.target);
    indexRegistry = await upgrades.deployProxy(IndexRegistry, [roleManagerCore.target]);
    console.log("IndexRegistry", indexRegistry.target);
    stakeRegistry = await upgrades.deployProxy(StakeRegistry, [
      roleManagerCore.target,
      delegationMock.target,
    ]);
    console.log("StakeRegistry.target",stakeRegistry.target);
    operatorStateRetriever = await OperatorStateRetriever.deploy();

    cedroServiceManager = await upgrades.deployProxy(CedroServiceManager, [
      roleManagerCore.target,
      avsDirectoryMock.target,
      rewardsCoordinatorMock.target,
      stakeRegistry.target,
    ]);

    console.log("cedroservice manager:",cedroServiceManager.target);

    registryCoordinator = await upgrades.deployProxy(RegistryCoordinator, [
      roleManagerCore.target,
      admin.address,
      admin.address,
      cedroServiceManager.target,
      stakeRegistry.target,
      blsApkRegistryMock.target,
      indexRegistry.target,
    ]);

    console.log("Registry coordinator",registryCoordinator.target);

    blsAvsMock = await BLSAvsMock.deploy(
      cedroServiceManager.target,
      operatorStateRetriever.target,
      registryCoordinator.target,
      blsApkRegistryMock.target,
      delegationMock.target,
      stakeRegistry.target
    );

    blsApkRegistryMock
      .connect(admin)
      .setRegisteryCoordinator(registryCoordinator.target);
    indexRegistry
      .connect(admin)
      .setRegisteryCoordinator(registryCoordinator.target);
    stakeRegistry
      .connect(admin)
      .setRegisteryCoordinator(registryCoordinator.target);
    await cedroServiceManager
      .connect(admin)
      .setRegistryCoordinatorAddress(registryCoordinator.target);

    const blockNumber = await ethers.provider.getBlockNumber();
    block = await ethers.provider.getBlock(blockNumber);
  });

  async function grantRoles() {
    await roleManagerCore.grantRole(
      await roleManagerCore.AVS_ADMIN(),
      admin.address
    );
  }
  async function createQuorum() {
    const numQuorumsToAdd = 1;

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

    const WEIGHTING_DIVISOR = 1e12;
    const quorumStrategiesConsideredAndMultipliers = [];
    for (let i = 0; i < numQuorumsToAdd; i++) {
      strategyId = operator1.address;

      // Here, you would typically call a method on the contract to add a strategy
      // Since we're mocking, we just push the strategy details to our array
      quorumStrategiesConsideredAndMultipliers[i] = [
        [strategyId, WEIGHTING_DIVISOR],
      ];

      for (let i = 0; i < operatorSetParams.length; i++) {
        await registryCoordinator.createQuorum(
          operatorSetParams[i],
          minimumStakeForQuorum[i],
          quorumStrategiesConsideredAndMultipliers[i]
        );
      }
    }
  }

  async function registerOperatorOne() {
    const defaultSocket = "69.69.69.69:420";
    await delegationMock.setIsOperator(operator1.address, true);
    await cedroServiceManager.addOperatorsToAllowlist([operator1.address]);
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
        cedroServiceManager.target,
        salt,
        expiry
      );
    signature = signingKey.sign(digestHash);
    signatureConcat = ethers.Signature.from(signature).serialized;
    let quorumNumbers = await blsAvsMock.converBitmapToBytesArray(1);

    blsPrivateKeyOne = 853;

    let pubkeyRegistrationSignature = await blsAvsMock.signMessage(
      blsPrivateKeyOne,
      operator1.address,
      registryCoordinator.target
    );

    let pubkeysignG1x =
      21008703310028844141620162130587521182514488709533704622337026710714836074802n;
    let pubkeysignG1y =
      1145283673551783492487423178395020830991754869805176606412560778073024405948n;
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
    let operatorSignatureParams = [signatureConcat, salt, expiry];
    await blsApkRegistryMock.setBLSPublicKey(operator1, [pubkeyG1x, pubkeyG1y]);

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
      defaultSocket,
      [
        [pubkeysignG1x, pubkeysignG1y],
        [pubkeyG1x, pubkeyG1y],
        [pubkeyG2x, pubkeyG2y],
      ],
      operatorSignatureParams
    );
  }

  async function registerOperatorTwo() {
    const defaultSocket = "69.69.69.69:420";
    await delegationMock.setIsOperator(operator2.address, true);
    await cedroServiceManager.addOperatorsToAllowlist([operator2.address]);
    privateKey =
      "0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a";
    wallet = new ethers.Wallet(privateKey);
    signingKey = new ethers.SigningKey(wallet.privateKey);
    let salt = ethers.keccak256(
      ethers.solidityPackedKeccak256(["string"], ["default"])
    );
    let expiry = block.timestamp + 100;
    let digestHash =
      await avsDirectoryMock.calculateOperatorAVSRegistrationDigestHash(
        operator2.address,
        cedroServiceManager.target,
        salt,
        expiry
      );
    signature = signingKey.sign(digestHash);
    signatureConcat = ethers.Signature.from(signature).serialized;
    let quorumNumbers = await blsAvsMock.converBitmapToBytesArray(1);

    let pubkeysignG1x =
      2668383907684820844565850416859618418902109493143488579964360725317467055487n;
    let pubkeysignG1y =
      17927264696147669850570728845637944087636679572365331845288700107037638136136n;
    let pubkeyG1x =
      11909330550897052025764217059883091939734032688455609568393118274846029799265n;
    let pubkeyG1y =
      18867945654008209292625458880315629237005264805214925261372079909097324926955n;
    let pubkeyG2x = [
      9400612866821152435661064963197029962476299469519836188003852408000559552133n,
      11833437650394266069432899676145012972823865862729442346860339009616152430781n,
    ];
    let pubkeyG2y = [
      15242336091651110087316026339756071637335622513479676657962837035471261121292n,
      15043272256647926691811884987704334488846179895139852981679303523937946859278n,
    ];

    let operatorSignatureParams = [signatureConcat, salt, expiry];
    await blsApkRegistryMock.setBLSPublicKey(operator2, [pubkeyG1x, pubkeyG1y]);

    let weight = ethers.parseEther("1.0");
    defaultStake = 1;
    let WEIGHTING_DIVISORS = BigInt(1e18);

    // _setOperatorWeight(operator, quorumNumbers[i], defaultStake);
    let strategyDetails = await stakeRegistry.strategyParams(0, 0);
    let actualWeight = (weight * WEIGHTING_DIVISORS) / WEIGHTING_DIVISORS;
    delegationMock.setOperatorShares(
      operator2.address,
      strategyDetails[0],
      actualWeight
    );

    await registryCoordinator.connect(operator2).registerOperator(
      quorumNumbers,
      defaultSocket,
      [
        [pubkeysignG1x, pubkeysignG1y],
        [pubkeyG1x, pubkeyG1y],
        [pubkeyG2x, pubkeyG2y],
      ],
      operatorSignatureParams
    );
  }

  describe("1. Testing CedroServiceManager Functionality ", () => {
    it("Creating Quorum in AVS", async () => {
      const numQuorumsToAdd = 2;

      let minimumStakeForQuorum = new Array(numQuorumsToAdd);

      let operatorSetParams = new Array(numQuorumsToAdd);

      for (let i = 0; i < numQuorumsToAdd; i++) {
        minimumStakeForQuorum[i] = BigInt(i + 1);
      }

      for (let i = 0; i < numQuorumsToAdd; i++) {
        operatorSetParams[i] = [10, 15000, 150];
      }

      for (let i = 0; i < numQuorumsToAdd.length; i++) {
        minimumStakeForQuorum[i] = BigInt(i + 1);
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

      for (let i = 0; i < numQuorumsToAdd; i++) {
        await registryCoordinator
          .connect(admin)
          .createQuorum(
            operatorSetParams[i],
            minimumStakeForQuorum[i],
            quorumStrategiesConsideredAndMultipliers[i]
          );
      }
      await ethers.provider.send("evm_increaseTime", [100]);
      await ethers.provider.send("evm_mine");

      expect(await registryCoordinator.quorumCount()).to.equal(2);
    });

    it("Register operator and enroll in the quorum", async () => {
      const defaultSocket = "69.69.69.69:420";
      await delegationMock.setIsOperator(operator1.address, true);
      await cedroServiceManager.addOperatorsToAllowlist([operator1.address]);

      const numQuorumsToAdd = 2; // Example value

      let minimumStakeForQuorum = new Array(numQuorumsToAdd);

      let operatorSetParams = new Array(numQuorumsToAdd);

      for (let i = 0; i < numQuorumsToAdd; i++) {
        minimumStakeForQuorum[i] = BigInt(i + 1);
      }

      for (let i = 0; i < numQuorumsToAdd; i++) {
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

      for (let i = 0; i < numQuorumsToAdd; i++) {
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
          cedroServiceManager.target,
          salt,
          expiry
        );
      signature = signingKey.sign(digestHash);
      signatureConcat = ethers.Signature.from(signature).serialized;
      let quorumNumbers = await blsAvsMock.converBitmapToBytesArray(1);

      blsPrivateKeyOne = 853;

      let pubkeyRegistrationSignature = await blsAvsMock.signMessage(
        blsPrivateKeyOne,
        operator1.address,
        registryCoordinator.target
      );

      let pubkeysignG1x =
        21008703310028844141620162130587521182514488709533704622337026710714836074802n;
      let pubkeysignG1y =
        1145283673551783492487423178395020830991754869805176606412560778073024405948n;
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
      let operatorSignatureParams = [signatureConcat, salt, expiry];
      await blsApkRegistryMock.setBLSPublicKey(operator1, [
        pubkeyG1x,
        pubkeyG1y,
      ]);

      let weight = ethers.parseEther("1.0");
      defaultStake = 1;
      let WEIGHTING_DIVISORS = BigInt(1e18);

      //  _setOperatorWeight(operator, quorumNumbers[i], defaultStake);
      let strategyDetails = await stakeRegistry.strategyParams(0, 0);
      let actualWeight = (weight * WEIGHTING_DIVISORS) / WEIGHTING_DIVISORS;
      delegationMock.setOperatorShares(
        operator1.address,
        strategyDetails[0],
        actualWeight
      );

      await registryCoordinator.connect(operator1).registerOperator(
        quorumNumbers,
        defaultSocket,
        [
          [pubkeysignG1x, pubkeysignG1y],
          [pubkeyG1x, pubkeyG1y],
          [pubkeyG2x, pubkeyG2y],
        ],
        operatorSignatureParams
      );
    });
  });
});

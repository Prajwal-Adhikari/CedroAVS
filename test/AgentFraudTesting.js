//Testing of connection with the Eigen Layer
const { assert, expect } = require("chai");
const chai = require("chai");
const chaiAlmost = require("chai-almost");

chai.use(chaiAlmost());
const { ethers, upgrades, config } = require("hardhat");
const { anyValue } = require("@nomicfoundation/hardhat-chai-matchers/withArgs");

const { mine } = require("@nomicfoundation/hardhat-network-helpers");
let cedroServiceManager,
  agentManager,
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
    const AgentManager = await ethers.getContractFactory("AgentManagerMock");
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

    console.log("Bls apk", blsApkRegistryMock.target);
    indexRegistry = await upgrades.deployProxy(IndexRegistry, [
      roleManagerCore.target,
    ]);
    console.log("IndexRegistry", indexRegistry.target);
    stakeRegistry = await upgrades.deployProxy(StakeRegistry, [
      roleManagerCore.target,
      delegationMock.target,
    ]);
    console.log("StakeRegistry.target", stakeRegistry.target);
    operatorStateRetriever = await OperatorStateRetriever.deploy();
    agentManager = await AgentManager.deploy();

    cedroServiceManager = await upgrades.deployProxy(CedroServiceManager, [
      roleManagerCore.target,
      avsDirectoryMock.target,
      rewardsCoordinatorMock.target,
      stakeRegistry.target,
    ]);

    console.log("cedroservice manager:", cedroServiceManager.target);
    cedroServiceManager.setAgentManagerAddress(agentManager.target);
    agentManager.setDeviationPercentage(
      "0xdF1Cf915dB7A658dFDEb75aF3873B4e55F25Aa42",
      1e18
    );

    registryCoordinator = await upgrades.deployProxy(RegistryCoordinator, [
      roleManagerCore.target,
      admin.address,
      admin.address,
      cedroServiceManager.target,
      stakeRegistry.target,
      blsApkRegistryMock.target,
      indexRegistry.target,
    ]);

    console.log("Registry coordinator", registryCoordinator.target);

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

  async function registerOperatorOne(quorumBitmap) {
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
    let quorumNumbers = await blsAvsMock.converBitmapToBytesArray(quorumBitmap);

    blsPrivateKeyOne =
      15436643226594533498284352790494521296965667247196462338996131163212826097054n;

    let pubkeyRegistrationSignature = await blsAvsMock.signMessage(
      blsPrivateKeyOne,
      operator1.address,
      registryCoordinator.target
    );

    let pubkeysignG1x = pubkeyRegistrationSignature[0];
    let pubkeysignG1y = pubkeyRegistrationSignature[1];
    let pubkeyG1x =
      14110567285676345389244725775337526580401753248023284461771486167266462355103n;
    let pubkeyG1y =
      9720413144746361199461732946893043367171487505859650901952835867397436923362n;
    let pubkeyG2x = [
      1230730797123182311850807131235073178385821783983372276639854957091263955384n,
      7034359866883434113182571131287568261842779761851733188885430766201745859762n,
    ];
    let pubkeyG2y = [
      17453839935845982565550211636126433446255050520010422383769263780951237157828n,
      16271686786459442426557467568799518089585926782040698082337113377336362878238n,
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
    blsPrivateKeyTwo =
      19502352088942698780298462309344617947373182784327086247271244087233782943208n;

    let pubkeyRegistrationSignature = await blsAvsMock.signMessage(
      blsPrivateKeyTwo,
      operator2.address,
      registryCoordinator.target
    );

    let pubkeysignG1x = pubkeyRegistrationSignature[0];
    let pubkeysignG1y = pubkeyRegistrationSignature[1];

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
    it.only("Test the signature verification flow", async () => {
      await createQuorum();
      await createQuorum();
      let referenceBlockNumbers = 1866122;
      let quorumNumbers = await blsAvsMock.converBitmapToBytesArray(1);
      let quorumBitmap = 1;
      await registerOperatorOne(quorumBitmap);

      await registerOperatorTwo();
      let nonSigneroperatorAddress = [operator2.address];
      let NonSignerStakesAndSignature =
        await blsAvsMock.getFraudNonSignerStakeAndSignature(
          referenceBlockNumbers,
          quorumBitmap,
          nonSigneroperatorAddress
        );

      console.log(NonSignerStakesAndSignature);
      await mine(1866122);
      await blsAvsMock
        .connect(admin)
        .testFraud(
          referenceBlockNumbers,
          quorumBitmap,
          nonSigneroperatorAddress
        );
    });
  });
});

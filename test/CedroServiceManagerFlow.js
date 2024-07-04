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

    it.only("Creating Quorum in AVS", async () => {
      const numQuorumsToAdd = 1; // Example value

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
          console.log("ewr");
          await registryCoordinator.createQuorum(
            operatorSetParams[i],
            minimumStakeForQuorum[i],
            quorumStrategiesConsideredAndMultipliers[i]
          );
        }
      }
      await ethers.provider.send("evm_increaseTime", [100]);
      await ethers.provider.send("evm_mine");

      console.log(await registryCoordinator.quorumCount());

      console.log("params:",await stakeRegistry.strategyParamsByIndex(0,0));
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
      let quorumNumbers = await blsAvsMock.converBitmapToBytesArray(2);
   

      blsPrivateKey = 853;

      let pubkeyRegistrationSignature = await blsAvsMock.signMessage(
        blsPrivateKeyOne,
        operator1.address,
        registryCoordinator.target
      );
      // let pubkeyRegistrationSignature =[13540890549107538365317826005663225429149650894201009774295406557850371632651n,2966724547965158980192270741744292034179150792486555967152513253931665577541n]
      // let pubkeyG1 = await blsAvsMock.getG1PubKey(blsPrivateKeyOne);
      console.log("**Main Public key**", pubkeyG1);
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

    it("Test the signature verification flow", async () => {
      await delegationMock.setIsOperator(operator1.address, true);
      await cedroServiceManger.addOperatorsToAllowlist([operator1.address]);
      await delegationMock.setIsOperator(operator2.address, true);
      await cedroServiceManger.addOperatorsToAllowlist([operator2.address]);

      console.log("Oper2 Address", operator2.address);

      const numQuorumsToAdd = 2; // Example value

      let minimumStakeForQuorum = new Array(numQuorumsToAdd);

      let operatorSetParams = new Array(numQuorumsToAdd);

      for (let i = 0; i < minimumStakeForQuorum.length; i++) {
        minimumStakeForQuorum[i] = BigInt(i + 1);
      }

      for (let i = 0; i < operatorSetParams.length; i++) {
        operatorSetParams[i] = [10, 15000, 150];
      }

      const WEIGHTING_DIVISOR = 1e10;
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
      privateKeyTwo =
        "0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a";

      wallet = new ethers.Wallet(privateKey);
      signingKey = new ethers.SigningKey(wallet.privateKey);

      wallet2 = new ethers.Wallet(privateKeyTwo);
      signingKey2 = new ethers.SigningKey(wallet2.privateKey);

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
      let digestHashTwo =
        await avsDirectoryMock.calculateOperatorAVSRegistrationDigestHash(
          operator2.address,
          cedroServiceManger.target,
          salt,
          expiry
        );
      signature = signingKey.sign(digestHash);
      signatureConcat = ethers.Signature.from(signature).serialized;

      signatureTwo = signingKey2.sign(digestHashTwo);
      signatureConcatTwo = ethers.Signature.from(signatureTwo).serialized;

      let quorumNumbers = await blsAvsMock.converBitmapToBytesArray(1);
      console.log("quorum Length : ", quorumNumbers.length);

      blsPrivateKeyOne = 853;
      blsPrivateKeyTwo = 655;

      // let pubkeyRegistrationSignature = await blsAvsMock.signMessage(
      //   blsPrivateKeyTwo,
      //   operator2.address,
      //   registryCoordinator.target
      // );
      // let pubkeyG1 = await blsAvsMock.getG1PubKey(blsPrivateKeyTwo);
      // console.log("**Main Public key**", pubkeyG1);
      let pubkey1SigG1x =
        17409514112937982184884384646770072864005066810947009403975675305346550516857n;
      let pubkey1SigG1y =
        16378308196958834164882488914719045785044027868371311005910777189078541893110n;

      let pubkey1G1x =
        21008703310028844141620162130587521182514488709533704622337026710714836074802n;
      let pubkey1G1y =
        1145283673551783492487423178395020830991754869805176606412560778073024405948n;

      let pubkey2SigG1x =
        13312587445422295126177139014651984144514712086661404227691088207480792536745n;
      let pubkey2SigG1y =
        15363240117222918788648648388562209952076267180616555977408454637004630486090n;

      let pubkey2G1x =
        9424467304309542109140998916676104359278524971990760318275336488052797649388n;
      let pubkey2G1y =
        1778726798519900263032462255183166302989403596381730183535085586499987568162n;

      let pubkeyG2x = [
        12302775405227555806769733102102076488786351091672743244822361848998753983454n,
        12302775405227555806769733102102076488786351091672743244822361848998753983454n,
      ];
      let pubkeyG2y = [
        12302775405227555806769733102102076488786351091672743244822361848998753983454n,
        12302775405227555806769733102102076488786351091672743244822361848998753983454n,
      ];

      let operatorSignatureParams = [signatureConcat, salt, expiry];
      let operatorSignatureParamsTwo = [signatureConcatTwo, salt, expiry];
      await blsApkRegistryMock.setBLSPublicKey(operator1.address, [
        pubkey1G1x,
        pubkey1G1y,
      ]);
      await blsApkRegistryMock.setBLSPublicKey(operator2.address, [
        pubkey2G1x,
        pubkey2G1y,
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
      delegationMock.setOperatorShares(
        operator2.address,
        strategyDetails[0],
        actualWeight
      );

      await registryCoordinator.connect(operator1).registerOperator(
        quorumNumbers,
        [
          [pubkey1SigG1x, pubkey1SigG1y],
          [pubkey1G1x, pubkey1G1y],
          [pubkeyG2x, pubkeyG2y],
        ],
        operatorSignatureParams
      );

      await registryCoordinator.connect(operator2).registerOperator(
        quorumNumbers,
        [
          [pubkey2SigG1x, pubkey2SigG1y],
          [pubkey2G1x, pubkey2G1y],
          [pubkeyG2x, pubkeyG2y],
        ],
        operatorSignatureParamsTwo
      );
      /* ==> De register Operator <== */
      // await registryCoordinator.connect(operator1).deregisterOperator(quorumNumbers);

      /* ==> confirming signature Check <== */
      let nonRandomNumber = 111;
      let numNonSigners = 0;
      let quorumBitmap = 1;
      let referenceBlockNumbers = 123;
      let nonsigenrOperatorIds = await blsApkRegistryMock.getOperatorId(
        operator2.address
      );
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

      let NonSignerStakesAndSignature =
        await blsAvsMock.getNonSignerStakeAndSignature(
          registryCoordinator.target,
          referenceBlockNumbers,
          quorumNumbers,
          [nonsigenrOperatorIds]
        );

      let nonSignerQuorumBitmapIndicesData = [
        NonSignerStakesAndSignature[0][0],
      ];

      let nonSignerPubkeysData = [
        [
          NonSignerStakesAndSignature[1][0].X,
          NonSignerStakesAndSignature[1][0].Y,
        ],
      ];
      console.log("quorum apks:", NonSignerStakesAndSignature[2]);
      let quorumApksData = [
        [
          NonSignerStakesAndSignature[2][0].X,
          NonSignerStakesAndSignature[2][0].Y,
        ]
      ];
      console.log("quorum apks Final:", quorumApksData);
     

      let apkG2Data = [
        [
          NonSignerStakesAndSignature[3].X[0],
          NonSignerStakesAndSignature[3].X[1],
        ],
        [
          NonSignerStakesAndSignature[3].Y[0],
          NonSignerStakesAndSignature[3].Y[1],
        ],
      ];

      let sigmaDetailsData = [
        NonSignerStakesAndSignature[4].X,
        NonSignerStakesAndSignature[4].Y,
      ];

      let quorumApkIndicesData = [NonSignerStakesAndSignature[5][0]];
      console.log(" quorumApkIndicesData",NonSignerStakesAndSignature[5])

      let totalStakeIndicesData = [NonSignerStakesAndSignature[6][0]];

      let nonSignerStakeIndices = [[NonSignerStakesAndSignature[7][0][0]]];

      let nonSignerStakeAndSignatureData = [
        nonSignerQuorumBitmapIndicesData,
        nonSignerPubkeysData,
        quorumApksData,
        apkG2Data,
        sigmaDetailsData,
        quorumApkIndicesData,
        totalStakeIndicesData,
        nonSignerStakeIndices,
      ];
      // console.log(NonSignerStakeAndSignatureData);
      await cedroServiceManger.confirmAgentDeviation(
        deviationHeaderData,
        nonSignerStakeAndSignatureData
      );
    });
  });
});

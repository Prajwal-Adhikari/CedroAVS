const hre = require("hardhat");

let { ethers, BigNumber } = require("hardhat");
avsDirectory = "0x055733000064333CaDDbC92763c58BF0192fFeBf";
operatorAddress = "0x394204E5ad767C87Ed47e12BFc1f09E6aC1f0FfD";

let deployed = require("../configurations/deployedAddresses.json");
async function main() {
  this.roleManagerCore = await hre.ethers.getContractAt(
    "RegistryCoordinator",
    deployed.roleManagerCore
  );
  this.blsApkRegistry = await hre.ethers.getContractAt(
    "BLSApkRegistry",
    deployed.blsApkRegistry
  );
  this.indexRegistry = await hre.ethers.getContractAt(
    "IndexRegistry",
    deployed.indexRegistry
  );
  this.stakeRegistry = await hre.ethers.getContractAt(
    "StakeRegistry",
    deployed.stakeRegistry
  );
  this.registryCoordinator = await hre.ethers.getContractAt(
    "RegistryCoordinator",
    deployed.registryCoordinator
  );
  this.cedroServiceManger = await hre.ethers.getContractAt(
    "CedroServiceManager",
    deployed.cedroServiceManager
  );
  this.avsDirectory = await hre.ethers.getContractAt(
    "IAVSDirectory",
    avsDirectory
  );
  privateKey =
    "7868afc56ff43938d92893dc217af0ce3b5124862cd7381505f19e5936d2cf90";
  wallet = new ethers.Wallet(privateKey);
  signingKey = new ethers.SigningKey(wallet.privateKey);
  let salt = ethers.keccak256(
    ethers.solidityPackedKeccak256(["string"], ["default"])
  );
  const blockNumber = await ethers.provider.getBlockNumber();
  block = await ethers.provider.getBlock(blockNumber);
  let expiry = block.timestamp + 10000000;
  let digestHash =
    await this.avsDirectory.calculateOperatorAVSRegistrationDigestHash(
      operatorAddress,
      cedroServiceManger.target,
      salt,
      expiry
    );
  signature = signingKey.sign(digestHash);
  signatureConcat = ethers.Signature.from(signature).serialized;
  let quorumNumbers = "0x00";
  let pubkeysignG1x = 5589022125577323934391431522572184039754807391478425702365075567092427810701n;
  let pubkeysignG1y = 4129025964625052252918938246113487525405756480484787245066509899697380762111n;
  let pubkeyG1x = 5423693149551641599464266040516750971468077004404591118057968424978252918921n;
  let pubkeyG1y = 15921742250186927282153244655129763748330272386599071944264681425191837573957n;
  let pubkeyG2x = [
    18622777961082427940629631196160544718932092620862893173939930284569768749369n,
    17531172364552774445360311089534592256177587830728359887801399947701128718877n,
  ];
  let pubkeyG2y = [
    14272664704219643683273900845860826666400648937855680261109840650498411827454n,
    15580571232102728152592263831513954367224122288182475756457719875787640640905n,
  ];
  // pubkeyRegistrationParams = [
  //   [pubkeysignG1x, pubkeysignG1y],
  //   [pubkeyG1x, pubkeyG1y],
  //   [pubkeyG2x, pubkeyG2y],
  // ];
  let operatorSignatureParams = [signatureConcat, salt, expiry];
  await registryCoordinator.registerOperator(
    quorumNumbers,
    [
      [pubkeysignG1x, pubkeysignG1y],
      [pubkeyG1x, pubkeyG1y],
      [pubkeyG2x, pubkeyG2y],
    ],
    operatorSignatureParams
  );
  // console.log(quorumNumbers);
}

main().catch((error) => {
  console.error(error);
});

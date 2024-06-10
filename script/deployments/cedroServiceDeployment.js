const { ethers, upgrades, config } = require("hardhat");

async function main() {
  const delegationManager = "0xA44151489861Fe9e3055d95adC98FbD462B948e7";
  const avsDirectory = "0x055733000064333CaDDbC92763c58BF0192fFeBf";
  const adminAddress = "0x19B78FF82C94b5E517f2279f3fBF10498B039179";
  const CedroServiceManager = await ethers.getContractFactory(
    "CedroServiceManager"
  );
  const RegistryCoordinator = await ethers.getContractFactory(
    "RegistryCoordinator"
  );

  const AVSDirectoryMock = await ethers.getContractFactory("AVSDirectoryMock");
  const PaymentCoordinatorMock = await ethers.getContractFactory(
    "PaymentCoordinatorMock"
  );
  const BlsApkRegistry = await ethers.getContractFactory("BLSApkRegistry");
  const IndexRegistry = await ethers.getContractFactory("IndexRegistry");
  const StakeRegistry = await ethers.getContractFactory("StakeRegistry");

  const OperatorStateRetriever = await ethers.getContractFactory(
    "OperatorStateRetriever"
  );
  const RoleManagerCore = await ethers.getContractFactory("RoleManagerCore");

  this.roleManagerCore = await upgrades.deployProxy(RoleManagerCore, [
    adminAddress,
  ]);
  console.log("Role manager Core Address : ", this.roleManagerCore.target);

  this.blsApkRegistry = await BlsApkRegistry.deploy(adminAddress);
  console.log("Bls Apk Registry Address : ", this.blsApkRegistry.target);

  this.indexRegistry = await IndexRegistry.deploy(adminAddress);
  console.log("Index Registry AAddress : ", this.indexRegistry.target);

  this.stakeRegistry = await StakeRegistry.deploy(
    adminAddress,
    delegationManager
  );
  console.log("Stake Registry Address : ", this.stakeRegistry.target);

  this.operatorStateRetriever = await OperatorStateRetriever.deploy();
  console.log(
    "Operator State Retriever Address : ",
    this.operatorStateRetriever.target
  );
  this.cedroServiceManger = await upgrades.deployProxy(CedroServiceManager, [
    this.roleManagerCore.target,
    avsDirectory,
    avsDirectory, //as holesky testnet data not found
    this.stakeRegistry.target,
  ]);
}

main().catch((error) => {
  console.error(error);
});

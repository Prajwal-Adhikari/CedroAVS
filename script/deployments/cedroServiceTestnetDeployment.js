const { ethers, upgrades, config } = require("hardhat");

async function main() {
  const delegationManager = "0xA44151489861Fe9e3055d95adC98FbD462B948e7";
  const avsDirectory = "0x055733000064333CaDDbC92763c58BF0192fFeBf";
  const adminAddress = "0x19B78FF82C94b5E517f2279f3fBF10498B039179";
  const RewardsCoordinatorAddress =
    "0xAcc1fb458a1317E886dB376Fc8141540537E68fE";
  const CedroServiceManager = await ethers.getContractFactory(
    "CedroServiceManager"
  );
  const RegistryCoordinator = await ethers.getContractFactory(
    "RegistryCoordinator"
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

  this.blsApkRegistry = await upgrades.deployProxy(BlsApkRegistry, [
    this.roleManagerCore.target,
  ]);
  console.log("Bls Apk Registry Address : ", this.blsApkRegistry.target);

  this.indexRegistry = await upgrades.deployProxy(IndexRegistry, [
    this.roleManagerCore.target,
  ]);
  console.log("Index Registry AAddress : ", this.indexRegistry.target);

  this.stakeRegistry = await upgrades.deployProxy(StakeRegistry, [
    this.roleManagerCore.target,
    delegationManager,
  ]);
  console.log("Stake Registry Address : ", this.stakeRegistry.target);

  this.operatorStateRetriever = await OperatorStateRetriever.deploy();
  console.log(
    "Operator State Retriever Address : ",
    this.operatorStateRetriever.target
  );
  this.cedroServiceManger = await upgrades.deployProxy(CedroServiceManager, [
    this.roleManagerCore.target,
    avsDirectory,
    RewardsCoordinatorAddress,
    this.stakeRegistry.target,
  ]);
  console.log("Cedro Service Manager Address:", this.cedroServiceManger.target);

  this.registryCoordinator = await upgrades.deployProxy(RegistryCoordinator, [
    this.roleManagerCore.target,
    adminAddress,
    adminAddress,
    this.cedroServiceManger.target,
    this.stakeRegistry.target,
    this.blsApkRegistry.target,
    this.indexRegistry.target,
  ]);
  console.log("Registry Coordinator Address:", this.registryCoordinator.target);
}

main().catch((error) => {
  console.error(error);
});

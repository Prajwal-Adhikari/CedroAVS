const hre = require("hardhat");

let { ethers } = require("hardhat");

let deployed = require("../configurations/deployedAddresses.json");

adminAddress = "0x394204e5ad767c87ed47e12bfc1f09e6ac1f0ffd";

async function main() {
  this.roleManagerCore = await hre.ethers.getContractAt(
    "RoleManagerCore",
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

  await this.roleManagerCore.grantRole(
    await this.roleManagerCore.AVS_ADMIN(),
    adminAddress
  );

  await this.roleManagerCore.grantRole(
    await this.roleManagerCore.VERIFIER_ADMIN(),
    adminAddress
  );

}
main().catch((error) => {
  console.error(error);
});

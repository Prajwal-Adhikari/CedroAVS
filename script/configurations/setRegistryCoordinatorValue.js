const hre = require("hardhat");

let { ethers } = require("hardhat");

let deployed = require("./deployedAddresses.json");
let registryCoordinator = "0xbc5DB5e02F8ff5C2dBB69784e318348d3d85E2e7";
let adminAddress = "0x19B78FF82C94b5E517f2279f3fBF10498B039179";

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
  this.cedroServiceManager = await hre.ethers.getContractAt(
    "CedroServiceManager",
    deployed.cedroServiceManager
  );

  // await this.blsApkRegistry.setRegisteryCoordinator(registryCoordinator);
  await this.indexRegistry.setRegisteryCoordinator(registryCoordinator);
  await this.stakeRegistry.setRegisteryCoordinator(registryCoordinator);
  await this.cedroServiceManager.setRegistryCoordinatorAddress(
    registryCoordinator
  );
  // await this.cedroServiceManager.setRegistryCoordinator(registryCoordinator);
  console.log("Value set");

  // console.log(await this.blsApkRegistry.registryCoordinator());
  // console.log(await this.indexRegistry.registryCoordinator());
  // console.log(await this.stakeRegistry.registryCoordinator());
  // console.log(await this.cedroServiceManager.registryCoordinator());
}
main().catch((error) => {
  console.error(error);
});

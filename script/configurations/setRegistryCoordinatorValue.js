const hre = require("hardhat");

let { ethers } = require("hardhat");

let deployed = require("./deployedAddresses.json");
let registryCoordinator = "0x17db11dAE61D0d9c66717bCac08e3B146d07B4fF";

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

  // await this.blsApkRegistry.setRegisterCoordinator(registryCoordinator);
  // await this.indexRegistry.setRegisterCoordinator(registryCoordinator);
  // await this.stakeRegistry.setRegisterCoordinator(registryCoordinator);
  // await this.cedroServiceManager.setRegistryCoordinatorAddress(
  //   registryCoordinator
  // );
  // await this.cedroServiceManager.setRegistryCoordinator(registryCoordinator);

  // console.log(await this.blsApkRegistry.registryCoordinator());
  // console.log(await this.indexRegistry.registryCoordinator());
  // console.log(await this.stakeRegistry.registryCoordinator());
  console.log(await this.cedroServiceManager.registryCoordinator());
}
main().catch((error) => {
  console.error(error);
});

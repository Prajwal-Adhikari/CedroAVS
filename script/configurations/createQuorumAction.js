const hre = require("hardhat");

let { ethers } = require("hardhat");

let deployed = require("../configurations/deployedAddresses.json");
// let registryCoordinator = "0xCD79AC7615f7C5e4dfc90e86F5D5c2bC66Fd184A";

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

//   const numQuorumsToAdd = 1;

//   let minimumStakeForQuorum = new Array(numQuorumsToAdd);

//   let operatorSetParams = new Array(numQuorumsToAdd);

//   for (let i = 0; i < minimumStakeForQuorum.length; i++) {
//     minimumStakeForQuorum[i] = BigInt(i + 1);
//   }

//   for (let i = 0; i < operatorSetParams.length; i++) {
//     operatorSetParams[i] = [10, 15000, 150];
//   }

//   for (let i = 0; i < numQuorumsToAdd.length; i++) {
//     minimumStakeForQuorum[i] = BigInt(i + 1);
//   }

//   const WEIGHTING_DIVISOR = 1000000000000000000n;
//   const quorumStrategiesConsideredAndMultipliers = [];
//   for (let i = 0; i < numQuorumsToAdd; i++) {
//     strategyId = "0xbeaC0eeEeeeeEEeEeEEEEeeEEeEeeeEeeEEBEaC0";

//     quorumStrategiesConsideredAndMultipliers[i] = [
//       [strategyId, WEIGHTING_DIVISOR],
//     ];

//     for (let i = 0; i < operatorSetParams.length; i++) {
//       await registryCoordinator.createQuorum(
//         operatorSetParams[i],
//         minimumStakeForQuorum[i],
//         quorumStrategiesConsideredAndMultipliers[i]
//       );
//     }
  //   }
  console.log("params:",await stakeRegistry.strategyParamsByIndex(0,0));
}
main().catch((error) => {
  console.error(error);
});

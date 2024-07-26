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

  const numQuorumsToAdd = 1;

  let minimumStakeForQuorum = new Array(numQuorumsToAdd);

  let operatorSetParams = new Array(numQuorumsToAdd);

  for (let i = 0; i < minimumStakeForQuorum.length; i++) {
    minimumStakeForQuorum[i] = BigInt(i + 1);
  }

  for (let i = 0; i < numQuorumsToAdd; i++) {
    operatorSetParams[i] = [10, 15000, 150];
  }

  for (let i = 0; i < numQuorumsToAdd; i++) {
    minimumStakeForQuorum[i] = BigInt(i + 1);
  }

  const WEIGHTING_DIVISOR = 1000000000000000000n;
  const quorumStrategiesConsideredAndMultipliers = [];
  for (let i = 0; i < numQuorumsToAdd; i++) {
    strategyId = "0x7D704507b76571a51d9caE8AdDAbBFd0ba0e63d3";
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

  console.log("Created Quorum")
}
main().catch((error) => {
  console.error(error);
});

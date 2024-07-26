const { ethers,upgrades } = require("hardhat");

async function main() {
 const CedroServiceManagerV2 = await ethers.getContractFactory("CedroServiceManager");
 const CedroServiceManagerV2Deploy = await upgrades.upgradeProxy("0xC78b7146629614280dF3dBffcEefce699636BdA9", CedroServiceManagerV2);
 console.log("Contract upgraded to:", CedroServiceManagerV2Deploy);
}

main().catch((error) => {
 console.error(error);
 process.exitCode = 1;
});




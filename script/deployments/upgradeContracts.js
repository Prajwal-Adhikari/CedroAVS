const { ethers,defender,upgrades } = require("hardhat");

async function main() {
 const CedroServiceManagerV2 = await ethers.getContractFactory("CedroServiceManager");
 const CedroServiceManagerV2Deploy = await upgrades.upgradeProxy("0xdBbE5334CfC0C98397dEf343fD0f29B48AD124DE", CedroServiceManagerV2);
 console.log("Contract upgraded to:", CedroServiceManagerV2Deploy);
}

main().catch((error) => {
 console.error(error);
 process.exitCode = 1;
});




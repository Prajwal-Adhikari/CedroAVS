require("@nomicfoundation/hardhat-toolbox");
require('hardhat-contract-sizer');
require('@openzeppelin/hardhat-upgrades');
require('@nomicfoundation/hardhat-chai-matchers');
require("@nomicfoundation/hardhat-verify");
require('dotenv').config();
const PRIVATE_KEY = process.env.PRIVATE_KEY;
const API_KEY = process.env.API_KEY
module.exports = {
  solidity: {
    version: "0.8.24",
    settings: {
      optimizer: {
        enabled: true,
        runs: 200,
      },
    }
  }, 
  contractSizer: {
    runOnCompile: false,
    alphaSort: false,
    disambiguatePaths: false,
    strictMode: true,
    ignoreFiles: [],
    reportGas: false,
    maxMethodDiff: 10,
    targets: [
      { path: "contracts/server/Pool.sol", name: "Pool" }
    ]
  },
  etherscan: {
    apiKey: 'EJFFCD84CTTH4PDY8Q1JUCKJCNA1VRWBGC',
  },
  sourcify: {
    enabled: true
  },
  networks: {
   holesky: {
      url: `https://rpc.holesky.ethpandaops.io/`,
      accounts: [`0x${PRIVATE_KEY}`],
      gas: 6000000
    },
    mainnet: {
      url: `https://jsonrpc.maiziqianbao.net`,
      accounts: [`0x${PRIVATE_KEY}`],
    },
    bsc: {
      url: `https://bsc-testnet-rpc.publicnode.com`,
      // url: `https://bsc-rpc.publicnode.com`,
      accounts: [`0x${PRIVATE_KEY}`],
      gas: 4000000
    },
    sepolia: {
      url: `https://eth-sepolia.public.blastapi.io`,
      accounts: [`0x${PRIVATE_KEY}`],
      gas: 4000000
    },
   
  }
};

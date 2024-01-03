import * as dotenv from 'dotenv';
import {HardhatUserConfig, task} from 'hardhat/config';
import '@nomiclabs/hardhat-ethers';
import '@nomiclabs/hardhat-etherscan';
import '@nomiclabs/hardhat-waffle';
import '@starboardventures/hardhat-verify';
import '@typechain/hardhat';
import 'hardhat-gas-reporter';
import 'solidity-coverage';
import 'hardhat-deploy';
import 'hardhat-watcher';

dotenv.config();

const infuraToken = process.env.INFURA_TOKEN;
const goerliDeployerPK = process.env.GOERLI_DEPLOYER_PK;
const wallabyDeployerPk = process.env.WALLABY_DEPLOYER_PK;
const calibrationDeployerPk = process.env.CALIBRATION_DEPLOYER_PK;
// This is a sample Hardhat task. To learn how to create your own go to
// https://hardhat.org/guides/create-task.html
task('accounts', 'Prints the list of accounts', async (taskArgs, hre) => {
  const accounts = await hre.ethers.getSigners();

  for (const account of accounts) {
    console.log(account.address);
  }
});

// You need to export an object to set up your config
// Go to https://hardhat.org/config/ to learn more

const config: HardhatUserConfig & {watcher: any} = {
  solidity: {
    compilers: [
      {
        version: '0.8.20',
        settings: {
          optimizer: {
            enabled: true,
            runs: 20_000,
          },
          viaIR: true,
        },
      },
    ],
  },
  namedAccounts: {
    deployer: {
      default: 0,
    },
  },
  paths: {
    sources: 'contracts',
    deploy: 'hardhat-deploy',
    deployments: 'hardhat-deployments',
  },
  watcher: {
    compilation: {
      tasks: ['compile'],
      verbose: true,
    },
  },
  networks: {
    hardhat: {
      chainId: 1337,
    },
    goerli: {
      url: infuraToken ? 'https://goerli.infura.io/v3/' + infuraToken : '',
      accounts: goerliDeployerPK ? [goerliDeployerPK] : [],
      chainId: 5,
    },
    wallaby: {
      url: 'https://wallaby.node.glif.io/rpc/v0',
      accounts: wallabyDeployerPk ? [wallabyDeployerPk] : [],
      chainId: 31415,
    },
    hyperspace: {
      url: 'https://api.hyperspace.node.glif.io/rpc/v1',
      accounts: wallabyDeployerPk ? [wallabyDeployerPk] : [],
      chainId: 3141,
    },
    calibration: {
      url: 'https://api.calibration.node.glif.io/rpc/v1',
      accounts: calibrationDeployerPk ? [calibrationDeployerPk] : [],
      chainId: 314159,
    },
  },
  starboardConfig: {
    baseURL: 'https://fvm-calibration-api.starboard.ventures',
    network: 'Calibration', // if there's no baseURL, url will depend on the network.  Mainnet || Calibration
  },
};

export default config;

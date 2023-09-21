import 'hardhat-deploy';
import 'hardhat-deploy-ethers';

import {HardhatRuntimeEnvironment} from 'hardhat/types';

module.exports = async (hre: HardhatRuntimeEnvironment) => {
  const {deployments, getNamedAccounts, getChainId} = hre;
  const {deploy} = deployments;
  const {deployer} = await getNamedAccounts();

  console.log('Working on chain id #', await getChainId());
  console.log('deployer', deployer);

  try {
    await deploy('NitroAdjudicator', {
      from: deployer,
      args: [],
      // since Ethereum's legacy transaction format is not supported on FVM, we need to specify
      // maxPriorityFeePerGas to instruct hardhat to use EIP-1559 tx format
      maxPriorityFeePerGas: hre.ethers.BigNumber.from(1500000000),
      skipIfAlreadyDeployed: false,
      log: true,
    });
  } catch (err) {
    const msg = err instanceof Error ? err.message : JSON.stringify(err);
    console.error(`Error when deploying contract: ${msg}`);
  }

  try {
    await deploy('ConsensusApp', {
      from: deployer,
      args: [],
      // since Ethereum's legacy transaction format is not supported on FVM, we need to specify
      // maxPriorityFeePerGas to instruct hardhat to use EIP-1559 tx format
      maxPriorityFeePerGas: hre.ethers.BigNumber.from(1500000000),
      skipIfAlreadyDeployed: false,
      log: true,
    });
  } catch (err) {
    const msg = err instanceof Error ? err.message : JSON.stringify(err);
    console.error(`Error when deploying contract: ${msg}`);
  }

  try {
    await deploy('VirtualPaymentApp', {
      from: deployer,
      args: [],
      // since Ethereum's legacy transaction format is not supported on FVM, we need to specify
      // maxPriorityFeePerGas to instruct hardhat to use EIP-1559 tx format
      maxPriorityFeePerGas: hre.ethers.BigNumber.from(1500000000),
      skipIfAlreadyDeployed: false,
      log: true,
    });
  } catch (err) {
    const msg = err instanceof Error ? err.message : JSON.stringify(err);
    console.error(`Error when deploying contract: ${msg}`);
  }
};
module.exports.tags = ['deploy-fvm'];

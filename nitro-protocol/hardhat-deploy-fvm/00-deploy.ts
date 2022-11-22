import 'hardhat-deploy';
import 'hardhat-deploy-ethers';

import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {BigNumber} from 'ethers';

module.exports = async (hre: HardhatRuntimeEnvironment) => {
  const {deployments, getNamedAccounts, getChainId} = hre;
  const {deploy} = deployments;
  const {deployer} = await getNamedAccounts();

  try {
    console.log('Working on chain id #', await getChainId());

    await deploy('NitroAdjudicator', {
      from: deployer,
      args: [],
      // since it's difficult to estimate the gas limit before f4 address is launched, it's safer to manually set
      // a large gasLimit. This should be addressed in the following releases.
      gasLimit: 1000000000, // BlockGasLimit / 10
      // since Ethereum's legacy transaction format is not supported on FVM, we need to specify
      // maxPriorityFeePerGas to instruct hardhat to use EIP-1559 tx format
      maxPriorityFeePerGas: BigNumber.from(1500000000),
      skipIfAlreadyDeployed: false,
      log: true,
    });
  } catch (err) {
    const msg = err instanceof Error ? err.message : JSON.stringify(err);
    console.error(`Error when deploying contract: ${msg}`);
  }
};

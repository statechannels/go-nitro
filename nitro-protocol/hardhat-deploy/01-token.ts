import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  const {
    deployments: {deploy},
    getChainId,
    getNamedAccounts,
  } = hre;

  const {deployer} = await getNamedAccounts();

  await deploy('Token', {
    from: deployer,
    log: true,
    args: [deployer],
    deterministicDeployment: true,
  });
};
export default func;
func.tags = ['Token'];

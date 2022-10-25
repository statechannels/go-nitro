import 'hardhat-deploy';
import 'hardhat-deploy-ethers';

import RpcEngine from '@glif/filecoin-rpc-client';
import fa, {newDelegatedEthAddress} from '@glif/filecoin-address';
import {ethers} from 'hardhat';
import {HttpNetworkConfig} from 'hardhat/types';

module.exports = async (hre: any) => {
  const deploy = hre.deployments.deploy;

  try {
    const config = hre.network.config as HttpNetworkConfig;
    // generate the f1 address equivalent from the same private key
    // note this method of extracting private key from hre might be unsafe...
    const w = new ethers.Wallet((config.accounts as string[])[0]);
    const pubKey = Uint8Array.from(Buffer.from(w.publicKey.slice(2), 'hex'));
    const f1addr = fa.newSecp256k1Address(pubKey).toString();
    console.log('Native actor addr: ', f1addr);
    console.log('Eth addr: ', w.address);

    const filRpc = new RpcEngine({apiAddress: config.url});
    const ethRpc = new RpcEngine({
      apiAddress: config.url,
      namespace: 'eth',
      delimeter: '_',
    });

    const nonce = await filRpc.request('MpoolGetNonce', f1addr);
    console.log('Nonce: ', nonce);
    const priorityFee = await ethRpc.request('maxPriorityFeePerGas');

    const {address} = await deploy('NitroAdjudicator', {
      from: w.address,
      args: [],
      // since it's difficult to estimate the gas limit before f4 address is launched, it's safer to manually set
      // a large gasLimit. This should be addressed in the following releases.
      gasLimit: 1000000000, // BlockGasLimit / 10
      // since Ethereum's legacy transaction format is not supported on FVM, we need to specify
      // maxPriorityFeePerGas to instruct hardhat to use EIP-1559 tx format
      maxPriorityFeePerGas: priorityFee,
      nonce,
      log: true,
    });

    console.log(address, newDelegatedEthAddress(address).toString());
  } catch (err) {
    const msg = err instanceof Error ? err.message : JSON.stringify(err);
    console.error(`Error when deploying contract: ${msg}`);
  }
};

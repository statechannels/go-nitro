require('hardhat-deploy');

const util = require('util');

const ethers = require('ethers');
const fa = require('@glif/filecoin-address');
const request = util.promisify(require('request'));

const wallet = ethers.Wallet.fromMnemonic(network.config.accounts.mnemonic);

function hexToBytes(hex) {
  // ref: https://stackoverflow.com/a/34356351
  for (var bytes = [], c = 0; c < hex.length; c += 2) bytes.push(parseInt(hex.substr(c, 2), 16));
  return new Uint8Array(bytes);
}

async function callRpc(method, params) {
  const options = {
    method: 'POST',
    url: 'https://wallaby.node.glif.io/rpc/v0',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      jsonrpc: '2.0',
      method: method,
      params: params,
      id: 1,
    }),
  };
  const res = await request(options);
  return JSON.parse(res.body).result;
}

module.exports = async ({deployments}) => {
  const {deploy} = deployments;

  const pubKey = hexToBytes(wallet.publicKey.slice(2));
  const f1addr = fa.newSecp256k1Address(pubKey).toString();

  const priorityFee = await callRpc('eth_maxPriorityFeePerGas');
  const nonce = await callRpc('Filecoin.MpoolGetNonce', [f1addr]);

  console.log('Ethereum deployer address:', wallet.address);
  console.log('Filecoin deployer address (f1):', f1addr);
  console.log('Nonce:', nonce);

  await deploy('NitroAdjudicator', {
    from: wallet.address,
    args: [],
    // since it's difficult to estimate the gas limit before f4 address is launched, it's safer to manually set
    // a large gasLimit. This should be addressed in the following releases.
    gasLimit: 1000000000, // BlockGasLimit / 10
    // since Ethereum's legacy transaction format is not supported on FVM, we need to specify
    // maxPriorityFeePerGas to instruct hardhat to use EIP-1559 tx format
    maxPriorityFeePerGas: priorityFee,
    nonce: nonce,
    log: true,
  });

  // await deploy('SimpleCoin', {
  //   from: deployer.address,
  //   args: [],
  //   // since it's difficult to estimate the gas limit before f4 address is launched, it's safer to manually set
  //   // a large gasLimit. This should be addressed in the following releases.
  //   gasLimit: 10000000000, // BlockGasLimit / 10
  //   // since Ethereum's legacy transaction format is not supported on FVM, we need to specify
  //   // maxPriorityFeePerGas to instruct hardhat to use EIP-1559 tx format
  //   maxPriorityFeePerGas: priorityFee,
  //   nonce: nonce,
  //   log: true,
  // });

  // await deploy('ERC20', {
  //   from: deployer.address,
  //   args: ['gold', 'GLD'],
  //   // since it's difficult to estimate the gas limit before f4 address is launched, it's safer to manually set
  //   // a large gasLimit. This should be addressed in the following releases.
  //   gasLimit: 10000000000, // BlockGasLimit / 10
  //   // since Ethereum's legacy transaction format is not supported on FVM, we need to specify
  //   // maxPriorityFeePerGas to instruct hardhat to use EIP-1559 tx format
  //   maxPriorityFeePerGas: priorityFee,
  //   nonce: nonce,
  //   log: true,
  // });
};
module.exports.tags = ['NitroAdjudicator'];

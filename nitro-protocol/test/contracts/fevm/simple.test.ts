import {Contract, utils, providers, Wallet} from 'ethers';
// eslint-disable-next-line import/order
import RpcEngine from '@glif/filecoin-rpc-client';
import {FeeMarketEIP1559Transaction} from '@ethereumjs/tx';
import {newSecp256k1Address} from '@glif/filecoin-address';

import SimpleCoinArtifact from '../../../artifacts/contracts/SimpleCoin.sol/SimpleCoin.json';
import {SimpleCoin} from '../../../typechain-types';

const wallabyUrl = 'https://wallaby.node.glif.io/rpc/v0';
const pk = '9182b5bf5b9c966e001934ebaf008f65516290cef6e3069d11e718cbd4336aae';
const addressWithFunds = '0xff00000000000000000000000000000000000415';
const contractAddress = '0xFf000000000000000000000000000000000004C8';

const provider = new providers.JsonRpcProvider(wallabyUrl);
export const simpleCoinAbi = new utils.Interface(SimpleCoinArtifact.abi);

const ethRpc = new RpcEngine({
  apiAddress: wallabyUrl,
  namespace: 'eth',
  delimeter: '_',
});
const filRpc = new RpcEngine({apiAddress: wallabyUrl});

const simpleCoinContract = new Contract(
  contractAddress,
  simpleCoinAbi,
  provider
) as unknown as SimpleCoin & Contract;

it.skip('submits a transaction', async () => {
  const txPromise = simpleCoinContract.getBalance(addressWithFunds);
  console.log(await txPromise);
});

it('submits a transaction', async () => {
  const data = simpleCoinContract.interface.encodeFunctionData('sendCoin', [
    '0xff00000000000000000000000000000000000485',
    1,
  ]);

  const {secpActor} = await deriveAddrsFromPk(pk, wallabyUrl);
  const priorityFee = await ethRpc.request('maxPriorityFeePerGas');
  const nonce = await filRpc.request('MpoolGetNonce', secpActor);

  const txObject = {
    nonce,
    gasLimit: 1000000000, // BlockGasLimit / 10
    to: contractAddress,
    maxPriorityFeePerGas: priorityFee,
    maxFeePerGas: '0x2E90EDD000',
    chainId: 31415,
    data,
    type: 2,
  };
  const tx = FeeMarketEIP1559Transaction.fromTxData(txObject);
  const sig = tx.sign(Buffer.from(pk, 'hex'));
  const serializedTx = sig.serialize();
  const rawTxHex = '0x' + serializedTx.toString('hex');
  const res = await ethRpc.request('sendRawTransaction', rawTxHex);
  console.log(res);
  // const txRequest: providers.TransactionRequest = {
  //   data,
  //   type: 2,
  //   chainId: 31415,
  //   to: '0xFf000000000000000000000000000000000004C8',
  //   gasLimit: '0x76c0000',
  //   maxFeePerGas: '0x9184e72a000',
  //   maxPriorityFeePerGas: '0x9184e72a000',
  //   nonce: 54,
  // };
  // const signedTx = await wallet.signTransaction(txRequest);
  // const txPromise = provider.sendTransaction(signedTx);
  // console.log(await txPromise);
});

function hexlify(id: string) {
  const hexId = Number(id.slice(1)).toString(16);
  return '0xff' + '0'.repeat(38 - hexId.length) + hexId;
}

async function deriveAddrsFromPk(pk: string, apiAddress: string) {
  const w = new Wallet(pk);
  const pubKey = Uint8Array.from(Buffer.from(w.publicKey.slice(2), 'hex'));
  const secpActor = newSecp256k1Address(pubKey).toString();
  const filRpc = new RpcEngine({apiAddress});

  const idActor = await filRpc.request('StateLookupID', secpActor, null);
  const idActorHex = hexlify(idActor);

  return {secpActor, idActor, idActorHex};
}

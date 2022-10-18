import {Contract, utils, providers, Wallet} from 'ethers';
// eslint-disable-next-line import/order
import RpcEngine from '@glif/filecoin-rpc-client';
import {FeeMarketEIP1559Transaction} from '@ethereumjs/tx';
import {newSecp256k1Address} from '@glif/filecoin-address';

import SimpleCoinArtifact from '../../../artifacts/contracts/SimpleCoin.sol/SimpleCoin.json';
import {SimpleCoin} from '../../../typechain-types';

const wallabyUrl = 'https://wallaby.node.glif.io/rpc/v0';
const pk = '716b7161580785bc96a4344eb52d23131aea0caf42a52dcf9f8aee9eef9dc3cd';
const addressWithFunds = '0xff000000000000000000000000000000000003f7';
const contractAddress = '0xff000000000000000000000000000000000003f9';

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

it('submits a transaction', async () => {
  const txPromise = simpleCoinContract.getBalance(addressWithFunds);
  console.log((await txPromise).toString());
});

it.skip('submits a transaction', async () => {
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

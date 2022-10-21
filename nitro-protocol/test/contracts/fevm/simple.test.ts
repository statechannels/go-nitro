import {Contract, utils, providers, Wallet, constants} from 'ethers';
// eslint-disable-next-line import/order
import RpcEngine from '@glif/filecoin-rpc-client';
import {FeeMarketEIP1559Transaction} from '@ethereumjs/tx';
import {newSecp256k1Address} from '@glif/filecoin-address';

import SimpleCoinArtifact from '../../../artifacts/contracts/SimpleCoin.sol/SimpleCoin.json';
import NitroArtifact from '../../../artifacts/contracts/NitroAdjudicator.sol/NitroAdjudicator.json';
import {NitroAdjudicator, SimpleCoin} from '../../../typechain-types';

const wallabyUrl = 'https://wallaby.node.glif.io/rpc/v0';
const pk = '716b7161580785bc96a4344eb52d23131aea0caf42a52dcf9f8aee9eef9dc3cd';
const simpleCoinAddress = '0xff000000000000000000000000000000000003f9';
const nitroAddress = '0xFF000000000000000000000000000000000003fA';
const channelId = '0xd9b535b686bcae01a00da8767de21d8bfc9915d513833160e5f15044fb4a3644';

const provider = new providers.JsonRpcProvider(wallabyUrl);
const simpleCoinAbi = new utils.Interface(SimpleCoinArtifact.abi);
const nitroAbi = new utils.Interface(NitroArtifact.abi);

const ethRpc = new RpcEngine({
  apiAddress: wallabyUrl,
  namespace: 'eth',
  delimeter: '_',
});
const filRpc = new RpcEngine({apiAddress: wallabyUrl});

const simpleCoinContract = new Contract(
  simpleCoinAddress,
  simpleCoinAbi,
  provider
) as unknown as SimpleCoin & Contract;

const nitroContract = new Contract(
  nitroAddress,
  nitroAbi,
  provider
) as unknown as NitroAdjudicator & Contract;

it.skip('reads balance', async () => {
  const {idActorHex} = await deriveAddrsFromPk(pk, wallabyUrl);
  const txPromise = simpleCoinContract.getBalance(idActorHex);
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
    to: simpleCoinAddress,
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

it('nitro deposit transaction', async () => {
  const data = nitroContract.interface.encodeFunctionData('deposit', [
    constants.AddressZero,
    channelId,
    0,
    1,
  ]);

  const {secpActor} = await deriveAddrsFromPk(pk, wallabyUrl);
  const priorityFee = await ethRpc.request('maxPriorityFeePerGas');
  console.log(priorityFee);
  const nonce = await filRpc.request('MpoolGetNonce', secpActor);

  const txObject = {
    nonce,
    gasLimit: 1000000000, // BlockGasLimit / 10
    to: nitroAddress,
    maxPriorityFeePerGas: priorityFee,
    maxFeePerGas: '0x2E90EDD000',
    chainId: 31415,
    data,
    type: 2,
    value: constants.One.toHexString(),
  };
  const tx = FeeMarketEIP1559Transaction.fromTxData(txObject);
  const sig = tx.sign(Buffer.from(pk, 'hex'));
  const serializedTx = sig.serialize();
  const rawTxHex = '0x' + serializedTx.toString('hex');
  const res = await ethRpc.request('sendRawTransaction', rawTxHex);
  console.log(res);
});

it('reads balance', async () => {
  const txPromise = nitroContract.holdings(constants.AddressZero, channelId);
  console.log((await txPromise).toString());
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

// // Copied from https://stackoverflow.com/questions/58325771/how-to-generate-random-hex-string-in-javascript
// function genRanHex(size: number) {
//   return [...Array(size)].map(() => Math.floor(Math.random() * 16).toString(16)).join('');
// }
// function randomChannelId() {
//   return '0x' + genRanHex(64);
// }

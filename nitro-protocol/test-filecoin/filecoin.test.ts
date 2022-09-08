// @ts-ignore
import {LotusRPC} from '@filecoin-shipyard/lotus-client-rpc';
// @ts-ignore
import {NodejsProvider} from '@filecoin-shipyard/lotus-client-provider-nodejs';
// @ts-ignore
import {mainnet} from '@filecoin-shipyard/lotus-client-schema';
import {utils} from 'ethers';
import cbor from 'cbor';
// @ts-ignore
import {FilecoinSigner, Message, FilecoinClient} from '@blitslabs/filecoin-js-signer';
import axios from 'axios';

const filecoin_signer = new FilecoinSigner();

import TESTNitroAdjudicatorArtifact from '../artifacts/contracts/test/TESTNitroAdjudicator.sol/TESTNitroAdjudicator.json';
import {SignedMessage} from '@blitslabs/filecoin-js-signer/dist/core/types/types';

// const localGanache = 'ws://localhost:7777/rpc/v0';
// const localProvider = new NodejsProvider(localGanache);
// const localClient = new LotusRPC(localProvider, {schema: mainnet.fullNode});

// // follow https://github.com/Zondax/filecoin-signing-tools/blob/master/examples/wasm_node/payment_channel.js

const wallaby = 'https://wallaby.node.glif.io/rpc/v0';
const jimpicknet = 'https://fvm-4.default.knative.hex.camp/rpc/v0';
const jimpicktoken =
  'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJyZWFkIiwid3JpdGUiLCJzaWduIiwiYWRtaW4iXX0.Fp-tjQ3fbUEKSXKDAybCR_vEE7MdPQ1Uqw2_2MVClfk';

const endpoint = wallaby;
// const provider = new NodejsProvider(wallaby);
// const client = new LotusRPC(provider, {schema: mainnet.fullNode});

// const filecoin_client = new FilecoinClient(wallaby, '');

// const myAddress =
// 't3w3ycieznjdzczqknlroo7aruchh5e373wz76su66ae4if3jrqjux346hyffu3l6htawazyiagbk6rne44cya';
// const privateKey = '380070f27bf48d2c9e50952553371171ea8b7ce0eb16a19b978821183baef811';

// const myAddress = 't1bl36yn6vr5mruyvaehidna2ijdencfgpkflflmy';
// const privateKey = '6645aa9129061ccef190e1bb1e11319b3d716b3140eec27595d045dbd565733b';
// const privateKeyLotus =
//   '7b2254797065223a22736563703235366b31222c22507269766174654b6579223a225a6b57716b536b47484d37786b4f4737486845786d7a3178617a464137734a316c6442463239566c637a733d227d';

const myAddress = 't1vdc3bjjejtdipqri77xtpc4cozmrglspe7qumfi';
const privateKey = Buffer.from('/+4JRNot2aej55QWksaAVUER02GkDvdS3eOlmrPt5XM=', 'base64').toString(
  'hex'
);

// const foo = filecoin_signer.wallet.keyRecover(privateKey);
// console.log(foo);

describe('Connects to a local filecoin-flavoured ganache instance', () => {
  // the instance must be started manually before the test is run

  const simplecoinhex =
    '0x608060405234801561001057600080fd5b506127106000803273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550610556806100656000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c80637bd703e81461004657806390b98a1114610076578063f8b2cb4f146100a6575b600080fd5b610060600480360381019061005b919061030a565b6100d6565b60405161006d9190610350565b60405180910390f35b610090600480360381019061008b9190610397565b6100f4565b60405161009d91906103f2565b60405180910390f35b6100c060048036038101906100bb919061030a565b61025f565b6040516100cd9190610350565b60405180910390f35b600060026100e38361025f565b6100ed919061043c565b9050919050565b6000816000803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205410156101455760009050610259565b816000803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546101939190610496565b92505081905550816000808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546101e891906104ca565b925050819055508273ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8460405161024c9190610350565b60405180910390a3600190505b92915050565b60008060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050919050565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006102d7826102ac565b9050919050565b6102e7816102cc565b81146102f257600080fd5b50565b600081359050610304816102de565b92915050565b6000602082840312156103205761031f6102a7565b5b600061032e848285016102f5565b91505092915050565b6000819050919050565b61034a81610337565b82525050565b60006020820190506103656000830184610341565b92915050565b61037481610337565b811461037f57600080fd5b50565b6000813590506103918161036b565b92915050565b600080604083850312156103ae576103ad6102a7565b5b60006103bc858286016102f5565b92505060206103cd85828601610382565b9150509250929050565b60008115159050919050565b6103ec816103d7565b82525050565b600060208201905061040760008301846103e3565b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600061044782610337565b915061045283610337565b9250817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff048311821515161561048b5761048a61040d565b5b828202905092915050565b60006104a182610337565b91506104ac83610337565b9250828210156104bf576104be61040d565b5b828203905092915050565b60006104d582610337565b91506104e083610337565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff038211156105155761051461040d565b5b82820190509291505056fea2646970667358221220d6b7446ff7b783a98dd1d68f516949438f2d578a05c05722e283b81a7774246564736f6c634300080e0033';
  const bytecodeHexString = TESTNitroAdjudicatorArtifact.bytecode;
  const evmBytes = utils.arrayify(bytecodeHexString);
  const evmBytesCbor = cbor.encode([evmBytes, new Uint8Array(0)]);

  const evmActorCidBytes = Buffer.from(
    '0155a0e40220aad04bf2cd6189c13a4594bde58718bb26d7f64ec8c2c4fee4085118625bc467',
    'hex'
  );

  const params = cbor.encode([new cbor.Tagged(42, evmActorCidBytes), evmBytesCbor]);

  // Sending create actor message...
  const messageBody: Message = {
    To: 't01',
    From: myAddress,
    Value: '0',
    Method: 2,
    Params: Buffer.from(simplecoinhex, 'hex').toString('base64'),
    // Params: bytecodeHexString,
    Version: 0,
    Nonce: 42,
    GasLimit: 1_000_000,
    GasFeeCap: '10000',
    GasPremium: '0',
  };

  it('reads the chain height', async () => {
    let response = await axios.post(wallaby, {
      jsonrpc: '2.0',
      method: 'Filecoin.Version',
      params: [],
      id: 1,
      auth: jimpicktoken,
    });
    console.log(response.data);

    console.log('signing message...');
    // const signature = filecoin_signer.utils.signMessage(messageBody, privateKey);
    const signature = filecoin_signer.tx.transactionSignLotus(messageBody, privateKey);
    // expect(filecoin_signer.utils.verifySignature(messageBody, signature, myAddress)).toBe(true);

    const signedMessage: SignedMessage = {
      Message: messageBody,
      Signature: {Type: 1, Data: Buffer.from(signature, 'hex').toString('base64')},
    };

    console.log('pushing signed message to mempool');
    const config = {
      headers: {Authorization: `Bearer ${jimpicktoken}`},
    };
    const responsePromise = axios.post(
      endpoint,
      {
        jsonrpc: '2.0',
        method: 'Filecoin.MpoolPush',
        params: [signedMessage],
        id: 2,
      },
      config
    );
    await Promise.race([
      responsePromise.catch(err => {
        console.log(err.response.data);
        throw err;
      }),
      responsePromise.then(response => console.log(response.data)),
    ]);
    // const response = await filecoin_client.tx.sendMessage(messageBody, privateKey);
    // // const response = await client.mpoolPush(signedMessage);
    // //   const waitResponse = await client.stateWaitMsg(response.CID, 0)

    // expect(response).toEqual('foo');

    // localClient.destroy();
    // client.destroy();
  });
});

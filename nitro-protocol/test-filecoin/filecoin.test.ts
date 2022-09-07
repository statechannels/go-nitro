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

const filecoin_signer = new FilecoinSigner();

import TESTNitroAdjudicatorArtifact from '../artifacts/contracts/test/TESTNitroAdjudicator.sol/TESTNitroAdjudicator.json';

const localGanache = 'ws://localhost:7777/rpc/v0';
const localProvider = new NodejsProvider(localGanache);
const localClient = new LotusRPC(localProvider, {schema: mainnet.fullNode});

// follow https://github.com/Zondax/filecoin-signing-tools/blob/master/examples/wasm_node/payment_channel.js

const wallaby = 'https://wallaby.node.glif.io/rpc/v1';
const provider = new NodejsProvider(wallaby);
const client = new LotusRPC(provider, {schema: mainnet.fullNode});

const filecoin_client = new FilecoinClient(wallaby, '');

const myAddress =
  't3w3ycieznjdzczqknlroo7aruchh5e373wz76su66ae4if3jrqjux346hyffu3l6htawazyiagbk6rne44cya';
const privateKey = '380070f27bf48d2c9e50952553371171ea8b7ce0eb16a19b978821183baef811';

describe('Connects to a local filecoin-flavoured ganache instance', () => {
  // the instance must be started manually before the test is run

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
    Value: 0,
    Method: 2,
    Params: params.toString('base64'),
    Version: 1,
    Nonce: 0,
    GasLimit: 1_000_000,
    GasFeeCap: '0',
    GasPremium: '0',
  };

  it('reads the chain height', async () => {
    const head = await localClient.chainHead();
    console.log(head);
    expect(head.Height).toEqual(0);

    // console.log('signing message...');
    // const signedMessage = filecoin_signer.tx.transactionSignLotus(messageBody, privateKey);

    // console.log(signedMessage);
    console.log('pushing signed message to mempool');
    const response = await filecoin_client.tx.sendMessage(messageBody, privateKey);
    // const response = await client.mpoolPush(signedMessage);
    //   const waitResponse = await client.stateWaitMsg(response.CID, 0)

    expect(response).toEqual('foo');

    localClient.destroy();
    client.destroy();
  });
});

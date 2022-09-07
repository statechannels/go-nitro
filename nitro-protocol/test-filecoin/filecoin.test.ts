// @ts-ignore
import {LotusRPC} from '@filecoin-shipyard/lotus-client-rpc';
// @ts-ignore
import {NodejsProvider} from '@filecoin-shipyard/lotus-client-provider-nodejs';
// @ts-ignore
import {mainnet} from '@filecoin-shipyard/lotus-client-schema';
import ethers from 'ethers';

import TESTNitroAdjudicatorArtifact from '../artifacts/contracts/test/TESTNitroAdjudicator.sol/TESTNitroAdjudicator.json';

// const wsUrl = 'ws://127.0.0.1:7777/rpc/v0';
const wallaby = 'https://wallaby.node.glif.io/rpc/v1';
const provider = new NodejsProvider(wallaby);
const client = new LotusRPC(provider, {schema: mainnet.fullNode});

const myAddress = 
'f1bl36yn6vr5mruyvaehidna2ijdencfgpkflflmy'
const privateKey = '6645aa9129061ccef190e1bb1e11319b3d716b3140eec27595d045dbd565733b'

describe('Connects to a local filecoin-flavoured ganache instance', () => {
  // @typesthe instance must be started manually before the test is run

  const bytecodeHexString = TESTNitroAdjudicatorArtifact.bytecode;
  const evmBytes = ethers.utils.arrayify(bytecodeHexString);
  
  const evmBytesCbor = cbor.encode([evmBytes, new Uint8Array(0)])


  const evmActorCidBytes = Buffer.from('0155a0e40220aad04bf2cd6189c13a4594bde58718bb26d7f64ec8c2c4fee4085118625bc467', 'hex');

  const params = cbor.encode([new cbor.Tagged(42, evmActorCidBytes), evmBytesCbor])

// Sending create actor message...
const messageBody = {
    To: 't01',
    From: myAddress,
    Value: "0",
    Method: 2,
    Params: params.toString('base64')
    };

const signedMessage = client.walletSignMessage(myAddress,messageBody);

      const response = await client.mpoolPush(messageBody)
    //   const waitResponse = await client.stateWaitMsg(response.CID, 0)
    }

  it('reads the chain height', async () => {
    const head = await client.chainHead();
    console.log(head);
    expect(head.Height).toEqual(0);
    client.destroy();
  });
});

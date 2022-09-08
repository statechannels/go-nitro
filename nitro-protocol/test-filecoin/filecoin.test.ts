// @ts-ignore
import {Message, FilecoinClient} from '@blitslabs/filecoin-js-signer';

import TESTNitroAdjudicatorArtifact from '../artifacts/contracts/test/TESTNitroAdjudicator.sol/TESTNitroAdjudicator.json';

const wallaby = 'https://wallaby.node.glif.io/rpc/v0';
const endpoint = wallaby;

// Generated on https://wallet.glif.io/?network=wallaby
// And hydrated with test FIL at https://wallaby.network/#faucet
const myAddress = 't1vdc3bjjejtdipqri77xtpc4cozmrglspe7qumfi';
const privateKey = Buffer.from('/+4JRNot2aej55QWksaAVUER02GkDvdS3eOlmrPt5XM=', 'base64').toString(
  'hex'
);

describe('Connects to wallaby testnet', () => {
  const bytecodeHexString = TESTNitroAdjudicatorArtifact.bytecode;

  it('sends a message deploying EVM bytecode', async () => {
    const filecoin_client = new FilecoinClient(endpoint);

    // Sending create actor message...
    const messageBody: Message = {
      To: 't01',
      From: myAddress,
      Value: '0',
      Method: 2,
      Params: Buffer.from(bytecodeHexString, 'hex').toString('base64'),
      Version: 0,
      Nonce: 42, // This will be overwritten by tx.sendMessage below
      GasLimit: 1_000_000,
      GasFeeCap: '30000',
      GasPremium: '20000',
    };

    // Internally, this will get a new nonce, sign the message, do some encoding and call the MpoolPush RPC method.
    const responsePromise = filecoin_client.tx.sendMessage(messageBody, privateKey);

    let responseCID: string;

    await Promise.race([
      responsePromise.catch(err => {
        console.log(err);
        throw err;
      }),
      responsePromise.then(response => {
        responseCID = (
          response as {
            '/': string;
          }
        )['/'];
        console.log('response CID is', responseCID);
      }),
    ]);
  });
});

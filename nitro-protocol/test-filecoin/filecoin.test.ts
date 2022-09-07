// @ts-ignore
import {LotusRPC} from '@filecoin-shipyard/lotus-client-rpc';
// @ts-ignore
import {NodejsProvider} from '@filecoin-shipyard/lotus-client-provider-nodejs';
// @ts-ignore
import {mainnet} from '@filecoin-shipyard/lotus-client-schema';

const wsUrl = 'ws://127.0.0.1:7777/rpc/v0';
const provider = new NodejsProvider(wsUrl);
const client = new LotusRPC(provider, {schema: mainnet.fullNode});

describe('Connects to a local filecoin-flavoured ganache instance', () => {
  // @typesthe instance must be started manually before the test is run

  it('reads the chain height', async () => {
    const head = await client.chainHead();
    expect(head.Height).toEqual(0);
    client.destroy();
  });
});

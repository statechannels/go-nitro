# Nitro RPC Client

The nitro RPC client is a typescript client that can make RPC calls against a go-nitro RPC server.

## Using the RPC Client

```typescript
import { NitroRpcClient } from "./rpc-client";

const rpcPort = 4222;

const rpcClient = await NitroRpcClient.CreateNatsRpcClient(
  `127.0.0.1:${rpcPort}`
);

const counterParty = `0xDEADBEEF`;

const paymentChannelInfo = await rpcClient.DirectFund(counterParty);

console.log(
  `Created channel ${paymentChannelInfo.ChannelId} with counterparty ${counterParty}`
);

await rpcClient.Close();
```

## CLI Tool

The Nitro RPC comes with a CLI tool to trigger calls through the nitro RPC client.

```shell
npm exec -c 'nitro-rpc-client version'

```

The full list of API calls can be seen using the `--help` flag.

```shell
❯ npm exec -c 'nitro-rpc-client --help'
nitro-rpc-client <command>

Commands:
  nitro-rpc-client version                  Get the version of the Nitro RPC
                                            server
  nitro-rpc-client address                  Get the address of the Nitro RPC
                                            server
  nitro-rpc-client direct-fund              Creates a directly funded ledger
  <counterparty>                            channel
  nitro-rpc-client direct-defund            Defunds a directly funded ledger
  <channelId>                               channel
  nitro-rpc-client virtual-fund             Creates a virtually funded payment
  <counterparty> [intermediaries...]        channel
  nitro-rpc-client virtual-defund           Defunds a virtually funded payment
  <channelId>                               channel
  nitro-rpc-client get-ledger-channel       Gets information about a ledger
  <channelId>                               channel
  nitro-rpc-client get-payment-channel      Gets information about a payment
  <channelId>                               channel
  nitro-rpc-client pay <channelId>          Sends a payment on the given channel
  <amount>

Options:
      --help     Show help                                             [boolean]
      --version  Show version number                                   [boolean]
  -p, --port                                            [number] [default: 4005]
```

### Using the create-channels script

A test script is available to easily create channels. It requires 3 running RPC servers for Alice,Bob, and Irene. The `go-nitro` repository contains a [test script to start and run the required RPC servers](https://github.com/statechannels/go-nitro#start-rpc-servers-test-script)

The `create-channels` script will do the following (using the `nitro-rpc-client`):

1. Create a ledger channel between Alice and Irene.
2. Create a ledger channel between Bon and Irene.
3. Create some virtual channels.
4. Make some payments.
5. Close some virtual channels.

The script can be run from the `packages/nitro-rpc-client` folder with `npx ts-node ./scripts/client-runner.ts create-channels`.

The script also accepts a few options for the amount of channels to create/close.

```
❯ npx ts-node ./scripts/client-runner.ts create-channels --help
client-runner create-channels

Creates some virtual channels and makes some payments

Options:
  --help             Show help                                         [boolean]
  --version          Show version number                               [boolean]
  --createledgers    Whether we attempt to create new ledger channels.
                     Set to false if you already have some ledger channels
                     created.                          [boolean] [default: true]
  --numvirtual       The number of virtual channels to create between Alice and
                     Bob.                                  [number] [default: 5]
  --numclosevirtual  The number of virtual channels to close and defund.
                                                           [number] [default: 2]
  --numpayments      The number of payments to make from Alice to Bob.Each
                     payment is made on a random virtual channel
                                                           [number] [default: 5]
```

The `--createledgers` option is helpful when you want to create some additional virtual channels using existing ledger channels. Setting it to `false` means the script will just use the existing ledger channels when creating new virtual channels.

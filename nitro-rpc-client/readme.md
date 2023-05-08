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

#!/usr/bin/env ts-node
/* eslint-disable @typescript-eslint/no-empty-function */
/* eslint-disable @typescript-eslint/no-shadow */

import yargs from 'yargs/yargs';

import { NitroRpcClient } from './rpc-client';

yargs(process.argv.slice(2))
  .scriptName('nitro-rpc-client')
  .option({
    p: { alias: 'port', default: 4005, type: 'number' },
  })
  .command(
    'version',
    'Get the version of the Nitro RPC server',
    async () => {},
    async (yargs) => {
      const rpcPort = yargs.p;

      const rpcClient = await NitroRpcClient.CreateHttpNitroClient(
        getLocalRPCUrl(rpcPort),
      );
      const version = await rpcClient.GetVersion();
      console.log(version);
      await rpcClient.Close();
      process.exit(0);
    },
  )
  .command(
    'address',
    'Get the address of the Nitro RPC server',
    async () => {},
    async (yargs) => {
      const rpcPort = yargs.p;

      const rpcClient = await NitroRpcClient.CreateHttpNitroClient(
        getLocalRPCUrl(rpcPort),
      );
      const address = await rpcClient.GetAddress();
      console.log(address);
      await rpcClient.Close();
      process.exit(0);
    },
  )
  .command(
    'direct-fund <counterparty>',
    'Creates a directly funded ledger channel',
    (yargsBuilder) => {
      return yargsBuilder.positional('counterparty', {
        describe: "The counterparty's address",
        type: 'string',
        demandOption: true,
      });
    },
    async (yargs) => {
      const rpcPort = yargs.p;

      const rpcClient = await NitroRpcClient.CreateHttpNitroClient(
        getLocalRPCUrl(rpcPort),
      );

      const dfObjective = await rpcClient.DirectFund(yargs.counterparty);
      const { Id } = dfObjective;
      console.log(`Objective started ${Id}`);
      await rpcClient.WaitForObjective(Id);
      console.log(`Objective complete ${Id}`);
      await rpcClient.Close();
      process.exit(0);
    },
  )
  .command(
    'direct-defund <channelId>',
    'Defunds a directly funded ledger channel',
    (yargsBuilder) => {
      return yargsBuilder.positional('channelId', {
        describe: 'The id of the ledger channel to defund',
        type: 'string',
        demandOption: true,
      });
    },
    async (yargs) => {
      const rpcPort = yargs.p;

      const rpcClient = await NitroRpcClient.CreateHttpNitroClient(
        getLocalRPCUrl(rpcPort),
      );

      const id = await rpcClient.DirectDefund(yargs.channelId);
      console.log(`Objective started ${id}`);
      await rpcClient.WaitForObjective(id);
      console.log(`Objective complete ${id}`);
      await rpcClient.Close();
      process.exit(0);
    },
  )
  .command(
    'virtual-fund <counterparty> [intermediaries...]',
    'Creates a virtually funded payment channel',
    (yargsBuilder) => {
      return yargsBuilder
        .positional('counterparty', {
          describe: "The counterparty's address",
          type: 'string',
          demandOption: true,
        })
        .array('intermediaries');
    },
    async (yargs) => {
      const rpcPort = yargs.p;

      const rpcClient = await NitroRpcClient.CreateHttpNitroClient(
        getLocalRPCUrl(rpcPort),
      );

      // Parse all intermediary args to strings
      const intermediaries =
        yargs.intermediaries?.map((intermediary) => {
          if (typeof intermediary === 'string') {
            return intermediary;
          }
          return intermediary.toString(16);
        }) ?? [];

      const vfObjective = await rpcClient.VirtualFund(
        yargs.counterparty,
        intermediaries,
      );

      const { Id } = vfObjective.result;
      console.log(`Objective started ${Id}`);
      await rpcClient.WaitForObjective(Id);
      console.log(`Objective complete ${Id}`);
      await rpcClient.Close();
      process.exit(0);
    },
  )
  .command(
    'virtual-defund <channelId>',
    'Defunds a virtually funded payment channel',
    (yargsBuilder) => {
      return yargsBuilder.positional('channelId', {
        describe: 'The id of the payment channel to defund',
        type: 'string',
        demandOption: true,
      });
    },
    async (yargs) => {
      const rpcPort = yargs.p;

      const rpcClient = await NitroRpcClient.CreateHttpNitroClient(
        getLocalRPCUrl(rpcPort),
      );

      const id = await rpcClient.VirtualDefund(yargs.channelId);

      console.log(`Objective started ${id}`);
      await rpcClient.WaitForObjective(id);
      console.log(`Objective complete ${id}`);
      await rpcClient.Close();
      process.exit(0);
    },
  )
  .command(
    'get-ledger-channel <channelId>',
    'Gets information about a ledger channel',
    (yargsBuilder) => {
      return yargsBuilder.positional('channelId', {
        describe: 'The channel ID of the ledger channel',
        type: 'string',
        demandOption: true,
      });
    },
    async (yargs) => {
      const rpcPort = yargs.p;

      const rpcClient = await NitroRpcClient.CreateHttpNitroClient(
        getLocalRPCUrl(rpcPort),
      );

      const ledgerInfo = await rpcClient.GetLedgerChannel(yargs.channelId);
      console.log(ledgerInfo);
      await rpcClient.Close();
      process.exit(0);
    },
  )
  .command(
    'get-payment-channel <channelId>',
    'Gets information about a payment channel',
    (yargsBuilder) => {
      return yargsBuilder.positional('channelId', {
        describe: 'The channel ID of the payment channel',
        type: 'string',
        demandOption: true,
      });
    },
    async (yargs) => {
      const rpcPort = yargs.p;

      const rpcClient = await NitroRpcClient.CreateHttpNitroClient(
        getLocalRPCUrl(rpcPort),
      );
      const paymentChannelInfo = await rpcClient.GetPaymentChannel(
        yargs.channelId,
      );
      console.log(paymentChannelInfo);
      await rpcClient.Close();
      process.exit(0);
    },
  )
  .command(
    'pay <channelId> <amount>',
    'Sends a payment on the given channel',
    (yargsBuilder) => {
      return yargsBuilder
        .positional('channelId', {
          describe: 'The channel ID of the payment channel',
          type: 'string',
          demandOption: true,
        })
        .positional('amount', {
          describe: 'The amount to pay',
          type: 'number',
          demandOption: true,
        });
    },
    async (yargs) => {
      const rpcPort = yargs.p;

      const rpcClient = await NitroRpcClient.CreateHttpNitroClient(
        getLocalRPCUrl(rpcPort),
      );
      const paymentChannelInfo = await rpcClient.Pay(
        yargs.channelId,
        yargs.amount,
      );
      console.log(paymentChannelInfo);
      await rpcClient.Close();
      process.exit(0);
    },
  )

  .demandCommand(1, 'You need at least one command before moving on')
  .parserConfiguration({ 'parse-numbers': false })
  .parse();

function getLocalRPCUrl(port: number): string {
  return `127.0.0.1:${port}`;
}

#!/usr/bin/env ts-node
/* eslint-disable @typescript-eslint/no-empty-function */
/* eslint-disable @typescript-eslint/no-shadow */

import yargs from "yargs";

import { NitroRpcClient } from "../src/rpc-client";
yargs(process.argv.slice(2))
  .scriptName("client-runner")
  .command(
    "create-channels",
    "Creates some virtual channels and makes some payments",
    (yargsBuilder) => {
      return yargsBuilder
        .option("createledgers", {
          describe: `Whether we attempt to create new ledger channels.
            Set to false if you already have some ledger channels created.`,
          type: "boolean",
          default: true,
        })
        .option("numvirtual", {
          describe:
            "The number of virtual channels to create between Alice and Bob.",
          type: "number",
          default: 5,
        })
        .option("numpayments", {
          describe:
            "The number of payments to make from Alice to Bob.Each payment is made on a random virtual channel",
          type: "number",
          default: 5,
        });
    },
    async (yargs) => {
      const aliceClient = await NitroRpcClient.CreateHttpNitroClient(
        getLocalRPCUrl(4005)
      );

      const ireneClient = await NitroRpcClient.CreateHttpNitroClient(
        getLocalRPCUrl(4006)
      );

      const ireneAddress = await ireneClient.GetAddress();
      const bobClient = await NitroRpcClient.CreateHttpNitroClient(
        getLocalRPCUrl(4007)
      );
      const bobAddress = await bobClient.GetAddress();

      if (yargs.createledgers) {
        // Setup ledger channels
        console.log("Constructing ledger channels");
        const aliceLedger = await aliceClient.DirectFund(ireneAddress);
        console.log(`Ledger channel ${aliceLedger.ChannelId} created`);
        const bobLedger = await ireneClient.DirectFund(bobAddress);
        console.log(`Ledger channel ${bobLedger.ChannelId} created`);
        await wait(1000);
      }

      const virtualChannels: string[] = [];
      // Setup virtual channels
      for (let i = 0; i < yargs.numvirtual; i++) {
        const res = await aliceClient.VirtualFund(bobAddress, [ireneAddress]);
        console.log(`Virtual channel ${res.ChannelId} created`);
        virtualChannels.push(res.ChannelId);
      }
      await wait(1000);

      // Make payments
      for (let i = 0; i < yargs.numpayments; i++) {
        const channelId = getRandomElement(virtualChannels);
        const res = await aliceClient.Pay(channelId, 1);
        console.log(`Paid ${res.Amount} on channel ${res.Channel}`);
      }

      aliceClient.Close();
      ireneClient.Close();
      bobClient.Close();
      process.exit(0);
    }
  )
  .demandCommand(1, "You need at least one command before moving on")
  .parse();

function getLocalRPCUrl(port: number): string {
  return `127.0.0.1:${port}`;
}
function getRandomElement(col: any[]) {
  return col[Math.floor(Math.random() * col.length)];
}
async function wait(ms: number) {
  await new Promise((res) => setTimeout(res, ms));
}

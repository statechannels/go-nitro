#!/usr/bin/env ts-node
/* eslint-disable @typescript-eslint/no-empty-function */
/* eslint-disable @typescript-eslint/no-shadow */

import yargs from "yargs";
import { hideBin } from "yargs/helpers";
import axios, { AxiosResponse } from "axios";

import { NitroRpcClient } from "../src/rpc-client";
import {
  compactJson,
  generateRequest,
  getLocalRPCUrl,
  logOutChannelUpdates,
} from "../src/utils";

yargs(hideBin(process.argv))
  .scriptName("client-runner")
  .command(
    "print-channels",
    "Prints all channels",
    async () => {},
    async () => {
      const aliceClient = await NitroRpcClient.CreateHttpNitroClient(
        getLocalRPCUrl(4005)
      );

      const ireneClient = await NitroRpcClient.CreateHttpNitroClient(
        getLocalRPCUrl(4006)
      );

      const bobClient = await NitroRpcClient.CreateHttpNitroClient(
        getLocalRPCUrl(4007)
      );

      for (const client of [aliceClient, ireneClient, bobClient]) {
        const ledgers = await client.GetAllLedgerChannels();

        console.log(
          `Client ${await client.GetAddress()} found ${ledgers.length} ledgers`
        );

        if (ledgers.length > 0) {
          console.log(`LEDGERS`);
        }
        for (const ledger of ledgers) {
          console.log(`${compactJson(ledger)}`);

          const paymentChans = await client.GetPaymentChannelsByLedger(
            ledger.ID
          );
          if (paymentChans.length > 0) {
            console.log(`PAYMENT CHANNELS FUNDED BY ${ledger.ID}`);

            for (const paymentChan of paymentChans) {
              console.log(`${compactJson(paymentChan)}`);
            }
          }
        }
      }
      await aliceClient.Close();
      await ireneClient.Close();
      await bobClient.Close();
      process.exit(0);
    }
  )
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
        .option("numclosevirtual", {
          describe: "The number of virtual channels to close and defund.",
          type: "number",
          default: 2,
        })
        .option("numpayments", {
          describe:
            "The number of payments to make from Alice to Bob.Each payment is made on a random virtual channel",
          type: "number",
          default: 5,
        })
        .option("printnotifications", {
          alias: "n",
          describe: "Whether channel notifications are printed to the console",
          type: "boolean",
          default: false,
        })
        .option("waitforservers", {
          alias: "w",
          describe:
            "Whether to wait for the RPC servers to be ready. If set to true we wait for 1 minute for the servers to be ready",
          type: "boolean",
          default: false,
        });
    },
    async (yargs) => {
      if (yargs.waitforservers) {
        const TWO_MIN = 120_000;
        await Promise.all([
          waitForRPCServer(4005, TWO_MIN),
          waitForRPCServer(4006, TWO_MIN),
          waitForRPCServer(4007, TWO_MIN),
        ]);
      }
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

      if (yargs.printnotifications) {
        logOutChannelUpdates(aliceClient);
        logOutChannelUpdates(ireneClient);
        logOutChannelUpdates(bobClient);
      }

      if (yargs.createledgers) {
        // Setup ledger channels
        console.log("Constructing ledger channels");
        const aliceLedger = await aliceClient.DirectFund(ireneAddress);
        console.log(`Ledger channel ${aliceLedger.ChannelId} created`);

        const bobLedger = await ireneClient.DirectFund(bobAddress);
        console.log(`Ledger channel ${bobLedger.ChannelId} created`);

        await wait(1000);
      }

      // Setup virtual channels
      const virtualChannels: string[] = [];
      console.log(`Constructing ${yargs.numvirtual} virtual channels`);
      for (let i = 0; i < yargs.numvirtual; i++) {
        const res = await aliceClient.VirtualFund(bobAddress, [ireneAddress]);
        console.log(`Virtual channel ${res.ChannelId} created`);
        virtualChannels.push(res.ChannelId);
      }
      await wait(1000);

      // Make payments
      console.log(`Making ${yargs.numpayments} payments`);
      for (let i = 0; i < yargs.numpayments; i++) {
        const channelId = getRandomElement(virtualChannels);
        const res = await aliceClient.Pay(channelId, 1);
        console.log(`Paid ${res.Amount} on channel ${res.Channel}`);
      }

      await wait(1000);

      // Close virtual channels
      console.log(`Closing ${yargs.numclosevirtual} virtual channels`);
      let closeCount = 0;
      for (const channelId of virtualChannels) {
        if (closeCount >= yargs.numclosevirtual) {
          break;
        }

        const res = await aliceClient.VirtualDefund(channelId);
        console.log(
          `Virtual channel ${getChannelIdFromObjectiveId(res)} closed`
        );
        closeCount++;
      }

      aliceClient.Close();
      ireneClient.Close();
      bobClient.Close();
      process.exit(0);
    }
  )
  .demandCommand(1, "You need at least one command before moving on")
  .strict()
  .parse();

function getRandomElement<T>(col: T[]): T {
  return col[Math.floor(Math.random() * col.length)];
}
async function wait(ms: number) {
  await new Promise((res) => setTimeout(res, ms));
}

function getChannelIdFromObjectiveId(objectiveId: string): string {
  return objectiveId.split("-")[1];
}

// Waits for the RPC server to be available by sending a simple get_address POST request until we get a response
async function waitForRPCServer(
  port: number,
  waitDuration: number
): Promise<void> {
  const startTime = Date.now();
  while (Date.now() - startTime < waitDuration) {
    const isUp = await isServerUp(port);
    if (isUp) {
      console.log(`RPC server ${getLocalRPCUrl(port)} is responding!`);
      return;
    } else {
      console.log(`RPC server ${getLocalRPCUrl(port)} not available, waiting!`);
      await wait(1000);
    }
  }
  throw new Error(
    `RPC server ${getLocalRPCUrl(port)} not reachable in ${waitDuration} ms`
  );
}

// Checks if the server is up by sending a simple get_address POST request
// This is specific to the HTTP/WS RPC transport
async function isServerUp(port: number): Promise<boolean> {
  let result: AxiosResponse<unknown, unknown>;

  try {
    const req = generateRequest("get_address", {});
    result = await axios.post(
      `http://${getLocalRPCUrl(port)}`,
      JSON.stringify(req)
    );
  } catch (e) {
    return false;
  }
  return result.status === 200;
}

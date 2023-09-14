#!/usr/bin/env ts-node
/* eslint-disable @typescript-eslint/no-empty-function */

import yargs from "yargs";
import { hideBin } from "yargs/helpers";
import axios, { AxiosResponse } from "axios";

import { NitroRpcClient } from "../src/rpc-client";
import {
  compactJson,
  getLocalRPCUrl,
  logOutChannelUpdates,
} from "../src/utils";

const clientNames = ["alice", "irene", "bob", "ivan"] as const;
const clientPortMap: Record<ClientNames, number> = {
  alice: 4005,
  irene: 4006,
  bob: 4007,
  ivan: 4008,
};

type ClientNames = (typeof clientNames)[number];
type Clients = Map<ClientNames, NitroRpcClient>;

async function initializeClients(): Promise<Clients> {
  let port = 4005;

  const clients: Clients = new Map<ClientNames, NitroRpcClient>();
  for (const clientName of clientNames) {
    const client = await NitroRpcClient.CreateHttpNitroClient(
      getLocalRPCUrl(port)
    );
    clients.set(clientName, client);
    port++;
  }
  return clients;
}

function closeClients(clients: Clients): Promise<void[]> {
  return Promise.all(
    Array.from(clients.values()).map((client) => client.Close())
  );
}

function isValidClientName(name: string): name is ClientNames {
  return clientNames.includes(name as ClientNames);
}

yargs(hideBin(process.argv))
  .scriptName("client-runner")
  .command(
    "print-channels",
    "Prints all channels",
    async () => {},
    async () => {
      const clients = await initializeClients();

      for (const client of clients.values()) {
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
      await closeClients(clients);
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
        .option("ledgerdeposit", {
          describe:
            "The amount of wei for each participant to deposit into the ledger channel",
          type: "number",
          default: 1_000_000_000_000,
        })
        .option("virtualdeposit", {
          describe:
            "The amount of wei for each participant to deposit into a virtual channel",
          type: "number",
          default: 1_000_000,
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
        .option("waitduration", {
          alias: "w",
          describe:
            "The amount of time, in milliseconds, to wait for the RPC servers to be responsive",
          type: "number",
          default: 0,
        });
    },
    async (yargs) => {
      if (yargs.waitduration > 0) {
        await Promise.all(
          Object.values(clientPortMap).map((port) =>
            waitForRPCServer(port, yargs.waitduration)
          )
        );
      }

      const clients = await initializeClients();

      if (yargs.printnotifications) {
        clients.forEach(logOutChannelUpdates);
      }

      const aliceClient = clients.get("alice");
      const ireneClient = clients.get("irene");
      const bobClient = clients.get("bob");
      if (!aliceClient || !ireneClient || !bobClient) {
        throw new Error("An client is undefined");
      }

      console.log("Retrieving client addresses");
      const ireneAddress = await ireneClient.GetAddress();
      const bobAddress = await bobClient.GetAddress();

      if (yargs.createledgers) {
        // Setup ledger channels
        console.log("Constructing ledger channels");
        const aliceLedger = await aliceClient.CreateLedgerChannel(
          ireneAddress,
          yargs.ledgerdeposit
        );
        const bobLedger = await ireneClient.CreateLedgerChannel(
          bobAddress,
          yargs.ledgerdeposit
        );

        await Promise.all([
          aliceClient.WaitForObjective(aliceLedger.Id),
          bobClient.WaitForObjective(bobLedger.Id),
        ]);
        console.log(`Ledger channel ${bobLedger.ChannelId} created`);
        console.log(`Ledger channel ${aliceLedger.ChannelId} created`);
      }

      // Setup virtual channels
      const virtualChannels: string[] = [];
      console.log(`Constructing ${yargs.numvirtual} virtual channels`);
      for (let i = 0; i < yargs.numvirtual; i++) {
        const res = await aliceClient.CreatePaymentChannel(
          bobAddress,
          [ireneAddress],
          yargs.virtualdeposit
        );
        await aliceClient.WaitForObjective(res.Id);
        console.log(`Virtual channel ${res.ChannelId} created`);
        virtualChannels.push(res.ChannelId);
      }

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

        const res = await aliceClient.ClosePaymentChannel(channelId);
        await aliceClient.WaitForObjective(res);
        console.log(
          `Virtual channel ${getChannelIdFromObjectiveId(res)} closed`
        );
        closeCount++;
      }

      await closeClients(clients);
      process.exit(0);
    }
  )
  .command(
    "create-ledger",
    "Create a ledger channel",
    (yargsBuilder) => {
      return yargsBuilder
        .option("left", {
          describe: "Left rpc node",
          type: "string",
          default: "alice",
        })
        .option("right", {
          describe: "Right rpc node",
          type: "string",
          default: "irene",
        });
    },
    async (yargs) => {
      const { left, right } = yargs;
      if (!isValidClientName(left) || !isValidClientName(right)) {
        throw new Error("Invalid client name");
      }

      const clients = await initializeClients();
      const leftClient = clients.get(left);
      const rightClient = clients.get(right);
      if (!leftClient || !rightClient) {
        throw new Error("A client is undefined");
      }
      const rightAddress = await rightClient.GetAddress();

      if (yargs.printnotifications) {
        [leftClient, rightClient].forEach(logOutChannelUpdates);
      }

      // Setup ledger channels
      console.log("Constructing ledger channels");
      const ledger = await leftClient.CreateLedgerChannel(
        rightAddress,
        1_000_000
      );
      await leftClient.WaitForObjective(ledger.Id);
      console.log(`Ledger channel ${ledger.ChannelId} created`);

      await closeClients(clients);
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
      console.log(
        `RPC server ${getLocalRPCUrl(port)} is responding!
        Waited ${
          Date.now() - startTime
        } milliseconds for server to be responsive`
      );

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
  const url = new URL(`https://${getLocalRPCUrl(port)}/health`).toString();

  try {
    result = await axios.get(url);
  } catch (e) {
    return false;
  }
  return result.status === 200;
}

#!/usr/bin/env ts-node
import { ChildProcess, exec, spawn } from "child_process";

enum Participant {
  Alice = "alice",
  Irene = "irene",
  Bob = "bob",
}
enum Color {
  black = "[30m",
  red = "[31m",
  green = "[32m",
  yellow = "[33m",
  blue = "[34m",
  magenta = "[35m",
  cyan = "[36m",
  white = "[37m",
  gray = "[90m",
}

// This is a simple script that:
// - starts up a local anvil chain
// - starts up 3 local rpc servers for alice,bob, and irene
// It can be run from the go-nitro folder with the following command:
// npx ts-node scripts/start-rpc-servers.ts
const chain = exec("anvil --chain-id 1337");
chain.stdout?.on("data", (data) => {
  console.log(data.toString());
});
chain.stderr?.on("data", (data) => {
  console.log(data.toString());
});

setupRPCServer(Participant.Alice, Color.blue);
setupRPCServer(Participant.Irene, Color.green);
setupRPCServer(Participant.Bob, Color.yellow);

function printWithColor(message: string, color: Color) {
  console.log(`\x1b${color.toString()}${message}\x1b[0m`);
}

function setupRPCServer(participant: Participant, color: Color): ChildProcess {
  let rpcServerProcess: ChildProcess;
  switch (participant) {
    case Participant.Alice:
      rpcServerProcess = exec("go run . -autodeploy -usedurablestore");
      break;
    case Participant.Irene:
      rpcServerProcess = exec(
        "go run . -autodeploy -usedurablestore -msgport 3006 -rpcport 4006 -pk febb3b74b0b52d0976f6571d555f4ac8b91c308dfa25c7b58d1e6a7c3f50c781"
      );
      break;
    case Participant.Bob:
      rpcServerProcess = exec(
        "go run . -autodeploy -usedurablestore -msgport 3007 -rpcport 4007 -pk 0279651921cd800ac560c21ceea27aab0107b67daf436cdd25ce84cad30159b4"
      );
      break;
  }

  rpcServerProcess.stdout?.on("data", (data) => {
    printWithColor(data.toString(), color);
  });
  rpcServerProcess.stderr?.on("data", (data) => {
    printWithColor(data.toString(), color);
  });
  return rpcServerProcess;
}

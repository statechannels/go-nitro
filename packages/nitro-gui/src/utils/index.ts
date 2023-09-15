import bigDecimal from "js-big-decimal-esm";

import { NitroRpcClient } from "@statechannels/nitro-rpc-client";
import { LedgerChannelInfo } from "@statechannels/nitro-rpc-client/src/types";

export const prettyPrintWei = (wei: bigint): string => {
  const PRECISION = 1;
  const names = ["wei", "kwei", "Mwei", "Gwei", "szabo", "finney", "ether"];
  const decimals = [0n, 3n, 6n, 9n, 12n, 15n, 18n];

  if (wei === 0n) {
    return "0 wei";
  } else if (!wei) {
    return "unknown";
  }
  let formattedString = "";
  decimals.forEach((decimal, index) => {
    if (wei >= 10n ** decimal) {
      formattedString = `${bigDecimal.divide(wei, 10n ** decimal, PRECISION)} ${
        names[index]
      }`;
    }
  });
  return formattedString;
};

export async function fetchAndSetLedgerChannels(
  nitroClient: NitroRpcClient,
  setLedgerChannels: (l: LedgerChannelInfo[]) => void
) {
  setLedgerChannels(await nitroClient.GetAllLedgerChannels());
}

export async function getRpcPort(): Promise<string> {
  const xhr = new XMLHttpRequest();
  xhr.open("GET", new URL("/rpc-port", window.location.href));
  xhr.send();
  return new Promise((res, rej) => {
    xhr.onload = function () {
      if (xhr.status === 200) {
        res(xhr.responseText);
      } else if (xhr.status === 404) {
        rej("could not get rpc port");
      }
    };
  });
}

export async function getRpcHost(): Promise<string> {
  if (import.meta.env.VITE_RPC_HOST) {
    return import.meta.env.VITE_RPC_HOST;
  } else {
    return getRpcPort().then(
      (rpcPort) => window.location.hostname + ":" + rpcPort
    );
  }
}

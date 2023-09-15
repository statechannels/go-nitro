import { useEffect, useState } from "react";
import { NitroRpcClient } from "@statechannels/nitro-rpc-client";
import { LedgerChannelInfo } from "@statechannels/nitro-rpc-client/src/types";

import "./App.css";
import TopBar from "./components/TopBar";
import LedgerChannelDetails from "./components/LedgerChannelDetails";
import PaymentChannelContainer from "./components/PaymentChannelContainer";

async function fetchAndSetLedgerChannels(
  nitroClient: NitroRpcClient,
  setLedgerChannels: (l: LedgerChannelInfo[]) => void
) {
  setLedgerChannels(await nitroClient.GetAllLedgerChannels());
}

async function getRpcPort(): Promise<string> {
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

async function getRpcHost(): Promise<string> {
  if (import.meta.env.VITE_RPC_HOST) {
    return import.meta.env.VITE_RPC_HOST;
  } else {
    return getRpcPort().then(
      (rpcPort) => window.location.hostname + ":" + rpcPort
    );
  }
}

function App() {
  const [nitroClient, setNitroClient] = useState<NitroRpcClient | null>(null);
  const [ledgerChannels, setLedgerChannels] = useState<LedgerChannelInfo[]>([]);
  const [focusedLedgerChannel, setFocusedLedgerChannel] = useState<string>("");

  useEffect(() => {
    getRpcHost().then((rpcHost) => {
      NitroRpcClient.CreateHttpNitroClient(rpcHost + "/api/v1").then((c) => {
        setNitroClient(c);
        fetchAndSetLedgerChannels(c, setLedgerChannels);
        c.Notifications.on("objective_completed", () =>
          fetchAndSetLedgerChannels(c, setLedgerChannels)
        );
      });
    });
  }, []);

  const focusedChannelInLedgerChannels = ledgerChannels.some(
    (lc) => lc.ID === focusedLedgerChannel
  );
  if (!focusedChannelInLedgerChannels && ledgerChannels.length > 0) {
    setFocusedLedgerChannel(ledgerChannels[0].ID);
  }

  return (
    <>
      <TopBar
        ledgerChannels={ledgerChannels}
        focusedLedgerChannel={focusedLedgerChannel}
        setFocusedLedgerChannel={setFocusedLedgerChannel}
      />
      <div style={{ display: "flex", justifyContent: "space-around" }}>
        <LedgerChannelDetails
          nitroClient={nitroClient}
          channelId={focusedLedgerChannel}
        />
        <PaymentChannelContainer
          nitroClient={nitroClient}
          ledgerChannel={focusedLedgerChannel}
        />
      </div>
    </>
  );
}

export default App;

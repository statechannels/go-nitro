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

function App() {
  const rpcPort = String(Number(window.location.port) - 100); // TODO need a better solution for finding this. We could just use the same port.

  const [nitroClient, setNitroClient] = useState<NitroRpcClient | null>(null);
  const [ledgerChannels, setLedgerChannels] = useState<LedgerChannelInfo[]>([]);
  const [focusedLedgerChannel, setFocusedLedgerChannel] = useState<string>("");

  useEffect(() => {
    NitroRpcClient.CreateHttpNitroClient(
      window.location.hostname + ":" + rpcPort + "/api/v1"
    ).then((c) => setNitroClient(c));
  });

  useEffect(() => {
    if (nitroClient) {
      fetchAndSetLedgerChannels(nitroClient, setLedgerChannels);
      nitroClient?.Notifications.on("objective_completed", () =>
        fetchAndSetLedgerChannels(nitroClient, setLedgerChannels)
      );
    }
  }, [nitroClient]);

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

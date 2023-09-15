import { useEffect, useState } from "react";
import { NitroRpcClient } from "@statechannels/nitro-rpc-client";
import { LedgerChannelInfo } from "@statechannels/nitro-rpc-client/src/types";

import "./App.css";
import TopBar from "./components/TopBar";
import LedgerChannelDetails from "./components/LedgerChannelDetails";
import PaymentChannelContainer from "./components/PaymentChannelContainer";
import { getRpcHost, fetchAndSetLedgerChannels } from "./utils";

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

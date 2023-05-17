import { useEffect, useState } from "react";
import { NitroRpcClient } from "@statechannels/nitro-rpc-client";
import { LedgerChannelInfo } from "@statechannels/nitro-rpc-client/src/types";

import "./App.css";
import TopBar from "./components/TopBar";
import { QUERY_KEY } from "./constants";
import PaymentChannelList from "./components/PaymentChannelList";
import PaymentChannelDetails from "./components/PaymentChannelDetails";
import LedgerChannelDetails from "./components/LedgerChannelDetails";

const paymentChannels = [
  {
    ID: "0x9823fa3d37ec304f90d1bef2c03c1fc70f86b6417f022d5e9ab88902a874f0cc",
  },
  {
    ID: "0x06a508ca629080f81954bb4dcce6b71f1d8de0dded88d333c720d3b9d4067af0",
  },
  {
    ID: "0x06a508ca629080f81954bb4dcce6b71f1d8de0dded88d333c720d3b9d4067af1",
  },
];

function App() {
  const url =
    new URLSearchParams(window.location.search).get(QUERY_KEY) ??
    "localhost:4005";
  const [nitroClient, setNitroClient] = useState<NitroRpcClient | null>(null);
  const [version, setVersion] = useState("");
  const [address, setAddress] = useState("");
  const [ledgerChannels, setLedgerChannels] = useState<LedgerChannelInfo[]>([]);
  const [focusedLedgerChannel, setFocusedLedgerChannel] = useState<string>(
    "0x9823fa3d37ec304f90d1bef2c03c1fc70f86b6417f022d5e9ab88902a874f0cc"
  );

  const [focusedPaymentChannel, setFocusedPaymentChannel] = useState<string>(
    "0x9823fa3d37ec304f90d1bef2c03c1fc70f86b6417f022d5e9ab88902a874f0cc"
  );

  useEffect(() => {
    NitroRpcClient.CreateHttpNitroClient(url).then((c) => setNitroClient(c));
  }, [url]);

  useEffect(() => {
    if (nitroClient) {
      nitroClient.GetVersion().then((v) => setVersion(v));
      nitroClient.GetAddress().then((a) => setAddress(a));
      nitroClient.GetAllLedgerChannels().then((l) => setLedgerChannels(l));
    }
  }, [nitroClient]);

  return (
    <>
      <TopBar
        url={url}
        ledgerChannels={ledgerChannels}
        focusedLedgerChannel={focusedLedgerChannel}
        setFocusedLedgerChannel={setFocusedLedgerChannel}
      />
      <div style={{ display: "flex", justifyContent: "space-around" }}>
        <LedgerChannelDetails version={version} url={url} address={address} />
        <PaymentChannelList
          paymentChannels={paymentChannels}
          focusedPaymentChannel={focusedPaymentChannel}
          setFocusedPaymentChannel={setFocusedPaymentChannel}
        />
        <PaymentChannelDetails paymentChannel={focusedPaymentChannel} />
      </div>
    </>
  );
}

export default App;

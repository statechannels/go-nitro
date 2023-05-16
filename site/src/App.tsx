import { useEffect, useState } from "react";
import { NitroRpcClient } from "@statechannels/nitro-rpc-client";

import { NetworkBalance } from "./components/NetworkBalance";
import "./App.css";
import TopBar from "./components/TopBar";
import { QUERY_KEY } from "./constants";

function App() {
  const url =
    new URLSearchParams(window.location.search).get(QUERY_KEY) ??
    "localhost:4005";
  const [nitroClient, setNitroClient] = useState<NitroRpcClient | null>(null);
  const [version, setVersion] = useState("");
  const [address, setAddress] = useState("");

  useEffect(() => {
    NitroRpcClient.CreateHttpNitroClient(url).then((c) => setNitroClient(c));
  }, [url]);

  useEffect(() => {
    if (nitroClient) {
      nitroClient.GetVersion().then((v) => setVersion(v));
      nitroClient.GetAddress().then((a) => setAddress(a));
    }
  }, [nitroClient]);

  return (
    <>
      <TopBar url={url} />
      <div style={{ display: "flex", justifyContent: "space-around" }}>
        <div className="card">
          <NetworkBalance
            status="running"
            lockedBalances={[]}
            myBalanceFree={BigInt(50)}
            theirBalanceFree={BigInt(200)}
          ></NetworkBalance>
          <p>Version: {version}</p>
          <p> Url: {url}</p>
          <p> Address: {address}</p>
        </div>
        <div>Payment channel list</div>
        <div>Payment channel details</div>
      </div>
    </>
  );
}

export default App;

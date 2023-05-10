import { useEffect, useState } from "react";
import { NitroRpcClient } from "@statechannels/nitro-rpc-client";

import { NetworkBalance } from "./components/NetworkBalance";
import statechannelsLogo from "./assets/statechannels.svg";
import "./App.css";

function App() {
  const url = "localhost:4005";
  const [nitroClient, setNitroClient] = useState<NitroRpcClient | null>(null);
  const [version, setVersion] = useState("");
  const [address, setAddress] = useState("");

  useEffect(() => {
    NitroRpcClient.CreateHttpNitroClient(url).then((c) => setNitroClient(c));
  }, []);

  useEffect(() => {
    if (nitroClient) {
      nitroClient.GetVersion().then((v) => setVersion(v));
      nitroClient.GetAddress().then((a) => setAddress(a));
    }
  }, [nitroClient]);

  return (
    <>
      <div>
        <a href="http://statechannels.org" className="href">
          <img src={statechannelsLogo} className="logo" />
        </a>
      </div>
      <h1>Vite + React + StateChannels</h1>
      <div className="card">
        <NetworkBalance
          status="running"
          lockedBalances={[]}
          myBalanceFree={BigInt(50)}
          theirBalanceFree={BigInt(200)}
        ></NetworkBalance>

        <p>The nitro client version is {version}</p>
        <p>
          The nitro node at {url} has address {address}
        </p>
      </div>
      <p className="read-the-docs">
        Click on the Vite and React logos to learn more
      </p>
    </>
  );
}

export default App;

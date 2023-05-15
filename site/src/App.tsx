import { useEffect, useState } from "react";
import { NitroRpcClient } from "@statechannels/nitro-rpc-client";

import { NetworkBalance } from "./components/NetworkBalance";
import "./App.css";
import RpcConnect from "./components/RpcConnect";

function App() {
  const [url, setUrl] = useState("localhost:4005");
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
      <RpcConnect url={url} setUrl={setUrl} />
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
    </>
  );
}

export default App;

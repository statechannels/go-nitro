import { useEffect, useState } from "react";
import { NitroRpcClient } from "@statechannels/nitro-rpc-client";

import { NetworkBalance } from "./components/NetworkBalance";
import statechannelsLogo from "./assets/statechannels.svg";
import "./App.css";

function App() {
  const [nitroClient, setNitroClient] = useState<NitroRpcClient | null>(null);
  const [version, setVersion] = useState("");

  useEffect(() => {
    const setupClient = async () => {
      const nitroClient = await NitroRpcClient.CreateHttpNitroClient(
        "http://localhost:4005"
      );
      setNitroClient(nitroClient);
    };
    setupClient();
  }, []);

  useEffect(() => {
    const getVersion = async () => {
      if (nitroClient) {
        console.log("getting version");
        const version = await nitroClient.GetVersion();
        setVersion(version);
      }
    };
    getVersion();
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
      </div>
      <p className="read-the-docs">
        Click on the Vite and React logos to learn more
      </p>
    </>
  );
}

export default App;

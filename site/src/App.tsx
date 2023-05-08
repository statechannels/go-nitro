import reactLogo from "./assets/react.svg";
import statechannelsLogo from "./assets/statechannels.svg";
import "./App.css";
import { NetworkBalance } from "./components/NetworkBalance";

function App() {
  return (
    <>
      <div>
        <a href="https://react.dev" target="_blank">
          <img src={reactLogo} className="logo react" alt="React logo" />
        </a>
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

        <p>
          Edit <code>src/App.tsx</code> and save to test HMR
        </p>
      </div>
      <p className="read-the-docs">
        Click on the Vite and React logos to learn more
      </p>
    </>
  );
}

export default App;

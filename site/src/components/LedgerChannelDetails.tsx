import { NetworkBalance } from "./NetworkBalance";

export default function LedgerChannelDetails({
  version,
  url,
  address,
}: {
  version: string;
  url: string;
  address: string;
}) {
  return (
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
  );
}

import LedgerChannelList from "./LedgerChannelList";
import RpcConnect, { RPCConnectProps } from "./RpcConnect";

export default function TopBar(props: RPCConnectProps) {
  return (
    <div style={{ display: "flex", justifyContent: "space-between" }}>
      <LedgerChannelList
        ledgerChannels={[
          {
            ID: "0x9823fa3d37ec304f90d1bef2c03c1fc70f86b6417f022d5e9ab88902a874f0cc",
          },
          {
            ID: "0x06a508ca629080f81954bb4dcce6b71f1d8de0dded88d333c720d3b9d4067af0",
          },
          {
            ID: "0x06a508ca629080f81954bb4dcce6b71f1d8de0dded88d333c720d3b9d4067af1",
          },
        ]}
      />
      <RpcConnect {...props} />
    </div>
  );
}

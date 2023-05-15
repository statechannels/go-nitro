import RpcConnect, { RPCConnectProps } from "./RpcConnect";

export default function TopBar(props: RPCConnectProps) {
  return (
    <div>
      <RpcConnect {...props} />
    </div>
  );
}

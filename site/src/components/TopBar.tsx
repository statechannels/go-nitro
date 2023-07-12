import Box from "@mui/material/Box";

import LedgerChannelList, { LedgerChannelListProps } from "./LedgerChannelList";
import RpcConnect, { RPCConnectProps } from "./RpcConnect";

type Props = LedgerChannelListProps & RPCConnectProps;

export default function TopBar(props: Props) {
  return (
    <Box
      sx={{
        display: "flex",
        justifyContent: "space-between",
        borderBottom: 1,
        borderColor: "divider",
      }}
    >
      <LedgerChannelList
        ledgerChannels={props.ledgerChannels}
        focusedLedgerChannel={props.focusedLedgerChannel}
        setFocusedLedgerChannel={props.setFocusedLedgerChannel}
      />
      <RpcConnect {...props} />
    </Box>
  );
}

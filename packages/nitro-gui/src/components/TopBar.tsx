import Box from "@mui/material/Box";

import colors from "../styles/colors.module.css";

import LedgerChannelList, { LedgerChannelListProps } from "./LedgerChannelList";

type Props = LedgerChannelListProps;

export default function TopBar(props: Props) {
  return (
    <Box
      sx={{
        display: "flex",
        justifyContent: "space-between",
        borderBottom: 1,
        borderColor: "divider",
        backgroundColor: colors.cBlue,
      }}
    >
      <LedgerChannelList
        ledgerChannels={props.ledgerChannels}
        focusedLedgerChannel={props.focusedLedgerChannel}
        setFocusedLedgerChannel={props.setFocusedLedgerChannel}
      />
    </Box>
  );
}

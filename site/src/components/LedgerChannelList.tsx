import Tab from "@mui/material/Tab";
import Tabs from "@mui/material/Tabs";

type LedgerChannel = {
  ID: string;
};

type Props = {
  ledgerChannels: LedgerChannel[];
};

function formatId(id: string): string {
  return id.slice(0, 8);
}

export default function LedgerChannelList({ ledgerChannels }: Props) {
  return (
    <Tabs>
      {ledgerChannels.map((ledgerChannel) => (
        <Tab label={formatId(ledgerChannel.ID)} />
      ))}
    </Tabs>
  );
}

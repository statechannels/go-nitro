import Button from "@mui/material/Button";

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
    <ul style={{ display: "flex" }}>
      {ledgerChannels.map((ledgerChannel) => (
        <li key={ledgerChannel.ID}>
          <Button variant="text">{formatId(ledgerChannel.ID)}</Button>
        </li>
      ))}
    </ul>
  );
}

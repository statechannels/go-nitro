import Tab from "@mui/material/Tab";
import Tabs from "@mui/material/Tabs";

import { PaymentChannel } from "../types";

type Props = {
  paymentChannels: PaymentChannel[];
  focusedPaymentChannel: string;
  setFocusedPaymentChannel: (channel: string) => void;
};

function formatPaymentChannel(chan: PaymentChannel): string {
  return chan.ID.slice(0, 8);
}

function focusedIndex(id: string, ids: PaymentChannel[]): number {
  const index = ids.findIndex((c) => c.ID === id);
  if (index != -1) {
    return index;
  }
  // The channel id is not found in the channel list.
  return 0;
}

function handleChange(
  ledgerChannels: PaymentChannel[],
  setter: (id: string) => void
) {
  return (_: React.SyntheticEvent, newValue: number) => {
    setter(ledgerChannels[newValue].ID);
  };
}

export default function PaymentChannelList(props: Props) {
  return (
    <Tabs
      value={focusedIndex(props.focusedPaymentChannel, props.paymentChannels)}
      onChange={handleChange(
        props.paymentChannels,
        props.setFocusedPaymentChannel
      )}
      orientation="vertical"
    >
      {props.paymentChannels.map((chan) => (
        <Tab key={chan.ID} label={formatPaymentChannel(chan)} />
      ))}
    </Tabs>
  );
}

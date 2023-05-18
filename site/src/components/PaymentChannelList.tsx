import Tab from "@mui/material/Tab";
import Tabs from "@mui/material/Tabs";
import { PaymentChannelInfo } from "@statechannels/nitro-rpc-client/src/types";

type Props = {
  paymentChannels: PaymentChannelInfo[];
  focusedPaymentChannel: string;
  setFocusedPaymentChannel: (channel: string) => void;
};

function formatPaymentChannel(chan: PaymentChannelInfo): string {
  return chan.ID.slice(0, 8);
}

function focusedIndex(id: string, ids: PaymentChannelInfo[]): number {
  const index = ids.findIndex((c) => c.ID === id);
  if (index != -1) {
    return index;
  }
  // The channel id is not found in the channel list.
  return 0;
}

function handleChange(
  ledgerChannels: PaymentChannelInfo[],
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

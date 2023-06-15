import Tab from "@mui/material/Tab";
import Tabs from "@mui/material/Tabs";
import { PaymentChannelInfo } from "@statechannels/nitro-rpc-client/src/types";
import { FC } from "react";

interface PaymentChannelListProps {
  paymentChannels: PaymentChannelInfo[];
  focusedPaymentChannel: string;
  setFocusedPaymentChannel: (channel: PaymentChannelInfo) => void;
}

function formatPaymentChannel(chan: PaymentChannelInfo): string {
  return chan.ID.slice(0, 8);
}

const PaymentChannelList: FC<PaymentChannelListProps> = ({
  paymentChannels,
  focusedPaymentChannel,
  setFocusedPaymentChannel,
}: PaymentChannelListProps) => {
  const handleChange = (_: object, value: number) => {
    setFocusedPaymentChannel(paymentChannels[value]);
  };

  const focusedIndex = paymentChannels.findIndex(
    (c) => c.ID === focusedPaymentChannel
  );

  return (
    <Tabs value={focusedIndex} onChange={handleChange} orientation="vertical">
      {paymentChannels.map((chan) => (
        <Tab key={chan.ID} label={formatPaymentChannel(chan)} />
      ))}
    </Tabs>
  );
};

export default PaymentChannelList;

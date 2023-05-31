import Tab from "@mui/material/Tab";
import Tabs from "@mui/material/Tabs";
import { PaymentChannelInfo } from "@statechannels/nitro-rpc-client/src/types";
import { FC, useCallback, useMemo } from "react";

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
  const handleChange = useCallback(
    (_: object, value: number) => {
      setFocusedPaymentChannel(paymentChannels[value]);
    },
    [paymentChannels, setFocusedPaymentChannel]
  );

  const focusedIndex = useMemo(() => {
    const index = paymentChannels.findIndex(
      (c) => c.ID === focusedPaymentChannel
    );
    if (index != -1) {
      return index;
    }
    // The channel id is not found in the channel list.
    return 0;
  }, [paymentChannels, focusedPaymentChannel]);

  return (
    <Tabs value={focusedIndex} onChange={handleChange} orientation="vertical">
      {paymentChannels.map((chan) => (
        <Tab key={chan.ID} label={formatPaymentChannel(chan)} />
      ))}
    </Tabs>
  );
};

export default PaymentChannelList;

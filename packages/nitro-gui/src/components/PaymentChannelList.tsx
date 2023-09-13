import Tab from "@mui/material/Tab";
import Tabs from "@mui/material/Tabs";
import { PaymentChannelInfo } from "@statechannels/nitro-rpc-client/src/types";
import { FC } from "react";
import { makeStyles } from "tss-react/mui";

import colors from "../styles/colors.module.css";

const useStyles = makeStyles()(() => ({
  paymentList: {
    backgroundColor: colors.cBlue,
  },
}));

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
  const { classes } = useStyles();

  const handleChange = (_: object, value: number) => {
    setFocusedPaymentChannel(paymentChannels[value]);
  };

  const focusedIndex = paymentChannels.findIndex(
    (c) => c.ID === focusedPaymentChannel
  );

  return (
    <Tabs
      className={classes.paymentList}
      value={focusedIndex}
      onChange={handleChange}
      orientation="vertical"
    >
      {paymentChannels.map((chan) => (
        <Tab key={chan.ID} label={formatPaymentChannel(chan)} />
      ))}
    </Tabs>
  );
};

export default PaymentChannelList;

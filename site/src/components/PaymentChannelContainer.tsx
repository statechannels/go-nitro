import { NitroRpcClient } from "@statechannels/nitro-rpc-client";
import { PaymentChannelInfo } from "@statechannels/nitro-rpc-client/src/types";
import { useEffect, useState, FC } from "react";
import { makeStyles } from "tss-react/mui";

import PaymentChannelList from "./PaymentChannelList";
import PaymentChannelDetails from "./PaymentChannelDetails";

interface Props {
  nitroClient: NitroRpcClient | null;
  ledgerChannel: string;
}

const useStyles = makeStyles()(() => ({
  paymentDetails: {
    padding: "4rem",
  },
}));

const PaymentChannelContainer: FC<Props> = ({
  nitroClient,
  ledgerChannel,
}: Props) => {
  const [paymentChannels, setPaymentChannels] = useState<PaymentChannelInfo[]>(
    []
  );
  const { classes } = useStyles();

  const [focusedPaymentChannel, setFocusedPaymentChannel] =
    useState<string>("");

  useEffect(() => {
    if (nitroClient && ledgerChannel) {
      nitroClient.GetPaymentChannelsByLedger(ledgerChannel).then((p) => {
        setPaymentChannels(p);
        if (p.length > 0) {
          setFocusedPaymentChannel(p[0].ID);
        }
      });
    }
  }, [nitroClient, ledgerChannel]);

  return (
    <>
      <PaymentChannelList
        paymentChannels={paymentChannels}
        focusedPaymentChannel={focusedPaymentChannel}
        setFocusedPaymentChannel={setFocusedPaymentChannel}
      />
      <div className={classes.paymentDetails}>
        <PaymentChannelDetails
          channelID={focusedPaymentChannel}
          counterparty={"0x123"}
          capacity={1000}
          myBalance={150}
          status={"running"}
        />
      </div>
    </>
  );
};

export default PaymentChannelContainer;

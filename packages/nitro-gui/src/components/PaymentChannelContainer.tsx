import { NitroRpcClient } from "@statechannels/nitro-rpc-client";
import {
  ChannelStatus,
  PaymentChannelBalance,
  PaymentChannelInfo,
} from "@statechannels/nitro-rpc-client/src/types";
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

const DEFAULT_CHANNEL: PaymentChannelInfo = {
  ID: "",
  Status: "Complete" as ChannelStatus,
  Balance: {
    AssetAddress: "",
    Payee: "",
    Payer: "",
    PaidSoFar: BigInt(0),
    RemainingFunds: BigInt(0),
  } as PaymentChannelBalance,
};

const PaymentChannelContainer: FC<Props> = ({
  nitroClient,
  ledgerChannel,
}: Props) => {
  const [paymentChannels, setPaymentChannels] = useState<PaymentChannelInfo[]>(
    []
  );
  const { classes } = useStyles();

  const [focusedPaymentChannel, setFocusedPaymentChannel] =
    useState<PaymentChannelInfo>(DEFAULT_CHANNEL);

  useEffect(() => {
    async function handleGetPaymentChannels(
      client: NitroRpcClient,
      channel: string
    ) {
      const paymentChannelsList = await client.GetPaymentChannelsByLedger(
        channel
      );

      setPaymentChannels(paymentChannelsList);
      if (paymentChannelsList.length) {
        setFocusedPaymentChannel(paymentChannelsList[0]);
      }
    }

    if (nitroClient && ledgerChannel) {
      handleGetPaymentChannels(nitroClient, ledgerChannel);
    }
  }, [nitroClient, ledgerChannel]);

  return (
    <>
      <PaymentChannelList
        paymentChannels={paymentChannels}
        focusedPaymentChannel={focusedPaymentChannel?.ID || ""}
        setFocusedPaymentChannel={setFocusedPaymentChannel}
      />
      <div className={classes.paymentDetails}>
        <PaymentChannelDetails
          myAddress={nitroClient?.GetAddress}
          channelID={focusedPaymentChannel.ID}
          payer={focusedPaymentChannel.Balance.Payer}
          payee={focusedPaymentChannel.Balance.Payee}
          remainingFunds={focusedPaymentChannel.Balance.RemainingFunds}
          paidSoFar={focusedPaymentChannel.Balance.PaidSoFar}
          status={focusedPaymentChannel?.Status}
        />
      </div>
    </>
  );
};

export default PaymentChannelContainer;

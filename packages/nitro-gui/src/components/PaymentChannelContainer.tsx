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

export enum PaymentChannelType {
  inbound,
  outbound,
  mediated,
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

  const [myAddress, setMyAddress] = useState("");

  useEffect(() => {
    if (nitroClient) {
      nitroClient.GetAddress().then((a) => setMyAddress(a));
    }
  }, [nitroClient]);

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

  const payer = focusedPaymentChannel.Balance.Payer;
  const payee = focusedPaymentChannel.Balance.Payee;

  const inferType = () => {
    if (myAddress.toLowerCase() == payer.toLowerCase()) {
      return PaymentChannelType.outbound;
    } else if (myAddress.toLowerCase() == payee.toLowerCase()) {
      return PaymentChannelType.inbound;
    } else return PaymentChannelType.mediated;
  };
  const pcT: PaymentChannelType = inferType();

  return (
    <>
      <PaymentChannelList
        paymentChannels={paymentChannels}
        focusedPaymentChannel={focusedPaymentChannel?.ID || ""}
        setFocusedPaymentChannel={setFocusedPaymentChannel}
      />
      <div className={classes.paymentDetails}>
        <PaymentChannelDetails
          type={pcT}
          channelID={focusedPaymentChannel.ID}
          payer={payer}
          payee={payee}
          remainingFunds={focusedPaymentChannel.Balance.RemainingFunds}
          paidSoFar={focusedPaymentChannel.Balance.PaidSoFar}
          status={focusedPaymentChannel?.Status}
        />
      </div>
    </>
  );
};

export default PaymentChannelContainer;

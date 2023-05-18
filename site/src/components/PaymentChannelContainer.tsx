import { NitroRpcClient } from "@statechannels/nitro-rpc-client";
import { PaymentChannelInfo } from "@statechannels/nitro-rpc-client/src/types";
import { useEffect, useState } from "react";

import PaymentChannelList from "./PaymentChannelList";
import PaymentChannelDetails from "./PaymentChannelDetails";

type Props = {
  nitroClient: NitroRpcClient | null;
  ledgerChannel: string;
};
export default function PaymentChannelContainer({
  nitroClient,
  ledgerChannel,
}: Props) {
  const [paymentChannels, setPaymentChannels] = useState<PaymentChannelInfo[]>(
    []
  );

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
      <PaymentChannelDetails paymentChannel={focusedPaymentChannel} />
    </>
  );
}

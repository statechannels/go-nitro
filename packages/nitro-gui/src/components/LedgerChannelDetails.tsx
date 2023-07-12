import { NitroRpcClient } from "@statechannels/nitro-rpc-client";
import { useEffect, useState } from "react";
import { LedgerChannelBalance } from "@statechannels/nitro-rpc-client/src/types";

import { NetworkBalance, VirtualChannelBalanceProps } from "./NetworkBalance";

type Props = {
  nitroClient: NitroRpcClient | null;
  channelId: string;
};

type LedgerDetails = {
  ledgerBalance: LedgerChannelBalance;
  lockedBalances: VirtualChannelBalanceProps[];
};

async function getLedgerDetails(
  nitroClient: NitroRpcClient,
  channelId: string
): Promise<LedgerDetails> {
  const ledgerChannel = await nitroClient.GetLedgerChannel(channelId);
  const paymentChannels = await nitroClient.GetPaymentChannelsByLedger(
    channelId
  );

  const lockedBalances = paymentChannels.map((pc) => {
    const total = pc.Balance.PaidSoFar + pc.Balance.RemainingFunds;

    return {
      budget: total,
      // bigint division outputs a bigint. To have decimal precision,
      // we multiply the numerator by a factor, divide, convert to JS number,
      // and then divide by the same factor.
      myPercentage: Number((pc.Balance.RemainingFunds * 1_000n) / total) / 1000,
    };
  });

  return {
    ledgerBalance: ledgerChannel.Balance,
    lockedBalances,
  };
}

export default function LedgerChannelDetails({
  nitroClient,
  channelId,
}: Props) {
  const [ledgerDetails, setLedgerDetails] = useState<LedgerDetails | null>(
    null
  );
  const [myAddress, setMyAddress] = useState("");

  useEffect(() => {
    if (nitroClient && channelId != "") {
      getLedgerDetails(nitroClient, channelId).then(setLedgerDetails);
    }
  }, [nitroClient, channelId]);

  useEffect(() => {
    if (nitroClient) {
      nitroClient.GetAddress().then((a) => setMyAddress(a));
    }
  }, [nitroClient]);

  const myBalance = ledgerDetails?.ledgerBalance.MyBalance ?? BigInt(0);
  const theirBalance = ledgerDetails?.ledgerBalance.TheirBalance ?? BigInt(0);

  return (
    <div className="card">
      <div> My Address: {myAddress}</div>
      <NetworkBalance
        status="running"
        lockedBalances={ledgerDetails?.lockedBalances ?? []}
        myBalanceFree={myBalance}
        theirBalanceFree={theirBalance}
      />
    </div>
  );
}

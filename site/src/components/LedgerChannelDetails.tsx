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

  // todo: remove this construction once big number over rpc is fixed
  // At the moment, fields that are defined as big numbers are populated as numbers
  const ledgerBalance = {
    ...ledgerChannel.Balance,
    ClientBalance: BigInt(ledgerChannel.Balance.ClientBalance),
    HubBalance: BigInt(ledgerChannel.Balance.HubBalance),
  };

  const lockedBalances = paymentChannels.map((pc) => {
    const total = pc.Balance.PaidSoFar + pc.Balance.RemainingFunds;
    return {
      // todo: this conversion should be removed once big number over rpc is fixed
      budget: BigInt(total),
      myPercentage: Number(pc.Balance.RemainingFunds / total),
    };
  });

  return {
    ledgerBalance,
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

  const myBalance = ledgerDetails?.ledgerBalance.ClientBalance ?? BigInt(0);
  const theirBalance = ledgerDetails?.ledgerBalance.HubBalance ?? BigInt(0);

  return (
    <div className="card">
      <div> My Address: {myAddress}</div>
      <NetworkBalance
        status="running"
        lockedBalances={ledgerDetails?.lockedBalances ?? []}
        myBalanceFree={myBalance}
        theirBalanceFree={theirBalance}
      ></NetworkBalance>
    </div>
  );
}

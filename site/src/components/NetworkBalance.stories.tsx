import { Meta, StoryObj } from "@storybook/react";

import {
  NetworkBalance,
  NetworkBalanceProps,
  VirtualChannelBalanceProps,
} from "./NetworkBalance";

const meta: Meta<typeof NetworkBalance> = {
  title: "NetworkBalance",
  component: NetworkBalance,
};
export default meta;

type NB = StoryObj<NetworkBalanceProps>;

export const Zeros: NB = {
  args: {
    myBalanceFree: 0n,
    theirBalanceFree: 0n,
    status: "running",
    lockedBalances: [],
  },
};

export const EvenStart: NB = {
  args: {
    myBalanceFree: 10n ** 18n,
    theirBalanceFree: 10n ** 18n,
    status: "running",
    lockedBalances: [],
  },
};

export const ClientStart: NB = {
  args: {
    myBalanceFree: 100n,
    theirBalanceFree: 0n,
    status: "running",
    lockedBalances: [],
  },
};

export const ClientMid: NB = {
  args: {
    myBalanceFree: 70n,
    theirBalanceFree: 10n,
    status: "running",
    lockedBalances: [],
  },
};

export const ProviderStart: NB = {
  args: {
    myBalanceFree: 0n,
    theirBalanceFree: 100n,
    status: "running",
    lockedBalances: [],
  },
};

export const ProviderMid: NB = {
  args: {
    myBalanceFree: 15n,
    theirBalanceFree: 60n,
    status: "running",
    lockedBalances: [],
  },
};

export const TwoChannels: NB = {
  args: {
    myBalanceFree: 47n,
    theirBalanceFree: 100n,
    status: "running",
    lockedBalances: [
      {
        budget: 10n,
        myPercentage: 0.5,
      },
      {
        budget: 30n,
        myPercentage: 0.2,
      },
    ],
  },
};

export const SomeChannels: NB = {
  args: {
    myBalanceFree: 50n,
    theirBalanceFree: 100n,
    status: "running",
    lockedBalances: randomChannels(5, 100n),
  },
};

export const ManyChannels: NB = {
  args: {
    myBalanceFree: 345n,
    theirBalanceFree: 123n,
    status: "running",
    lockedBalances: randomChannels(15, 150n),
  },
};

export const UnresponsivePeer: NB = {
  args: {
    myBalanceFree: 25n,
    theirBalanceFree: 65n,
    status: "unresponsive-peer",
    lockedBalances: randomChannels(5, 100n),
  },
};

export const UnderChallenge: NB = {
  args: {
    myBalanceFree: 83n,
    theirBalanceFree: 24n,
    status: "under-challenge",
    lockedBalances: randomChannels(5, 100n),
  },
};

function randomChannels(
  numChannels: number,
  budgetCeiling: bigint
): VirtualChannelBalanceProps[] {
  const channels = [];
  for (let i = 0; i < numChannels; i++) {
    channels.push(randomChannel(budgetCeiling));
  }
  return channels;
}

function randomChannel(budgetCeiling: bigint): VirtualChannelBalanceProps {
  return {
    budget: BigInt(Math.floor(Math.random() * Number(budgetCeiling))),
    myPercentage: Math.random(),
  };
}

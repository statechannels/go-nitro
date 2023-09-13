import type { Meta, StoryObj } from "@storybook/react";
import { ChannelStatus } from "@statechannels/nitro-rpc-client/src/types";
import { FC } from "react";

import PaymentChannelDetails from "./PaymentChannelDetails";
import { PaymentChannelType } from "./PaymentChannelContainer";

const wrappedComponent: FC<{ t: PaymentChannelType }> = (p) => {
  return (
    <PaymentChannelDetails
      type={p.t}
      channelID="fa745d81208c3f9f394a04db57a27f11c46be1d6dce0f81dd2852347d83fe4e4"
      payer="b25e8dc6f4795e9441b3e0b2519f2c9c827eb734"
      payee="b55e8dc6f4795e9441b3e0b2519f2c9c827eb734"
      remainingFunds={BigInt(850)}
      paidSoFar={BigInt(150)}
      status={"Ready" as ChannelStatus}
    />
  );
};

const meta: Meta<typeof wrappedComponent> = {
  title: "PaymentChannelDetails",
  component: wrappedComponent,
};

export default meta;

type Story = StoryObj<typeof wrappedComponent>;

export const PaymentChannelDetailsComponent: Story = {
  argTypes: {
    t: {
      options: [
        PaymentChannelType.inbound,
        PaymentChannelType.outbound,
        PaymentChannelType.mediated,
      ],
      control: { type: "radio" },
      defaultValue: PaymentChannelType.inbound,
    },
  },
};

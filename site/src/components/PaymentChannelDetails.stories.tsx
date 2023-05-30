import type { Meta, StoryObj } from "@storybook/react";

import PaymentChannelDetails from "./PaymentChannelDetails";

const meta: Meta<typeof PaymentChannelDetails> = {
  title: "PaymentChannelDetails",
  component: PaymentChannelDetails,
};

export default meta;

type Story = StoryObj<typeof PaymentChannelDetails>;

export const PaymentChannelDetailsComponent: Story = {
  render: () => (
    <PaymentChannelDetails
      channelID="0x1234"
      counterparty="0x123"
      capacity={1000}
      myBalance={150}
      status="running"
    />
  ),
};

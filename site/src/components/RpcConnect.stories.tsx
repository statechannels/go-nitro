import type { Meta, StoryObj } from "@storybook/react";

import RpcConnect from "./RpcConnect";

const meta: Meta<typeof RpcConnect> = {
  title: "RpcConnect",
  component: RpcConnect,
};

export default meta;

type Story = StoryObj<typeof RpcConnect>;

export const Primary: Story = {
  render: () => <RpcConnect url="localhost:8545" />,
};

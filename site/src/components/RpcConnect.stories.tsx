import type { Meta, StoryObj } from "@storybook/react";

import RpcConnect from "./RpcConnect";

const meta: Meta<typeof RpcConnect> = {
  title: "RpcConnect",
  component: RpcConnect,
};

export default meta;

type Story = StoryObj<typeof RpcConnect>;

// eslint-disable-next-line @typescript-eslint/no-empty-function
const setUrl = function (_url: string) {};

export const Primary: Story = {
  render: () => <RpcConnect url="localhost:8545" setUrl={setUrl} />,
};

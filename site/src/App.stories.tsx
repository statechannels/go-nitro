import { rest } from "msw";
import { Meta } from "@storybook/react";
import { Server } from "mock-socket";
import { useEffect } from "react";

import App from "./App";
import {
  allLedgerChannelMock,
  getAddressMock,
  getLedgerChannelMock,
  getPaymentChannelsByLedgerMock,
} from "./mocks/request";

function createMockServer() {
  const mockServer = new Server("ws://localhost:4005/api/subscribe");
  // eslint-disable-next-line @typescript-eslint/no-empty-function
  mockServer.on("connection", () => {});
  return mockServer;
}

const meta: Meta<typeof App> = {
  title: "App",
  component: App,
};

export default meta;
export const AppPopulated = () => {
  const mockServer = createMockServer();

  // Clean up after the story is unmounted
  useEffect(() => {
    return () => {
      mockServer.stop();
    };
  }, [mockServer]);

  return <App />;
};

AppPopulated.parameters = {
  msw: {
    handlers: [
      rest.post("http://localhost:4005/api", async (req, res, ctx) => {
        const json = await req.json();
        let retVal = {};
        switch (json.method) {
          case "get_address":
            retVal = getAddressMock;
            break;
          case "get_all_ledger_channels":
            retVal = allLedgerChannelMock;
            break;
          case "get_ledger_channel":
            retVal = getLedgerChannelMock;
            break;
          case "get_payment_channels_by_ledger":
            retVal = getPaymentChannelsByLedgerMock;
            break;
          default:
            ctx.status(403);
            retVal = { errorMessage: "Method not found" };
        }
        return res(ctx.json(retVal));
      }),
    ],
  },
};

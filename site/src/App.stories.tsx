import { rest } from "msw";
import { Meta } from "@storybook/react";

import App from "./App";
import {
  allLedgerChannelMock,
  getAddressMock,
  getLedgerChannelMock,
  getPaymentChannelsByLedgerMock,
} from "./mocks/request";

const meta: Meta<typeof App> = {
  title: "App",
  component: App,
};
export default meta;
export const AppPopulated = () => <App />;
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

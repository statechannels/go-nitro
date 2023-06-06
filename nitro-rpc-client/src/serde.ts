import Ajv, { JTDDataType } from "ajv/dist/jtd";
const ajv = new Ajv();

import {
  ChannelStatus,
  LedgerChannelInfo,
  RPCRequestAndResponses,
  RequestMethod,
} from "./types";

const schema = {
  properties: {
    jsonrpc: { type: "string" },
    id: { type: "uint32" },
    result: {
      properties: {
        ID: { type: "string" },
        Status: { type: "string" },
        Balance: {
          properties: {
            AssetAddress: { type: "string" },
            Hub: { type: "string" },
            Client: { type: "string" },
            HubBalance: { type: "string" },
            ClientBalance: { type: "string" },
          },
        },
      },
    },
  },
} as const;
type LedgerChannelResponse = JTDDataType<typeof schema>;

export function validateResponse<T extends RequestMethod>(
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  response: any,
  method: T
): RPCRequestAndResponses[T][1]["result"] {
  let errors;
  switch (method) {
    case "get_ledger_channel": {
      const validate = ajv.compile<LedgerChannelResponse>(schema);
      if (validate(response)) {
        const convertedResponse = convertBalance(response);
        return convertedResponse;
      }
      errors = validate.errors;
      break;
    }
    default:
      throw new Error(`Unknown method: ${method}`);
  }
  throw new Error(`Invalid response: ${JSON.stringify(errors)}`);
}

function convertBalance(response: LedgerChannelResponse): LedgerChannelInfo {
  const result = response.result;
  // todo: validate channel status
  return {
    ...result,
    Status: result.Status as ChannelStatus,
    Balance: {
      ...result.Balance,
      HubBalance: BigInt(result.Balance.HubBalance),
      ClientBalance: BigInt(result.Balance.ClientBalance),
    },
  };
}

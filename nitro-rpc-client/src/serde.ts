import Ajv, { JTDDataType } from "ajv/dist/jtd";

import {
  ChannelStatus,
  LedgerChannelInfo,
  PaymentChannelInfo,
  RPCRequestAndResponses,
  RequestMethod,
} from "./types";

const ajv = new Ajv();

const jsonRpcSchema = {
  properties: {
    jsonrpc: { type: "string" },
    id: { type: "uint32" },
    result: {
      properties: {},
      additionalProperties: true,
    },
  },
  optionalProperties: { error: { type: "string", nullable: true } },
} as const;
type JsonRpcSchemaType = JTDDataType<typeof jsonRpcSchema>;

const objectiveSchema = {
  properties: {
    Id: { type: "string" },
    ChannelId: { type: "string" },
  },
} as const;
type ObjectiveSchemaType = JTDDataType<typeof objectiveSchema>;

const stringSchema = { type: "string" } as const;
type StringSchemaType = JTDDataType<typeof stringSchema>;

const ledgerChannelSchema = {
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
} as const;
type LedgerChannelSchemaType = JTDDataType<typeof ledgerChannelSchema>;

const paymentChannelSchema = {
  properties: {
    ID: { type: "string" },
    Status: { type: "string" },
    Balance: {
      properties: {
        AssetAddress: { type: "string" },
        Payee: { type: "string" },
        Payer: { type: "string" },
        PaidSoFar: { type: "string" },
        RemainingFunds: { type: "string" },
      },
    },
  },
} as const;
type PaymentChannelSchemaType = JTDDataType<typeof paymentChannelSchema>;

const paymentChannelsSchema = {
  elements: {
    ...paymentChannelSchema,
  },
} as const;
type PaymentChannelsSchemaType = JTDDataType<typeof paymentChannelsSchema>;

type ResponseSchema =
  | typeof objectiveSchema
  | typeof stringSchema
  | typeof ledgerChannelSchema
  | typeof paymentChannelSchema
  | typeof paymentChannelsSchema;
type ResponseSchemaType =
  | ObjectiveSchemaType
  | StringSchemaType
  | LedgerChannelSchemaType
  | PaymentChannelSchemaType
  | PaymentChannelsSchemaType;

export function validateResponse<T extends RequestMethod>(
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  response: any,
  method: T
): RPCRequestAndResponses[T][1]["result"] {
  const result = getJsonRpcResult(response);
  switch (method) {
    case "direct_fund":
    case "virtual_fund":
      return validateResult(
        objectiveSchema,
        result,
        (result: ObjectiveSchemaType) => result
      );
    case "direct_defund":
    case "version":
    case "get_address":
    case "virtual_defund":
      return validateResult(
        stringSchema,
        result,
        (result: StringSchemaType) => result
      );
    case "get_ledger_channel":
      return validateResult(
        ledgerChannelSchema,
        result,
        convertToInternalLedgerChannelType
      );
    case "get_payment_channel":
      return validateResult(
        paymentChannelSchema,
        result,
        convertToInternalPaymentChannelType
      );
    case "get_payment_channels_by_ledger":
      return validateResult(
        paymentChannelsSchema,
        result,
        convertToInternalPaymentChannelsType
      );
    default:
      throw new Error(`Unknown method: ${method}`);
  }
}

function validateResult<
  U extends ResponseSchemaType,
  S extends ResponseSchema,
  T extends RequestMethod
>(
  schema: S,
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  result: any,
  converstionFn: (result: U) => RPCRequestAndResponses[T][1]["result"]
): RPCRequestAndResponses[T][1]["result"] {
  const validate = ajv.compile<U>(schema);
  if (validate(result)) {
    return converstionFn(result);
  }
  throw new Error(
    `Error parsing json rpc result: ${JSON.stringify(
      validate.errors
    )}. The result is ${JSON.stringify(result)}`
  );
}

function convertToInternalLedgerChannelType(
  result: LedgerChannelSchemaType
): LedgerChannelInfo {
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

function convertToInternalPaymentChannelType(
  result: PaymentChannelSchemaType
): PaymentChannelInfo {
  // todo: validate channel status
  return {
    ...result,
    Status: result.Status as ChannelStatus,
    Balance: {
      ...result.Balance,
      PaidSoFar: BigInt(result.Balance.PaidSoFar),
      RemainingFunds: BigInt(result.Balance.RemainingFunds),
    },
  };
}

function convertToInternalPaymentChannelsType(
  result: PaymentChannelsSchemaType
): PaymentChannelInfo[] {
  return result.map((pc) => convertToInternalPaymentChannelType(pc));
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
function getJsonRpcResult(response: any): any {
  const validate = ajv.compile<JsonRpcSchemaType>(jsonRpcSchema);
  if (validate(response)) {
    return response.result;
  }
  throw new Error(
    `Invalid json rpc response: ${JSON.stringify(
      validate.errors
    )}. The response is ${JSON.stringify(response)}`
  );
}

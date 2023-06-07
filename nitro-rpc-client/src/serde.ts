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

const ledgerChannelsSchema = {
  elements: {
    ...ledgerChannelSchema,
  },
} as const;
type LedgerChannelsSchemaType = JTDDataType<typeof ledgerChannelsSchema>;

const paymentChannelsSchema = {
  elements: {
    ...paymentChannelSchema,
  },
} as const;
type PaymentChannelsSchemaType = JTDDataType<typeof paymentChannelsSchema>;

const paymentSchema = {
  properties: {
    Amount: { type: "uint32" },
    Channel: { type: "string" },
  },
} as const;
type PaymentSchemaType = JTDDataType<typeof paymentSchema>;

type ResponseSchema =
  | typeof objectiveSchema
  | typeof stringSchema
  | typeof ledgerChannelSchema
  | typeof ledgerChannelsSchema
  | typeof paymentChannelSchema
  | typeof paymentChannelsSchema
  | typeof paymentSchema;

type ResponseSchemaType =
  | ObjectiveSchemaType
  | StringSchemaType
  | LedgerChannelSchemaType
  | LedgerChannelsSchemaType
  | PaymentChannelSchemaType
  | PaymentChannelsSchemaType
  | PaymentSchemaType;

/**
 * Validates that the response is a valid JSON RPC response with a valid result
 * @param response - JSON RPC response
 * @param method - JSON RPC method
 * @returns The validated result of the JSON RPC response
 */
export function getAndValidateResult<T extends RequestMethod>(
  response: unknown,
  method: T
): RPCRequestAndResponses[T][1]["result"] {
  const result = getJsonRpcResult(response);
  switch (method) {
    case "direct_fund":
    case "virtual_fund":
      return validateAndConvertResult(
        objectiveSchema,
        result,
        (result: ObjectiveSchemaType) => result
      );
    case "direct_defund":
    case "version":
    case "get_address":
    case "virtual_defund":
      return validateAndConvertResult(
        stringSchema,
        result,
        (result: StringSchemaType) => result
      );
    case "get_ledger_channel":
      return validateAndConvertResult(
        ledgerChannelSchema,
        result,
        convertToInternalLedgerChannelType
      );

    case "get_all_ledger_channels":
      return validateAndConvertResult(
        ledgerChannelsSchema,
        result,
        convertToInternalLedgerChannelsType
      );
    case "get_payment_channel":
      return validateAndConvertResult(
        paymentChannelSchema,
        result,
        convertToInternalPaymentChannelType
      );
    case "get_payment_channels_by_ledger":
      return validateAndConvertResult(
        paymentChannelsSchema,
        result,
        convertToInternalPaymentChannelsType
      );
    case "pay":
      return validateAndConvertResult(
        paymentSchema,
        result,
        (result: PaymentSchemaType) => result
      );
    default:
      throw new Error(`Unknown method: ${method}`);
  }
}

/**
 * Validates that the response is a valid JSON RPC response and pulls out the result
 * @param response - JSON RPC response
 * @returns The result of the response
 */
function getJsonRpcResult(response: unknown): unknown {
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

/**
 * validateAndConvertResult validates that the result object conforms to the schema and converts it to the internal type
 *
 * @param schema - JSON Type Definition
 * @param result - Object to validate
 * @param converstionFn - Function to convert the valiated object to internal type
 * @returns A validated object of internal type
 */
function validateAndConvertResult<
  U extends ResponseSchemaType,
  S extends ResponseSchema,
  T extends RequestMethod
>(
  schema: S,
  result: unknown,
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

function convertToInternalLedgerChannelsType(
  result: LedgerChannelsSchemaType
): LedgerChannelInfo[] {
  return result.map((lc) => convertToInternalLedgerChannelType(lc));
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

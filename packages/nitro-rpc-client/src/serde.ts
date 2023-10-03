import Ajv, { JTDDataType } from "ajv/dist/jtd";

import {
  ChannelStatus,
  LedgerChannelInfo,
  PaymentChannelInfo,
  RPCNotification,
  RPCRequestAndResponses,
  RequestMethod,
} from "./types";

const ajv = new Ajv();

const jsonRpcSchema = {
  properties: {
    jsonrpc: { type: "string" },
    id: { type: "uint32" },
  },
  optionalProperties: {
    result: {
      nullable: true,
    },
    error: {
      properties: {
        code: { type: "int32" },
        message: { type: "string" },
      },
      additionalProperties: true,
      nullable: true,
    },
  },
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
        Them: { type: "string" },
        Me: { type: "string" },
        MyBalance: { type: "string" },
        TheirBalance: { type: "string" },
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

const voucherSchema = {
  properties: {
    ChannelId: { type: "string" },
    Amount: { type: "uint32" },
    Signature: {
      type: "string",
    },
  },
} as const;
type VoucherSchemaType = JTDDataType<typeof voucherSchema>;

const receiveVoucherSchema = {
  properties: {
    Total: { type: "string" },
    Delta: { type: "string" },
  },
} as const;

type ReceiveVoucherSchemaType = JTDDataType<typeof receiveVoucherSchema>;

type ResponseSchema =
  | typeof objectiveSchema
  | typeof stringSchema
  | typeof ledgerChannelSchema
  | typeof ledgerChannelsSchema
  | typeof paymentChannelSchema
  | typeof paymentChannelsSchema
  | typeof paymentSchema
  | typeof voucherSchema
  | typeof receiveVoucherSchema;

type ResponseSchemaType =
  | ObjectiveSchemaType
  | StringSchemaType
  | LedgerChannelSchemaType
  | LedgerChannelsSchemaType
  | PaymentChannelSchemaType
  | PaymentChannelsSchemaType
  | PaymentSchemaType
  | VoucherSchemaType
  | ReceiveVoucherSchemaType;

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
  const { result, error } = getJsonRpcResult(response);
  if (error) {
    throw new Error("jsonrpc response: " + error.message);
  }
  switch (method) {
    case "create_ledger_channel":
    case "create_payment_channel":
      return validateAndConvertResult(
        objectiveSchema,
        result,
        (result: ObjectiveSchemaType) => result
      );
    case "get_auth_token":
    case "close_ledger_channel":
    case "version":
    case "get_address":
    case "close_payment_channel":
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
    case "receive_voucher":
      return validateAndConvertResult(
        receiveVoucherSchema,
        result,
        (result: ReceiveVoucherSchemaType) => ({
          Total: BigInt(result.Total),
          Delta: BigInt(result.Delta),
        })
      );
    case "create_voucher":
      return validateAndConvertResult(
        voucherSchema,
        result,
        (result: VoucherSchemaType) => {
          return {
            ...result,
          };
        }
      );

    default:
      throw new Error(`Unknown method: ${method}`);
  }
}

export function getAndValidateNotification<T extends RPCNotification["method"]>(
  data: unknown,
  method: T
): RPCNotification["params"]["payload"] {
  switch (method) {
    case "payment_channel_updated":
      return convertToInternalPaymentChannelType(
        data as PaymentChannelSchemaType
      );
    case "ledger_channel_updated":
      return convertToInternalPaymentChannelType(
        data as PaymentChannelSchemaType
      );
    case "objective_completed":
      return data as string;
    default:
      throw new Error(`Unknown method: ${method}`);
  }
}
/**
 * Validates that the response is a valid JSON RPC response and pulls out the result
 * @param response - JSON RPC response
 * @returns The result of the response
 */
function getJsonRpcResult(response: unknown): JsonRpcSchemaType {
  const validate = ajv.compile<JsonRpcSchemaType>(jsonRpcSchema);
  if (validate(response)) {
    return response as JsonRpcSchemaType;
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
      TheirBalance: BigInt(result.Balance.TheirBalance),
      MyBalance: BigInt(result.Balance.MyBalance),
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
      PaidSoFar: BigInt(result.Balance.PaidSoFar ?? 0),
      RemainingFunds: BigInt(result.Balance.RemainingFunds ?? 0),
    },
  };
}

function convertToInternalPaymentChannelsType(
  result: PaymentChannelsSchemaType
): PaymentChannelInfo[] {
  return result.map((pc) => convertToInternalPaymentChannelType(pc));
}

/**
 * JSON RPC Types
 */
export type JsonRpcRequest<MethodName extends RequestMethod, RequestParams> = {
  id: number; // in the json-rpc spec this is optional, but we require it for all our requests
  jsonrpc: "2.0";
  method: MethodName;
  params: RequestParams;
};
export type JsonRpcResponse<ResultType> = {
  id: number;
  jsonrpc: "2.0";
  result: ResultType;
};

export type JsonRpcNotification<NotificationName, NotificationParams> = {
  jsonrpc: "2.0";
  method: NotificationName;
  params: NotificationParams;
};

export type JsonRpcError<Code, Message, Data = undefined> = {
  id: number;
  jsonrpc: "2.0";
  error: Data extends undefined
    ? { code: Code; message: Message }
    : { code: Code; message: Message; data: Data };
};

/**
 * Objective params and responses
 */
export type DirectFundParams = {
  CounterParty: string;
  ChallengeDuration: number;
  Outcome: Outcome;
  Nonce: number;
  AppDefinition: string;
  AppData: string;
};
export type VirtualFundParams = {
  Intermediaries: string[];
  CounterParty: string;
  ChallengeDuration: number;
  Outcome: Outcome;
  Nonce: number;
  AppDefinition: string;
};
export type PaymentParams = {
  // todo: this should be a bigint
  Amount: number;
  Channel: string;
};

export type Voucher = {
  ChannelId: string;
  // todo: this should be a bigint
  Amount: number;

  Signature: string;
};
type GetChannelRequest = {
  Id: string;
};

type GetByLedgerRequest = {
  LedgerId: string;
};
export type DefundObjectiveRequest = {
  ChannelId: string;
};
export type ObjectiveResponse = {
  Id: string;
  ChannelId: string;
};
export type ReceiveVoucherResult = {
  Total: bigint;
  Delta: bigint;
};

/**
 * RPC Requests
 */
export type GetAuthTokenRequest = JsonRpcRequest<
  "get_auth_token",
  Record<string, never>
>;
export type GetAddressRequest = JsonRpcRequest<
  "get_address",
  Record<string, never>
>;
export type DirectFundRequest = JsonRpcRequest<
  "create_ledger_channel",
  DirectFundParams
>;
export type PaymentRequest = JsonRpcRequest<"pay", PaymentParams>;
export type VirtualFundRequest = JsonRpcRequest<
  "create_payment_channel",
  VirtualFundParams
>;
export type GetLedgerChannelRequest = JsonRpcRequest<
  "get_ledger_channel",
  GetChannelRequest
>;
export type GetAllLedgerChannelsRequest = JsonRpcRequest<
  "get_all_ledger_channels",
  Record<string, never>
>;
export type GetPaymentChannelRequest = JsonRpcRequest<
  "get_payment_channel",
  GetChannelRequest
>;
export type GetPaymentChannelsByLedgerRequest = JsonRpcRequest<
  "get_payment_channels_by_ledger",
  GetByLedgerRequest
>;

export type VersionRequest = JsonRpcRequest<"version", Record<string, never>>;
export type DirectDefundRequest = JsonRpcRequest<
  "close_ledger_channel",
  DefundObjectiveRequest
>;
export type VirtualDefundRequest = JsonRpcRequest<
  "close_payment_channel",
  DefundObjectiveRequest
>;

export type CreateVoucherRequest = JsonRpcRequest<
  "create_voucher",
  PaymentParams
>;

export type ReceiveVoucherRequest = JsonRpcRequest<"receive_voucher", Voucher>;

/**
 * RPC Responses
 */
export type GetAuthTokenResponse = JsonRpcResponse<string>;
export type GetPaymentChannelResponse = JsonRpcResponse<PaymentChannelInfo>;
export type PaymentResponse = JsonRpcResponse<PaymentParams>;
export type GetLedgerChannelResponse = JsonRpcResponse<LedgerChannelInfo>;
export type VirtualFundResponse = JsonRpcResponse<ObjectiveResponse>;
export type VersionResponse = JsonRpcResponse<string>;
export type GetAddressResponse = JsonRpcResponse<string>;
export type DirectFundResponse = JsonRpcResponse<ObjectiveResponse>;
export type DirectDefundResponse = JsonRpcResponse<string>;
export type VirtualDefundResponse = JsonRpcResponse<string>;
export type GetAllLedgerChannelsResponse = JsonRpcResponse<LedgerChannelInfo[]>;
export type GetPaymentChannelsByLedgerResponse = JsonRpcResponse<
  PaymentChannelInfo[]
>;
export type CreateVoucherResponse = JsonRpcResponse<Voucher>;
export type ReceiveVoucherResponse = JsonRpcResponse<ReceiveVoucherResult>;
/**
 * RPC Request/Response map
 * This is a map of all the RPC methods to their request and response types
 */
export type RPCRequestAndResponses = {
  get_auth_token: [GetAuthTokenRequest, GetAuthTokenResponse];
  create_ledger_channel: [DirectFundRequest, DirectFundResponse];
  close_ledger_channel: [DirectDefundRequest, DirectDefundResponse];
  version: [VersionRequest, VersionResponse];
  create_payment_channel: [VirtualFundRequest, VirtualFundResponse];
  get_address: [GetAddressRequest, GetAddressResponse];
  get_ledger_channel: [GetLedgerChannelRequest, GetLedgerChannelResponse];
  get_payment_channel: [GetPaymentChannelRequest, GetPaymentChannelResponse];
  pay: [PaymentRequest, PaymentResponse];
  close_payment_channel: [VirtualDefundRequest, VirtualDefundResponse];
  get_all_ledger_channels: [
    GetAllLedgerChannelsRequest,
    GetAllLedgerChannelsResponse
  ];
  get_payment_channels_by_ledger: [
    GetPaymentChannelsByLedgerRequest,
    GetPaymentChannelsByLedgerResponse
  ];
  create_voucher: [CreateVoucherRequest, CreateVoucherResponse];
  receive_voucher: [ReceiveVoucherRequest, ReceiveVoucherResponse];
};

export type RequestMethod = keyof RPCRequestAndResponses;

export type RPCRequest =
  RPCRequestAndResponses[keyof RPCRequestAndResponses][0];
export type RPCResponse =
  RPCRequestAndResponses[keyof RPCRequestAndResponses][1];

/**
 * RPC Notifications
 */
export type RPCNotification =
  | ObjectiveCompleteNotification
  | PaymentChannelUpdatedNotification
  | LedgerChannelUpdatedNotification;
export type NotificationMethod = RPCNotification["method"];
export type NotificationPayload = RPCNotification["params"];
export type PaymentChannelUpdatedNotification = JsonRpcNotification<
  "payment_channel_updated",
  PaymentChannelInfo
>;

export type LedgerChannelUpdatedNotification = JsonRpcNotification<
  "ledger_channel_updated",
  LedgerChannelInfo
>;

export type ObjectiveCompleteNotification = JsonRpcNotification<
  "objective_completed",
  string
>;

/**
 * Outcome related types
 */
export type LedgerChannelInfo = {
  ID: string;
  Status: ChannelStatus;
  Balance: LedgerChannelBalance;
};

export type LedgerChannelBalance = {
  AssetAddress: string;
  Them: string;
  Me: string;
  TheirBalance: bigint;
  MyBalance: bigint;
};

export type PaymentChannelBalance = {
  AssetAddress: string;
  Payee: string;
  Payer: string;
  PaidSoFar: bigint;
  RemainingFunds: bigint;
};

export type PaymentChannelInfo = {
  ID: string;
  Status: ChannelStatus;
  Balance: PaymentChannelBalance;
};

export type Outcome = SingleAssetOutcome[];

export type SingleAssetOutcome = {
  Asset: string;
  AssetMetadata: AssetMetadata;
  Allocations: Allocation[];
};

export type Allocation = {
  Destination: string;
  Amount: number;
  AllocationType: number;
  Metadata: null;
};
export type AssetMetadata = {
  AssetType: number;
  Metadata: null;
};

export type ChannelStatus = "Proposed" | "Ready" | "Closing" | "Complete";

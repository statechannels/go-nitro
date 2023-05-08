/**
 * JSON RPC Types
 */
export type JsonRpcRequest<MethodName extends RPCMethod, RequestParams> = {
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
  Amount: number;
  Channel: string;
};
type GetChannelRequest = {
  Id: string;
};
export type DefundObjectiveRequest = {
  ChannelId: string;
};
export type ObjectiveResponse = {
  Id: string;
  ChannelId: string;
};

/**
 * RPC Requests
 */
export type GetAddressRequest = JsonRpcRequest<
  "get_address",
  Record<string, never>
>;
export type DirectFundRequest = JsonRpcRequest<"direct_fund", DirectFundParams>;
export type PaymentRequest = JsonRpcRequest<"pay", PaymentParams>;
export type VirtualFundRequest = JsonRpcRequest<
  "virtual_fund",
  VirtualFundParams
>;
export type GetLedgerChannelRequest = JsonRpcRequest<
  "get_ledger_channel",
  GetChannelRequest
>;
export type GetPaymentChannelRequest = JsonRpcRequest<
  "get_payment_channel",
  GetChannelRequest
>;
export type VersionRequest = JsonRpcRequest<"version", Record<string, never>>;
export type DirectDefundRequest = JsonRpcRequest<
  "direct_defund",
  DefundObjectiveRequest
>;
export type VirtualDefundRequest = JsonRpcRequest<
  "virtual_defund",
  DefundObjectiveRequest
>;

/**
 * RPC Responses
 */
export type GetPaymentChannelResponse = JsonRpcResponse<PaymentChannelInfo>;
export type PaymentResponse = JsonRpcResponse<PaymentParams>;
export type GetLedgerChannelResponse = JsonRpcResponse<LedgerChannelInfo>;
export type VirtualFundResponse = JsonRpcResponse<ObjectiveResponse>;
export type VersionResponse = JsonRpcResponse<string>;
export type GetAddressResponse = JsonRpcResponse<string>;
export type DirectFundResponse = JsonRpcResponse<ObjectiveResponse>;
export type DirectDefundResponse = JsonRpcResponse<string>;
export type VirtualDefundResponse = JsonRpcResponse<string>;

/**
 * RPC Request/Response map
 * This is a map of all the RPC methods to their request and response types
 */
export type RPCRequestAndResponses = {
  direct_fund: [DirectFundRequest, DirectFundResponse];
  direct_defund: [DirectDefundRequest, DirectDefundResponse];
  version: [VersionRequest, VersionResponse];
  virtual_fund: [VirtualFundRequest, VirtualFundResponse];
  get_address: [GetAddressRequest, GetAddressResponse];
  get_ledger_channel: [GetLedgerChannelRequest, GetLedgerChannelResponse];
  get_payment_channel: [GetPaymentChannelRequest, GetPaymentChannelResponse];
  pay: [PaymentRequest, PaymentResponse];
  virtual_defund: [VirtualDefundRequest, VirtualDefundResponse];
};

export type RPCNotification = ObjectiveCompleteNotification;
export type RPCMethod = keyof RPCRequestAndResponses;
export type RPCRequest =
  RPCRequestAndResponses[keyof RPCRequestAndResponses][0];
export type RPCResponse =
  RPCRequestAndResponses[keyof RPCRequestAndResponses][1];

/**
 * RPC Notifications
 */
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
  Hub: string;
  Client: string;
  HubBalance: bigint;
  ClientBalance: bigint;
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

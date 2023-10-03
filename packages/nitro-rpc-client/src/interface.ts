import {
  LedgerChannelInfo,
  ObjectiveResponse,
  PaymentChannelInfo,
  PaymentPayload,
  ReceiveVoucherResult,
  Voucher,
} from "./types";

interface ledgerChannelApi {
  CreateLedgerChannel(
    counterParty: string,
    amount: number
  ): Promise<ObjectiveResponse>;
  CloseLedgerChannel(channelId: string): Promise<string>;
  GetLedgerChannel(channelId: string): Promise<LedgerChannelInfo>;
  GetAllLedgerChannels(): Promise<LedgerChannelInfo[]>;
}
interface paymentChannelApi {
  CreatePaymentChannel(
    counterParty: string,
    intermediaries: string[],
    amount: number
  ): Promise<ObjectiveResponse>;
  ClosePaymentChannel(channelId: string): Promise<string>;
  GetPaymentChannel(channelId: string): Promise<PaymentChannelInfo>;
  GetPaymentChannelsByLedger(ledgerId: string): Promise<PaymentChannelInfo[]>;
}

interface paymentApi {
  CreateVoucher(channelId: string, amount: number): Promise<Voucher>;
  ReceiveVoucher(voucher: Voucher): Promise<ReceiveVoucherResult>;
  Pay(channelId: string, amount: number): Promise<PaymentPayload>;
}

interface syncAPI {
  WaitForObjective(objectiveId: string): Promise<void>;
}

export interface RpcClientApi
  extends ledgerChannelApi,
    paymentChannelApi,
    paymentApi,
    syncAPI {
  GetVersion(): Promise<string>;
  GetAddress(): Promise<string>;
  Close(): Promise<void>;
}

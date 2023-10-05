import {
  LedgerChannelInfo,
  ObjectiveResponse,
  PaymentChannelInfo,
  PaymentPayload,
  ReceiveVoucherResult,
  Voucher,
} from "./types";

interface ledgerChannelApi {
  /**
   * CreateLedgerChannel creates a directly funded ledger channel with the counterparty.
   *
   * @param counterParty - The counterparty to create the channel with
   * @returns A promise that resolves to an objective response, containing the ID of the objective and the channel id.
   */
  CreateLedgerChannel(
    counterParty: string,
    amount: number
  ): Promise<ObjectiveResponse>;
  /**
   * CloseLedgerChannel defunds a directly funded ledger channel.
   *
   * @param channelId - The ID of the channel to defund
   * @returns The ID of the objective that was created
   */
  CloseLedgerChannel(channelId: string): Promise<string>;
  /**
   * GetLedgerChannel queries the RPC server for a payment channel.
   *
   * @param channelId - The ID of the channel to query for
   * @returns A `LedgerChannelInfo` object containing the channel's information
   */
  GetLedgerChannel(channelId: string): Promise<LedgerChannelInfo>;
  /**
   * GetAllLedgerChannels queries the RPC server for all ledger channels.
   * @returns A `LedgerChannelInfo` object containing the channel's information for each ledger channel
   */
  GetAllLedgerChannels(): Promise<LedgerChannelInfo[]>;
}
interface paymentChannelApi {
  /**
   * CreatePaymentChannel creates a virtually funded payment channel with the counterparty, using the given intermediaries.
   *
   * @param counterParty - The counterparty to create the channel with
   * @param intermediaries - The intermerdiaries to use
   * @returns A promise that resolves to an objective response, containing the ID of the objective and the channel id.
   */
  CreatePaymentChannel(
    counterParty: string,
    intermediaries: string[],
    amount: number
  ): Promise<ObjectiveResponse>;
  /**
   * ClosePaymentChannel defunds a virtually funded payment channel.
   *
   * @param channelId - The ID of the channel to defund
   * @returns The ID of the objective that was created
   */
  ClosePaymentChannel(channelId: string): Promise<string>;
  /**
   * GetPaymentChannel queries the RPC server for a payment channel.
   *
   * @param channelId - The ID of the channel to query for
   * @returns A `PaymentChannelInfo` object containing the channel's information
   */
  GetPaymentChannel(channelId: string): Promise<PaymentChannelInfo>;
  /**
   * GetPaymentChannelsByLedger queries the RPC server for any payment channels that are actively funded by the given ledger.
   *
   * @param ledgerId - The ID of the ledger to find payment channels for
   * @returns A `PaymentChannelInfo` object containing the channel's information for each payment channel
   */
  GetPaymentChannelsByLedger(ledgerId: string): Promise<PaymentChannelInfo[]>;
}

interface paymentApi {
  /**
   * Creates a payment voucher for the given channel and amount.
   * The voucher does not get sent to the other party automatically.
   * @param channelId The payment channel to use for the voucher
   * @param amount The amount for the voucher
   * @returns A signed voucher
   */
  CreateVoucher(channelId: string, amount: number): Promise<Voucher>;
  /**
   * Adds a voucher to the go-nitro node that was received from the other party to the channel.
   * @param voucher The voucher to add
   * @returns The total amount of the channel and the delta of the voucher
   */
  ReceiveVoucher(voucher: Voucher): Promise<ReceiveVoucherResult>;
  /**
   * Pay sends a payment on a virtual payment chanel.
   *
   * @param channelId - The ID of the payment channel to use
   * @param amount - The amount to pay
   */
  Pay(channelId: string, amount: number): Promise<PaymentPayload>;
}

interface syncAPI {
  /**
   * WaitForObjective blocks until the objective with the given ID to complete.
   *
   * @param objectiveId - The id objective to wait for
   */
  WaitForObjective(objectiveId: string): Promise<void>;
  /**
   * PaymentChannelUpdated attaches a callback which is triggered when the channel with supplied ID is updated.
   * Returns a cleanup function which can be used to remove the subscription.
   *
   * @param objectiveId - The id objective to wait for
   */
  onPaymentChannelUpdated(
    channelId: string,
    callback: (info: PaymentChannelInfo) => void
  ): () => void;
}

export interface RpcClientApi
  extends ledgerChannelApi,
    paymentChannelApi,
    paymentApi,
    syncAPI {
  /**
   * GetVersion queries the API server for it's version.
   *
   * @returns The version of the RPC server
   */
  GetVersion(): Promise<string>;
  /**
   * GetAddress queries the RPC server for it's state channel address.
   *
   * @returns The address of the wallet connected to the RPC server
   */
  GetAddress(): Promise<string>;
  /**
   * Close closes the RPC client and stops listening for notifications.
   */
  Close(): Promise<void>;
}

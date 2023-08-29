import {
  DefundObjectiveRequest,
  DirectFundParams,
  LedgerChannelInfo,
  PaymentChannelInfo,
  PaymentParams,
  VirtualFundParams,
  RequestMethod,
  RPCRequestAndResponses,
  ObjectiveResponse,
  ObjectiveCompleteNotification,
  Voucher,
  ReceiveVoucherResult,
} from "./types";
import { Transport } from "./transport";
import {
  createOutcome,
  createPaymentChannelOutcome,
  generateRequest,
} from "./utils";
import { HttpTransport } from "./transport/http";
import { getAndValidateResult } from "./serde";

export class NitroRpcClient {
  private transport: Transport;

  // We fetch the address from the RPC server on first use
  private myAddress: string | undefined;

  public get Notifications() {
    return this.transport.Notifications;
  }

  /**
   * Creates a payment voucher for the given channel and amount.
   * The voucher does not get sent to the other party automatically.
   * @param channelId The payment channel to use for the voucher
   * @param amount The amount for the voucher
   * @returns A signed voucher
   */
  public async CreateVoucher(
    channelId: string,
    amount: number
  ): Promise<Voucher> {
    const params = {
      Amount: amount,
      Channel: channelId,
    };
    const request = generateRequest("create_voucher", params);
    const res = await this.transport.sendRequest<"create_voucher">(request);
    return getAndValidateResult(res, "create_voucher");
  }

  /**
   * Adds a voucher to the go-nitro node that was received from the other party to the channel.
   * @param voucher The voucher to add
   * @returns The total amount of the channel and the delta of the voucher
   */
  public async ReceiveVoucher(voucher: Voucher): Promise<ReceiveVoucherResult> {
    const request = generateRequest("receive_voucher", voucher);
    const res = await this.transport.sendRequest<"receive_voucher">(request);
    return getAndValidateResult(res, "receive_voucher");
  }

  /**
   * WaitForObjective blocks until the objective with the given ID to complete.
   *
   * @param objectiveId - The id objective to wait for
   */
  public async WaitForObjective(objectiveId: string): Promise<void> {
    return new Promise((resolve) => {
      this.transport.Notifications.on(
        "objective_completed",
        (params: ObjectiveCompleteNotification["params"]) => {
          if (params === objectiveId) {
            resolve();
          }
        }
      );
    });
  }

  /**
   * CreateLedgerChannel creates a directly funded ledger channel with the counterparty.
   *
   * @param counterParty - The counterparty to create the channel with
   * @returns A promise that resolves to an objective response, containing the ID of the objective and the channel id.
   */
  public async CreateLedgerChannel(
    counterParty: string,
    amount: number
  ): Promise<ObjectiveResponse> {
    const asset = `0x${"00".repeat(20)}`;
    const params: DirectFundParams = {
      CounterParty: counterParty,
      ChallengeDuration: 0,
      Outcome: createOutcome(
        asset,
        await this.GetAddress(),
        counterParty,
        amount
      ),
      AppDefinition: asset,
      AppData: "0x00",
      Nonce: Date.now(),
    };
    return this.sendRequest("create_ledger_channel", params);
  }

  /**
   * CreatePaymentChannel creates a virtually funded payment channel with the counterparty, using the given intermediaries.
   *
   * @param counterParty - The counterparty to create the channel with
   * @param intermediaries - The intermerdiaries to use
   * @returns A promise that resolves to an objective response, containing the ID of the objective and the channel id.
   */
  public async CreatePaymentChannel(
    counterParty: string,
    intermediaries: string[],
    amount: number
  ): Promise<ObjectiveResponse> {
    const asset = `0x${"00".repeat(20)}`;
    const params: VirtualFundParams = {
      CounterParty: counterParty,
      Intermediaries: intermediaries,
      ChallengeDuration: 0,
      Outcome: createPaymentChannelOutcome(
        asset,
        await this.GetAddress(),
        counterParty,
        amount
      ),
      AppDefinition: asset,
      Nonce: Date.now(),
    };

    return this.sendRequest("create_payment_channel", params);
  }

  /**
   * Pay sends a payment on a virtual payment chanel.
   *
   * @param channelId - The ID of the payment channel to use
   * @param amount - The amount to pay
   */
  public async Pay(channelId: string, amount: number): Promise<PaymentParams> {
    const params = {
      Amount: amount,
      Channel: channelId,
    };
    const request = generateRequest("pay", params);
    const res = await this.transport.sendRequest<"pay">(request);
    return getAndValidateResult(res, "pay");
  }

  /**
   * CloseLedgerChannel defunds a directly funded ledger channel.
   *
   * @param channelId - The ID of the channel to defund
   * @returns The ID of the objective that was created
   */
  public async CloseLedgerChannel(channelId: string): Promise<string> {
    const params: DefundObjectiveRequest = { ChannelId: channelId };
    return this.sendRequest("close_ledger_channel", params);
  }
  /**
   * ClosePaymentChannel defunds a virtually funded payment channel.
   *
   * @param channelId - The ID of the channel to defund
   * @returns The ID of the objective that was created
   */

  public async ClosePaymentChannel(channelId: string): Promise<string> {
    const params: DefundObjectiveRequest = { ChannelId: channelId };
    return this.sendRequest("close_payment_channel", params);
  }

  /**
   * GetVersion queries the API server for it's version.
   *
   * @returns The version of the RPC server
   */
  public async GetVersion(): Promise<string> {
    return this.sendRequest("version", {});
  }

  /**
   * GetAddress queries the RPC server for it's state channel address.
   *
   * @returns The address of the wallet connected to the RPC server
   */
  public async GetAddress(): Promise<string> {
    if (this.myAddress) {
      return this.myAddress;
    }

    this.myAddress = await this.sendRequest("get_address", {});
    return this.myAddress;
  }

  /**
   * GetLedgerChannel queries the RPC server for a payment channel.
   *
   * @param channelId - The ID of the channel to query for
   * @returns A `LedgerChannelInfo` object containing the channel's information
   */
  public async GetLedgerChannel(channelId: string): Promise<LedgerChannelInfo> {
    return this.sendRequest("get_ledger_channel", { Id: channelId });
  }

  /**
   * GetAllLedgerChannels queries the RPC server for all ledger channels.
   * @returns A `LedgerChannelInfo` object containing the channel's information for each ledger channel
   */
  public async GetAllLedgerChannels(): Promise<LedgerChannelInfo[]> {
    return this.sendRequest("get_all_ledger_channels", {});
  }

  /**
   * GetPaymentChannel queries the RPC server for a payment channel.
   *
   * @param channelId - The ID of the channel to query for
   * @returns A `PaymentChannelInfo` object containing the channel's information
   */
  public async GetPaymentChannel(
    channelId: string
  ): Promise<PaymentChannelInfo> {
    return this.sendRequest("get_payment_channel", { Id: channelId });
  }
  /**
   * GetPaymentChannelsByLedger queries the RPC server for any payment channels that are actively funded by the given ledger.
   *
   * @param ledgerId - The ID of the ledger to find payment channels for
   * @returns A `PaymentChannelInfo` object containing the channel's information for each payment channel
   */
  public async GetPaymentChannelsByLedger(
    ledgerId: string
  ): Promise<PaymentChannelInfo[]> {
    return this.sendRequest("get_payment_channels_by_ledger", {
      LedgerId: ledgerId,
    });
  }

  async sendRequest<K extends RequestMethod>(
    method: K,
    params: RPCRequestAndResponses[K][0]["params"]
  ): Promise<RPCRequestAndResponses[K][1]["result"]> {
    const request = generateRequest(method, params);
    const res = await this.transport.sendRequest<K>(request);
    return getAndValidateResult(res, method);
  }

  /**
   * Close closes the RPC client and stops listening for notifications.
   */
  public async Close(): Promise<void> {
    return this.transport.Close();
  }

  private constructor(transport: Transport) {
    this.transport = transport;
  }

  /**
   * Creates an RPC client that uses HTTP/WS as the transport.
   *
   * @param url - The URL of the HTTP/WS server
   * @returns A NitroRpcClient that uses WS as the transport
   */
  public static async CreateHttpNitroClient(
    url: string
  ): Promise<NitroRpcClient> {
    const transport = await HttpTransport.createTransport(url);
    return new NitroRpcClient(transport);
  }
}

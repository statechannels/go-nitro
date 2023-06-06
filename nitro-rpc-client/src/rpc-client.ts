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
} from "./types";
import { Transport } from "./transport";
import { createOutcome, generateRequest } from "./utils";
import { HttpTransport } from "./transport/http";
import { validateResponse } from "./serde";

export class NitroRpcClient {
  private transport: Transport;

  // We fetch the address from the RPC server on first use
  private myAddress: string | undefined;

  public get Notifications() {
    return this.transport.Notifications;
  }

  /**
   * WaitForObjective blocks until the objective with the given ID to complete.
   *
   * @param objectiveId - The id objective to wait for
   */
  public async WaitForObjective(objectiveId: string): Promise<void> {
    return new Promise((resolve) => {
      this.transport.Notifications.addListener(
        "objective_completed",
        (notif) => {
          if (notif.params === objectiveId) {
            resolve();
          }
        }
      );
    });
  }

  /**
   * DirectFund creates a directly funded ledger channel with the counterparty.
   *
   * @param counterParty - The counterparty to create the channel with
   * @returns A promise that resolves to an objective response, containing the ID of the objective and the channel id.
   */
  public async DirectFund(counterParty: string): Promise<ObjectiveResponse> {
    const asset = `0x${"00".repeat(20)}`;
    const params: DirectFundParams = {
      CounterParty: counterParty,
      ChallengeDuration: 0,
      Outcome: createOutcome(
        asset,
        await this.GetAddress(),
        counterParty,
        1_000_000
      ),
      AppDefinition: asset,
      AppData: "0x00",
      Nonce: Date.now(),
    };
    return this.sendRequest("direct_fund", params);
  }

  /**
   * VirtualFund creates a virtually funded channel with the counterparty, using the given intermediaries.
   *
   * @param counterParty - The counterparty to create the channel with
   * @param intermediaries - The intermerdiaries to use
   * @returns A promise that resolves to an objective response, containing the ID of the objective and the channel id.
   */
  public async VirtualFund(
    counterParty: string,
    intermediaries: string[]
  ): Promise<ObjectiveResponse> {
    const asset = `0x${"00".repeat(20)}`;
    const params: VirtualFundParams = {
      CounterParty: counterParty,
      Intermediaries: intermediaries,
      ChallengeDuration: 0,
      Outcome: createOutcome(
        asset,
        await this.GetAddress(),
        counterParty,
        1_000
      ),
      AppDefinition: asset,
      Nonce: Date.now(),
    };

    return this.sendRequest("virtual_fund", params);
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
    return res.result;
  }

  /**
   * DirectDefund defunds a directly funded ledger channel.
   *
   * @param channelId - The ID of the channel to defund
   * @returns The ID of the objective that was created
   */
  public async DirectDefund(channelId: string): Promise<string> {
    const params: DefundObjectiveRequest = { ChannelId: channelId };
    return this.sendRequest("direct_defund", params);
  }
  /**
   * VirtualDefund defunds a virtually funded payment channel.
   *
   * @param channelId - The ID of the channel to defund
   * @returns The ID of the objective that was created
   */

  public async VirtualDefund(channelId: string): Promise<string> {
    const params: DefundObjectiveRequest = { ChannelId: channelId };
    return this.sendRequest("virtual_defund", params);
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
    const blacklist = ["pay", "get_all_ledger_channels"];
    if (!blacklist.includes(method)) {
      return validateResponse(res, method);
    }

    return res.result;
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

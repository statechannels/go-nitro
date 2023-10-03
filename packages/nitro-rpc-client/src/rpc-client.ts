import {
  DefundObjectiveRequest,
  DirectFundPayload,
  LedgerChannelInfo,
  PaymentChannelInfo,
  PaymentPayload,
  VirtualFundPayload,
  RequestMethod,
  RPCRequestAndResponses,
  ObjectiveResponse,
  ObjectiveCompleteNotification,
  Voucher,
  ReceiveVoucherResult,
} from "./types";
import { Transport } from "./transport";
import { createOutcome, generateRequest } from "./utils";
import { HttpTransport } from "./transport/http";
import { getAndValidateResult } from "./serde";
import { RpcClientApi } from "./interface";

export class NitroRpcClient implements RpcClientApi {
  private transport: Transport;

  // We fetch the address from the RPC server on first use
  private myAddress: string | undefined;

  private authToken: string | undefined;

  public get Notifications() {
    return this.transport.Notifications;
  }

  public async CreateVoucher(
    channelId: string,
    amount: number
  ): Promise<Voucher> {
    const payload = {
      Amount: amount,
      Channel: channelId,
    };
    const request = generateRequest(
      "create_voucher",
      payload,
      this.authToken || ""
    );
    const res = await this.transport.sendRequest<"create_voucher">(request);
    return getAndValidateResult(res, "create_voucher");
  }

  public async ReceiveVoucher(voucher: Voucher): Promise<ReceiveVoucherResult> {
    const request = generateRequest(
      "receive_voucher",
      voucher,
      this.authToken || ""
    );
    const res = await this.transport.sendRequest<"receive_voucher">(request);
    return getAndValidateResult(res, "receive_voucher");
  }

  public async WaitForObjective(objectiveId: string): Promise<void> {
    return new Promise((resolve) => {
      this.transport.Notifications.on(
        "objective_completed",
        (params: ObjectiveCompleteNotification["params"]) => {
          if (params["payload"] === objectiveId) {
            resolve();
          }
        }
      );
    });
  }

  public async PaymentChannelUpdated(
    channelId: string,
    callback: (info: PaymentChannelInfo) => void
  ): Promise<void> {
    this.transport.Notifications.on(
      "payment_channel_updated",
      (info: PaymentChannelInfo) => {
        if (info.ID == channelId) {
          callback(info);
        }
      }
    );
  }

  public async CreateLedgerChannel(
    counterParty: string,
    amount: number
  ): Promise<ObjectiveResponse> {
    const asset = `0x${"00".repeat(20)}`;
    const payload: DirectFundPayload = {
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
    return this.sendRequest("create_ledger_channel", payload);
  }

  public async CreatePaymentChannel(
    counterParty: string,
    intermediaries: string[],
    amount: number
  ): Promise<ObjectiveResponse> {
    const asset = `0x${"00".repeat(20)}`;
    const payload: VirtualFundPayload = {
      CounterParty: counterParty,
      Intermediaries: intermediaries,
      ChallengeDuration: 0,
      Outcome: createOutcome(
        asset,
        await this.GetAddress(),
        counterParty,
        amount
      ),
      AppDefinition: asset,
      Nonce: Date.now(),
    };

    return this.sendRequest("create_payment_channel", payload);
  }

  public async Pay(channelId: string, amount: number): Promise<PaymentPayload> {
    const payload = {
      Amount: amount,
      Channel: channelId,
    };
    const request = generateRequest("pay", payload, this.authToken || "");
    const res = await this.transport.sendRequest<"pay">(request);
    return getAndValidateResult(res, "pay");
  }

  public async CloseLedgerChannel(channelId: string): Promise<string> {
    const payload: DefundObjectiveRequest = { ChannelId: channelId };
    return this.sendRequest("close_ledger_channel", payload);
  }

  public async ClosePaymentChannel(channelId: string): Promise<string> {
    const payload: DefundObjectiveRequest = { ChannelId: channelId };
    return this.sendRequest("close_payment_channel", payload);
  }

  public async GetVersion(): Promise<string> {
    return this.sendRequest("version", {});
  }

  public async GetAddress(): Promise<string> {
    if (this.myAddress) {
      return this.myAddress;
    }

    this.myAddress = await this.sendRequest("get_address", {});
    return this.myAddress;
  }

  public async GetLedgerChannel(channelId: string): Promise<LedgerChannelInfo> {
    return this.sendRequest("get_ledger_channel", { Id: channelId });
  }

  public async GetAllLedgerChannels(): Promise<LedgerChannelInfo[]> {
    return this.sendRequest("get_all_ledger_channels", {});
  }

  public async GetPaymentChannel(
    channelId: string
  ): Promise<PaymentChannelInfo> {
    return this.sendRequest("get_payment_channel", { Id: channelId });
  }

  public async GetPaymentChannelsByLedger(
    ledgerId: string
  ): Promise<PaymentChannelInfo[]> {
    return this.sendRequest("get_payment_channels_by_ledger", {
      LedgerId: ledgerId,
    });
  }

  private async getAuthToken(): Promise<string> {
    return this.sendRequest("get_auth_token", {});
  }

  private async sendRequest<K extends RequestMethod>(
    method: K,
    payload: RPCRequestAndResponses[K][0]["params"]["payload"]
  ): Promise<RPCRequestAndResponses[K][1]["result"]> {
    const request = generateRequest(method, payload, this.authToken || "");
    const res = await this.transport.sendRequest<K>(request);
    return getAndValidateResult(res, method);
  }

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
    const rpcClient = new NitroRpcClient(transport);
    rpcClient.authToken = await rpcClient.getAuthToken();
    return rpcClient;
  }
}

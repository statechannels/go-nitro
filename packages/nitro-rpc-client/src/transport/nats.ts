import { NatsConnection, connect, JSONCodec, Subscription } from "nats";
import { EventEmitter } from "eventemitter3";

import {
  JsonRpcNotification,
  NotificationMethod,
  NotificationPayload,
  RequestMethod,
  RPCRequestAndResponses,
} from "../types";

import { Transport } from ".";

const NITRO_REQUEST_TOPIC = "nitro-request";
const NITRO_NOTIFICATION_TOPIC = "nitro-notify";

export class NatsTransport {
  private natsConn: NatsConnection;

  private natsSub: Subscription;

  private notifications = new EventEmitter<
    NotificationMethod,
    NotificationPayload
  >();

  public static async createTransport(server: string): Promise<Transport> {
    const natConn = await connect({ servers: server });
    const natsSub = natConn.subscribe(NITRO_NOTIFICATION_TOPIC);
    const transport = new NatsTransport(natConn, natsSub);

    // Start listening for messages without blocking
    transport.listenForMessages(transport.natsSub);

    return transport;
  }

  private constructor(natsConn: NatsConnection, natsSub: Subscription) {
    this.natsConn = natsConn;
    this.natsSub = natsSub;
  }

  public get Notifications(): EventEmitter<
    NotificationMethod,
    NotificationPayload
  > {
    return this.notifications;
  }

  private async listenForMessages(sub: Subscription) {
    for await (const msg of sub) {
      msg.data;
      const notif = JSONCodec().decode(msg.data) as JsonRpcNotification<
        NotificationMethod,
        NotificationPayload
      >;

      switch (notif.method) {
        case "objective_completed":
          this.notifications.emit(notif.method, notif);
          break;
        case "ledger_channel_updated":
          this.notifications.emit(notif.method, notif);
          break;
        case "payment_channel_updated":
          this.notifications.emit(notif.method, notif);
          break;
      }
    }
  }

  public async sendRequest<K extends RequestMethod>(
    req: RPCRequestAndResponses[K][0]
  ): Promise<RPCRequestAndResponses[K][1]> {
    const natsRes = await this.natsConn?.request(
      NITRO_REQUEST_TOPIC,
      JSONCodec().encode(req)
    );

    if (!natsRes) {
      throw new Error("No response");
    }
    const decoded = JSONCodec().decode(natsRes?.data);

    return decoded as RPCRequestAndResponses[K][1];
  }

  public async Close() {
    this.natsSub.unsubscribe();
    await this.natsConn.close();
  }
}

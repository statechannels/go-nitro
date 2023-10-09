import https from "https";

import axios from "axios";
import { w3cwebsocket } from "websocket";
import { EventEmitter } from "eventemitter3";

import {
  NotificationMethod,
  NotificationParams,
  RequestMethod,
  RPCRequestAndResponses,
} from "../types";
import { getAndValidateNotification } from "../serde";

import { Transport } from ".";

export class HttpTransport {
  Notifications: EventEmitter<NotificationMethod, NotificationParams>;

  public static async createTransport(server: string): Promise<Transport> {
    // eslint-disable-next-line new-cap
    const ws = new w3cwebsocket(`wss://${server}/subscribe`);

    // throw any websocket errors so we don't fail silently
    ws.onerror = (e) => {
      console.error("Error with websocket connection to server: " + e);
      throw e;
    };

    // Wait for onopen to fire so we know the connection is ready
    await new Promise<void>((resolve) => (ws.onopen = () => resolve()));

    const transport = new HttpTransport(ws, server);
    return transport;
  }

  public async sendRequest<K extends RequestMethod>(
    req: RPCRequestAndResponses[K][0]
  ): Promise<unknown> {
    const url = new URL(`https://${this.server}`).toString();

    const result = await axios.post(url.toString(), JSON.stringify(req));

    return result.data;
  }

  public async Close(): Promise<void> {
    this.ws.close(1000);
  }

  private ws: w3cwebsocket;

  private server: string;

  private constructor(ws: w3cwebsocket, server: string) {
    this.ws = ws;
    this.server = server;

    this.Notifications = new EventEmitter();
    this.ws.onmessage = (event) => {
      const data = JSON.parse(event.data.toString());
      const validatedResult = getAndValidateNotification(
        data.params.payload,
        data.method
      );

      this.Notifications.emit(data.method, validatedResult);
    };
  }
}

// For testing with self-signed certs, ignore certificate errors. DO NOT use in production.
export function unsecureHttpsAgent(): https.Agent {
  // For testing with self-signed certs, ignore certificate errors. DO NOT use in production.
  const httpsAgent = new https.Agent({
    rejectUnauthorized: false,
  });
  return httpsAgent;
}

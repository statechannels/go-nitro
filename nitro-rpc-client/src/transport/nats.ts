import { NatsConnection, connect, JSONCodec, Subscription } from 'nats';
import { EventEmitter } from 'eventemitter3';
import {
  ObjectiveCompleteNotification,
  RPCMethod,
  RPCRequestAndResponses,
} from '../types';
import { Transport } from '.';

const NITRO_REQUEST_TOPIC = 'nitro-request';
const NITRO_NOTIFICATION_TOPIC = 'nitro-notify';

export class NatsTransport {
  private natsConn: NatsConnection;

  private natsSub: Subscription;

  private objectiveCompleteEmitter = new EventEmitter();

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

  public get Notifications(): EventEmitter {
    return this.objectiveCompleteEmitter;
  }

  private async listenForMessages(sub: Subscription) {
    for await (const msg of sub) {
      const notif = JSONCodec().decode(
        msg.data,
      ) as ObjectiveCompleteNotification;

      this.objectiveCompleteEmitter.emit('objective_completed', notif);
    }
  }

  public async sendRequest<K extends RPCMethod>(
    req: RPCRequestAndResponses[K][0],
  ): Promise<RPCRequestAndResponses[K][1]> {
    const natsRes = await this.natsConn?.request(
      NITRO_REQUEST_TOPIC,
      JSONCodec().encode(req),
    );

    if (!natsRes) {
      throw new Error('No response');
    }
    const decoded = JSONCodec().decode(natsRes?.data);

    return decoded as RPCRequestAndResponses[K][1];
  }

  public async Close() {
    this.natsSub.unsubscribe();
    await this.natsConn.close();
  }
}

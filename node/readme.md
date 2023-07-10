## Nitro node architecture

A nitro node may be instantiated by calling `New()` and passing in a chain service, a messaging service, and a store.

The flow of data through the node is shown in this diagram:

![go-nitro architecture](./go-nitro%20architecture.png)

0. An API can originate as a remote procedure call (RPC) over https/ws or nats. The RPC call is handled by a server which consumes a go-nitro node as a library.
1. The go-nitro `engine` runs in its own goroutine and has a select statement listening for a message from one of:
   - the consuming application (the go-nitro `node` translates an API call into a send on a go channel, and returns a go channel where a response will later be received)
   - the `message` service
   - the `chain` service.
2. The engine reads channels and objectives from the `store`, and computes updates and side effects.
3. The engine gets payment info from the `voucher manager`, and computes updates and side effects.
4. The updates are committed to the `store`.
5. The side effects are sent on go channels to:
   - the `message` service
   - the `chain` service
   - _back_ to the `engine` (e.g. when an update declares further progress can be made)
6. The consuming application is informed about updates (which may return over the RPC connection)

The `chain` and `message` services are responsible for communicating with the blockchain and with counterparties (respectively).

---

To edit the diagram, paste this code into www.sequencediagram.org:

```sequencediagram
title go-nitro Architecture
fontawesome f109 RPC client
participantgroup #azure **go-nitro process**
participantgroup #honeydew **networked components**
fontawesome f233 RPC server
fontawesome f0c1 chain
fontawesome f0e0 msg
end
fontawesome f1e0 node
fontawesome f013 engine
fontawesome f1c0 store
end

alt User Request
RPC client -#red>RPC server: <background:#yellow>JSON-RPC request</background>
RPC server -> node:
node --> engine:<background:#yellow>API Request Triggered
else Blockain Event
chain-->engine: <background:#yellow>Blockchain Event Triggered
else Peer Request
msg-->engine: <background:#yellow>Message Received
end
group handler
engine<->store: <background:#yellow>Read
engine->engine: <background:#yellow>Compute effects
engine<->store: <background:#yellow>Write
engine-->msg: <background:#yellow>Send messages

engine-->chain: <background:#yellow>Submit transactions
engine-->engine: <background:#yellow>Trigger another loop

end

alt Notify User
engine-->node:
node->RPC server:
RPC server-#red>RPC client: <background:#yellow>JSON-RPC response
end
```

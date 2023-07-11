## go-nitro architecture

The flow of data through go-nitro (running as an independent process) is shown in this diagram:

![go-nitro architecture](./go-nitro%20architecture.png)

0. An API call can originate as a remote procedure call (RPC) over https/ws or nats. The RPC call is handled by a server which consumes a go-nitro node as a library.
1. The go-nitro `engine` runs in its own goroutine and listens for messages from one of:
   - the consuming application - the end-user wishes to perform some action
   - the `message` service - a network peer has requested some cooperative action
   - the `chain` service - a channel we are involved with has had a state change on chain 
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

---

An alternate view of the architecture is shown here:

![go-nitro architecture](./go-nitro%20architecture2.png)

[Click here to edit the diagram](https://excalidraw.com/#json=8yRbfXLrLq5sWPp5JD2hQ,l4-oScoyZV6QxGINus0qCQ)

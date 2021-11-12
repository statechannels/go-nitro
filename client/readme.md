## Client architecture

A nitro client may be instantiated by calling `New()` and passing in a chain service and messaging service.

The flow of data through the client is shown in this diagram:

![architecture](architecture.png)

1. The `engine` runs in its own goroutine and has a select statement listening for a message from one of:
   - the consuming application (the go-nitro `client` translates an API call into a send on a go channel, and returns a go channel where a response will later be recieved)
   - the `message` service
   - the `chain` service.
2. The engine reads channels and objectives from the `store`, and computes updates and side effects.
3. The updates arre committed to the `store`.
4. The side effects are sent on go channels to:
   - the `message` service
   - the `chain` service
5. The consuming application is informed about updates

The `chain` and `message` services are responsible for communicating with the blockchain and with counterparties (respectively).

---

[Click here to edit the diagram](https://excalidraw.com/#json=boyMvd14JkaqjD3cRSg0s,xnxQRDKynDM-h28-cWq2mA)

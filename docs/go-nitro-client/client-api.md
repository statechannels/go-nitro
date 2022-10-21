# `go-nitro` Client API

`go-nitro` is a Go library which handles execution of state channels. Applications which want to run on Nitro can import the library, and instantiate a client. Here's an example of how to use it:

```Go
 import nc "github.com/statechannels/go-nitro/client"


 nitroClient := nc.New( // (1)
                    messageservice,
                    chain,
                    storeA,
                    logDestination,
                    nil,
                    nil
                )



response := nitroClient.CreateLedgerChannel(hub.Address, 0, someOutcome)
nitroClient.WaitForCompletedObjective(response.objectiveId) // this is currently only in our test code but we could expose it

// alternatively, we could have nitroClient.waitFor.CreateLedgerChannel, and bundle up the "await" logic in there?

response = nitroClient.CreateVirtualPaymentChannel([hub.Address],bob.Address, defaultChallengeDuration, someOtherOutcome)
nitroClient.WaitForCompletedObjective(response.objectiveId)

for i := 0; i < len(10); i++ {
    clientA.Pay(response.ChannelId, big.NewInt(int64(5)))
}

response = nitroClient.CloseVirtualChannel(response.ChannelId)
nitroClient.WaitForCompletedObjective(response.objectiveId)

// note a lot of this is actually pseudocode, I haven't constructed proper outcome etc.

```

1. This constructor requires several injected dependencies. We are working on a more developer-friendly constructor!

See [here](https://pkg.go.dev/github.com/statechannels/go-nitro@v0.0.0-20221013015616-00c5614be2d2/client#Client) for the godoc description of the `go-nitro.Client` API.

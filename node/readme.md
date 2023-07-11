# `go-nitro` Node

Our [integration tests](./node_test/readme.md) give the best idea of how to use the API. Another useful resource is [the godoc](https://pkg.go.dev/github.com/statechannels/go-nitro/node#Node) description of the `go-nitro.Node` API.

Broadly, consumers will construct a go-nitro `Node`, possibly using injected dependencies. Then, they can create channels and send payments:

```Go
 import nc "github.com/statechannels/go-nitro/node"

 nitroNode := nc.New(
                    messageservice,
                    chain,
                    storeA,
                    logDestination,
                    nil,
                    nil
                )
response := nitroNode.CreateLedgerChannel(hub.Address, 0, someOutcome)
nitroNode.WaitForCompletedObjective(response.objectiveId)

response = nitroNode.CreateVirtualPaymentChannel([hub.Address],bob.Address, defaultChallengeDuration, someOtherOutcome)
nitroNode.WaitForCompletedObjective(response.objectiveId)

for i := 0; i < len(10); i++ {
    nitroNode.Pay(response.ChannelId, big.NewInt(int64(5)))
}

response = nitroNode.CloseVirtualChannel(response.ChannelId)
nitroNode.WaitForCompletedObjective(response.objectiveId)
```

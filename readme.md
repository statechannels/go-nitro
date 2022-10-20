<h1 align="center">
<div><img src="https://statechannels.org/favicon.ico"><br>
go-nitro
</h1>
Implementation of the [Nitro State Channels Framework](https://docs.statechannels.org) in Golang and solidity.

---

## Usage

Go-nitro is work-in-progress library code with an evolving API.

Our [integration tests](./client_test/readme.md) give the best idea of how to use the API. Another useful resource is [the godoc](https://pkg.go.dev/github.com/statechannels/go-nitro@v0.0.0-20221013015616-00c5614be2d2/client#Client) description of the `go-nitro.Client` API (please check for the latest version).

Broadly, consumers will construct a go-nitro `Client`, possibly using injected dependencies. Then, they can create channels and send payments:

```Go
 import nc "github.com/statechannels/go-nitro/client"

 nitroClient := nc.New(
                    messageservice,
                    chain,
                    storeA,
                    logDestination,
                    nil,
                    nil
                )
response := nitroClient.CreateLedgerChannel(hub.Address, 0, someOutcome)
nitroClient.WaitForCompletedObjective(response.objectiveId)

response = nitroClient.CreateVirtualPaymentChannel([hub.Address],bob.Address, defaultChallengeDuration, someOtherOutcome)
nitroClient.WaitForCompletedObjective(response.objectiveId)

for i := 0; i < len(10); i++ {
    clientA.Pay(response.ChannelId, big.NewInt(int64(5)))
}

response = nitroClient.CloseVirtualChannel(response.ChannelId)
nitroClient.WaitForCompletedObjective(response.objectiveId)
```

## Contributing

Please see [contributing.md](./contributing.md)

## ADRs

Architectural decision records may be viewed [here](./.adr/0000-adrs.md).

## Roadmap

The following roadmap gives an idea of the various packages that compose the `go-nitro` module, and their implementation status:

```bash
├── abi ✅                     # types for abi encoding and decoding.
├── channel ✅                 # query the latest supported state of a channel
│   ├── consensus_channel ✅    # manage a running ledger channel.
│   └── state ✅               # generate and recover signatures on state updates
│       ├── outcome ✅         # define how funds are dispersed when a channel closes
├── client 🚧                  # exposes an API to the consuming application
│   └── engine ✅              # coordinate the client components, runs the protocols
│       ├── chainservice 🚧    # watch the chain and submit transactions
│       ├── messageservice ✅  # send and receives messages from peers
│       └── store 🚧           # store keys, state updates and other critical data
├── client_test ✅             # integration tests involving multiple clients
├── crypto  ✅                 # create Ethereum accounts, create & recover signatures
├── internal
│   ├── testactors ✅          # peers with vanity addresses (Alice, Bob, Irene, ... )
│   ├── testdata ✅            # literals and utility functions used by other test packages
│   ├── testhelpers ✅         # pretty-print test failures
|
├── protocols ✅               # functional core of the go-nitro client
│   ├── direct-fund ✅         # fund a channel on-chain
│   ├── direct-defund ✅       # defund a channel on-chain
│   ├── virtual-fund ✅        # fund a channel off-chain through one or more  intermediaries
│   └── virtual-defund ✅      # defund a channel off-chain through one or more intermediaries
└── types ✅                   # basic types and utility methods
```

## On-chain code

The on-chain component of Nitro (i.e. the solidity contracts) are housed in the [`nitro-protocol`](./nitro-protocol/readme.md) directory. This directory contains an npm package with a hardhat / typechain / jest toolchain.

## License

Dual-licensed under [MIT](https://opensource.org/licenses/MIT) + [Apache 2.0](http://www.apache.org/licenses/LICENSE-2.0)

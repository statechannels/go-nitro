<h1 align="center">
<div><img src="https://statechannels.org/favicon.ico"><br>
go-nitro
</h1>

<p align="center">Implementation of the <a href="https://docs.statechannels.org">Nitro State Channels Framework</a> in Golang and Solidity.</p>

## Usage

> ⚠️ Go-nitro is pre-production software ⚠️

### As a Service

Go-nitro can be run as a system service with an RPC api. Go-nitro's default configuration looks for a local blockchain network on port `8545` with chainid `1337`.

A suitably configured node as a docker container is maintained here: https://github.com/statechannels/hardhat-docker, but default hardhat nodes work as well.

After a hardhat node is running, go-nitro can be started from the root directory with

```
go run .
```

Or, built to an executable binary with

```
go build -o gonitro
```

Go nitro accepts the following command flags, which can also be displayed via `go run . -help` (or `gonitro -help` for the build binary).

```
Usage of ./nitro-rpc-server:
  -chainid int
        Specifies the chain id of the chain. (default 1337)
  -chainurl string
        Specifies the url of a RPC endpoint for the chain. (default "ws://127.0.0.1:8545")
  -deploycontracts
        Specifies whether to deploy the adjudicator and create2deployer contracts.
  -msgport int
        Specifies the tcp port for the  message service. (default 3005)
  -naaddress string
        Specifies the address of the nitro adjudicator contract. Default is the address computed by the Create2Deployer contract. (default "0xC6A55E07566416274dBF020b5548eecEdB56290c")
  -pk string
        Specifies the private key used by the node. Default is Alice's private key. (default "2d999770f7b5d49b694080f987b82bbc9fc9ac2b4dcc10b0f8aba7d700f69c6d")
  -rpcport int
        Specifies the tcp port for the rpc server. (default 4005)
  -usedurablestore
        Specifies whether to use a durable store or an in-memory store.
  -usenats
        Specifies whether to use NATS or http/ws for the rpc server.
```

You can make remote procedure calls like so:

```shell
curl -X POST \
  -H 'Content-Type: application/json' \
  -d '{"jsonrpc":"2.0","id":1,"method":"get_address","params":{}}' \
  http://localhost:4005/api/v1
```

but see https://github.com/statechannels/nitro-gui for an RPC client to do so programmatically.

### As a Library

Go-nitro is also work-in-progress library code with an evolving API.

Our [integration tests](./node_test/readme.md) give the best idea of how to use the API. Another useful resource is [the godoc](https://pkg.go.dev/github.com/statechannels/go-nitro@v0.0.0-20221013015616-00c5614be2d2/client#Client) description of the `go-nitro.Node` API (please check for the latest version).

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

### Start RPC servers with Docker

To spin up a docker image with 3 rpc servers and channels pre-populated, run the following:

1. `make docker/build`
1. `make docker/start`

Three rpc go-nitro servers will be available on ports 4005, 4006, and 4007 for Alice, Irene, and Bob. A ledger channel is created between Alice and Irene, and another ledger channel is created between Irene and Bob. A virtual channel is created between Alice and Bob.

### Start RPC servers test script

A [test script](./scripts/start-rpc-servers.go) is available to start up multiple RPC servers and a test chain. This is used to easily and quickly spin up a test environment. The script requires that `foundry` is installed locally; `foundry` installation instructions are available [here](https://book.getfoundry.sh/getting-started/installation).

The script will:

1. Start an `foundry anvil` test chain
2. Deploy the adjudicator contract to the test chain
3. Start an RPC server for Alice (`0xAAA6628Ec44A8a742987EF3A114dDFE2D4F7aDCE`) listening for RPCs on port `4005`
4. Start an RPC server for Irene (`0x111A00868581f73AB42FEEF67D235Ca09ca1E8db`) listening for RPCs on port `4006`
5. Start an RPC server for Bob (`0xBBB676f9cFF8D242e9eaC39D063848807d3D1D94`) listening for RPCs on port `4007`

Stopping the test script will shutdown all RPC servers and `anvil`.

To run the script from the `go-nitro` directory run `go run ./scripts/start-rpc-servers.go`

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
├── crypto  ✅                 # create Ethereum accounts, create & recover signatures
├── node 🚧                    # exposes an API to the consuming application
│   └── engine ✅              # coordinate the node components, runs the protocols
│       ├── chainservice 🚧    # watch the chain and submit transactions
│       ├── messageservice ✅  # send and receives messages from peers
│       └── store 🚧           # store keys, state updates and other critical data
├── node_test ✅               # integration tests involving multiple nodes
├── internal
│   ├── testactors ✅          # peers with vanity addresses (Alice, Bob, Irene, ... )
│   ├── testdata ✅            # literals and utility functions used by other test packages
│   ├── testhelpers ✅         # pretty-print test failures
|
├── protocols ✅               # functional core of the go-nitro node
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

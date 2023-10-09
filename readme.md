<h1 align="center">
<div><img src="https://statechannels.org/favicon.ico"><br>
go-nitro
</h1>

<p align="center">Implementation of the <a href="https://docs.statechannels.org">Nitro State Channels Framework</a> in Golang and Solidity.</p>

`go-nitro` is an implementation of a node in a nitro state channel network. It is software that:

- manages a secret "channel" key
- crafts blockchain transactions (to allow the user to join and exit the network)
- crafts, signs, and sends state channel updates to counterparties in the network
- listens to blockchain events
- listens for counterparty messages
- stores important data to allow for recovery from counterparty inactivity / malice
- understands how to perform these functions safely without risking any funds

## Usage

> ⚠️ Go-nitro is pre-production software ⚠️

Go-nitro can be consumed either as [library code](./node/readme.md) or run as an [independent process](./doc.go) and interfaced with remote procedure calls (recommended).

## Contributing

Please see [contributing.md](./contributing.md)

## ADRs

Architectural decision records may be viewed [here](./.adr/0000-adrs.md).

## Testing

To run unit tests locally, you will need to generate a TLS certificate. Details are [here](./tls/readme.md).

## On-chain code

The on-chain component of Nitro (i.e. the solidity contracts) are housed in the [`nitro-protocol`](./packages/nitro-protocol/readme.md) directory. This directory contains an yarn workspace with a hardhat / typechain / jest toolchain.

## License

Dual-licensed under [MIT](https://opensource.org/licenses/MIT) + [Apache 2.0](http://www.apache.org/licenses/LICENSE-2.0)

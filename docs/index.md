---
description: Ultra-low cost, near-zero-latency conditional transfer of cryptoassets
---

# Nitro State Channel Framework

Nitro channels are a technology which allows for **ultra-low cost, near-zero-latency conditional transfer** of cryptoassets.

They work by allowing users to deposit their funds into a "layer 2" network which sits above a blockchain. This layer 2 network inherits many of the properties of the underlying chain -- for example, **permisionlessness** and the lack of a need to place **trust** in any counterparty.

The network benefits from significant advantages, however -- in a fraction of a second, a direct, private connection can be forged with anyone else in the network. Then, payments can be streamed incredibly fast: up to **hundreds of times per second**. What's more, there is no per-payment fee to cover (just a small fee per unit time for maintaining the connection).

This is in contrast to the experience of transacting publicly on a "layer 1" blockchain or even on a rollup chain. There, fees and latency apply to each and every action.

## What's on offer

These docs cover the open source code at [https://github.com/statechannels/go-nitro](https://github.com/statechannels/go-nitro). This includes

- :fontawesome-brands-npm: [@statechannels/nitro-protocol](https://www.npmjs.com/package/@statechannels/nitro-protocol): smart contracts :simple-solidity: and utilities :simple-typescript:.
- :simple-go: [`go-nitro`](https://github.com/statechannels/go-nitro) fully-featured off-chain client libary. [Inline documentation](https://pkg.go.dev/github.com/statechannels/go-nitro) generated with godoc.

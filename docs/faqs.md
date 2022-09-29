# FAQs

## What are state channels?

State channels are a technology which allows for peer-to-peer games, payments and other few-user applications to safely trade blockchain assets at extremely low latency, low cost and high throughput without requiring trust in a third-party. State channels can be programmed such that assets are redistributed according to arbitary logic, allowing for applications such as poker, conditional payments, atomic swaps and more.

## What is Nitro Protocol?

Nitro protocol is a state of the art state channel protocol which is focussed on security, performance and extensibility. It has been developed over several years of research: the first version was announced in a whitepaper in 2019, and v2 is described in a second paper.

One of the key features of Nitro are _virtual_ channels, where peers can setup a direct connection with each other entirely off-chain.

## How is it implemented?

The on-chain components of Nitro protocol are implemented in solidity, and are published alongside lightweight off-chain support in Typescript in the npm package @statechannels/nitro-protocol.

The off-chain component of the protocol are implemented in `go-nitro`, a reference client for Nitro Protocol v2.

## Where is it being used?

The maintainers of `nitro-protocol` and `go-nitro` are working towards integrating the system into the Filecoin Retrieval Market.

## How can I find out more?

This website covers all the material you need to understand whether Nitro Protocol is a good fit for your use case.

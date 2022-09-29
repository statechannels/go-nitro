# FAQs

## What is a state channel?

A state channel can be thought of as an account with multiple balances (commonly just two). The owners of that account can update those balances according to some rules which they agree on beforehand and which can be enforced on a blockchain.

State channels therefore allow for peer-to-peer games, payments and other few-user applications to safely trade blockchain assets at extremely low latency, low cost and high throughput without requiring trust in a third-party.

## What kind of applications are there?

State channels can be programmed such that assets are redistributed according to arbitary logic, allowing for applications such as poker, conditional payments, atomic swaps and more.

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

## Is it a good fit for my use case?

State channels are not a panacea. If you can answer "yes" to these questions, then they could be a good solution for your application:

- Do you require very low latency transactions?
- Do you require the cost per transaction to be extremely low or zero?
- Do you require some level of privacy?

# 0008 - Data types in Solidity, Go and Typescript

## Status

Accepted

## Context

In a statechannel system, data exchanged and processed off-chain may need to be posted on chain. Therefore it is important to be able to losslessly convert between off-chain and on-chain types.

The Ethereum Virtual Machine uses 256 bit words, and `uint256` (unsigned 256 bit integer) is a very commonly used type. 256 bit integers are generally too large to be represented as primitive types in either Javascript or Golang. We therefore typically use:

JS/TS: a hex `string`, which we manipulate using ethers.js' BigNumber class.
Go: a `*big.Int`.

In either case we may define an alias to these types, and call it something like `Uint256`. This helps us understand which Solidity type we are approximating with the off-chain types.

We also make use of smaller Solidity types, such as `uint48`. There are a couple of reasons why we may have made that choice:

1. Because we do some highly-specialized / custom data packing in order to efficiently store several variables in a single EVM word. Example: `finalizesAt`, which depends on `FixedPart.challengeDuration`
2. Because this as big a type as possible which is still representable by a Javascript `number`. Example: `FixedPart.channelNonce`.

Reason 1 is more fundamental than reason 2. This is truer still now that we have a second (primary) off-chain language in play.

Reason 2 is simply about convenience.

## Poor decisions

When starting out with Golang, we chose to represent `FixedPart.channelNonce` as a `*big.Int`. This had two implications:

- being a non-primitive type, basic arithmetic operations are more verbose to write
- this choice allowed us to represent channelNonces _larger than the maximium possible on-chain channelNonce_ i.e. overflowing the Solidity type we were trying to approximate. We hit upon that issue when we chose to use randomly generated channelNonces.

## Decision 1

`channelNonce` is to have the following type:

Solidity: `uint64` (previously `uint48`)
Go: `uint64` (prevously `*big.Int`)
Typescript: hex `string` (previously `number`)

Implications for Go: This allows us to use the _same_ type in Solidity and Go (a big advantage -- no overflows) as well as having the sweetness of a primitive type in our primary off-chain language.

Implications for Solidity: There are some very small changes in gas consumption. Since most or all of the encoding of this variable results in it being padded to 256 bits, enlarging it does no harm at all. The implications for the semantics of this variable remain essentially unchanged or slightly improved (it allows for a truly huge number of channels).

Implications for Typescript: We have lost the use of a primitive type (unfortunate, since now it is more awkard to do arithmetic). There is an implicit heuristic with the `number` type in Javascript to not try and represent integers larger than 52 bits. Moving to a hex `string` might encourage programmers to represent even larger integers -- but they should always bear in mind that they should not exceed 64 bits because of the solidity type. They are unlikely to do this unless using a random number gererator.

## Decision 2

`challengeDuration` is to have the following type:

Solidity: `uint48` (previously `uint48`)
Go: `uint32` (prevously `*big.Int`)
Typescript: `number` (previously `number`)

Implications for Go: This allows us to use a _safe_ `challengeDuration` in our off-chain code, in the sense that it less than `max(uint48) - <timestamp> ` for any unix timestamp for another few million years. This is important since we pack `finalizesAt = timestamp + challengeDuration` into a slot only 48 bits wide. The Go type system will prevent this overflow happening on chain. A downside is that we aren't necessarily allowing the use of as-large-a-value as possible for this variable.

Implications for Solidity: No change. There is still the chance of an overflow when computing `timestamp + challengeDuration`, so this must be checked off-chain (as ever)!

Implications for Typescript: No change. There is a chance of the same overflow as for Solidity.

## Decision 3

We acknowledge that the type systems of the various runtimes discussed here is limited in the extent to which it can protect from things like overflows or mismatches between runtimes. Future solutions may include runtime checks / validation of variables during serde.

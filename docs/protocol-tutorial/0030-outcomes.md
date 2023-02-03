---
description: State component governing redistribution of cryptoassets.
---

# Outcomes

The **outcome** of a state is the part which dictates where funds are disbursed to once the channel has finalized.

Nitro protocol uses the [L2 exit format](https://github.com/statechannels/exit-format), which is designed as standard for supporting a multitude of token types.

!!! tip

    At the current time, Nitro supports native (e.g. ETH) and ERC20 tokens.

An `Outcome` is an array of `SingleAssetExits`, each specifying:

- an `asset`
- optional `assetMetadata`
- an ordered list of `allocations`

The optional `assetMetadata` is used only by more exotic asset types, and is zero-ed out for native and ERC20 assets.

## Allocations

An `allocation` is

- a `destination` and an `amount`
- an `allocationType` identifier
- optional `metadata`

The `allocationType` identifier is usually set to 0, meaning "simple".

### Simple Allocations

Simple allocations do not have any `metadata`, and allow for funds to be moved on-chain using the `transfer` method. See the section on [defunding](./0080-defunding-a-channel.md).

### Guarantees

When `allocationType` is set to `guarantee`, funds cannot be transferred in the usual way. Instead, they may be moved to another channel on chain using the `reclaim` method. This is explained further in the section on virtual channels. See the section on [defunding](./0080-defunding-a-channel.md).

The `metadata` is an encoding of the following struct:

```solidity
  struct Guarantee {
        bytes32 left;
        bytes32 right;
    }
```

where `left` and `right` correspond to channel participants.

## Destinations

A `destination` is 32 byte identifier which may either denote a [channel ID](./0010-states-channels.md#channel-ids) or a so-called "external destination" (a 20 byte Ethereum address left-padded with zeros).

## A simple example

Putting these elements together, a simple outcome such as "5 ETH to Alice, 5 ETH to Bob" is expressed like this:

=== "Typescript"

    ```ts
    import {
      Exit,
      SingleAssetExit,
      NullAssetMetadata,
    } from "@statechannels/exit-format";

    const ethExit: SingleAssetExit = {
      asset: "0x0", // this implies the native token (e.g. ETH)
      assetMetadata: NullAssetMetadata, // Intentionally left blank
      allocations: [
        {
          destination: "0x00000000000000000000000096f7123E3A80C9813eF50213ADEd0e4511CB820f", // Alice
          amount: "0x05",
          allocationType: AllocationType.simple, // a regular ETH transfer
          metadata: "0x",
        },
        {
          destination: "0x0000000000000000000000000737369d5F8525D039038Da1EdBAC4C4f161b949", // Bob
          amount: "0x05",
          allocationType: AllocationType.simple, // a regular ETH transfer
          metadata: "0x",
        },
      ],
    };

    const exit = [ethExit];
    ```

=== "Go"

    ```Go
      import (
        "math/big"

        "github.com/ethereum/go-ethereum/common"
        "github.com/statechannels/go-nitro/channel/state/outcome"
        "github.com/statechannels/go-nitro/types"
      )

      var ethExit = outcome.SingleAssetExit{
          Allocations: outcome.Allocations{
            outcome.Allocation{
              Destination: types.Destination(common.HexToHash("0x00000000000000000000000096f7123E3A80C9813eF50213ADEd0e4511CB820f")), // Alice
              Amount:      big.NewInt(5),
              // Other fields implicitly zero-ed out
            },
            outcome.Allocation{
              Destination: types.Destination(common.HexToHash("0x0000000000000000000000000737369d5F8525D039038Da1EdBAC4C4f161b949")), // Bob
              Amount:      big.NewInt(5),
              // Other fields implicitly zero-ed out
            },
          },
        }

      var exit = outcome.Exit{{ethExit}}
    ```

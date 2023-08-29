---
description: Unlocking and redistributing money which has been staked.
---

# Defunding a channel

Defunding can only happen after a channel [finalizes](./0040-lifecycle-of-a-channel.md#finalized). Broadly speaking, it is the inverse of [channel funding](./0060-funding-a-channel.md): it can therefore happen on- or off-chain.

## On-chain defunding using `transfer`

If the channel is funded on-chain, it can be defunded using the `transfer` method. It must first be [finalized on chain](./0070-finalizing-a-channel.md), via either a happy or sad path. Because of the information is [stored on chain](./0040-lifecycle-of-a-channel.md#adjudicator-storage), it is necesary to supply both the [`stateHash`](./0010-states-channels.md#state-commitments) and encoded outcome of the channel when calling `transfer`. Furthermore, an asset index is required for slicing the [outcome](./0030-outcomes.md) into a single asset outcome, as well as
a list of indices to "target".

```typescript hl_lines="22 23 24 25 26 27 28"
import { BigNumber, ethers } from "ethers";
import {
  encodeAllocation,
  MAGIC_ADDRESS_INDICATING_ETH,
} from "@statechannels/nitro-protocol";

const amount = "0x03";
const EOA = ethers.Wallet.createRandom().address;
const destination = hexZeroPad(EOA, 32);
const assetOutcome: AllocationAssetOutcome = {
  asset: MAGIC_ADDRESS_INDICATING_ETH,
  allocationItems: [{ destination, amount }],
};
const outcomeBytes = encodeOutcome([
  { asset: MAGIC_ADDRESS_INDICATING_ETH, allocationItems: allocation },
]);

const assetIndex = 0; // (1)
const stateHash = constants.HashZero; // (2)
const indices = []; // (3)

const tx2 = NitroAdjudicator.transfer(
  assetIndex,
  channelId,
  outcomeBytes,
  stateHash,
  indices
);
```

1. This implies we are paying out the 0th asset (in this case the only asset, ETH)
2. If the channel was concluded on the happy path, we can use this default value
3. This magic value (a zero length array) implies we want to pay out all of the allocationItems (in this case there is only one)

Visually, we can see some of the on-chain funding for the channel has been transferred directly to one the channel's participants:

=== "Before"

    ```mermaid
    graph TD;
    linkStyle default interpolate basis;
    ETHAssetHolder( )
    ledger((L))
    me(( )):::me
    hub(( )):::hub
    ETHAssetHolder-->|10|ledger;
    ledger-->|5|me;
    ledger-->|5|hub;
    classDef me fill:#4287f5
    classDef hub fill:#85e69f
    classDef bob fill:#d93434
    ```

=== "After"

    ```mermaid
    graph TD;
    linkStyle default interpolate basis;
    ETHAssetHolder( )
    me(( )):::me
    hub(( )):::hub
    ETHAssetHolder-->|5|me;
    ETHAssetHolder-->|5|hub;
    classDef me fill:#4287f5
    classDef hub fill:#85e69f
    classDef bob fill:#d93434
    ```

!!! info

    There is a convenience method `concludeAndTransferAllAssets` which combines concluding with transferring for every asset --  batching them to save gas.

### Tracking on-chain storage

When a channel has been finalized, and also when a channel has been (partially) defunded using `transfer`, the [on-chain storage](./0040-lifecycle-of-a-channel.md#adjudicator-storage) is updated. To continue the defunding process, it is necessary to track sufficient information to supply the new outcome for the next call to `transfer` (which could be made by a different party).

To do this, it is necessary to listen for `AllocationUpdated` events and to compute the new outcome using `computeTransferEffectsAndInteractions` off-chain helper function

## Off-chain defunding

If the channel in question is funded off chain, it can usually be cooperatively defunded. If that fails, the channel must be finalized on chain and then defunded using the on-chain `transfer` or `reclaim` methods.

### Cooperate

Here, participants in the parent channel make an update which reverses (or reverts) the update they made when [funding the channel](./0060-funding-a-channel.md#fund-virtually). This involves:

1. Removing the allocation which targets the channel in question.
2. Appropriately incrementing the allocations to the parent channel's participants' external destinations.

For step 2, this is very simple if the channel in question was funded with a [simple allocation](./0030-outcomes.md#simple-allocations). Each participant in the parent channel is awarded the funds which were allocated to them in the child channel.

If the channel in question was [virtually funded](./0060-funding-a-channel.md#fund-virtually) with a [guarantee](./0030-outcomes.md#guarantees), each participant in the parent channel is awarded the funds which were allocated to a possibly-distinct participant in the child channel, according to the mapping encoded in the guarantee metadata. The operation should mirror the on-chain `reclaim` method.

### Transfer in, transfer out

If cooperation is not possible, the parent and child channels (let's call them `L` and `X` respectively) must both be finalized on chain. If the child channel is funded with a [simple allocation](./0030-outcomes.md#simple-allocations) like so:

funds may be `transferred` from the parent channel **in** to the child channel. Now the child channel is funded on chain. It can now be defunded as [above](#on-chain-defunding-using-transfer) by transferring money **out** of `X`.

Visually, the following transformation has been applied:

=== "Before"

    ```mermaid
    graph TD;
    linkStyle default interpolate basis;
    ETHAssetHolder( )
    ledger((L))
    channel((X))
    me(( )):::me
    hub(( )):::hub
    ETHAssetHolder-->|10|ledger;
    ledger-->|2|me;
    ledger-->|2|hub;
    ledger-->|6|channel;
    classDef me fill:#4287f5
    classDef hub fill:#85e69f
    classDef bob fill:#d93434
    ```

=== "After"

    ```mermaid
    graph TD;
    linkStyle default interpolate basis;
    ETHAssetHolder( )
    ledger((L))
    channel((X))
    me(( )):::me
    hub(( )):::hub
    ETHAssetHolder-->|4|ledger;
    ledger-->|2|me;
    ledger-->|2|hub;
    ETHAssetHolder-->|6|channel;
    classDef me fill:#4287f5
    classDef hub fill:#85e69f
    classDef bob fill:#d93434
    ```

### Reclaim and transfer out

If cooperation is not possible, the parent and child channels must both be finalized on chain. If the child channel is funded with a [guarantee](./0030-outcomes.md#guarantees), funds may be transferred from the child channel to the parent channel using `reclaim`. Next, the parent channel may be defunded as [above](#on-chain-defunding-using-transfer).

Visually, the following transformation has been applied:

=== "Before"

    ```mermaid
    graph TD;
    linkStyle default interpolate basis;
    ETHAssetHolder( )
    ledger((L))
    channel((X))
    me(( )):::me
    hub(( )):::hub
    ETHAssetHolder-->|10|ledger;
    ledger-->|2|me;
    ledger-->|2|hub;
    ledger-.->|6|channel;
    classDef me fill:#4287f5
    classDef hub fill:#85e69f
    classDef bob fill:#d93434
    ```

=== "After"

    ```mermaid
    graph TD;
    linkStyle default interpolate basis;
    ETHAssetHolder( )
    ledger((L))
    me(( )):::me
    hub(( )):::hub
    ETHAssetHolder-->|10|ledger;
    ledger-->|6|me;
    ledger-->|4|hub;
    classDef me fill:#4287f5
    classDef hub fill:#85e69f
    classDef bob fill:#d93434
    ```

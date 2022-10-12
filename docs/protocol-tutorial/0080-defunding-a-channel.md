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

!!! info

    There is a convenience method `oncludeAndTransferAllAssets` which combines concluding with transferring for every asset --  batching them to save gas.

## Off-chain defunding

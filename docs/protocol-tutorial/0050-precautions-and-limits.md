---
description: What to be aware of!
---

# Precautions and Limits

## Precautions

As a state channel participant, it is advised to check the [`FixedPart`](./0010-states-channels.md#channel-ids) of any channel before joining it. The section on [states](./0010-states-channels.md#states) explains these checks.

## Limits

There are also some limits to be aware of, which apply to the [`VariablePart`](./0010-states-channels.md#states). We describe these limits via some exported constants:

```typescript
import {
  MAX_TX_DATA_SIZE,
  MAX_OUTCOME_ITEMS,
  NITRO_MAX_GAS,
} from "@statechannels/nitro-protocol";
```

| Constant            | Notes                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                          |
| ------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `MAX_TX_DATA_SIZE`  | Reflects the typical effective maximum size for ethereum transaction data. This is set by ethereum clients such as [geth](https://github.com/ethereum/go-ethereum). At the time of writing this is 128KB.                                                                                                                                                                                                                                                                                                                                                                                                                      |
| `NITRO_MAX_GAS`     | An upper limit on the gas consumed by a transaction that we consider "safe" in the sense that it is below the block gas limit on mainnet and most testnets. At the time of writing this constant is set to 6M gas.                                                                                                                                                                                                                                                                                                                                                                                                             |
| `MAX_OUTCOME_ITEMS` | Denotes a safe upper limit on the number of `allocationItems` that may be stored in an [outcome](./0030-outcomes.md). We deem this number safe because the resulting transaction size is less than `MAX_TX_DATA_SIZE` and the transaction consumes less than `NITRO_MAX_GAS` (as confirmed by our test suite). This is for the [`challenge`](./0070-finalizing-a-channel.md#call-challenge) transaction, with the other fields in the state set to modest values (e.g. 2 participants). If those fields grow, `MAX_OUTCOME_ITEMS` may no longer be safe. At the time of writing this constant is set to 2000 allocation items. |

Paying out tokens from a state channel is potentially one of the most expensive operations from a gas perspective (if the recipient does not have any already, the transaction will consume 20K gas per pay out). The same is true of channels paying out (ETH or tokens) to other channels on chain. Bear this in mind when deciding whether to transfer one, many-at-a-time or all-at-once of the tokens from a finalized channel outcome. `NITRO_MAX_GAS / 20000` would be a sensible choice. Remember to leave some headroom for the `transfer` method's intrinsic gas costs: our test suite confirms that at least 100 Token payouts are possible.

TLDR: stick to outcomes withe fewer than `MAX_OUTCOME_ITEMS` entries, and don't try to `transfer` many more than `NITRO_MAX_GAS` / 20000 tokens in one `transfer` transaction.

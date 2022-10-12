# Precautions and Limits

TODO fix links in this file.

## Precautions

As a state channel participant, it is advised to check the [`FixedPart`](./0001-states-channels-execution-rules.md#channel-ids) of any channel before participating in it. A good state channels wallet will perform these checks for you:

| Field               | Precautions                                                                                                                                                                                                                                                                                                                                                 |
| ------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `chainId`           | This needs to match the id of the chain where assets are to be locked (i.e. the 'root' of the funding graph for this channel). In the event of a mismatch, the channel cannot be concluded and funds cannot be unlocked.                                                                                                                                    |
| `participants`      | This should have length at least 2, but no more than 255, and include a public key (account) that you control. Each entry should be a nonzero ethereum address.                                                                                                                                                                                             |
| `channelNonce`      | This should be different to any previous `channelNonce` used by the same `participants` and `chainId`. This is to prevent states from previous channels being "replayed" to conclude subsequent channels with unintended outcomes. Must be less than `2**48-1`.                                                                                             |
| `appDefinition`     | There should be an [`IForceMoveApp`](../../contract-api/natspec/IForceMoveApp) compliant contract deployed at this address, and you should have confidence that it is not malicious or suffering from security flaws. You should inspect the source code (which should be publically available and verifiable) or appeal to a trusted authority to do this. |
| `challengeDuration` | In the extreme, this should be at least 1 block time (15 seconds on mainnet) and less than `2**48-1` seconds. Whatever it is set to, the channel should be closed long before `2**48 - 1 - challengeDuration`. In practice we recommend somewhere between 5 minutes and 5 months.                                                                           |

## Limits

There are also some limits to be aware of, which apply to the [`VariablePart`](./0001-states-channels-execution-rules.md#channel-ids). We describe these limits via some exported constants:

```typescript
import {
  MAX_TX_DATA_SIZE,
  MAX_OUTCOME_ITEMS,
  NITRO_MAX_GAS,
} from "@statechannels/nitro-protocol";
```

| Constant            | Notes                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                          |
| ------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `MAX_TX_DATA_SIZE`  | Reflects the typical effective maximum size for ethereum transaction data. This is set by ethereum clients such as [geth](https://github.com/ethereum/go-ethereum). At the time of writing this is 128KB.                                                                                                                                                                                                                                                                                                                                                                                                                                      |
| `NITRO_MAX_GAS`     | An upper limit on the gas consumed by a transaction that we consider "safe" in the sense that it is below the block gas limit on mainnet and most testnets. At the time of writing this constant is set to 6M gas.                                                                                                                                                                                                                                                                                                                                                                                                                             |
| `MAX_OUTCOME_ITEMS` | Denotes a safe upper limit on the number of `allocationItems` that may be stored in an [outcome](./outcomes#outcomes-that-allocate). We deem this number safe because the resulting transaction size is less than `MAX_TX_DATA_SIZE` and the transaction consumes less than `NITRO_MAX_GAS` (as confirmed by our test suite). This is for the [`challenge`](../../contract-api/natspec/ForceMove#challenge) transaction, with the other fields in the state set to modest values (e.g. 2 participants). If those fields grow, `MAX_OUTCOME_ITEMS` may no longer be safe. At the time of writing this constant is set to 2000 allocation items. |

Paying out tokens from a state channel is potentially one of the most expensive operations from a gas perspective (if the recipient does not have any already, the transaction will consume 20K gas per pay out). The same is true of channels paying out (ETH or tokens) to other channels on chain. Bear this in mind when deciding whether to transfer one, many-at-a-time or all-at-once of the tokens from a finalized channel outcome. `NITRO_MAX_GAS / 20000` would be a sensible choice. Remember to leave some headroom for the `transfer` method's intrinsic gas costs: our test suite confirms that at least 100 Token payouts are possible.

TLDR: stick to outcomes withe fewer than `MAX_OUTCOME_ITEMS` entries, and don't try to `transfer` many more than `NITRO_MAX_GAS` / 20000 tokens in one `transfer` transaction.

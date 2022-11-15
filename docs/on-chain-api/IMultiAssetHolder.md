# Solidity API

## IMultiAssetHolder

_The IMultiAssetHolder interface calls for functions that allow assets to be transferred from one channel to other channel and/or external destinations, as well as for guarantees to be claimed._

### deposit

```solidity
function deposit(address asset, bytes32 destination, uint256 expectedHeld, uint256 amount) external payable
```

Deposit ETH or erc20 assets against a given destination.

_Deposit ETH or erc20 assets against a given destination._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| asset | address | erc20 token address, or zero address to indicate ETH |
| destination | bytes32 | ChannelId to be credited. |
| expectedHeld | uint256 | The number of wei the depositor believes are _already_ escrowed against the channelId. |
| amount | uint256 | The intended number of wei to be deposited. |

### transfer

```solidity
function transfer(uint256 assetIndex, bytes32 fromChannelId, bytes outcomeBytes, bytes32 stateHash, uint256[] indices) external
```

Transfers as many funds escrowed against `channelId` as can be afforded for a specific destination. Assumes no repeated entries.

_Transfers as many funds escrowed against `channelId` as can be afforded for a specific destination. Assumes no repeated entries._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| assetIndex | uint256 | Will be used to slice the outcome into a single asset outcome. |
| fromChannelId | bytes32 | Unique identifier for state channel to transfer funds *from*. |
| outcomeBytes | bytes | The encoded Outcome of this state channel |
| stateHash | bytes32 | The hash of the state stored when the channel finalized. |
| indices | uint256[] | Array with each entry denoting the index of a destination to transfer funds to. An empty array indicates "all". |

### ReclaimArgs

```solidity
struct ReclaimArgs {
  bytes32 sourceChannelId;
  bytes32 sourceStateHash;
  bytes sourceOutcomeBytes;
  uint256 sourceAssetIndex;
  uint256 indexOfTargetInSource;
  bytes32 targetStateHash;
  bytes targetOutcomeBytes;
  uint256 targetAssetIndex;
}
```

### reclaim

```solidity
function reclaim(struct IMultiAssetHolder.ReclaimArgs reclaimArgs) external
```

Reclaim moves money from a target channel back into a ledger channel which is guaranteeing it. The guarantee is removed from the ledger channel.

_Reclaim moves money from a target channel back into a ledger channel which is guaranteeing it. The guarantee is removed from the ledger channel._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| reclaimArgs | struct IMultiAssetHolder.ReclaimArgs | arguments used in the claim function. Used to avoid stack too deep error. |

### Deposited

```solidity
event Deposited(bytes32 destination, address asset, uint256 amountDeposited, uint256 destinationHoldings)
```

_Indicates that `amountDeposited` has been deposited into `destination`._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| destination | bytes32 | The channel being deposited into. |
| asset | address |  |
| amountDeposited | uint256 | The amount being deposited. |
| destinationHoldings | uint256 | The new holdings for `destination`. |

### AllocationUpdated

```solidity
event AllocationUpdated(bytes32 channelId, uint256 assetIndex, uint256 initialHoldings)
```

_Indicates the assetOutcome for this channelId and assetIndex has changed due to a transfer. Includes sufficient data to compute:
- the new assetOutcome
- the new holdings for this channelId and any others that were transferred to
- the payouts to external destinations
when combined with the calldata of the transaction causing this event to be emitted._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| channelId | bytes32 | The channelId of the funds being withdrawn. |
| assetIndex | uint256 |  |
| initialHoldings | uint256 | holdings[asset][channelId] **before** the allocations were updated. The asset in question can be inferred from the calldata of the transaction (it might be "all assets") |

### Reclaimed

```solidity
event Reclaimed(bytes32 channelId, uint256 assetIndex)
```

_Indicates the assetOutcome for this channelId and assetIndex has changed due to a reclaim. Includes sufficient data to compute:
- the new assetOutcome
when combined with the calldata of the transaction causing this event to be emitted._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| channelId | bytes32 | The channelId of the funds being withdrawn. |
| assetIndex | uint256 |  |


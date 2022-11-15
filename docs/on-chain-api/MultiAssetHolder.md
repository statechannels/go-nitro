# Solidity API

## MultiAssetHolder

_An implementation of the IMultiAssetHolder interface. The AssetHolder contract escrows ETH or tokens against state channels. It allows assets to be internally accounted for, and ultimately prepared for transfer from one channel to other channels and/or external destinations, as well as for guarantees to be reclaimed._

### holdings

```solidity
mapping(address => mapping(bytes32 => uint256)) holdings
```

holdings[asset][channelId] is the amount of asset held against channel channelId. 0 address implies ETH

### deposit

```solidity
function deposit(address asset, bytes32 channelId, uint256 expectedHeld, uint256 amount) external payable
```

Deposit ETH or erc20 tokens against a given channelId.

_Deposit ETH or erc20 tokens against a given channelId._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| asset | address | erc20 token address, or zero address to indicate ETH |
| channelId | bytes32 | ChannelId to be credited. |
| expectedHeld | uint256 | The number of wei/tokens the depositor believes are _already_ escrowed against the channelId. |
| amount | uint256 | The intended number of wei/tokens to be deposited. |

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

### _apply_transfer_checks

```solidity
function _apply_transfer_checks(uint256 assetIndex, uint256[] indices, bytes32 channelId, bytes32 stateHash, bytes outcomeBytes) internal view returns (struct ExitFormat.SingleAssetExit[] outcome, address asset, uint256 initialAssetHoldings)
```

### compute_transfer_effects_and_interactions

```solidity
function compute_transfer_effects_and_interactions(uint256 initialHoldings, struct ExitFormat.Allocation[] allocations, uint256[] indices) public pure returns (struct ExitFormat.Allocation[] newAllocations, bool allocatesOnlyZeros, struct ExitFormat.Allocation[] exitAllocations, uint256 totalPayouts)
```

### _apply_transfer_effects

```solidity
function _apply_transfer_effects(uint256 assetIndex, address asset, bytes32 channelId, bytes32 stateHash, struct ExitFormat.SingleAssetExit[] outcome, struct ExitFormat.Allocation[] newAllocations, uint256 initialHoldings, uint256 totalPayouts) internal
```

### _apply_transfer_interactions

```solidity
function _apply_transfer_interactions(struct ExitFormat.SingleAssetExit singleAssetExit, struct ExitFormat.Allocation[] exitAllocations) internal
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
| reclaimArgs | struct IMultiAssetHolder.ReclaimArgs | arguments used in the reclaim function. Used to avoid stack too deep error. |

### _apply_reclaim_checks

```solidity
function _apply_reclaim_checks(struct IMultiAssetHolder.ReclaimArgs reclaimArgs) internal view returns (struct ExitFormat.SingleAssetExit[] sourceOutcome, struct ExitFormat.SingleAssetExit[] targetOutcome)
```

_Checks that the source and target channels are finalized; that the supplied outcomes match the stored fingerprints; that the asset is identical in source and target. Computes and returns the decoded outcomes._

### compute_reclaim_effects

```solidity
function compute_reclaim_effects(struct ExitFormat.Allocation[] sourceAllocations, struct ExitFormat.Allocation[] targetAllocations, uint256 indexOfTargetInSource) public pure returns (struct ExitFormat.Allocation[])
```

_Computes side effects for the reclaim function. Returns updated allocations for the source, computed by finding the guarantee in the source for the target, and moving money out of the guarantee and back into the ledger channel as regular allocations for the participants._

### _apply_reclaim_effects

```solidity
function _apply_reclaim_effects(struct IMultiAssetHolder.ReclaimArgs reclaimArgs, struct ExitFormat.SingleAssetExit[] sourceOutcome, struct ExitFormat.Allocation[] newSourceAllocations) internal
```

_Updates the fingerprint of the outcome for the source channel and emit an event for it._

### _executeSingleAssetExit

```solidity
function _executeSingleAssetExit(struct ExitFormat.SingleAssetExit singleAssetExit) internal
```

Executes a single asset exit by paying out the asset and calling external contracts, as well as updating the holdings stored in this contract.

_Executes a single asset exit by paying out the asset and calling external contracts, as well as updating the holdings stored in this contract._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| singleAssetExit | struct ExitFormat.SingleAssetExit | The single asset exit to be paid out. |

### _transferAsset

```solidity
function _transferAsset(address asset, address destination, uint256 amount) internal
```

Transfers the given amount of this AssetHolders's asset type to a supplied ethereum address.

_Transfers the given amount of this AssetHolders's asset type to a supplied ethereum address._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| asset | address |  |
| destination | address | ethereum address to be credited. |
| amount | uint256 | Quantity of assets to be transferred. |

### _isExternalDestination

```solidity
function _isExternalDestination(bytes32 destination) internal pure returns (bool)
```

Checks if a given destination is external (and can therefore have assets transferred to it) or not.

_Checks if a given destination is external (and can therefore have assets transferred to it) or not._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| destination | bytes32 | Destination to be checked. |

#### Return Values

| Name | Type | Description |
| ---- | ---- | ----------- |
| [0] | bool | True if the destination is external, false otherwise. |

### _addressToBytes32

```solidity
function _addressToBytes32(address participant) internal pure returns (bytes32)
```

Converts an ethereum address to a nitro external destination.

_Converts an ethereum address to a nitro external destination._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| participant | address | The address to be converted. |

#### Return Values

| Name | Type | Description |
| ---- | ---- | ----------- |
| [0] | bytes32 | The input address left-padded with zeros. |

### _bytes32ToAddress

```solidity
function _bytes32ToAddress(bytes32 destination) internal pure returns (address payable)
```

Converts a nitro destination to an ethereum address.

_Converts a nitro destination to an ethereum address._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| destination | bytes32 | The destination to be converted. |

#### Return Values

| Name | Type | Description |
| ---- | ---- | ----------- |
| [0] | address payable | The rightmost 160 bits of the input string. |

### _requireMatchingFingerprint

```solidity
function _requireMatchingFingerprint(bytes32 stateHash, bytes32 outcomeHash, bytes32 channelId) internal view
```

Checks that a given variables hash to the data stored on chain.

_Checks that a given variables hash to the data stored on chain._

### _requireChannelFinalized

```solidity
function _requireChannelFinalized(bytes32 channelId) internal view
```

Checks that a given channel is in the Finalized mode.

_Checks that a given channel is in the Finalized mode._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| channelId | bytes32 | Unique identifier for a channel. |

### _updateFingerprint

```solidity
function _updateFingerprint(bytes32 channelId, bytes32 stateHash, bytes32 outcomeHash) internal
```

### _requireIncreasingIndices

```solidity
function _requireIncreasingIndices(uint256[] indices) internal pure
```

Checks that the supplied indices are strictly increasing.

_Checks that the supplied indices are strictly increasing. This allows us allows us to write a more efficient claim function._

### min

```solidity
function min(uint256 a, uint256 b) internal pure returns (uint256)
```

### Guarantee

```solidity
struct Guarantee {
  bytes32 left;
  bytes32 right;
}
```

### decodeGuaranteeData

```solidity
function decodeGuaranteeData(bytes data) internal pure returns (struct MultiAssetHolder.Guarantee)
```


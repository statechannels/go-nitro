# Solidity API

## StatusManager

_The StatusManager is responsible for on-chain storage of the status of active channels_

### statusOf

```solidity
mapping(bytes32 => bytes32) statusOf
```

### _mode

```solidity
function _mode(bytes32 channelId) internal view returns (enum IStatusManager.ChannelMode)
```

Computes the ChannelMode for a given channelId.

_Computes the ChannelMode for a given channelId._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| channelId | bytes32 | Unique identifier for a channel. |

### _generateStatus

```solidity
function _generateStatus(struct IStatusManager.ChannelData channelData) internal pure returns (bytes32 status)
```

Formats the input data for on chain storage.

_Formats the input data for on chain storage._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| channelData | struct IStatusManager.ChannelData | ChannelData data. |

### _generateFingerprint

```solidity
function _generateFingerprint(bytes32 stateHash, bytes32 outcomeHash) internal pure returns (uint160)
```

### _unpackStatus

```solidity
function _unpackStatus(bytes32 channelId) internal view returns (uint48 turnNumRecord, uint48 finalizesAt, uint160 fingerprint)
```

Unpacks turnNumRecord, finalizesAt and fingerprint from the status of a particular channel.

_Unpacks turnNumRecord, finalizesAt and fingerprint from the status of a particular channel._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| channelId | bytes32 | Unique identifier for a state channel. |

#### Return Values

| Name | Type | Description |
| ---- | ---- | ----------- |
| turnNumRecord | uint48 | A turnNum that (the adjudicator knows) is supported by a signature from each participant. |
| finalizesAt | uint48 | The unix timestamp when `channelId` will finalize. |
| fingerprint | uint160 | The last 160 bits of kecca256(stateHash, outcomeHash) |


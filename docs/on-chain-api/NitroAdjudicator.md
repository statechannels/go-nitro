# Solidity API

## NitroAdjudicator

_The NitroAdjudicator contract extends MultiAssetHolder and ForceMove_

### concludeAndTransferAllAssets

```solidity
function concludeAndTransferAllAssets(struct INitroTypes.FixedPart fixedPart, struct INitroTypes.SignedVariablePart[] proof, struct INitroTypes.SignedVariablePart candidate) public
```

Finalizes a channel by providing a finalization proof, and liquidates all assets for the channel.

_Finalizes a channel by providing a finalization proof, and liquidates all assets for the channel._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| fixedPart | struct INitroTypes.FixedPart | Data describing properties of the state channel that do not change with state updates. |
| proof | struct INitroTypes.SignedVariablePart[] | Variable parts of the states with signatures in the support proof. The proof is a validation for the supplied candidate. |
| candidate | struct INitroTypes.SignedVariablePart | Variable part of the state to change to. The candidate state is supported by proof states. |

### transferAllAssets

```solidity
function transferAllAssets(bytes32 channelId, struct ExitFormat.SingleAssetExit[] outcome, bytes32 stateHash) public
```

Liquidates all assets for the channel

_Liquidates all assets for the channel_

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| channelId | bytes32 | Unique identifier for a state channel |
| outcome | struct ExitFormat.SingleAssetExit[] | An array of SingleAssetExit[] items. |
| stateHash | bytes32 | stored state hash for the channel |

### requireStateSupported

```solidity
function requireStateSupported(struct INitroTypes.FixedPart fixedPart, struct INitroTypes.SignedVariablePart[] proof, struct INitroTypes.SignedVariablePart candidate) external pure
```

Encodes application-specific rules for a particular ForceMove-compliant state channel.

_Encodes application-specific rules for a particular ForceMove-compliant state channel._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| fixedPart | struct INitroTypes.FixedPart | Fixed part of the state channel. |
| proof | struct INitroTypes.SignedVariablePart[] | Variable parts of the states with signatures in the support proof. The proof is a validation for the supplied candidate. |
| candidate | struct INitroTypes.SignedVariablePart | Variable part of the state to change to. The candidate state is supported by proof states. |

### _executeExit

```solidity
function _executeExit(struct ExitFormat.SingleAssetExit[] exit) internal
```

Executes an exit by paying out assets and calling external contracts

_Executes an exit by paying out assets and calling external contracts_

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| exit | struct ExitFormat.SingleAssetExit[] | The exit to be paid out. |


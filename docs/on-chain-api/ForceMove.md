# Solidity API

## ForceMove

_An implementation of ForceMove protocol, which allows state channels to be adjudicated and finalized._

### unpackStatus

```solidity
function unpackStatus(bytes32 channelId) external view returns (uint48 turnNumRecord, uint48 finalizesAt, uint160 fingerprint)
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

### challenge

```solidity
function challenge(struct INitroTypes.FixedPart fixedPart, struct INitroTypes.SignedVariablePart[] proof, struct INitroTypes.SignedVariablePart candidate, struct INitroTypes.Signature challengerSig) external
```

Registers a challenge against a state channel. A challenge will either prompt another participant into clearing the challenge (via one of the other methods), or cause the channel to finalize at a specific time.

_Registers a challenge against a state channel. A challenge will either prompt another participant into clearing the challenge (via one of the other methods), or cause the channel to finalize at a specific time._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| fixedPart | struct INitroTypes.FixedPart | Data describing properties of the state channel that do not change with state updates. |
| proof | struct INitroTypes.SignedVariablePart[] | An ordered array of structs, that can be signed by any number of participants, each struct describing the properties of the state channel that may change with each state update. The proof is a validation for the supplied candidate. |
| candidate | struct INitroTypes.SignedVariablePart | A struct, that can be signed by any number of participants, describing the properties of the state channel to change to. The candidate state is supported by proof states. |
| challengerSig | struct INitroTypes.Signature | The signature of a participant on the keccak256 of the abi.encode of (supportedStateHash, 'forceMove'). |

### checkpoint

```solidity
function checkpoint(struct INitroTypes.FixedPart fixedPart, struct INitroTypes.SignedVariablePart[] proof, struct INitroTypes.SignedVariablePart candidate) external
```

Overwrites the `turnNumRecord` stored against a channel by providing a proof with higher turn number.

_Overwrites the `turnNumRecord` stored against a channel by providing a proof with higher turn number._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| fixedPart | struct INitroTypes.FixedPart | Data describing properties of the state channel that do not change with state updates. |
| proof | struct INitroTypes.SignedVariablePart[] | An ordered array of structs, that can be signed by any number of participants, each struct describing the properties of the state channel that may change with each state update. The proof is a validation for the supplied candidate. |
| candidate | struct INitroTypes.SignedVariablePart | A struct, that can be signed by any number of participants, describing the properties of the state channel to change to. The candidate state is supported by proof states. |

### conclude

```solidity
function conclude(struct INitroTypes.FixedPart fixedPart, struct INitroTypes.SignedVariablePart[] proof, struct INitroTypes.SignedVariablePart candidate) external
```

Finalizes a channel by providing a finalization proof. External wrapper for _conclude.

_Finalizes a channel by providing a finalization proof. External wrapper for _conclude._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| fixedPart | struct INitroTypes.FixedPart | Data describing properties of the state channel that do not change with state updates. |
| proof | struct INitroTypes.SignedVariablePart[] | An ordered array of structs, that can be signed by any number of participants, each struct describing the properties of the state channel that may change with each state update. The proof is a validation for the supplied candidate. |
| candidate | struct INitroTypes.SignedVariablePart | A struct, that can be signed by any number of participants, describing the properties of the state channel to change to. The candidate state is supported by proof states. |

### _conclude

```solidity
function _conclude(struct INitroTypes.FixedPart fixedPart, struct INitroTypes.SignedVariablePart[] proof, struct INitroTypes.SignedVariablePart candidate) internal returns (bytes32 channelId)
```

Finalizes a channel by providing a finalization proof. Internal method.

_Finalizes a channel by providing a finalization proof. Internal method._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| fixedPart | struct INitroTypes.FixedPart | Data describing properties of the state channel that do not change with state updates. |
| proof | struct INitroTypes.SignedVariablePart[] | An ordered array of structs, that can be signed by any number of participants, each struct describing the properties of the state channel that may change with each state update. The proof is a validation for the supplied candidate. |
| candidate | struct INitroTypes.SignedVariablePart | A struct, that can be signed by any number of participants, describing the properties of the state channel to change to. The candidate state is supported by proof states. |

### getChainID

```solidity
function getChainID() public view returns (uint256)
```

### _requireCorrectChainId

```solidity
function _requireCorrectChainId(uint256 declaredChainId) internal view
```

Checks that the supplied chain id matches the chain id of this contract.

_Checks that the supplied chain id matches the chain id of this contract._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| declaredChainId | uint256 | The chain id to check. |

### _requireChallengerIsParticipant

```solidity
function _requireChallengerIsParticipant(bytes32 supportedStateHash, address[] participants, struct INitroTypes.Signature challengerSignature) internal pure
```

Checks that the challengerSignature was created by one of the supplied participants.

_Checks that the challengerSignature was created by one of the supplied participants._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| supportedStateHash | bytes32 | Forms part of the digest to be signed, along with the string 'forceMove'. |
| participants | address[] | A list of addresses representing the participants of a channel. |
| challengerSignature | struct INitroTypes.Signature | The signature of a participant on the keccak256 of the abi.encode of (supportedStateHash, 'forceMove'). |

### _isAddressInArray

```solidity
function _isAddressInArray(address suspect, address[] addresses) internal pure returns (bool)
```

Tests whether a given address is in a given array of addresses.

_Tests whether a given address is in a given array of addresses._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| suspect | address | A single address of interest. |
| addresses | address[] | A line-up of possible perpetrators. |

#### Return Values

| Name | Type | Description |
| ---- | ---- | ----------- |
| [0] | bool | true if the address is in the array, false otherwise |

### _requireStateSupported

```solidity
function _requireStateSupported(struct INitroTypes.FixedPart fixedPart, struct INitroTypes.SignedVariablePart[] proof, struct INitroTypes.SignedVariablePart candidate) internal pure
```

Check that the submitted data constitute a support proof, revert if not.

_Check that the submitted data constitute a support proof, revert if not._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| fixedPart | struct INitroTypes.FixedPart | Fixed Part of the states in the support proof. |
| proof | struct INitroTypes.SignedVariablePart[] | Variable parts of the states with signatures in the support proof. The proof is a validation for the supplied candidate. |
| candidate | struct INitroTypes.SignedVariablePart | Variable part of the state to change to. The candidate state is supported by proof states. |

### recoverVariableParts

```solidity
function recoverVariableParts(struct INitroTypes.FixedPart fixedPart, struct INitroTypes.SignedVariablePart[] signedVariableParts) internal pure returns (struct INitroTypes.RecoveredVariablePart[])
```

Recover signatures for each variable part in the supplied array.

_Recover signatures for each variable part in the supplied array._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| fixedPart | struct INitroTypes.FixedPart | Fixed Part of the states in the support proof. |
| signedVariableParts | struct INitroTypes.SignedVariablePart[] | Signed variable parts of the states in the support proof. |

#### Return Values

| Name | Type | Description |
| ---- | ---- | ----------- |
| [0] | struct INitroTypes.RecoveredVariablePart[] | An array of recoveredVariableParts, identical to the supplied signedVariableParts array, but with the signatures replaced with a signedBy bitmask. |

### recoverVariablePart

```solidity
function recoverVariablePart(struct INitroTypes.FixedPart fixedPart, struct INitroTypes.SignedVariablePart signedVariablePart) internal pure returns (struct INitroTypes.RecoveredVariablePart)
```

Recover signatures for a variable part.

_Recover signatures for a variable part._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| fixedPart | struct INitroTypes.FixedPart | Fixed Part of the states in the support proof. |
| signedVariablePart | struct INitroTypes.SignedVariablePart | A signed variable part. |

#### Return Values

| Name | Type | Description |
| ---- | ---- | ----------- |
| [0] | struct INitroTypes.RecoveredVariablePart | RecoveredVariablePart, identical to the supplied signedVariablePart, but with the signatures replaced with a signedBy bitmask. |

### _clearChallenge

```solidity
function _clearChallenge(bytes32 channelId, uint48 newTurnNumRecord) internal
```

Clears a challenge by updating the turnNumRecord and resetting the remaining channel storage fields, and emits a ChallengeCleared event.

_Clears a challenge by updating the turnNumRecord and resetting the remaining channel storage fields, and emits a ChallengeCleared event._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| channelId | bytes32 | Unique identifier for a channel. |
| newTurnNumRecord | uint48 | New turnNumRecord to overwrite existing value |

### _requireIncreasedTurnNumber

```solidity
function _requireIncreasedTurnNumber(bytes32 channelId, uint48 newTurnNumRecord) internal view
```

Checks that the submitted turnNumRecord is strictly greater than the turnNumRecord stored on chain.

_Checks that the submitted turnNumRecord is strictly greater than the turnNumRecord stored on chain._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| channelId | bytes32 | Unique identifier for a channel. |
| newTurnNumRecord | uint48 | New turnNumRecord intended to overwrite existing value |

### _requireNonDecreasedTurnNumber

```solidity
function _requireNonDecreasedTurnNumber(bytes32 channelId, uint48 newTurnNumRecord) internal view
```

Checks that the submitted turnNumRecord is greater than or equal to the turnNumRecord stored on chain.

_Checks that the submitted turnNumRecord is greater than or equal to the turnNumRecord stored on chain._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| channelId | bytes32 | Unique identifier for a channel. |
| newTurnNumRecord | uint48 | New turnNumRecord intended to overwrite existing value |

### _requireChannelNotFinalized

```solidity
function _requireChannelNotFinalized(bytes32 channelId) internal view
```

Checks that a given channel is NOT in the Finalized mode.

_Checks that a given channel is in the Challenge mode._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| channelId | bytes32 | Unique identifier for a channel. |

### _requireChannelOpen

```solidity
function _requireChannelOpen(bytes32 channelId) internal view
```

Checks that a given channel is in the Open mode.

_Checks that a given channel is in the Challenge mode._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| channelId | bytes32 | Unique identifier for a channel. |

### _matchesStatus

```solidity
function _matchesStatus(struct IStatusManager.ChannelData data, bytes32 s) internal pure returns (bool)
```

Checks that a given ChannelData struct matches a supplied bytes32 when formatted for storage.

_Checks that a given ChannelData struct matches a supplied bytes32 when formatted for storage._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| data | struct IStatusManager.ChannelData | A given ChannelData data structure. |
| s | bytes32 | Some data in on-chain storage format. |


# Solidity API

## IForceMove

_The IForceMove interface defines the interface that an implementation of ForceMove should implement. ForceMove protocol allows state channels to be adjudicated and finalized._

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
| proof | struct INitroTypes.SignedVariablePart[] | Additional proof material (in the form of an array of signed states) which completes the support proof. |
| candidate | struct INitroTypes.SignedVariablePart | A candidate state (along with signatures) which is being claimed to be supported. |
| challengerSig | struct INitroTypes.Signature | The signature of a participant on the keccak256 of the abi.encode of (supportedStateHash, 'forceMove'). |

### checkpoint

```solidity
function checkpoint(struct INitroTypes.FixedPart fixedPart, struct INitroTypes.SignedVariablePart[] proof, struct INitroTypes.SignedVariablePart candidate) external
```

Overwrites the `turnNumRecord` stored against a channel by providing a candidate with higher turn number.

_Overwrites the `turnNumRecord` stored against a channel by providing a candidate with higher turn number._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| fixedPart | struct INitroTypes.FixedPart | Data describing properties of the state channel that do not change with state updates. |
| proof | struct INitroTypes.SignedVariablePart[] | Additional proof material (in the form of an array of signed states) which completes the support proof. |
| candidate | struct INitroTypes.SignedVariablePart | A candidate state (along with signatures) which is being claimed to be supported. |

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
| proof | struct INitroTypes.SignedVariablePart[] | Additional proof material (in the form of an array of signed states) which completes the support proof. |
| candidate | struct INitroTypes.SignedVariablePart | A candidate state (along with signatures) which is being claimed to be supported. |

### ChallengeRegistered

```solidity
event ChallengeRegistered(bytes32 channelId, uint48 turnNumRecord, uint48 finalizesAt, bool isFinal, struct INitroTypes.FixedPart fixedPart, struct INitroTypes.SignedVariablePart[] proof, struct INitroTypes.SignedVariablePart candidate)
```

_Indicates that a challenge has been registered against `channelId`._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| channelId | bytes32 | Unique identifier for a state channel. |
| turnNumRecord | uint48 | A turnNum that (the adjudicator knows) is supported by a signature from each participant. |
| finalizesAt | uint48 | The unix timestamp when `channelId` will finalize. |
| isFinal | bool | Boolean denoting whether the challenge state is final. |
| fixedPart | struct INitroTypes.FixedPart | Data describing properties of the state channel that do not change with state updates. |
| proof | struct INitroTypes.SignedVariablePart[] | Additional proof material (in the form of an array of signed states) which completes the support proof. |
| candidate | struct INitroTypes.SignedVariablePart | A candidate state (along with signatures) which is being claimed to be supported. |

### ChallengeCleared

```solidity
event ChallengeCleared(bytes32 channelId, uint48 newTurnNumRecord)
```

_Indicates that a challenge, previously registered against `channelId`, has been cleared._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| channelId | bytes32 | Unique identifier for a state channel. |
| newTurnNumRecord | uint48 | A turnNum that (the adjudicator knows) is supported by a signature from each participant. |

### Concluded

```solidity
event Concluded(bytes32 channelId, uint48 finalizesAt)
```

_Indicates that a challenge has been registered against `channelId`._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| channelId | bytes32 | Unique identifier for a state channel. |
| finalizesAt | uint48 | The unix timestamp when `channelId` finalized. |


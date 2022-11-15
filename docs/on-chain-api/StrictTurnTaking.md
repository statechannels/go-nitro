# Solidity API

## StrictTurnTaking

### requireValidTurnTaking

```solidity
function requireValidTurnTaking(struct INitroTypes.FixedPart fixedPart, struct INitroTypes.RecoveredVariablePart[] proof, struct INitroTypes.RecoveredVariablePart candidate) internal pure
```

Require supplied arguments to comply with turn taking logic, i.e. each participant signed the one state, they were mover for.

_Require supplied arguments to comply with turn taking logic, i.e. each participant signed the one state, they were mover for._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| fixedPart | struct INitroTypes.FixedPart | Data describing properties of the state channel that do not change with state updates. |
| proof | struct INitroTypes.RecoveredVariablePart[] | Array of recovered variable parts which constitutes a support proof for the candidate. The proof is a validation for the supplied candidate. |
| candidate | struct INitroTypes.RecoveredVariablePart | Recovered variable part the proof was supplied for. The candidate state is supported by proof states. |

### isSignedByMover

```solidity
function isSignedByMover(struct INitroTypes.FixedPart fixedPart, struct INitroTypes.RecoveredVariablePart recoveredVariablePart) internal pure
```

Require supplied state is signed by its corresponding moving participant.

_Require supplied state is signed by its corresponding moving participant._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| fixedPart | struct INitroTypes.FixedPart | Data describing properties of the state channel that do not change with state updates. |
| recoveredVariablePart | struct INitroTypes.RecoveredVariablePart | A struct describing dynamic properties of the state channel, that must be signed by moving participant. |

### requireHasTurnNum

```solidity
function requireHasTurnNum(struct INitroTypes.VariablePart variablePart, uint48 turnNum) internal pure
```

Require supplied variable part has specified turn number.

_Require supplied variable part has specified turn number._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| variablePart | struct INitroTypes.VariablePart | Variable part to check turn number of. |
| turnNum | uint48 | Turn number to compare with. |

### _moverAddress

```solidity
function _moverAddress(address[] participants, uint48 turnNum) internal pure returns (address)
```

Find moving participant address based on state turn number.

_Find moving participant address based on state turn number._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| participants | address[] | Array of participant addresses. |
| turnNum | uint48 | State turn number. |

#### Return Values

| Name | Type | Description |
| ---- | ---- | ----------- |
| [0] | address | address Moving partitipant address. |

### _requireValidInput

```solidity
function _requireValidInput(uint256 numParticipants, uint256 numProofStates) internal pure
```

Validate input for turn taking logic.

_Validate input for turn taking logic._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| numParticipants | uint256 | Number of participants in a channel. |
| numProofStates | uint256 | Number of proof states submitted. |


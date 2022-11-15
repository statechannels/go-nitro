# Solidity API

## SingleAssetPayments

_The SingleAssetPayments contract complies with the ForceMoveApp interface, uses strict turn taking logic and implements a simple payment channel with a single asset type only._

### requireStateSupported

```solidity
function requireStateSupported(struct INitroTypes.FixedPart fixedPart, struct INitroTypes.RecoveredVariablePart[] proof, struct INitroTypes.RecoveredVariablePart candidate) external pure
```

Encodes application-specific rules for a particular ForceMove-compliant state channel. Must revert when invalid support proof and a candidate are supplied.

_Encodes application-specific rules for a particular ForceMove-compliant state channel. Must revert when invalid support proof and a candidate are supplied._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| fixedPart | struct INitroTypes.FixedPart | Fixed part of the state channel. |
| proof | struct INitroTypes.RecoveredVariablePart[] | Array of recovered variable parts which constitutes a support proof for the candidate. |
| candidate | struct INitroTypes.RecoveredVariablePart | Recovered variable part the proof was supplied for. |

### _requireValidOutcome

```solidity
function _requireValidOutcome(uint256 nParticipants, struct ExitFormat.SingleAssetExit[] outcome) internal pure
```

Require specific rules in outcome are followed.

_Require specific rules in outcome are followed._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| nParticipants | uint256 | Number of participants in a channel. |
| outcome | struct ExitFormat.SingleAssetExit[] | Outcome to check. |

### _requireValidTransition

```solidity
function _requireValidTransition(uint256 nParticipants, struct INitroTypes.VariablePart a, struct INitroTypes.VariablePart b) internal pure
```

Require specific rules in variable parts are followed when progressing state.

_Require specific rules in variable parts are followed when progressing state._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| nParticipants | uint256 | Number of participants in a channel. |
| a | struct INitroTypes.VariablePart | Variable part to progress from. |
| b | struct INitroTypes.VariablePart | Variable part to progress to. |


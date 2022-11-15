# Solidity API

## IForceMoveApp

_The IForceMoveApp interface calls for its children to implement an application-specific requireStateSupported function, defining the state machine of a ForceMove state channel DApp._

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
| proof | struct INitroTypes.RecoveredVariablePart[] | Array of recovered variable parts which constitutes a support proof for the candidate. May be omitted when `candidate` constitutes a support proof itself. |
| candidate | struct INitroTypes.RecoveredVariablePart | Recovered variable part the proof was supplied for. Also may constitute a support proof itself. |


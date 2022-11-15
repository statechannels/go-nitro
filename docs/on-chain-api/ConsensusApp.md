# Solidity API

## ConsensusApp

_The ConsensusApp contracts complies with the ForceMoveApp interface and uses consensus signatures logic._

### requireStateSupported

```solidity
function requireStateSupported(struct INitroTypes.FixedPart fixedPart, struct INitroTypes.RecoveredVariablePart[] proof, struct INitroTypes.RecoveredVariablePart candidate) external pure
```

Encodes application-specific rules for a particular ForceMove-compliant state channel.

_Encodes application-specific rules for a particular ForceMove-compliant state channel._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| fixedPart | struct INitroTypes.FixedPart | Fixed part of the state channel. |
| proof | struct INitroTypes.RecoveredVariablePart[] | Array of recovered variable parts which constitutes a support proof for the candidate. |
| candidate | struct INitroTypes.RecoveredVariablePart | Recovered variable part the proof was supplied for. |


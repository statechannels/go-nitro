# Solidity API

## Consensus

_Library for consensus signatures logic, which implies that all participants have signed the candidate state, while supplying proof as empty._

### requireConsensus

```solidity
function requireConsensus(struct INitroTypes.FixedPart fixedPart, struct INitroTypes.RecoveredVariablePart[] proof, struct INitroTypes.RecoveredVariablePart candidate) internal pure
```

Require supplied arguments to comply with consensus signatures logic, i.e. each participant has signed the candidate state.

_Require supplied arguments to comply with consensus signatures logic, i.e. each participant has signed the candidate state._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| fixedPart | struct INitroTypes.FixedPart | Data describing properties of the state channel that do not change with state updates. |
| proof | struct INitroTypes.RecoveredVariablePart[] | Array of recovered variable parts which constitutes a support proof for the candidate. The proof is a validation for the supplied candidate. Must be empty. |
| candidate | struct INitroTypes.RecoveredVariablePart | Recovered variable part the proof was supplied for. The candidate state is supported by proof states. All participants must have signed this state. |


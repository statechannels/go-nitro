# Solidity API

## CountingApp

_The CountingApp contract complies with the ForceMoveApp interface and strict turn taking logic and allows only for a simple counter to be incremented. Used for testing purposes._

### CountingAppData

```solidity
struct CountingAppData {
  uint256 counter;
}
```

### appData

```solidity
function appData(bytes appDataBytes) internal pure returns (struct CountingApp.CountingAppData)
```

Decodes the appData.

_Decodes the appData._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| appDataBytes | bytes | The abi.encode of a CountingAppData struct describing the application-specific data. |

#### Return Values

| Name | Type | Description |
| ---- | ---- | ----------- |
| [0] | struct CountingApp.CountingAppData | A CountingAppData struct containing the application-specific data. |

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

### _requireIncrementedCounter

```solidity
function _requireIncrementedCounter(struct INitroTypes.RecoveredVariablePart b, struct INitroTypes.RecoveredVariablePart a) internal pure
```

Checks that counter encoded in first variable part equals an incremented counter in second variable part.

_Checks that counter encoded in first variable part equals an incremented counter in second variable part._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| b | struct INitroTypes.RecoveredVariablePart | RecoveredVariablePart with incremented counter. |
| a | struct INitroTypes.RecoveredVariablePart | RecoveredVariablePart with counter before incrementing. |

### _requireEqualOutcomes

```solidity
function _requireEqualOutcomes(struct INitroTypes.RecoveredVariablePart a, struct INitroTypes.RecoveredVariablePart b) internal pure
```

Checks that supplied signed variable parts contain the same outcome.

_Checks that supplied signed variable parts contain the same outcome._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| a | struct INitroTypes.RecoveredVariablePart | First RecoveredVariablePart. |
| b | struct INitroTypes.RecoveredVariablePart | Second RecoveredVariablePart. |


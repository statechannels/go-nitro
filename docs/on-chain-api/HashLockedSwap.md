# Solidity API

## HashLockedSwap

_The HashLockedSwap contract complies with the ForceMoveApp interface, uses strict turn taking logic and implements a HashLockedSwapped payment._

### AppData

```solidity
struct AppData {
  bytes32 h;
  bytes preImage;
}
```

### appData

```solidity
function appData(bytes appDataBytes) internal pure returns (struct HashLockedSwap.AppData)
```

Decodes the appData.

_Decodes the appData._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| appDataBytes | bytes | The abi.encode of a AppData struct describing the application-specific data. |

#### Return Values

| Name | Type | Description |
| ---- | ---- | ----------- |
| [0] | struct HashLockedSwap.AppData | AppData struct containing the application-specific data. |

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

### decode2PartyAllocation

```solidity
function decode2PartyAllocation(struct ExitFormat.SingleAssetExit[] outcome) private pure returns (struct ExitFormat.Allocation[] allocations)
```


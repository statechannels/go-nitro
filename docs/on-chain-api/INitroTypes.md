# Solidity API

## INitroTypes

### Signature

```solidity
struct Signature {
  uint8 v;
  bytes32 r;
  bytes32 s;
}
```

### FixedPart

```solidity
struct FixedPart {
  uint256 chainId;
  address[] participants;
  uint64 channelNonce;
  address appDefinition;
  uint48 challengeDuration;
}
```

### VariablePart

```solidity
struct VariablePart {
  struct ExitFormat.SingleAssetExit[] outcome;
  bytes appData;
  uint48 turnNum;
  bool isFinal;
}
```

### SignedVariablePart

```solidity
struct SignedVariablePart {
  struct INitroTypes.VariablePart variablePart;
  struct INitroTypes.Signature[] sigs;
}
```

### RecoveredVariablePart

```solidity
struct RecoveredVariablePart {
  struct INitroTypes.VariablePart variablePart;
  uint256 signedBy;
}
```


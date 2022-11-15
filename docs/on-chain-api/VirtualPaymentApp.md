# Solidity API

## VirtualPaymentApp

_The VirtualPaymentApp contract complies with the ForceMoveApp interface and allows payments to be made virtually from Alice to Bob (participants[0] to participants[n+1], where n is the number of intermediaries)._

### VoucherAmountAndSignature

```solidity
struct VoucherAmountAndSignature {
  uint256 amount;
  struct INitroTypes.Signature signature;
}
```

### AllocationIndices

```solidity
enum AllocationIndices {
  Alice,
  Bob
}
```

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

### requireProofOfUnanimousConsensusOnPostFund

```solidity
function requireProofOfUnanimousConsensusOnPostFund(struct INitroTypes.RecoveredVariablePart rVP, uint256 numParticipants) internal pure
```

### requireValidVoucher

```solidity
function requireValidVoucher(bytes appData, struct INitroTypes.FixedPart fixedPart) internal pure returns (uint256)
```

### requireCorrectAdjustments

```solidity
function requireCorrectAdjustments(struct ExitFormat.SingleAssetExit[] oldOutcome, struct ExitFormat.SingleAssetExit[] newOutcome, uint256 voucherAmount) internal pure
```


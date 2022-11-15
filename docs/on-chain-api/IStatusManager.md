# Solidity API

## IStatusManager

### ChannelMode

```solidity
enum ChannelMode {
  Open,
  Challenge,
  Finalized
}
```

### ChannelData

```solidity
struct ChannelData {
  uint48 turnNumRecord;
  uint48 finalizesAt;
  bytes32 stateHash;
  bytes32 outcomeHash;
}
```


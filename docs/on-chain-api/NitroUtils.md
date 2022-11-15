# Solidity API

## NitroUtils

### isSignedBy

```solidity
function isSignedBy(bytes32 stateHash, struct INitroTypes.Signature sig, address signer) internal pure returns (bool)
```

Require supplied stateHash is signed by signer.

_Require supplied stateHash is signed by signer._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| stateHash | bytes32 | State hash to check. |
| sig | struct INitroTypes.Signature | Signed state signature. |
| signer | address | Address which must have signed the state. |

#### Return Values

| Name | Type | Description |
| ---- | ---- | ----------- |
| [0] | bool | true if signer with sig has signed stateHash. |

### isClaimedSignedBy

```solidity
function isClaimedSignedBy(uint256 signedBy, uint8 participantIndex) internal pure returns (bool)
```

Check if supplied participantIndex bit is set to 1 in signedBy bit mask.

_Check if supplied partitipationIndex bit is set to 1 in signedBy bit mask._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| signedBy | uint256 | Bit mask field to check. |
| participantIndex | uint8 | Bit to check. |

#### Return Values

| Name | Type | Description |
| ---- | ---- | ----------- |
| [0] | bool | true if supplied partitipationIndex bit is set to 1 in signedBy bit mask. |

### isClaimedSignedOnlyBy

```solidity
function isClaimedSignedOnlyBy(uint256 signedBy, uint8 participantIndex) internal pure returns (bool)
```

Check if supplied participantIndex is the only bit set to 1 in signedBy bit mask.

_Check if supplied participantIndex is the only bit set to 1 in signedBy bit mask._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| signedBy | uint256 | Bit mask field to check. |
| participantIndex | uint8 | Bit to check. |

#### Return Values

| Name | Type | Description |
| ---- | ---- | ----------- |
| [0] | bool | true if supplied partitipationIndex bit is the only bit set to 1 in signedBy bit mask. |

### recoverSigner

```solidity
function recoverSigner(bytes32 _d, struct INitroTypes.Signature sig) internal pure returns (address)
```

Given a digest and ethereum digital signature, recover the signer.

_Given a digest and digital signature, recover the signer._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| _d | bytes32 | message digest. |
| sig | struct INitroTypes.Signature | ethereum digital signature. |

#### Return Values

| Name | Type | Description |
| ---- | ---- | ----------- |
| [0] | address | signer |

### getClaimedSignersNum

```solidity
function getClaimedSignersNum(uint256 signedBy) internal pure returns (uint8)
```

Count number of bits set to '1', specifying the number of participants which have signed the state.

_Count number of bits set to '1', specifying the number of participants which have signed the state._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| signedBy | uint256 | Bit mask field specifying which participants have signed the state. |

#### Return Values

| Name | Type | Description |
| ---- | ---- | ----------- |
| [0] | uint8 | amount of signers, which have signed the state. |

### getClaimedSignersIndices

```solidity
function getClaimedSignersIndices(uint256 signedBy) internal pure returns (uint8[])
```

Determine indices of participants who have signed the state.

_Determine indices of participants who have signed the state._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| signedBy | uint256 | Bit mask field specifying which participants have signed the state. |

#### Return Values

| Name | Type | Description |
| ---- | ---- | ----------- |
| [0] | uint8[] | signerIndices |

### getChannelId

```solidity
function getChannelId(struct INitroTypes.FixedPart fixedPart) internal pure returns (bytes32 channelId)
```

Computes the unique id of a channel.

_Computes the unique id of a channel._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| fixedPart | struct INitroTypes.FixedPart | Part of the state that does not change |

#### Return Values

| Name | Type | Description |
| ---- | ---- | ----------- |
| channelId | bytes32 |  |

### getChainID

```solidity
function getChainID() internal view returns (uint256)
```

### hashState

```solidity
function hashState(bytes32 channelId, bytes appData, struct ExitFormat.SingleAssetExit[] outcome, uint48 turnNum, bool isFinal) internal pure returns (bytes32)
```

Computes the hash of the state corresponding to the input data.

_Computes the hash of the state corresponding to the input data._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| channelId | bytes32 | Unique identifier for the channel |
| appData | bytes | Application specific data. |
| outcome | struct ExitFormat.SingleAssetExit[] | Outcome structure. |
| turnNum | uint48 | Turn number |
| isFinal | bool | Is the state final? |

#### Return Values

| Name | Type | Description |
| ---- | ---- | ----------- |
| [0] | bytes32 | The stateHash |

### hashState

```solidity
function hashState(struct INitroTypes.FixedPart fp, struct INitroTypes.VariablePart vp) internal pure returns (bytes32)
```

Computes the hash of the state corresponding to the input data.

_Computes the hash of the state corresponding to the input data._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| fp | struct INitroTypes.FixedPart | The FixedPart of the state |
| vp | struct INitroTypes.VariablePart | The VariablePart of the state |

#### Return Values

| Name | Type | Description |
| ---- | ---- | ----------- |
| [0] | bytes32 | The stateHash |

### hashOutcome

```solidity
function hashOutcome(struct ExitFormat.SingleAssetExit[] outcome) internal pure returns (bytes32)
```

Hashes the outcome structure. Internal helper.

_Hashes the outcome structure. Internal helper._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| outcome | struct ExitFormat.SingleAssetExit[] | Outcome structure to encode hash. |

#### Return Values

| Name | Type | Description |
| ---- | ---- | ----------- |
| [0] | bytes32 | bytes32 Hash of encoded outcome structure. |

### bytesEqual

```solidity
function bytesEqual(bytes _preBytes, bytes _postBytes) internal pure returns (bool)
```

Check for equality of two byte strings

_Check for equality of two byte strings_

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| _preBytes | bytes | One bytes string |
| _postBytes | bytes | The other bytes string |

#### Return Values

| Name | Type | Description |
| ---- | ---- | ----------- |
| [0] | bool | true if the bytes are identical, false otherwise. |


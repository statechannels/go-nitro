# Solidity API

## Token

_This contract extends an ERC20 implementation, and mints 10,000,000,000 tokens to the deploying account. Used for testing purposes._

### constructor

```solidity
constructor(address owner) public
```

_Constructor function minting 10 billion tokens to the owner. Do not use msg.sender for default owner as that will not work with CREATE2_

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| owner | address | Tokens are minted to the owner address |


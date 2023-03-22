<h1 align="center">
<div><img src="https://statechannels.org/favicon.ico"> </div>
Nitro Protocol
</h1>

Smart contracts which implement nitro protocol for state channel networks on Ethereum, Filecoin and other EVM-compatible chains. Includes javascript and typescript support.

:new: There is an accompanying documentation [website](https://docs.statechannels.org/).

## Installation

```
.../my-statechannel-app> npm install --save @statechannels/nitro-protocol
```

## Getting started

### Building your state channel application contract against our interface:

Please see [this section of our docs](https://docs.statechannels.org/protocol-tutorial/0020-execution-rules/#core-protocol-rules). 

### Import precompiled artifacts for deployment/testing

```typescript
const {NitroAdjudicatorArtifact, ConsensusAppArtifact, VirtualPaymentAppArtifact} =
  require('@statechannels/nitro-protocol').ContractArtifacts;
```

### Import typescript types

```typescript
import {State} from '@statechannels/nitro-protocol';

const state: State = {
  channelNonce: 0,
  participants: ['0xalice...', '0xbob...'],
  appDefinition: '0xabc...',
  challengeDuration: '0x258',
  outcome: [],
  appData: '0x',
  turnNum: 0,
  isFinal: false,
};
```

For more information see [this section of our docs](https://docs.statechannels.org/protocol-tutorial/0010-states-channels/)

### Import javascript helper functions

```typescript
import {getChannelId, getFixedPart} from '@statechannels/nitro-protocol';

const channelId = getChannelId(getFixedPart(state));
```

## Development (GitHub)

We use hardhat to develop smart contracts. You can run the solidity compiler in watch mode like this:

```
npx hardhat watch compilation
```

### For the goerli testnet:

After successfully deploying you should see some changes to `addresses.json`. Please raise a pull request with this updated file.

```
INFURA_TOKEN=[your token here] RINKEBY_DEPLOYER_PK=[private key used for rinkeby deploy] yarn contract:deploy-goerli
```

### For mainnet

WARNING: This can be expensive. Each contract will take several million gas to deploy. Choose your moment and [gas price](https://ethgas.watch/) wisely!

```
INFURA_TOKEN=[your token here] MAINNET_DEPLOYER_PK=[private key used for mainnet deploy] yarn contract:deploy-mainnet --gasprice [your-chosen-gasPrice-here]
```

### To a local blockchain (for testing)

Contract deployment is handled automatically by our test setup scripts. Note that a **different** set of contracts is deployed when testing. Those contracts expose some helper functions that should not exist on production contracts.

## Verifying on etherscan

This is a somewhat manual process, but easier than using the etherscan GUI.

After deployment, run

```
ETHERSCAN_API_KEY=<a-secret> INFURA_TOKEN=<another-secret> yarn hardhat --network rinkeby verify <DeployedContractAddress> 'ConstructorArgs'
```

for each contract you wish to verify. Swap rinkeby for mainnet as appropriate.

You need to provide both `ETHERSCAN_API_KEY` and `INFURA_TOKEN` for this to work. For more info, see the [docs](https://hardhat.org/plugins/nomiclabs-hardhat-etherscan.html).

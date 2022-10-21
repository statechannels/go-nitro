import pick from 'lodash.pick';

import FULLNitroAdjudicatorArtifact from '../artifacts/contracts/NitroAdjudicator.sol/NitroAdjudicator.json';
import FULLConsensusAppArtifact from '../artifacts/contracts/ConsensusApp.sol/ConsensusApp.json';
import FULLVirtualPaymentAppArtifact from '../artifacts/contracts/VirtualPaymentApp.sol/VirtualPaymentApp.json';

interface ArtifactT {
  _format: string;
  contractName: string;
  sourceName: string;
  abi: object;
  bytecode: string;
  deployedBytecode: string;
  linkReferences: object;
  deployedLinkReferences: object;
}

// https://hardhat.org/guides/compile-contracts.html#artifacts
const fields = [
  'contractName',
  'abi',
  'bytecode',
  'deployedBytecode',
  'linkReferences',
  'deployedLinkReferences',
];

interface MinimalArtifact {
  contractName: string;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  abi: any;
  bytecode: string;
  deployedBytecode: string;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  linkReferences: any;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  deployedLinkReferences: any;
}

const minimize = (artifact: ArtifactT) => pick(artifact, fields) as MinimalArtifact;

export const ContractArtifacts = {
  NitroAdjudicatorArtifact: minimize(FULLNitroAdjudicatorArtifact),
  ConsensusAppArtifact: minimize(FULLConsensusAppArtifact),
  VirtualPaymentAppArtifact: minimize(FULLVirtualPaymentAppArtifact),
};

export {
  DepositedEvent,
  getDepositedEvent,
  convertBytes32ToAddress,
  convertAddressToBytes32,
} from './contract/multi-asset-holder';
export {
  getChallengeRegisteredEvent,
  getChallengeClearedEvent,
  ChallengeRegisteredEvent,
} from './contract/challenge';
export {getChannelId, isExternalDestination} from './contract/channel';
export {
  validTransition,
  ForceMoveAppContractInterface,
  createValidTransitionTransaction,
} from './contract/force-move-app';
export {encodeOutcome, decodeOutcome, Outcome, AssetOutcome, hashOutcome} from './contract/outcome';
export {channelDataToStatus} from './contract/channel-storage';

export {State, VariablePart, getVariablePart, getFixedPart, hashState} from './contract/state';

export * from './signatures';
export * from './transactions';
export * from './contract/vouchers';

// types
export {Uint256, Bytes32} from './contract/types';

export * from './channel-mode';

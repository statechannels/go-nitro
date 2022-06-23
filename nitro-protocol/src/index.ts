/**
 * @packageDocumentation Smart contracts that implement nitro protocol for state channel networks on ethereum. Includes javascript and typescript support.
 *
 * @remarks
 *
 * Building your state channel application contract against our interface:
 *
 * ```solidity
 * pragma solidity 0.7.6;
 * pragma experimental ABIEncoderV2;
 *
 * import {ExitFormat as Outcome} from '@statechannels/exit-format/contracts/ExitFormat.sol';
 * import '@statechannels/go-nitro/nitro-protocol/contracts/interfaces/IForceMoveApp.sol';
 *
 * contract MyStateChannelApp is IForceMoveApp {
 *   function latestSupportedState(
 *     FixedPart fixedPart,
 *     SignedVariablePart[] calldata signedVariableParts
 *   ) external override pure returns (VariablePart memory) {

*     // Your logic ...
 *
 *     return signedVariableParts[signedVariableParts.length - 1].variablePart;
 *   }
 * }
 * ```
 *
 * Import precompiled artifacts for deployment/testing
 *
 * ```typescript
 * const {
 *   NitroAdjudicatorArtifact,
 *   TrivialAppArtifact,
 *   TokenArtifact,
 * } = require('@statechannels/nitro-protocol').ContractArtifacts;
 * ```
 *
 * Import typescript types
 *
 * ```typescript
 * import {Channel} from '@statechannels/nitro-protocol';
 *
 * const channel: Channel = {
 *   chainId: '0x1',
 *   channelNonce: 0,
 *   participants: ['0xalice...', '0xbob...'],
 * };
 * ```
 *
 * Import javascript helper functions
 *
 * ```typescript
 * import {getChannelId} from '@statechannels/nitro-protocol';
 *
 * const channelId = getChannelId(channel);
 * ```
 */

import pick from 'lodash.pick';

import FULLTokenArtifact from '../artifacts/contracts/Token.sol/Token.json';
import FULLNitroAdjudicatorArtifact from '../artifacts/contracts/NitroAdjudicator.sol/NitroAdjudicator.json';
import FULLTestNitroAdjudicatorArtifact from '../artifacts/contracts/test/TESTNitroAdjudicator.sol/TESTNitroAdjudicator.json';
import FULLCountingAppArtifact from '../artifacts/contracts/CountingApp.sol/CountingApp.json';
import FULLHashLockedSwapArtifact from '../artifacts/contracts/examples/HashLockedSwap.sol/HashLockedSwap.json';

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
  HashLockedSwapArtifact: minimize(FULLHashLockedSwapArtifact),
};

/*
 * Various test contract artifacts used for testing.
 * They expose helper functions to allow for easier testing.
 * They should NEVER be used in a production environment.
 */
export const TestContractArtifacts = {
  CountingAppArtifact: minimize(FULLCountingAppArtifact),
  TestNitroAdjudicatorArtifact: minimize(FULLTestNitroAdjudicatorArtifact),
  TokenArtifact: minimize(FULLTokenArtifact),
};

export {
  AssetOutcomeShortHand,
  getTestProvider,
  OutcomeShortHand,
  randomChannelId,
  randomExternalDestination,
  replaceAddressesAndBigNumberify,
} from '../test/test-helpers';
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
export {Channel, getChannelId, isExternalDestination} from './contract/channel';
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

// types
export {Uint256, Bytes32} from './contract/types';

export * from './channel-mode';

import {utils} from 'ethers';

import {Channel, getChannelId} from './channel';
import {encodeOutcome, Outcome} from './outcome';
import {Address, Bytes, Bytes32, Uint256, Uint48} from './types';

/**
 * Holds all of the data defining the state of a channel
 */
export interface State {
  turnNum: number; // TODO: This should maybe be a string b/c it is uint256 in solidity
  isFinal: boolean;
  channel: Channel;
  challengeDuration: number;
  outcome: Outcome;
  appDefinition: string;
  appData: string;
}

/**
 * The part of a State which does not ordinarily change during state channel updates
 */
export interface FixedPart {
  chainId: Uint256;
  participants: Address[];
  channelNonce: Uint48;
  appDefinition: Address;
  challengeDuration: Uint48;
}
/**
 * Extracts the FixedPart of a state
 * @param state a State
 * @returns the FixedPart, which does not ordinarily change during state channel updates
 */
export function getFixedPart(state: State): FixedPart {
  const {appDefinition, challengeDuration, channel} = state;
  const {chainId, participants, channelNonce} = channel;
  return {chainId, participants, channelNonce, appDefinition, challengeDuration};
}

/**
 * The part of a State which usually changes during state channel updates
 */
export interface VariablePart {
  outcome: Bytes;
  appData: Bytes; // any encoded app-related type encoded once more as bytes
  //(e.g. if in SC App uint256 is used, firstly enode appData as uint256, then as bytes)
}

/**
 * Extracts the VariablePart of a state
 * @param state a State
 * @returns the VariablePart, which usually changes during state channel updates
 */
export function getVariablePart(state: State): VariablePart {
  return {outcome: encodeOutcome(state.outcome), appData: encodeAppData(state.appData)};
}

/**
 * Encodes appData
 * @param appData appData of the state
 * @returns an array of bytes of apppData
 */
export function encodeAppData(appData: string): Bytes {
  return utils.defaultAbiCoder.encode(['bytes'], [appData]);
}

/**
 * Encodes and hashes the AppPart of a state
 * @param state a State
 * @returns a 32 byte keccak256 hash
 */
 export function hashAppPart(state: State): Bytes32 {
  const {challengeDuration, appDefinition, appData} = state;
  return utils.keccak256(
    utils.defaultAbiCoder.encode(
      ['uint256', 'address', 'bytes'],
      [challengeDuration, appDefinition, appData]
    )
  );
}

/**
 * Encodes a state
 * @param state a State
 * @returns bytes array encoding
 */
export function encodeState(state: State): Bytes {
  const {turnNum, isFinal, appData, outcome} = state;
  const channelId = getChannelId(getFixedPart(state));

  const appDataBytes = encodeAppData(appData);
  return utils.defaultAbiCoder.encode(
    [
      'bytes32',
      'bytes',
      {
        type: "tuple[]",
        components: [
            // @ts-ignore - reference ethers.utils.ParamType for more info on why certain properties are not present
            { name: "asset", type: "address" },
            // @ts-ignore
            { name: "metadata", type: "bytes" },
            {
                type: "tuple[]",
                name: "allocations",
                components: [
                    // @ts-ignore
                    { name: "destination", type: "bytes32" },
                    // @ts-ignore
                    { name: "amount", type: "uint256" },
                    // @ts-ignore
                    { name: "allocationType", type: "uint8" },
                    // @ts-ignore
                    { name: "metadata", type: "bytes" },
                ],
            },
        ],
      },
      'uint256',
      'bool',
    ],
    [channelId, appDataBytes, outcome, turnNum, isFinal]
  );
}

/**
 * Hashes a state
 * @param state a State
 * @returns a 32 byte keccak256 hash
 */
export function hashState(state: State): Bytes32 {
  return utils.keccak256(
    encodeState(state)
  );
}

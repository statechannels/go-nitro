import {defaultAbiCoder, ParamType} from '@ethersproject/abi';
import {Signature} from '@ethersproject/bytes';

import {encodeOutcome} from './outcome';
import {encodeAppData, FixedPart, State} from './state';
import {Bytes, Bytes32} from './types';

// redefinition to support EmbeddedApplication.sol logic
export interface VariablePart {
  outcome: string;
  appData: string;
}
// redefinition to support EmbeddedApplication.sol logic
export function getVariablePart(state: State): VariablePart {
  return {
    outcome: encodeOutcome(state.outcome),
    appData: encodeAppData(state.appData),
  };
}

export interface SupportProof {
  fixedPart: FixedPart;
  variableParts: [VariablePart, VariablePart] | [VariablePart];
  turnNumTo: number;
  sigs: [Signature, Signature];
  whoSignedWhat: [number, number];
}

export enum AlreadyMoved {
  'None',
  'A',
  'B',
  'AB',
}
export interface EmbeddedApplicationData {
  channelIdForX: Bytes32;
  supportProofForX: SupportProof;
  alreadyMoved: AlreadyMoved;
}

export function encodeEmbeddedApplicationData(data: EmbeddedApplicationData): Bytes {
  return defaultAbiCoder.encode(
    [
      {
        type: 'tuple',
        components: [
          {name: 'channelIdForX', type: 'bytes32'},
          {
            name: 'supportProofForX',
            type: 'tuple',
            components: [
              {
                name: 'fixedPart',
                type: 'tuple',
                components: [
                  {name: 'chainId', type: 'uint256'},
                  {name: 'participants', type: 'address[]'},
                  {name: 'channelNonce', type: 'uint48'},
                  {name: 'appDefinition', type: 'address'},
                  {name: 'challengeDuration', type: 'uint48'},
                ],
              },
              {
                name: 'variableParts',
                type: 'tuple[]',
                components: [
                  {name: 'outcome', type: 'bytes'},
                  {name: 'appData', type: 'bytes'},
                ],
              },
              {name: 'turnNumTo', type: 'uint48'},
              {
                name: 'sigs',
                type: 'tuple[2]',
                components: [
                  {name: 'v', type: 'uint8'},
                  {name: 'r', type: 'bytes32'},
                  {name: 's', type: 'bytes32'},
                ],
              },
              {name: 'whoSignedWhat', type: 'uint8[2]'},
            ],
          },
          {name: 'alreadyMoved', type: 'uint8'},
        ],
      } as ParamType,
    ],
    [data]
  );
}

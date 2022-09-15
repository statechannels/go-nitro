import {utils} from 'ethers';

import {FixedPart} from './state';
import {Address, Bytes32, Uint256, Uint64} from './types';

/**
 * Holds the parameters that define a channel (in particular, its id)
 */
export interface Channel {
  channelNonce: Uint64; // Unique identifier for each new channel created by the same participants on the same chain
  participants: Address[]; // List of participant addresses (corresponding to ECDSA signing keys used to sign state channel updates)
  chainId: Uint256; // Identifier of the chain where this channel is adjudicated and where assets are held
}

/**
 * Determines if the supplied 32 byte hex string represents an external destination (meaning funds will be paid _out_ of the adjudicator on chain)
 * @param bytes32 a destination
 * @returns true if the destination has 12 leading bytes as zero, false otherwise
 */
export function isExternalDestination(bytes32: Bytes32): boolean {
  return /^0x(0{24})([a-fA-F0-9]{40})$/.test(bytes32);
}

/**
 * Computes the unique id for the supplied channel
 * @param channel Parameters which determine the id
 * @returns a 32 byte hex string representing the id
 */
export function getChannelId(fixedPart: FixedPart): Bytes32 {
  const {chainId, participants, channelNonce, appDefinition, challengeDuration} = fixedPart;
  const channelId = utils.keccak256(
    utils.defaultAbiCoder.encode(
      ['uint256', 'address[]', 'uint256', 'address', 'uint48'],
      [chainId, participants, channelNonce, appDefinition, challengeDuration]
    )
  );
  if (isExternalDestination(channelId))
    throw Error('This channel would have an external destination as an id');
  return channelId;
}

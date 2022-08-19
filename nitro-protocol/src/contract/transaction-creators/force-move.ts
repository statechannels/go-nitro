import {Signature, ethers} from 'ethers';

import ForceMoveArtifact from '../../../artifacts/contracts/ForceMove.sol/ForceMove.json';
import {bindSignatures, signChallengeMessage} from '../../signatures';
import {getFixedPart, getVariablePart, separateProofAndCandidate, State} from '../state';

// https://github.com/ethers-io/ethers.js/issues/602#issuecomment-574671078
export const ForceMoveContractInterface = new ethers.utils.Interface(ForceMoveArtifact.abi);

interface CheckpointData {
  challengeState?: State;
  states: State[];
  signatures: Signature[];
  whoSignedWhat: number[];
}

export function createChallengeTransaction(
  states: State[], // in turnNum order [..,state-with-largestTurnNum]
  signatures: Signature[], // in participant order: [sig-from-p0, sig-from-p1, ...]
  whoSignedWhat: number[],
  challengerPrivateKey: string
): ethers.providers.TransactionRequest {
  // Sanity checks on expected lengths
  if (states.length === 0) {
    throw new Error('No states provided');
  }
  const {participants} = states[0].channel;
  if (participants.length !== signatures.length) {
    throw new Error(
      `Participants (length:${participants.length}) and signatures (length:${signatures.length}) need to be the same length`
    );
  }

  const fixedPart = getFixedPart(states[0]);
  const variableParts = states.map(s => getVariablePart(s));
  const {proof, candidate} = separateProofAndCandidate(
    bindSignatures(variableParts, signatures, whoSignedWhat)
  );

  // Q: Is there a reason why createForceMoveTransaction accepts a State[] and a Signature[]
  // Argument rather than a SignedState[] argument?
  // A: Yes, because the signatures must be passed in participant order: [sig-from-p0, sig-from-p1, ...]
  // and SignedStates[] won't comply with that in general. This function accepts the re-ordered sigs.
  const signedStates = states.map(s => ({
    state: s,
    signature: {v: 0, r: '', s: '', _vs: '', recoveryParam: 0},
  }));
  const challengerSignature = signChallengeMessage(signedStates, challengerPrivateKey);

  const data = ForceMoveContractInterface.encodeFunctionData('challenge', [
    fixedPart,
    proof,
    candidate,
    challengerSignature,
  ]);
  return {data};
}

export function createCheckpointTransaction({
  states,
  signatures,
  whoSignedWhat,
}: CheckpointData): ethers.providers.TransactionRequest {
  const data = ForceMoveContractInterface.encodeFunctionData(
    'checkpoint',
    checkpointArgs({states, signatures, whoSignedWhat})
  );

  return {data};
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export function checkpointArgs({states, signatures, whoSignedWhat}: CheckpointData): any[] {
  const fixedPart = getFixedPart(states[0]);
  const variableParts = states.map(s => getVariablePart(s));
  const {proof, candidate} = separateProofAndCandidate(
    bindSignatures(variableParts, signatures, whoSignedWhat)
  );

  return [fixedPart, proof, candidate];
}

export function createConcludeTransaction(
  states: State[],
  signatures: Signature[],
  whoSignedWhat: number[]
): ethers.providers.TransactionRequest {
  const data = ForceMoveContractInterface.encodeFunctionData(
    'conclude',
    concludeArgs(states, signatures, whoSignedWhat)
  );
  return {data};
}

export function concludeArgs(
  states: State[],
  signatures: Signature[],
  whoSignedWhat: number[]
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
): any[] {
  // Sanity checks on expected lengths
  if (states.length === 0) {
    throw new Error('No states provided');
  }
  const {participants} = states[0].channel;
  if (participants.length !== signatures.length) {
    throw new Error(
      `Participants (length:${participants.length}) and signatures (length:${signatures.length}) need to be the same length`
    );
  }

  const fixedPart = getFixedPart(states[0]);

  const variableParts = states.map(s => getVariablePart(s));
  const {proof, candidate} = separateProofAndCandidate(
    bindSignatures(variableParts, signatures, whoSignedWhat)
  );

  return [fixedPart, proof, candidate];
}

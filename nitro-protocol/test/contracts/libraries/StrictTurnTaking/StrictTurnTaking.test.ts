import {BigNumber, Contract, Wallet} from 'ethers';
import {it} from '@jest/globals';
import {expectRevert} from '@statechannels/devtools';

import {
  shortenedToRecoveredVariableParts,
  TurnNumToShortenedVariablePart,
} from '../../../../src/signatures';
import testStrictTurnTakingArtifact from '../../../../artifacts/contracts/test/TESTStrictTurnTaking.sol/TESTStrictTurnTaking.json';
import {
  generateParticipants,
  getCountingAppContractAddress,
  getTestProvider,
  setupContract,
} from '../../../test-helpers';
import {TESTStrictTurnTaking} from '../../../../typechain-types';
import {getFixedPart, getRandomNonce, getVariablePart, Outcome, State} from '../../../../src';
import {
  INVALID_NUMBER_OF_PROOF_STATES,
  INVALID_SIGNED_BY,
  TOO_MANY_PARTICIPANTS,
  WRONG_TURN_NUM,
} from '../../../../src/contract/transaction-creators/revert-reasons';
import {RecoveredVariablePart, separateProofAndCandidate} from '../../../../src/contract/state';
import {getSignedBy} from '../../../../src/bitfield-utils';
import {expectSucceedWithNoReturnValues} from '../../../tx-expect-wrappers';
const provider = getTestProvider();
let StrictTurnTaking: Contract & TESTStrictTurnTaking;

const challengeDuration = 0x1000;
const asset = Wallet.createRandom().address;
const defaultOutcome: Outcome = [
  {asset, allocations: [], assetMetadata: {assetType: 0, metadata: '0x'}},
];
const appDefinition = getCountingAppContractAddress();

const nParticipants = 3;
const {wallets, participants} = generateParticipants(nParticipants);

beforeAll(async () => {
  StrictTurnTaking = setupContract(
    provider,
    testStrictTurnTakingArtifact,
    process.env.TEST_STRICT_TURN_TAKING_ADDRESS
  ) as Contract & TESTStrictTurnTaking;
});

let channelNonce = getRandomNonce('StrictTurnTaking');
beforeEach(() => (channelNonce = BigNumber.from(channelNonce).add(1).toHexString()));

describe('isSignedByMover', () => {
  const accepts1 = 'should not revert when signed only by mover';

  const reverts1 = 'should revert when not signed by mover';
  const reverts2 = 'should revert when signed not only by mover';

  it.each`
    description | turnNum | signedBy  | reason
    ${accepts1} | ${3}    | ${[0]}    | ${undefined}
    ${reverts1} | ${3}    | ${[2]}    | ${INVALID_SIGNED_BY}
    ${reverts2} | ${3}    | ${[0, 1]} | ${INVALID_SIGNED_BY}
  `('$description', async ({turnNum, signedBy, reason}) => {
    const state: State = {
      turnNum,
      isFinal: false,
      participants,
      channelNonce,
      challengeDuration,
      outcome: defaultOutcome,
      appDefinition,
      appData: '0x',
    };

    const variablePart = getVariablePart(state);
    const fixedPart = getFixedPart(state);

    const rvp: RecoveredVariablePart = {
      variablePart,
      signedBy: BigNumber.from(getSignedBy(signedBy)).toHexString(),
    };

    if (reason) {
      await expectRevert(() => StrictTurnTaking.isSignedByMover(fixedPart, rvp), reason);
    } else {
      await expectSucceedWithNoReturnValues(() => StrictTurnTaking.isSignedByMover(fixedPart, rvp));
    }
  });
});

describe('moverAddress', () => {
  const accepts1 = 'return correct mover';
  const accepts2 = 'return correct mover for turnNum >= numParticipants';

  it.each`
    description | turnNum | expectedParticipantIdx
    ${accepts1} | ${0}    | ${0}
    ${accepts1} | ${1}    | ${1}
    ${accepts1} | ${2}    | ${2}
    ${accepts2} | ${3}    | ${0}
    ${accepts2} | ${7}    | ${1}
  `(
    '$description',
    async ({
      turnNum,
      expectedParticipantIdx,
    }: {
      turnNum: number;
      expectedParticipantIdx: number;
    }) => {
      expect(await StrictTurnTaking.moverAddress(participants, turnNum)).toEqual(
        wallets[expectedParticipantIdx].address
      );
    }
  );
});

describe('requireValidInput', () => {
  const accepts1 = 'accept when all rules are satisfied';

  const reverts1 = 'revert when supplied zero proof states';
  const reverts2 = 'revert when supplied not enough proof states';
  const reverts3 = 'revert when supplied excessive proof states';
  const reverts4 = 'revert when too many participants';

  it.each`
    description | nParticipants | numProof | reason
    ${accepts1} | ${2}          | ${1}     | ${undefined}
    ${accepts1} | ${4}          | ${3}     | ${undefined}
    ${reverts1} | ${2}          | ${0}     | ${INVALID_NUMBER_OF_PROOF_STATES}
    ${reverts2} | ${4}          | ${1}     | ${INVALID_NUMBER_OF_PROOF_STATES}
    ${reverts3} | ${2}          | ${2}     | ${INVALID_NUMBER_OF_PROOF_STATES}
    ${reverts4} | ${256}        | ${255}   | ${TOO_MANY_PARTICIPANTS}
  `(
    '$description',
    async ({
      nParticipants,
      numProof,
      reason,
    }: {
      description: string;
      nParticipants: number;
      numProof: number;
      reason: undefined | string;
    }) => {
      if (reason) {
        await expectRevert(
          () => StrictTurnTaking.requireValidInput(nParticipants, numProof),
          reason
        );
      } else {
        await expectSucceedWithNoReturnValues(() => StrictTurnTaking.requireValidInput(nParticipants, numProof));
      }
    }
  );
});

describe('requireValidTurnTaking', () => {
  const accepts1 = 'accept when strict turn taking from 0';
  const accepts2 = 'accept when strict turn taking not from 0';

  const reverts1 = 'revert when insufficient states';
  const reverts2 = 'revert when excess states';
  const reverts3 = 'revert when a state is signed by multiple participants';
  const reverts4 = 'revert when a state is not signed';
  const reverts5 = 'revert when a state signed by non mover';
  const reverts6 = 'revert when a turn number is skipped';

  it.each`
    description | turnNumToShortenedVariablePart                       | reason
    ${accepts1} | ${new Map([[0, [0]], [1, [1]], [2, [2]]])}           | ${undefined}
    ${accepts2} | ${new Map([[3, [0]], [4, [1]], [5, [2]]])}           | ${undefined}
    ${reverts1} | ${new Map([[0, [0]], [1, [1]]])}                     | ${INVALID_NUMBER_OF_PROOF_STATES}
    ${reverts2} | ${new Map([[0, [0]], [1, [1]], [2, [2]], [3, [0]]])} | ${INVALID_NUMBER_OF_PROOF_STATES}
    ${reverts3} | ${new Map([[0, [0]], [1, [1, 2]], [2, [2]]])}        | ${INVALID_SIGNED_BY}
    ${reverts4} | ${new Map([[0, [0]], [1, []], [2, [2]]])}            | ${INVALID_SIGNED_BY}
    ${reverts5} | ${new Map([[0, [0]], [1, [2]], [2, [1]]])}           | ${INVALID_SIGNED_BY}
    ${reverts6} | ${new Map([[0, [0]], [2, [1]], [3, [2]]])}           | ${WRONG_TURN_NUM}
  `(
    '$description',
    async ({
      turnNumToShortenedVariablePart,
      reason,
    }: {
      turnNumToShortenedVariablePart: TurnNumToShortenedVariablePart;
      reason: undefined | string;
    }) => {
      const state: State = {
        turnNum: 0,
        isFinal: false,
        participants,
        channelNonce,
        challengeDuration,
        outcome: defaultOutcome,
        appDefinition,
        appData: '0x',
      };

      const fixedPart = getFixedPart(state);

      const recoveredVP = shortenedToRecoveredVariableParts(turnNumToShortenedVariablePart);
      const {proof, candidate} = separateProofAndCandidate(recoveredVP);

      if (reason) {
        await expectRevert(() =>
          StrictTurnTaking.requireValidTurnTaking(fixedPart, proof, candidate)
        );
      } else {
        await expectSucceedWithNoReturnValues(() =>
          StrictTurnTaking.requireValidTurnTaking(fixedPart, proof, candidate)
        );
      }
    }
  );
});

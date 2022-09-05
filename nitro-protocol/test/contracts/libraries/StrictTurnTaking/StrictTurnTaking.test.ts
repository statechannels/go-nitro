import {BigNumber, Contract, Wallet} from 'ethers';
import {it} from '@jest/globals';
import {expectRevert} from '@statechannels/devtools';

import testStrictTurnTakingArtifact from '../../../../artifacts/contracts/test/TESTStrictTurnTaking.sol/TESTStrictTurnTaking.json';
import {
  getCountingAppContractAddress,
  getRandomNonce,
  getTestProvider,
  setupContract,
} from '../../../test-helpers';
import {TESTStrictTurnTaking} from '../../../../typechain-types';
import {Channel, getFixedPart, getVariablePart, Outcome, State} from '../../../../src';
import {
  INVALID_NUMBER_OF_PROOF,
  SIGNED_BY_NON_MOVER,
  TOO_MANY_PARTICIPANTS,
} from '../../../../src/contract/transaction-creators/revert-reasons';
import {RecoveredVariablePart} from '../../../../src/contract/state';
import {getSignedBy} from '../../../../src/bitfield-utils';
const provider = getTestProvider();
let StrictTurnTaking: Contract & TESTStrictTurnTaking;

const chainId = process.env.CHAIN_NETWORK_ID;
const challengeDuration = 0x1000;
const asset = Wallet.createRandom().address;
const defaultOutcome: Outcome = [{asset, allocations: [], metadata: '0x'}];
const appDefinition = getCountingAppContractAddress();
const participants = ['', '', ''];
const wallets = new Array(3);

// Populate wallets and participants array
for (let i = 0; i < 3; i++) {
  wallets[i] = Wallet.createRandom();
  participants[i] = wallets[i].address;
}

beforeAll(async () => {
  StrictTurnTaking = setupContract(
    provider,
    testStrictTurnTakingArtifact,
    process.env.TEST_STRICT_TURN_TAKING_ADDRESS
  ) as Contract & TESTStrictTurnTaking;
});

let channelNonce = getRandomNonce('StrictTurnTaking');
beforeEach(() => (channelNonce += 1));

describe('isSignedByMover', () => {
  const accepts1 = 'should not revert when signed only by mover';

  const reverts1 = 'should revert when not signed by mover';
  const INVALID_SIGNED_BY = SIGNED_BY_NON_MOVER;

  const reverts2 = 'should revert when signed not only by mover';

  it.each`
    description | turnNum | signedBy  | reason
    ${accepts1} | ${3}    | ${[0]}    | ${undefined}
    ${reverts1} | ${3}    | ${[2]}    | ${INVALID_SIGNED_BY}
    ${reverts2} | ${3}    | ${[0, 1]} | ${INVALID_SIGNED_BY}
  `('$description', async ({turnNum, signedBy, reason}) => {
    const channel: Channel = {
      chainId,
      participants,
      channelNonce,
    };

    const state: State = {
      turnNum,
      isFinal: false,
      channel,
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
      const txResult = (await StrictTurnTaking.isSignedByMover(fixedPart, rvp)) as any;

      // As 'requireStateSupported' method is constant (view or pure), if it succeedes, it returns an object with returned values
      // which in this case should be empty
      expect(txResult.length).toBe(0);
    }
  });
});

describe('moverAddress', () => {
  const accepts1 = 'return correct mover';
  const accepts2 = 'return correct mover for turnNum >= numParticipants';

  const participants = wallets.map(w => w.address);

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
  const accepts1 = 'accept when all rules are preserved';
  const reverts1 = 'revert when supplied zero proof states';
  const reverts2 = 'revert when supplied not enough proof states';
  const reverts3 = 'revert when supplied excessive proof states';
  const reverts4 = 'revert when too many participants';

  it.each`
    description | numParticipants | numProof | reason
    ${accepts1} | ${2}            | ${1}     | ${undefined}
    ${accepts1} | ${4}            | ${3}     | ${undefined}
    ${reverts1} | ${2}            | ${0}     | ${INVALID_NUMBER_OF_PROOF}
    ${reverts2} | ${4}            | ${1}     | ${INVALID_NUMBER_OF_PROOF}
    ${reverts3} | ${2}            | ${2}     | ${INVALID_NUMBER_OF_PROOF}
    ${reverts4} | ${256}          | ${255}   | ${TOO_MANY_PARTICIPANTS}
  `(
    '$description',
    async ({
      numParticipants,
      numProof,
      reason,
    }: {
      description: string;
      numParticipants: number;
      numProof: number;
      reason: undefined | string;
    }) => {
      if (reason) {
        await expectRevert(
          () => StrictTurnTaking.requireValidInput(numParticipants, numProof),
          reason
        );
      } else {
        const txResult = (await StrictTurnTaking.requireValidInput(
          numParticipants,
          numProof
        )) as any;

        // As 'requireStateSupported' method is constant (view or pure), if it succeedes, it returns an object with returned values
        // which in this case should be empty
        expect(txResult.length).toBe(0);
      }
    }
  );
});

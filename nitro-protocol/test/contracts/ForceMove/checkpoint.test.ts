import {expectRevert} from '@statechannels/devtools';
import {Contract, Wallet, ethers, BigNumber} from 'ethers';
import {it} from '@jest/globals';

const {HashZero} = ethers.constants;

const {defaultAbiCoder} = ethers.utils;

import ForceMoveArtifact from '../../../artifacts/contracts/test/TESTForceMove.sol/TESTForceMove.json';
import {getChannelId} from '../../../src/contract/channel';
import {channelDataToStatus} from '../../../src/contract/channel-storage';
import {Outcome} from '../../../src/contract/outcome';
import {
  getFixedPart,
  getVariablePart,
  separateProofAndCandidate,
  State,
} from '../../../src/contract/state';
import {
  CHANNEL_FINALIZED,
  TURN_NUM_RECORD_NOT_INCREASED,
  INVALID_SIGNED_BY,
  COUNTING_APP_INVALID_TRANSITION,
} from '../../../src/contract/transaction-creators/revert-reasons';
import {
  generateParticipants,
  getCountingAppContractAddress,
  getTestProvider,
  setupContract,
} from '../../test-helpers';
import {bindSignatures, getRandomNonce, signStates} from '../../../src';

import {testParams} from './types';

const provider = getTestProvider();
let ForceMove: Contract;

const participantsNum = 3;
const {wallets, participants} = generateParticipants(participantsNum);

const challengeDuration = 0x1000;
const asset = Wallet.createRandom().address;
const defaultOutcome: Outcome = [{asset, allocations: [], metadata: '0x'}];
let appDefinition: string;

beforeAll(async () => {
  ForceMove = setupContract(provider, ForceMoveArtifact, process.env.TEST_FORCE_MOVE_ADDRESS);
  appDefinition = getCountingAppContractAddress();
});

const valid = {
  whoSignedWhat: [0, 1, 2],
  appDatas: [0, 1, 2],
};
const invalidTransition = {
  whoSignedWhat: [0, 1, 2],
  appDatas: [0, 2, 1],
};
const unsupported = {
  whoSignedWhat: [0, 0, 0],
  appDatas: [0, 1, 2],
};

const itOpensTheChannelIf = 'It accepts valid input, and clears any existing challenge, if ';
const accepts1 = itOpensTheChannelIf + 'the slot is empty';
const accepts2 =
  itOpensTheChannelIf + 'there is a challenge and the existing turnNumRecord is increased';
const accepts3 =
  itOpensTheChannelIf + 'there is no challenge and the existing turnNumRecord is increased';

const itRevertsWhenOpenBut = 'It reverts when the channel is open, but ';
const reverts1 = itRevertsWhenOpenBut + 'the turnNumRecord is not increased.';
const reverts2 = itRevertsWhenOpenBut + 'there is an invalid transition';
const reverts3 = itRevertsWhenOpenBut + 'the final state is not supported';

const itRevertsWithChallengeBut = 'It reverts when there is an ongoing challenge, but ';
const reverts4 = itRevertsWithChallengeBut + 'the turnNumRecord is not increased.';
const reverts5 = itRevertsWithChallengeBut + 'there is an invalid transition';
const reverts6 = itRevertsWithChallengeBut + 'the final state is not supported';

const reverts7 = 'It reverts when a challenge has expired';

const future = 1e12;
const past = 1;
const never = '0x00';
const turnNumRecord = 7;

describe('checkpoint', () => {
  let channelNonce = getRandomNonce('checkpoint');
  beforeEach(() => (channelNonce = BigNumber.from(channelNonce).add(1).toHexString()));
  it.each`
    description | largestTurnNum                         | support              | finalizesAt  | reason
    ${accepts1} | ${turnNumRecord + 1}                   | ${valid}             | ${undefined} | ${undefined}
    ${accepts2} | ${turnNumRecord + 1}                   | ${valid}             | ${never}     | ${undefined}
    ${accepts3} | ${turnNumRecord + 1 + participantsNum} | ${valid}             | ${future}    | ${undefined}
    ${reverts1} | ${turnNumRecord}                       | ${valid}             | ${never}     | ${TURN_NUM_RECORD_NOT_INCREASED}
    ${reverts2} | ${turnNumRecord + 1}                   | ${invalidTransition} | ${never}     | ${COUNTING_APP_INVALID_TRANSITION}
    ${reverts3} | ${turnNumRecord + 1}                   | ${unsupported}       | ${never}     | ${INVALID_SIGNED_BY}
    ${reverts4} | ${turnNumRecord}                       | ${valid}             | ${future}    | ${TURN_NUM_RECORD_NOT_INCREASED}
    ${reverts5} | ${turnNumRecord + 1}                   | ${invalidTransition} | ${future}    | ${COUNTING_APP_INVALID_TRANSITION}
    ${reverts6} | ${turnNumRecord + 1}                   | ${unsupported}       | ${future}    | ${INVALID_SIGNED_BY}
    ${reverts7} | ${turnNumRecord + 1}                   | ${valid}             | ${past}      | ${CHANNEL_FINALIZED}
  `('$description', async ({largestTurnNum, support, finalizesAt, reason}: testParams) => {
    const {appDatas, whoSignedWhat} = support;

    const states: State[] = appDatas.map((data, idx) => ({
      turnNum: largestTurnNum - appDatas.length + 1 + idx,
      isFinal: false,
      channelNonce,
      participants,
      challengeDuration,
      outcome: defaultOutcome,
      appData: defaultAbiCoder.encode(['uint256'], [data]),
      appDefinition,
    }));

    const variableParts = states.map(state => getVariablePart(state));
    const fixedPart = getFixedPart(states[0]);
    const channelId = getChannelId(fixedPart);

    // Sign the states
    const signatures = await signStates(states, wallets, whoSignedWhat);
    const {proof, candidate} = separateProofAndCandidate(
      bindSignatures(variableParts, signatures, whoSignedWhat)
    );

    const isOpen = !!finalizesAt;
    const outcome = isOpen ? [] : defaultOutcome;

    const challengeState: State | undefined = isOpen
      ? undefined
      : {
          turnNum: turnNumRecord,
          isFinal: false,
          channelNonce,
          participants,
          outcome,
          appData: defaultAbiCoder.encode(['uint256'], [appDatas[0]]),
          appDefinition,
          challengeDuration,
        };

    const fingerprint = finalizesAt
      ? channelDataToStatus({
          turnNumRecord,
          finalizesAt,
          state: challengeState,
          outcome,
        })
      : HashZero;

    // Call public wrapper to set state (only works on test contract)
    await (await ForceMove.setStatus(channelId, fingerprint)).wait();
    expect(await ForceMove.statusOf(channelId)).toEqual(fingerprint);

    const tx = ForceMove.checkpoint(fixedPart, proof, candidate);
    if (reason) {
      await expectRevert(() => tx, reason);
    } else {
      const receipt = await (await tx).wait();
      const event = receipt.events.pop();
      expect(event.args).toMatchObject({
        channelId,
        newTurnNumRecord: largestTurnNum,
      });

      const expectedChannelStorageHash = channelDataToStatus({
        turnNumRecord: largestTurnNum,
        finalizesAt: 0x0,
      });

      // Check channelStorageHash against the expected value
      expect(await ForceMove.statusOf(channelId)).toEqual(expectedChannelStorageHash);
    }
  });
});

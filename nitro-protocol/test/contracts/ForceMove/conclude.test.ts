import {expectRevert} from '@statechannels/devtools';
import {ethers, Contract, Wallet} from 'ethers';
const {HashZero} = ethers.constants;
const {defaultAbiCoder} = ethers.utils;
import {it} from '@jest/globals';

import ForceMoveArtifact from '../../../artifacts/contracts/test/TESTForceMove.sol/TESTForceMove.json';
import {Channel, getChannelId} from '../../../src/contract/channel';
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
  MOVER_SIGNED_EARLIER_STATE,
} from '../../../src/contract/transaction-creators/revert-reasons';
import {
  clearedChallengeFingerprint,
  finalizedFingerprint,
  getCountingAppContractAddress,
  getRandomNonce,
  getTestProvider,
  ongoingChallengeFingerprint,
  setupContract,
} from '../../test-helpers';
import {bindSignatures, signStates} from '../../../src';

const provider = getTestProvider();
let ForceMove: Contract;
const chainId = process.env.CHAIN_NETWORK_ID;
const participants = ['', '', ''];
const wallets = new Array(3);
const challengeDuration = 0x1000;
const asset = Wallet.createRandom().address;
const outcome: Outcome = [{asset, allocations: [], metadata: '0x'}];
let appDefinition: string;

const nParticipants = 3;
// Populate wallets and participants array
for (let i = 0; i < nParticipants; i++) {
  wallets[i] = Wallet.createRandom();
  participants[i] = wallets[i].address;
}
beforeAll(async () => {
  ForceMove = setupContract(provider, ForceMoveArtifact, process.env.TEST_FORCE_MOVE_ADDRESS);
  appDefinition = getCountingAppContractAddress();
});

const acceptsWhenOpenIf =
  'It accepts when the channel is open, and sets the channel storage correctly, if ';
const accepts1 = acceptsWhenOpenIf + 'passed n states, and the slot is empty';
const accepts2 = acceptsWhenOpenIf + 'passed one state, and the slot is empty';
const accepts3 = acceptsWhenOpenIf + 'the largestTurnNum is large enough';
const accepts6 =
  acceptsWhenOpenIf + 'despite the largest turn number being less than turnNumRecord';
const accepts7 = acceptsWhenOpenIf + 'the largest turn number is not large enough';

const acceptsWhenChallengeOngoingIf =
  'It accepts when there is an ongoing challenge, and sets the channel storage correctly, if ';
const accepts4 = acceptsWhenChallengeOngoingIf + 'passed n states';
const accepts5 = acceptsWhenChallengeOngoingIf + 'passed one state';

const reverts1 = 'It reverts when the channel is open, but the final state is not supported';
const reverts2 =
  'It reverts when there is an ongoing challenge, but the final state is not supported';
const reverts3 = 'It reverts when the outcome is already finalized';

const threeStates = {
  whoSignedWhat: [0, 1, 2],
  appData: [0, 1, 2],
};
const oneState = {
  whoSignedWhat: [0, 0, 0],
  appData: [0],
};
const unsupported = {
  whoSignedWhat: [0, 0, 0],
  appData: [0, 0, 0],
};
const turnNumRecord = 5;
const channelOpen = clearedChallengeFingerprint(turnNumRecord);
const challengeOngoing = ongoingChallengeFingerprint(turnNumRecord);
const finalized = finalizedFingerprint(turnNumRecord);

let channelNonce = getRandomNonce('conclude');
describe('conclude', () => {
  beforeEach(() => (channelNonce += 1));
  it.each`
    description | initialFingerprint  | largestTurnNum                   | support        | reasonString
    ${accepts1} | ${HashZero}         | ${turnNumRecord - nParticipants} | ${threeStates} | ${undefined}
    ${accepts2} | ${HashZero}         | ${turnNumRecord - 1}             | ${oneState}    | ${undefined}
    ${accepts2} | ${HashZero}         | ${turnNumRecord + 1}             | ${oneState}    | ${undefined}
    ${accepts3} | ${channelOpen}      | ${turnNumRecord + 2}             | ${oneState}    | ${undefined}
    ${accepts4} | ${challengeOngoing} | ${turnNumRecord + 3}             | ${oneState}    | ${undefined}
    ${accepts5} | ${challengeOngoing} | ${turnNumRecord + 4}             | ${oneState}    | ${undefined}
    ${accepts6} | ${channelOpen}      | ${turnNumRecord - 1}             | ${oneState}    | ${undefined}
    ${accepts7} | ${challengeOngoing} | ${turnNumRecord - 1}             | ${oneState}    | ${undefined}
    ${reverts1} | ${channelOpen}      | ${turnNumRecord + nParticipants} | ${unsupported} | ${MOVER_SIGNED_EARLIER_STATE}
    ${reverts2} | ${challengeOngoing} | ${turnNumRecord + nParticipants} | ${unsupported} | ${MOVER_SIGNED_EARLIER_STATE}
    ${reverts3} | ${finalized}        | ${turnNumRecord + 1}             | ${oneState}    | ${CHANNEL_FINALIZED}
  `(
    '$description', // For the purposes of this test, chainId and participants are fixed, making channelId 1-1 with channelNonce
    async ({initialFingerprint, largestTurnNum, support, reasonString}) => {
      const channel: Channel = {chainId, participants, channelNonce};
      const {appData, whoSignedWhat} = support;
      const numStates = appData.length;

      const states: State[] = [];
      for (let i = 1; i <= numStates; i++) {
        states.push({
          isFinal: true,
          channel,
          outcome,
          appDefinition,
          appData: defaultAbiCoder.encode(['uint256'], [appData[i - 1]]),
          challengeDuration,
          turnNum: largestTurnNum + i - numStates,
        });
      }

      const channelId = getChannelId({...channel, appDefinition, challengeDuration});
      const variableParts = states.map(state => getVariablePart(state));
      const fixedPart = getFixedPart(states[0]);

      // Call public wrapper to set state (only works on test contract)
      await (await ForceMove.setStatus(channelId, initialFingerprint)).wait();
      expect(await ForceMove.statusOf(channelId)).toEqual(initialFingerprint);

      // Sign the states
      const signatures = await signStates(states, wallets, whoSignedWhat);
      const {proof, candidate} = separateProofAndCandidate(
        bindSignatures(variableParts, signatures, whoSignedWhat)
      );

      const tx = ForceMove.conclude(fixedPart, proof, candidate);
      if (reasonString) {
        await expectRevert(() => tx, reasonString);
      } else {
        const receipt = await (await tx).wait();
        const event = receipt.events.pop();
        const finalizesAt = (await provider.getBlock(receipt.blockNumber)).timestamp;
        expect(event.args).toMatchObject({channelId, finalizesAt});

        // Compute expected ChannelDataHash
        const blockTimestamp = (await provider.getBlock(receipt.blockNumber)).timestamp;
        const expectedFingerprint = channelDataToStatus({
          turnNumRecord: 0,
          finalizesAt: blockTimestamp,
          outcome,
        });

        // Check fingerprint against the expected value
        expect(await ForceMove.statusOf(channelId)).toEqual(expectedFingerprint);
      }
    }
  );
});

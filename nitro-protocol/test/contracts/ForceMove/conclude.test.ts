import {expectRevert} from '@statechannels/devtools';
import {ethers, Contract, Wallet, BigNumber} from 'ethers';
const {HashZero} = ethers.constants;
const {defaultAbiCoder} = ethers.utils;
import {it} from '@jest/globals';

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
  NONFINAL_STATE,
} from '../../../src/contract/transaction-creators/revert-reasons';
import {
  clearedChallengeFingerprint,
  finalizedFingerprint,
  generateParticipants,
  getCountingAppContractAddress,
  getRandomNonce,
  getTestProvider,
  ongoingChallengeFingerprint,
  setupContract,
} from '../../test-helpers';
import {bindSignatures, signStates} from '../../../src';

let ForceMove: Contract;
const provider = getTestProvider();
const chainId = process.env.CHAIN_NETWORK_ID;

const nParticipants = 3;
const {wallets, participants} = generateParticipants(nParticipants);

const challengeDuration = 0x1000;
const asset = Wallet.createRandom().address;
const outcome: Outcome = [{asset, allocations: [], metadata: '0x'}];
let appDefinition: string;

beforeAll(async () => {
  ForceMove = setupContract(provider, ForceMoveArtifact, process.env.TEST_FORCE_MOVE_ADDRESS);
  appDefinition = getCountingAppContractAddress();
});

const acceptsWhenOpenIf =
  'It accepts when the channel is open, and sets the channel storage correctly, if ';
const acceptsWhenChallengeOngoingIf =
  'It accepts when there is an ongoing challenge, and sets the channel storage correctly, if ';

const accepts1 = acceptsWhenOpenIf + 'passed one state, and the slot is empty';
const accepts2 = acceptsWhenOpenIf + 'the largestTurnNum is large enough';
const accepts3 = acceptsWhenChallengeOngoingIf + 'passed one state';
const accepts4 =
  acceptsWhenOpenIf + 'despite the largest turn number being less than turnNumRecord';
const accepts5 = acceptsWhenOpenIf + 'the largest turn number is not large enough';

const reverts1 = 'It reverts when the outcome is already finalized';
const reverts2 = 'It reverts when the state is not final';

const threeStates = {
  whoSignedWhat: [0, 1, 2],
  appData: [0, 1, 2],
};
const oneState = {
  whoSignedWhat: [0, 0, 0],
  appData: [0],
};
const turnNumRecord = 5;
const channelOpen = clearedChallengeFingerprint(turnNumRecord);
const challengeOngoing = ongoingChallengeFingerprint(turnNumRecord);
const finalized = finalizedFingerprint(turnNumRecord);

let channelNonce = getRandomNonce('conclude');
describe('conclude', () => {
  beforeEach(() => (channelNonce = BigNumber.from(channelNonce).add(1).toHexString()));
  it.each`
    description | initialFingerprint  | isFinal  | largestTurnNum       | support     | reasonString
    ${accepts1} | ${HashZero}         | ${true}  | ${turnNumRecord - 1} | ${oneState} | ${undefined}
    ${accepts1} | ${HashZero}         | ${true}  | ${turnNumRecord + 1} | ${oneState} | ${undefined}
    ${accepts2} | ${channelOpen}      | ${true}  | ${turnNumRecord + 2} | ${oneState} | ${undefined}
    ${accepts3} | ${challengeOngoing} | ${true}  | ${turnNumRecord + 4} | ${oneState} | ${undefined}
    ${accepts4} | ${channelOpen}      | ${true}  | ${turnNumRecord - 1} | ${oneState} | ${undefined}
    ${accepts5} | ${challengeOngoing} | ${true}  | ${turnNumRecord - 1} | ${oneState} | ${undefined}
    ${reverts1} | ${finalized}        | ${true}  | ${turnNumRecord + 1} | ${oneState} | ${CHANNEL_FINALIZED}
    ${reverts2} | ${HashZero}         | ${false} | ${turnNumRecord - 1} | ${oneState} | ${NONFINAL_STATE}
  `(
    '$description', // For the purposes of this test, chainId and participants are fixed, making channelId 1-1 with channelNonce
    async ({initialFingerprint, isFinal, largestTurnNum, support, reasonString}) => {
      const {appData, whoSignedWhat} = support;
      const numStates = appData.length;

      const states: State[] = [];
      for (let i = 1; i <= numStates; i++) {
        states.push({
          isFinal,
          chainId,
          participants,
          channelNonce,
          outcome,
          appDefinition,
          appData: defaultAbiCoder.encode(['uint256'], [appData[i - 1]]),
          challengeDuration,
          turnNum: largestTurnNum + i - numStates,
        });
      }

      const channelId = getChannelId({
        chainId,
        participants,
        channelNonce,
        appDefinition,
        challengeDuration,
      });
      const variableParts = states.map(state => getVariablePart(state));
      const fixedPart = getFixedPart(states[0]);

      // Call public wrapper to set state (only works on test contract)
      await (await ForceMove.setStatus(channelId, initialFingerprint)).wait();
      expect(await ForceMove.statusOf(channelId)).toEqual(initialFingerprint);

      // Sign the states
      const signatures = await signStates(states, wallets, whoSignedWhat);
      const {candidate} = separateProofAndCandidate(
        bindSignatures(variableParts, signatures, whoSignedWhat)
      );

      const tx = ForceMove.conclude(fixedPart, candidate);
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

import {expectRevert} from '@statechannels/devtools';
import {Contract, Wallet, ethers, Signature, BigNumber} from 'ethers';
import {it} from '@jest/globals';

const {HashZero} = ethers.constants;
const {defaultAbiCoder} = ethers.utils;

import ForceMoveArtifact from '../../../artifacts/contracts/test/TESTForceMove.sol/TESTForceMove.json';
import {getChannelId} from '../../../src/contract/channel';
import {channelDataToStatus, ChannelData} from '../../../src/contract/channel-storage';
import {
  getFixedPart,
  getVariablePart,
  separateProofAndCandidate,
  State,
} from '../../../src/contract/state';
import {
  CHALLENGER_NON_PARTICIPANT,
  CHANNEL_FINALIZED,
  INVALID_NUMBER_OF_PROOF_STATES,
  INVALID_SIGNATURE,
  TURN_NUM_RECORD_DECREASED,
  TURN_NUM_RECORD_NOT_INCREASED,
  COUNTING_APP_INVALID_TRANSITION,
} from '../../../src/contract/transaction-creators/revert-reasons';
import {getRandomNonce, Outcome, SignedState} from '../../../src/index';
import {
  bindSignatures,
  signChallengeMessage,
  signData,
  signState,
  signStates,
} from '../../../src/signatures';
import {
  clearedChallengeFingerprint,
  finalizedFingerprint,
  getCountingAppContractAddress,
  getTestProvider,
  largeOutcome,
  nonParticipant,
  ongoingChallengeFingerprint,
  parseOutcomeEventResult,
  setupContract,
} from '../../test-helpers';
import {createChallengeTransaction, NITRO_MAX_GAS} from '../../../src/transactions';
import {hashChallengeMessage} from '../../../src/contract/challenge';
import {MAX_OUTCOME_ITEMS} from '../../../src/contract/outcome';

import {transitionType} from './types';

let ForceMove: Contract;
const provider = getTestProvider();

const participants = ['', '', ''];
const wallets = new Array<Wallet>(3);

const challengeDuration = 86400; // 1 day
const outcome: Outcome = [
  {
    allocations: [],
    asset: Wallet.createRandom().address,
    assetMetadata: {assetType: 0, metadata: '0x'},
  },
];

const appDefinition = getCountingAppContractAddress();
const keys = [
  '0x8624ebe7364bb776f891ca339f0aaa820cc64cc9fca6a28eec71e6d8fc950f29',
  '0x275a2e2cd9314f53b42246694034a80119963097e3adf495fbf6d821dc8b6c8e',
  '0x1b7598002c59e7d9131d7e7c9d0ec48ed065a3ed04af56674497d6b0048f2d84',
];

// Populate wallets and participants array
for (let i = 0; i < 3; i++) {
  wallets[i] = new Wallet(keys[i]);
  participants[i] = wallets[i].address;
}

async function createTwoPartySignedCountingAppState(
  appData: number,
  turnNum: number,
  outcome: Outcome = []
) {
  return await signState(
    {
      turnNum,
      isFinal: false,
      appDefinition: getCountingAppContractAddress(),
      appData: defaultAbiCoder.encode(['uint256'], [appData]),
      outcome,
      channelNonce: '0x1',
      participants: [wallets[0].address, wallets[1].address],
      challengeDuration: 0xfff,
    },
    wallets[turnNum % 2].privateKey
  );
}

beforeAll(async () => {
  ForceMove = setupContract(provider, ForceMoveArtifact, process.env.TEST_FORCE_MOVE_ADDRESS|| "");
});

// Scenarios are synonymous with channelNonce:

const acceptsWhenOpen = 'It accepts for an open channel, and updates storage correctly, ';
const accepts1 = acceptsWhenOpen + 'when the slot is empty, 3 states submitted';
const accepts2 = acceptsWhenOpen + 'when the slot is not empty, 3 states submitted';
const accepts3 =
  acceptsWhenOpen + 'when the slot is not empty, 3 states submitted, open at largestTurnNum';

const acceptsWhenChallengePresent =
  'It accepts when a challenge is present, and updates storage correctly, ';
const accepts4 = acceptsWhenChallengePresent + 'when the turnNumRecord increases, 3 states';

const revertsWhenOpenIf = 'It reverts for an open channel if ';
const reverts1 = revertsWhenOpenIf + 'the turnNumRecord does not increase';
const reverts2a = revertsWhenOpenIf + 'the challengerSig is incorrect';
const reverts2b = revertsWhenOpenIf + 'the challengerSig is invalid';
const reverts3 = revertsWhenOpenIf + 'the states do not form a validTransition chain';

const reverts4 = 'It reverts when a challenge is present if the turnNumRecord does not increase';
const reverts5 = 'It reverts when the channel is finalized';
const reverts6 = 'It reverts when too few states are submitted';
const reverts7 = 'It reverts when too many states are submitted';

describe('challenge', () => {
  const threeStates = {appDatas: [0, 1, 2], whoSignedWhat: [0, 1, 2]};
  const fourStates = {appDatas: [0, 1, 2, 3], whoSignedWhat: [0, 1, 2, 0]};
  const oneState = {appDatas: [2], whoSignedWhat: [0, 0, 0]};
  const invalid = {appDatas: [0, 2, 1], whoSignedWhat: [0, 1, 2]};
  const largestTurnNum = 8;
  const isFinalCount = 0;
  const challenger = wallets[2];
  const empty = HashZero; // Equivalent to openAtZero
  const openAtFive = clearedChallengeFingerprint(5);
  const openAtLargestTurnNum = clearedChallengeFingerprint(largestTurnNum);
  const openAtTwenty = clearedChallengeFingerprint(20);
  const challengeAtFive = ongoingChallengeFingerprint(5);
  const challengeAtLargestTurnNum = ongoingChallengeFingerprint(largestTurnNum);
  const challengeAtTwenty = ongoingChallengeFingerprint(20);
  const finalizedAtFive = finalizedFingerprint(5);

  let channelNonce = getRandomNonce('challenge');
  beforeEach(() => (channelNonce = BigNumber.from(channelNonce).add(1).toHexString()));
  it.each`
    description  | initialFingerprint           | stateData      | challengeSignatureType | reasonString
    ${accepts1}  | ${empty}                     | ${threeStates} | ${'correct'}           | ${undefined}
    ${accepts2}  | ${openAtFive}                | ${threeStates} | ${'correct'}           | ${undefined}
    ${accepts3}  | ${openAtLargestTurnNum}      | ${threeStates} | ${'correct'}           | ${undefined}
    ${accepts4}  | ${challengeAtFive}           | ${threeStates} | ${'correct'}           | ${undefined}
    ${reverts1}  | ${openAtTwenty}              | ${threeStates} | ${'correct'}           | ${TURN_NUM_RECORD_DECREASED}
    ${reverts2a} | ${empty}                     | ${threeStates} | ${'incorrect'}         | ${CHALLENGER_NON_PARTICIPANT}
    ${reverts2b} | ${empty}                     | ${threeStates} | ${'invalid'}           | ${INVALID_SIGNATURE}
    ${reverts3}  | ${empty}                     | ${invalid}     | ${'correct'}           | ${COUNTING_APP_INVALID_TRANSITION}
    ${reverts4}  | ${challengeAtTwenty}         | ${threeStates} | ${'correct'}           | ${TURN_NUM_RECORD_NOT_INCREASED}
    ${reverts4}  | ${challengeAtLargestTurnNum} | ${threeStates} | ${'correct'}           | ${TURN_NUM_RECORD_NOT_INCREASED}
    ${reverts5}  | ${finalizedAtFive}           | ${threeStates} | ${'correct'}           | ${CHANNEL_FINALIZED}
    ${reverts6}  | ${empty}                     | ${oneState}    | ${'correct'}           | ${INVALID_NUMBER_OF_PROOF_STATES}
    ${reverts7}  | ${empty}                     | ${fourStates}  | ${'correct'}           | ${INVALID_NUMBER_OF_PROOF_STATES}
  `(
    '$description', // For the purposes of this test, participants are fixed, making channelId 1-1 with channelNonce
    async (tc) => {


      const {reasonString,challengeSignatureType,stateData,initialFingerprint} = tc as unknown as {initialFingerprint:string, stateData:transitionType, challengeSignatureType:string, reasonString:undefined|string}
      const {appDatas,whoSignedWhat} = stateData

      const states: State[] = appDatas.map((data, idx) => ({
        turnNum: largestTurnNum - appDatas.length + 1 + idx,
        isFinal: idx > appDatas.length - isFinalCount,
        participants,
        channelNonce,
        challengeDuration,
        outcome,
        appDefinition,
        appData: defaultAbiCoder.encode(['uint256'], [data]),
      }));
      const variableParts = states.map(state => getVariablePart(state));
      const fixedPart = getFixedPart(states[0]);
      const channelId = getChannelId(fixedPart);

      // Sign the states
      const signatures = await signStates(states, wallets, whoSignedWhat);
      const {proof, candidate} = separateProofAndCandidate(
        bindSignatures(variableParts, signatures, whoSignedWhat)
      );

      const challengeState: SignedState = {
        state: states[states.length - 1],
        signature: {v: 0, r: '', s: '', _vs: '', recoveryParam: 0} as Signature,
      };

      const correctChallengeSignature = signChallengeMessage(
        [challengeState],
        challenger.privateKey
      );
      let challengeSignature: ethers.Signature;

      switch (challengeSignatureType) {
        case 'incorrect':
          challengeSignature = signChallengeMessageByNonParticipant([challengeState]);
          break;
        case 'invalid':
          challengeSignature = {v: 1, s: HashZero, r: HashZero} as ethers.Signature;
          break;
        case 'correct':
        default:
          challengeSignature = correctChallengeSignature;
      }

      // Set current channelStorageHashes value
      await (await ForceMove.setStatus(channelId, initialFingerprint)).wait();

      const tx = ForceMove.challenge(fixedPart, proof, candidate, challengeSignature);
      if (reasonString) {
        await expectRevert(() => tx, reasonString);
      } else {
        const receipt = await (await tx).wait();
        const event = receipt.events.pop();

        // Catch ChallengeRegistered event
        const {
          channelId: eventChannelId,
          finalizesAt: eventFinalizesAt,

          proof: eventProof,
          candidate: eventCandidate,
        } = event.args;

        // Check this information is enough to respond
        expect(eventChannelId).toEqual(channelId);

        if (proof.length > 0) {
          expect(
            parseOutcomeEventResult(eventProof[eventProof.length - 1].variablePart.outcome)
          ).toEqual(proof[proof.length - 1].variablePart.outcome);
          expect(eventProof[eventProof.length - 1].variablePart.appData).toEqual(
            proof[proof.length - 1].variablePart.appData
          );
        }

        expect(parseOutcomeEventResult(eventCandidate.variablePart.outcome)).toEqual(
          candidate.variablePart.outcome
        );
        expect(eventCandidate.variablePart.appData).toEqual(candidate.variablePart.appData);

        const expectedChannelStorage: ChannelData = {
          turnNumRecord: largestTurnNum,
          finalizesAt: eventFinalizesAt,
          state: states[states.length - 1],
          outcome,
        };
        const expectedFingerprint = channelDataToStatus(expectedChannelStorage);

        // Check channelStorageHash against the expected value
        expect(await ForceMove.statusOf(channelId)).toEqual(expectedFingerprint);
      }
    }
  );
});

describe('challenge with transaction generator', () => {
  const twoPartyFixedPart = {
    channelNonce: '0x1',
    participants: [wallets[0].address, wallets[1].address],
    appDefinition,
    challengeDuration,
  };

  beforeEach(async () => {
    await (await ForceMove.setStatus(getChannelId(twoPartyFixedPart), HashZero)).wait();
  });
  // FIX: even if dropping channel status before each test, turn nums from prev tests are saved and can cause reverts
  it.each`
    description                                     | appData   | outcome                            | turnNums  | challenger
    ${'challenge(0,1) accepted'}                    | ${[0, 1]} | ${[]}                              | ${[1, 2]} | ${1}
    ${'challenge(1,2) accepted'}                    | ${[0, 1]} | ${[]}                              | ${[2, 3]} | ${0}
    ${'challenge(2,3) accepted, MAX_OUTCOME_ITEMS'} | ${[0, 1]} | ${largeOutcome(MAX_OUTCOME_ITEMS)} | ${[3, 4]} | ${0}
  `('$description', async (tc) => {
    const {appData,  turnNums, challenger} = tc as unknown as {appData:number[], turnNums:number[], challenger:number}
    const transactionRequest: ethers.providers.TransactionRequest = createChallengeTransaction(
      [
        await createTwoPartySignedCountingAppState(appData[0], turnNums[0]),
        await createTwoPartySignedCountingAppState(appData[1], turnNums[1]),
      ],
      wallets[challenger].privateKey
    );

    const signer = provider.getSigner();
    const response = await signer.sendTransaction({
      to: ForceMove.address,
      ...transactionRequest,
    });
    expect(BigNumber.from((await response.wait()).gasUsed).lt(BigNumber.from(NITRO_MAX_GAS))).toBe(
      true
    );
  });
});

function signChallengeMessageByNonParticipant(signedStates: SignedState[]): Signature {
  if (signedStates.length === 0) {
    throw new Error('At least one signed state must be provided');
  }
  const challengeState = signedStates[signedStates.length - 1].state;
  const challengeHash = hashChallengeMessage(challengeState);
  return signData(challengeHash, nonParticipant.privateKey);
}

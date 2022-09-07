import {it} from '@jest/globals';
import {Contract, Wallet} from 'ethers';
import {expectRevert} from '@statechannels/devtools';

import testConsensusArtifact from '../../../../artifacts/contracts/test/TESTConsensus.sol/TESTConsensus.json';
import {
  generateParticipants,
  getRandomNonce,
  getTestProvider,
  setupContract,
} from '../../../test-helpers';
import {TESTConsensus} from '../../../../typechain-types';
import {
  Channel,
  getFixedPart,
  Outcome,
  shortenedToRecoveredVariableParts,
  State,
  TurnNumToShortenedVariablePart,
} from '../../../../src';
import {
  NOT_UNANIMOUS,
  PROOF_SUPPLIED,
} from '../../../../src/contract/transaction-creators/revert-reasons';
import {separateProofAndCandidate} from '../../../../src/contract/state';
import {expectSucceed} from '../../../expect-succeed';
const provider = getTestProvider();
let Consensus: Contract & TESTConsensus;

const chainId = process.env.CHAIN_NETWORK_ID;
const challengeDuration = 0x1000;
const asset = Wallet.createRandom().address;
const defaultOutcome: Outcome = [{asset, allocations: [], metadata: '0x'}];
const appDefinition = process.env.CONSENSUS_APP_ADDRESS;

const nParticipants = 3;
const {participants} = generateParticipants(nParticipants);

beforeAll(async () => {
  Consensus = setupContract(
    provider,
    testConsensusArtifact,
    process.env.TEST_CONSENSUS_ADDRESS
  ) as Contract & TESTConsensus;
});

let channelNonce = getRandomNonce('Consensus');
beforeEach(() => (channelNonce += 1));

describe('requireConsensus', () => {
  const accepts1 = 'accept when signed by all (one turnNum)';
  const accepts2 = 'accept when signed by all (other turnNum)';

  const reverts1 = 'revert when not signed by all';
  const reverts2 = 'revert when not signed at all';
  const reverts3 = 'revert when supplied proof state';

  it.each`
    description | turnNumToShortenedVariablePart         | reason
    ${accepts1} | ${new Map([[0, [0, 1, 2]]])}           | ${undefined}
    ${accepts2} | ${new Map([[2, [0, 1, 2]]])}           | ${undefined}
    ${reverts1} | ${new Map([[0, [0, 1]]])}              | ${NOT_UNANIMOUS}
    ${reverts2} | ${new Map([[0, []]])}                  | ${NOT_UNANIMOUS}
    ${reverts3} | ${new Map([[0, [0]], [1, [0, 1, 2]]])} | ${PROOF_SUPPLIED}
  `(
    '$description',
    async ({
      turnNumToShortenedVariablePart,
      reason,
    }: {
      turnNumToShortenedVariablePart: TurnNumToShortenedVariablePart;
      reason: undefined | string;
    }) => {
      const channel: Channel = {
        chainId,
        participants,
        channelNonce,
      };

      const state: State = {
        turnNum: 0,
        isFinal: false,
        channel,
        challengeDuration,
        outcome: defaultOutcome,
        appDefinition,
        appData: '0x',
      };

      const fixedPart = getFixedPart(state);

      const recoveredVP = shortenedToRecoveredVariableParts(turnNumToShortenedVariablePart);
      const {proof, candidate} = separateProofAndCandidate(recoveredVP);

      if (reason) {
        await expectRevert(() => Consensus.requireConsensus(fixedPart, proof, candidate));
      } else {
        await expectSucceed(() => Consensus.requireConsensus(fixedPart, proof, candidate));
      }
    }
  );
});

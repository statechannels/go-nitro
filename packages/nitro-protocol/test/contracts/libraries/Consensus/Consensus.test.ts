import {it} from '@jest/globals';
import {BigNumber, Contract, Wallet} from 'ethers';
import {expectRevert} from '@statechannels/devtools';

import testConsensusArtifact from '../../../../artifacts/contracts/test/TESTConsensus.sol/TESTConsensus.json';
import {generateParticipants, getTestProvider, setupContract} from '../../../test-helpers';
import {TESTConsensus} from '../../../../typechain-types';
import {
  getFixedPart,
  getRandomNonce,
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
import {expectSucceedWithNoReturnValues} from '../../../tx-expect-wrappers';
const provider = getTestProvider();
let Consensus: Contract & TESTConsensus;

const challengeDuration = 0x1000;
const asset = Wallet.createRandom().address;
const defaultOutcome: Outcome = [
  {asset, allocations: [], assetMetadata: {assetType: 0, metadata: '0x'}},
];
const appDefinition = process.env.CONSENSUS_APP_ADDRESS;

const nParticipants = 3;
const {participants} = generateParticipants(nParticipants);

beforeAll(async () => {
  Consensus = setupContract(
    provider,
    testConsensusArtifact,
    process.env.TEST_CONSENSUS_ADDRESS ||""
  ) as Contract & TESTConsensus;
});

let channelNonce = getRandomNonce('Consensus');
beforeEach(() => (channelNonce = BigNumber.from(channelNonce).add(1).toHexString()));

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
    }: any) => {
      const state: State = {
        turnNum: 0,
        isFinal: false,
        participants,
        channelNonce,
        challengeDuration,
        outcome: defaultOutcome,
        appDefinition: appDefinition ||"",
        appData: '0x',
      };

      const fixedPart = getFixedPart(state);

      const recoveredVP = shortenedToRecoveredVariableParts(turnNumToShortenedVariablePart);
      const {proof, candidate} = separateProofAndCandidate(recoveredVP);

      if (reason) {
        await expectRevert(() => Consensus.requireConsensus(fixedPart, proof, candidate));
      } else {
        await expectSucceedWithNoReturnValues(() =>
          Consensus.requireConsensus(fixedPart, proof, candidate)
        );
      }
    }
  );
});

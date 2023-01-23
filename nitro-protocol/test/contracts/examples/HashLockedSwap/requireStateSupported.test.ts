import {expectRevert} from '@statechannels/devtools';
import {Allocation, AllocationType} from '@statechannels/exit-format';
import {BigNumber, Contract, ethers, utils} from 'ethers';
import {it} from '@jest/globals';

const {HashZero} = ethers.constants;
import HashLockedSwapArtifact from '../../../../artifacts/contracts/examples/HashLockedSwap.sol/HashLockedSwap.json';
import {
  AssetOutcomeShortHand,
  bindSignaturesWithSignedByBitfield,
  Bytes32,
  getRandomNonce,
  signStates,
} from '../../../../src';
import {Outcome} from '../../../../src/contract/outcome';
import {
  getFixedPart,
  getVariablePart,
  separateProofAndCandidate,
  State,
} from '../../../../src/contract/state';
import {Bytes} from '../../../../src/contract/types';
import {
  getTestProvider,
  randomExternalDestination,
  setupContract,
  generateParticipants,
} from '../../../test-helpers';
import {expectSucceed} from '../../../expect-succeed';
import {replaceAddressesAndBigNumberify} from '../../../../src/helpers';

// Utilities
// TODO: move to a src file
interface HashLockedSwapData {
  h: Bytes32;
  preImage: Bytes;
}

function encodeHashLockedSwapData(data: HashLockedSwapData): string {
  return utils.defaultAbiCoder.encode(['tuple(bytes32 h, bytes preImage)'], [data]);
}
// *****

let hashTimeLock: Contract;

const addresses = {
  // Participants
  Sender: randomExternalDestination(),
  Receiver: randomExternalDestination(),
};

const provider = getTestProvider();

const nParticipants = 2;
const {wallets, participants} = generateParticipants(nParticipants);

const challengeDuration = 0x100;
const whoSignedWhat = [1, 0];

beforeAll(async () => {
  hashTimeLock = setupContract(provider, HashLockedSwapArtifact, process.env.HASH_LOCK_ADDRESS);
});

const preImage = '0xdeadbeef';
const conditionalPayment: HashLockedSwapData = {
  h: utils.sha256(preImage),
  // ^^^^ important field (SENDER)
  preImage: HashZero,
};

const correctPreImage: HashLockedSwapData = {
  preImage: preImage,
  // ^^^^ important field (RECEIVER)
  h: HashZero,
};

const incorrectPreImage: HashLockedSwapData = {
  preImage: '0xdeadc0de',
  // ^^^^ important field (RECEIVER)
  h: HashZero,
};

describe('requireStateSupported', () => {
  let channelNonce = getRandomNonce('HashLockedSwap');
  beforeEach(() => (channelNonce = BigNumber.from(channelNonce).add(1).toHexString()));
  it.each`
    isValid  | dataA                 | balancesA                   | turnNumB | dataB                | balancesB                   | description
    ${true}  | ${conditionalPayment} | ${{Sender: 1, Receiver: 0}} | ${4}     | ${correctPreImage}   | ${{Sender: 0, Receiver: 1}} | ${'Receiver unlocks the conditional payment'}
    ${false} | ${conditionalPayment} | ${{Sender: 1, Receiver: 0}} | ${4}     | ${incorrectPreImage} | ${{Sender: 0, Receiver: 1}} | ${'Receiver cannot unlock with incorrect preimage'}
  `(
    '$description',
    async ({
      isValid,
      dataA,
      balancesA,
      turnNumB,
      dataB,
      balancesB,
    }: {
      isValid: boolean;
      dataA: HashLockedSwapData;
      balancesA: AssetOutcomeShortHand;
      turnNumB: number;
      dataB: HashLockedSwapData;
      balancesB: AssetOutcomeShortHand;
    }) => {
      const turnNumA = turnNumB - 1;
      balancesA = replaceAddressesAndBigNumberify(balancesA, addresses) as AssetOutcomeShortHand;
      const allocationsA: Allocation[] = [];
      Object.keys(balancesA).forEach(key =>
        allocationsA.push({
          destination: key,
          amount: balancesA[key].toString(),
          allocationType: AllocationType.simple,
          metadata: '0x',
        })
      );
      const outcomeA: Outcome = [
        {
          asset: ethers.constants.AddressZero,
          allocations: allocationsA,
          metadata: '0x',
        },
      ];
      balancesB = replaceAddressesAndBigNumberify(balancesB, addresses) as AssetOutcomeShortHand;
      const allocationsB: Allocation[] = [];
      Object.keys(balancesB).forEach(key =>
        allocationsB.push({
          destination: key,
          amount: balancesB[key].toString(),
          allocationType: AllocationType.simple,
          metadata: '0x',
        })
      );
      const outcomeB: Outcome = [
        {asset: ethers.constants.AddressZero, allocations: allocationsB, metadata: '0x'},
      ];
      const states: State[] = [
        {
          turnNum: turnNumA,
          isFinal: false,
          channelNonce,
          participants,
          challengeDuration,
          outcome: outcomeA,
          appData: encodeHashLockedSwapData(dataA),
          appDefinition: hashTimeLock.address,
        },
        {
          turnNum: turnNumB,
          isFinal: false,
          channelNonce,
          participants,
          challengeDuration,
          outcome: outcomeB,
          appData: encodeHashLockedSwapData(dataB),
          appDefinition: hashTimeLock.address,
        },
      ];
      const fixedPart = getFixedPart(states[0]);
      const variableParts = states.map(s => getVariablePart(s));

      // Sign the states
      const signatures = await signStates(states, wallets, whoSignedWhat);
      const {proof, candidate} = separateProofAndCandidate(
        bindSignaturesWithSignedByBitfield(variableParts, signatures, whoSignedWhat)
      );

      if (isValid) {
        await expectSucceed(() => hashTimeLock.requireStateSupported(fixedPart, proof, candidate));
      } else {
        await expectRevert(
          () => hashTimeLock.requireStateSupported(fixedPart, proof, candidate),
          'incorrect preimage'
        );
      }
    }
  );
});

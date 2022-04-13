import {expectRevert} from '@statechannels/devtools';
import {Allocation, AllocationType} from '@statechannels/exit-format';
import {Contract, ethers, utils} from 'ethers';
import {it} from '@jest/globals'

const {HashZero} = ethers.constants;
import HashLockedSwapArtifact from '../../../../artifacts/contracts/examples/HashLockedSwap.sol/HashLockedSwap.json';
import {Bytes32} from '../../../../src';
import {Outcome} from '../../../../src/contract/outcome';
import {VariablePart} from '../../../../src/contract/state';
import {Bytes} from '../../../../src/contract/types';
import {
  getTestProvider,
  randomExternalDestination,
  replaceAddressesAndBigNumberify,
  setupContract,
  AssetOutcomeShortHand,
} from '../../../test-helpers';

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

const numParticipants = 2;
const addresses = {
  // Participants
  Sender: randomExternalDestination(),
  Receiver: randomExternalDestination(),
};
const provider = getTestProvider();

beforeAll(async () => {
  hashTimeLock = setupContract(provider, HashLockedSwapArtifact, process.env.HASH_LOCK_ADDRESS);
});

const preImage = '0xdeadbeef';
const conditionalPayment: HashLockedSwapData = {
  h: utils.sha256(preImage),
  // ^^^^ important field (SENDER)
  preImage: '0x',
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

describe('validTransition', () => {
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
      let turnNumA = turnNumB - 1;
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
      const variablePartA: VariablePart = {
        outcome: outcomeA,
        appData: encodeHashLockedSwapData(dataA),
        turnNum: turnNumA,
        isFinal: false,
      };
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
      const variablePartB: VariablePart = {
        outcome: outcomeB,
        appData: encodeHashLockedSwapData(dataB),
        turnNum: turnNumB,
        isFinal: false,
      };

      if (isValid) {
        const isValidFromCall = await hashTimeLock.validTransition(
          variablePartA,
          variablePartB,
          numParticipants
        );
        expect(isValidFromCall).toBe(true);
      } else {
        await expectRevert(
          () =>
            hashTimeLock.validTransition(variablePartA, variablePartB, numParticipants),
          'Incorrect preimage'
        );
      }
    }
  );
});

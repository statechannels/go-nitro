import {expectRevert} from '@statechannels/devtools';
import {Allocation, AllocationType} from '@statechannels/exit-format';
import {Contract, ethers, utils, Wallet} from 'ethers';
import {it} from '@jest/globals';

const {HashZero} = ethers.constants;
import HashLockedSwapArtifact from '../../../../artifacts/contracts/examples/HashLockedSwap.sol/HashLockedSwap.json';
import {bindSignaturesWithSignedByBitfield, Bytes32, Channel, signStates} from '../../../../src';
import {Outcome} from '../../../../src/contract/outcome';
import {getFixedPart, getVariablePart} from '../../../../src/contract/state';
import {Bytes} from '../../../../src/contract/types';
import {
  getTestProvider,
  randomExternalDestination,
  replaceAddressesAndBigNumberify,
  setupContract,
  AssetOutcomeShortHand,
  getRandomNonce,
  parseVariablePartEventResult,
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

const addresses = {
  // Participants
  Sender: randomExternalDestination(),
  Receiver: randomExternalDestination(),
};
const participants = ['', ''];
const wallets = new Array(2);
const provider = getTestProvider();
const chainId = process.env.CHAIN_NETWORK_ID;
const challengeDuration = 0x100;
const whoSignedWhat = [1, 0];

// Populate wallets and participants array
for (let i = 0; i < 2; i++) {
  wallets[i] = Wallet.createRandom();
  participants[i] = wallets[i].address;
}
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

describe('validTransition', () => {
  let channelNonce = getRandomNonce('HashLockedSwap');
  beforeEach(() => (channelNonce += 1));
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
      const channel: Channel = {chainId, channelNonce, participants};

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
      const states = [
        {
          turnNum: turnNumA,
          isFinal: false,
          channel,
          challengeDuration,
          outcome: outcomeA,
          appData: encodeHashLockedSwapData(dataA),
          appDefinition: hashTimeLock.address,
        },
        {
          turnNum: turnNumB,
          isFinal: false,
          channel,
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
      const recoveredVariableParts = bindSignaturesWithSignedByBitfield(
        variableParts,
        signatures,
        whoSignedWhat
      );

      if (isValid) {
        const latestSupportedState = await hashTimeLock.latestSupportedState(
          fixedPart,
          recoveredVariableParts
        );
        expect(parseVariablePartEventResult(latestSupportedState)).toStrictEqual(
          variableParts[variableParts.length - 1]
        );
      } else {
        await expectRevert(
          () => hashTimeLock.latestSupportedState(fixedPart, recoveredVariableParts),
          'incorrect preimage'
        );
      }
    }
  );
});

import {expectRevert} from '@statechannels/devtools';
import {Allocation, AllocationType} from '@statechannels/exit-format';
import {Contract, ethers, Wallet} from 'ethers';
import {it} from '@jest/globals';

const {HashZero} = ethers.constants;
import SingleAssetPaymentsArtifact from '../../../../artifacts/contracts/examples/SingleAssetPayments.sol/SingleAssetPayments.json';
import {encodeGuaranteeData, Outcome} from '../../../../src/contract/outcome';
import {getFixedPart, getVariablePart} from '../../../../src/contract/state';
import {
  getRandomNonce,
  getTestProvider,
  parseVariablePartEventResult,
  randomExternalDestination,
  replaceAddressesAndBigNumberify,
  setupContract,
} from '../../../test-helpers';
import {bindSignatures, Channel, signStates} from '../../../../src';
import {MOVER_SIGNED_EARLIER_STATE} from '../../../../src/contract/transaction-creators/revert-reasons';

const provider = getTestProvider();
let singleAssetPayments: Contract;

const addresses = {
  // Participants
  A: randomExternalDestination(),
  B: randomExternalDestination(),
  C: randomExternalDestination(),
};
const guaranteeDestinations = [addresses.A];

const participants = ['', '', ''];
const wallets = new Array(3);
const chainId = process.env.CHAIN_NETWORK_ID;
const challengeDuration = 0x100;

// Populate wallets and participants array
for (let i = 0; i < 3; i++) {
  wallets[i] = Wallet.createRandom();
  participants[i] = wallets[i].address;
}

beforeAll(async () => {
  singleAssetPayments = setupContract(
    provider,
    SingleAssetPaymentsArtifact,
    process.env.SINGLE_ASSET_PAYMENT_ADDRESS
  );
});

const whoSignedWhatA = [1, 0, 0];
const whoSignedWhatB = [0, 1, 0];
const whoSignedWhatC = [0, 0, 1];

const reason1 = MOVER_SIGNED_EARLIER_STATE;
const reason2 = 'not a simple allocation';
const reason3 = 'Total allocated cannot change';
const reason4 = 'outcome: Only one asset allowed';

describe('validTransition', () => {
  let channelNonce = getRandomNonce('HashLockedSwap');
  beforeEach(() => (channelNonce += 1));
  it.each`
    isValid  | numAssets | isAllocation      | balancesA             | turnNumB | balancesB             | whoSignedWhat     | reason       | description
    ${true}  | ${[1, 1]} | ${[true, true]}   | ${{A: 1, B: 1, C: 1}} | ${3}     | ${{A: 0, B: 2, C: 1}} | ${whoSignedWhatA} | ${undefined} | ${'A pays B 1 wei'}
    ${true}  | ${[1, 1]} | ${[true, true]}   | ${{A: 1, B: 1, C: 1}} | ${4}     | ${{A: 1, B: 0, C: 2}} | ${whoSignedWhatB} | ${undefined} | ${'B pays C 1 wei'}
    ${true}  | ${[1, 1]} | ${[true, true]}   | ${{A: 1, B: 1, C: 1}} | ${5}     | ${{A: 1, B: 2, C: 0}} | ${whoSignedWhatC} | ${undefined} | ${'C pays B 1 wei'}
    ${false} | ${[1, 1]} | ${[true, true]}   | ${{A: 1, B: 1, C: 1}} | ${5}     | ${{A: 0, B: 2, C: 1}} | ${whoSignedWhatA} | ${reason1}   | ${'A pays B 1 wei (not their move)'}
    ${false} | ${[1, 1]} | ${[false, false]} | ${{A: 1, B: 1, C: 1}} | ${3}     | ${{A: 0, B: 2, C: 1}} | ${whoSignedWhatA} | ${reason2}   | ${'Guarantee'}
    ${false} | ${[1, 1]} | ${[true, true]}   | ${{A: 1, B: 1, C: 1}} | ${3}     | ${{A: 1, B: 2, C: 1}} | ${whoSignedWhatA} | ${reason3}   | ${'Total amounts increase'}
    ${false} | ${[2, 2]} | ${[true, true]}   | ${{A: 1, B: 1, C: 1}} | ${3}     | ${{A: 2, B: 0, C: 1}} | ${whoSignedWhatA} | ${reason4}   | ${'More than one asset'}
  `(
    '$description',
    async ({
      isValid,
      isAllocation,
      numAssets,
      balancesA,
      turnNumB,
      balancesB,
      whoSignedWhat,
      reason,
    }: {
      isValid: boolean;
      isAllocation: boolean[];
      numAssets: number[];
      balancesA: any;
      turnNumB: number;
      balancesB: any;
      whoSignedWhat: number[];
      reason?: string;
    }) => {
      const channel: Channel = {chainId, channelNonce, participants};

      const turnNumA = turnNumB - 1;
      balancesA = replaceAddressesAndBigNumberify(balancesA, addresses);
      const allocationsA: Allocation[] = [];
      Object.keys(balancesA).forEach(key =>
        allocationsA.push({
          destination: key,
          amount: balancesA[key].toHexString(),
          allocationType: isAllocation[0] ? AllocationType.simple : AllocationType.guarantee,
          metadata: isAllocation[0] ? '0x' : encodeGuaranteeData(guaranteeDestinations),
        })
      );
      const outcomeA: Outcome = [
        {asset: ethers.constants.AddressZero, metadata: '0x', allocations: allocationsA},
      ];

      if (numAssets[0] === 2) {
        outcomeA.push(outcomeA[0]);
      }

      balancesB = replaceAddressesAndBigNumberify(balancesB, addresses);
      const allocationsB: Allocation[] = [];

      Object.keys(balancesB).forEach(key =>
        allocationsB.push({
          destination: key,
          amount: balancesB[key].toHexString(),
          allocationType: isAllocation[1] ? AllocationType.simple : AllocationType.guarantee,
          metadata: isAllocation[1] ? '0x' : encodeGuaranteeData(guaranteeDestinations),
        })
      );

      const outcomeB: Outcome = [
        {asset: ethers.constants.AddressZero, metadata: '0x', allocations: allocationsB},
      ];

      if (numAssets[1] === 2) {
        outcomeB.push(outcomeB[0]);
      }

      const states = [
        {
          turnNum: turnNumA,
          isFinal: false,
          channel,
          challengeDuration,
          outcome: outcomeA,
          appData: HashZero,
          appDefinition: singleAssetPayments.address,
        },
        {
          turnNum: turnNumB,
          isFinal: false,
          channel,
          challengeDuration,
          outcome: outcomeB,
          appData: HashZero,
          appDefinition: singleAssetPayments.address,
        },
      ];
      const fixedPart = getFixedPart(states[0]);
      const variableParts = states.map(s => getVariablePart(s));

      // Sign the states
      const signatures = await signStates(states, wallets, whoSignedWhat);
      const signedVariableParts = bindSignatures(variableParts, signatures, whoSignedWhat);

      if (isValid) {
        const latestSupportedState = await singleAssetPayments.latestSupportedState(
          fixedPart,
          signedVariableParts
        );
        expect(parseVariablePartEventResult(latestSupportedState)).toEqual(
          variableParts[variableParts.length - 1]
        );
      } else {
        await expectRevert(
          () => singleAssetPayments.latestSupportedState(fixedPart, signedVariableParts),
          reason
        );
      }
    }
  );
});

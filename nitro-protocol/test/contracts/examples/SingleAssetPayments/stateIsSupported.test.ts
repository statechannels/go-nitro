import {expectRevert} from '@statechannels/devtools';
import {Allocation, AllocationType} from '@statechannels/exit-format';
import {BigNumber, Contract, ethers} from 'ethers';
import {it} from '@jest/globals';

const {HashZero} = ethers.constants;
import SingleAssetPaymentsArtifact from '../../../../artifacts/contracts/examples/SingleAssetPayments.sol/SingleAssetPayments.json';
import {encodeGuaranteeData, Outcome} from '../../../../src/contract/outcome';
import {
  getFixedPart,
  getVariablePart,
  separateProofAndCandidate,
  State,
} from '../../../../src/contract/state';
import {
  generateParticipants,
  getTestProvider,
  randomExternalDestination,
  setupContract,
} from '../../../test-helpers';
import {
  AssetOutcomeShortHand,
  bindSignaturesWithSignedByBitfield,
  getRandomNonce,
  signStates,
} from '../../../../src';
import {INVALID_SIGNED_BY} from '../../../../src/contract/transaction-creators/revert-reasons';
import {expectSupportedState} from '../../../tx-expect-wrappers';
import {replaceAddressesAndBigNumberify} from '../../../../src/helpers';

const provider = getTestProvider();
let singleAssetPayments: Contract;

const addresses = {
  // Participants
  A: randomExternalDestination(),
  B: randomExternalDestination(),
};

const nParticipants = 2;
const {wallets, participants} = generateParticipants(nParticipants);

const challengeDuration = 0x100;
const guaranteeData = {left: addresses.A, right: addresses.B};

beforeAll(async () => {
  singleAssetPayments = setupContract(
    provider,
    SingleAssetPaymentsArtifact,
    process.env.SINGLE_ASSET_PAYMENTS_ADDRESS
  );
});

const whoSignedWhatA = [1, 0];
const whoSignedWhatB = [0, 1];

const reason1 = INVALID_SIGNED_BY;
const reason2 = 'not a simple allocation';
const reason3 = 'Total allocated cannot change';
const reason4 = 'outcome: Only one asset allowed';

describe('stateIsSupported', () => {
  let channelNonce = getRandomNonce('SingleAssetPayments');
  beforeEach(() => (channelNonce = BigNumber.from(channelNonce).add(1).toHexString()));
  it.each`
    numAssets | isAllocation      | balancesA       | turnNums  | balancesB       | whoSignedWhat     | reason       | description
    ${[1, 1]} | ${[true, true]}   | ${{A: 1, B: 1}} | ${[3, 4]} | ${{A: 0, B: 2}} | ${whoSignedWhatA} | ${undefined} | ${'A pays B 1 wei'}
    ${[1, 1]} | ${[true, true]}   | ${{A: 1, B: 1}} | ${[2, 3]} | ${{A: 2, B: 0}} | ${whoSignedWhatB} | ${undefined} | ${'B pays A 1 wei'}
    ${[1, 1]} | ${[true, true]}   | ${{A: 1, B: 1}} | ${[2, 3]} | ${{A: 0, B: 2}} | ${whoSignedWhatA} | ${reason1}   | ${'A pays B 1 wei (not their move)'}
    ${[1, 1]} | ${[false, false]} | ${{A: 1, B: 1}} | ${[3, 4]} | ${{A: 0, B: 2}} | ${whoSignedWhatA} | ${reason2}   | ${'Guarantee'}
    ${[1, 1]} | ${[true, true]}   | ${{A: 1, B: 1}} | ${[3, 4]} | ${{A: 1, B: 2}} | ${whoSignedWhatA} | ${reason3}   | ${'Total amounts increase'}
    ${[2, 2]} | ${[true, true]}   | ${{A: 1, B: 1}} | ${[3, 4]} | ${{A: 2, B: 0}} | ${whoSignedWhatA} | ${reason4}   | ${'More than one asset'}
  `(
    '$description',
    async ({
      isAllocation,
      numAssets,
      balancesA,
      turnNums,
      balancesB,
      whoSignedWhat,
      reason,
    }: {
      isAllocation: boolean[];
      numAssets: number[];
      balancesA: AssetOutcomeShortHand;
      turnNums: number[];
      balancesB: AssetOutcomeShortHand;
      whoSignedWhat: number[];
      reason?: string;
    }) => {
      balancesA = replaceAddressesAndBigNumberify(balancesA, addresses) as AssetOutcomeShortHand;
      const allocationsA: Allocation[] = [];
      Object.keys(balancesA).forEach(key =>
        allocationsA.push({
          destination: key,
          amount: balancesA[key].toString(),
          allocationType: isAllocation[0] ? AllocationType.simple : AllocationType.guarantee,
          metadata: isAllocation[0] ? '0x' : encodeGuaranteeData(guaranteeData),
        })
      );
      const outcomeA: Outcome = [
        {
          asset: ethers.constants.AddressZero,
          assetMetadata: {assetType: 0, metadata: '0x'},
          allocations: allocationsA,
        },
      ];

      if (numAssets[0] === 2) {
        outcomeA.push(outcomeA[0]);
      }

      balancesB = replaceAddressesAndBigNumberify(balancesB, addresses) as AssetOutcomeShortHand;
      const allocationsB: Allocation[] = [];

      Object.keys(balancesB).forEach(key =>
        allocationsB.push({
          destination: key,
          amount: balancesB[key].toString(),
          allocationType: isAllocation[1] ? AllocationType.simple : AllocationType.guarantee,
          metadata: isAllocation[1] ? '0x' : encodeGuaranteeData(guaranteeData),
        })
      );

      const outcomeB: Outcome = [
        {
          asset: ethers.constants.AddressZero,
          assetMetadata: {assetType: 0, metadata: '0x'},
          allocations: allocationsB,
        },
      ];

      if (numAssets[1] === 2) {
        outcomeB.push(outcomeB[0]);
      }

      const states: State[] = [
        {
          turnNum: turnNums[0],
          isFinal: false,
          channelNonce,
          participants,
          challengeDuration,
          outcome: outcomeA,
          appData: HashZero,
          appDefinition: singleAssetPayments.address,
        },
        {
          turnNum: turnNums[1],
          isFinal: false,
          channelNonce,
          participants,
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
      const recoveredVariableParts = bindSignaturesWithSignedByBitfield(
        variableParts,
        signatures,
        whoSignedWhat
      );

      const {proof, candidate} = separateProofAndCandidate(recoveredVariableParts);

      if (reason) {
        await expectRevert(
          () => singleAssetPayments.stateIsSupported(fixedPart, proof, candidate),
          reason
        );
      } else {
        await expectSupportedState(() =>
          singleAssetPayments.stateIsSupported(fixedPart, proof, candidate)
        );
      }
    }
  );
});

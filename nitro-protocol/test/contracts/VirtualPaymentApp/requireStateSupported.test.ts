import {expectRevert} from '@statechannels/devtools';
import {Contract, Wallet, ethers, BigNumber} from 'ethers';

import VirtualPaymentAppArtifact from '../../../artifacts/contracts/VirtualPaymentApp.sol/VirtualPaymentApp.json';
import {
  Channel,
  convertAddressToBytes32,
  encodeVoucherAmountAndSignature,
  getChannelId,
  signVoucher,
  Voucher,
} from '../../../src';
import {
  getFixedPart,
  getVariablePart,
  RecoveredVariablePart,
  State,
} from '../../../src/contract/state';
import {computeOutcome, getTestProvider, setupContract} from '../../test-helpers';
const {HashZero} = ethers.constants;

const provider = getTestProvider();
let virtualPaymentApp: Contract;

const participants = ['', '', ''];
const wallets = new Array(3);
const chainId = process.env.CHAIN_NETWORK_ID;
const challengeDuration = 0x100;
const MAGIC_ETH_ADDRESS = '0x0000000000000000000000000000000000000000';

// Populate wallets and participants array
for (let i = 0; i < 3; i++) {
  wallets[i] = Wallet.createRandom();
  participants[i] = wallets[i].address;
}
const channel: Channel = {chainId, channelNonce: 8, participants};

const baseState: State = {
  turnNum: 0,
  isFinal: false,
  channel,
  challengeDuration,
  outcome: [],
  appData: HashZero,
  appDefinition: process.env.VIRTUAL_PAYMENT_APP_ADDRESS,
};
const fixedPart = getFixedPart(baseState);
const channelId = getChannelId(fixedPart);

const alice = convertAddressToBytes32(participants[0]); // NOTE these desinations do not necessarily need to be related to participant addresses
const bob = convertAddressToBytes32(participants[2]);

beforeAll(async () => {
  virtualPaymentApp = setupContract(
    provider,
    VirtualPaymentAppArtifact,
    process.env.VIRTUAL_PAYMENT_APP_ADDRESS
  );
});

describe('requireStateSupported (lone candidate route)', () => {
  interface TestCase {
    turnNum: number;
    isFinal: boolean;
    reason?: string;
  }

  const testcases: TestCase[] = [
    {turnNum: 0, isFinal: false, reason: undefined},
    {turnNum: 1, isFinal: false, reason: undefined},
    {turnNum: 2, isFinal: false, reason: 'bad candidate turnNum'},
    {turnNum: 3, isFinal: false, reason: '!final; turnNum=3 && |proof|=0'},
    {turnNum: 3, isFinal: true, reason: undefined},
    {turnNum: 4, isFinal: false, reason: 'bad candidate turnNum'},
  ];

  testcases.map(async tc => {
    it(`${tc.reason ? 'reverts        ' : 'does not revert'} for unaninmous consensus on ${
      tc.isFinal ? 'final' : 'nonfinal'
    } state with turnNum ${tc.turnNum}`, async () => {
      const state: State = {
        ...baseState,
        turnNum: tc.turnNum,
        isFinal: tc.isFinal,
      };

      const variablePart = getVariablePart(state);

      const candidate: RecoveredVariablePart = {
        variablePart,
        signedBy: BigNumber.from(0b111).toHexString(),
      };

      if (tc.reason) {
        await expectRevert(
          () => virtualPaymentApp.requireStateSupported(fixedPart, [], candidate),
          tc.reason
        );
      } else {
        await virtualPaymentApp.requireStateSupported(fixedPart, [], candidate);
      }
    });
  });
});

describe('requireStateSupported (candidate plus single proof state route)', () => {
  interface TestCase {
    proofTurnNum: number;
    candidateTurnNum: number;
    unanimityOnProof: boolean;
    bobSignedCandidate: boolean;
    voucherForThisChannel: boolean;
    voucherSignedByAlice: boolean;
    aliceAdjustedCorrectly: boolean;
    bobAdjustedCorrectly: boolean;
    nativeAsset: boolean;
    multipleAssets: boolean;
    aliceUnderflow: boolean;
    reason?: string;
  }

  const vVR: TestCase = {
    // valid voucher redemption
    proofTurnNum: 1,
    candidateTurnNum: 2,
    unanimityOnProof: true,
    bobSignedCandidate: true,
    voucherForThisChannel: true,
    voucherSignedByAlice: true,
    aliceAdjustedCorrectly: true,
    bobAdjustedCorrectly: true,
    nativeAsset: true,
    multipleAssets: false,
    aliceUnderflow: false,
    reason: undefined,
  };
  const testcases: TestCase[] = [
    vVR,
    {...vVR, proofTurnNum: 0, reason: 'bad proof[0].turnNum; |proof|=1'},
    {...vVR, unanimityOnProof: false, reason: 'postfund !unanimous; |proof|=1'},
    {...vVR, bobSignedCandidate: false, reason: 'redemption not signed by Bob'},
    {...vVR, voucherSignedByAlice: false, reason: 'irrelevant voucher'},
    {...vVR, voucherForThisChannel: false, reason: 'irrelevant voucher'},
    {...vVR, aliceAdjustedCorrectly: false, reason: 'Alice not adjusted correctly'},
    {...vVR, bobAdjustedCorrectly: false, reason: 'Bob not adjusted correctly'},
    {...vVR, nativeAsset: false, reason: 'only native asset allowed'},
    {...vVR, multipleAssets: true, reason: 'only native asset allowed'},
    {...vVR, aliceUnderflow: true, reason: 'underflow'},
  ];

  testcases.map(async tc => {
    it(`${
      tc.reason ? 'reverts        ' : 'does not revert'
    } for a redemption transition with ${JSON.stringify(tc)}`, async () => {
      const proofState: State = {
        ...baseState,
        turnNum: tc.proofTurnNum,
        isFinal: false,
        outcome: computeOutcome({
          [MAGIC_ETH_ADDRESS]: {[alice]: 10, [bob]: 10},
        }),
      };

      // construct voucher with the (in)appropriate channelId
      const amount = tc.aliceUnderflow
        ? BigNumber.from(999_999_999_999).toHexString() // much larger than Alice's original balance
        : BigNumber.from(7).toHexString();

      const voucher: Voucher = {
        channelId: tc.voucherForThisChannel
          ? channelId
          : convertAddressToBytes32(MAGIC_ETH_ADDRESS),
        amount,
      };

      // make an (in)valid signature
      const signature = await signVoucher(voucher, wallets[0]);
      if (!tc.voucherSignedByAlice) signature.s = signature.r; // (conditionally) corrupt the signature

      // embed voucher into candidate state
      const encodedVoucherAmountAndSignature = encodeVoucherAmountAndSignature(amount, signature);
      const candidateState: State = {
        ...proofState,
        outcome: computeOutcome({
          [MAGIC_ETH_ADDRESS]: {
            [alice]: tc.aliceAdjustedCorrectly ? 3 : 2,
            [bob]: tc.bobAdjustedCorrectly ? 7 : 99,
          },
        }),
        turnNum: tc.candidateTurnNum,
        appData: encodedVoucherAmountAndSignature,
      };

      if (!tc.nativeAsset)
        candidateState.outcome[0].asset = process.env.VIRTUAL_PAYMENT_APP_ADDRESS;

      if (tc.multipleAssets) candidateState.outcome.push(candidateState.outcome[0]);

      // Sign the proof state (should be everyone)
      const proof: RecoveredVariablePart[] = [
        {
          variablePart: getVariablePart(proofState),
          signedBy: BigNumber.from(tc.unanimityOnProof ? 0b111 : 0b101).toHexString(),
        },
      ];

      // Sign the candidate state (should be just Bob)
      const candidate: RecoveredVariablePart = {
        variablePart: getVariablePart(candidateState),
        signedBy: BigNumber.from(tc.bobSignedCandidate ? 0b100 : 0b000).toHexString(), // 0b100 signed by Bob obly
      };

      if (tc.reason) {
        await expectRevert(
          () => virtualPaymentApp.requireStateSupported(fixedPart, proof, candidate),
          tc.reason
        );
      } else {
        await virtualPaymentApp.requireStateSupported(fixedPart, proof, candidate);
      }
    });
  });
});

describe('requireStateSupported (longer proof state route)', () => {
  it(`reverts for |support|>1`, async () => {
    const variablePart = getVariablePart(baseState);

    const candidate: RecoveredVariablePart = {
      variablePart,
      signedBy: BigNumber.from(0b111).toHexString(),
    };

    await expectRevert(
      () => virtualPaymentApp.requireStateSupported(fixedPart, [candidate, candidate], candidate),
      'bad proof length'
    );
  });
});

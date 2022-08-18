import {expectRevert} from '@statechannels/devtools';
import {Contract, Wallet, ethers, BigNumber, Signature} from 'ethers';

import VirtualPaymentAppArtifact from '../../../artifacts/contracts/VirtualPaymentApp.sol/VirtualPaymentApp.json';
import {
  bindSignaturesWithSignedByBitfield,
  Channel,
  convertAddressToBytes32,
  encodeVoucherAmountAndSignature,
  getChannelId,
  signState,
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

const alice = convertAddressToBytes32(participants[0]); // NOTE these desinations do not necessarily need to be related to participant addresses
const bob = convertAddressToBytes32(participants[2]);

const channel: Channel = {chainId, channelNonce: 8, participants};

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
    revertString?: string;
  }

  const testcases: TestCase[] = [
    {turnNum: 0, isFinal: false, revertString: undefined},
    {turnNum: 1, isFinal: false, revertString: undefined},
    {turnNum: 2, isFinal: false, revertString: 'bad candidate turnNum'},
    {turnNum: 3, isFinal: false, revertString: '!final; turnNum=3 && |proof|=0'},
    {turnNum: 3, isFinal: true, revertString: undefined},
    {turnNum: 4, isFinal: false, revertString: 'bad candidate turnNum'},
  ];

  testcases.map(async tc => {
    it(`${tc.revertString ? 'reverts        ' : 'does not revert'} for unaninmous consensus on ${
      tc.isFinal ? 'final' : 'nonfinal'
    } state with turnNum ${tc.turnNum}`, async () => {
      const state: State = {
        turnNum: tc.turnNum,
        isFinal: tc.isFinal,
        channel,
        challengeDuration,
        outcome: [],
        appData: HashZero,
        appDefinition: process.env.VIRTUAL_PAYMENT_APP_ADDRESS,
      };

      const fixedPart = getFixedPart(state);
      const variablePart = getVariablePart(state);

      // Sign the states
      const sigs = wallets.map((w: Wallet) => signState(state, w.privateKey).signature);

      const candidate: RecoveredVariablePart = bindSignaturesWithSignedByBitfield(
        [variablePart],
        sigs,
        [0, 0, 0]
      )[0];

      if (tc.revertString) {
        await expectRevert(
          () => virtualPaymentApp.requireStateSupported(fixedPart, [], candidate),
          tc.revertString
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
    revertString?: string;
  }

  const testcases: TestCase[] = [
    {
      proofTurnNum: 1,
      candidateTurnNum: 2,
      unanimityOnProof: true,
      bobSignedCandidate: true,
      voucherForThisChannel: true,
      voucherSignedByAlice: true,
      aliceAdjustedCorrectly: true,
      bobAdjustedCorrectly: true,
      revertString: undefined,
    }, // valid voucher redemption
    {
      proofTurnNum: 0, // incorrect
      candidateTurnNum: 2,
      unanimityOnProof: true,
      bobSignedCandidate: true,
      voucherForThisChannel: true,
      voucherSignedByAlice: true,
      aliceAdjustedCorrectly: true,
      bobAdjustedCorrectly: true,
      revertString: 'bad proof[0].turnNum; |proof|=1',
    },
    {
      proofTurnNum: 1,
      candidateTurnNum: 2,
      unanimityOnProof: false, // incorrect
      bobSignedCandidate: true,
      voucherForThisChannel: true,
      voucherSignedByAlice: true,
      aliceAdjustedCorrectly: true,
      bobAdjustedCorrectly: true,
      revertString: 'postfund !unanimous; |proof|=1',
    },
    {
      proofTurnNum: 1,
      candidateTurnNum: 2,
      unanimityOnProof: true,
      bobSignedCandidate: false, // incorrect
      voucherForThisChannel: true,
      voucherSignedByAlice: true,
      aliceAdjustedCorrectly: true,
      bobAdjustedCorrectly: true,
      revertString: 'redemption not signed by Bob',
    }, // valid voucher redemption
    {
      proofTurnNum: 1,
      candidateTurnNum: 2,
      unanimityOnProof: true,
      bobSignedCandidate: true,
      voucherForThisChannel: true,
      voucherSignedByAlice: false, // incorrect
      aliceAdjustedCorrectly: true,
      bobAdjustedCorrectly: true,
      revertString: 'irrelevant voucher',
    }, // valid voucher redemption
    {
      proofTurnNum: 1,
      candidateTurnNum: 2,
      unanimityOnProof: true,
      bobSignedCandidate: true,
      voucherForThisChannel: false, // incorrect
      voucherSignedByAlice: true,
      aliceAdjustedCorrectly: true,
      bobAdjustedCorrectly: true,
      revertString: 'irrelevant voucher',
    },
    {
      proofTurnNum: 1,
      candidateTurnNum: 2,
      unanimityOnProof: true,
      bobSignedCandidate: true,
      voucherForThisChannel: true,
      voucherSignedByAlice: true,
      aliceAdjustedCorrectly: false, // incorrect
      bobAdjustedCorrectly: true,
      revertString: 'Alice not adjusted correctly',
    },
    {
      proofTurnNum: 1,
      candidateTurnNum: 2,
      unanimityOnProof: true,
      bobSignedCandidate: true,
      voucherForThisChannel: true,
      voucherSignedByAlice: true,
      aliceAdjustedCorrectly: true,
      bobAdjustedCorrectly: false,
      revertString: 'Bob not adjusted correctly',
    },
  ];

  testcases.map(async tc => {
    it(`${
      tc.revertString ? 'reverts        ' : 'does not revert'
    } for a redemption transition with ${JSON.stringify(tc)}`, async () => {
      const proofState: State = {
        turnNum: tc.proofTurnNum,
        isFinal: false,
        channel,
        challengeDuration,
        outcome: computeOutcome({
          [MAGIC_ETH_ADDRESS]: {[alice]: 10, [bob]: 10},
        }),
        appData: HashZero,
        appDefinition: process.env.VIRTUAL_PAYMENT_APP_ADDRESS,
      };

      // construct voucher, sign it, and encode it into the appdata

      const fixedPart = getFixedPart(proofState);
      const channelId = getChannelId(fixedPart);
      const amount = BigNumber.from(7).toHexString();
      const voucher: Voucher = {
        channelId: tc.voucherForThisChannel
          ? channelId
          : convertAddressToBytes32(MAGIC_ETH_ADDRESS),
        amount,
      };

      const signature = await signVoucher(voucher, wallets[0]);
      if (!tc.voucherSignedByAlice) signature.s = signature.r; // corrupt the signature
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

      const candidateVariablePart = getVariablePart(candidateState);
      const proofVariablePart = getVariablePart(proofState);

      // Sign the proof state (should be everyone)
      const proof: RecoveredVariablePart[] = [
        {
          variablePart: proofVariablePart,
          signedBy: BigNumber.from(tc.unanimityOnProof ? 0b111 : 0b101).toHexString(),
        },
      ];

      // Sign the candidate state (just Bob)
      const candidate: RecoveredVariablePart = {
        variablePart: candidateVariablePart,
        signedBy: BigNumber.from(tc.bobSignedCandidate ? 0b100 : 0b000).toHexString(), // 0b100 signed by Bob obly
      };

      if (tc.revertString) {
        await expectRevert(
          () => virtualPaymentApp.requireStateSupported(fixedPart, proof, candidate),
          tc.revertString
        );
      } else {
        await virtualPaymentApp.requireStateSupported(fixedPart, proof, candidate);
      }
    });
  });
});

describe('requireStateSupported (longer proof state route)', () => {
  it(`reverts for |support|>1`, async () => {
    const state: State = {
      turnNum: 2,
      isFinal: false,
      channel,
      challengeDuration,
      outcome: [],
      appData: HashZero,
      appDefinition: process.env.VIRTUAL_PAYMENT_APP_ADDRESS,
    };

    const fixedPart = getFixedPart(state);
    const variablePart = getVariablePart(state);

    // Sign the states
    const sigs = wallets.map((w: Wallet) => signState(state, w.privateKey).signature);

    const candidate: RecoveredVariablePart = bindSignaturesWithSignedByBitfield(
      [variablePart],
      sigs,
      [0, 0, 0]
    )[0];

    await expectRevert(
      () => virtualPaymentApp.requireStateSupported(fixedPart, [candidate, candidate], candidate),
      'bad proof length'
    );
  });
});

// TODO we do not actually need to generate any signatures in tests like this. All we need is the signedBy bitfield declaration.

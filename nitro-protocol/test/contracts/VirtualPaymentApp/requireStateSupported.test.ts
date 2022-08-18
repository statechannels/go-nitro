import {expectRevert} from '@statechannels/devtools';
import {Contract, Wallet, ethers, BigNumber, Signature} from 'ethers';
import {defaultAbiCoder, keccak256, ParamType} from 'ethers/lib/utils';

import VirtualPaymentAppArtifact from '../../../artifacts/contracts/VirtualPaymentApp.sol/VirtualPaymentApp.json';
import {
  bindSignaturesWithSignedByBitfield,
  Channel,
  convertAddressToBytes32,
  getChannelId,
  sign,
  signState,
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
      interface SignedVoucher {
        channelId: string;
        amount: string;
        signature: Signature;
      }
      const fixedPart = getFixedPart(proofState);

      const channelId = getChannelId(fixedPart);
      const amount = BigNumber.from(7).toHexString();

      const voucherTy = {
        type: 'tuple',
        components: [
          {name: 'channelId', type: 'bytes32'},
          {
            name: 'amount',
            type: 'uint256',
          },
        ],
      } as ParamType;

      const voucher: SignedVoucher = {
        channelId,
        amount,
        signature: await sign(
          wallets[0],
          keccak256(defaultAbiCoder.encode([voucherTy], [{channelId, amount}]))
        ),
      };

      const signedVoucherTy = {
        type: 'tuple',
        components: [
          {name: 'channelId', type: 'bytes32'},
          {
            name: 'amount',
            type: 'uint256',
          },
          {
            type: 'tuple',
            name: 'signature',
            components: [
              {name: 'v', type: 'uint8'},
              {name: 'r', type: 'bytes32'},
              {name: 's', type: 'bytes32'},
            ],
          } as ParamType,
        ],
      } as ParamType;
      const encodedVoucher = defaultAbiCoder.encode([signedVoucherTy], [voucher]);

      const candidateState: State = {
        ...proofState,
        outcome: computeOutcome({
          [MAGIC_ETH_ADDRESS]: {[alice]: 3, [bob]: 7},
        }),
        turnNum: tc.candidateTurnNum,
        appData: encodedVoucher,
      };

      const candidateVariablePart = getVariablePart(candidateState);
      const proofVariablePart = getVariablePart(proofState);

      // Sign the proof state (everyone)
      const proofSigs = wallets.map((w: Wallet) => signState(proofState, w.privateKey).signature);
      const proof: RecoveredVariablePart[] = bindSignaturesWithSignedByBitfield(
        [proofVariablePart],
        proofSigs,
        [0, 0, 0]
      );

      // Sign the candidate state (just Bob)
      const candidate: RecoveredVariablePart = {
        variablePart: candidateVariablePart,
        signedBy: BigNumber.from(0b100).toHexString(), // 0b100 signed by Bob obly
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

// TODO
// describe('requireStateSupported (longer proof state route)', () => {});

// TODO we do not actually need to generate any signatures in tests like this. All we need is the signedBy bitfield declaration.

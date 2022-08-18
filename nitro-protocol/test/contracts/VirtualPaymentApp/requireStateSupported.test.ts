import {expectRevert} from '@statechannels/devtools';
import {Contract, Wallet, ethers, BigNumber} from 'ethers';

import VirtualPaymentAppArtifact from '../../../artifacts/contracts/VirtualPaymentApp.sol/VirtualPaymentApp.json';
import {bindSignaturesWithSignedByBitfield, Channel, signState} from '../../../src';
import {
  getFixedPart,
  getVariablePart,
  RecoveredVariablePart,
  State,
} from '../../../src/contract/state';
import {getTestProvider, setupContract} from '../../test-helpers';
const {HashZero} = ethers.constants;

const provider = getTestProvider();
let virtualPaymentApp: Contract;

const participants = ['', '', ''];
const wallets = new Array(3);
const chainId = process.env.CHAIN_NETWORK_ID;
const challengeDuration = 0x100;

// Populate wallets and participants array
for (let i = 0; i < 3; i++) {
  wallets[i] = Wallet.createRandom();
  participants[i] = wallets[i].address;
}

const channel: Channel = {chainId, channelNonce: 8, participants};

beforeAll(async () => {
  virtualPaymentApp = setupContract(
    provider,
    VirtualPaymentAppArtifact,
    process.env.VIRTUAL_PAYMENT_APP_ADDRESS
  );
});

describe('requireStateSupported (unanimous consensus route)', () => {
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

import {expectRevert} from '@statechannels/devtools';
import {Contract, Wallet, constants} from 'ethers';

import TESTShortcuttingTurnTakingArtifact from '../../../../artifacts/contracts/test/TESTShortcuttingTurnTaking.sol/TESTShortcuttingTurnTaking.json';
import {
  bindSignatures,
  Channel,
  getFixedPart,
  getVariablePart,
  signState,
  State,
} from '../../../../src';
const {HashZero} = constants;
import {TESTShortcuttingTurnTaking} from '../../../../typechain-types';
import {getTestProvider, setupContract} from '../../../test-helpers';
const provider = getTestProvider();
let TESTShortcuttingTurnTaking: Contract & TESTShortcuttingTurnTaking;

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

const state: State = {
  turnNum: 6,
  isFinal: false,
  channel,
  challengeDuration,
  outcome: [],
  appData: HashZero,
  appDefinition: process.env.CONSENSUS_APP_ADDRESS as string,
};

beforeAll(async () => {
  TESTShortcuttingTurnTaking = setupContract(
    provider,
    TESTShortcuttingTurnTakingArtifact,
    process.env.TEST_SHORTCUTTING_TURN_TAKING_ADDRESS
  ) as Contract & TESTShortcuttingTurnTaking;
});

describe('requireValidTurnTaking', () => {
  it('handles round robin signing', async () => {
    interface TestCase {
      turnNums: number[];
      allowed: boolean;
    }
    const testcases: TestCase[] = [
      {turnNums: [0, 1, 2], allowed: true},
      {turnNums: [3, 4, 5], allowed: true},
      {turnNums: [6, 7, 8], allowed: true},
      {turnNums: [9, 10, 11], allowed: true},
      {turnNums: [5, 6, 7], allowed: false},
      {turnNums: [1, 2, 3], allowed: false},
      {turnNums: [0, 1, 3], allowed: false},
    ];

    await Promise.all(
      testcases.map(async (testcase: TestCase) => {
        const states = [
          {...state, turnNum: testcase.turnNums[0]},
          {...state, turnNum: testcase.turnNums[1]},
          {...state, turnNum: testcase.turnNums[2]},
        ];
        const variableParts = states.map(getVariablePart);
        const fixedPart = getFixedPart(state);

        // Sign the states
        const sigs = wallets.map(
          (w: Wallet, idx: number) => signState(states[idx], w.privateKey).signature
        );
        const signedVariableParts = bindSignatures(variableParts, sigs, [0, 1, 2]);

        if (testcase.allowed == true) {
          const result = await TESTShortcuttingTurnTaking.requireValidTurnTaking(
            fixedPart,
            signedVariableParts
          );
          expect(result).toBe(true);
        } else {
          await expectRevert(async () =>
            TESTShortcuttingTurnTaking.requireValidTurnTaking(fixedPart, signedVariableParts)
          );
        }
      })
    );

    // it('permits consensus proofs', async () => {});
  });
});

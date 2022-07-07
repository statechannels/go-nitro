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

const states = [state, {...state, turnNum: 7}, {...state, turnNum: 8}];
const variableParts = states.map(getVariablePart);
const fixedPart = getFixedPart(state);

// Sign the states
const sigs = wallets.map(
  (w: Wallet, idx: number) => signState(states[idx], w.privateKey).signature
);

const signedVariableParts = bindSignatures(variableParts, sigs, [0, 1, 2]);

beforeAll(async () => {
  TESTShortcuttingTurnTaking = setupContract(
    provider,
    TESTShortcuttingTurnTakingArtifact,
    process.env.TEST_SHORTCUTTING_TURN_TAKING_ADDRESS
  ) as Contract & TESTShortcuttingTurnTaking;
});

describe('_recoverSigner', () => {
  it('permits round robin signing', async () => {
    const result = await TESTShortcuttingTurnTaking.requireValidTurnTaking(
      fixedPart,
      signedVariableParts
    );
    expect(result).toBe(true);
  });
});

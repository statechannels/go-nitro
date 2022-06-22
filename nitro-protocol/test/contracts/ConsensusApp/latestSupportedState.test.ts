import {expectRevert} from '@statechannels/devtools';
import {Contract, Wallet, ethers, BigNumber} from 'ethers';

import ConsensusAppArtifact from '../../../artifacts/contracts/ConsensusApp.sol/ConsensusApp.json';
import {Channel, signState} from '../../../src';
import {
  getFixedPart,
  getVariablePart,
  SignedVariablePart,
  State,
} from '../../../src/contract/state';
import {
  deepEquals,
  getTestProvider,
  parseVariablePartEventResult,
  setupContract,
} from '../../test-helpers';
const {HashZero} = ethers.constants;

const provider = getTestProvider();
let consensusApp: Contract;

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
  consensusApp = setupContract(
    provider,
    ConsensusAppArtifact,
    process.env.CONSENSUS_APP_ADDRESS as string
  );
});

const state: State = {
  turnNum: 5,
  isFinal: false,
  channel,
  challengeDuration,
  outcome: [],
  appData: HashZero,
  appDefinition: process.env.CONSENSUS_APP_ADDRESS as string,
};

const fixedPart = getFixedPart(state);
const variablePart = getVariablePart(state);

// Sign the states
const sigs = wallets.map((w: Wallet) => signState(state, w.privateKey).signature);

describe('latestSupportedState', () => {
  it('A single state signed by everyone is considered supported', async () => {
    expect.assertions(1);

    const signedVariablePart: SignedVariablePart = {
      variablePart,
      signedBy: BigNumber.from(0b111).toHexString(),
      sigs,
    };

    const latestSupportedState = await consensusApp.latestSupportedState(fixedPart, [
      signedVariablePart,
    ]);
    expect(deepEquals(parseVariablePartEventResult(latestSupportedState), variablePart)).toBe(true);
  });

  it('Submitting more than one state does NOT constitute a support proof', async () => {
    expect.assertions(1);

    const signedVariablePart: SignedVariablePart = {
      variablePart,
      signedBy: BigNumber.from(0b011).toHexString(),
      sigs: sigs.slice(0, 1),
    };

    await expectRevert(() =>
      consensusApp.latestSupportedState(fixedPart, [signedVariablePart, signedVariablePart])
    );
  });

  it('A single state signed by less than everyone is NOT considered supported', async () => {
    expect.assertions(1);

    const signedVariablePart: SignedVariablePart = {
      variablePart,
      signedBy: BigNumber.from(0b011).toHexString(),
      sigs: sigs.slice(0, 1),
    };

    await expectRevert(() => consensusApp.latestSupportedState(fixedPart, [signedVariablePart]));
  });

  it('A single state signed by less than everyone is NOT considered supported, even if we claim it is signed by everyone', async () => {
    expect.assertions(1);

    const signedVariablePart: SignedVariablePart = {
      variablePart,
      signedBy: BigNumber.from(0b111).toHexString(),
      sigs: sigs.slice(0, 1),
    };

    await expectRevert(() => consensusApp.latestSupportedState(fixedPart, [signedVariablePart]));
  });
});

import {expectRevert} from '@statechannels/devtools';
import {Contract, Wallet, ethers, BigNumber} from 'ethers';

import ConsensusAppArtifact from '../../../artifacts/contracts/ConsensusApp.sol/ConsensusApp.json';
import {bindSignaturesWithSignedByBitfield, Channel, signState} from '../../../src';
import {
  getFixedPart,
  getVariablePart,
  RecoveredVariablePart,
  State,
} from '../../../src/contract/state';
import {generateParticipants, getTestProvider, setupContract} from '../../test-helpers';
const {HashZero} = ethers.constants;

const provider = getTestProvider();
let consensusApp: Contract;

const nParticipants = 3;
const {wallets, participants} = generateParticipants(nParticipants);
const chainId = process.env.CHAIN_NETWORK_ID;
const challengeDuration = 0x100;

const channel: Channel = {chainId, channelNonce: 8, participants};

beforeAll(async () => {
  consensusApp = setupContract(provider, ConsensusAppArtifact, process.env.CONSENSUS_APP_ADDRESS);
});

const state: State = {
  turnNum: 5,
  isFinal: false,
  channel,
  challengeDuration,
  outcome: [],
  appData: HashZero,
  appDefinition: process.env.CONSENSUS_APP_ADDRESS,
};

const fixedPart = getFixedPart(state);
const variablePart = getVariablePart(state);

// Sign the states
const sigs = wallets.map((w: Wallet) => signState(state, w.privateKey).signature);

describe('requireStateSupported', () => {
  const candidate: RecoveredVariablePart = bindSignaturesWithSignedByBitfield(
    [variablePart],
    sigs,
    [0, 0, 0]
  )[0];
  it('A single state signed by everyone is considered supported', async () => {
    expect.assertions(1);
    const txResult = await consensusApp.requireStateSupported(fixedPart, [], candidate);

    // As 'requireStateSupported' method is constant (view or pure), if it succeedes, it returns an object/array with returned values
    // which in this case should be empty
    expect(txResult.length).toBe(0);
  });

  it('Submitting more than one state does NOT constitute a support proof', async () => {
    expect.assertions(1);

    await expectRevert(() => consensusApp.requireStateSupported(fixedPart, [candidate], candidate));
  });

  it('A single state signed by less than everyone is NOT considered supported', async () => {
    expect.assertions(1);

    const candidate: RecoveredVariablePart = {
      variablePart,
      signedBy: BigNumber.from(0b011).toHexString(),
    };
    await expectRevert(() => consensusApp.requireStateSupported(fixedPart, [], candidate));
  });
});

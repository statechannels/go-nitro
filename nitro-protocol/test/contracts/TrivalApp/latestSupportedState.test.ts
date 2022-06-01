import {Contract, Wallet, ethers, utils} from 'ethers';

import TrivialAppArtifact from '../../../artifacts/contracts/TrivialApp.sol/TrivialApp.json';
import {Channel} from '../../../src/contract/channel';
import {
  FixedPart,
  getFixedPart,
  getVariablePart,
  SignedVariablePart,
  State,
  VariablePart,
} from '../../../src/contract/state';
import {getRandomNonce, getTestProvider, setupContract} from '../../test-helpers';

const provider = getTestProvider();
let trivialApp: Contract;

function computeSaltedHash(salt: string, num: number) {
  return utils.solidityKeccak256(['bytes32', 'uint256'], [salt, num]);
}

function getRandomSignedVariablePart(): SignedVariablePart {
  const randomNum = Math.floor(Math.random() * 100);
  const salt = ethers.constants.MaxUint256.toHexString();
  const hash = computeSaltedHash(salt, randomNum);

  const signedVariablePart: SignedVariablePart = {
    variablePart: {
      outcome: [],
      appData: hash,
      turnNum: 1,
      isFinal: false,
    },
    sigs: [],
    signedBy: '0',
  };
  return signedVariablePart;
}

function getMockedFixedPart(): FixedPart {
  const fixedPart: FixedPart = {
    chainId: '',
    participants: [],
    channelNonce: 0,
    appDefinition: '',
    challengeDuration: 0,
  };
  return fixedPart;
}

function mockSigs(vp: VariablePart): SignedVariablePart {
  return {
    variablePart: vp,
    sigs: [],
    signedBy: '0',
  };
}

beforeAll(async () => {
  trivialApp = setupContract(provider, TrivialAppArtifact, process.env.TRIVIAL_APP_ADDRESS);
});

describe('latestSupportedState', () => {
  it('Transitions between random VariableParts are valid', async () => {
    expect.assertions(5);
    for (let i = 0; i < 5; i++) {
      const from: SignedVariablePart = getRandomSignedVariablePart();
      const to: SignedVariablePart = getRandomSignedVariablePart();
      const latestSupportedState = await trivialApp.latestSupportedState(getMockedFixedPart(), [
        from,
        to,
      ]);
      expect(latestSupportedState).toBe(to);
    }
  });

  it('Transitions between States with mocked-up data are valid', async () => {
    const channel: Channel = {
      participants: [Wallet.createRandom().address, Wallet.createRandom().address],
      chainId: process.env.CHAIN_NETWORK_ID,
      channelNonce: getRandomNonce('trivialApp'),
    };
    const fromState: State = {
      channel,
      outcome: [],
      turnNum: 1,
      isFinal: false,
      challengeDuration: 0x0,
      appDefinition: trivialApp.address,
      appData: '0x00',
    };
    const toState: State = {...fromState, turnNum: 2};

    const from: SignedVariablePart = mockSigs(getVariablePart(fromState));
    const to: SignedVariablePart = mockSigs(getVariablePart(toState));

    const latestSupportedState = await trivialApp.latestSupportedState(getFixedPart(fromState), [
      from,
      to,
    ]);
    expect(latestSupportedState).toBe(to);
  });
});

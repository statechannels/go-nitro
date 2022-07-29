import {Contract, Wallet, ethers, utils, BigNumber} from 'ethers';

import TrivialAppArtifact from '../../../artifacts/contracts/TrivialApp.sol/TrivialApp.json';
import {NITRO_MAX_GAS} from '../../../src';
import {Channel} from '../../../src/contract/channel';
import {
  FixedPart,
  getFixedPart,
  getVariablePart,
  RecoveredVariablePart,
  State,
  VariablePart,
} from '../../../src/contract/state';
import {getRandomNonce, getTestProvider, setupContract} from '../../test-helpers';

const provider = getTestProvider();
let trivialApp: Contract;

function computeSaltedHash(salt: string, num: number) {
  return utils.solidityKeccak256(['bytes32', 'uint256'], [salt, num]);
}

function getRandomRecoveredVariablePart(): RecoveredVariablePart {
  const randomNum = Math.floor(Math.random() * 100);
  const salt = ethers.constants.MaxUint256.toHexString();
  const hash = computeSaltedHash(salt, randomNum);

  const recoveredVariablePart: RecoveredVariablePart = {
    variablePart: {
      outcome: [],
      appData: hash,
      turnNum: 1,
      isFinal: false,
    },
    signedBy: '0',
  };
  return recoveredVariablePart;
}

function getMockedFixedPart(): FixedPart {
  const fixedPart: FixedPart = {
    chainId: process.env.CHAIN_NETWORK_ID,
    participants: [Wallet.createRandom().address, Wallet.createRandom().address],
    channelNonce: 0,
    appDefinition: trivialApp.address,
    challengeDuration: 0,
  };
  return fixedPart;
}

function mockSigs(vp: VariablePart): RecoveredVariablePart {
  return {
    variablePart: vp,
    signedBy: '0',
  };
}

beforeAll(async () => {
  trivialApp = setupContract(provider, TrivialAppArtifact, process.env.TRIVIAL_APP_ADDRESS);
});

describe('requireStateSupported', () => {
  it('Transitions between random VariableParts are valid', async () => {
    expect.assertions(5);
    for (let i = 0; i < 5; i++) {
      const from: RecoveredVariablePart = getRandomRecoveredVariablePart();
      const to: RecoveredVariablePart = getRandomRecoveredVariablePart();
      const requireStateSupported = await trivialApp.requireStateSupported(
        getMockedFixedPart(),
        [from],
        to
      );
      expect(
        BigNumber.from((await requireStateSupported.wait()).gasUsed).lt(
          BigNumber.from(NITRO_MAX_GAS)
        )
      ).toBe(true);
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

    const from: RecoveredVariablePart = mockSigs(getVariablePart(fromState));
    const to: RecoveredVariablePart = mockSigs(getVariablePart(toState));

    const tx = trivialApp.requireStateSupported(getFixedPart(fromState), [from], to);

    expect(
      BigNumber.from((await (await tx).wait()).gasUsed).lt(BigNumber.from(NITRO_MAX_GAS))
    ).toBe(true);
  });
});

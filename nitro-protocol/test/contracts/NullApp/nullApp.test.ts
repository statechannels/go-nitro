import {expectRevert} from '@statechannels/devtools';
import {Contract, Wallet, ethers} from 'ethers';

import {getTestProvider, setupContract} from '../../test-helpers';
import NitroAdjudicatorArtifact from '../../../artifacts/contracts/test/TESTNitroAdjudicator.sol/TESTNitroAdjudicator.json';
import {getVariablePart, State, getFixedPart, getRandomNonce} from '../../../src';
import {FixedPart, SignedVariablePart} from '../../../src/contract/state';

const provider = getTestProvider();

let NitroAdjudicator: Contract;
beforeAll(async () => {
  NitroAdjudicator = setupContract(
    provider,
    NitroAdjudicatorArtifact,
    process.env.TEST_NITRO_ADJUDICATOR_ADDRESS
  );
});

describe('null app', () => {
  it('should revert when stateIsSupported is called', async () => {
    const fromState: State = {
      participants: [Wallet.createRandom().address, Wallet.createRandom().address],
      channelNonce: getRandomNonce('nullApp'),
      outcome: [],
      turnNum: 1,
      isFinal: false,
      challengeDuration: 0x0,
      appDefinition: ethers.constants.AddressZero,
      appData: '0x00',
    };
    const toState: State = {...fromState, turnNum: 2};

    const fixedPart: FixedPart = getFixedPart(fromState);
    const from: SignedVariablePart = {
      variablePart: getVariablePart(fromState),
      sigs: [],
    };
    const to: SignedVariablePart = {
      variablePart: getVariablePart(toState),
      sigs: [],
    };

    await expectRevert(async () => {
      await NitroAdjudicator.stateIsSupported(fixedPart, [from], to);
    }, 'VM Exception while processing transaction: revert');
  });
});

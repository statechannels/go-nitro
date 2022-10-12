import {ethers, Contract, Wallet, BigNumber, constants} from 'ethers';
import {it} from '@jest/globals';

import {getRandomNonce} from '../../test-helpers';
import {TESTNitroAdjudicator} from '../../../typechain-types';
// eslint-disable-next-line import/order
import TESTNitroAdjudicatorArtifact from '../../../artifacts/contracts/test/TESTNitroAdjudicator.sol/TESTNitroAdjudicator.json';

const provider = new ethers.providers.JsonRpcProvider('https://wallaby.node.glif.io/rpc/v0');

let wallet = new Wallet('9182b5bf5b9c966e001934ebaf008f65516290cef6e3069d11e718cbd4336aae');
wallet = wallet.connect(provider);

const testNitroAdjudicator = new ethers.Contract(
  '0xFF00000000000000000000000000000000000485',
  TESTNitroAdjudicatorArtifact.abi,
  wallet
) as unknown as TESTNitroAdjudicator & Contract;

describe('deposit', () => {
  let channelNonce = getRandomNonce('deposit');
  afterEach(() => {
    channelNonce = BigNumber.from(channelNonce).add(1).toHexString();
  });
  it('submits a deposit transaction', async () => {
    const txPromise = testNitroAdjudicator.deposit(
      constants.AddressZero, // asset
      constants.HashZero, // destination
      BigNumber.from(0), // expectedHeld
      BigNumber.from(1), // amount
      {
        value: BigNumber.from(1),
      }
    );
    console.log(await txPromise);
  });
});

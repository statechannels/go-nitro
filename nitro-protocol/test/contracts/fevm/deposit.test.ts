import {expectRevert} from '@statechannels/devtools';
import {ethers, Contract, Wallet, BigNumber, utils} from 'ethers';
import {it} from '@jest/globals';
const {AddressZero} = ethers.constants;

import {getChannelId} from '../../../src/contract/channel';
import {getRandomNonce} from '../../test-helpers';
import {TESTNitroAdjudicator} from '../../../typechain-types';
// eslint-disable-next-line import/order
import TESTNitroAdjudicatorArtifact from '../../../artifacts/contracts/test/TESTNitroAdjudicator.sol/TESTNitroAdjudicator.json';
import {MAGIC_ADDRESS_INDICATING_ETH} from '../../../src/transactions';

const provider = new ethers.providers.JsonRpcProvider('https://wallaby.node.glif.io/rpc/v0');
const wallet = new Wallet('9182b5bf5b9c966e001934ebaf008f65516290cef6e3069d11e718cbd4336aae');

const testNitroAdjudicator = new ethers.Contract(
  't2ncif6kqekarww5vartvfwa6xxwjhg3l34grgmoi',
  TESTNitroAdjudicatorArtifact.abi,
  wallet
) as unknown as TESTNitroAdjudicator & Contract;

// const token = setupContract(
//   provider,
//   TokenArtifact,
//   process.env.TEST_TOKEN_ADDRESS
// ) as unknown as Token & Contract;

const chainId = '31415';
const participants: string[] = [];
const challengeDuration = 0x1000;

const ETH = MAGIC_ADDRESS_INDICATING_ETH;
//const ERC20 = token.address;

// Populate destinations array
for (let i = 0; i < 3; i++) {
  participants[i] = Wallet.createRandom({extraEntropy: utils.id('erc20-deposit-test')}).address;
}

const description0 = 'Deposits ETH (msg.value = amount , expectedHeld = 0)';

describe('deposit', () => {
  let channelNonce = getRandomNonce('deposit');
  afterEach(() => {
    channelNonce = BigNumber.from(channelNonce).add(1).toHexString();
  });
  it.each`
    description     | asset  | held | expectedHeld | amount | heldAfter | reasonString
    ${description0} | ${ETH} | ${0} | ${0}         | ${1}   | ${1}      | ${undefined}
  `('$description', async ({asset, held, expectedHeld, amount, reasonString, heldAfter}) => {
    held = BigNumber.from(held);
    expectedHeld = BigNumber.from(expectedHeld);
    amount = BigNumber.from(amount);
    heldAfter = BigNumber.from(heldAfter);

    const destination = getChannelId({
      chainId,
      channelNonce,
      participants,
      appDefinition: AddressZero,
      challengeDuration,
    });

    // if (asset === ERC20) {
    //   // Check msg.sender has enough tokens
    //   const balance = await token.balanceOf(signer0Address);
    //   await expect(balance.gte(held.add(amount))).toBe(true);

    //   // Increase allowance
    //   await (await token.increaseAllowance(testNitroAdjudicator.address, held.add(amount))).wait(); // Approve enough for setup and main test

    //   // Check allowance updated
    //   const allowance = BigNumber.from(
    //     await token.allowance(signer0Address, testNitroAdjudicator.address)
    //   );
    //   expect(allowance.sub(amount).sub(held).gte(0)).toBe(true);
    // }

    if (held > 0) {
      // Set holdings by depositing in the 'safest' way
      const {events} = await (
        await testNitroAdjudicator.deposit(asset, destination, 0, held, {
          value: asset === ETH ? held : 0,
        })
      ).wait();

      expect(events).not.toBe(undefined);
      if (events === undefined) {
        return;
      }

      expect(await testNitroAdjudicator.holdings(asset, destination)).toEqual(held);
      // if (asset === ERC20) {
      //   const {data: amountTransferred} = getTransferEvent(events);
      //   expect(held.eq(amountTransferred)).toBe(true);
      // }
    }

    const balanceBefore = await getBalance(asset, wallet.address);

    const tx = testNitroAdjudicator.deposit(asset, destination, expectedHeld, amount, {
      value: asset === ETH ? amount : 0,
    });

    if (reasonString) {
      await expectRevert(() => tx, reasonString);
    } else {
      const {events} = await (await tx).wait();
      expect(events).not.toBe(undefined);
      if (events === undefined) {
        return;
      }

      // if (asset == ERC20) {
      //   const amountTransferred = BigNumber.from(getTransferEvent(events).data);
      //   expect(heldAfter.sub(held).eq(amountTransferred)).toBe(true);
      //   const balanceAfter = await getBalance(asset, signer0Address);
      //   expect(balanceAfter.eq(balanceBefore.sub(heldAfter.sub(held)))).toBe(true);
      // }

      const allocatedAmount = await testNitroAdjudicator.holdings(asset, destination);
      await expect(allocatedAmount).toEqual(heldAfter);
    }
  });
});

async function getBalance(asset: string, address: string) {
  return BigNumber.from(await provider.getBalance(address));
}

import {expectRevert} from '@statechannels/devtools';
import {ethers, Contract, Wallet, BigNumber, utils} from 'ethers';
import {it} from '@jest/globals';
const {AddressZero} = ethers.constants;

import TokenArtifact from '../../../artifacts/contracts/Token.sol/Token.json';
import BadTokenArtifact from '../../../artifacts/contracts/test/BadToken.sol/BadToken.json';
import {getChannelId} from '../../../src/contract/channel';
import {getCountingAppContractAddress, getTestProvider, setupContract} from '../../test-helpers';
import {Token, BadToken, TESTNitroAdjudicator} from '../../../typechain-types';
// eslint-disable-next-line import/order
import TESTNitroAdjudicatorArtifact from '../../../artifacts/contracts/test/TESTNitroAdjudicator.sol/TESTNitroAdjudicator.json';
import {MAGIC_ADDRESS_INDICATING_ETH} from '../../../src/transactions';
import {getRandomNonce} from '../../../src/helpers';
const provider = getTestProvider();
const testNitroAdjudicator = setupContract(
  provider,
  TESTNitroAdjudicatorArtifact,
  process.env.TEST_NITRO_ADJUDICATOR_ADDRESS
) as unknown as TESTNitroAdjudicator & Contract;

const token = setupContract(
  provider,
  TokenArtifact,
  process.env.TEST_TOKEN_ADDRESS
) as unknown as Token & Contract;

const badToken = setupContract(
  provider,
  BadTokenArtifact,
  process.env.BAD_TOKEN_ADDRESS
) as unknown as BadToken & Contract;

const signer0 = getTestProvider().getSigner(0); // Convention matches setupContract function
let signer0Address: string;
const participants: string[] = [];
const challengeDuration = 0x1000;
let appDefinition: string;

const ETH = MAGIC_ADDRESS_INDICATING_ETH;
const ERC20 = token.address;
const BadERC20 = badToken.address;

// Populate destinations array
for (let i = 0; i < 3; i++) {
  participants[i] = Wallet.createRandom({extraEntropy: utils.id('erc20-deposit-test')}).address;
}

beforeAll(async () => {
  signer0Address = await signer0.getAddress();
  appDefinition = getCountingAppContractAddress();
});

const description0 = 'Deposits Tokens (expectedHeld = 0)';
const description1 = 'Deposits Tokens (expectedHeld = 1)';
const description2 = 'Reverts deposit of Tokens (expectedHeld > holdings)';
const description3 = 'Reverts deposit of Tokens (expectedHeld < holdings)';
const description5 = 'Deposits ETH (msg.value = amount , expectedHeld = 0)';
const description6 = 'Reverts deposit of ETH (msg.value = amount, expectedHeld > holdings)';
const description9 = 'Deposits a Bad token (expectedHeld = 0)';

describe('deposit', () => {
  let channelNonce = getRandomNonce('deposit');
  afterEach(() => {
    channelNonce = BigNumber.from(channelNonce).add(1).toHexString();
  });
  it.each`
    description     | asset       | held | expectedHeld | amount | heldAfter | reasonString
    ${description0} | ${ERC20}    | ${0} | ${0}         | ${1}   | ${1}      | ${undefined}
    ${description1} | ${ERC20}    | ${1} | ${1}         | ${1}   | ${2}      | ${undefined}
    ${description2} | ${ERC20}    | ${0} | ${1}         | ${2}   | ${0}      | ${'held != expectedHeld'}
    ${description3} | ${ERC20}    | ${1} | ${0}         | ${2}   | ${2}      | ${'held != expectedHeld'}
    ${description5} | ${ETH}      | ${0} | ${0}         | ${1}   | ${1}      | ${undefined}
    ${description6} | ${ETH}      | ${0} | ${1}         | ${2}   | ${0}      | ${'held != expectedHeld'}
    ${description9} | ${BadERC20} | ${0} | ${0}         | ${1}   | ${1}      | ${undefined}
  `('$description', async ({asset, held, expectedHeld, amount, reasonString, heldAfter}) => {
    held = BigNumber.from(held);
    expectedHeld = BigNumber.from(expectedHeld);
    amount = BigNumber.from(amount);
    heldAfter = BigNumber.from(heldAfter);

    const destination = getChannelId({
      channelNonce,
      participants,
      appDefinition,
      challengeDuration,
    });

    if (asset === ERC20) {
      // Check msg.sender has enough tokens
      const balance = await token.balanceOf(signer0Address);
      await expect(balance.gte(held.add(amount))).toBe(true);

      // Increase allowance
      await (await token.increaseAllowance(testNitroAdjudicator.address, held.add(amount))).wait(); // Approve enough for setup and main test

      // Check allowance updated
      const allowance = BigNumber.from(
        await token.allowance(signer0Address, testNitroAdjudicator.address)
      );
      expect(allowance.sub(amount).sub(held).gte(0)).toBe(true);
    }

    if (asset === BadERC20) {
      // Check msg.sender has enough tokens
      const balance = await badToken.balanceOf(signer0Address);
      await expect(balance.gte(held.add(amount))).toBe(true);

      // Increase allowance
      await (
        await badToken.increaseAllowance(testNitroAdjudicator.address, held.add(amount))
      ).wait(); // Approve enough for setup and main test

      // Check allowance updated
      const allowance = BigNumber.from(
        await badToken.allowance(signer0Address, testNitroAdjudicator.address)
      );
      expect(allowance.sub(amount).sub(held).gte(0)).toBe(true);
    }

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
      if (asset === ERC20 || asset == BadERC20) {
        const {data: amountTransferred} = getTransferEvent(events);
        expect(held.eq(amountTransferred)).toBe(true);
      }
    }

    const balanceBefore = await getBalance(asset, signer0Address);

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
      const depositedEvent = getDepositedEvent(events);
      expect(depositedEvent).toMatchObject({
        destination,
        amountDeposited: heldAfter.sub(held),
        destinationHoldings: heldAfter,
      });

      if (asset == ERC20 || asset == BadERC20) {
        const amountTransferred = BigNumber.from(getTransferEvent(events).data);
        expect(heldAfter.sub(held).eq(amountTransferred)).toBe(true);
        const balanceAfter = await getBalance(asset, signer0Address);
        expect(balanceAfter.eq(balanceBefore.sub(heldAfter.sub(held)))).toBe(true);
      }

      const allocatedAmount = await testNitroAdjudicator.holdings(asset, destination);
      await expect(allocatedAmount).toEqual(heldAfter);
    }
  });
});

const getDepositedEvent = (events: ethers.Event[]) =>
  // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
  events.find(({event}) => event === 'Deposited')!.args;

const getTransferEvent = (events: ethers.Event[]) =>
  // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
  events.find(({topics}) => topics[0] === token.filters.Transfer(AddressZero).topics![0])!;

async function getBalance(asset: string, address: string) {
  switch (asset) {
    case ETH:
      return BigNumber.from(await provider.getBalance(address));
    case ERC20:
      return BigNumber.from(await token.balanceOf(address));
    case BadERC20:
      return BigNumber.from(await badToken.balanceOf(address));
  }
  throw Error('unrecognized asset');
}

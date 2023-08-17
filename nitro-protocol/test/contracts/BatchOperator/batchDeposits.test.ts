import {ethers, Contract, Wallet, BigNumber, utils, BigNumberish} from 'ethers';
import {it} from '@jest/globals';
import {expectRevert} from '@statechannels/devtools';

import {MAGIC_ADDRESS_INDICATING_ETH, getChannelId, getRandomNonce} from '../../../src';
import {getTestProvider, setupContract} from '../../test-helpers';
// artifacts
import TokenArtifact from '../../../artifacts/contracts/Token.sol/Token.json';
import BadTokenArtifact from '../../../artifacts/contracts/test/BadToken.sol/BadToken.json';
import TESTNitroAdjudicatorArtifact from '../../../artifacts/contracts/test/TESTNitroAdjudicator.sol/TESTNitroAdjudicator.json';
import BatchOperatorArtifact from '../../../artifacts/contracts/auxiliary/BatchOperator.sol/BatchOperator.json';
import {Token, BadToken, TESTNitroAdjudicator, BatchOperator} from '../../../typechain-types';
const provider = getTestProvider();

const batchOperator = setupContract(
  provider,
  BatchOperatorArtifact,
  process.env.BATCH_OPERATOR_ADDRESS!
) as unknown as BatchOperator & Contract;

const testNitroAdjudicator = setupContract(
  provider,
  TESTNitroAdjudicatorArtifact,
  process.env.TEST_NITRO_ADJUDICATOR_ADDRESS!
) as unknown as TESTNitroAdjudicator & Contract;

const token = setupContract(
  provider,
  TokenArtifact,
  process.env.TEST_TOKEN_ADDRESS!
) as unknown as Token & Contract;

const badToken = setupContract(
  provider,
  BadTokenArtifact,
  process.env.BAD_TOKEN_ADDRESS!
) as unknown as BadToken & Contract;

const ETH = MAGIC_ADDRESS_INDICATING_ETH;
const ERC20 = token.address;
const BadERC20 = badToken.address;
const consensusApp = process.env.CONSENSUS_APP_ADDRESS!;

const signer = provider.getSigner(0);
let signerAddress: string;

const batchSize = 3;
const counterparties: string[] = [];
for (let i = 0; i < batchSize; i++) {
  counterparties[i] = Wallet.createRandom({
    extraEntropy: utils.id('multi-asset-holder-deposit-test'),
  }).address;
}

beforeAll(async () => {
  signerAddress = await signer.getAddress();
});

const description0 = 'Deposits Eth to Multiple Channels (expectedHeld = 0)';
const description1 = 'Deposits Eth to Multiple Channels (expectedHeld = 1)';
const description2 = 'Deposits Eth to Multiple Channels (mixed expectedHeld)';
const description3 =
  'Reverts deposit of Eth to Multiple Channels (mismatched expectedHeld, zero expected)';
const description4 =
  'Reverts deposit of Eth to Multiple Channels (mismatched expectedHeld, nonzero expected)';

const description5 = 'Deposits Tokens to Multiple Channels (expectedHeld = 0)';
const description6 = 'Deposits Tokens to Multiple Channels (expectedHeld = 1)';
const description7 = 'Deposits Tokens to Multiple Channels (mixed expectedHeld)';
const description8 =
  'Reverts deposit of Tokens to Multiple Channels (mismatched expectedHeld, zero expected)';
const description9 =
  'Reverts deposit of Tokens to Multiple Channels (mismatched expectedHeld, nonzero expected)';

const description10 = 'Deposits BadToken to Multiple Channels (expectedHeld = 0)';

const description11 = 'Reverts if input lengths do not match';

const unexpectedHeld = 'held != expectedHeld';

type testParams = {
  description: string;
  assetId: string;
  expectedHelds: number[];
  amounts: number[];
  heldAfters: number[];
  reasonString: string;
};

function sum(x: BigNumber[]): BigNumber {
  return x.reduce((s, n) => s.add(n));
}

describe('deposit_batch', () => {
  it.each`
    description     | assetId | expectedHelds | amounts      | heldAfters   | reasonString
    ${description1} | ${ETH}  | ${[1, 1, 1]}  | ${[2, 2, 2]} | ${[3, 3, 3]} | ${''}
  `(
    // ${description0} | ${ETH}  | ${[0, 0, 0]}  | ${[1, 2, 3]} | ${[1, 2, 3]} | ${''}
    // ${description2} | ${ETH}  | ${[0, 1, 2]}  | ${[1, 1, 1]} | ${[1, 2, 3]} | ${''}
    // ${description3}  | ${ETH}      | ${[0, 0, 0]}  | ${[1, 1, 1]} | ${[1, 1, 1]} | ${unexpectedHeld}
    // ${description4}  | ${ETH}      | ${[1, 1, 1]}  | ${[1, 1, 1]} | ${[2, 2, 2]} | ${unexpectedHeld}
    // ${description5} | ${ERC20} | ${[0, 0, 0]}  | ${[1, 1, 1]} | ${[1, 1, 1]} | ${''}
    // ${description6}  | ${ERC20}    | ${[1, 1, 1]}  | ${[1, 1, 1]} | ${[2, 2, 2]} | ${''}
    // ${description7}  | ${ERC20}    | ${[0, 1, 0]}  | ${[1, 1, 1]} | ${[1, 2, 1]} | ${''}
    // ${description8}  | ${ERC20}    | ${[0, 0, 0]}  | ${[1, 1, 1]} | ${[1, 1, 1]} | ${unexpectedHeld}
    // ${description9}  | ${ERC20}    | ${[1, 1, 1]}  | ${[1, 1, 1]} | ${[2, 2, 2]} | ${unexpectedHeld}
    // ${description10} | ${BadERC20} | ${[0, 0, 0]}  | ${[1, 1, 1]} | ${[1, 1, 1]} | ${''}
    // ${description11} | ${ETH}      | ${[0, 0]}     | ${[1, 1, 1]} | ${[1, 1, 1]} | ${'Array lengths must match'}
    '$description',
    async ({
      description,
      assetId,
      expectedHelds,
      amounts,
      heldAfters,
      reasonString,
    }: testParams) => {
      console.log('asset', assetId);

      ///////////////////////////////////////
      //
      // Construct deposit_batch parameters
      //
      ///////////////////////////////////////
      const channelIds = counterparties.map(counterparty =>
        getChannelId({
          channelNonce: getRandomNonce(description),
          participants: [signerAddress, counterparty],
          appDefinition: consensusApp,
          challengeDuration: 100,
        })
      );
      const expectedHeldsBN = expectedHelds.map(x => BigNumber.from(x));
      const amountsBN = amounts.map(x => BigNumber.from(x));
      const heldAftersBN = heldAfters.map(x => BigNumber.from(x));
      const totalValue = sum(amountsBN);
      const totalExpectedHeld = sum(expectedHeldsBN);

      if (assetId === ERC20) {
        const balance = await token.balanceOf(signerAddress);
        console.log('erc20 balance:', balance);
        await (
          await token.increaseAllowance(batchOperator.address, totalValue.add(totalExpectedHeld))
        ).wait();
        await (
          await token.increaseAllowance(
            testNitroAdjudicator.address,
            totalValue.add(totalExpectedHeld)
          )
        ).wait();
        // Check Balance Updated
        const allowance = BigNumber.from(
          await token.allowance(signerAddress, batchOperator.address)
        );
        console.log(`Allowance: `, allowance);
      }

      if (assetId === BadERC20) {
        await (
          await badToken.increaseAllowance(batchOperator.address, totalValue.add(totalExpectedHeld))
        ).wait();
        await (
          await badToken.increaseAllowance(
            testNitroAdjudicator.address,
            totalValue.add(totalExpectedHeld)
          )
        ).wait();
      }

      ///////////////////////////////////////
      //
      // Set up preexisting holdings (if any)
      //
      ///////////////////////////////////////

      await Promise.all(
        channelIds.map(async (channelId, i) => {
          // apply incorrect amount if unexpectedHeld reasonString is set
          const value =
            reasonString == unexpectedHeld ? expectedHeldsBN[i].add(1) : expectedHeldsBN[i];
          const {events} = await (
            await testNitroAdjudicator.deposit(assetId, channelId, 0, value, {
              value: assetId === ETH ? value : 0,
            })
          ).wait();
          expect(events).not.toBe(undefined);
        })
      );

      for (const c of channelIds) {
        const holdings = await testNitroAdjudicator.holdings(assetId, c);
        console.log(`pre-holdings[${assetId}][${c}]`, holdings);
      }

      ///////////////////////////////////////
      //
      // Execute deposit
      //
      ///////////////////////////////////////

      console.log(`expectedHeldsBN`, expectedHeldsBN);

      const tx =
        assetId === ETH
          ? batchOperator.deposit_batch_eth(channelIds, expectedHeldsBN, amountsBN, {
              value: totalValue,
            })
          : batchOperator.deposit_batch_erc(
              assetId,
              channelIds,
              expectedHeldsBN,
              amountsBN,
              totalValue
            );

      ///////////////////////////////////////
      //
      // Check postconditions
      //
      ///////////////////////////////////////
      if (reasonString != '') {
        await expectRevert(() => tx, reasonString);
      } else {
        const {events} = await (await tx).wait();
        console.log('events', events);

        const holdings: BigNumber[] = [];
        for (let i = 0; i < channelIds.length; i++) {
          holdings.push(await testNitroAdjudicator.holdings(assetId, channelIds[i]));
          console.log(`post-holdings[${assetId}][${channelIds[i]}]`, holdings[i]);
        }
        for (let i = 0; i < channelIds.length; i++) {
          expect(holdings[i]).toEqual(heldAftersBN[i]);
        }
      }
    }
  );
});

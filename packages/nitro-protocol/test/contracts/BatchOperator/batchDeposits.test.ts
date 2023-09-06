import {it} from '@jest/globals';
import {expectRevert} from '@statechannels/devtools';
import {BigNumber, Contract, Wallet, utils} from 'ethers';

import {MAGIC_ADDRESS_INDICATING_ETH, getChannelId, getRandomNonce} from '../../../src';
import {getTestProvider, setupContract} from '../../test-helpers';
// artifacts
import NitroAdjudicatorArtifact from '../../../artifacts/contracts/NitroAdjudicator.sol/NitroAdjudicator.json';
import TokenArtifact from '../../../artifacts/contracts/Token.sol/Token.json';
import BatchOperatorArtifact from '../../../artifacts/contracts/auxiliary/BatchOperator.sol/BatchOperator.json';
import BadTokenArtifact from '../../../artifacts/contracts/test/BadToken.sol/BadToken.json';
import {BadToken, BatchOperator, NitroAdjudicator, Token} from '../../../typechain-types';

const provider = getTestProvider();

const batchOperator = setupContract(
  provider,
  BatchOperatorArtifact,
  process.env.BATCH_OPERATOR_ADDRESS || ''
) as unknown as BatchOperator & Contract;

const nitroAdjudicator = setupContract(
  provider,
  NitroAdjudicatorArtifact,
  process.env.NITRO_ADJUDICATOR_ADDRESS || ''
) as unknown as NitroAdjudicator & Contract;

const token = setupContract(
  provider,
  TokenArtifact,
  process.env.TEST_TOKEN_ADDRESS || ''
) as unknown as Token & Contract;

const badToken = setupContract(
  provider,
  BadTokenArtifact,
  process.env.BAD_TOKEN_ADDRESS || ''
) as unknown as BadToken & Contract;

const ETH = MAGIC_ADDRESS_INDICATING_ETH;
const ERC20 = token.address;
const BadERC20 = badToken.address;
const consensusApp = process.env.CONSENSUS_APP_ADDRESS;

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
    description      | assetId     | expectedHelds | amounts      | heldAfters   | reasonString
    ${description0}  | ${ETH}      | ${[0, 0, 0]}  | ${[1, 2, 3]} | ${[1, 2, 3]} | ${''}
    ${description1}  | ${ETH}      | ${[1, 1, 1]}  | ${[2, 2, 2]} | ${[3, 3, 3]} | ${''}
    ${description2}  | ${ETH}      | ${[0, 1, 2]}  | ${[1, 1, 1]} | ${[1, 2, 3]} | ${''}
    ${description3}  | ${ETH}      | ${[0, 0, 0]}  | ${[1, 1, 1]} | ${[1, 1, 1]} | ${unexpectedHeld}
    ${description4}  | ${ETH}      | ${[1, 1, 1]}  | ${[1, 1, 1]} | ${[2, 2, 2]} | ${unexpectedHeld}
    ${description5}  | ${ERC20}    | ${[0, 0, 0]}  | ${[1, 1, 1]} | ${[1, 1, 1]} | ${''}
    ${description6}  | ${ERC20}    | ${[1, 1, 1]}  | ${[1, 1, 1]} | ${[2, 2, 2]} | ${''}
    ${description7}  | ${ERC20}    | ${[0, 1, 0]}  | ${[1, 1, 1]} | ${[1, 2, 1]} | ${''}
    ${description8}  | ${ERC20}    | ${[0, 0, 0]}  | ${[1, 1, 1]} | ${[1, 1, 1]} | ${unexpectedHeld}
    ${description9}  | ${ERC20}    | ${[1, 1, 1]}  | ${[1, 1, 1]} | ${[2, 2, 2]} | ${unexpectedHeld}
    ${description10} | ${BadERC20} | ${[0, 0, 0]}  | ${[1, 1, 1]} | ${[1, 1, 1]} | ${''}
    ${description11} | ${ETH}      | ${[0, 0]}     | ${[1, 1, 1]} | ${[1, 1, 1]} | ${'Array lengths must match'}
  `('$description', async tc => {
    const {description, assetId, expectedHelds, amounts, heldAfters, reasonString} =
      tc as testParams;
    ///////////////////////////////////////
    //
    // Construct deposit_batch parameters
    //
    ///////////////////////////////////////
    const channelIds = counterparties.map(counterparty =>
      getChannelId({
        channelNonce: getRandomNonce(description),
        participants: [signerAddress, counterparty],
        appDefinition: consensusApp || '',
        challengeDuration: 100,
      })
    );
    const expectedHeldsBN = expectedHelds.map(x => BigNumber.from(x));
    const amountsBN = amounts.map(x => BigNumber.from(x));
    const heldAftersBN = heldAfters.map(x => BigNumber.from(x));
    const totalValue = sum(amountsBN);

    if (assetId === ERC20) {
      await (
        await token.increaseAllowance(batchOperator.address, totalValue.add(totalValue))
      ).wait();
    }

    if (assetId === BadERC20) {
      await (
        await badToken.increaseAllowance(batchOperator.address, totalValue.add(totalValue))
      ).wait();
    }

    ///////////////////////////////////////
    //
    // Set up preexisting holdings (if any)
    //
    ///////////////////////////////////////

    await Promise.all(
      expectedHeldsBN.map(async (expected, i) => {
        const channelID = channelIds[i];
        // apply incorrect amount if unexpectedHeld reasonString is set
        const value = reasonString == unexpectedHeld ? expected.add(1) : expected;

        if (assetId === ERC20) {
          await (await token.increaseAllowance(nitroAdjudicator.address, value)).wait();
        }
        if (assetId === BadERC20) {
          await (await badToken.increaseAllowance(nitroAdjudicator.address, value)).wait();
        }

        const {events} = await (
          await nitroAdjudicator.deposit(assetId, channelID, 0, value, {
            value: assetId === ETH ? value : 0,
          })
        ).wait();
        expect(events).not.toBe(undefined);
      })
    );

    ///////////////////////////////////////
    //
    // Execute deposit
    //
    ///////////////////////////////////////

    const tx =
      assetId === ETH
        ? batchOperator.deposit_batch_eth(channelIds, expectedHeldsBN, amountsBN, {
            value: totalValue,
          })
        : batchOperator.deposit_batch_erc20(
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
      await (await tx).wait();

      for (let i = 0; i < channelIds.length; i++) {
        const channelId = channelIds[i];
        const expectedHoldings = heldAftersBN[i];

        const holdings = await nitroAdjudicator.holdings(assetId, channelId);
        expect(holdings).toEqual(expectedHoldings);
      }
    }
  });
});

import {readFileSync, existsSync} from 'fs';

import {BigNumber} from 'ethers';

import {encodeOutcome, Outcome} from '../src';
import {computeReclaimEffects} from '../src/contract/multi-asset-holder';
import {MAGIC_ADDRESS_INDICATING_ETH} from '../src/transactions';

import {
  waitForChallengesToTimeOut,
  Y,
  X,
  LforX,
  LforV,
  getChannelBatch,
  assertEthBalancesAndHoldings,
  amountForAlice,
  amountForAliceAndBob,
  V,
  challengeVirtualPaymentChannelWithVoucher,
  Alice,
  Bob,
  paymentAmount,
} from './fixtures';
import {batchSizes, GasResults} from './gas';
import {challengeChannelAndExpectGas} from './jestSetup';
import {nitroAdjudicator, token} from './localSetup';

/**
 * Ensures the asset holding contract always has a nonzero token balance.
 */
async function addResidualTokenBalance(asset: string) {
  /**
   * Funding someOtherChannel with tokens, as well as the channel in question
   * makes the benchmark more realistic. In practice many other
   * channels are funded by the nitro adjudicator. If we didn't reflect
   * that, our benchmark might reflect a gas refund for clearing storage
   * in the token contract (setting the token balance of the asset holder to 0)
   * which we would only expect in rare cases.
   */
  await (await nitroAdjudicator.deposit(asset, Y.channelId, 0, 1)).wait();
}

if (!existsSync(__dirname + '/gasResults.json')) {
  throw new Error('Error: file "gasResults.json" with previous benchmark results does not exist');
}

const gasRequiredTo = JSON.parse(
  readFileSync(__dirname + '/gasResults.json', 'utf-8')
) as GasResults;

describe('Consumes the expected gas for deployments', () => {
  it(`when deploying the NitroAdjudicator`, async () => {
    await expect(await nitroAdjudicator.deployTransaction).toConsumeGas(
      gasRequiredTo.deployInfrastructureContracts.satp.NitroAdjudicator
    );
  });
});
describe('Consumes the expected gas for deposits', () => {
  it(`when directly funding a channel with ETH (first deposit)`, async () => {
    await expect(
      await nitroAdjudicator.deposit(MAGIC_ADDRESS_INDICATING_ETH, X.channelId, 0, 5, {value: 5})
    ).toConsumeGas(gasRequiredTo.directlyFundAChannelWithETHFirst.satp);
  });

  it(`when directly funding a channel with ETH (second deposit)`, async () => {
    // begin setup
    const setupTX = nitroAdjudicator.deposit(MAGIC_ADDRESS_INDICATING_ETH, X.channelId, 0, 5, {
      value: 5,
    });
    await (await setupTX).wait();
    // end setup
    await expect(
      await nitroAdjudicator.deposit(MAGIC_ADDRESS_INDICATING_ETH, X.channelId, 5, 5, {value: 5})
    ).toConsumeGas(gasRequiredTo.directlyFundAChannelWithETHSecond.satp);
  });

  for (const batchSize of batchSizes) {
    const batch = getChannelBatch(batchSize);

    it(`when batch funding ${batchSize} channels with ETH (first deposit)`, async () => {
      await expect(
        await nitroAdjudicator.deposit_batch(
          MAGIC_ADDRESS_INDICATING_ETH,
          batch.map(c => c.channelId),
          batch.map(() => 0),
          batch.map(() => 5),
          {value: 5 * batch.length}
        )
      ).toConsumeGas(gasRequiredTo.batchFundChannelsWithETHFirst.satp[batchSize]);
    });

    it(`when batch funding ${batchSize} channels with ETH (second deposit)`, async () => {
      // begin setup
      await (
        await nitroAdjudicator.deposit_batch(
          MAGIC_ADDRESS_INDICATING_ETH,
          batch.map(c => c.channelId),
          batch.map(() => 0),
          batch.map(() => 5),
          {value: 5 * batch.length}
        )
      ).wait();
      // end setup

      await expect(
        await nitroAdjudicator.deposit_batch(
          MAGIC_ADDRESS_INDICATING_ETH,
          batch.map(c => c.channelId),
          batch.map(() => 5),
          batch.map(() => 5),
          {value: 5 * batch.length}
        )
      ).toConsumeGas(gasRequiredTo.batchFundChannelsWithETHSecond.satp[batchSize]);
    });
  }

  it(`when directly funding a channel with an ERC20 (first deposit)`, async () => {
    // begin setup
    await (await token.transfer(nitroAdjudicator.address, 1)).wait(); // The asset holder already has some tokens (for other channels)
    // end setup
    await expect(await token.increaseAllowance(nitroAdjudicator.address, 100)).toConsumeGas(
      gasRequiredTo.directlyFundAChannelWithERC20First.satp.approve
    );
    await expect(await nitroAdjudicator.deposit(token.address, X.channelId, 0, 5)).toConsumeGas(
      gasRequiredTo.directlyFundAChannelWithERC20First.satp.deposit
    );
  });

  it(`when directly funding a channel with an ERC20 (second deposit)`, async () => {
    // begin setup
    await (await token.increaseAllowance(nitroAdjudicator.address, 100)).wait();
    await (await nitroAdjudicator.deposit(token.address, X.channelId, 0, 5)).wait(); // The asset holder already has some tokens *for this channel*
    await (await token.decreaseAllowance(nitroAdjudicator.address, 95)).wait(); // reset allowance to zero
    // end setup
    await expect(await token.increaseAllowance(nitroAdjudicator.address, 100)).toConsumeGas(
      gasRequiredTo.directlyFundAChannelWithERC20Second.satp.approve
    );
    await expect(await nitroAdjudicator.deposit(token.address, X.channelId, 5, 5)).toConsumeGas(
      gasRequiredTo.directlyFundAChannelWithERC20Second.satp.deposit
    );
  });
});

describe('Consumes the expected gas for happy-path exits', () => {
  it(`when exiting a directly funded (with ETH) channel`, async () => {
    // begin setup
    await (
      await nitroAdjudicator.deposit(MAGIC_ADDRESS_INDICATING_ETH, X.channelId, 0, 10, {value: 10})
    ).wait();
    // end setup
    await expect(await X.concludeAndTransferAllAssetsTx(MAGIC_ADDRESS_INDICATING_ETH)).toConsumeGas(
      gasRequiredTo.ETHexit.satp
    );
  });

  it(`when exiting a directly funded (with ERC20s) channel`, async () => {
    // begin setup
    await (await token.increaseAllowance(nitroAdjudicator.address, 100)).wait();
    await (await nitroAdjudicator.deposit(token.address, X.channelId, 0, 10)).wait();
    await addResidualTokenBalance(token.address);
    // end setup
    await expect(await X.concludeAndTransferAllAssetsTx(token.address)).toConsumeGas(
      gasRequiredTo.ERC20exit.satp
    );
  });
});

describe('Consumes the expected gas for sad-path exits', () => {
  it(`when exiting a directly funded (with ETH) channel`, async () => {
    // begin setup
    await (
      await nitroAdjudicator.deposit(MAGIC_ADDRESS_INDICATING_ETH, X.channelId, 0, 10, {value: 10})
    ).wait();
    // end setup
    // initially                 ⬛ ->  X  -> 👩
    const {proof, finalizesAt} = await challengeChannelAndExpectGas(
      X,
      MAGIC_ADDRESS_INDICATING_ETH,
      gasRequiredTo.ETHexitSad.satp.challenge
    );
    // begin wait
    await waitForChallengesToTimeOut([finalizesAt]);
    // end wait
    // challenge + timeout       ⬛ -> (X) -> 👩
    await expect(
      await nitroAdjudicator.transferAllAssets(
        X.channelId,
        proof.outcome, // outcome,
        proof.stateHash // stateHash
      )
    ).toConsumeGas(gasRequiredTo.ETHexitSad.satp.transferAllAssets);
    // transferAllAssets ⬛ --------> 👩
    expect(
      gasRequiredTo.ETHexitSad.satp.challenge + gasRequiredTo.ETHexitSad.satp.transferAllAssets
    ).toEqual(gasRequiredTo.ETHexitSad.satp.total);
  });

  it(`when exiting a ledger funded (with ETH) channel`, async () => {
    // begin setup
    await (
      await nitroAdjudicator.deposit(MAGIC_ADDRESS_INDICATING_ETH, LforX.channelId, 0, 10, {
        value: 10,
      })
    ).wait();
    // end setup
    // initially                   ⬛ ->  L  ->  X  -> 👩
    const {proof: ledgerProof, finalizesAt: ledgerFinalizesAt} = await challengeChannelAndExpectGas(
      LforX,
      MAGIC_ADDRESS_INDICATING_ETH,
      gasRequiredTo.ETHexitSadLedgerFunded.satp.challengeL
    );
    const {proof, finalizesAt} = await challengeChannelAndExpectGas(
      X,
      MAGIC_ADDRESS_INDICATING_ETH,
      gasRequiredTo.ETHexitSadLedgerFunded.satp.challengeX
    );
    // begin wait
    await waitForChallengesToTimeOut([ledgerFinalizesAt, finalizesAt]); // just go to the max one
    // end wait
    // challenge X, L and timeout  ⬛ -> (L) -> (X) -> 👩
    await expect(
      await nitroAdjudicator.transferAllAssets(
        LforX.channelId,
        ledgerProof.outcome, // outcome
        ledgerProof.stateHash // stateHash
      )
    ).toConsumeGas(gasRequiredTo.ETHexitSadLedgerFunded.satp.transferAllAssetsL);
    // transferAllAssetsL  ⬛ --------> (X) -> 👩
    await expect(
      await nitroAdjudicator.transferAllAssets(
        X.channelId,
        proof.outcome, // outcome
        proof.stateHash // stateHash
      )
    ).toConsumeGas(gasRequiredTo.ETHexitSadLedgerFunded.satp.transferAllAssetsX);
    // transferAllAssetsX  ⬛ ---------------> 👩

    // meta-test here to confirm the total recorded in gas.ts is up to date
    // with the recorded costs of each step
    expect(
      gasRequiredTo.ETHexitSadLedgerFunded.satp.challengeL +
        gasRequiredTo.ETHexitSadLedgerFunded.satp.transferAllAssetsL +
        gasRequiredTo.ETHexitSadLedgerFunded.satp.challengeX +
        gasRequiredTo.ETHexitSadLedgerFunded.satp.transferAllAssetsX
    ).toEqual(gasRequiredTo.ETHexitSadLedgerFunded.satp.total);
  });

  it(`when exiting a virtual funded (with ETH) channel`, async () => {
    // begin setup
    await (
      await nitroAdjudicator.deposit(
        // This deposit represents what in reality would likely be two deposits (one from Bob, one from Ingrid)
        MAGIC_ADDRESS_INDICATING_ETH,
        LforV.channelId,
        0,
        amountForAliceAndBob,
        {
          value: amountForAliceAndBob,
        }
      )
    ).wait();
    // end setup
    // initially                   ⬛ ->  L  ->  V  -> 👨
    // challenge L
    const {proof: ledgerProof, finalizesAt: ledgerFinalizesAt} = await challengeChannelAndExpectGas(
      LforV,
      MAGIC_ADDRESS_INDICATING_ETH,
      gasRequiredTo.ETHexitSadVirtualFunded.satp.challengeL
    );

    // challenge V with a voucher
    const {
      stateHash: vStateHash,
      outcome: vOutcome,
      gasUsed: vGasUsed,
      finalizesAt: vFinalizesAt,
    } = await challengeVirtualPaymentChannelWithVoucher(
      V,
      MAGIC_ADDRESS_INDICATING_ETH,
      BigNumber.from(paymentAmount).toNumber(),
      Alice,
      Bob
    );

    expect(vGasUsed).toEqual(gasRequiredTo.ETHexitSadVirtualFunded.satp.challengeV);

    // begin wait
    await waitForChallengesToTimeOut([ledgerFinalizesAt, vFinalizesAt]);
    // end wait
    // challenge L,V   + timeout   ⬛ -> (L) -> (V) -> 👨
    await assertEthBalancesAndHoldings(
      {Alice: 0, Bob: 0, Ingrid: 0},
      {LforV: amountForAliceAndBob, V: 0}
    );
    await expect(
      await nitroAdjudicator.reclaim({
        sourceChannelId: LforV.channelId,
        sourceStateHash: ledgerProof.stateHash,
        sourceOutcomeBytes: encodeOutcome(ledgerProof.outcome),
        sourceAssetIndex: 0,
        indexOfTargetInSource: 2,
        targetStateHash: vStateHash,
        targetOutcomeBytes: encodeOutcome(vOutcome),
        targetAssetIndex: 0,
      })
    ).toConsumeGas(gasRequiredTo.ETHexitSadVirtualFunded.satp.reclaimL);
    // reclaim L                   ⬛ -- (L) --------> 👨

    await assertEthBalancesAndHoldings(
      {Alice: 0, Bob: 0, Ingrid: 0},
      {LforV: amountForAliceAndBob, V: 0}
    );

    // track change to ledger outcome caused by calling reclaim
    const updatedAllocations = computeReclaimEffects(
      ledgerProof.outcome[0].allocations,
      vOutcome[0].allocations,
      2
    );
    const updatedOutcome: Outcome = [
      {
        ...ledgerProof.outcome[0],
        allocations: updatedAllocations,
      },
    ];

    await expect(
      await nitroAdjudicator.transferAllAssets(
        LforV.channelId,
        updatedOutcome,
        ledgerProof.stateHash // stateHash
      )
    ).toConsumeGas(gasRequiredTo.ETHexitSadVirtualFunded.satp.transferAllAssetsL);
    // transferAllAssetsL          ⬛ ---------------> 👨

    await assertEthBalancesAndHoldings(
      {
        Bob: BigNumber.from(paymentAmount), // Bob gets his paymennt
        Ingrid: BigNumber.from(amountForAlice).sub(BigNumber.from(paymentAmount)), // Ingrid is adjusted down, she will be compensated in the other ledger channel
      },
      {LforV: 0, V: 0}
    );

    // meta-test here to confirm the total recorded in gas.ts is up to date
    // with the recorded costs of each step
    expect(
      (Object.values(gasRequiredTo.ETHexitSadVirtualFunded.satp) as number[]).reduce(
        (a, b) => a + b
      ) - gasRequiredTo.ETHexitSadVirtualFunded.satp.total
    ).toEqual(gasRequiredTo.ETHexitSadVirtualFunded.satp.total);
  });
});

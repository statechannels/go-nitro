import {writeFileSync} from 'fs';

import {BigNumber} from 'ethers';

import {encodeOutcome, Outcome} from '../src';
import {computeReclaimEffects} from '../src/contract/multi-asset-holder';
import {MAGIC_ADDRESS_INDICATING_ETH} from '../src/transactions';

import {
  waitForChallengesToTimeOut,
  challengeChannel,
  Y,
  X,
  LforX,
  gasUsed,
  executeAndRevert,
  LforV,
  V,
  Alice,
  Bob,
  challengeVirtualPaymentChannelWithVoucher,
  paymentAmount,
  getChannelBatch,
  checkpointChannel,
} from './fixtures';
import {batchSizes, emptyGasResults} from './gas';
import {deployContracts, nitroAdjudicator, batchOperator, token} from './localSetup';
import {challengeChannelAndExpectGas, respondWithChallengeAndExpectGas} from './jestSetup';

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

// The channel being benchmarked is a 2 party null app funded with 5 wei / tokens each.
// KEY
// ---
// ⬛ -> funding on chain (from Alice)
//  C    channel not yet on chain
// (C)   channel finalized on chain
// 👩    Alice's external destination (e.g. her EOA)
async function main() {
  await deployContracts();

  const gasResults = emptyGasResults;

  // *****************
  // Deployments:
  // *****************

  // deploy NitroAdjudicator
  // Singleton
  await executeAndRevert(async () => {
    gasResults.deployInfrastructureContracts.satp.NitroAdjudicator = await gasUsed(
      await nitroAdjudicator.deployTransaction
    );
  });

  // *****************
  // Deposits:
  // *****************

  // directly funding a channel with ETH (first deposit)
  await executeAndRevert(async () => {
    gasResults.directlyFundAChannelWithETHFirst.satp = await gasUsed(
      await nitroAdjudicator.deposit(MAGIC_ADDRESS_INDICATING_ETH, X.channelId, 0, 5, {value: 5})
    );
  });

  // directly funding a channel with ETH (second deposit)
  // meaning the second participant in the channel
  await executeAndRevert(async () => {
    // begin setup
    const setupTX = nitroAdjudicator.deposit(MAGIC_ADDRESS_INDICATING_ETH, X.channelId, 0, 5, {
      value: 5,
    });
    await (await setupTX).wait();
    // end setup
    gasResults.directlyFundAChannelWithETHSecond.satp = await gasUsed(
      await nitroAdjudicator.deposit(MAGIC_ADDRESS_INDICATING_ETH, X.channelId, 5, 5, {value: 5})
    );
  });

  for (const batchSize of batchSizes) {
    const batch = getChannelBatch(batchSize);
    const totalValue = 5 * batchSize;
    await executeAndRevert(async () => {
      // batch funding channels with ETH (first deposit)
      gasResults.batchFundChannelsWithETHFirst.satp['' + batchSize] = await gasUsed(
        await batchOperator.deposit_batch_eth(
          batch.map(c => c.channelId),
          batch.map(() => 0),
          batch.map(() => 5),
          {value: totalValue}
        )
      );
      // batch funding channels with ETH (second deposit)
      gasResults.batchFundChannelsWithETHSecond.satp['' + batchSize] = await gasUsed(
        await batchOperator.deposit_batch_eth(
          batch.map(c => c.channelId),
          batch.map(() => 5),
          batch.map(() => 5),
          {value: totalValue}
        )
      );
    });
    await executeAndRevert(async () => {
      await token.increaseAllowance(batchOperator.address, 3 * totalValue); // over-approve to avoid gas "refund" when approval returns to 0
      await (await token.transfer(nitroAdjudicator.address, 1)).wait(); // The asset holder already has some tokens (for other channels)

      // batch funding channels with ERC20 (first deposit)
      gasResults.batchFundChannelsWithERCFirst.satp['' + batchSize] = await gasUsed(
        await batchOperator.deposit_batch_erc20(
          token.address,
          batch.map(c => c.channelId),
          batch.map(() => 0),
          batch.map(() => 5),
          totalValue
        )
      );

      // batch funding channels with ERC20 (second deposit)
      gasResults.batchFundChannelsWithERCSecond.satp['' + batchSize] = await gasUsed(
        await batchOperator.deposit_batch_erc20(
          token.address,
          batch.map(c => c.channelId),
          batch.map(() => 5),
          batch.map(() => 5),
          totalValue
        )
      );
    });
  }

  // directly funding a channel with an ERC20 (first deposit)
  // The depositor begins with zero tokens approved for the AssetHolder
  // The AssetHolder begins with some token balance already
  // The depositor retains a nonzero balance of tokens after depositing
  // The depositor retains some tokens approved for the AssetHolder after depositing
  await executeAndRevert(async () => {
    // begin setup
    await (await token.transfer(nitroAdjudicator.address, 1)).wait(); // The asset holder already has some tokens (for other channels)
    // end setup
    gasResults.directlyFundAChannelWithERC20First.satp.approve = await gasUsed(
      await token.increaseAllowance(nitroAdjudicator.address, 100)
    );
    // ^^^^^
    // In principle this only needs to be done once per account
    // (the cost may be amortized over several deposits into this AssetHolder)
    gasResults.directlyFundAChannelWithERC20First.satp.deposit = await gasUsed(
      await nitroAdjudicator.deposit(token.address, X.channelId, 0, 5)
    );
  });

  // directly funding a channel with an ERC20 (second deposit)
  // meaning the second participant in the channel
  await executeAndRevert(async () => {
    // begin setup
    await (await token.increaseAllowance(nitroAdjudicator.address, 100)).wait();
    await (await nitroAdjudicator.deposit(token.address, X.channelId, 0, 5)).wait(); // The asset holder already has some tokens *for this channel*
    await (await token.decreaseAllowance(nitroAdjudicator.address, 95)).wait(); // reset allowance to zero
    // end setup
    gasResults.directlyFundAChannelWithERC20Second.satp.approve = await gasUsed(
      await token.increaseAllowance(nitroAdjudicator.address, 100)
    );
    // ^^^^^
    // In principle this only needs to be done once per account
    // (the cost may be amortized over several deposits into this AssetHolder)
    gasResults.directlyFundAChannelWithERC20Second.satp.deposit = await gasUsed(
      await nitroAdjudicator.deposit(token.address, X.channelId, 5, 5)
    );
  });

  // *****************
  // Happy-path exits:
  // *****************

  // exiting a directly funded (with ETH) channel
  // We completely liquidate the channel (paying out both parties)
  await executeAndRevert(async () => {
    // begin setup
    await (
      await nitroAdjudicator.deposit(MAGIC_ADDRESS_INDICATING_ETH, X.channelId, 0, 10, {value: 10})
    ).wait();
    // end setup
    gasResults.ETHexit.satp = await gasUsed(
      await X.concludeAndTransferAllAssetsTx(MAGIC_ADDRESS_INDICATING_ETH)
    );
  });

  // exiting a directly funded (with ERC20s) channel
  // We completely liquidate the channel (paying out both parties)
  await executeAndRevert(async () => {
    // begin setup
    await (await token.increaseAllowance(nitroAdjudicator.address, 100)).wait();
    await (await nitroAdjudicator.deposit(token.address, X.channelId, 0, 10)).wait();
    await addResidualTokenBalance(token.address);
    // end setup
    gasResults.ERC20exit.satp = await gasUsed(
      await X.concludeAndTransferAllAssetsTx(token.address)
    );
  });

  // *****************
  // Sad-path exits:
  // *****************

  // exiting a directly funded (with ETH) channel
  // Scenario: Counterparty Bob goes offline
  // initially                 ⬛ ->  X  -> 👩
  // challenge + timeout       ⬛ -> (X) -> 👩
  // transferAllAssets         ⬛ --------> 👩
  await executeAndRevert(async () => {
    // begin setup
    await (
      await nitroAdjudicator.deposit(MAGIC_ADDRESS_INDICATING_ETH, X.channelId, 0, 10, {value: 10})
    ).wait();
    // end setup
    // initially                 ⬛ ->  X  -> 👩
    const {challengeTx, proof, finalizesAt} = await challengeChannel(
      X,
      MAGIC_ADDRESS_INDICATING_ETH
    );
    gasResults.ETHexitSad.satp.challenge = await gasUsed(challengeTx);
    // begin wait
    await waitForChallengesToTimeOut([finalizesAt]);
    // end wait
    // challenge + timeout       ⬛ -> (X) -> 👩\
    gasResults.ETHexitSad.satp.transferAllAssets = await gasUsed(
      await nitroAdjudicator.transferAllAssets(
        X.channelId,
        proof.outcome, // outcome,
        proof.stateHash // stateHash
      )
    );
    // transferAllAssets ⬛ --------> 👩
    gasResults.ETHexitSad.satp.total =
      gasResults.ETHexitSad.satp.challenge + gasResults.ETHexitSad.satp.transferAllAssets;
  });

  // exiting a ledger funded (with ETH) channel
  // Scenario: Counterparty Bob goes offline
  // initially                   ⬛ ->  L  ->  X  -> 👩
  // challenge X, L and timeout  ⬛ -> (L) -> (X) -> 👩
  // transferAllAssetsL          ⬛ --------> (X) -> 👩
  // transferAllAssetsX          ⬛ ---------------> 👩
  await executeAndRevert(async () => {
    // begin setup
    await (
      await nitroAdjudicator.deposit(MAGIC_ADDRESS_INDICATING_ETH, LforX.channelId, 0, 10, {
        value: 10,
      })
    ).wait();
    // end setup
    // initially                   ⬛ ->  L  ->  X  -> 👩
    const {
      challengeTx: ledgerChallengeTx,
      proof: ledgerProof,
      finalizesAt: ledgerFinalizesAt,
    } = await challengeChannel(LforX, MAGIC_ADDRESS_INDICATING_ETH);
    gasResults.ETHexitSadLedgerFunded.satp.challengeL = await gasUsed(ledgerChallengeTx);

    const {challengeTx, proof, finalizesAt} = await challengeChannel(
      X,
      MAGIC_ADDRESS_INDICATING_ETH
    );
    gasResults.ETHexitSadLedgerFunded.satp.challengeX = await gasUsed(challengeTx);

    // begin wait
    await waitForChallengesToTimeOut([ledgerFinalizesAt, finalizesAt]); // just go to the max one
    // end wait
    // challenge X, L and timeout  ⬛ -> (L) -> (X) -> 👩
    gasResults.ETHexitSadLedgerFunded.satp.transferAllAssetsL = await gasUsed(
      await nitroAdjudicator.transferAllAssets(
        LforX.channelId,
        ledgerProof.outcome, // outcome
        ledgerProof.stateHash // stateHash
      )
    );
    // transferAllAssetsL  ⬛ --------> (X) -> 👩
    gasResults.ETHexitSadLedgerFunded.satp.transferAllAssetsX = await gasUsed(
      await nitroAdjudicator.transferAllAssets(
        X.channelId,
        proof.outcome, // outcome
        proof.stateHash // stateHash
      )
    );
    // transferAllAssetsX  ⬛ ---------------> 👩

    // record total
    gasResults.ETHexitSadLedgerFunded.satp.total =
      gasResults.ETHexitSadLedgerFunded.satp.challengeL +
      gasResults.ETHexitSadLedgerFunded.satp.transferAllAssetsL +
      gasResults.ETHexitSadLedgerFunded.satp.challengeX +
      gasResults.ETHexitSadLedgerFunded.satp.transferAllAssetsX;
  });

  // exiting a virtual funded (with ETH) channel
  // Scenario: Intermediary Ingrid goes offline
  // initially                   ⬛ ->  L  ->  V  -> 👩
  // challenge L,V   + timeout   ⬛ -> (L) -> (V) -> 👩
  // reclaim L                   ⬛ -- (L) --------> 👩
  // transferAllAssetsL          ⬛ ---------------> 👩
  await executeAndRevert(async () => {
    // begin setup
    await (
      await nitroAdjudicator.deposit(MAGIC_ADDRESS_INDICATING_ETH, LforV.channelId, 0, 10, {
        value: 10,
      })
    ).wait();
    // end setup
    // initially                   ⬛ ->  L  ->  V  -> 👨
    // challenge L
    const {
      challengeTx: ledgerChallengeTx,
      proof: ledgerProof,
      finalizesAt: ledgerFinalizesAt,
    } = await challengeChannel(LforV, MAGIC_ADDRESS_INDICATING_ETH);
    gasResults.ETHexitSadVirtualFunded.satp.challengeL = await gasUsed(ledgerChallengeTx);

    // challenge V ...
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
    gasResults.ETHexitSadVirtualFunded.satp.challengeV = vGasUsed;

    // begin wait
    await waitForChallengesToTimeOut([ledgerFinalizesAt, vFinalizesAt]);
    // end wait
    // challenge L,V   + timeout   ⬛ -> (L) -> (V) -> 👨

    gasResults.ETHexitSadVirtualFunded.satp.reclaimL = await gasUsed(
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
    );
    // reclaim L                   ⬛ -- (L) --------> 👨

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
    gasResults.ETHexitSadVirtualFunded.satp.transferAllAssetsL = await gasUsed(
      await nitroAdjudicator.transferAllAssets(
        LforV.channelId,
        updatedOutcome,
        ledgerProof.stateHash // stateHash
      )
    );
    // transferAllAssetsL          ⬛ ---------------> 👨

    // record total
    gasResults.ETHexitSadVirtualFunded.satp.total =
      gasResults.ETHexitSadVirtualFunded.satp.challengeL +
      gasResults.ETHexitSadVirtualFunded.satp.challengeV +
      gasResults.ETHexitSadVirtualFunded.satp.reclaimL +
      gasResults.ETHexitSadVirtualFunded.satp.transferAllAssetsL;
  });

  // Scenario: Clearing a challenge with a challenge response
  // initially                   ⬛ -> X -> 👩
  // challenge X                 ⬛ -> (X) -> 👩
  // challenge X                 ⬛ -> (X) -> 👩
  await executeAndRevert(async () => {
    await (
      await nitroAdjudicator.deposit(MAGIC_ADDRESS_INDICATING_ETH, X.channelId, 0, 10, {value: 10})
    ).wait();

    await challengeChannel(X, MAGIC_ADDRESS_INDICATING_ETH);

    const {challengeTx} = await challengeChannel(X, MAGIC_ADDRESS_INDICATING_ETH, true);
    gasResults.ETHClearChallenge.satp.challengeResponseX = await gasUsed(challengeTx);
  });

  // Scenario: Clearing a challenge with a checkpoint response
  // initially                   ⬛ -> X -> 👩
  // challenge X                 ⬛ -> (X) -> 👩
  // checkpoint X                ⬛ -> X -> 👩
  await executeAndRevert(async () => {
    await (
      await nitroAdjudicator.deposit(MAGIC_ADDRESS_INDICATING_ETH, X.channelId, 0, 10, {value: 10})
    ).wait();

    await challengeChannel(X, MAGIC_ADDRESS_INDICATING_ETH);

    const {checkpointTx} = await checkpointChannel(X, MAGIC_ADDRESS_INDICATING_ETH);
    gasResults.ETHClearChallenge.satp.checkpointX = await gasUsed(checkpointTx);
  });

  writeFileSync(__dirname + '/gasResults.json', JSON.stringify(gasResults, null, 2));
  console.log('Benchmark results updated successfully!');
}

main().catch(error => {
  console.error(error);
  process.exitCode = 1;
});

import {writeFileSync} from 'fs';

import {MAGIC_ADDRESS_INDICATING_ETH} from '../src/transactions';

import {
  waitForChallengesToTimeOut,
  challengeChannel,
  Y,
  X,
  LforX,
  gasUsed,
  executeAndRevert,
} from './fixtures';
import {emptyGasResults} from './gas';
import {deployContracts, nitroAdjudicator, token} from './localSetup';

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

    // meta-test here to confirm the total recorded in gas.ts is up to date
    // with the recorded costs of each step
    gasResults.ETHexitSadLedgerFunded.satp.total =
      gasResults.ETHexitSadLedgerFunded.satp.challengeL +
      gasResults.ETHexitSadLedgerFunded.satp.transferAllAssetsL +
      gasResults.ETHexitSadLedgerFunded.satp.challengeX +
      gasResults.ETHexitSadLedgerFunded.satp.transferAllAssetsX;
  });

  // TODO uncomment this test when VirtualApp is developed
  // Scenario: Intermediary Ingrid goes offline
  // initially                   ⬛ ->  L  ->  V  -> 👩
  // challenge L,V   + timeout   ⬛ -> (L) -> (V) -> 👩
  // reclaim L                   ⬛ -- (L) --------> 👩
  // transferAllAssetsL          ⬛ ---------------> 👩
  // exiting a virtual funded (with ETH) channel
  // await executeAndRevert(async () => {
  //   // begin setup
  //   await (
  //     await nitroAdjudicator.deposit(MAGIC_ADDRESS_INDICATING_ETH, LforJ.channelId, 0, 10, {
  //       value: 10,
  //     })
  //   ).wait();
  //   // end setup
  //   // initially                   ⬛ ->  L  ->  J  ->  X  -> 👩
  //   // challenge L
  //   const {
  //     challengeTx: ledgerChallengeTx,
  //     proof: ledgerProof,
  //     finalizesAt: ledgerFinalizesAt,
  //   } = await challengeChannel(LforJ, MAGIC_ADDRESS_INDICATING_ETH);
  //   gasResults.ETHexitSadVirtualFunded.satp.challengeL = await gasUsed(ledgerChallengeTx);

  //   // challenge J
  //   const {
  //     challengeTx: jointChallengeTx,
  //     proof: jointProof,
  //     finalizesAt: jointChannelFinalizesAt,
  //   } = await challengeChannel(J, MAGIC_ADDRESS_INDICATING_ETH);
  //   gasResults.ETHexitSadVirtualFunded.satp.challengeJ = await gasUsed(jointChallengeTx);

  //   // challenge X
  //   const {challengeTx, proof, finalizesAt} = await challengeChannel(
  //     X,
  //     MAGIC_ADDRESS_INDICATING_ETH
  //   );
  //   gasResults.ETHexitSadVirtualFunded.satp.challengeX = await gasUsed(challengeTx);

  //   // begin wait
  //   await waitForChallengesToTimeOut([ledgerFinalizesAt, jointChannelFinalizesAt, finalizesAt]);
  //   // end wait
  //   // challenge L,J,X + timeout   ⬛ -> (L) -> (J) -> (X) -> 👩
  //   await assertEthBalancesAndHoldings(
  //     {Alice: 0, Bob: 0, Ingrid: 0},
  //     {LforJ: amountForAliceAndBob, J: 0, X: 0}
  //   );

  //   gasResults.ETHexitSadVirtualFunded.satp.claimL = await gasUsed(
  //     await nitroAdjudicator.claim({
  //       sourceChannelId: LforJ.channelId,
  //       sourceStateHash: ledgerProof.stateHash,
  //       sourceOutcomeBytes: encodeOutcome(ledgerProof.outcome),
  //       sourceAssetIndex: 0,
  //       indexOfTargetInSource: 0,
  //       targetStateHash: jointProof.stateHash,
  //       targetOutcomeBytes: encodeOutcome(jointProof.outcome),
  //       targetAssetIndex: 0,
  //       targetAllocationIndicesToPayout: [], // meaning "all"
  //     })
  //   );
  //   // claimL                      ⬛ ---------------> (X) -> 👩
  //   await assertEthBalancesAndHoldings(
  //     {Alice: 0, Bob: 0, Ingrid: 0},
  //     {LforJ: 0, J: 0, X: amountForAliceAndBob}
  //   );

  //   gasResults.ETHexitSadVirtualFunded.satp.transferAllAssetsX = await gasUsed(
  //     await nitroAdjudicator.transferAllAssets(
  //       X.channelId,
  //       proof.outcome, // outcomeBytes
  //       proof.stateHash // stateHash
  //     )
  //   );
  //   // transferAllAssetsX          ⬛ ----------------------> 👩
  //   await assertEthBalancesAndHoldings(
  //     {Alice: amountForAlice, Bob: amountForBob, Ingrid: 0},
  //     {LforJ: 0, J: 0, X: 0}
  //   );

  //   // meta-test here to confirm the total recorded in gas.ts is up to date
  //   // with the recorded costs of each step
  //   gasResults.ETHexitSadVirtualFunded.satp.total =
  //     (Object.values(gasResults.ETHexitSadVirtualFunded.satp) as number[]).reduce((a, b) => a + b) -
  //     gasResults.ETHexitSadVirtualFunded.satp.total;
  // });

  writeFileSync(__dirname + '/gasResults.json', JSON.stringify(gasResults, null, 2));
}

main().catch(error => {
  console.error(error);
  process.exitCode = 1;
});

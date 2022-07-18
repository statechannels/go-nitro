import {writeFileSync} from 'fs';

import {encodeOutcome} from '../src';
import {MAGIC_ADDRESS_INDICATING_ETH} from '../src/transactions';

import {
  waitForChallengesToTimeOut,
  challengeChannel,
  Y,
  X,
  LforX,
  LforJ,
  J,
  assertEthBalancesAndHoldings,
  amountForAlice,
  amountForBob,
  amountForAliceAndBob,
  gasUsed,
} from './fixtures';
import {gasRequiredTo} from './gas';
import {nitroAdjudicator, token} from './vanillaSetup';

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

afterAll(async () => {
  writeFileSync(__dirname + '/gasResults.json', JSON.stringify(gasRequiredTo, null, 2));
});

describe('Consumes the expected gas for deployments', () => {
  it(`when deploying the NitroAdjudicator`, async () => {
    gasRequiredTo.deployInfrastructureContracts.satp.NitroAdjudicator = await gasUsed(
      await nitroAdjudicator.deployTransaction
    );
  });
});
describe('Consumes the expected gas for deposits', () => {
  it(`when directly funding a channel with ETH (first deposit)`, async () => {
    gasRequiredTo.directlyFundAChannelWithETHFirst.satp = await gasUsed(
      await nitroAdjudicator.deposit(MAGIC_ADDRESS_INDICATING_ETH, X.channelId, 0, 5, {value: 5})
    );
  });

  it(`when directly funding a channel with ETH (second deposit)`, async () => {
    // begin setup
    const setupTX = nitroAdjudicator.deposit(MAGIC_ADDRESS_INDICATING_ETH, X.channelId, 0, 5, {
      value: 5,
    });
    await (await setupTX).wait();
    // end setup
    gasRequiredTo.directlyFundAChannelWithETHSecond.satp = await gasUsed(
      await nitroAdjudicator.deposit(MAGIC_ADDRESS_INDICATING_ETH, X.channelId, 5, 5, {value: 5})
    );
  });

  it(`when directly funding a channel with an ERC20 (first deposit)`, async () => {
    // begin setup
    await (await token.transfer(nitroAdjudicator.address, 1)).wait(); // The asset holder already has some tokens (for other channels)
    // end setup
    gasRequiredTo.directlyFundAChannelWithERC20First.satp.approve = await gasUsed(
      await token.increaseAllowance(nitroAdjudicator.address, 100)
    );
    gasRequiredTo.directlyFundAChannelWithERC20First.satp.deposit = await gasUsed(
      await nitroAdjudicator.deposit(token.address, X.channelId, 0, 5)
    );
  });

  it(`when directly funding a channel with an ERC20 (second deposit)`, async () => {
    // begin setup
    await (await token.increaseAllowance(nitroAdjudicator.address, 100)).wait();
    await (await nitroAdjudicator.deposit(token.address, X.channelId, 0, 5)).wait(); // The asset holder already has some tokens *for this channel*
    await (await token.decreaseAllowance(nitroAdjudicator.address, 95)).wait(); // reset allowance to zero
    // end setup
    gasRequiredTo.directlyFundAChannelWithERC20Second.satp.approve = await gasUsed(
      await token.increaseAllowance(nitroAdjudicator.address, 100)
    );
    gasRequiredTo.directlyFundAChannelWithERC20Second.satp.deposit = await gasUsed(
      await nitroAdjudicator.deposit(token.address, X.channelId, 5, 5)
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
    gasRequiredTo.ETHexit.satp = await gasUsed(
      await X.concludeAndTransferAllAssetsTx(MAGIC_ADDRESS_INDICATING_ETH)
    );
  });

  it(`when exiting a directly funded (with ERC20s) channel`, async () => {
    // begin setup
    await (await token.increaseAllowance(nitroAdjudicator.address, 100)).wait();
    await (await nitroAdjudicator.deposit(token.address, X.channelId, 0, 10)).wait();
    await addResidualTokenBalance(token.address);
    // end setup
    gasRequiredTo.ERC20exit.satp = await gasUsed(
      await X.concludeAndTransferAllAssetsTx(token.address)
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
    const {challengeTx, proof, finalizesAt} = await challengeChannel(
      X,
      MAGIC_ADDRESS_INDICATING_ETH
    );
    gasRequiredTo.ETHexitSad.satp.challenge = await gasUsed(challengeTx);
    // begin wait
    await waitForChallengesToTimeOut([finalizesAt]);
    // end wait
    // challenge + timeout       ⬛ -> (X) -> 👩\
    gasRequiredTo.ETHexitSad.satp.transferAllAssets = await gasUsed(
      await nitroAdjudicator.transferAllAssets(
        X.channelId,
        proof.outcome, // outcome,
        proof.stateHash // stateHash
      )
    );
    // transferAllAssets ⬛ --------> 👩
    gasRequiredTo.ETHexitSad.satp.total =
      gasRequiredTo.ETHexitSad.satp.challenge + gasRequiredTo.ETHexitSad.satp.transferAllAssets;
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
    const {
      challengeTx: ledgerChallengeTx,
      proof: ledgerProof,
      finalizesAt: ledgerFinalizesAt,
    } = await challengeChannel(LforX, MAGIC_ADDRESS_INDICATING_ETH);
    gasRequiredTo.ETHexitSadLedgerFunded.satp.challengeL = await gasUsed(ledgerChallengeTx);

    const {challengeTx, proof, finalizesAt} = await challengeChannel(
      X,
      MAGIC_ADDRESS_INDICATING_ETH
    );
    gasRequiredTo.ETHexitSadLedgerFunded.satp.challengeX = await gasUsed(challengeTx);

    // begin wait
    await waitForChallengesToTimeOut([ledgerFinalizesAt, finalizesAt]); // just go to the max one
    // end wait
    // challenge X, L and timeout  ⬛ -> (L) -> (X) -> 👩
    gasRequiredTo.ETHexitSadLedgerFunded.satp.transferAllAssetsL = await gasUsed(
      await nitroAdjudicator.transferAllAssets(
        LforX.channelId,
        ledgerProof.outcome, // outcome
        ledgerProof.stateHash // stateHash
      )
    );
    // transferAllAssetsL  ⬛ --------> (X) -> 👩
    gasRequiredTo.ETHexitSadLedgerFunded.satp.transferAllAssetsX = await gasUsed(
      await nitroAdjudicator.transferAllAssets(
        X.channelId,
        proof.outcome, // outcome
        proof.stateHash // stateHash
      )
    );
    // transferAllAssetsX  ⬛ ---------------> 👩

    // meta-test here to confirm the total recorded in gas.ts is up to date
    // with the recorded costs of each step
    gasRequiredTo.ETHexitSadLedgerFunded.satp.total =
      gasRequiredTo.ETHexitSadLedgerFunded.satp.challengeL +
      gasRequiredTo.ETHexitSadLedgerFunded.satp.transferAllAssetsL +
      gasRequiredTo.ETHexitSadLedgerFunded.satp.challengeX +
      gasRequiredTo.ETHexitSadLedgerFunded.satp.transferAllAssetsX;
  });

  // TODO unskip this test when the contracts are satp compatible
  it.skip(`when exiting a virtual funded (with ETH) channel`, async () => {
    // begin setup
    await (
      await nitroAdjudicator.deposit(MAGIC_ADDRESS_INDICATING_ETH, LforJ.channelId, 0, 10, {
        value: 10,
      })
    ).wait();
    // end setup
    // initially                   ⬛ ->  L  ->  J  ->  X  -> 👩
    // challenge L
    const {
      challengeTx: ledgerChallengeTx,
      proof: ledgerProof,
      finalizesAt: ledgerFinalizesAt,
    } = await challengeChannel(LforJ, MAGIC_ADDRESS_INDICATING_ETH);
    gasRequiredTo.ETHexitSadVirtualFunded.satp.challengeL = await gasUsed(ledgerChallengeTx);

    // challenge J
    const {
      challengeTx: jointChallengeTx,
      proof: jointProof,
      finalizesAt: jointChannelFinalizesAt,
    } = await challengeChannel(J, MAGIC_ADDRESS_INDICATING_ETH);
    gasRequiredTo.ETHexitSadVirtualFunded.satp.challengeJ = await gasUsed(jointChallengeTx);

    // challenge X
    const {challengeTx, proof, finalizesAt} = await challengeChannel(
      X,
      MAGIC_ADDRESS_INDICATING_ETH
    );
    gasRequiredTo.ETHexitSadVirtualFunded.satp.challengeX = await gasUsed(challengeTx);

    // begin wait
    await waitForChallengesToTimeOut([ledgerFinalizesAt, jointChannelFinalizesAt, finalizesAt]);
    // end wait
    // challenge L,J,X + timeout   ⬛ -> (L) -> (J) -> (X) -> 👩
    await assertEthBalancesAndHoldings(
      {Alice: 0, Bob: 0, Ingrid: 0},
      {LforJ: amountForAliceAndBob, J: 0, X: 0}
    );

    gasRequiredTo.ETHexitSadVirtualFunded.satp.claimL = await gasUsed(
      await nitroAdjudicator.claim({
        sourceChannelId: LforJ.channelId,
        sourceStateHash: ledgerProof.stateHash,
        sourceOutcomeBytes: encodeOutcome(ledgerProof.outcome),
        sourceAssetIndex: 0,
        indexOfTargetInSource: 0,
        targetStateHash: jointProof.stateHash,
        targetOutcomeBytes: encodeOutcome(jointProof.outcome),
        targetAssetIndex: 0,
        targetAllocationIndicesToPayout: [], // meaning "all"
      })
    );
    // claimL                      ⬛ ---------------> (X) -> 👩
    await assertEthBalancesAndHoldings(
      {Alice: 0, Bob: 0, Ingrid: 0},
      {LforJ: 0, J: 0, X: amountForAliceAndBob}
    );

    gasRequiredTo.ETHexitSadVirtualFunded.satp.transferAllAssetsX = await gasUsed(
      await nitroAdjudicator.transferAllAssets(
        X.channelId,
        proof.outcome, // outcomeBytes
        proof.stateHash // stateHash
      )
    );
    // transferAllAssetsX          ⬛ ----------------------> 👩
    await assertEthBalancesAndHoldings(
      {Alice: amountForAlice, Bob: amountForBob, Ingrid: 0},
      {LforJ: 0, J: 0, X: 0}
    );

    // meta-test here to confirm the total recorded in gas.ts is up to date
    // with the recorded costs of each step
    gasRequiredTo.ETHexitSadVirtualFunded.satp.total =
      (Object.values(gasRequiredTo.ETHexitSadVirtualFunded.satp) as number[]).reduce(
        (a, b) => a + b
      ) - gasRequiredTo.ETHexitSadVirtualFunded.satp.total;
  });
});

import {expectRevert} from '@statechannels/devtools';
import {Contract, constants, BigNumber} from 'ethers';
import {it} from '@jest/globals';
import {Allocation, AllocationType} from '@statechannels/exit-format';

import {
  getTestProvider,
  randomChannelId,
  randomExternalDestination,
  replaceAddressesAndBigNumberify,
  setupContract,
  AssetOutcomeShortHand,
} from '../../test-helpers';
import {TESTNitroAdjudicator} from '../../../typechain-types/TESTNitroAdjudicator';
// eslint-disable-next-line import/order
import TESTNitroAdjudicatorArtifact from '../../../artifacts/contracts/test/TESTNitroAdjudicator.sol/TESTNitroAdjudicator.json';
import {
  channelDataToStatus,
  convertAddressToBytes32,
  convertBytes32ToAddress,
  encodeOutcome,
  hashOutcome,
  Outcome,
} from '../../../src';
import {MAGIC_ADDRESS_INDICATING_ETH} from '../../../src/transactions';
import {encodeGuaranteeData} from '../../../src/contract/outcome';
const provider = getTestProvider();

const testNitroAdjudicator: TESTNitroAdjudicator & Contract = setupContract(
  provider,
  TESTNitroAdjudicatorArtifact,
  process.env.TEST_NITRO_ADJUDICATOR_ADDRESS
) as unknown as TESTNitroAdjudicator & Contract;

// Amounts are valueString representations of wei
describe('reclaim', () => {
  it('handles a simpe case as expected', async () => {
    const targetId = randomChannelId();
    const sourceId = randomChannelId();
    const Alice = randomExternalDestination();
    const Bob = randomExternalDestination();
    const Irene = randomExternalDestination();

    // prepare an appropriate virtual channel outcome and finalize

    const vAllocations: Allocation[] = [
      {
        destination: Alice,
        amount: BigNumber.from(7).toHexString(),
        allocationType: AllocationType.simple,
        metadata: '0x',
      },
      {
        destination: Bob,
        amount: BigNumber.from(3).toHexString(),
        allocationType: AllocationType.simple,
        metadata: '0x',
      },
    ];

    const vOutcome: Outcome = [
      {asset: MAGIC_ADDRESS_INDICATING_ETH, allocations: vAllocations, metadata: '0x'},
    ];
    const vOutcomeHash = hashOutcome(vOutcome);
    await (
      await testNitroAdjudicator.setStatusFromChannelData(targetId, {
        turnNumRecord: 99,
        finalizesAt: 0,
        stateHash: constants.HashZero, // not realistic, but OK for purpose of this test
        outcomeHash: vOutcomeHash,
      })
    ).wait();

    // prepare an appropriate ledger channel outcome and finalize

    const lAllocations: Allocation[] = [
      {
        destination: Alice,
        amount: BigNumber.from(10).toHexString(),
        allocationType: AllocationType.simple,
        metadata: '0x',
      },
      {
        destination: Irene,
        amount: BigNumber.from(10).toHexString(),
        allocationType: AllocationType.simple,
        metadata: '0x',
      },
      {
        destination: targetId,
        amount: BigNumber.from(10).toHexString(),
        allocationType: AllocationType.guarantee,
        metadata: encodeGuaranteeData({left: Alice, right: Irene}),
      },
    ];

    const lOutcome: Outcome = [
      {asset: MAGIC_ADDRESS_INDICATING_ETH, allocations: lAllocations, metadata: '0x'},
    ];
    const lOutcomeHash = hashOutcome(lOutcome);
    await (
      await testNitroAdjudicator.setStatusFromChannelData(sourceId, {
        turnNumRecord: 99,
        finalizesAt: 0,
        stateHash: constants.HashZero, // not realistic, but OK for purpose of this test
        outcomeHash: lOutcomeHash,
      })
    ).wait();

    // call reclaim

    const tx = testNitroAdjudicator.reclaim({
      sourceChannelId: sourceId,
      sourceStateHash: constants.HashZero,
      sourceOutcomeBytes: encodeOutcome(lOutcome),
      sourceAssetIndex: 0, // TODO: introduce test cases with multiple-asset Source and Targets
      indexOfTargetInSource: 2,
      targetStateHash: constants.HashZero,
      targetOutcomeBytes: encodeOutcome(vOutcome),
      targetAssetIndex: 0,
    });

    // Extract logs
    const {events: eventsFromTx} = await (await tx).wait();

    // Compile event expectations

    // Check that each expectedEvent is contained as a subset of the properies of each *corresponding* event: i.e. the order matters!
    expect(eventsFromTx).toMatchObject([]);

    // assert on updated ledger channel

    // Check new outcomeHash
    const allocationAfter: Allocation[] = [];
    const outcomeAfter: Outcome = [
      {asset: MAGIC_ADDRESS_INDICATING_ETH, allocations: allocationAfter, metadata: '0x'},
    ];
    const expectedStatusAfter = channelDataToStatus({
      turnNumRecord: 99,
      finalizesAt: 0,
      // stateHash will be set to HashZero by this helper fn
      // if state property of this object is undefined
      outcome: outcomeAfter,
    });
    expect(await testNitroAdjudicator.statusOf(sourceId)).toEqual(expectedStatusAfter);

    // assert that virtual channel did not change.

    expect(await testNitroAdjudicator.statusOf(targetId)).toEqual({
      turnNumRecord: 99,
      finalizesAt: 0,
      stateHash: constants.HashZero, // not realistic, but OK for purpose of this test
      outcomeHash: vOutcomeHash,
    });
  });
});

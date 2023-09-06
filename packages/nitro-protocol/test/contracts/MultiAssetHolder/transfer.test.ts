import {expectRevert} from '@statechannels/devtools';
import {BigNumber, BigNumberish, constants, Contract} from 'ethers';
import {it} from '@jest/globals';
import {Allocation, AllocationType} from '@statechannels/exit-format';

import {
  getTestProvider,
  randomChannelId,
  randomExternalDestination,
  setupContract,
} from '../../test-helpers';
import {encodeOutcome, hashOutcome, Outcome} from '../../../src/contract/outcome';
import {TESTNitroAdjudicator} from '../../../typechain-types/TESTNitroAdjudicator';
// eslint-disable-next-line import/order
import TESTNitroAdjudicatorArtifact from '../../../artifacts/contracts/test/TESTNitroAdjudicator.sol/TESTNitroAdjudicator.json';
import {channelDataToStatus, isExternalDestination} from '../../../src';
import {MAGIC_ADDRESS_INDICATING_ETH} from '../../../src/transactions';
import {AssetOutcomeShortHand, OutcomeShortHand, replaceAddressesAndBigNumberify} from '../../../src/helpers';

const testProvider = getTestProvider();

const testNitroAdjudicator = setupContract(
  testProvider,
  TESTNitroAdjudicatorArtifact,
  process.env.TEST_NITRO_ADJUDICATOR_ADDRESS ||""
) as unknown as TESTNitroAdjudicator & Contract;

const addresses = {
  // Channels
  c: undefined as string | undefined,
  C: randomChannelId(),
  X: randomChannelId(),
  // Externals
  A: randomExternalDestination(),
  B: randomExternalDestination(),
};

const reason0 = 'Channel not finalized';
const reason1 = 'Indices must be sorted';
const reason2 = 'incorrect fingerprint';
const reason3 = 'cannot transfer a guarantee';

// c is the channel we are transferring from.
describe('transfer', () => {
  it.each`
    name                                   | heldBefore | isSimple | setOutcome            | indices      | newOutcome            | heldAfter       | payouts         | reason
    ${' 0. channel not finalized        '} | ${{c: 1}}  | ${true}  | ${{}}                 | ${[0]}       | ${{}}                 | ${{}}           | ${{A: 1}}       | ${reason0}
    ${' 1. funded          -> 1 EOA'}      | ${{c: 1}}  | ${true}  | ${{A: 1}}             | ${[0]}       | ${{A: 0}}             | ${{}}           | ${{A: 1}}       | ${undefined}
    ${' 2. overfunded      -> 1 EOA'}      | ${{c: 2}}  | ${true}  | ${{A: 1}}             | ${[0]}       | ${{A: 0}}             | ${{c: 1}}       | ${{A: 1}}       | ${undefined}
    ${' 3. underfunded     -> 1 EOA'}      | ${{c: 1}}  | ${true}  | ${{A: 2}}             | ${[0]}       | ${{A: 1}}             | ${{}}           | ${{A: 1}}       | ${undefined}
    ${' 4. funded      -> 1 channel'}      | ${{c: 1}}  | ${true}  | ${{C: 1}}             | ${[0]}       | ${{C: 0}}             | ${{c: 0, C: 1}} | ${{}}           | ${undefined}
    ${' 5. overfunded  -> 1 channel'}      | ${{c: 2}}  | ${true}  | ${{C: 1}}             | ${[0]}       | ${{C: 0}}             | ${{c: 1, C: 1}} | ${{}}           | ${undefined}
    ${' 6. underfunded -> 1 channel'}      | ${{c: 1}}  | ${true}  | ${{C: 2}}             | ${[0]}       | ${{C: 1}}             | ${{c: 0, C: 1}} | ${{}}           | ${undefined}
    ${' 7. -> 2 EOA         1 index'}      | ${{c: 2}}  | ${true}  | ${{A: 1, B: 1}}       | ${[0]}       | ${{A: 0, B: 1}}       | ${{c: 1}}       | ${{A: 1}}       | ${undefined}
    ${' 8. -> 2 EOA         1 index'}      | ${{c: 1}}  | ${true}  | ${{A: 1, B: 1}}       | ${[0]}       | ${{A: 0, B: 1}}       | ${{c: 0}}       | ${{A: 1}}       | ${undefined}
    ${' 9. -> 2 EOA         partial'}      | ${{c: 3}}  | ${true}  | ${{A: 2, B: 2}}       | ${[1]}       | ${{A: 2, B: 1}}       | ${{c: 2}}       | ${{B: 1}}       | ${undefined}
    ${'10. -> 2 chan             no'}      | ${{c: 1}}  | ${true}  | ${{C: 1, X: 1}}       | ${[1]}       | ${{C: 1, X: 1}}       | ${{c: 1}}       | ${{}}           | ${undefined}
    ${'11. -> 2 chan           full'}      | ${{c: 1}}  | ${true}  | ${{C: 1, X: 1}}       | ${[0]}       | ${{C: 0, X: 1}}       | ${{c: 0, C: 1}} | ${{}}           | ${undefined}
    ${'12. -> 2 chan        partial'}      | ${{c: 3}}  | ${true}  | ${{C: 2, X: 2}}       | ${[1]}       | ${{C: 2, X: 1}}       | ${{c: 2, X: 1}} | ${{}}           | ${undefined}
    ${'13. -> 2 indices'}                  | ${{c: 3}}  | ${true}  | ${{C: 2, X: 2}}       | ${[0, 1]}    | ${{C: 0, X: 1}}       | ${{c: 0, X: 1}} | ${{C: 2}}       | ${undefined}
    ${'14. -> 3 indices'}                  | ${{c: 5}}  | ${true}  | ${{A: 1, C: 2, X: 2}} | ${[0, 1, 2]} | ${{A: 0, C: 0, X: 0}} | ${{c: 0, X: 2}} | ${{A: 1, C: 2}} | ${undefined}
    ${'15. -> reverse order (see 13)'}     | ${{c: 3}}  | ${true}  | ${{C: 2, X: 2}}       | ${[1, 0]}    | ${{C: 2, X: 1}}       | ${{c: 2, X: 1}} | ${{}}           | ${reason1}
    ${'16. incorrect fingerprint        '} | ${{c: 1}}  | ${true}  | ${{}}                 | ${[0]}       | ${{}}                 | ${{}}           | ${{A: 1}}       | ${reason2}
    ${'17. guarantee allocationType'}      | ${{c: 1}}  | ${false} | ${{A: 1}}             | ${[0]}       | ${{A: 0}}             | ${{}}           | ${{A: 1}}       | ${reason3}
  `(
    `$name: isSimple: $isSimple, heldBefore: $heldBefore, setOutcome: $setOutcome, newOutcome: $newOutcome, heldAfter: $heldAfter, payouts: $payouts`,
    async (tc) => {
      
        let heldBefore = tc.heldBefore as AssetOutcomeShortHand
        let isSimple = tc.isSimple as boolean;
        let setOutcome = tc.setOutcome   as AssetOutcomeShortHand
        let indices = tc.indices as number[];
        let newOutcome = tc.newOutcome  as AssetOutcomeShortHand
        let heldAfter = tc.heldAfter  as AssetOutcomeShortHand
        let payouts = tc.payouts  as AssetOutcomeShortHand
        let reason = tc.reason as string;
      // Compute channelId
      addresses.c = randomChannelId();
      const channelId = addresses.c;
      addresses.C = randomChannelId();
      addresses.X = randomChannelId();
      addresses.A = randomExternalDestination();
      addresses.B = randomExternalDestination();

      // Transform input data (unpack addresses and BigNumberify amounts)
      heldBefore = replaceAddressesAndBigNumberify(heldBefore, addresses) as AssetOutcomeShortHand;
      setOutcome = replaceAddressesAndBigNumberify(setOutcome, addresses) as AssetOutcomeShortHand;
      newOutcome = replaceAddressesAndBigNumberify(newOutcome, addresses) as AssetOutcomeShortHand;
      heldAfter = replaceAddressesAndBigNumberify(heldAfter, addresses) as AssetOutcomeShortHand;
      payouts = replaceAddressesAndBigNumberify(payouts, addresses) as AssetOutcomeShortHand;

      // Deposit into channels

      await Promise.all(
        Object.keys(heldBefore).map(async key => {
          // Key must be either in heldBefore or heldAfter or both
          const amount = heldBefore[key];
          await (
            await testNitroAdjudicator.deposit(MAGIC_ADDRESS_INDICATING_ETH, key, 0, amount, {
              value: amount,
            })
          ).wait();
          expect(
            (await testNitroAdjudicator.holdings(MAGIC_ADDRESS_INDICATING_ETH, key)).eq(amount)
          ).toBe(true);
        })
      );

      // Compute an appropriate allocation.
      const allocations: Allocation[] = [];
      Object.keys(setOutcome).forEach(key =>
        allocations.push({
          destination: key,
          amount: BigNumber.from(setOutcome[key]).toHexString(),
          metadata: '0x',
          allocationType: isSimple ? AllocationType.simple : AllocationType.guarantee,
        })
      );
      const outcomeHash = hashOutcome([
        {
          asset: MAGIC_ADDRESS_INDICATING_ETH,
          assetMetadata: {assetType: 0, metadata: '0x'},
          allocations,
        },
      ]);
      const outcomeBytes = encodeOutcome([
        {
          asset: MAGIC_ADDRESS_INDICATING_ETH,
          assetMetadata: {assetType: 0, metadata: '0x'},
          allocations,
        },
      ]);

      // Set adjudicator status
      const stateHash = constants.HashZero; // not realistic, but OK for purpose of this test
      const finalizesAt = 42;
      const turnNumRecord = 7;

      if (reason != 'Channel not finalized') {
        await (
          await testNitroAdjudicator.setStatusFromChannelData(channelId, {
            turnNumRecord,
            finalizesAt,
            stateHash,
            outcomeHash,
          })
        ).wait();
      }

      const tx = testNitroAdjudicator.transfer(
        MAGIC_ADDRESS_INDICATING_ETH,
        channelId,
        reason == 'incorrect fingerprint' ? '0xdeadbeef' : outcomeBytes,
        stateHash,
        indices
      );

      // Call method in a slightly different way if expecting a revert
      if (reason) {
        await expectRevert(() => tx, reason);
      } else {
        const {events: eventsFromTx} = await (await tx).wait();
        // Check new holdings
        await Promise.all(
          Object.keys(heldAfter).map(async key =>
            expect(await testNitroAdjudicator.holdings(MAGIC_ADDRESS_INDICATING_ETH, key)).toEqual(
              heldAfter[key]
            )
          )
        );

        // Check new status
        const allocationsAfter: Allocation[] = [];
        Object.keys(newOutcome).forEach(key => {
          allocationsAfter.push({
            destination: key,
            amount: BigNumber.from(newOutcome[key]).toHexString(),
            metadata: '0x',
            allocationType: AllocationType.simple,
          });
        });
        const outcomeAfter: Outcome = [
          {
            asset: MAGIC_ADDRESS_INDICATING_ETH,
            assetMetadata: {assetType: 0, metadata: '0x'},
            allocations: allocationsAfter,
          },
        ];
        const expectedStatusAfter = channelDataToStatus({
          turnNumRecord,
          finalizesAt,
          // stateHash will be set to HashZero by this helper fn
          // if state property of this object is undefined
          outcome: outcomeAfter,
        });
        expect(await testNitroAdjudicator.statusOf(channelId)).toEqual(expectedStatusAfter);

        const expectedEvents = [
          {
            event: 'AllocationUpdated',
            args: {
              channelId,
              assetIndex: BigNumber.from(0),
              initialHoldings: heldBefore[addresses.c],
            },
          },
        ];

        expect(eventsFromTx).toMatchObject(expectedEvents);

        // Check payouts
        for (const destination of Object.keys(payouts)) {
          if (isExternalDestination(destination)) {
            const asAddress = '0x' + destination.substring(26);
            const balance = await testProvider.getBalance(asAddress);
            console.log(`checking balance of ${destination}: ${balance.toString()}`);
            expect(balance).toEqual(payouts[destination]);
          } else {
            const holdings = await testNitroAdjudicator.holdings(
              MAGIC_ADDRESS_INDICATING_ETH,
              destination
            );
            console.log(`checking holdings of ${destination}: ${holdings.toString()}`);
            expect(holdings).toEqual(payouts[destination]);
          }
        }
      }
    }
  );
});

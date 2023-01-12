import {expectRevert} from '@statechannels/devtools';
import {Contract, ethers, BigNumber, constants} from 'ethers';
import {it} from '@jest/globals';

import TokenArtifact from '../../../artifacts/contracts/Token.sol/Token.json';
import {getChannelId} from '../../../src/contract/channel';
import {Outcome} from '../../../src/contract/outcome';
import {
  FixedPart,
  getFixedPart,
  getVariablePart,
  separateProofAndCandidate,
  State,
} from '../../../src/contract/state';
import {
  computeOutcome,
  generateParticipants,
  getCountingAppContractAddress,
  getRandomNonce,
  getTestProvider,
  OutcomeShortHand,
  randomChannelId,
  randomExternalDestination,
  replaceAddressesAndBigNumberify,
  setupContract,
} from '../../test-helpers';
import {signStates, channelDataToStatus, bindSignatures} from '../../../src';
import {MAGIC_ADDRESS_INDICATING_ETH, NITRO_MAX_GAS} from '../../../src/transactions';
import {TESTNitroAdjudicator} from '../../../typechain-types/TESTNitroAdjudicator';
import {Token} from '../../../typechain-types/Token';
// eslint-disable-next-line import/order
import TESTNitroAdjudicatorArtifact from '../../../artifacts/contracts/test/TESTNitroAdjudicator.sol/TESTNitroAdjudicator.json';

const testNitroAdjudicator = setupContract(
  getTestProvider(),
  TESTNitroAdjudicatorArtifact,
  process.env.TEST_NITRO_ADJUDICATOR_ADDRESS
) as unknown as TESTNitroAdjudicator & Contract;

const token = setupContract(
  getTestProvider(),
  TokenArtifact,
  process.env.TEST_TOKEN_ADDRESS
) as unknown as Token & Contract;

const provider = getTestProvider();
const chainId = process.env.CHAIN_NETWORK_ID;

const nParticipants = 3;
const {wallets, participants} = generateParticipants(nParticipants);

const challengeDuration = 0x1000;
let appDefinition: string;

interface addressesT {
  [index: string]: string | undefined;
  At: string;
  Bt: string;
}

const addresses: addressesT = {
  // Channels
  c: undefined,
  C: randomChannelId(),
  X: randomChannelId(),
  // Externals
  A: randomExternalDestination(),
  B: randomExternalDestination(),
  // // Externals preloaded with TOK (cheaper to pay to)
  At: randomExternalDestination(),
  Bt: randomExternalDestination(),
  // Asset Holders
  ETH: undefined,
  ETH2: undefined,
  ERC20: undefined,
};

interface payoutsT {
  [index: string]: number;
}

const tenPayouts = {ERC20: {} as payoutsT};
const fiftyPayouts = {ERC20: {} as payoutsT};
const oneHundredPayouts = {ERC20: {} as payoutsT};
for (let i = 0; i < 100; i++) {
  addresses[i.toString()] =
    '0x000000000000000000000000e0c3b40fdff77c786dd3737837887c85' + (0x2392fa22 + i).toString(16); // they need to be distinct because JS objects
  if (i < 10) tenPayouts.ERC20[i.toString()] = 1;
  if (i < 50) fiftyPayouts.ERC20[i.toString()] = 1;
  if (i < 100) oneHundredPayouts.ERC20[i.toString()] = 1;
}

beforeAll(async () => {
  addresses.ETH = MAGIC_ADDRESS_INDICATING_ETH;
  addresses.ERC20 = token.address;
  appDefinition = getCountingAppContractAddress();
  // Preload At and Bt with TOK
  await (await token.transfer('0x' + addresses.At.slice(26), BigNumber.from(1))).wait();
  await (await token.transfer('0x' + addresses.Bt.slice(26), BigNumber.from(1))).wait();
});

const accepts1 = '{ETH: {A: 1}}';
const accepts2 = '{ETH: {A: 1}}';
const accepts3 = '{ETH: {A: 1, B: 1}}';
const accepts4 = '{ERC20: {A: 1, B: 1}}';
const accepts4a = '{ERC20: {A: 1}}';
const accepts5 = '{ERC20: {At: 1, Bt: 1}} (At and Bt already have some TOK)';
const accepts6 = '10 TOK payouts';
const accepts7 = '50 TOK payouts';
const accepts8 = '100 TOK payouts';

const oneState = {
  whoSignedWhat: [0, 0, 0],
  appData: [ethers.constants.HashZero],
};
const turnNumRecord = 5;
let channelNonce = getRandomNonce('concludeAndTransferAllAssets');
describe('concludeAndTransferAllAssets', () => {
  beforeEach(() => (channelNonce = BigNumber.from(channelNonce).add(1).toHexString()));
  it.each`
    description  | outcomeShortHand           | heldBefore           | heldAfter          | newOutcome | payouts                    | reasonString
    ${accepts1}  | ${{ETH: {A: 1}}}           | ${{ETH: {c: 1}}}     | ${{ETH: {c: 0}}}   | ${{}}      | ${{ETH: {A: 1}}}           | ${undefined}
    ${accepts2}  | ${{ETH: {A: 1}}}           | ${{ETH: {c: 1}}}     | ${{ETH: {c: 0}}}   | ${{}}      | ${{ETH: {A: 1}}}           | ${undefined}
    ${accepts3}  | ${{ETH: {A: 1, B: 1}}}     | ${{ETH: {c: 2}}}     | ${{ETH: {c: 0}}}   | ${{}}      | ${{ETH: {A: 1, B: 1}}}     | ${undefined}
    ${accepts4}  | ${{ERC20: {A: 1, B: 1}}}   | ${{ERC20: {c: 2}}}   | ${{ERC20: {c: 0}}} | ${{}}      | ${{ERC20: {A: 1, B: 1}}}   | ${undefined}
    ${accepts4a} | ${{ERC20: {A: 1}}}         | ${{ERC20: {c: 1}}}   | ${{ERC20: {c: 0}}} | ${{}}      | ${{ERC20: {A: 1}}}         | ${undefined}
    ${accepts5}  | ${{ERC20: {At: 1, Bt: 1}}} | ${{ERC20: {c: 2}}}   | ${{ERC20: {c: 0}}} | ${{}}      | ${{ERC20: {At: 1, Bt: 1}}} | ${undefined}
    ${accepts6}  | ${tenPayouts}              | ${{ERC20: {c: 10}}}  | ${{ERC20: {c: 0}}} | ${{}}      | ${tenPayouts}              | ${undefined}
    ${accepts7}  | ${fiftyPayouts}            | ${{ERC20: {c: 50}}}  | ${{ERC20: {c: 0}}} | ${{}}      | ${fiftyPayouts}            | ${undefined}
    ${accepts8}  | ${oneHundredPayouts}       | ${{ERC20: {c: 100}}} | ${{ERC20: {c: 0}}} | ${{}}      | ${oneHundredPayouts}       | ${undefined}
  `(
    '$description', // For the purposes of this test, chainId and participants are fixed, making channelId 1-1 with channelNonce
    async ({
      outcomeShortHand,
      heldBefore,
      heldAfter,
      newOutcome,
      payouts,
      reasonString,
    }: {
      description: string;
      outcomeShortHand: OutcomeShortHand;
      heldBefore: OutcomeShortHand;
      heldAfter: OutcomeShortHand;
      newOutcome: OutcomeShortHand;
      payouts: OutcomeShortHand;
      reasonString: string;
    }) => {
      const fixedPart: FixedPart = {
        chainId,
        participants,
        channelNonce,
        appDefinition,
        challengeDuration,
      };
      const channelId = getChannelId(fixedPart);
      addresses.c = channelId;
      const support = oneState;
      const {appData, whoSignedWhat} = support;
      const numStates = appData.length;
      const largestTurnNum = turnNumRecord + 1;

      // Transfer some tokens into the relevant AssetHolder
      // Do this step before transforming input data (easier)
      if ('ERC20' in heldBefore) {
        await (
          await token.increaseAllowance(testNitroAdjudicator.address, heldBefore.ERC20.c)
        ).wait();
        await (
          await testNitroAdjudicator.deposit(token.address, channelId, '0x00', heldBefore.ERC20.c)
        ).wait();
      }
      if ('ETH' in heldBefore) {
        await (
          await testNitroAdjudicator.deposit(
            MAGIC_ADDRESS_INDICATING_ETH,
            channelId,
            '0x00',
            heldBefore.ETH.c,
            {
              value: heldBefore.ETH.c,
            }
          )
        ).wait();
      }

      // Transform input data (unpack addresses and BigNumberify amounts)
      [heldBefore, outcomeShortHand, newOutcome, heldAfter, payouts] = [
        heldBefore,
        outcomeShortHand,
        newOutcome,
        heldAfter,
        payouts,
      ].map(object => replaceAddressesAndBigNumberify(object, addresses) as OutcomeShortHand);

      // Compute the outcome.
      const outcome: Outcome = computeOutcome(outcomeShortHand);

      // Construct states
      const states: State[] = [];
      for (let i = 1; i <= numStates; i++) {
        states.push({
          isFinal: true,
          chainId,
          participants,
          channelNonce,
          outcome,
          appDefinition,
          appData: appData[i - 1],
          challengeDuration,
          turnNum: largestTurnNum + i - numStates,
        });
      }

      const variableParts = states.map(state => getVariablePart(state));

      // Sign the states
      const signatures = await signStates(states, wallets, whoSignedWhat);
      const {candidate} = separateProofAndCandidate(
        bindSignatures(variableParts, signatures, whoSignedWhat)
      );

      // Form transaction
      const tx = testNitroAdjudicator.concludeAndTransferAllAssets(
        getFixedPart(states[0]),
        candidate,
        {gasLimit: NITRO_MAX_GAS}
      );

      // Switch on overall test expectation
      if (reasonString) {
        await expectRevert(() => tx, reasonString);
      } else {
        const receipt = await (await tx).wait();

        expect(BigNumber.from(receipt.gasUsed).lt(BigNumber.from(NITRO_MAX_GAS))).toBe(true);

        // Compute expected ChannelDataHash
        const blockTimestamp = (await provider.getBlock(receipt.blockNumber)).timestamp;
        const expectedFingerprint = newOutcome.length
          ? channelDataToStatus({
              turnNumRecord: 0,
              finalizesAt: blockTimestamp,
              outcome: computeOutcome(newOutcome),
            })
          : constants.HashZero;

        // Check fingerprint against the expected value
        expect(await testNitroAdjudicator.statusOf(channelId)).toEqual(expectedFingerprint);

        // Extract logs
        await (await tx).wait();

        // Check new holdings
        await Promise.all(
          // For each asset
          Object.keys(heldAfter).map(async asset => {
            await Promise.all(
              Object.keys(heldAfter[asset]).map(async destination => {
                // for each channel
                const amount = heldAfter[asset][destination];
                expect((await testNitroAdjudicator.holdings(asset, destination)).eq(amount)).toBe(
                  true
                );
              })
            );
          })
        );
      }
    }
  );
});

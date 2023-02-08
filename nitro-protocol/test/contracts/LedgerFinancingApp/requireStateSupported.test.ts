import {expectRevert} from '@statechannels/devtools';
import {Contract, ethers, BigNumber} from 'ethers';
import {ParamType} from 'ethers/lib/utils';

import LedgerFinancingAppArtifact from '../../../artifacts/contracts/LedgerFinancingApp.sol/LedgerFinancingApp.json';
import {computeOutcome, convertAddressToBytes32} from '../../../src';
import {
  getFixedPart,
  getVariablePart,
  RecoveredVariablePart,
  State,
} from '../../../src/contract/state';
import {generateParticipants, getTestProvider, setupContract} from '../../test-helpers';

let ledgerFinancingApp: Contract;
const provider = getTestProvider();

const {participants} = generateParticipants(2);
const challengeDuration = 0x100;
const MAGIC_NATIVE_ASSET_ADDRESS = '0x0000000000000000000000000000000000000000';
const APPDEF = process.env.LEDGER_FINANCING_APP_ADDRESS
  ? process.env.LEDGER_FINANCING_APP_ADDRESS
  : 'failfast';

const intermediary = convertAddressToBytes32(participants[0]);
const merchant = convertAddressToBytes32(participants[1]);

interface Funds {
  asset: string[]; // asset token address
  amount: number[]; // amount of each asset with shared index
}

interface LedgerFinancingAppData {
  dpyNum: number;
  dpyDen: number;
  blocknumber: number;
  principal: Funds;
  collectedInterest: Funds;
}
const ledgerFinancingAppDataTy: ParamType = {
  type: 'tuple',
  components: [
    {type: 'uint128', name: 'dpyNum'},
    {type: 'uint128', name: 'dpyDen'},
    {type: 'uint256', name: 'blocknumber'},
    {
      type: 'tuple',
      name: 'principal',
      components: [
        {type: 'address[]', name: 'asset'},
        {type: 'uint256[]', name: 'amount'},
      ],
    },
    {
      type: 'tuple',
      name: 'collectedInterest',
      components: [
        {type: 'address[]', name: 'asset'},
        {type: 'uint256[]', name: 'amount'},
      ],
    },
  ],
} as ParamType;

function appDataABIEncode(appData: LedgerFinancingAppData): string {
  return ethers.utils.defaultAbiCoder.encode(
    // ['uint128', 'uint128', 'uint256', 'tuple(address[], uint256[])', 'tuple(address[], uint256[])'],
    [ledgerFinancingAppDataTy],
    [appData]
  );
}

const initialOutcome = computeOutcome({
  [MAGIC_NATIVE_ASSET_ADDRESS]: {[intermediary]: 500, [merchant]: 500},
});

const baseAppData: LedgerFinancingAppData = {
  // 1/100 -> 1% daily percentage yield.
  dpyNum: 1,
  dpyDen: 100,
  blocknumber: 1,
  principal: {
    asset: [MAGIC_NATIVE_ASSET_ADDRESS],
    amount: [500],
  },
  collectedInterest: {
    asset: [MAGIC_NATIVE_ASSET_ADDRESS],
    amount: [0],
  },
};

const baseState: State = {
  turnNum: 0,
  isFinal: false, // intermediary wants to force finalization
  channelNonce: '0x8',
  participants,
  challengeDuration,
  outcome: initialOutcome,
  appData: appDataABIEncode(baseAppData),
  appDefinition: APPDEF,
};

const variablePart = getVariablePart(baseState);

const signedByMerchant: RecoveredVariablePart = {
  variablePart,
  signedBy: BigNumber.from(0b10).toHexString(),
};
const signedByIntermediary: RecoveredVariablePart = {
  variablePart,
  signedBy: BigNumber.from(0b01).toHexString(),
};
const signedByBoth: RecoveredVariablePart = {
  variablePart,
  signedBy: BigNumber.from(0b11).toHexString(),
};

const fixedPart = getFixedPart(baseState);

beforeAll(async () => {
  ledgerFinancingApp = setupContract(provider, LedgerFinancingAppArtifact, APPDEF);
});

describe('requireStateSupported', () => {
  it('accepts unanimous states', async () => {
    await ledgerFinancingApp.requireStateSupported(fixedPart, [], signedByBoth);
  });

  it('accepts legitimate interest calculations', async () => {
    // construct a proof+candidate test case with fair interest calculation, assert passing    // with the interest rate applied.
    // test case:
    // - appdata interest rate is 1% per day
    // - initial outcome is 500:500, with 500 principal
    // - 1 day passes
    // - challenge outcome is 505:495
    advanceOneDay();

    const challengeState: State = {
      ...baseState,
      turnNum: baseState.turnNum + 1,
      outcome: computeOutcome({
        [MAGIC_NATIVE_ASSET_ADDRESS]: {[intermediary]: 505, [merchant]: 495}, // intermediary picks up 1% of the principal
      }),
    };
    const challengeWithIntermediarySignature: RecoveredVariablePart = {
      variablePart: getVariablePart(challengeState),
      signedBy: BigNumber.from(0b01).toHexString(),
    };

    await ledgerFinancingApp.requireStateSupported(
      fixedPart,
      [signedByBoth],
      challengeWithIntermediarySignature
    );
  });

  it('reverts if the proof block number is in the future', async () => {
    // create a challenge where the intermediary fakes the blocknumber
    // (in order to get a higher interest rate).
    const futureState: State = {
      ...baseState,
      appData: appDataABIEncode({...baseAppData, blocknumber: baseAppData.blocknumber + 1000000}),
    };
    const signedByIntermediary: RecoveredVariablePart = {
      variablePart: getVariablePart(futureState),
      signedBy: BigNumber.from(0b01).toHexString(),
    };
    await expectRevert(() =>
      ledgerFinancingApp.requireStateSupported(fixedPart, [signedByBoth], signedByIntermediary)
    );
  });

  it('rejects excessive interest calulations', async () => {
    // construct proof+candidate test case with unfair interest calculation, assert failure.
    // test case:
    //  - the appData's interest rate is 1% per day
    //  - the initial outcome is 500:500
    //  - the chain is advanced by 1 day
    //  - challenge outcome is 506:494. Fraud!

    const currentBlockNumber = await provider.getBlockNumber();
    const supportproofAppData = {
      ...baseAppData,
      blocknumber: currentBlockNumber,
    };

    const supportproofState: State = {
      ...baseState,
      appData: appDataABIEncode(supportproofAppData),
    };
    const pfStateSignedByBoth: RecoveredVariablePart = {
      variablePart: getVariablePart(supportproofState),
      signedBy: BigNumber.from(0b11).toHexString(),
    };

    advanceOneDay();

    const challengeState: State = {
      ...baseState,
      appData: appDataABIEncode(supportproofAppData),
      turnNum: baseState.turnNum + 1,
      outcome: computeOutcome({
        [MAGIC_NATIVE_ASSET_ADDRESS]: {[intermediary]: 506, [merchant]: 494}, // 506 is unfair: should be 505
      }),
    };
    const updatedWithIntermediarySignature: RecoveredVariablePart = {
      variablePart: getVariablePart(challengeState),
      signedBy: BigNumber.from(0b01).toHexString(),
    };
    await expectRevert(
      () =>
        ledgerFinancingApp.requireStateSupported(
          fixedPart,
          [pfStateSignedByBoth],
          updatedWithIntermediarySignature
        ),
      'earned<claimed'
    );
  });

  it('rejects unilateral unsupported candidates', async () => {
    // construct challenges with unsupported candidates signed by:
    // - only the intermediary
    // - only the merchant
    // assert failure

    await expectRevert(
      () => ledgerFinancingApp.requireStateSupported(fixedPart, [], signedByMerchant),
      '!unanimous; |proof|=0'
    );
    await expectRevert(
      () => ledgerFinancingApp.requireStateSupported(fixedPart, [], signedByIntermediary),
      '!unanimous; |proof|=0'
    );
  });

  it('rejects unilateral support proof states', async () => {
    // construct challenges where the proof state is signed by
    //  - only the intermediary
    //  - only the merchant
    // assert failure.

    await expectRevert(
      () =>
        ledgerFinancingApp.requireStateSupported(
          fixedPart,
          [signedByMerchant],
          signedByIntermediary
        ),
      '!unanimous proof state'
    );

    await expectRevert(
      () =>
        ledgerFinancingApp.requireStateSupported(
          fixedPart,
          [signedByIntermediary],
          signedByMerchant
        ),
      '!unanimous proof state'
    );
  });

  it('rejects too-long proofs', async () => {
    // construct a challenge with two proof states, assert failure.

    await expectRevert(
      () =>
        ledgerFinancingApp.requireStateSupported(
          fixedPart,
          [signedByMerchant, signedByIntermediary], // two proof states - should fail
          signedByIntermediary
        ),
      '|proof| > 1'
    );
  });
});

/**
 * increase blocknumber on provider by 7500.
 * each day is ~7200= blocks (24*60*60/12)
 */
function advanceOneDay() {
  // note: this is a hacky way to advance so many blocks, and results
  // in a slower test.
  // The 'hardhat_mine' method is better, but it causes a different error
  // in the test.
  for (let i = 0; i < 7500; i++) {
    provider.send('evm_mine', []);
  }
}

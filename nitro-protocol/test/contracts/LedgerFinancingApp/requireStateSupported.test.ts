import {expectRevert} from '@statechannels/devtools';
import {Contract, ethers, BigNumber} from 'ethers';

import LedgerFinancingAppArtifact from '../../../artifacts/contracts/LedgerFinancingApp.sol/LedgerFinancingApp.json';
import {
  computeOutcome,
  convertAddressToBytes32,
  encodeVoucherAmountAndSignature,
  getChannelId,
  signState,
  Voucher,
} from '../../../src';
import {
  getFixedPart,
  getVariablePart,
  RecoveredVariablePart,
  State,
} from '../../../src/contract/state';
import {generateParticipants, getTestProvider, setupContract} from '../../test-helpers';
const {HashZero} = ethers.constants;

let ledgerFinancingApp: Contract;
const provider = getTestProvider();

const {wallets, participants} = generateParticipants(2);
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

interface AppData {
  dpyNum: number;
  dpyDen: number;
  blocknumber: number;
  principal: Funds;
  collectedInterest: Funds;
}

function fundsABIEncode(funds: Funds): string {
  return ethers.utils.defaultAbiCoder.encode(
    ['address[]', 'uint256[]'],
    [funds.asset, funds.amount]
  );
}

function appDataABIEncode(appData: AppData): string {
  return ethers.utils.defaultAbiCoder.encode(
    ['uint256', 'uint256', 'uint256', 'tuple(address[], uint256[])', 'tuple(address[], uint256[])'],
    [
      appData.dpyNum,
      appData.dpyDen,
      appData.blocknumber,
      [appData.principal.asset, appData.principal.amount],
      [appData.collectedInterest.asset, appData.collectedInterest.amount],
    ]
  );
}

const baseState: State = {
  turnNum: 0,
  isFinal: false,
  channelNonce: '0x8',
  participants,
  challengeDuration,
  outcome: [],
  appData: HashZero,
  appDefinition: APPDEF,
};

const fixedPart = getFixedPart(baseState);
const channelId = getChannelId(fixedPart);

beforeAll(async () => {
  ledgerFinancingApp = setupContract(provider, LedgerFinancingAppArtifact, APPDEF);
});

describe('requireStateSupported accepts unanimous states', () => {
  // construct candidate-only test case with unanimous state, assert success.
  const state = baseState;
  const variablePart = getVariablePart(state);
  const unanimousCandidate: RecoveredVariablePart = {
    variablePart,
    signedBy: BigNumber.from(0b11).toHexString(),
  };
  (async () => {
    await ledgerFinancingApp.requireStateSupported(fixedPart, [], unanimousCandidate);
  })();
});

describe('requireStateSupported', () => {
  it('accepts legitimate interest calculations', () => {
    // test case:
    // - proof state w/ some outcome + appdata
    // - candidate state with interest rate
  });

  it('rejects excessive interest calulations', () => {
    // construct proof+candidate test case with unfair interest calculation, assert failure.
  });

  it('rejects unilateral unsupported candidates', () => {
    // test cases:
    // - signed by intermediary only
    // - signed by merchant only
    // proof.length == 0
    // failure for both
    const state = baseState;
    const variablePart = getVariablePart(state);
    const signedByMerchant: RecoveredVariablePart = {
      variablePart,
      signedBy: BigNumber.from(0b10).toHexString(),
    };
    const signedByIntermediary: RecoveredVariablePart = {
      variablePart,
      signedBy: BigNumber.from(0b01).toHexString(),
    };

    (async () => {
      await expectRevert(
        () => ledgerFinancingApp.requireStateSupported(fixedPart, [], signedByMerchant),
        '!unanimous; |proof|=0'
      );
      await expectRevert(
        () => ledgerFinancingApp.requireStateSupported(fixedPart, [], signedByIntermediary),
        '!unanimous; |proof|=0'
      );
    })();
  });

  it('rejects unilateral support proof states', () => {
    // construct proof[0] with:
    // - signed by intermediary only
    // - signed by merchant only
    // assert failure for both

    const state = baseState;
    const variablePart = getVariablePart(state);
    const signedByMerchant: RecoveredVariablePart = {
      variablePart,
      signedBy: BigNumber.from(0b10).toHexString(),
    };
    const signedByIntermediary: RecoveredVariablePart = {
      variablePart,
      signedBy: BigNumber.from(0b01).toHexString(),
    };

    (async () => {
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
    })();
  });

  it('rejects too-long proofs', () => {
    // construct a challenge w/ two proof states, assert failure
    const state = baseState;
    const variablePart = getVariablePart(state);
    const signedByMerchant: RecoveredVariablePart = {
      variablePart,
      signedBy: BigNumber.from(0b10).toHexString(),
    };
    const signedByIntermediary: RecoveredVariablePart = {
      variablePart,
      signedBy: BigNumber.from(0b01).toHexString(),
    };

    (async () => {
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
            [signedByIntermediary, signedByIntermediary],
            signedByMerchant
          ),
        '!unanimous proof state'
      );
    })();
  });
});

import hre from 'hardhat';
import {BigNumber, BigNumberish, constants, ContractReceipt, ethers} from 'ethers';
import {Signature} from '@ethersproject/bytes';
import {Wallet} from '@ethersproject/wallet';
import {AllocationType} from '@statechannels/exit-format';
import {takeSnapshot} from '@nomicfoundation/hardhat-network-helpers';

import {
  Bytes32,
  convertAddressToBytes32,
  encodeVoucherAmountAndSignature,
  getChannelId,
  getFixedPart,
  MAGIC_ADDRESS_INDICATING_ETH,
  signChallengeMessage,
  SignedState,
  signState,
  signVoucher,
  State,
  Voucher,
} from '../src';
import {
  encodeGuaranteeData,
  GuaranteeAllocation,
  Outcome,
  SimpleAllocation,
} from '../src/contract/outcome';
import {FixedPart, getVariablePart, hashState, SignedVariablePart} from '../src/contract/state';

import {nitroAdjudicator, consensusAppAddress, virtualPaymentAppAddress} from './localSetup';

export const chainId = '0x7a69'; // 31337 in hex (hardhat network default)

export const Alice = new Wallet(
  '0x277fb9e0ad81dc836c60294e385b10dfcc0a9586eeb0b1d31da92e384a0d2efa'
);
export const Bob = new Wallet('0xc8774aa98410b3e3281ff1ec40ea2637d2b9280328c4d1ff00d06cd95dd42cbd');
export const Ingrid = new Wallet(
  '0x558789345da13a7ac1d6d6ac9275ba66836eb4a088efc1920db0f5d092d6ee71'
);
export const participants = [Alice.address, Bob.address];

export const amountForAlice = BigNumber.from(10).toHexString();
export const amountForBob = BigNumber.from(0).toHexString(); // We will use a unidirectional payment channel, i.e. Bob starts with nothing!
export const paymentAmount = BigNumber.from(1).toHexString();
export const amountForAliceAndBob = BigNumber.from(amountForAlice).add(amountForBob).toHexString();

export class TestChannel {
  constructor(
    channelNonce: string,
    wallets: ethers.Wallet[],
    allocations: Array<GuaranteeAllocation | SimpleAllocation>,
    appDefinition?: string
  ) {
    this.wallets = wallets;
    this.appDefinition = appDefinition ?? consensusAppAddress;
    this.fixedPart = {
      channelNonce,
      participants: wallets.map(w => w.address),
      appDefinition: this.appDefinition,
      challengeDuration: 600,
    };
    this.allocations = allocations;
  }
  appDefinition: string;
  wallets: ethers.Wallet[];
  fixedPart: FixedPart;
  private allocations: Array<GuaranteeAllocation | SimpleAllocation>;
  outcome(asset: string) {
    const outcome: Outcome = [
      {
        asset,
        allocations: Array.from(this.allocations, a => ({...a})),
        assetMetadata: {assetType: 0, metadata: '0x'},
      },
    ];
    return outcome;
  }
  get channelId() {
    return getChannelId(this.fixedPart);
  }
  someState(asset: string): State {
    return {
      ...this.fixedPart,
      turnNum: 6,
      isFinal: false,
      outcome: this.outcome(asset),
      appData: '0x', // TODO choose a more representative example
    };
  }

  finalState(asset: string): State {
    return {
      ...this.someState(asset),
      isFinal: true,
    };
  }

  counterSignedSupportProof(
    // for challenging and outcome pushing
    state: State
  ): {
    fixedPart: FixedPart;
    proof: SignedVariablePart[];
    candidate: SignedVariablePart;
    challengeSignature: Signature;
    outcome: Outcome;
    stateHash: string;
  } {
    return {
      fixedPart: getFixedPart(state),
      proof: [],
      candidate: {
        variablePart: getVariablePart(state),
        sigs: this.wallets.map(w => signState(state, w.privateKey).signature),
      },
      challengeSignature: signChallengeMessage(
        [{state} as SignedState],
        this.wallets[0].privateKey
      ),
      outcome: state.outcome,
      stateHash: hashState(state),
    };
  }

  supportProof(
    // for concluding
    state: State
  ): {
    fixedPart: FixedPart;
    proof: SignedVariablePart[];
    candidate: SignedVariablePart;
  } {
    return {
      fixedPart: getFixedPart(state),
      proof: [],
      candidate: {
        variablePart: getVariablePart(state),
        sigs: this.wallets.map(w => signState(state, w.privateKey).signature),
      },
    };
  }

  async concludeAndTransferAllAssetsTx(asset: string) {
    const fP = this.supportProof(this.finalState(asset));
    return await nitroAdjudicator.concludeAndTransferAllAssets(fP.fixedPart, fP.candidate);
  }

  async challengeTx(asset: string) {
    const proof = this.counterSignedSupportProof(this.someState(asset));
    return await nitroAdjudicator.challenge(
      proof.fixedPart,
      proof.proof,
      proof.candidate,
      proof.challengeSignature
    );
  }
}

/** An application channel between Alice and Bob */
export const X = new TestChannel(
  '0x2',
  [Alice, Bob],
  [
    {
      destination: convertAddressToBytes32(Alice.address),
      amount: amountForAlice,
      metadata: '0x',
      allocationType: AllocationType.simple,
    },
    {
      destination: convertAddressToBytes32(Bob.address),
      amount: amountForBob,
      metadata: '0x',
      allocationType: AllocationType.simple,
    },
  ]
);

/** Another application channel between Alice and Bob */
export const Y = new TestChannel(
  '0x3',
  [Alice, Bob],
  [
    {
      destination: convertAddressToBytes32(Alice.address),
      amount: amountForAlice,
      metadata: '0x',
      allocationType: AllocationType.simple,
    },
    {
      destination: convertAddressToBytes32(Bob.address),
      amount: amountForBob,
      metadata: '0x',
      allocationType: AllocationType.simple,
    },
  ]
);

/** Ledger channel between Alice and Bob, providing funds to channel X */
export const LforX = new TestChannel(
  '0x4',
  [Alice, Bob],
  [
    {
      destination: X.channelId,
      amount: amountForAliceAndBob,
      metadata: '0x',
      allocationType: AllocationType.simple,
    },
  ]
);

/** Virtual payment channel between Alice and Bob with Ingrid as intermediary*/
export const V = new TestChannel(
  '0x5',
  [Alice, Ingrid, Bob],
  [
    {
      destination: convertAddressToBytes32(Alice.address),
      amount: amountForAlice,
      metadata: '0x',
      allocationType: AllocationType.simple,
    },
    {
      destination: convertAddressToBytes32(Bob.address),
      amount: amountForBob,
      metadata: '0x',
      allocationType: AllocationType.simple,
    },
  ],
  virtualPaymentAppAddress
);

/** Ledger channel between Bob and Ingrid, with Guarantee targeting virtual channel V */
export const LforV = new TestChannel(
  '0x7',
  [Bob, Ingrid],
  [
    {
      destination: convertAddressToBytes32(Bob.address),
      amount: '0x0',
      metadata: '0x',
      allocationType: AllocationType.simple,
    },
    {
      destination: convertAddressToBytes32(Ingrid.address),
      amount: '0x0',
      metadata: '0x',
      allocationType: AllocationType.simple,
    },
    {
      destination: V.channelId,
      amount: amountForAliceAndBob,
      metadata: encodeGuaranteeData({
        left: convertAddressToBytes32(Ingrid.address),
        right: convertAddressToBytes32(Bob.address),
      }),
      allocationType: AllocationType.guarantee,
    },
  ]
);

// Utils
export async function getFinalizesAtFromTransactionHash(hash: string): Promise<number> {
  // eslint-disable-next-line @typescript-eslint/ban-ts-comment
  // @ts-ignore
  const provider = hre.ethers.provider;
  const receipt = (await provider.getTransactionReceipt(hash)) as ContractReceipt;
  return nitroAdjudicator.interface.decodeEventLog('ChallengeRegistered', receipt.logs[0].data)[2];
}

export async function waitForChallengesToTimeOut(finalizesAtArray: number[]): Promise<void> {
  const finalizesAt = Math.max(...finalizesAtArray);
  // eslint-disable-next-line @typescript-eslint/ban-ts-comment
  // @ts-ignore
  const provider = hre.ethers.provider;
  await provider.send('evm_setNextBlockTimestamp', [finalizesAt + 1]);
  await provider.send('evm_mine', []);
}

/**
 * Constructs a support proof for the supplied channel and calls challenge
 * @returns Challenge transaction, the proof and finalizesAt
 */
export async function challengeChannel(
  channel: TestChannel,
  asset: string
): Promise<{
  challengeTx: ethers.ContractTransaction;
  proof: ReturnType<typeof channel.counterSignedSupportProof>;
  finalizesAt: number;
}> {
  const proof = channel.counterSignedSupportProof(channel.someState(asset)); // TODO use a nontrivial app with a state transition
  const challengeTx = await nitroAdjudicator.challenge(
    proof.fixedPart,
    proof.proof,
    proof.candidate,
    proof.challengeSignature
  );

  const finalizesAt = await getFinalizesAtFromTransactionHash(challengeTx.hash);
  return {challengeTx, proof, finalizesAt};
}

interface ETHBalances {
  Alice: BigNumberish;
  Bob: BigNumberish;
  Ingrid: BigNumberish;
}

interface ETHHoldings {
  LforV: BigNumberish;
  V: BigNumberish;
  X: BigNumberish;
}

/**
 * Asserts the ETH balance of the supplied ethereum account addresses and the ETH holdings in the statechannels asset holding contract for the supplied channelIds.
 */
export async function assertEthBalancesAndHoldings(
  ethBalances: Partial<ETHBalances>,
  ethHoldings: Partial<ETHHoldings>
): Promise<void> {
  const provider = hre.ethers.provider;
  const internalDestinations: {[Property in keyof ETHHoldings]: string} = {
    LforV: LforV.channelId,
    V: V.channelId,
    X: X.channelId,
  };
  const externalDestinations: {[Property in keyof ETHBalances]: string} = {
    Alice: Alice.address,
    Bob: Bob.address,
    Ingrid: Ingrid.address,
  };
  await Promise.all([
    ...Object.keys(ethHoldings).map(async key => {
      expect(
        // eslint-disable-next-line @typescript-eslint/ban-ts-comment
        // @ts-ignore
        (await nitroAdjudicator.holdings(constants.AddressZero, internalDestinations[key])).eq(
          // eslint-disable-next-line @typescript-eslint/ban-ts-comment
          // @ts-ignore
          BigNumber.from(ethHoldings[key])
        )
      ).toBe(true);
    }),
    ...Object.keys(ethBalances).map(async key => {
      expect(
        // eslint-disable-next-line @typescript-eslint/ban-ts-comment
        // @ts-ignore
        (await provider.getBalance(externalDestinations[key])).eq(BigNumber.from(ethBalances[key]))
      ).toBe(true);
    }),
  ]);
}

/**
 * Calculates the gas used by a transaction supplied.
 */
export async function gasUsed(
  txRes: ethers.ContractTransaction // TransactionResponse
): Promise<number> {
  const {gasUsed: gasUsedBN} = await txRes.wait();
  return (gasUsedBN as BigNumber).toNumber();
}

/**
 * Takes a snapshot of the state, execute supplied function and revert the state to the taken snapshot.
 */
export async function executeAndRevert(fnc: () => void) {
  const snapshot = await takeSnapshot();
  await fnc();
  await snapshot.restore();
}

/**
 * Constructs a support proof for the supplied channel which includes a payment voucher, and calls challenge,
 * @returns The state hash, outcome, finalizesAt and gas consumed
 */
export async function challengeVirtualPaymentChannelWithVoucher(
  channel: TestChannel,
  asset: string,
  amount: number,
  payerWallet: Wallet,
  challengerWallet: Wallet
): Promise<{
  stateHash: Bytes32;
  outcome: Outcome;
  finalizesAt: number;
  gasUsed: number;
}> {
  const postFund = channel.someState(asset);
  postFund.appData = '0x';
  postFund.turnNum = 1;

  const proof = [channel.counterSignedSupportProof(postFund).candidate];
  const redemption = channel.someState(asset);
  const voucher: Voucher = {
    channelId: channel.channelId,
    amount: BigNumber.from(amount).toHexString(),
  };
  const voucherSignature = await signVoucher(voucher, payerWallet);
  redemption.appData = encodeVoucherAmountAndSignature(voucher.amount, voucherSignature);
  redemption.turnNum = 2;

  const outcome = channel.outcome(MAGIC_ADDRESS_INDICATING_ETH);
  outcome[0].allocations[0].amount = BigNumber.from(outcome[0].allocations[0].amount)
    .sub(amount)
    .toHexString();
  outcome[0].allocations[1].amount = BigNumber.from(amount).toHexString();
  redemption.outcome = outcome;

  const candidate: SignedVariablePart = {
    variablePart: getVariablePart(redemption),
    sigs: [signState(redemption, challengerWallet.privateKey).signature],
  };

  const challengeSignature = signChallengeMessage(
    [{state: redemption} as SignedState],
    challengerWallet.privateKey
  );

  const challengeTx = await nitroAdjudicator.challenge(
    channel.fixedPart,
    proof,
    candidate,
    challengeSignature
  );

  const finalizesAt = await getFinalizesAtFromTransactionHash(challengeTx.hash);

  const gasUsed = (await challengeTx.wait()).gasUsed.toNumber();

  return {
    stateHash: hashState(redemption),
    finalizesAt,
    outcome: redemption.outcome,
    gasUsed,
  };
}

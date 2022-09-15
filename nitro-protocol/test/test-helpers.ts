import {Contract, ethers, BigNumberish, BigNumber, providers, Wallet} from 'ethers';
import {BytesLike} from '@ethersproject/bytes';
import {Allocation, AllocationType} from '@statechannels/exit-format';
import {isBigNumberish} from '@ethersproject/bignumber/lib/bignumber';
import {LogDescription} from '@ethersproject/abi';

import {ChallengeClearedEvent, ChallengeRegisteredStruct} from '../src/contract/challenge';
import {channelDataToStatus} from '../src/contract/channel-storage';
import {Outcome} from '../src/contract/outcome';
import {Bytes32, VariablePart} from '../src';

// Interfaces

/**
 * A mapping from destination to BigNumberish. E.g. {ALICE:2, BOB:3}. Only used in testing.
 */
export interface AssetOutcomeShortHand {
  [destination: string]: BigNumberish;
}

/**
 * A mapping from asset to AssetOutcomeShorthand. E.g. {ETH: {ALICE:2, BOB:3}, DAI: {ALICE:1, BOB:4}}. Only used in testing.
 */
export interface OutcomeShortHand {
  [assetHolder: string]: AssetOutcomeShortHand;
}

export interface AddressesLookup {
  [shorthand: string]: string | undefined;
}

// Functions
export const getTestProvider = (): ethers.providers.JsonRpcProvider => {
  if (!process.env.GANACHE_PORT) {
    throw new Error('Missing environment variable GANACHE_PORT required');
  }
  return new ethers.providers.JsonRpcProvider(`http://localhost:${process.env.GANACHE_PORT}`);
};

/**
 * Get a rich object representing an on-chain contract
 * @param provider an ethers JsonRpcProvider
 * @param artifact an object containing the abi of the contract in question
 * @param address the ethereum address of the contract, once it is deployed
 * @returns a rich (ethers) Contract object with a connected signer (ther 0th signer of the supplied provider)
 */
export function setupContract(
  provider: ethers.providers.JsonRpcProvider,
  artifact: {abi: ethers.ContractInterface},
  address: string
): Contract {
  return new ethers.Contract(address, artifact.abi, provider.getSigner(0));
}

export function getCountingAppContractAddress(): string {
  return process.env.COUNTING_APP_ADDRESS;
}

export const nonParticipant = ethers.Wallet.createRandom();

/**
 * Generate `n` wallets and return them alongside with array of their addresses.
 * @param n Number of wallets to create.
 * @returns N wallets with their addresses.
 */
export function generateParticipants(n: number): {
  wallets: Wallet[];
  participants: string[];
} {
  const participants = new Array(n).fill('');
  const wallets = new Array(n);

  for (let i = 0; i < n; i++) {
    wallets[i] = Wallet.createRandom();
    participants[i] = wallets[i].address;
  }

  return {participants, wallets};
}

export const clearedChallengeFingerprint = (turnNumRecord = 5): Bytes32 =>
  channelDataToStatus({
    turnNumRecord,
    finalizesAt: 0,
  });

export const ongoingChallengeFingerprint = (turnNumRecord = 5): Bytes32 =>
  channelDataToStatus({
    turnNumRecord,
    finalizesAt: 1e12,
    outcome: [],
  });

export const finalizedFingerprint = (
  turnNumRecord = 5,
  finalizesAt = 1,
  outcome: Outcome = [],
  state = undefined
): Bytes32 =>
  channelDataToStatus({
    turnNumRecord,
    finalizesAt,
    outcome,
    state,
  });

export const parseOutcomeEventResult = (eventOutcomeResult: any[]): Outcome => {
  eventOutcomeResult = Array.from(eventOutcomeResult);
  const outcome: Outcome = [];

  if (eventOutcomeResult.length == 0) {
    return outcome;
  }

  eventOutcomeResult.forEach((eventSingleAssetExit: any[]) => {
    const asset: string = eventSingleAssetExit[0];
    const metadata: BytesLike = eventSingleAssetExit[1];
    const eventAllocations: any[] = Array.from(eventSingleAssetExit[2]);
    const allocations: Allocation[] = [];

    if (eventAllocations.length != 0) {
      eventAllocations.forEach((eventAllocation: any[]) => {
        const destination: string = eventAllocation[0];
        const amount: string = BigNumber.from(eventAllocation[1]['_hex']).toString();
        const allocationType: number = eventAllocation[2];
        const metadata: string = eventAllocation[3];

        allocations.push({destination, amount, allocationType, metadata});
      });
    }

    outcome.push({asset, metadata, allocations});
  });

  return outcome;
};

export const parseVariablePartEventResult = (vpEventResult: any[]): VariablePart => {
  vpEventResult = Array.from(vpEventResult);
  return {
    outcome: parseOutcomeEventResult(vpEventResult[0]),
    appData: vpEventResult[1],
    turnNum: vpEventResult[2],
    isFinal: vpEventResult[3],
  };
};

export const newChallengeRegisteredEvent = (
  contract: ethers.Contract,
  channelId: string
): Promise<ChallengeRegisteredStruct[keyof ChallengeRegisteredStruct]> => {
  const filter = contract.filters.ChallengeRegistered(channelId);
  return new Promise(resolve => {
    contract.on(
      filter,
      (
        eventChannelIdArg,
        eventTurnNumRecordArg,
        eventFinalizesAtArg,
        eventChallengerArg,
        eventIsFinalArg,
        eventFixedPartArg,
        eventChallengeVariablePartArg
      ) => {
        contract.removeAllListeners(filter);
        resolve([
          eventChannelIdArg,
          eventTurnNumRecordArg,
          eventFinalizesAtArg,
          eventChallengerArg,
          eventIsFinalArg,
          eventFixedPartArg,
          eventChallengeVariablePartArg,
        ]);
      }
    );
  });
};

export const newChallengeClearedEvent = (
  contract: ethers.Contract,
  channelId: string
): Promise<ChallengeClearedEvent[keyof ChallengeClearedEvent]> => {
  const filter = contract.filters.ChallengeCleared(channelId);
  return new Promise(resolve => {
    contract.on(filter, (eventChannelId, eventTurnNumRecord) => {
      // Match event for this channel only
      contract.removeAllListeners(filter);
      resolve([eventChannelId, eventTurnNumRecord]);
    });
  });
};

export const newConcludedEvent = (
  contract: ethers.Contract,
  channelId: string
): Promise<[Bytes32]> => {
  const filter = contract.filters.Concluded(channelId);
  return new Promise(resolve => {
    contract.on(filter, () => {
      // Match event for this channel only
      contract.removeAllListeners(filter);
      resolve([channelId]);
    });
  });
};

export const newDepositedEvent = (
  contract: ethers.Contract,
  destination: string
): Promise<[string, BigNumber, BigNumber]> => {
  const filter = contract.filters.Deposited(destination);
  return new Promise(resolve => {
    contract.on(
      filter,
      (eventDestination: string, amountDeposited: BigNumber, amountHeld: BigNumber) => {
        // Match event for this destination only
        contract.removeAllListeners(filter);
        resolve([eventDestination, amountDeposited, amountHeld]);
      }
    );
  });
};

// Copied from https://stackoverflow.com/questions/58325771/how-to-generate-random-hex-string-in-javascript
const genRanHex = (size: number) =>
  [...Array(size)].map(() => Math.floor(Math.random() * 16).toString(16)).join('');

export const randomChannelId = (): Bytes32 => '0x' + genRanHex(64);
export const randomExternalDestination = (): Bytes32 => '0x' + genRanHex(40).padStart(64, '0');

export async function sendTransaction(
  provider: ethers.providers.JsonRpcProvider,
  contractAddress: string,
  transaction: providers.TransactionRequest
): Promise<providers.TransactionReceipt> {
  const signer = provider.getSigner();
  const response = await signer.sendTransaction({to: contractAddress, ...transaction});
  return await response.wait();
}

// Recursively replaces any key with the value of that key in the addresses object
// BigNumberify all numbers
/**
 * Recursively replaces any key in a copy of the supplied object with the value of that key in the supplied addresses object. Also BigNumberifies all numbers.
 * Used in testing only.
 * @param object Object to be copied and modified
 * @param addresses Key-value address lookup
 * @returns suitably modified copy of object
 */
export function replaceAddressesAndBigNumberify(
  object: AssetOutcomeShortHand | OutcomeShortHand | BigNumberish,
  addresses: AddressesLookup
): AssetOutcomeShortHand | OutcomeShortHand | BigNumberish {
  if (isBigNumberish(object)) {
    return BigNumber.from(object);
  }
  const newObject: AssetOutcomeShortHand | OutcomeShortHand = {};
  Object.keys(object).forEach(key => {
    if (isBigNumberish(object[key])) {
      newObject[addresses[key] as string] = BigNumber.from(object[key]);
    } else if (typeof object[key] === 'object') {
      // Recurse
      newObject[addresses[key] as string] = replaceAddressesAndBigNumberify(
        object[key],
        addresses
      ) as AssetOutcomeShortHand | BigNumberish;
    }
  });
  return newObject;
}

// Sets the holdings defined in the multipleHoldings object. Requires an array of the relevant contracts to be passed in.
export function resetMultipleHoldings(
  multipleHoldings: OutcomeShortHand,
  contractsArray: Contract[]
): void {
  Object.keys(multipleHoldings).forEach(assetHolder => {
    const holdings = multipleHoldings[assetHolder];
    Object.keys(holdings).forEach(async destination => {
      const amount = holdings[destination];
      contractsArray.forEach(async contract => {
        if (contract.address === assetHolder) {
          await (await contract.setHoldings(destination, amount)).wait();
          expect((await contract.holdings(destination)).eq(amount)).toBe(true);
        }
      });
    });
  });
}

// Check the holdings defined in the multipleHoldings object. Requires an array of the relevant contracts to be passed in.
export function checkMultipleHoldings(
  multipleHoldings: OutcomeShortHand,
  contractsArray: Contract[]
): void {
  Object.keys(multipleHoldings).forEach(assetHolder => {
    const holdings = multipleHoldings[assetHolder];
    Object.keys(holdings).forEach(async destination => {
      const amount = holdings[destination];
      contractsArray.forEach(async contract => {
        if (contract.address === assetHolder) {
          expect((await contract.holdings(destination)).eq(amount)).toBe(true);
        }
      });
    });
  });
}

/** Computes an Outcome from a shorthand description */
export function computeOutcome(outcomeShortHand: OutcomeShortHand): Outcome {
  const outcome: Outcome = [];
  Object.keys(outcomeShortHand).forEach(asset => {
    const allocations: Allocation[] = [];
    Object.keys(outcomeShortHand[asset]).forEach(destination =>
      allocations.push({
        destination,
        amount: BigNumber.from(outcomeShortHand[asset][destination]).toHexString(),
        metadata: '0x',
        allocationType: AllocationType.simple,
      })
    );
    outcome.push({asset, metadata: '0x', allocations});
  });
  return outcome;
}

interface Event extends LogDescription {
  contract: string;
}

export function compileEventsFromLogs(logs: any[], contractsArray: Contract[]): Event[] {
  const events: Event[] = [];
  logs.forEach(log => {
    contractsArray.forEach(contract => {
      if (log.address === contract.address) {
        events.push({...contract.interface.parseLog(log), contract: log.address});
      }
    });
  });
  return events;
}

export function getRandomNonce(seed: string): string {
  return ethers.utils.id(seed);
}

export const largeOutcome = (
  numAllocationItems: number,
  asset: string = ethers.Wallet.createRandom().address
): Outcome => {
  const randomDestination = '0x8595a84df2d81430f6213ece3d8519c77daf98f04fe54e253a2caeef4d2add39';
  return numAllocationItems > 0
    ? [
        {
          allocations: Array(numAllocationItems).fill({
            destination: randomDestination,
            amount: '0x01',
            allocationType: AllocationType.simple,
            metadata: '0x',
          }),
          asset,
          metadata: '0x',
        },
      ]
    : [];
};

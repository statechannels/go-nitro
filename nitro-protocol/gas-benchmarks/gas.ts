export type GasResults = Record<
  Path,
  {
    satp: any;
  }
>;

type Path =
  | 'deployInfrastructureContracts'
  | 'directlyFundAChannelWithETHFirst'
  | 'directlyFundAChannelWithETHSecond'
  | 'directlyFundAChannelWithERC20First'
  | 'directlyFundAChannelWithERC20Second'
  | 'ETHexit'
  | 'ERC20exit'
  | 'ETHexitSad'
  | 'ETHexitSadLedgerFunded'
  | 'ETHexitSadVirtualFunded'
  | 'ETHexitSadLedgerFunded';

// The channel being benchmarked is a 2 party null app funded with 5 wei / tokens each.
// KEY
// ---
// ⬛ -> funding on chain (from Alice)
//  C    channel not yet on chain
// (C)   channel finalized on chain
// 👩    Alice's external destination (e.g. her EOA)
export const emptyGasResults: GasResults = {
  deployInfrastructureContracts: {
    satp: {
      NitroAdjudicator: 0, // Singleton
    },
  },
  directlyFundAChannelWithETHFirst: {
    satp: 0,
  },
  directlyFundAChannelWithETHSecond: {
    // meaning the second participant in the channel
    satp: 0,
  },
  directlyFundAChannelWithERC20First: {
    // The depositor begins with zero tokens approved for the AssetHolder
    // The AssetHolder begins with some token balance already
    // The depositor retains a nonzero balance of tokens after depositing
    // The depositor retains some tokens approved for the AssetHolder after depositing
    satp: {
      approve: 0,
      // ^^^^^
      // In principle this only needs to be done once per account
      // (the cost may be amortized over several deposits into this AssetHolder)
      deposit: 0,
    },
  },
  directlyFundAChannelWithERC20Second: {
    // meaning the second participant in the channel
    satp: {
      approve: 0,
      // ^^^^^
      // In principle this only needs to be done once per account
      // (the cost may be amortized over several deposits into this AssetHolder)
      deposit: 0,
    },
  },
  ETHexit: {
    // We completely liquidate the channel (paying out both parties)
    satp: 0,
  },
  ERC20exit: {
    // We completely liquidate the channel (paying out both parties)
    satp: 0,
  },
  ETHexitSad: {
    // Scenario: Counterparty Bob goes offline
    // initially                 ⬛ ->  X  -> 👩
    // challenge + timeout       ⬛ -> (X) -> 👩
    // transferAllAssets         ⬛ --------> 👩
    satp: {
      challenge: 0,
      transferAllAssets: 0,
      total: 0,
    },
  },
  ETHexitSadLedgerFunded: {
    // Scenario: Counterparty Bob goes offline
    // initially                   ⬛ ->  L  ->  X  -> 👩
    // challenge X, L and timeout  ⬛ -> (L) -> (X) -> 👩
    // transferAllAssetsL          ⬛ --------> (X) -> 👩
    // transferAllAssetsX          ⬛ ---------------> 👩
    satp: {
      challengeX: 0,
      challengeL: 0,
      transferAllAssetsL: 0,
      transferAllAssetsX: 0,
      total: 0,
    },
  },
  ETHexitSadVirtualFunded: {
    // Scenario: Intermediary Ingrid goes offline
    // initially                   ⬛ ->  L  ->  V  -> 👩
    // challenge L,V   + timeout   ⬛ -> (L) -> (V) -> 👩
    // reclaim L                   ⬛ -- (L) --------> 👩
    // transferAllAssetsL          ⬛ ---------------> 👩
    // TODO
    satp: {},
  },
};

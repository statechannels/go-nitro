type GasRequiredTo = Record<
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
export const gasRequiredTo: GasRequiredTo = {
  deployInfrastructureContracts: {
    satp: {
      NitroAdjudicator: 4_389_105, // Singleton
    },
  },
  directlyFundAChannelWithETHFirst: {
    satp: 47_797,
  },
  directlyFundAChannelWithETHSecond: {
    // meaning the second participant in the channel
    satp: 30_709,
  },
  directlyFundAChannelWithERC20First: {
    // The depositor begins with zero tokens approved for the AssetHolder
    // The AssetHolder begins with some token balance already
    // The depositor retains a nonzero balance of tokens after depositing
    // The depositor retains some tokens approved for the AssetHolder after depositing
    satp: {
      approve: 46_383,
      // ^^^^^
      // In principle this only needs to be done once per account
      // (the cost may be amortized over several deposits into this AssetHolder)
      deposit: 71_280,
    },
  },
  directlyFundAChannelWithERC20Second: {
    // meaning the second participant in the channel
    satp: {
      approve: 46_383,
      // ^^^^^
      // In principle this only needs to be done once per account
      // (the cost may be amortized over several deposits into this AssetHolder)
      deposit: 54_192,
    },
  },
  ETHexit: {
    // We completely liquidate the channel (paying out both parties)
    satp: 144_462,
  },
  ERC20exit: {
    // We completely liquidate the channel (paying out both parties)
    satp: 134_816,
  },
  ETHexitSad: {
    // Scenario: Counterparty Bob goes offline
    // initially                 ⬛ ->  X  -> 👩
    // challenge + timeout       ⬛ -> (X) -> 👩
    // transferAllAssets         ⬛ --------> 👩
    satp: {
      challenge: 107_405,
      transferAllAssets: 110_149,
      total: 217_554,
    },
  },
  ETHexitSadLedgerFunded: {
    // Scenario: Counterparty Bob goes offline
    satp: {
      // initially                   ⬛ ->  L  ->  X  -> 👩
      // challenge X, L and timeout  ⬛ -> (L) -> (X) -> 👩
      // transferAllAssetsL          ⬛ --------> (X) -> 👩
      // transferAllAssetsX          ⬛ ---------------> 👩
      challengeX: 107_405,
      challengeL: 101_397,
      transferAllAssetsL: 59_047,
      transferAllAssetsX: 110_149,
      total: 377_998,
    },
  },
  ETHexitSadVirtualFunded: {
    // Scenario: Intermediary Ingrid goes offline
    satp: {
      // initially                   ⬛ ->  L  ->  V  -> 👩
      // challenge L,V   + timeout   ⬛ -> (L) -> (V) -> 👩
      // reclaim L                   ⬛ -- (L) --------> 👩
      // transferAllAssetsL          ⬛ ---------------> 👩
      // TODO
    },
  },
};

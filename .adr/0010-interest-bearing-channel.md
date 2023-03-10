# 0010 - Interest Bearing Channel Application

# Status

Accepted

# Context

For a permissionless state channel network (SCN) to be sustainable, incentive structures must exist for each class of system actor (service providers, service consumers, & liquidity providers).

Our work generally presumes the existence of incentives that motivate service providers & consumers. The design work for us is in mechanisms to incentivize intermediaries (liquidity providers) into the network.

# Prior Art

The Bitcoin Lightning Network (LN) is an in-production payment channel network (PCN) with a fee mechanism.

## Similarities to Nitro-SCN

In a PCN, payments between A and B are routed through intermediaries I_1, ..., I_n, where balances are pairwise adjusted along the route. Capital in intermediary channels is locked for the duration of the payment - ie, until either the payment completes or times out.

In an SCN, collateral to secure channels between A and B is locked along a similar route of intermediaries, and closing balance adjusments between A and B are pairwise adjusted along the route. Capital in intermediary channels is locked for the duration

## Differences from Nitro-SCN

### Ecosystem / Context

The Bitcoin ecosystem has a strong intrinsic incentive toward participating in LN, as LN is the most (only?) viable solution to scaling BTC denominated payments with security anchored to bitcoin.

Further: holders of BTC have relatively few options for activating their capital to produce a yield. Assets in EVM environments can generate yield via any number of defi protocols, including Ethereum proof-of-stake (staking pools offer approx %4 yield with near zero overhead).

### Technical

In the dominant use case, PCN payments are brief. The capital lockup time for an individual payment is `O(n)`, where `n` is length of the chain of intermediaries. LN strives toward a mesh network topology so that

- scale requirements can be met (a single hub would struggle to be involved in _all_ payments)
- privacy claims are defensible (payments are onion routed - each node passes the payment one step further but cannot read the entire route. Would be moot with single hub, and less effective with dense network of hubs)

By contrast, the dominant case in an SCN is for individual capital lockup events (channel fundings) to have non-trivial duration.

For an SCN, the capital lockup time is equal to `fundingTime + channelRuningTime + defundingTime`. Each of `fundingTime` and `defundingTime` is expected to be smaller than LN's `O(n)`, because Nitro makes weaker privacy claims & is happy to disclose the entire intermediary chain at funding time. Despite this, the _total lockup time_ in Nitro will be much bigger on average because of the run-time of funded channels.

## LN Fee mechanism + Commentary

In LN, intermediaries charge a base fee per transaction, plus a scalar fee based on the size of the transaction. In practice, this translates to the base fee acting as a "networking fee" and the scalar acting as a "liquidity rental fee".

Fees on liquidity are required because:

- a given intermediary channel has hard limits on the number of pending transactions it can process at a time (~400)
- a given intermediary channel has only so much capacity to adjust its balance in one direction or the other (IE, to forward payments)

A channel in a PCN is capacity limited according to its current balance. The limitation is asymmetric: a one BTC channel whose balance is 0.5-0.5 can forward payments up to 0.5 BTC in either direction, but the same channel at a different time with a balance of 0.1-0.9 can forward only 0.1 to the right (and 0.9 to the left). Most LN channels _expect_ to be bi-directional in practice, and for the general flow of payments back and forth to keep the channel balanced on average.

There is something of a mismatch between LN's fee model & the state of its network at any given moment.

# Decision Criteria

There are a number of decision dimensions here:

1. purpose-fit of introduced economic model (does it drive pro-network behaviours?)
2. security implications
3. efficiency / gas consumption implications
4. programmability / tunability of implementation
5. ergonomics for
   - channel application developers
   - SCN wallet developers
   - network participants
6. software lock-in, protocol flexibility, composability

... plus the usual meta-considerations such as engineering effort.

# Considered Options

## Fee Buckets as a Protocol Primitive

(See related [Notion Doc](https://www.notion.so/statechannels/RFC-20-Introduce-a-fee-model-for-intermediaries-by-modifying-reclaim-dce28ac74d764fcf93db7d6b5cf05b3a))

The linked document describes a mechanism where a channel _payer_ (Alice) fronts a specified (potentiall run-time negotiated) fee for each funded virtual channel, which is distributed among route intermediaries at channel closing time. This is enforced at the protocol level via a modification to `reclaim`.

Against the above decision criteria: (on 😀 🙂 😐 🙁 😢
scale)

1. Fit: 😐 Per channel fees which aren't time-based recreate some of the existing issues from LN, which, as discussed, are still more important in the SCN context.
2. Security: 🙁 Modification to `reclaim` affects security of the protocol at large - more code in the public API of the adjudicator is greater surface area for bugs, exploits, etc.
3. Efficiency: 🙂 Probably reasonbly strong here - the changes to `reclaim` are only relevent during the sad path, and in this case the accounting is pretty straightforward - checking that outcomes match the prefund commitment for fee distribution.
4. Flexibility: 🙂 It's relatively easy in this model to adjust fees on the fly and also to adjust based on the number of hops involved in the virtual fund.
5. Ergonomics:
   - AppDev: 😀 no effect on channel application development
   - Wallet: 😐 adds some complexity to the virtualfund / vurtual defund operations
   - Users: 🙁 Networking fees originating at the service consumer-side can potentially introduce extra UI hurdles (eg, think of setting gas fees for an L1 transaction). This can be obfuscated from the user by building network fees into the total service fee, but this creates potential for information assymetries in the marketplace which aren't to the user's benefit.
6. software-lock: 🙁 Modifications to the adjudicator are the hardest lock-in possible for the network. If this implementation is revealed to be broken, or unfit for purpose, then migration to a separate in-protocol implementation would require a redeployment of the adjudicator, hard shut-down of the entire running SCN, etc. Software-lock could be improved by, eg, implementing claim via upgradable contract, but this introduces tradeoffs in the credible trustlessness of the network.

Meta commentary: sketches toward implementation of this model showed that the engineering, security and auditability effort would be reasonably high for this option.

## Interest Bearing Application Channels

Against the above decision criteria: (on 😀 🙂 😐 🙁 😢
scale)

1. Fit: 🙂 Interest yields on committed collateral strongly mimic the "actual cost centers" of the network. A potential limitation is that hub-hub channels are not financialized, but this is at worse a "mixed curse". Intermediaries have a natural incentive to maintain their ledger connections (it connects their clients to other clients). Intermediaries are, in expectation, the most financially savvy class of participants, so they ought to be well situated to price these uncompensated channels into their offerings to service providers.
2. Security: 😀 Adding application channels has no (obvious) effect on existing security guarantees of the protocol. All moving / connected parts of the network interact in the same way as before. Introduced threats are confined to the interest-bearing ledger channels themselves.
3. Efficiency: 😐 Similar to above, the added code is only invoked in the sad case so should be a non-issue most of the time. The accounting in this structure is a little more complicated & gas intensive than in the fee bucket scenario, but also presents itself in fewer channels (only ledger channels).
4. Flexibility: 🙂
5. Ergonomics:
   - AppDev: 😀 no effect on channel application development
   - Wallet: 😐 adds some complexity to the direct fund / defund operations, and to the accounting inside a running ledger channel (including a requirement to track current block number)
   - Users: 🙂 "extra" work to accomodate fees under this structure is done only by the relatively professionalized users: service providers and liquidity providers.
6. software-lock: 😀 users commit to this fee scheme only over the duration of individual ledger channels. If the scheme is broken or unfit, then only those channels need to be wound down. The application-channel architecture for fees also implicitly invites future refinement, supplimentation, & innovation.

Meta commentary: sketches toward implementation of this model were reasonably straight forward, and security analysis likewise seems straightforward.

## Decision

We have implemented an interest-bearing application channel as first-pass on incentivizing liquiidty providers.

## Conclusion

- fee models are subject to future modification (open market)
- nice to gain experience building against our own API

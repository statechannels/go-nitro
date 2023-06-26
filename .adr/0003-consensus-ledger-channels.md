# 0003 -- Consensus ledger channel updates

## Status

Accepted

## Definitions

- **L(x,y)** is the ledger channel between x and y.
- references to **virtualfunding** operations are intended as generic references to either `virtualfund` or `virtualdefund`. (Similar for **directfunding**)

## Context

A running ledger channel must serve concurrent funding requests originating from many different particpants. In the below channel network:

```
Alice ----.              .-- Doris
           \            /
Bob   ---- Irene --- Ivy --- Eric
           /           \
Charlie --/             \--- Frannie
```

Any client participant (A-F) can initate a virtualfund protocol with any other participant at any time, which will require updates to each ledger channel on the path between them. Updates to ledger channels must be ordered in some capacity so that (for example) the channel is not "overdrawn".

The hub-hub channel L(Irene, Ivy) is particularly challenging, since each of Irene & Ivy can recieve a request that will require updates to L(Irene, Ivy) at any time, from any participant.

## Considered Options

### Complete Consensus Ledger Update

Consensus channel runs the null-app. State updates must be by unanimous consensus (ie, signed by all participants).

pros:

- security guarantees as strong as the base protocol
- no bespoke on-chain application logic
- lower off-chain implementation complexity

cons:

- sacrificing some performance optimizations

### Async Ledger Updates

- (https://www.notion.so/statechannels/RFC-13-Async-virtual-funding-b2b6ed9e39b34a7fbd362026dc248b0f)

Above document outlines a protocol for ledger channel in which the channel application logic:

- contains independent queues for updates:
  - signed by participantA
  - signed by participantB
  - cosigned
- contains update rules which allow participants to unilaterally include items from their counterparty's queue into the current outcome

pros:

- potential for a reduction in total network round-trips on virtualfund / defund operations
- potential for reduced "blocking" time while waiting on network round-trips

cons:

- greater off-chain implementation complexity
- introduces on-chain implementation requirements
- more difficult security analysis

## Decision

go-nitro updates its ledger channel via **complete consensus**.

Update **ordering** is managed by a designated ledger channel leader (`participants[0]`) and updates are countersigned by the follower (`participants[1]`).

## Implementation Detail

### Introduced Data Structures

The consensus update ledger channel implementation exists as the `struct ConsensusChannel` export from `package consensus_channel`.

`ConsensusChannel` represents an API into a running channel that is specific to ledger channels. In addition to default channel data like `state.FixedPart` (package state), it:

- defines the data structure `LedgerOutcome`, tailored to represent the variablepart of a ledger channel (each party's balance + a map of guarantees for virtualchannels)
- defines `Proposal` update structures `Add` and `Remove`, which are used to update `LedgerOutcome` when funding or defunding a given channel
- contains a queue of `SignedProposal`s, which is used to order updates
- contains various helpers to translate between canonical channel `state.State` data and the local `LedgerOutcome` representation

Further, the package exposes constructors `NewLeaderChannel()` and `NewFollowerChannel()` which return role-specific APIs and offer assurances against nodes mistakenly assuming incorrect roles.

```
LeaderChannel:
  Propose(proposal) -> process the funding or defunding of a channel, enqueue it, sign it (prepare for sending)
  leaderReceive(signedProposal) -> inspect a proposal that has been countersigned and returned by the counterparty, dequeue it & apply to current state

FollowerChannel:
  followerReceive(signedProposal) -> process a proposal received from the channel leader & enqueue it
  SignNextProposal(expectedProposal) -> inspect the proposalQueue, and if item[0] matches the supplied expectation, sign, dequeue, and prepare for sending to the counterparty
```

**Note**: state channel security guarantees depend on signed channel states. To that end, proposals sent over the wire are labelled with their turn number, and the signature sent with a proposal is a signature on the **resultant channel state** after the proposed `Add` or `Remove` is applied - not on the proposal data itself.

### Integration with virtualfunding protocols

A client `X` initializing a new virtual channel via the `virtualfund` package (`./protocols/virtualfund`) is provided access to its ledger channel via a getter.

The virtualfund objective is responsible for crafting the `Add` proposal for its target channel, and

- in the case of a leader ledger channel, calling `Propose()` with this proposal to generate a message for the counterparty
- in the case of a follower ledger channel, applying this constructed `Add` as the `expectedProposal` for `SignNextProposal()`

Each of these produces a SignedProposal, which is returned to the engine as a protocol side effect and sent to the ledger counterpaty.

For `virtualdefund`, the above applies with a `Remove` proposal in place of the `Add`.

### Integration with directfunding protocols

Existing `directfund` and `directdefund` protocols operate on the `Channel` struct from `package channel`. The `Channel` struct is the canonical view of a channel.

This requires a conversion steps on each end:

- `directfund.Objective` contains converter method `CreateConsensusChannel()`, which prepares a ledger channel for usage by the virtualfunding protocols after the directfund is complete
- package `directdefund` contains converter method `CreateChannelFromConsensusChannel(cc)`, which prepares a `Channel` struct from the ledger channel.

This is not ideal (why create one temp data-type only to immediatly convert it for storage?), and is subject to future refactoring. See [this notion doc](https://www.notion.so/statechannels/Channels-ConsensusChannels-Objectives-9701e24e2bc1491a83fae375e4e4c64a#9c41424d07694150bdc6cf373309757a) for further detail.

### Integration with engine

Because individual protocols are implemented as pure functions, reordering recieved proposals and managing sequential updates is handled by the engine (`package engine`). See [ADR-0006](0006-proposal-processing-ledger-effects.md) for details on how proposal ordering is handled by the engine.

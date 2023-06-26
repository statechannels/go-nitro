# 0004 -- Proposal Messaging

## Status

Accepted with misgivings (future considerations)

## Context

Consensus ledger channels update via a designated leader crafting update proposals, which are responded to by the channel follower.

Because each proposal implies a specific, sequentially numbered state, the sending and receiving of proposals must be reliable in context of a real-world network. Safeguards must exist against both **dropped messages** and **message reordering**.

See [ADR-0003](0003-consensus-ledger-channels.md) for additional relevant context.

## Considered Options

### Redundancy

Nodes include all unacknowledged proposals with each message. In the context of a healthy network, this might be expected to result in occasional double transmissions and rare triple transmissions.

In an unhealthy network, or in the case of an offline counterparty, nodes will require a mechanism to "give up" on sending messages to the unavailable counterparty for some cooldown period.

### MessageService with guaranteed delivery

A messageservice which recieves an acknowledgement on each transmission could effectively guarantee message delivery, and could be relied upon for correct message ordering as well by waiting for message `N` to be acknowledged before sending message `N+1`.

This approach simplifies the **message reordering** problem for both sender and receiver, but adds latency (network trips for each `ack` signal) to every funding operation. Ledger channel updates are the expected bottleneck in a busy channel network, so adding latency here is not preferred.

#### Receiever requests absent proposals

A follower node waiting on proposal `N` who recieves some out-of-order proposal `N+M` could explicitly request proposals `[N, ... ,N+M-1]` from the counterparty. In the case of out-of-order messages, this is well-behaved:

- it allows the channel leader to optimistically send all proposals individually
- where message `N` and `N+1` cross mid-flight (ie, `N+1` arrives first), the follower makes a retransmission request but can still proceed immediately when message `N` arrives

In the case of dropped messages, this scheme is expected to be less performant than the `Redundancy` scheme, as it requires the extra network trip with the request for a retransmission, which itself can not be triggered until the follower is "aware" that it is missing proposals.

## Decision

go-nitro currently implements redundant proposal broadcasts.

The current implementation comes with the following caveats:

- assumptions about the frequency of doubled / tripled (or worse) proposal transmissions are untested
- there are no mechanisms in place for applying cooldown periods for unreachable counterparties
- proposals are batched naively (sending `N` proposals results in a payload `N` times as large as a single proposal). Future work could improve on this by omitting redundant information (repeated labelling of the ledger channel, signatures associated with each individual proposal, etc).

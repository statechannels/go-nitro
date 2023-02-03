---
description: Wrapping up the channel execution.
---

# Finalizing a channel

Finalization of a state channel is a necessary step before defunding it. It can happen on- or off-chain.

## Happy path

In the so-called 'happy' case, all participants cooperate to achieve this.

A participant wishing to end the state channel will sign a state with `isFinal = true`. Then, the other participants may support that state. Once a full set of `n` such signatures exists (this set is known as a **finalization proof**) the channel is said to be 'closed' or 'finalized'.

### Off chain

In most cases, the channel would be finalized and defunded [off chain](./0060-funding-a-channel.md#fund-from-an-existing-channel), and no contract calls are necessary.

### On chain -- calling `conclude`

In the case where assets were deposited against the channel on chain (the case of direct funding), anyone in possession of a finalization proof may use it to finalize the `outcome` on-chain. They would do this by calling `conclude` on the adjudicator. This enables [assets to be released](./release-assets).

The conclude method allows anyone with sufficient off-chain state to immediately finalize an outcome for a channel without having to wait for a challenge to expire (more on that later).

The off-chain state(s) is submitted (in an optimized format), and once relevant checks have passed, an expired challenge is stored against the `channelId`. (This is an implementation detail -- the important point is that the chain shows that the channel has been finalized.)

```typescript
TODO example
```

## Sad path

When cooperation breaks down, it is possible to finalize a state channel without requiring on-demand cooperation of all `n` participants. This is the so-called 'sad' path to finalizing a channel, and it requires a supported (but not necessarily `isFinal`) state(s) being submitted to the chain.

The `challenge` function allows anyone holding the appropriate off-chain state(s) to register a _challenge state_ on chain. It is designed to ensure that a state channel can progress or be finalized in the event of inactivity on behalf of a participant (e.g. the current mover).

The required data for this method consists of a single state, along with `n` signatures. Once these are submitted (in an optimized format), and once relevant checks have passed, an `outcome` is registered against the `channelId`, with a finalization time set at some delay after the transaction is processed.

This delay allows the challenge to be cleared by a timely and well-formed [respond](./clear-a-challenge#call-respond) or [checkpoint](./clear-a-challenge#call-checkpoint) transaction. We'll get to those shortly. If no such transaction is forthcoming, the challenge will time out, allowing the `outcome` registered to be finalized. A finalized outcome can then be used to extract funds from the channel (more on that below, too).

```typescript
TODO example
```

!!! tip

    The `challengeDuration` is a [fixed parameter](./execute-state-transitions#construct-a-state-with-the-correct-format) expressed in seconds, that is set when a channel is proposed. It should be set to a value low enough that participants may close a channel in a reasonable amount of time in the case that a counterparty becomes unresponsive; but high enough that malicious challenges can be detected and responded to without placing unreasonable liveness requirements on responders. A `challengeDuration` of 1 day is a reasonable starting point, but the "right" value will likely depend on your application.

### Call `challenge`

```typescript
TODO example
```

!!!note

    The challenger needs to sign this data:

    ```
    keccak256(abi.encode(challengeStateHash, 'forceMove'))
    ```

    in order to form `challengerSig`. This signals their intent to challenge this channel with this particular state. This mechanism allows the challenge to be authorized only by a channel participant.


    We provide a handy utility function `signChallengeMessage` to form this signature.

A challenge being registered does _not_ mean that the channel will inexorably finalize. Participants have the timeout period in order to be able to respond. Perhaps they come back online after a brief spell of inactivity, or perhaps the challenger was trying to (maliciously) finalize the channel with a supported but outdated (or 'stale') state.

## Clear a challenge

### Call `checkpoint`

The `checkpoint` method allows anyone with a supported off-chain state to establish a new and higher `turnNumRecord` on chain, and leave the resulting channel in the "Open" mode. It can be used to clear a challenge.

### Call `challenge` again

It is important to understand that a challenge may be "cleared" by another more recent challenge. The channel will be left in challenge mode (so it has not really been 'cleared' in that sense), but some [on chain storage](./understand-adjudicator-status) will be updated, such as the deadline for responding.

## Extract info from Adjudicator Events

You may have noticed that to respond, the challenge state itself must be (re)submitted to the chain. To save gas, information is only stored on chain in a hashed format. Clients should, therefore, cache information emitted in Events emitted by the adjudicator, in order to be able to respond to challenges.

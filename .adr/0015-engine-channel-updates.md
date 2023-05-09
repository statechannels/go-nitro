# 0015 Engine Channel Updates

## Status

Review

## Context

We want to enhance the our go-nitro RPC client and server with the ability to receive notifications whenever a channel is updated. This allows us to design a UI that can easily update when a channel changes.

To support this we need some kind of mechanism for the engine to:

- determine when a channel has been updated
- return that information to the API client

## Decision

The engine now returns two new slices on an `EngineEvent`, `LedgerChannelUpdates` and `PaymentChannelUpdates`. These slices contain `query.PaymentChannelInfo` and `query.LedgerChannelInfo` for any updated channels.

After cranking the engine now iterates over the objective's [Related](../protocols/interfaces.go#L87), creating a `query.PaymentChannelInfo` or `query.LedgerChannelInfo` for each related channel. These are sent to the client via a engine event's `LedgerChannelUpdates`/`PaymentChannelUpdates` to notify the API client that updates have occurred. We return a channel info for every one of an objective's `Related`, even if they have not been changed. ** This may result in the engine returning duplicate updates to the API client; we rely on the API client to determine if an update has changed and should be dispatched as an event. **

## Alternatives considered

### Protocols returning Updated Channels

Instead of assuming that all channels have been updated when an objective is cranked, protocols could report which channels they have updated. `Crank` would return a collection of channels that were updated. This was the initial approach I took.

Unfortunately this requires adding a lot of boilerplate code throughout the protocols, making the protocols less clear and harder to understand. More importantly since a `query.PaymentChannelInfo` and `query.LedgerChannelInfo`'s status can [depend on the status of a virtual fund objective](../client/query/query.go#L180), the virtual fund protocol has to return a channel as updated when the objective completes. It ends up putting the burden on protocols to think about when we want notify the client, instead of when a channel is updated.

## Consequences

Return updates for every objective's `Related` when we crank an objective (and letting the client screen if an update is new) is less performant than the a protocol specifying which channels it has updated. However by paying this performance cost we can keep our protocols simpler and keep responsibilities separate.

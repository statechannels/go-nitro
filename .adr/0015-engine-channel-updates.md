# 0015 Engine Channel Updates

## Status

Review

## Context

To support channel update notifications the engine must communicate to the client that some channels have been updated.

## Decision

The engine now returns two new slices on an `EngineEvent`, `LedgerChannelUpdates` and `PaymentChannelUpdates`. These slices contain `query.PaymentChannelInfo` and `query.LedgerChannelInfo` for any updated channels.

The engine now populates these slices by inspecting an objective's `Related` channels after it's been cranked. We add an entry for each `Related` channel, without checking if it's been changed, which may result in duplicate updates being returned from the engine.

## Alternatives considered

### Protocols returning Updated Channels

Instead of assuming that all channels have been updated when an objective is cranked, protocols could report which channels they have updated.

Unfortunately this requires adding a lot of boilerplate code throughout the protocols. It can also be difficult to determine when a protocol should mark a channel as updated.

## Consequences

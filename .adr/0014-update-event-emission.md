# 0012 Channel Update Event Emission

## Status

Review

## Context

We want to enhance the our go-nitro RPC client and server with the ability to receive notifications whenever a channel is updated. This allows us to design a UI that can easily update when a channel changes.

Previously in [ADR 00012](./0012-event-emission.md) we decided to `close` channels to notify that objectives have been completed. While this approach work well for one-off events (that contain no information) with multiple subscribers. We need another approach that lets us broadcast updates to multiple consumers.

## Decision

A client store a slice of event listener `chans` for each channel id. A consumer can "subscribe" by calling a method which appends a `chan` to that slice and returns it. Then, Whenever the client receives a channel update from the engine it iterates on over the slice and sends an update to any event listener `chans`.

To prevent duplicate notifications being emitted, the client keeps track of the previous emitted channel update. A channel update notification is only emitted if the channel has changed from the previous emitted notification.

## Alternatives considered

### Using a Third Party Library

A third party library, like [EventBus](https://github.com/asaskevich/EventBus), could handle registering and notifying event listeners. However it wouldn't allow us to easily prevent duplicate notifications.

### Sync.Cond

While `sync.Cond` allows for broadcasting updates to multiple consumers, it requires the consumers and broadcaster to be goroutines. This means we'd have to spin up goroutines for each event listener.

## Consequences

[ADR 0012](./0012-event-emission.md) identified some drawbacks with this approach, specifically that if the subscriber is too late it can "miss" the event. This still remains an issue, a channel may be updated before a listener is subscribed.

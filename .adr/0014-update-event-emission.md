# 0014 Channel Update Event Emission

## Status

Review

## Context

We want to enhance the our go-nitro RPC client and server with the ability to receive notifications whenever a channel is updated. This allows us to design a UI that can easily update when a channel changes.

Previously in [ADR 0013](./0013-event-emission.md) we decided to `close` channels to notify consumers that an objective have been completed. Using `close` allowed us to broadcast a signal to multiple consumers at once. While this works well for the objective completed event, it does not work for channel updates. This is because:

1. `close` doesn't allow us to pass any information (like channel data)
2. `close` can only be used once as a signal on a channel.

To support channel updates we need another approach that let's us dispatch channel updates to multiple consumers.

## Decision

We will implement the "Slice of chans" pattern similiar to the one outlined in [ADR 0013](./0013-event-emission.md).

The nitro node will store a slice of event listener `chans` for each channel id. A consumer can "subscribe" by calling a method which appends a `chan` to that slice and returns it. Whenever the node receives a channel update from the engine it iterates on over the slice and sends an update to any event listener `chans`.

To prevent duplicate notifications being emitted, the node keeps track of the previous emitted channel update. A channel update notification is only emitted if the channel has changed from the previous emitted notification.

## Alternatives considered

### Using a Third Party Library

We could use a third-party library, like [EventBus](https://github.com/asaskevich/EventBus), to handle the work of pubsub. However by writing our own eventhandler logic we can adjust it to our needs. One example is comparing the new event against the previously emitted event to prevent duplicate events being fired.

### Sync.Cond

We considered this option again, but discarded it for the same reasons as in [ADR 0013](./0013-event-emission.md).

## Consequences

[ADR 0013](./0013-event-emission.md) identified some drawbacks with this approach, specifically that if the subscriber is too late it can "miss" the event. This still remains an issue, a channel may be updated before a listener is subscribed.

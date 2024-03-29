# 0013 Event Emission

## Status

Accepted

## Context

The existing mechanism for updates concerning objective updates is a go chan `ToApi()` on which the engine sends `EngineEvents`. The `go-nitro` RPC Client spawns a go routine which reads those events, and then sends the objective identifier down either a `completedObjectives` or `failedObjectives` chan.

The consuming application (h.f. consumer) would like to be able to synchronize its behaviour based on such updates. Currently, it is possible for the consumer to e.g. wait on a specific objective completing by reading from the `completedObjective` chan until the desired objective id is found.

The problem with this approach is that it is _not_ well suited for multiple consumers. If another module or goroutine in the consumer wants to wait on (say) a different objective, there is a likelihood that the first goroutine will read the identifier which the second goroutine is interested in.

For example, let objective with `id1` complete and then objective with `id2` complete:

```go
func waitFor(objectiveId) {
	for id := range client.completedObjectives {
		if id == objectiveId {
			return
		}
	}
}

waitFor(id2) // blocks until id1,id2 have been sent
waitFor(id1) // blocks forever since id1 will not be re-sent

```

## Alternatives considered

### Slice of `chans`

It is possible to make a sort-of Javascript-style event system by having the `go-nitro` node store a slice of `chans` for each objective id. Then, the consumer can "subscribe" by calling a method which appends a `chan` to that slice and returns it. Then, when the `node` learns from the engine that an objective has completed, it iterates through the slice and closes each `chan`.

This approach would not be idiomatic and has the following drawbacks:

- if the subscriber is too late it can "miss" the event
- trying to solve that issue with a cache leads to a race condition (a new subscription call reads the cache, infers the event has not fired, and then appends to the slice of `chans` _after_ the event fires and therefore still misses it)

### Condition variable

A `sync.Cond` https://pkg.go.dev/sync#Cond is "a rendezvous point for goroutines waiting for or announcing the occurrence of an event.". It would allow events be emitted more than once by using the `Broadcast()` method. For example, if we wanted to have an "objective updated' event stream.

There are many downsides to trying to build an event system with this type, however. Firstly, the event consumers would have to be goroutines. Secondly, no information can be passed with `Broadcast()`, so those goroutines would have to query information when they are woken up by `Broadcast()`. This introduces race conditions. Thirdly, usage of `sync.Cond` is [very low](https://lukechampine.com/cond.html) and there are proposals to deprecate it https://github.com/golang/go/issues/21165.

## Decision

We follow the hint above and use a single "broadcast `chan`" to unblock an abitrary number of consuming goroutines that have subscribed to (i.e. have a handle on) it.

Assuming a `chan` exists for each objective id, consumers simply do things like

```go
<-client.ObjectiveCompletedChan(id1) // unblocks when the chan corresponding to id1 is closed
<-client.ObjectiveCompletedChan(id2) // unblocks when the chan corresponding to id2 is closed
<-client.ObjectiveCompletedChan(id2) // unblocks when the chan corresponding to id2 is closed
```

This allows for may different patterns in the consumer. For example, each of these lines could be run in a separate goroutine -- they could be intentionally raced, or a `sync.WaitGroup` could be used to wait for them all in parallel and synchronize when they are all done.

### When to create broadcast `chans`

Above we assumed a `chan` exists for each objective id. We need to unpack this and consider when the `chan` will be allocated in memory. Since objective ids are not predictable in a permisionless protocol like Nitro, any per-channel object (such as the proposed broadcast `chans`) must be spun up on demand.

One approach might be to spin up the broadcast chan when the objective starts. But this will cause race conditions, when consumers want to subscribe before the `chan` exists.

Take the following example: a proposer may open a channel with a counterparty, and then _out of band_ send them a request for service provision referencing the channel of objective id. The proposee may then want to wait on that objective completing before starting to service the request. They subscribe to this "topic" before their `go-nitro` node knows anything about the objective.

The solution is to have the first reader or writer to the broadcast `chan` construct it (using e.g. `make(chan struct{})`). In practice this means using the `sync.Map` method `LoadOrStore`. Any later reader or writer can then (respectively) try to read or close the channel safely without risking a null pointer exception. Because `sync.Map` is concurrency safe, all of the necessary locking is taken care of automagically.

Please see code changes committed atomically with this ADR for the full reference code.

## Future considerations

1. We will want to roll this pattern out to other events in the codebase which will have multiple consumers.

2. An obvious thought might be: instead of closing the channel couldn't we just leave the channel open and send empty structs, allowing us to implement a stream of "objective updated" events? Unfotunately, we can't simply send empty structs on the channel, because they would be read by only one consumer at a time. That doesn't give us a true "broadcast" pattern -- more like a "signal" pattern. See the section above on "condition variable" for a pattern which similar to a chan which can be closed multiple times.

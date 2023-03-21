# 0012 Event Emission

## Status

Accepted

## Context

The existing mechanism for updates concerning objective updates is a go chan `ToApi()` on which the engine sends `EngineEvents`. The `go-nitro` Client spawns a go routine which reads those events, and then sends the objective identifier down either a `completedObjective` chan or `failedObjective`.

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

It is possible to make a sort-of Javascript-style event system by having the `go-nitro` Client store a slice of `chans` for each objective id. Then, the consumer can "subscribe" by calling a method which appends a `chan` to that slice and returns it. Then, when the `client` learns from the engine that an objective has copmleted, it iterates through the slice and closes each `chan`. 

This approach would not be idiomatic and has the following drawbacks:
* if the subscriber is too late it can "miss" the event
* trying to solve that issue with a cache leads to a race condition (a new subscription call reads the cache, infers the event has not fired, and then appends to the slice of `chans` _after_ the event fires and therefore still misses it)

### Condition variable

This is actually the most idiomatic approach if we want to emit events more than once. For example, if we wanted to have an "objective updated' event stream. See `sync.Cond` https://pkg.go.dev/sync#Cond .

Note in the godoc of this type: 

> For many simple use cases, users will be better off using channels than a Cond (Broadcast corresponds to closing a channel, and Signal corresponds to sending on a channel).



## Decision
We follow the hint above and use a single "broadcast `chan`" to unblock an abitrary number of consuming goroutines that have subscribed to (i.e. have a handle on) it. 

Assuming a `chan` exists for each objective id, consumers simply do things like

```go 
<-client.ObjectiveCompletedChan(id1) // unblocks when the chan corresponding to id1 is closed
<-client.ObjectiveCompletedChan(id2) // unblocks when the chan corresponding to id2 is closed
<-client.ObjectiveCompletedChan(id2) // unblocks when the chan corresponding to id2 is closed
```

This allows for may different patterns in the consumer. For example, each of these lines could be run in a separate goroutine -- they could be intentionally raced, or a `sync.WaitGroup` could be used to wait for them all in parallel and synchronize when they are all done. 

## Future considerations
We will want to roll this pattern out to other events in the codebase which will have multiple consumers. See the section above on "condition variable" for a pattern which can handle multiple emissions from the same event. 






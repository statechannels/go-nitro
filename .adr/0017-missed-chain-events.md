# 0017 Missed Chain Events

## Status

Review

## Context

It is critical that a go-nitro node can synchronise pertinent blockchain state, in order to successfully perform off chain protocols and protect against faulty challenges. For efficiency reasons, a node synchronises that state by listening to blockchain events. Missing those events can therefore cause the node's data to get out of sync with other nodes and the chain state. Events might be missed if the node has not been continuously running throughout the lifetime of any particular channel it is a participant in. For example, there may have been a period of downtime or loss of internet due to power loss or maintenance.

To address these concerns, a nitro node should check for any chain events it might have missed since it was last online, and queue those events to be processed as part of the node's initialization routines.

If the `chainservice` communicates with the `engine` when new blocks are processed, the node can store the `lastBlockSeen` in the `durablestore`. Then, when the node is initialized it can read the `lastBlockSeen` from the `store`, and search for any chain events that occurred between that block and the current block.

```go
		ecs.wg.Add(3)
		go ecs.listenForEventLogs(errChan, eventSub, eventChan, eventQuery)
		go ecs.listenForNewBlocks(errChan, newBlockSub, newBlockChan)
		go ecs.listenForErrors(errChan)

		// Search for any missed events emitted while this node was offline
		err = ecs.checkForMissedEvents(startBlock)
		if err != nil {
			return nil, err
		}
```

One concern with this design is how does the node ensure it processes chain events idempotently (i.e. how does it ensure the same chain event does not make changes to a `Channel`'s data multiple times?).

## Alternatives considered

### Put all `Channel` update protection logic in `chainservice`

One option is to trust the `chainservice` to always deliver `Channel` updates in order. However, putting data sanitation logic closer to the data itself is safer and less prone to duplicated code since `chainservice` is an interface that can have multiple implementations.

## Decision

Add logic to the `Channel` class to protect its data from stale updates rather than relying on an external component (`chainservice`) to only feed it sanitized data. This seems less error-prone and it could be useful to add a `Channel.LastUpdated` struct in case a user wants to know when the channel was last changed. The `chainservice` will still be expected to feed chain events in order, but the added `Channel` logic will act as another layer of protection against faulty data updates.

The `Channel` class should keep track of the `TxIndex` associated with each chain event in case there are multiple events related to a single channel within the same block. Also, an `ethereum.FilterQuery` cannot specify a `TxIndex` so it will be more computationally efficient to check the `TxIndex` in the `Channel` class rather than in the `chainservice` logic.

```go
type Channel struct {
	state.FixedPart
	Id      types.Destination
	MyIndex uint

	OnChain  OnChainData
	OffChain OffChainData

  LastChainUpdate ChainUpdateData
}

type ChainUpdateData struct {
        BlockNum uint64
        TxIndex uint64
}

func (c *Channel) UpdateWithChainEvent(event chainservice.Event) (*Channel, error) {
	if event.BlockNum() > c.LastChainUpdate.BlockNum && event.TxIndex() > c.LastChainUpdate.TxIndex {
		// Process event
		...

		// Update Channel.LastChainUpdate
		c.LastChainUpdate.BlockNum = event.BlockNum()
		c.LastChainUpdate.TxIndex  = event.TxIndex()
		return c, nil
	}
	return nil, fmt.Errorf("chain event older than channels last update")
}
```

To ensure no events are missed or processed out of order during the `chainservice` initialization, we can modify the `chainservice` init routine to acquire/release the `eventTracker`'s `mutex` as shown below:

```go
	// Prevent go routines from processing events before checkForMissedEvents completes
	ecs.eventTracker.mu.Lock()
	{
		ecs.wg.Add(3)
		go ecs.listenForEventLogs(errChan, eventSub, eventChan, eventQuery)
		go ecs.listenForNewBlocks(errChan, newBlockSub, newBlockChan)
		go ecs.listenForErrors(errChan)

		// Search for any missed events emitted while this node was offline
		err = ecs.checkForMissedEvents(startBlock)
		if err != nil {
			return nil, err
		}
	}
	ecs.eventTracker.mu.Unlock()
```

The `listenForEventLogs` go routine also tries to acquire the `mutex` before updating the `eventTracker.events` queue but will be blocked until `checkForMissedEvents` completes. This means that `listenForEventLogs` will still be able to detect new chain events, but will not add them to the queue or trigger event processing until `checkForMissedEvents` finishes queueing old events.

## Future considerations

1. Is there any situation where an older chain event SHOULD be allowed to update a `Channel` (i.e. chain re-org)?

# 0017 Missed Chain Events

## Status

Review

## Context

A user should be able to restart a nitro node without concern that it has missed chain events that caused the node's data to get out of sync with other nodes and the chain state. Similarly, a nitro node should be able to start for the first time and catch up on any relevant chain events that occurred before the node started. To address these concerns, a nitro node should check for any chain events it might have missed since it was last online, and queue those events to be processed as part of the node's initialization routines.

If the `chainservice` communicates with the `engine` when new blocks are processed, the node can store the `lastBlockSeen` in the `durablestore`. Then, when the node is initialized it can read the `lastBlockSeen` from the `store`, and search for any chain events that occurred between that block and the current block.

```go
// Search for any missed events emitted while this node was offline
err = ecs.checkForMissedEvents(startBlock)
if err != nil {
	return nil, err
}

ecs.wg.Add(3)
go ecs.listenForEventLogs(errChan, eventSub, eventChan, eventQuery)
go ecs.listenForNewBlocks(errChan, newBlockSub, newBlockChan)
go ecs.listenForErrors(errChan)
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

## Future considerations

1. Is there any situation where an older chain event SHOULD be allowed to update a `Channel` (i.e. chain re-org)?

2. Is there a chance any events are missed between `ecs.checkForMissedEvents(startBlock)` and `go ecs.listenForEventLogs(errChan, eventSub, eventChan, eventQuery)` during the `chainservice` initialization?

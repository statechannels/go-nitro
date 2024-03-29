# 0017 Missed Chain Events

## Status

Review

## Context

It is critical that a go-nitro node can synchronise pertinent blockchain state, in order to successfully perform off chain protocols and protect against faulty challenges. For efficiency reasons, a node synchronises that state by listening to blockchain events. Missing those events can therefore cause the node's data to get out of sync with other nodes and the chain state. Events might be missed if the node has not been continuously running throughout the lifetime of any particular channel it participates in. For example, there may have been a period of downtime or loss of internet due to power loss or maintenance.

To address these concerns, a nitro node should check for any chain events it might have missed since it was last online, and queue those events to be processed as part of the node's initialization routines.

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

There are a few decisions to make within this design:

1. How does the node protect against race conditions that could caused missed chain events during initialization?
2. How does the node ensure it processes chain events idempotently (i.e. how does it ensure the same chain event does not make changes to a `Channel`'s data multiple times?). This situation could occur if the node restarts, scans through blocks for events during initialization, and detects events that it has already processed.
3. By storing a `lastBlockNum` variable in the `store`, the `chainservice` will not have to scan through as many old blocks when it first initializes and calls `checkForMissedEvents`. When the node is initialized it can read the `lastBlockNumSeen` from the `store`, and search for any chain events that occurred between that block and the current block. If the node uses a `memstore`, then the `lastBlockNumSeen` will always be set to `0` when the node is first initialized. When/how should the `lastBlockNumSeen` be updated?

## Alternatives considered

The alternative design choices listed below relate to the same numbered design decisions from the previous section.

### [2] Put all `Channel` data protection logic in `chainservice`

The node could trust the `chainservice` to always deliver `Channel` updates in order. However, putting data sanitation logic closer to the data itself is safer and less prone to duplicated code since `chainservice` is an interface that can have multiple implementations.

### [3] Updating `lastBlockNum`: `chainservice` sends write on chan to `engine`

If the `chainservice` sends a new `blockNum` on a chan the `engine` each time a new block is mined, the `engine` can trigger a write to store to update the `lastBlockNum`. However, in environments where blocks are mined fast, writing to the `store` for every block could have a performance impact especially if the `store` updates are part of the `engine`'s main `run` loop.

## Decision

After considering the aforementioned alternatives, the following design decisions were made.

### [1] Use eventTracker mutex

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

### [2] Put `Channel` data protection logic in `Channel` class

Add logic to the `Channel` class to protect its data from stale updates rather than relying on an external component (`chainservice`) to only feed it sanitized data. This seems less error-prone and it could be useful to add a `Channel.LastUpdated` struct in case a user wants to know when the channel was last changed. The `chainservice` will still be expected to feed chain events in order, but the added `Channel` logic will act as another layer of protection against faulty data updates.

The `Channel` class should keep track of the `TxIndex` associated with each chain event in case there are multiple events related to a single channel within the same block.

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
  TxIndex  uint
}

func (c *Channel) isNewChainEvent(event chainservice.Event) bool {
	return event.BlockNum() > c.LastChainUpdate.BlockNum ||
		(event.BlockNum() == c.LastChainUpdate.BlockNum && event.TxIndex() > c.LastChainUpdate.TxIndex)
}

func (c *Channel) UpdateWithChainEvent(event chainservice.Event) (*Channel, error) {
	if !c.isNewChainEvent(event) {
		return nil, fmt.Errorf("chain event older than channel's last update")
	}
	// Process event
	...

	// Update Channel.LastChainUpdate
	c.LastChainUpdate.BlockNum = event.BlockNum()
	c.LastChainUpdate.TxIndex  = event.TxIndex()
	return c, nil
}
```

### [3] Updating `lastBlockNum`: `engine` periodically reads from `chainservice`

Instead of the `chainservice` alerting the `engine`, the `engine` can request the block number from the `chainservice`. This means the `lastBlockNumSeen` won't get updated every block, and could cause the node to try to process old events a second time during initialization. However, that's not a major concern since the node has protection in the `Channel` class (topic [2] above) to prevent it from making `Channel` updates based on stale data.

The benefit of this approach is simplicity and limited additional strain on the `engine`/`store`. Instead of adding a new chan or new event type on an existing `chainservice` --> `engine` chan, a `chainservice.GetLastConfirmedBlockNum()` function is called by the `engine` periodically. This means the node won't write to the `store` every block but every x seconds. The new `engine` run loop looks like this:

```go
	blockTicker := time.NewTicker(15 * time.Second)

	select {
		...
		case chainEvent := <-e.fromChain:
			res, err = e.handleChainEvent(chainEvent)
		case <-blockTicker.C:
			blockNum := e.chain.GetLastConfirmedBlockNum()
			err = e.store.SetLastBlockNumSeen(blockNum)
	}
```

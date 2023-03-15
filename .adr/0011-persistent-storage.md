# 0011 Implementing a Persistent Store

## Status

Accepted

## Context

For the go-nitro library to be useful in any kind of production setting it must have some way of persisting data to disk (or some form of storage). Our [memory store](../client/engine/store/memstore.go), is an in-memory store that makes use of `sync.Map` to store data.

We have a [well-defined interface](../client/engine/store/store.go) for our store; this means the implementation details of the store do not influence the rest of the client.

## Decision

The initial implementation of a persistent store will be done using [BuntDb](https://github.com/tidwall/buntdb).

> BuntDB is a low-level, in-memory, key/value store in pure Go. It persists to disk, is ACID compliant, and uses locking for multiple readers and a single writer.

The choice of BuntDB was motivated by a few factors:

- BuntDB uses a key/value interface, conceptually similar to the existing [memory store](../client/engine/store/memstore.go). This allows us to quickly implement a persistent store.
- BuntDB is ACID compliant and uses locking for reader/writers. We don't have to worry about creating malformed data and have some basic guarantees about the state of the data.
- BuntDB is file-based and lightweight. No additional services need to be installed or configured.

The selection of BuntDB is not a permanent decision. Once we establish some basic store benchmarks we can look at implementing a store using other kinds of storage/DBs (like SQL) to address performance.

### Sync Policy

BuntDB is an in-memory store, and how often it writes to disk is [configurable](https://github.com/tidwall/buntdb#durability-and-fsync). There are three options:

```
// SyncPolicy represents how often data is synced to disk.
type SyncPolicy int

const (
	// Never is used to disable syncing data to disk.
	// The faster and less safe method.
	Never SyncPolicy = 0
	// EverySecond is used to sync data to disk every second.
	// It's pretty fast and you can lose 1 second of data if there
	// is a disaster.
	// This is the recommended setting.
	EverySecond = 1s
	// Always is used to sync data after every write to disk.
	// Slow. Very safe.
	Always = 2
)
```

Using testground and the benchmark scenario we can see the difference between sync policies.

- [Sync always: 247 ms TTFP](http://34.168.92.245:3000/d/5OBBeW37k/time-to-first-payment?orgId=1&from=1678478595989&to=1678478699884)
- [Sync every second: 86.5 ms TTFP](http://34.168.92.245:3000/d/5OBBeW37k/time-to-first-payment?orgId=1&from=1678478669430&to=1678478825842&var-runId=cg5op0gnr2ghv5ug1m8g&var-jobCount=3&var-testDuration=1m0s&var-hubs=1&var-payees=5&var-jitter=1&var-latency=10&var-payers=2&var-payeepayers=0&var-nitroVersion=v0.0.0-20230310171721-486f70744942&var-storeSyncFrequency=1)
- [Never Sync Policy: 84.1 ms TTFP](http://34.168.92.245:3000/d/5OBBeW37k/time-to-first-payment?orgId=1&from=1678479249909&to=1678479339937&var-runId=cg5o96gnr2ghv5ug1m60&var-jobCount=3&var-testDuration=1m0s&var-hubs=1&var-payees=5&var-jitter=1&var-latency=10&var-payers=2&var-payeepayers=0&var-nitroVersion=v0.0.0-20230310171721-486f70744942&var-storeSyncFrequency=1)

**Note:** We could also implement our own sync policy by using the `Never SyncPolicy` and manually triggering the write to file ourselves (like at the end of the run loop.) Since we are using different BuntDB databases we could also set `SyncPolicy` per database.

Initially we will use the `EverySecond` policy as it is the default, and quite performant. This does expose us to possibility of losing data; we should investigate this further (with more comprehensive tests) and look at tuning this appropriately.

## Consequences

- We need to establish some benchmarks so we can understand how the choice of BuntDB affects performance. We should especially consider the role of the hub, who could be dealing with a large amount of traffic.
- We need to establish some tests that simulate a client crashing and see the impact of sync policies. We should especially consider the role of the hub, as it will probably be writing to the store more frequently.
- We should look at mostly deprecating the existing `MemStore` in favour of a store using BuntDB.

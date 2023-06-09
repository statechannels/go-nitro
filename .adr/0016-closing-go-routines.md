# 0016 Closing Go-routines

## Status

Review

## Context

In our codebase we have a few structs that use long-running go-routines to handle tasks asynchronously (often things like sending out messages/notifications). However we don't always consider how to stop and clean up after these go-routines, and we aren't always consistent on how we approach it. In some scenarios these long-running go-routines can continue running after `Close` has been called. This can easily introduce subtle race conditions.

Here's a list of structs that spin up long-running go-routines.

- [The RPC client](https://github.com/statechannels/go-nitro/blob/0b5fa37613363720c91c115c3de252a39b1b1f0a/rpc/client.go#L142)
- [The RPC server](https://github.com/statechannels/go-nitro/blob/0b5fa37613363720c91c115c3de252a39b1b1f0a/rpc/server.go#L223)
- [Eth Chain service](https://github.com/statechannels/go-nitro/blob/0b5fa37613363720c91c115c3de252a39b1b1f0a/client/engine/chainservice/eth_chainservice.go#L244)
- [The API client](https://github.com/statechannels/go-nitro/blob/0b5fa37613363720c91c115c3de252a39b1b1f0a/client/client.go#L87)

## Decision

**Note:** I use the term struct as a shorthand for a struct with a long-running go-routine, like the examples above.

When a struct's `Close` function is called, it should block and not return until:

- all go-routines it owns stop executing.
- any closeable resources it owns are closed.

By enforcing these constraints a running go-routine can be guaranteed that it's parent struct is in a "running" state. This rules out a large class of race conditions and errors such as a go-routine attempting to use a parent struct's resource that has been closed.

To enforce these constraints we should follow this pattern in a struct's `Close` function:

1. Signal any go-routines we own to exit.
2. Wait until all go-routines have completed execution.
3. Close any resources it owns.

### Step 1: Signal go-routines to exit

Long-running `go-routines` need some kind of trigger to stop executing. A common and simple pattern we often use is simply closing the chan the go-routine is consuming from.

```golang
	toRoutine := make(chan int)

	go func() {
		for v := range toRoutine {
			doSomething(v)
		}
		// This logic runs once the channel is closed
		doCleanup()
	}()

	// This closes the channel and signals the go routine to exit
	close(toRoutine)

```

This works, however a [cancelable context](https://cs.opensource.google/go/go/+/go1.20.5:src/context/context.go;l=238) provides some benefits over this:

- A buffered chan will only get closed once it's buffer is read. This means a go-routine will read all the buffered entries before it finishes executing. By using a context we can halt the execution almost immediately.
- It makes go-routine cleanup logic explicit and easy to see. It's now just a case statement for `ctx.Done`.
- Minor but it would allow us to use other context features, such as timeouts, in the future.

Due to these benefits, and the limited use of go-routines, we should update our structs to use a cancelable context to signal to a go-routine to stop executing.

```golang
	toRoutine := make(chan int)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for {
			select {
			case v := <-toRoutine:
				doSomething(v)
			case <-ctx.Done():
				doCleanup()
			}
		}
	}()

	// This triggers the goroutine to exit
	cancel()
```

### 2: Wait until all go-routines have completed execution.

After we have signalled our go-routines to exit we should wait for them to complete. The easiest way to accomplish this is with a `sync.WaitGroup`

```golang
wg := sync.WaitGroup{}
	toRoutine := make(chan int)
	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)
	go func() {
		for {
			select {
			case v := <-toRoutine:
				doSomething(v)
			case <-ctx.Done():
				doCleanup()
				wg.Done()
			}
		}
	}()

	// This triggers the goroutine to exit
	cancel()

	// This blocks until the goroutine has exited
	wg.Wait()
```

### Step 3: Close Resources

Once a struct has waited for all go-routines to finish executing, it can dispose of any resources like network connections or child structs. We do this by calling `Close` on any child structs that implement [io.Closer interface](https://pkg.go.dev/io#Closer). **In general, if a child struct implements the `Closer` interface, we should consider calling it in our struct's `Close`**

## Prior Art

A example of this pattern can be found in the libp2p codebase, such as the [mdns service Close function](https://github.com/libp2p/go-libp2p/blob/c9de1665054229bdfd40884cd0b893744ec8ef7e/p2p/discovery/mdns/mdns.go#L75).

```golang

func (s *mdnsService) Close() error {
	s.ctxCancel()
	if s.server != nil {
		s.server.Shutdown()
	}
	s.resolverWG.Wait()
	return nil
}
```

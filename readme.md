<h1 align="center">
<div><img src="https://protocol.statechannels.org/img/favicon.ico"><br>
go-nitro
</h1>
Implementation of nitro protocol in golang.

---

# Getting started (MacOS)

Install golang

```
brew install golang
```

Install [golangci-lint](https://golangci-lint.run):

```
brew install golangci-lint
brew upgrade golangci-lint
```

Make sure GOPATH is set:

```
echo $GOPATH
```

You should see `$HOME/go`.

### For developers

To format:

```shell
gofmt -w .
```

To lint:

```shell
golangci-lint run
```

To build:

```shell
go build ./...
```

To run tests:

```shell
go test ./...
```

To view docs website:

```shell
godoc --http :6060
```

and navigate to http://localhost:6060/pkg/github.com/statechannels/go-nitro/
To remove unused dependencies (CI checks will fail unless this is a no-op):

```shell
go mod tidy
```

### License

Dual-licensed under [MIT](https://opensource.org/licenses/MIT) + [Apache 2.0](http://www.apache.org/licenses/LICENSE-2.0)

---

## Roadmap

The following roadmap gives an idea of the various packages that compose the `go-nitro` module, and their implementation status:

```bash
├── channel 🚧                 # query the latest supported state of a channel
│   └── state ✅               # generate and recover signatures on state updates
│       ├── outcome ✅         # define how funds are dispersed when a channel closes
├── client 🚧                  # exposes an API to the consuming application
│   └── engine 🚧              # coordinate the client components, runs the protocols
│       ├── chainservice 🚧    # watch the chain and submit transactions
│       ├── messageservice 🚧  # send and recieves messages from peers
│       └── store 🚧           # store keys, state updates and other critical data
├── protocols 🚧
│   ├── interfaces.go ✅       # specify the interface of our protocols
│   ├── direct-fund ✅         # fund a channel on-chain
│   ├── direct-defund 🚧       # defund a channel on-chain
│   ├── virtual-fund 🚧        # fund a channel off-chain through an intermediary
│   └── virtual-defund 🚧      # defund a channel off-chain through an intermediary
└── types 🚧                   # basic types and utility methods
```

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

```bash
├── channel 🚧 # query the latest supported state of a channel
│   └── state ✅ # sign a state
│       ├── outcome ✅ # construct and manipulate an outcome
├── client 🚧 # runs the off chain protocols
│   └── engine 🚧 # coordinates the client components
│       ├── chainservice 🚧 # watches the chain and submits transactions
│       ├── messageservice 🚧 # sends and recieves messages from peers
│       └── store 🚧 # stores keys and state updates
├── protocols 🚧
│   └── interfaces.go ✅
    |__ direct-funder.go 🚧 # specifies how clients fund a channel on chain
└── types 🚧
```

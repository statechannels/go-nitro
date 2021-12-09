<h1 align="center">
<div><img src="https://protocol.statechannels.org/img/favicon.ico"><br>
go-nitro
</h1>
Implementation of nitro protocol in golang.

---

### Contributing

See [contributing.md](./contributing.md)
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
│   ├── virtual-fund 🚧        # fund a channel off-chain through one or more intermediaries
│   └── virtual-defund 🚧      # defund a channel off-chain through one or more intermediaries
└── types 🚧                   # basic types and utility methods
```

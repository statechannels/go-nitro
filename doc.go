// Go-nitro can be run as a system service with an RPC api. Go-nitro's default configuration looks for a local blockchain network on port 8545 with chainid 1337.
// If such a network is available is running, go-nitro can be started from the root directory with
//
//	go run .
//
// Or, built to an executable binary with
//
//	go build -o gonitro
//
// Go nitro accepts the following command flags, which can also be displayed via `go run . -help` (or `gonitro -help` for the build binary).
// Usage of ./nitro-rpc-server:
//
//	-chainid int
//	      Specifies the chain id of the chain. (default 1337)
//	-chainurl string
//	      Specifies the url of a RPC endpoint for the chain. (default "ws://127.0.0.1:8545")
//	-deploycontracts
//	      Specifies whether to deploy the adjudicator and create2deployer contracts.
//	-msgport int
//	      Specifies the tcp port for the  message service. (default 3005)
//	-naaddress string
//	      Specifies the address of the nitro adjudicator contract. Default is the address computed by the Create2Deployer contract. (default "0xC6A55E07566416274dBF020b5548eecEdB56290c")
//	-pk string
//	      Specifies the private key used by the node. Default is Alice's private key. (default "2d999770f7b5d49b694080f987b82bbc9fc9ac2b4dcc10b0f8aba7d700f69c6d")
//	-rpcport int
//	      Specifies the tcp port for the rpc server. (default 4005)
//	-usedurablestore
//	      Specifies whether to use a durable store or an in-memory store.
//	-usenats
//	      Specifies whether to use NATS or http/ws for the rpc server.
//
// You can make remote procedure calls like so:
//
//	curl -X POST \
//	  -H 'Content-Type: application/json' \
//	  -d '{"jsonrpc":"2.0","id":1,"method":"get_address","params":{}}' \
//	  http://localhost:4005/api/v1
//
// but see  [github.com/statechannels/go-nitro/rpc] or https://github.com/statechannels/nitro-gui/tree/main/packages/nitro-rpc-client for an RPC client to do so programmatically.
package main

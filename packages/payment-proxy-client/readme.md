# Payment Proxy Client

This is basic client UI designed to work with a [reverse payment proxy](../../cmd/start-reverse-payment-proxy/). It provides some basic functionality to request a payload from the payment proxy and handles downloading the file or displaying the error(such as a 402- Payment Required).

It relies on a go-nitro rpc server network (which can be started using [this script](https://github.com/statechannels/go-nitro/blob/5b8c876d34638f9c322cf332bf758f5e9c284907/scripts/start-rpc-servers.go))

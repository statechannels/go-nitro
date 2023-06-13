# Boost Integration Demo

This is basic demo designed to work with a [forked version of boost](https://github.com/statechannels/boost) that requires payments before serving up a file via `booster-http`.

It provides some basic functionality to request a payload from `booster-http` and handles downloading the file or displaying the error(such as a 402- Payment Required). The selected payment channel will be passed into the request to `booster-http` for the payment check. A `Pay` button makes a payment on selected payment channel.

It relies on a go-nitro rpc server network (which can be started using [this script](https://github.com/statechannels/go-nitro/blob/5b8c876d34638f9c322cf332bf758f5e9c284907/scripts/start-rpc-servers.go)) and a running instance of our [forked version of `booster-http` alongside a `lotus` devnet](https://github.com/statechannels/boost).

**Note** The [hardcoded default payload id](./src/App.tsx#L43) is based on a local file and will not work on your machine! You must enter the payload id of a locally deployed file.

## Version 0

Version 0 works against [this](https://github.com/statechannels/boost/releases/tag/nitro-boost-demo-v0) tag in our boost fork. It uses a very simple and not-secure payment model:
![](./version-0-seq.png)

- HTTP requests have a channel id as a query param
- Retrievals are charged at a flat rate
- The retrieval provider will serve the request as long as the channel with supplied id has had a payment
- This is vulernable to replays!
  ![](./version-0.gif)

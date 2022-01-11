# Application channels

An "application channel" is intended to be used for the purpose of exchanging assets within the context of some user application. (This is in contrast to a "ledger" channel, which is used as a "private ledger" between a small, fixed number of peers, which can be used to fund multi-hop "virtual" channels without interacting with a blockchain.)

## Payment channels
"Payment channels" are the simplest example of an application channel. Often, the only type of state update allowed is to change the channel participant balances by signing a (new) state with a different `outcome` and a larger `version`:

# Virtual fund off-chain protocol

## Single hop case

Take three actors Alice, Bob and Irene. Given a ledger channel `L` between Alice and Irene and a ledger channel `L'` between Bob and Irene, the clients send and wait on messages as shown in the below sequence diagram in order to create and fund a virtual channel `V`:

![](./virtual-fund-sequence-diagram.svg)

The diagram is generated at https://sequencediagram.org/. The source code for this diagram is co-located in this folder.

## Off-chain protocols

This package defines the off-chain protocols for managing channels. These protocols are abstract rules which nodes should follow in order to reach a desired goal.

To implement the protocols, we make use of an abstraction we call the `Objective`. Objectives are named for the goal which the off-chain protocol aspires to achieve.

The following table shows some example objectives and their implemetation status:

| Objective type                             | Implemented |
| ------------------------------------------ | ----------- |
| `direct-fund`                              | x           |
| `direct-defund`                            | x           |
| [`virtual-fund`](./virtual-fund/readme.md) | x           |
| `virtual-defund`                           | x           |
| `challenge`                                |             |

The set of objectives comprises the functional core of a go-nitro node. They expose only _pure_ functions -- but otherwise take on as much responsibility as possible, leaving only a small amount of responsibility to an imperative shell.

Each `Objective` type is implemented by a Go struct implementing the `Objective` interface. Each instance of an `Objective` type will hold its own data.

The imperative shell will typically be responsible for:

- spawning new objective instances of one type or another
- reading and persisting objectives to and from a store.
- listening to peers and the blockchain for events relevant to the objective
- "updating" objectives by passing events to their `Update(event)` function (an updated copy is returned)
- "cranking" objectives by calling `Crank()` (side effects are returned)
- executing side effects

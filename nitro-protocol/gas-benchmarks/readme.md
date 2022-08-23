### Benchmarking strategy

To gauge the efficiency of our smart contract(s) implementation, we have gas consumption figures which are easy to refer to -- they are committed in [`./gasResults.json`](./gasResults.json).

In order for these figures to be reliable, we want to have a deterministic process for generating them: we have defined several scenarios of interest and arrange for the relevant transactions to be applied to a consistent blockchain state without interfering with each other.

To achieve that, we spin up a local hardhat instance. After each test _case_, we revert the blockchain to a snapshot. That snapshot contains deployed contracts, and has mined no other transactions.

### Updating benchmark results file

Any change to the smart contract source code (i.e. the set of `.sol` files) is likely to lead to an increase or decrease in gas consumption. Run

`npm run benchmark:update`

to update [`./gasResults.json`](./gasResults.json) with the figures that describe the updated smart contract implementation. Failing to do this will result in a continuous integration failure on your pull request, and you will not be able to merge your changes into `main`.

(Under the hood, this runs `exportBenchmarkResults.ts`).

### Showing benchmarks diff

The continuous integration suite will regenerate the gas consumption figures (in memory) and report a failure if they do not match those committed in [`./gasResults.json`](./gasResults.json). You can run this process yourself with:

```
npm run benchmark:diff
```

Running this command will provide a color-coded report showing the difference in gas consumption between the current smart contract implementation and the implementation which `npm run benchmark:update` was previously run with.

(Under the hood, this runs `benchmark.test.ts`).

You can also use `npm run benchmark` will invoke both commands, thus showing benchmark difference and then updating `gasResults.json` (if necessary).

Test are run with `maxConcurrency = 1`. For maximum efficiency, tests are run in a single file to prevent having to restart hardhat.

### Scenario shorthand

We use simple diagrams / pictures to show the "exit strategy" of a certain actor in a particular scenario. This actor is using the on chain API to recover their own funds, which they are free to do at any point but must rely on in the case that one or more counterparties stops cooperating.

Ths diagrams show how money moves after each on chain operation. Take the example of a Bob, who is receving money from Alice in a virtual payment channel. Alice has gone offline during the virtual payment channel execution. Bob's stratgey is to

1.  Redeem a payment voucher on chain.
2.  Reclaim his money into his ledger channel.
3.  Transfer money out of the ledger channel.

The diagram show how the money owed to Bob is locked into a funding path, and gradually becomes unlocked as he performs the on chain operations. Channels which are off-chain are written with letters -- once they finalize on chain we put parentheses around the relevant letter: `(L)`.

Initially, Bob's deposit (as well as Ingrid the intermdiary's deposit) is is held in escrow against the ledger channel `L` in the adjudicator. We write this as `⬛ -> L`.
This channel is off chain, and funds a virtual payment channel `V` (by way of a guarantee). The virtual channel in turn allocated funds directly to Bob's external address which we write as 👨.

Therefore, initially we have

```
⬛ ->  L  ->  V  -> 👨
```

Bob's first step is to finalize both `L` and `V` (in any order). He does this by calling `challenge` for both channels. After a timeout, the resulting state is written:

```
⬛ -> (L) -> (V) -> 👨
```

Next Bob calls `reclaim`, which changes an allocation to Bob in `V` to an alloaction to Bob in `L`:

```
⬛ -> (L) --------> 👨
```

Finally, he calls `transferAllAssets` on `L`, resulting in him regaining complete control of his assets:

```
⬛ ---------------> 👨
```

Although this final operation may result in payouts to other parties, that is not reflected in the diagrams and is simply a side-effect of Bob's actions (and not his objective).

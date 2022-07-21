### Benchmarking strategy

To gauge the efficiency of our smart contract(s) implementation, we have gas consumption figures which are easy to refer to -- they are committed in [`./gasResults.json`](./gasResults.json).

In order for these figures to be reliable, we want to have a deterministic process for generating them: we have defined several scenarios of interest and arrange for the relevant transactions to be applied to a consistent blockchain state without interfering with each other.

To achieve that, we spin up a local hardhat instance. After each test _case_, we revert the blockchain to a snapshot. That snapshot contains deployed contracts, and has mined no other transactions.

### Updating benchmark results file

Any change to the smart contract source code (i.e. the set of `.sol` files) is likely to lead to an increase or decrease in gas consumption. Run

`npm run benchmark:update`

to update [`./gasResults.json`](./gasResults.json) with the figures that describe the updated smart contract implementation. Failing to do this will result in a continuous integration failure on your pull request, and you will not be able to merge your changes into `main`.

### Showing benchmarks diff

The continuous integration suite will regenerate the gas consumption figures (in memory) and report a failure if they do not match those committed in [`./gasResults.json`](./gasResults.json). You can run this process yourself with:

```
npm run benchmark:diff
```

Running this command will provide a color-coded report showing the difference in gas consumption between the current smart contract implementation and the implementation which `npm run benchmark:update` was previously run with.

You can also use `npm run benchmark` will invoke both commands, thus showing benchmark difference and then updating `gasResults.json` (if necessary).

Test are run with `maxConcurrency = 1`. For maximum efficiency, tests are run in a single file to prevent having to restart hardhat.

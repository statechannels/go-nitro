### Benchmarking options

Benchmarking is separated into two parts: checking the difference between current gas spendings and the previous ones (stored in `gasResults.json`) and updating `gasResults.json` with current benchmark results. This is done via `npm run benchmark:diff` and `npm run benchmark:update` respectively, whereas `npm run benchmark` will invoke both of them, thus showing benchmark difference and updating `gasResults.json` simultaneously.

### Benchmarking strategy

We want to have deterministic benchmark tests that always see a consistent blockchain state and do not interfere with each other.

To achieve that, we spin up a local hardhat instance. After each test _case_, we revert the blockchain to a snapshot. That snapshot just contains deployed contracts, and has mined no other transactions. Test are run with `maxConcurrency = 1`. For maximum efficiency, you should therefore run tests in a single file to prevent having to restart hardhat.

### Showing benchmarks diff

To get color-coded gas spending difference between current on-chain implementation and the previous one we use testing with `jest`. Tests are located in `benchmark.test.ts` file, which can only be invoked with specification of `../config/jest/jest.gas-benchmarks.config.js`as a config. This is done for `jestSetup.ts` to be called after `jest` has started, but before tests are executed. Basically `jestSetup.ts` set up snapshot reverting rules, extends `jest` with `toConsumeGas` check method and defining `challengeChannelAndExpectGas` function, which is used in tests. It is worth mentioning that `jestSetup.ts` also uses results from `localSetup.ts`, so latest file is also gets executed.

Functions from `fixtures.ts` are used in tests.

### Updating benchmark results file

The logic to update `gasResults.json` file is similar to `benchmark.test.ts`, except that write benchmark results instead of comparing them. Thus, we do not need the `jest` and all setup happens in `localSetup.ts`.

Functions from `fixtures.ts` are used in the script.

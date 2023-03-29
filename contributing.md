# Contributing

Contributions to go-nitro are welcome from anyone.

Magmo (at Consensys Mesh) are developing go-nitro under contract for the Filecoin Foundation. As a consequence, we reserve discretionary veto rights on feature additions, architectural decisions, etc., which may hinder progress toward the specific aims of this contract.

## Contribution lifecycle

Magmo uses [Github Projects](https://github.com/statechannels/go-nitro/projects/1) for go-nitro project management. We will aim to keep a handful of "good first issues" in [the pipeline](https://github.com/statechannels/go-nitro/issues?q=is%3Aopen+is%3Aissue+label%3A%22good+first+issue%22), which are a good place to start.

Before spending lots of time working on an issue, consider asking us for feedback via [Slack](https://statechannels.slack.com/archives/C02J81JFD3J), the [issue tracker](https://github.com/statechannels/go-nitro/issues), or [email](magmo@mesh.xyz). We would love to help make your contributions more successful!

Pull requests (internal or external) will be reviewed by one or two members of Magmo. We aim to review most contributions within one business day.

We include a PR template to serve as a **guideline**, not a **rule**, to communicate the code culture we wish to maintain in this repository.

## Style

When in doubt, defer to the [Effective Go](https://go.dev/doc/effective_go) document [CodeReviewComments](https://github.com/golang/go/wiki/CodeReviewComments) as "style guides".

# Prerequisites (MacOS instructions)

### Install golang

```
brew install golang
```

### Install [golangci-lint](https://golangci-lint.run):

```
brew install golangci-lint
brew upgrade golangci-lint
```

### Install [staticcheck](http://staticcheck.io)

```
go install honnef.co/go/tools/cmd/staticcheck@latest
```

### Make sure GOPATH is set:

```
echo $GOPATH
```

You should see `$HOME/go`.

# Building and Testing

To build:

```shell
go build ./...
```

To run tests:

```shell
go test ./...
```

# Documenation website

This is built using a variant of `mkdocs`. First, follow [these instructions](https://squidfunk.github.io/mkdocs-material/getting-started/#installation) to install. Then run

```
mkdocs serve
```

and navigate to http://localhost:8000 .

## Viewing Godocs website

Run

```shell
godoc --http :6060
```

and navigate to http://localhost:6060/pkg/github.com/statechannels/go-nitro/

# Pre PR checks:

Please execute the following on any branch that's ready for review or merge.

### staticcheck:

```shell
staticcheck ./...
```

### format:

```shell
gofumpt -w .
```

### lint:

```shell
golangci-lint run
```

### remove unused dependencies:

```shell
go mod tidy
```

# Debugging Tests

VS code is used to debug tests. To start a debugging session in VS code:

- Ensure you have the [go extension](https://marketplace.visualstudio.com/items?itemName=golang.Go) installed
- Open the test file.
- Open the `Run and Debug` section.
- Run the `Debug Test` configuration.

With the extension it is also possible to start a debugging session right from a test function.

# Generating Go bindings

We use `solc` and `abigen` to generate Go bindings from our `*.sol` source files.

This is achieved by running the `generate_adjudicator_bindings.sh` script at the top level of the repo. Because our `*.sol` files depend on external projects via `node_modules`, to run this script successfully you must:

- have successfully run `npm install` in the `nitro-protocol` directory.
- have [solc](https://docs.soliditylang.org/en/v0.8.17/installing-solidity.html) installed at the correct version (currently 0.8.17, see the CI or linting config for a hint if you think it may have changed)
- have [abigen](https://geth.ethereum.org/docs/install-and-build/installing-geth) (a tool shipped with go-ethereum) installed.

The resulting Go bindings file is _checked-in_ to the repository. Although it is auto-generated from on-chain source code, it effectively forms part of the off-chain source code.

If you alter the contracts, you should regenerate the bindings file at check it in. A github action will run which will check that you have done this correctly.

TIP: if you find that the github action still reports a mismatch despite regenerating the bindings, this may be due to the action using the "test merge" of your PR (rather than the tip of your branch). Try rebasing your branch.

# Testground

We run deeper tests of the code on your PR using a hosted [testground](https://docs.testground.ai/) runner. You may see a comment appear with links to dashboards containing performance metrics and statistics. The tests may also be run locally. For more information, see our testground [test plan repository](https://github.com/statechannels/go-nitro-testground).

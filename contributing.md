# Contributing

Contributions to go-nitro are welcome from anyone.

Magmo (at Consensys Mesh) are developing go-nitro under contract for the Filecoin Foundation. As a consequence, we are obligated to maintain a certain level of quality.

## Contribution lifecycle

Magmo uses [Github Projects](https://github.com/statechannels/go-nitro/projects/1) for go-nitro project management. We will aim to keep a handful of "good first issues" in [the pipeline](https://github.com/statechannels/go-nitro/issues?q=is%3Aopen+is%3Aissue+label%3A%22good+first+issue%22), which are a good place to start.

Before spending lots of time working on an issue, consider asking us for feedback. We would love to help make your contributions more successful!

Pull requests (internal or external) will be reviewed by one or two members of Magmo. We aim to review most contributions within one business day.

We include a PR template to serve as a **guideline**, not a **rule**, to communicate the code culture we wish to maintain in this repository.

## Style

As of now, Magmo is new to the go language, and has no specific opinions about go style. As a rule of thumb, we defer to the [Effective Go](https://go.dev/doc/effective_go) document [CodeReviewComments](https://github.com/golang/go/wiki/CodeReviewComments) as "style guides".

## Judgement

You are welcome to use your judgement and ignore our PR template, Effective Go, and CodeReviewComments when appropriate.

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

# Viewing Godocs website

```shell
godoc --http :6060
```

and navigate to http://localhost:6060/pkg/github.com/statechannels/go-nitro/

# Pre PR checks:

Please execute the following on any branch that's ready for review or merge.

### format:

```shell
gofmt -w .
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

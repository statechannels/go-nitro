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

## Pre-commit Hooks

There are pre-commit hooks to run `gofmt`, `golangci-lint` and `go mod tidy`. To get these to run `pre-commit` must be installed.

### Install pre-commit

```shell
brew install pre-commit
```

### Install git hook script

From the repository directory run

```shell
pre-commit install
```

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

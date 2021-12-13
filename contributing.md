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
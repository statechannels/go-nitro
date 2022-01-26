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

- Have the test file open in the editor.
- Open `Run and Debug` in VS code.
- Run the `Debug Test` configuration.

The [go extension](https://marketplace.visualstudio.com/items?itemName=golang.Go) also supports debugging and can be used to start debugging session right from a test function.

package main

import "github.com/statechannels/go-nitro/internal/infra"

func main() {
	err := infra.InitializeNitroNetwork()
	if err != nil {
		panic(err)
	}
}

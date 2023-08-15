package main

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	// Generate a new private key
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	fmt.Printf("privateKey: %x\n", privateKeyBytes)

	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	fmt.Printf("address:    %v\n", address)
}

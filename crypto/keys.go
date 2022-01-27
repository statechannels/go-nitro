package crypto

import (
	"crypto/ecdsa"
	"log"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/statechannels/go-nitro/types"
)

// GeneratePrivateKeyAndAddress generates a pseudo-random ECDSA and its corresponding Ethereum address.
func GeneratePrivateKeyAndAddress() (types.Bytes, types.Address) {
	channelSecretKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}
	channelSecretKeyBytes := crypto.FromECDSA(channelSecretKey)

	publicKey := channelSecretKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	return channelSecretKeyBytes, address
}

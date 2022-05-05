package crypto

import (
	"crypto/ecdsa"
	"log"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/statechannels/go-nitro/types"
)

// GeneratePrivateKeyAndAddress generates a pseudo-random ECDSA secret key and its corresponding Ethereum address.
func GeneratePrivateKeyAndAddress() (types.Bytes, types.Address) {
	secretKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}
	return crypto.FromECDSA(secretKey), getAddressFromSecretKey(*secretKey)
}

// GetAddressFromSecretKeyBytes computes the Ethereum address corresponding to the supplied private key.
func GetAddressFromSecretKeyBytes(secretKeyBytes []byte) types.Address {
	secretKey, err := crypto.ToECDSA(secretKeyBytes)
	if err != nil {
		log.Fatal(err)
	}
	return getAddressFromSecretKey(*secretKey)
}

// GetAddressFromSecretKey computes the Ethereum address corresponding to the supplied private key.
func getAddressFromSecretKey(secretKey ecdsa.PrivateKey) types.Address {
	publicKey := secretKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	return crypto.PubkeyToAddress(*publicKeyECDSA)
}

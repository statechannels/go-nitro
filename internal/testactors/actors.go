// Package testactors exports peers with vanity addresses: with corresponding keys, names and virtual funding protocol roles.
package testactors

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/crypto"
	"github.com/statechannels/go-nitro/types"
)

type ActorName string

type Actor struct {
	PrivateKey []byte
	Role       uint
	Name       ActorName
	Port       uint
}

func (a Actor) Destination() types.Destination {
	return types.AddressToDestination(a.Address())
}

func (a Actor) Address() types.Address {
	return crypto.GetAddressFromSecretKeyBytes(a.PrivateKey)
}

const START_PORT = 3200

// Alice has the address 0xAAA6628Ec44A8a742987EF3A114dDFE2D4F7aDCE
var Alice Actor = Actor{
	common.Hex2Bytes(`2d999770f7b5d49b694080f987b82bbc9fc9ac2b4dcc10b0f8aba7d700f69c6d`),
	0,
	"Alice",
	START_PORT + 0,
}

// Bob has the address 0xBBB676f9cFF8D242e9eaC39D063848807d3D1D94
var Bob Actor = Actor{
	common.Hex2Bytes(`0279651921cd800ac560c21ceea27aab0107b67daf436cdd25ce84cad30159b4`),
	2,
	"Bob",
	START_PORT + 1,
}

// Brian has the address 0xB2B22ec3889d11f2ddb1A1Db11e80D20EF367c01
var Brian Actor = Actor{
	common.Hex2Bytes("0aca28ba64679f63d71e671ab4dbb32aaa212d4789988e6ca47da47601c18fe2"),
	2,
	"Brian",
	START_PORT + 2,
}

// Irene has the address 0x111A00868581f73AB42FEEF67D235Ca09ca1E8db
var Irene Actor = Actor{
	common.Hex2Bytes(`febb3b74b0b52d0976f6571d555f4ac8b91c308dfa25c7b58d1e6a7c3f50c781`),
	1,
	"Irene",
	START_PORT + 3,
}

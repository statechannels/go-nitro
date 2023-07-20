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
// peerId: 16Uiu2HAmSjXJqsyBJgcBUU2HQmykxGseafSatbpq5471XmuaUqyv
var Alice Actor = Actor{
	common.Hex2Bytes(`2d999770f7b5d49b694080f987b82bbc9fc9ac2b4dcc10b0f8aba7d700f69c6d`),
	0,
	"alice",
	START_PORT + 0,
}

// Bob has the address 0xBBB676f9cFF8D242e9eaC39D063848807d3D1D94
// peerId: 16Uiu2HAmJDxLM8rSybX78FH51iZq9PdrwCoCyyHRBCndNzcAYMes
var Bob Actor = Actor{
	common.Hex2Bytes(`0279651921cd800ac560c21ceea27aab0107b67daf436cdd25ce84cad30159b4`),
	2,
	"bob",
	START_PORT + 1,
}

// Ivan has the address 0xA8d2D06aCE9c7FFc24Ee785C2695678aeCDfd7A0
// peerId: 16Uiu2HAm1hgN2MkrhGen8JPrBBYyXACbZtfqJmraN53XiHYCeoFi
var Ivan Actor = Actor{
	common.Hex2Bytes("1ea91a2724b40fb8fed6a3648d49e9431996c09744fb841b718377fb0700f3e7"),
	2,
	"ivan",
	START_PORT + 2,
}

// Irene has the address 0x111A00868581f73AB42FEEF67D235Ca09ca1E8db
// peerId: 16Uiu2HAmHntR3SGeS7iF2tdeNBefSahXBhmTrqVozVLHydxzkaZn
var Irene Actor = Actor{
	common.Hex2Bytes(`febb3b74b0b52d0976f6571d555f4ac8b91c308dfa25c7b58d1e6a7c3f50c781`),
	1,
	"irene",
	START_PORT + 3,
}

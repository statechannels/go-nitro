package client_test

import "github.com/ethereum/go-ethereum/common"

var aliceKey = common.Hex2Bytes(`2d999770f7b5d49b694080f987b82bbc9fc9ac2b4dcc10b0f8aba7d700f69c6d`)
var alice = common.HexToAddress(`0xAAA6628Ec44A8a742987EF3A114dDFE2D4F7aDCE`)

var ireneKey = common.Hex2Bytes(`febb3b74b0b52d0976f6571d555f4ac8b91c308dfa25c7b58d1e6a7c3f50c781`)
var irene = common.HexToAddress(`0x111A00868581f73AB42FEEF67D235Ca09ca1E8db`)

var bobKey = common.Hex2Bytes(`0279651921cd800ac560c21ceea27aab0107b67daf436cdd25ce84cad30159b4`)
var bob = common.HexToAddress(`0xBBB676f9cFF8D242e9eaC39D063848807d3D1D94`)

var brianKey = common.Hex2Bytes("0aca28ba64679f63d71e671ab4dbb32aaa212d4789988e6ca47da47601c18fe2") //nolint:unused
var brian = common.HexToAddress("0xB2B22ec3889d11f2ddb1A1Db11e80D20EF367c01")                       //nolint:unused

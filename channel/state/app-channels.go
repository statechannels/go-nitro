package state

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/types"
)

var ALICE = "Alice"

// var aliceAddr = types.Address()

// WANT:
// var aliceAddr = common.Hex2Bytes("0xaaaa4Ae81F235c5FeAD170cE30e0325E2ac331f9")
var aliceAddr = types.Address{}

// differentParticipants := []types.Address{
// 	common.HexToAddress(`0x760bf27cd45036a6C486802D30B5D90CfFBE31FE`),
// 	common.HexToAddress(`0xF5A1BB5607C9D079E46d1B3Dc33f257d937b43BD`),
// }
var alicePK = common.Hex2Bytes(`c347a80d6c938a862262bb2d0af8a812510ffc59ebf62556d7b2e7d31668d701`)

type (
	AppState struct {
		fp        FixedPart
		Alice     uint
		Bob       uint
		TurnNum   uint
		Signature Signature
	}
)

// AlicePaysBob mimics what would happen if a client controlling Alice's private key wanted to pay Bob some amount
func (s AppState) AlicePaysBob(amount uint, pk []byte) (AppState, error) {
	var next AppState

	next.fp = s.fp

	next.TurnNum += 1
	next.Alice = s.Alice - amount
	next.Bob = s.Bob + amount
	next.sign(pk)

	return next, nil
}

func (prev AppState) BobValidatesPayment(payment AppState) error {
	ch1, _ := prev.fp.ChannelId()
	ch2, _ := payment.fp.ChannelId()
	if ch1 != ch2 {
		return errors.New("bob receiving: wrong channelId")
	}

	if payment.TurnNum != prev.TurnNum+1 {
		return errors.New("alice sent the wrong turn number")
	}

	// TODO: Need to check for overflow/underflow errors
	received := payment.Bob - prev.Bob
	spent := prev.Alice - payment.Alice
	if spent != received {
		return errors.New("funds went missing")
	}

	// TODO: Would need
	signer, err := payment.asState().RecoverSigner(payment.Signature)
	if err != nil {
		return fmt.Errorf("bob: couldn't recover signer %w", err)
	}

	// if signer != aliceAddr {
	if signer != aliceAddr {
		return fmt.Errorf("bob: malicious signature detected")

	}

	return nil
}

func (s AppState) asState() State {
	fp := s.fp
	exit := outcome.Exit{}
	return State{
		fp.ChainId, fp.Participants, fp.ChannelNonce, fp.AppDefinition, fp.ChallengeDuration, types.Bytes{}, exit, big.NewInt(int64(s.TurnNum)), false}
}

func (s *AppState) sign(pk []byte) error {
	sig, err := s.asState().Sign(pk)

	if err != nil {
		return err
	}

	s.Signature = sig

	return nil
}

func Demo() {
	fp := FixedPart{}
	s0 := AppState{fp, 2, 10, 1, Signature{}}

	s1, err := s0.AlicePaysBob(10, alicePK)
	if err != nil {
		panic(err)
	}

	s0.BobValidatesPayment(s1)

	s2, err := s1.AlicePaysBob(20, alicePK)
	if err != nil {
		panic(err)
	}

	s1.BobValidatesPayment(s2)

	fmt.Print(s2)
}

// var bytesZero types.Bytes = types.Bytes{}
// var o0 outcome.SingleAssetExit = outcome.SingleAssetExit{}

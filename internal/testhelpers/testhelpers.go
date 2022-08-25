// Package testhelpers contains functions which pretty-print test failures.
package testhelpers

import (
	"bytes"
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/internal/testactors"
	"github.com/statechannels/go-nitro/payments"
	"github.com/statechannels/go-nitro/protocols"
)

// Copied from https://github.com/benbjohnson/testing

// makeRed sets the colour to red when printed
const makeRed = "\033[31m"

// makeBlack sets the colour to black when printed.
// as it is intended to be used at the end of a string, it also adds two linebreaks
const makeBlack = "\033[39m\n\n"

// Assert fails the test immediately if the condition is false.
// If the assertion fails the formatted message will be output to the console.
func Assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf(makeRed+"%s:%d: "+msg+makeBlack, append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// Ok fails the test immediately if an err is not nil.
// If the error is not nil the message containing the error will be outputted to the console
func Ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf(makeRed+"%s:%d: unexpected error: %s"+makeBlack, filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// Equals fails the test if want is not deeply equal to got.
// Equals uses reflect.DeepEqual to compare the two values.
func Equals(tb testing.TB, want, got interface{}) {
	if !reflect.DeepEqual(want, got) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf(makeRed+"%s:%d:\n\n\texp: %#v\n\n\tgot: %#v"+makeBlack, filepath.Base(file), line, want, got)
		tb.FailNow()
	}
}

// AssertStateSentToEveryone asserts that ses contains a message for every participant but from
func AssertStateSentToEveryone(t *testing.T, ses protocols.SideEffects, expected state.SignedState, from testactors.Actor, allActors []testactors.Actor) {
	for _, a := range allActors {
		if a.Role != from.Role {
			AssertStateSentTo(t, ses, expected, a)
		}
	}
}

// AssertStateSentTo asserts that ses contains a message for the participant
func AssertStateSentTo(t *testing.T, ses protocols.SideEffects, expected state.SignedState, to testactors.Actor) {
	for _, msg := range ses.MessagesToSend {
		toAddress := to.Address()
		if bytes.Equal(msg.To[:], toAddress[:]) {
			for _, ss := range msg.SignedStates() {
				Equals(t, ss.Payload, expected)
			}
		}
	}
}

func AssertProposalSent(t *testing.T, ses protocols.SideEffects, sp consensus_channel.SignedProposal, to testactors.Actor) {

	if len(ses.MessagesToSend) != 1 {
		fmt.Printf("%+v\n", ses.MessagesToSend)
	}
	Assert(t, len(ses.MessagesToSend) == 1, "expected one message")

	found := false

	msg := ses.MessagesToSend[0]
	for _, p := range msg.SignedProposals() {
		found = found || p.Payload.Proposal.Equal(&sp.Proposal) && p.Payload.TurnNum == sp.TurnNum
	}
	toAddress := to.Address()
	Assert(t, found, "proposal %+v not found in signed proposals %+v", sp.Proposal, msg.SignedProposals())
	Assert(t, bytes.Equal(msg.To[:], toAddress[:]), "exp: %+v\n\n\tgot%+v", msg.To.String(), to.Address().String())

}

// SignState generates a signature on the signed state with the supplied key, and adds that signature.
// If an error occurs the function panics
func SignState(ss *state.SignedState, secretKey *[]byte) {
	sig, err := ss.State().Sign(*secretKey)
	if err != nil {
		panic(fmt.Errorf("SignAndAdd failed to sign the state: %w", err))
	}
	err = ss.AddSignature(sig)
	if err != nil {
		panic(fmt.Errorf("SignAndAdd failed to sign the state: %w", err))
	}
}

// AssertVoucherSentToEveryone asserts that ses contains a message for every participant but from
func AssertVoucherSentToEveryone(t *testing.T, ses protocols.SideEffects, expected *payments.Voucher, from testactors.Actor, allActors []testactors.Actor) {

	for _, a := range allActors {
		if a.Role != from.Role {
			AssertVoucherSentTo(t, ses, *expected, a)
		}
	}
}

// AssertVoucherSentTo asserts that ses contains a message for the participant
func AssertVoucherSentTo(t *testing.T, ses protocols.SideEffects, expected payments.Voucher, to testactors.Actor) {
	for _, msg := range ses.MessagesToSend {
		toAddress := to.Address()
		if bytes.Equal(msg.To[:], toAddress[:]) {
			for _, v := range msg.Vouchers() {
				Equals(t, v.Payload, expected)
			}
		}
	}
}

// Package testhelpers contains functions which pretty-print test failures.
package testhelpers

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
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

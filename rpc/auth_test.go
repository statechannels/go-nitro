package rpc

import (
	"errors"
	"testing"
	"time"
)

func TestValidAuthToken(t *testing.T) {
	token, err := generateAuthToken("1", allPermissions)
	if err != nil {
		t.Fatal(err)
	}

	err = checkTokenValidity(token, permSign, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAuthTokenMissingPermission(t *testing.T) {
	token, err := generateAuthToken("1", []permission{permRead})
	if err != nil {
		t.Fatal(err)
	}

	err = checkTokenValidity(token, permSign, time.Hour)
	if !errors.Is(err, errMissingPermission) {
		t.Fatal("expected errMissingPermission, got", err)
	}
}

func TestExpiredAuthToken(t *testing.T) {
	token, err := generateAuthToken("1", allPermissions)
	if err != nil {
		t.Fatal(err)
	}

	err = checkTokenValidity(token, permSign, time.Duration(0))
	if !errors.Is(err, errExpiredToken) {
		t.Fatal("expected errExpiredToken, got", err)
	}
}

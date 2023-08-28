package rpc

import (
	"errors"
	"testing"
	"time"
)

func TestValidAuthToken(t *testing.T) {
	token, err := generateAuthToken(allPermissions)
	if err != nil {
		t.Fatal(err)
	}

	err = checkTokenValidity(token, permSign, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAuthTokenMissingPermission(t *testing.T) {
	token, err := generateAuthToken([]permission{permRead})
	if err != nil {
		t.Fatal(err)
	}

	err = checkTokenValidity(token, permSign, nil)
	if !errors.Is(err, errMissingPermission) {
		t.Fatal("expected errMissingPermission, got", err)
	}
}

func TestExpiredAuthToken(t *testing.T) {
	token, err := generateAuthToken(allPermissions)
	if err != nil {
		t.Fatal(err)
	}

	zeroDuration := time.Duration(0)
	err = checkTokenValidity(token, permSign, &zeroDuration)
	if !errors.Is(err, errExpiredToken) {
		t.Fatal("expected errExpiredToken, got", err)
	}
}

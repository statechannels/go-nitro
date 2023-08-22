package rpc

import (
	"errors"
	"testing"
)

func TestValidAuthToken(t *testing.T) {
	token, err := generateAuthToken(allPermissions)
	if err != nil {
		t.Fatal(err)
	}

	err = checkPermission(token, permSign)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAuthTokenMissingPermission(t *testing.T) {
	token, err := generateAuthToken([]permission{permRead})
	if err != nil {
		t.Fatal(err)
	}

	err = checkPermission(token, permSign)
	if !errors.Is(err, errMissingPermission) {
		t.Fatal("expected errMissingPermission, got", err)
	}
}

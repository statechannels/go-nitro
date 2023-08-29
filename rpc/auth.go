package rpc

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// todo: the private key should not be hardcoded
var rpcPK = []byte("rpcPK")

type permission string

const permissionKey = "perm"
const (
	permNone permission = "none"
	permRead permission = "read"
	permSign permission = "sign"
)

var allPermissions = []permission{permRead, permSign}

var (
	errInvalidSigningMethod = errors.New("invalid signing method")
	errInvalidToken         = errors.New("invalid token")
	errExpiredToken         = errors.New("token has expired")
	errInvalidPermissions   = errors.New("token has invalid permissions")
	errInvalidPermission    = errors.New("token has an invalid permission")
	errMissingPermission    = errors.New("token is missing permission")
)

var invalidIAtFormat = "invalid issued at: %w"

// generateAuthToken generates a JWT token that a client uses to authenticate with the server for restricted endpoints
// subject is the itentifier of the client for which the token is generated
func generateAuthToken(subject string, p []permission) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims[permissionKey] = p
	// the keys are defined by https://datatracker.ietf.org/doc/html/rfc7519
	claims["iat"] = time.Now().Unix()
	claims["sub"] = subject
	return token.SignedString(rpcPK)
}

// checkTokenValidity takes a JWT token, verifies that the token is valid and that the token contains the required permission
func checkTokenValidity(tokenString string, requiredPermission permission, validDuration time.Duration) error {
	if requiredPermission == permNone {
		return nil
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errInvalidSigningMethod
		}
		return rpcPK, nil
	})
	if err != nil {
		return err
	}

	if !token.Valid {
		return errInvalidToken
	}

	claims := token.Claims.(jwt.MapClaims)

	// Check expiration
	iAt, err := claims.GetIssuedAt()
	if err != nil {
		return fmt.Errorf(invalidIAtFormat, err)
	}

	if time.Now().After(iAt.Add(validDuration)) {
		return errExpiredToken
	}

	// Check permissions
	permissions, ok := claims[permissionKey].([]interface{})
	if !ok {
		return errInvalidPermissions
	}

	for _, p := range permissions {
		sp, ok := p.(string)
		if !ok {
			return errInvalidPermission
		}

		pp := permission(sp)

		if pp == requiredPermission {
			return nil
		}
	}

	return errMissingPermission
}

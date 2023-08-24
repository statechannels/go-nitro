package rpc

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

// todo: the private key should not be hardcoded
var rpcPK = []byte("rpcPK")

type permission string

const (
	permNone permission = "none"
	permRead permission = "read"
	permSign permission = "sign"
)

var allPermissions = []permission{permRead, permSign}

var (
	errInvalidSigningMethod = errors.New("invalid signing method")
	errInvalidToken         = errors.New("invalid token")
	errInvalidPermissions   = errors.New("token has invalid permissions")
	errInvalidPermission    = errors.New("token has an invalid permission")
	errMissingPermission    = errors.New("token is missing permission")
)

// generateAuthToken generates a JWT token that a client uses to authenticate with the server for restricted endpoints
func generateAuthToken(p []permission) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["permissions"] = p
	return token.SignedString(rpcPK)
}

// verifyPermission takes a JWT token, verifies that the token is valid and that the token contains the required permission
func checkPermission(tokenString string, requiredPermission permission) error {
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
	permissions, ok := claims["permissions"].([]interface{})
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

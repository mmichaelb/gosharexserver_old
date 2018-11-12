package user

import (
	"encoding/hex"
)

// AuthorizationToken wraps the byte representation of the authorization token to get a string representation.
type AuthorizationToken []byte

// String returns the string representation in the form of the hex encoded string.
func (token AuthorizationToken) String() string {
	return hex.EncodeToString(token)
}

// AuthorizationTokenFromString returns the parsed AuthorizationToken or nil if there was an error while parsing the
// token.
func AuthorizationTokenFromString(stringToken string) AuthorizationToken {
	if rawToken, err := hex.DecodeString(stringToken); err == nil {
		return AuthorizationToken(rawToken)
	}
	return nil
}

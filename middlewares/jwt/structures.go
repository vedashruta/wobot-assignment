package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	m *model
)

type model struct {
	subject       string
	issuer        string
	audience      jwt.ClaimStrings
	expiry        time.Duration
	signingMethod jwt.SigningMethod
	secret        []byte
	skip          map[string]struct{}
}

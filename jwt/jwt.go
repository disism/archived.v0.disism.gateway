package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

//var secret = []byte("new_disism.com.hvturingga.test_secret")

// https://datatracker.ietf.org/doc/html/rfc7519
// JWT https://datatracker.ietf.org/doc/html/rfc7519#section-4

var now = time.Now

type JWT struct {
	*jwt.RegisteredClaims
}

type Option func(j *JWT)

func WithJwtIssuer(issuer string) Option {
	return func(j *JWT) {
		j.Issuer = issuer
	}
}

func WithSubject(subject string) Option {
	return func(j *JWT) {
		j.Subject = subject
	}
}

func WithJwtAudience(audience []string) Option {
	return func(j *JWT) {
		j.Audience = audience
	}
}

func WithJwtExpiresAt(expiresAt time.Time) Option {
	return func(j *JWT) {
		j.ExpiresAt = &jwt.NumericDate{
			Time: expiresAt,
		}
	}
}

func WithJwtNotBefore(notBefore time.Time) Option {
	return func(j *JWT) {
		j.NotBefore = &jwt.NumericDate{
			Time: notBefore,
		}
	}
}

func WithJwtIssuedAt(issuedAt time.Time) Option {
	return func(j *JWT) {
		j.IssuedAt = &jwt.NumericDate{
			Time: issuedAt,
		}
	}
}

func WithJwtID(id string) Option {
	return func(j *JWT) {
		j.ID = id
	}
}

// NewJWT creates a new JWT instance with the provided options.
//
// It accepts variadic JwtOption arguments to configure the JWT instance.
// The function returns a pointer to the created JWT instance.
func NewJWT(opts ...Option) *JWT {
	s := &JWT{
		RegisteredClaims: &jwt.RegisteredClaims{
			IssuedAt: &jwt.NumericDate{Time: now()},
		},
	}
	for _, srv := range opts {
		srv(s)
	}
	return s
}

// GenerateToken generates a JWT token with the given secret.
//
// Parameters:
// - secret: The secret used to sign the token.
//
// Returns:
// - string: The generated JWT token.
// - error: An error if the token generation fails.
func (j *JWT) GenerateToken(secret string) (string, error) {
	s, err := jwt.
		NewWithClaims(
			jwt.SigningMethodHS256,
			*j.RegisteredClaims,
		).
		SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return s, nil
}

// JWTParse parses a JWT token and returns the registered claims.
//
// The function takes a `bearer` string as a parameter, which represents
// the JWT token. It returns a pointer to `jwt.RegisteredClaims` and an
// error. The `jwt.RegisteredClaims` struct represents the registered
// claims of a JWT token.
func JWTParse(bearer string) (*jwt.RegisteredClaims, error) {
	t, _, err := new(jwt.Parser).ParseUnverified(bearer, &jwt.RegisteredClaims{})
	if err != nil {
		return nil, fmt.Errorf("parse token error: %v", err)
	}
	claims, ok := t.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return nil, fmt.Errorf("failed to get JWT claims")
	}

	return claims, nil
}

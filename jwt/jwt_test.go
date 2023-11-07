package jwt

import (
	"testing"
)

func TestGenerateToken(t *testing.T) {
	// Test case: Generate token with a valid secret
	secret := "12345678901234"
	// Create a JWT instance
	j := NewJWT(WithJwtIssuer("disism.com"), WithJwtID("476947492886283552"))
	token, err := j.GenerateToken(secret)
	if err != nil {
		t.Errorf("Error generating token: %v", err)
	}
	// Verify that the token is not empty
	if token == "" {
		t.Error("Generated token is empty")
	}
	t.Log(token)
}

func TestJWTParse(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodnR1cmluZ2dhQGRpc2lzbS5jb20iLCJzdWIiOiJkaXNpc20uY29tIiwiaWF0IjoxNjkyMTc4MTU4LCJqdGkiOiI0NzQyMDU1OTE3MDM2NTE2MTYifQ.fXit2suXpAC8bZAfR0D1ckk3bzR6K356KzabjBcg8XA"

	claims, err := JWTParse(token)
	if err != nil {
		t.Errorf("Error parsing token: %v", err)
	}
	t.Log(claims)
}

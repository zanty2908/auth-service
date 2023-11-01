package token

import (
	"net/http"
)

// Maker is an interface for managing tokens
type Maker interface {

	// CreateToken creates a new token for a specific username and duration
	CreateTokenPair(aud Audience, platform Platform, userId, role, email, phone string, data map[string]interface{}) (access, refresh string, accessPayload, refreshPayload *Payload, err error)

	// VerifyToken checks if the token is valid or not
	VerifyToken(token string) (*Payload, error)
	VerifyRefreshToken(token string) (*Payload, error)

	// GetTokenPayload from header
	GetTokenPayload(r *http.Request) (*Payload, error)
}

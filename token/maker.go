package token

import (
	"net/http"
	"time"
)

// Maker is an interface for managing tokens
type Maker interface {

	// CreateToken creates a new token for a specific username and duration
	CreateTokenPair(data interface{}) (access, refresh string, accessPayload, refreshPayload *Payload, err error)

	// VerifyToken checks if the token is valid or not
	VerifyToken(token string) (*Payload, error)
	VerifyRefreshToken(token string) (*Payload, error)

	// GetTokenPayload from header
	GetTokenPayload(r *http.Request) (*Payload, error)

	// For custom key
	CreateTokenWithKey(key string, data interface{}, duration time.Duration) (string, *Payload, error)
	VerifyTokenWithKey(key, token string) (*Payload, error)
}

package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Different types of error returned by the VerifyToken function
var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

// Payload contains the payload data of the token
type Payload struct {
	ID        string      `json:"id"`
	IssuedAt  time.Time   `json:"iat"`
	ExpiresAt time.Time   `json:"exp"`
	Data      interface{} `json:"data"`
}

// NewPayload creates a new token payload with a specific username and duration
func NewPayload(data interface{}, duration time.Duration, id *string) (*Payload, error) {
	var tokenID string
	if id == nil {
		newID, err := uuid.NewRandom()
		if err != nil {
			return nil, err
		}
		tokenID = newID.String()
	} else {
		tokenID = *id
	}

	payload := &Payload{
		ID:        tokenID,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(duration),
		Data:      data,
	}
	return payload, nil
}

// Valid checks if the token payload is valid or not
func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiresAt) {
		return ErrExpiredToken
	}
	return nil
}

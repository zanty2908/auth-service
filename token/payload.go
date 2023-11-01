package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Different types of error returned by the VerifyToken function
var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token is expired")
)

// Payload contains the payload data of the token
type Payload struct {
	Id        string                 `json:"jti,omitempty"`
	Audience  string                 `json:"aud,omitempty"`
	Platform  string                 `json:"platform,omitempty"`
	ExpiresAt time.Time              `json:"exp,omitempty"`
	IssuedAt  time.Time              `json:"iat,omitempty"`
	Subject   string                 `json:"sub,omitempty"`
	Role      string                 `json:"role,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Email     string                 `json:"email,omitempty"`
	Phone     string                 `json:"phone,omitempty"`
}

// NewPayload creates a new token payload with a specific username and duration
func NewPayload(aud Audience, platform Platform, userId, role, email, phone string, data map[string]interface{}, expDuration time.Duration, id *string) (*Payload, error) {
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

	curTime := time.Now()
	expiresAt := curTime.Add(expDuration)

	payload := &Payload{
		Id:        tokenID,
		Platform:  string(platform),
		Audience:  string(aud),
		ExpiresAt: expiresAt,
		Subject:   userId,
		IssuedAt:  curTime,
		Data:      data,
		Role:      role,
		Email:     email,
		Phone:     phone,
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

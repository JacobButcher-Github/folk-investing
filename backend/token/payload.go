package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Payload contains payload data of the token
type Payload struct {
	ID        uuid.UUID `json:"id"`
	UserLogin string    `json:"user_login"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

// NewPayload creates a new token payload with a specific user login and duration
func NewPayload(userLogin string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	payload := &Payload{
		ID:        tokenID,
		UserLogin: userLogin,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload, nil
}

// Valid checks if token payload is valid
func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return errors.New("token has expired")
	}
	return nil
}

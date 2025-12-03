package token

import "time"

// Maker is interface for managing tokens

type Maker interface {
	// CreateToken creates a new token for a specific username and duration
	CreateToken(userLogin string, duration time.Duration) (string, error)

	// VerifyToken checks if token is valid
	VerifyToken(token string) (*Payload, error)
}

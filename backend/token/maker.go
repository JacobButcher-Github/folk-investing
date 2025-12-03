package token

import "time"

// Maker is interface for managing tokens

type Maker interface {
	// CreateToken creates a new token for a specific username and duration
	CreateToken(userID int64, userLogin string, duration time.Duration) (string, *Payload, error)

	// VerifyToken checks if token is valid
	VerifyToken(token string) (*Payload, error)
}

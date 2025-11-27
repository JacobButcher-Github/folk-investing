package db

import (
	"database/sql"
)

// Store provides all functions to execute db queries and transactions
type Store struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new Store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

//buying stock

//selling stock

//create user

//multiple stock information updates

//multiple settings at the same time (that shouldn't need this ngl)

package db

import (
	"context"
	"time"
)

// CreateUserTxParams contains input parameters of create user transaction
type CreateUserTxParams struct {
	CreateUserParams
}

// CreateUserTxResult is result of the create user transaction
type CreateUserTxResult struct {
	User User
}

// CreateUserTx performs create user action and AfterCreate function defined in the CraeteUserTxParams
func (store *Store) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error) {
	var result CreateUserTxResult
	err := store.retryableTx(ctx,
		5,
		5*time.Second,
		func(q *Queries) error {
			var err error

			result.User, err = q.CreateUser(ctx, arg.CreateUserParams)
			if err != nil {
				return err
			}

			return nil
		})
	return result, err
}

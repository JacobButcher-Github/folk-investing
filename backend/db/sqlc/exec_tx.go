package db

import (
	"context"
	"database/sql"
	"fmt"
)

// ExecTx executes a function within a database transaction
func (store *Store) execTx(ctx context.Context, txOption *sql.TxOptions, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, txOption)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

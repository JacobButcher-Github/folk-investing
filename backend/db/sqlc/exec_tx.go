package db

import (
	"context"
	"fmt"
	"time"
)

// execTx executes a function within a database transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
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

// retryableTx tries to execute fn within a transaction and retries on SQLITE_BUSY.
func (store *Store) retryableTx(
	ctx context.Context,
	maxRetries int,
	delay time.Duration,
	fn func(q *Queries) error,
) error {
	for i := 0; i <= maxRetries; i++ {

		// Try transaction
		err := store.execTx(ctx, fn)
		if err == nil {
			return nil
		}

		if err.Error() == "database is locked (5) (SQLITE_BUSY)" || err.Error() == "database is locked (517)" {
			if i == maxRetries {
				return fmt.Errorf("transaction failed after retries: %w", err)
			}
			time.Sleep(delay)
			continue
		}

		// Any other error: do not retry
		return err
	}

	return nil
}

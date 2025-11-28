package db

import (
	"context"
	"database/sql"
	"fmt"
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

// ExecTx executes a function within a database transaction
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

// BuyStockTxParams contains input parameters of buy stock transaction
type BuyStockTxParams struct {
	UserID    int64  `json:"user_id"`
	StockName string `json:"stock_name"`
	Amount    int64  `json:"amount"`
}

// BuyStockTxResult is result of the buy stock transaction
type BuyStockTxResult struct {
	User      User      `json:"user"`
	UserStock UserStock `json:"user_stock"`
	Stock     Stock     `json:"stock"`
}

// BuyStockTx performs a money subtraction from an account and adds stock to an account
// update stock amount in user stock and update user dollar and cents using latest stock data in a single transaction
func (store *Store) BuyStockTx(ctx context.Context, arg BuyStockTxParams) (BuyStockTxResult, error) {
	var result BuyStockTxResult

	return result, err
}

//selling stock

//create user

//multiple stock information updates

//multiple settings at the same time (that shouldn't need this ngl)

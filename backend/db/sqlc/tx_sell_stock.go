package db

import (
	"context"
	"database/sql"
)

// SellStockTxParams contains input parameteres of sell stock transaction
type SellStockTxParams struct {
	UserID  int64 `json:"user_id"`
	StockID int64 `json:"stock_id"`
	Amount  int64 `json:"amount"`
}

// SellStockTxResult is result of the sell stock transaction
type SellStockTxResult struct {
	User      User      `json:"user"`
	UserStock UserStock `json:"user_stock"`
}

// SellStockTx performs money addition to User and subtracts stocks from associated UserStock
func (store *Store) SellStockTx(ctx context.Context, arg SellStockTxParams) (SellStockTxResult, error) {
	var result SellStockTxResult
	err := store.execTx(ctx,
		&sql.TxOptions{
			Isolation: sql.LevelSerializable,
		},
		func(q *Queries) error {
			//TODO: implement sell user stock here.

			//get cost of stock being sold
			//get user selling stock
			//update amount of stock held
			//delete UserStock if hit 0 in amount.
			//update money of user
			return nil
		})

	return result, err
}

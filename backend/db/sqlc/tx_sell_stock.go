package db

import (
	"context"
	"database/sql"
	"time"
)

// SellStockTxParams contains input parameters of sell stock transaction
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
	err := store.retryableTx(ctx,
		5,
		2*time.Second,
		func(q *Queries) error {
			//get cost of stock being sold
			stockData, err := q.GetStockData(ctx, GetStockDataParams{
				StockID: arg.StockID,
				Limit:   1,
			})
			if err != nil {
				return err
			}

			stockCost := stockData.ValueDollars*100 + stockData.ValueCents

			//get user selling stock
			user, err := q.GetUserFromId(ctx, arg.UserID)
			if err != nil {
				return err
			}
			userMoney := user.Dollars*100 + user.Cents

			//update amount of stock held
			userStock, err := q.GetUserStock(ctx, GetUserStockParams{
				UserID:  arg.UserID,
				StockID: arg.StockID,
			})
			if err != nil {
				return err
			}

			updatedUserStock, err := q.UpdateUserStock(ctx, UpdateUserStockParams{
				Quantity: userStock.Quantity - arg.Amount,
				UserID:   arg.UserID,
				StockID:  arg.StockID,
			})

			result.UserStock = updatedUserStock

			//delete UserStock if hit 0 in amount.
			if updatedUserStock.Quantity == 0 {
				q.DeleteUserStock(ctx, DeleteUserStockParams{
					StockID: updatedUserStock.StockID,
					UserID:  updatedUserStock.UserID,
				})
			}

			//update money of user
			updatedAmount := userMoney + stockCost*arg.Amount
			updatedUser, err := q.UpdateUser(ctx, UpdateUserParams{
				HashedPassword: sql.NullString{String: "", Valid: false},
				Dollars:        sql.NullInt64{Int64: updatedAmount / 100, Valid: true},
				Cents:          sql.NullInt64{Int64: updatedAmount % 100, Valid: true},
				UserLogin:      user.UserLogin,
			})
			if err != nil {
				return err
			}

			result.User = updatedUser
			return nil
		})

	return result, err
}

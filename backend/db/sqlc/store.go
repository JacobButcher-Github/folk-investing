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

// BuyStockTxParams contains input parameters of buy stock transaction
type BuyStockTxParams struct {
	UserID  int64 `json:"user_id"`
	StockID int64 `json:"stock_id"`
	Amount  int64 `json:"amount"`
}

// BuyStockTxResult is result of the buy stock transaction
type BuyStockTxResult struct {
	User      User      `json:"user"`
	UserStock UserStock `json:"user_stock"`
}

// BuyStockTx performs a money subtraction from User and adds stock to associated UserStock
func (store *Store) BuyStockTx(ctx context.Context, arg BuyStockTxParams) (BuyStockTxResult, error) {
	var result BuyStockTxResult
	err := store.execTx(ctx,
		&sql.TxOptions{
			Isolation: sql.LevelSerializable,
		},
		func(q *Queries) error {
			//check if the associated UserStock exists.
			userStock, err := q.GetUserStock(ctx, GetUserStockParams{
				UserID:  arg.UserID,
				StockID: arg.StockID,
			})

			if err != nil {
				// Does not exist, create UserStock
				if err == sql.ErrNoRows {

					userStock, err = q.CreateUserStock(ctx, CreateUserStockParams{
						UserID:   arg.UserID,
						StockID:  arg.StockID,
						Quantity: 0,
					})
					if err != nil {
						return err
					}
				} else {
					//unknown error, handle by returning
					return err
				}
			}
			//Get the cost of the stock being purchased
			stockData, err := q.GetStockData(ctx, GetStockDataParams{
				StockID: userStock.StockID,
				Limit:   1,
			})
			if err != nil {
				return err
			}

			stockCost := stockData.ValueDollars*100 + stockData.ValueCents

			//Get user that's purchasing stock
			user, err := q.GetUserFromId(ctx, userStock.UserID)
			if err != nil {
				return err
			}

			userMoney := user.Dollars.Int64*100 + user.Cents.Int64

			//Update the money of the user
			newUserMoney := userMoney - stockCost*arg.Amount

			updatedUser, err := q.UpdateUser(ctx, UpdateUserParams{
				Dollars:   sql.NullInt64{Int64: newUserMoney / 100, Valid: true},
				Cents:     sql.NullInt64{Int64: newUserMoney % 100, Valid: true},
				UserLogin: user.UserLogin,
			})
			if err != nil {
				return err
			}

			//update value on UserStock.
			updatedUserStock, err := q.UpdateUserStock(ctx, UpdateUserStockParams{
				Quantity: userStock.Quantity + arg.Amount,
				UserID:   arg.UserID,
				StockID:  arg.StockID,
			})
			if err != nil {
				return err
			}
			result.User = updatedUser
			result.UserStock = updatedUserStock
			return nil
		})
	return result, err
}

//selling stock

//create user

//multiple stock information updates

//multiple settings at the same time (that shouldn't need this ngl)

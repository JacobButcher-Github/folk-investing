package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuyStockTx(t *testing.T) {
	store := NewStore(testDB)

	randUser := createRandomUser(t)
	randStock, randStockData := createRandomStock(t)
	//TODO: Made associated UserStock between randUser and randStock.

	userStartMoney := randUser.Dollars.Int64*100 + randUser.Cents.Int64

	stockCost := randStockData.ValueDollars*100 + randStockData.ValueCents

	//run n concurrent stock buy transactions
	n := 5
	amount := int64(1)
	errs := make(chan error)
	results := make(chan BuyStockTxResult)

	for range n {
		go func() {
			result, err := store.BuyStockTx(context.Background(), BuyStockTxParams{
				UserID:    randUser.ID,
				StockName: randStock.Name,
				Amount:    amount,
			})

			errs <- err
			results <- result
		}()
	}

	for range n {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		resUser := result.User
		require.NotEmpty(t, resUser)
		require.Equal(t, resUser.ID, randUser.ID)

		resStock := result.Stock
		require.NotEmpty(t, resStock)
		require.Equal(t, resStock.Name, randStock.Name)

		resUserStock := result.UserStock
		require.NotEmpty(t, resUserStock)
		require.Equal(t, resUserStock.UserID, resUser.ID)
		require.Equal(t, resUserStock.StockID, resStock.ID)
		require.True(t, int64(1) <= resUserStock.Quantity && resUserStock.Quantity <= int64(n))
	}

	updatedUser, err := testQueries.GetUser(context.Background(), randUser.UserLogin)
	require.NoError(t, err)

	updatedUserStock, err := testQueries.GetUserStock(context.Background(), GetUserStockParams{
		UserID:  randUser.ID,
		StockID: randStock.ID,
	})
	require.NoError(t, err)

	userUpdatedMoney := updatedUser.Dollars.Int64*100 + updatedUser.Cents.Int64
	require.Equal(t, userStartMoney-int64(n)*amount, userUpdatedMoney)
	//TODO: compare original UserStock amount and updated UserStock amount
}

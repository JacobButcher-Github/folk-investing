package db

import (
	//stl
	"context"
	"testing"

	//go package
	"github.com/stretchr/testify/require"

	//local
	"home/osarukun/repos/tower-investing/backend/util"
)

func TestBuyStockTx(t *testing.T) {
	store := NewStore(testDB)

	randUser := createRandomUser(t)
	randStock, randStockData := createRandomStock(t)
	userStock, err := testQueries.CreateUserStock(context.Background(), CreateUserStockParams{
		UserID:   randUser.ID,
		StockID:  randStock.ID,
		Quantity: util.RandomInt(0, 100),
	})

	userStartMoney := randUser.Dollars*100 + randUser.Cents
	stockCost := randStockData.ValueDollars*100 + randStockData.ValueCents

	//run n concurrent stock buy transactions of random amount
	n := 5

	cap := int64(userStartMoney / (int64(n) * stockCost))

	amount := util.RandomInt(1, cap)

	errs := make(chan error)
	results := make(chan BuyStockTxResult)

	for range n {
		go func() {
			result, err := store.BuyStockTx(context.Background(), BuyStockTxParams{
				UserID:  randUser.ID,
				StockID: randStock.ID,
				Amount:  amount,
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

		resUserStock := result.UserStock
		require.NotEmpty(t, resUserStock)
		require.Equal(t, resUserStock.UserID, resUser.ID)
		require.True(t, int64(1) <= resUserStock.Quantity && resUserStock.Quantity <= int64(amount*int64(n)))
	}

	updatedUser, err := testQueries.GetUserFromName(context.Background(), randUser.UserLogin)
	require.NoError(t, err)

	updatedUserStock, err := testQueries.GetUserStock(context.Background(), GetUserStockParams{
		UserID:  randUser.ID,
		StockID: randStock.ID,
	})
	require.NoError(t, err)

	userUpdatedMoney := updatedUser.Dollars*100 + updatedUser.Cents
	require.Equal(t, userStartMoney-int64(n)*amount*stockCost, userUpdatedMoney)
	require.Equal(t, updatedUserStock.Quantity-amount*int64(n), userStock.Quantity)
}

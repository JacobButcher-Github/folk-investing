package db

import (
	//stl
	"context"
	"testing"

	//go package
	"github.com/stretchr/testify/require"

	//local
	"github.com/JacobButcher-Github/folk-investing/backend/util"
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

	userStartMoney := randUser.Dollars.Int64*100 + randUser.Cents.Int64
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
		require.True(t, int64(1) <= resUserStock.Quantity && resUserStock.Quantity <= int64(amount*n))
	}

	updatedUser, err := testQueries.GetUserFromName(context.Background(), randUser.UserLogin)
	require.NoError(t, err)

	updatedUserStock, err := testQueries.GetUserStock(context.Background(), GetUserStockParams{
		UserID:  randUser.ID,
		StockID: randStock.ID,
	})
	require.NoError(t, err)

	userUpdatedMoney := updatedUser.Dollars.Int64*100 + updatedUser.Cents.Int64
	require.Equal(t, userStartMoney-int64(n)*amount*stockCost, userUpdatedMoney)
	require.Equal(t, updatedUserStock.Quantity-amount*int64(n), userStock.Quantity)
}

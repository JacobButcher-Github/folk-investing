package db

import (
	//stl
	"context"
	"database/sql"
	"testing"

	//go package
	"github.com/stretchr/testify/require"

	//local
	"github.com/JacobButcher-Github/folk-investing/backend/util"
)

func TestSellStockTx(t *testing.T) {
	store := NewStore(testDB)
	userStartQuantity := util.RandomInt(5, 1000)
	randUser := createRandomUser(t)
	randStock, randStockData := createRandomStock(t)
	userStock, err := testQueries.CreateUserStock(context.Background(), CreateUserStockParams{
		UserID:   randUser.ID,
		StockID:  randStock.ID,
		Quantity: userStartQuantity,
	})

	userStartMoney := randUser.Dollars*100 + randUser.Cents
	stockCost := randStockData.ValueDollars*100 + randStockData.ValueCents

	n := 5
	cap := int64(userStock.Quantity / int64(n))
	amount := util.RandomInt(1, cap)
	errs := make(chan error)
	results := make(chan SellStockTxResult)

	for range n {
		go func() {
			result, err := store.SellStockTx(context.Background(), SellStockTxParams{
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

		require.True(t, int64(0) <= resUserStock.Quantity && resUserStock.Quantity <= int64(amount*int64(n)+userStartQuantity))
	}

	updatedUser, err := testQueries.GetUserFromId(context.Background(), randUser.ID)
	require.NoError(t, err)

	updatedUserStock, err := testQueries.GetUserStock(context.Background(), GetUserStockParams{
		UserID:  randUser.ID,
		StockID: randStock.ID,
	})
	if int64(n)*amount == userStock.Quantity {
		require.Equal(t, err, sql.ErrNoRows)
	} else {
		require.NoError(t, err)
		require.Equal(t, updatedUserStock.Quantity+amount*int64(n), userStock.Quantity)
	}

	userUpdatedMoney := updatedUser.Dollars*100 + updatedUser.Cents
	require.Equal(t, userStartMoney+int64(n)*amount*stockCost, userUpdatedMoney)
}

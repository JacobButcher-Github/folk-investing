package db

import (
	"context"
	"testing"

	"github.com/JacobButcher-Github/folk-investing/backend/util"
	"github.com/stretchr/testify/require"
)

func createRandomUserStock(t *testing.T) UserStock {
	user := createRandomUser(t)
	stock, _ := createRandomStock(t)
	randQuantity := util.RandomInt(1, 99)

	arg := CreateUserStockParams{
		UserID:   user.ID,
		StockID:  stock.ID,
		Quantity: randQuantity,
	}

	userStock, err := testQueries.CreateUserStock(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, userStock.Quantity, randQuantity)

	return userStock
}

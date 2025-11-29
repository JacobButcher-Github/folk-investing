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

func TestCreateUserStock(t *testing.T) {
	createRandomUserStock(t)
}

func TestGetUserStock(t *testing.T) {
	userStock1 := createRandomUserStock(t)
	userStock2, err := testQueries.GetUserStock(context.Background(), GetUserStockParams{
		UserID:  userStock1.UserID,
		StockID: userStock1.StockID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, userStock2)

	require.Equal(t, userStock1.UserID, userStock2.UserID)
	require.Equal(t, userStock1.StockID, userStock2.StockID)
	require.Equal(t, userStock1.Quantity, userStock2.Quantity)
}

func TestUpdateUserStock(t *testing.T) {
	userStock := createRandomUserStock(t)
	var newQuantity int64
	for {
		newQuantity = util.randomInt(1, 99)
		if newQuantity != userStock.Quantity {
			break
		}
	}

	updatedUserStock, err := testQueries.UpdateUserStock(context.Background(), UpdateUserStockParams{
		Quantity: newQuantity,
		UserID:   userStock.UserID,
		StockID:  userStock.StockID,
	})

	require.NoError(t, err)
	require.NotEqual(t, userStock.Quantity, updatedUserStock.Quantity)
	require.Equal(t, updatedUserStock.Quantity, newQuantity)
	require.Equal(t, userStock.UserID, updatedUserStock.UserID)
	require.Equal(t, userStock.StockID, updatedUserStock.StockID)
}

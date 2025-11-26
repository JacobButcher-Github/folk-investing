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

func createRandomStock(t *testing.T) Stock {

	stockArg := CreateStockParams{
		Name:      util.RandomString(16),
		ImagePath: sql.NullString{String: util.RandomString(8), Valid: true},
	}

	stock, err := testQueries.CreateStock(context.Background(), stockArg)
	require.NoError(t, err)
	require.NotEmpty(t, stock)

	require.NotZero(t, stock.ID)
	require.Equal(t, stockArg.Name, stock.Name)
	require.Equal(t, stockArg.ImagePath, stock.ImagePath)

	stockDataArg := CreateStockDataParams{
		StockID:      stock.ID,
		EventLabel:   util.RandomString(5),
		ValueDollars: util.RandomInt(0, 100),
		ValueCents:   util.RandomInt(0, 99),
	}

	stockData, err := testQueries.CreateStockData(context.Background(), stockDataArg)
	require.NoError(t, err)
	require.NotEmpty(t, stockData)

	require.NotZero(t, stockData.ID)
	require.Equal(t, stockDataArg.StockID, stockData.StockID)
	require.Equal(t, stockDataArg.EventLabel, stockData.EventLabel)
	require.Equal(t, stockDataArg.ValueDollars, stockData.ValueDollars)
	require.Equal(t, stockDataArg.ValueDollars, stockData.ValueCents)

	return stock
}

func TestCreateStock(t *testing.T) {
	createRandomStock(t)
}

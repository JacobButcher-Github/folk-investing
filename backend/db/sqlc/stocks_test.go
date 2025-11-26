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

func createRandomStock(t *testing.T) (Stock, StockDatum) {

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

	return stock, stockData
}

func TestCreateStock(t *testing.T) {
	createRandomStock(t)
}

func TestGetStock(t *testing.T) {
	stock1, stockData1 := createRandomStock(t)
	stock2, err := testQueries.GetStock(context.Background(), stock1.Name)
	require.NoError(t, err)
	require.NotEmpty(t, stock2)
	stockData2, err := testQueries.GetStockData(context.Background(), GetStockDataParams{StockID: stock1.ID, Limit: 1})
	require.NoError(t, err)
	require.NotEmpty(t, stockData2)

	require.Equal(t, stock1.ID, stock2.ID)
	require.Equal(t, stock1.Name, stock2.Name)
	require.Equal(t, stock1.ImagePath, stock2.ImagePath)

	require.Equal(t, stockData1.ID, stockData2.ID)
	require.Equal(t, stockData1.StockID, stockData2.StockID)
	require.Equal(t, stockData1.EventLabel, stockData2.EventLabel)
	require.Equal(t, stockData1.ValueDollars, stockData2.ValueDollars)
	require.Equal(t, stockData1.ValueCents, stockData2.ValueCents)
}

func TestUpdateStockName(t *testing.T) {

}

func TestUpdateStockImagePath(t *testing.T) {

}

func TestUpdateStockAllFields(t *testing.T) {

}

func TestUpdateStockDataID(t *testing.T) {

}

func TestUpdateStockDataEventLabel(t *testing.T) {

}

func TestUpdateStockDataMoney(t *testing.T) {

}

func TestUpdateStockDataAllFields(t *testing.T) {

}

func TestDeleteStock(t *testing.T) {

}

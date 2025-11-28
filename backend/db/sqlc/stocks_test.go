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
	oldStock, _ := createRandomStock(t)
	var newName string
	for {
		newName = util.RandomString(16)
		if oldStock.Name != newName {
			break
		}
	}

	updatedStock, err := testQueries.UpdateStock(context.Background(), UpdateStockParams{
		NewName:   sql.NullString{String: newName, Valid: true},
		ImagePath: sql.NullString{String: "", Valid: false},
		Name:      oldStock.Name,
	})

	require.NoError(t, err)
	require.NotEqual(t, updatedStock.Name, oldStock.Name)
	require.Equal(t, updatedStock.Name, newName)
	require.Equal(t, oldStock.ID, updatedStock.ID)
	require.Equal(t, oldStock.ImagePath, updatedStock.ImagePath)
}

func TestUpdateStockImagePath(t *testing.T) {
	oldStock, _ := createRandomStock(t)
	var newImagePath sql.NullString

	for {
		newImagePath = sql.NullString{String: util.RandomString(8), Valid: true}
		if oldStock.ImagePath != newImagePath {
			break
		}
	}

	updatedStock, err := testQueries.UpdateStock(context.Background(), UpdateStockParams{
		NewName:   sql.NullString{String: "", Valid: false},
		ImagePath: newImagePath,
		Name:      oldStock.Name,
	})

	require.NoError(t, err)
	require.NotEqual(t, updatedStock.ImagePath, oldStock.ImagePath)
	require.Equal(t, updatedStock.ImagePath, newImagePath.String)
	require.Equal(t, oldStock.ID, updatedStock.ID)
	require.Equal(t, oldStock.Name, updatedStock.Name)
}

func TestUpdateStockAllFields(t *testing.T) {
	oldStock, _ := createRandomStock(t)
	var newName string
	var newImagePath sql.NullString
	for {
		newName = util.RandomString(16)
		newImagePath = sql.NullString{String: util.RandomString(8), Valid: true}
		if oldStock.ImagePath != newImagePath && oldStock.Name != newName {
			break
		}
	}

	updatedStock, err := testQueries.UpdateStock(context.Background(), UpdateStockParams{
		NewName:   sql.NullString{String: newName, Valid: true},
		ImagePath: newImagePath,
		Name:      oldStock.Name,
	})

	require.NoError(t, err)
	require.NotEqual(t, oldStock.Name, updatedStock.Name)
	require.NotEqual(t, oldStock.ImagePath, updatedStock.ImagePath)
	require.Equal(t, newName, updatedStock.Name)
	require.Equal(t, newImagePath.String, updatedStock.ImagePath)
}

func TestUpdateStockDataID(t *testing.T) {
	stock, oldStockData := createRandomStock(t)
	newStock, _ := createRandomStock(t)
	newId := newStock.ID

	updatedStockData, err := testQueries.UpdateStockData(context.Background(), UpdateStockDataParams{
		NewID:        sql.NullInt64{Int64: newId, Valid: true},
		NewLabel:     sql.NullString{String: "", Valid: false},
		ValueDollars: sql.NullInt64{Int64: 0, Valid: false},
		ValueCents:   sql.NullInt64{Int64: 0, Valid: false},
		StockID:      stock.ID,
		EventLabel:   oldStockData.EventLabel,
	})

	require.NoError(t, err)
	require.NotEqual(t, updatedStockData.StockID, oldStockData.StockID)
	require.Equal(t, updatedStockData.StockID, newId)
	require.Equal(t, updatedStockData.ID, oldStockData.ID)
	require.Equal(t, oldStockData.EventLabel, updatedStockData.EventLabel)
	require.Equal(t, oldStockData.ValueDollars, updatedStockData.ValueDollars)
	require.Equal(t, oldStockData.ValueCents, updatedStockData.ValueCents)
}

func TestUpdateStockDataEventLabel(t *testing.T) {
	stock, oldStockData := createRandomStock(t)

	var newEventLabel string

	for {
		newEventLabel = util.RandomString(5)
		if oldStockData.EventLabel != newEventLabel {
			break
		}
	}

	updatedStockData, err := testQueries.UpdateStockData(context.Background(), UpdateStockDataParams{
		NewID:        sql.NullInt64{Int64: 0, Valid: false},
		NewLabel:     sql.NullString{String: newEventLabel, Valid: true},
		ValueDollars: sql.NullInt64{Int64: 0, Valid: false},
		ValueCents:   sql.NullInt64{Int64: 0, Valid: false},
		StockID:      stock.ID,
		EventLabel:   oldStockData.EventLabel,
	})

	require.NoError(t, err)
	require.NotEqual(t, updatedStockData.EventLabel, oldStockData.EventLabel)
	require.Equal(t, updatedStockData.EventLabel, newEventLabel)
	require.Equal(t, oldStockData.ID, updatedStockData.ID)
	require.Equal(t, oldStockData.StockID, updatedStockData.StockID)
	require.Equal(t, oldStockData.ValueDollars, updatedStockData.ValueDollars)
	require.Equal(t, oldStockData.ValueCents, updatedStockData.ValueCents)
}

func TestUpdateStockDataMoney(t *testing.T) {
	stock, oldStockData := createRandomStock(t)
	var newDollars int64
	var newCents int64

	for {
		newDollars = util.RandomInt(0, 100000)
		newCents = util.RandomInt(0, 99)
		if oldStockData.ValueDollars != newDollars && oldStockData.ValueCents != newCents {
			break
		}
	}

	updatedStockData, err := testQueries.UpdateStockData(context.Background(), UpdateStockDataParams{
		NewID:        sql.NullInt64{Int64: 0, Valid: false},
		NewLabel:     sql.NullString{String: "", Valid: false},
		ValueDollars: sql.NullInt64{Int64: newDollars, Valid: false},
		ValueCents:   sql.NullInt64{Int64: newCents, Valid: false},
		StockID:      stock.ID,
		EventLabel:   oldStockData.EventLabel,
	})

	require.NoError(t, err)
	require.NotEqual(t, updatedStockData.ValueDollars, oldStockData.ValueDollars)
	require.NotEqual(t, updatedStockData.ValueCents, oldStockData.ValueCents)
	require.Equal(t, updatedStockData.ValueDollars, newDollars)
	require.Equal(t, updatedStockData.ValueCents, newCents)
	require.Equal(t, oldStockData.ID, updatedStockData.ID)
	require.Equal(t, oldStockData.StockID, updatedStockData.StockID)
	require.Equal(t, oldStockData.EventLabel, updatedStockData.EventLabel)
}

func TestUpdateStockDataAllFields(t *testing.T) {
	stock, oldStockData := createRandomStock(t)
	newStock, _ := createRandomStock(t)
	newId := newStock.ID
	var newEventLabel string
	var newDollars int64
	var newCents int64

	for {
		newEventLabel = util.RandomString(5)
		newDollars = util.RandomInt(0, 100000)
		newCents = util.RandomInt(0, 99)
		if newEventLabel != oldStockData.EventLabel && newDollars != oldStockData.ValueDollars && newCents != oldStockData.ValueCents {
			break
		}
	}

	updatedStockData, err := testQueries.UpdateStockData(context.Background(), UpdateStockDataParams{
		NewID:        sql.NullInt64{Int64: newId, Valid: true},
		NewLabel:     sql.NullString{String: newEventLabel, Valid: true},
		ValueDollars: sql.NullInt64{Int64: newDollars, Valid: true},
		ValueCents:   sql.NullInt64{Int64: newCents, Valid: true},
		StockID:      stock.ID,
		EventLabel:   oldStockData.EventLabel,
	})

	require.NoError(t, err)
	require.NotEqual(t, updatedStockData.EventLabel, oldStockData.EventLabel)
	require.Equal(t, updatedStockData.EventLabel, newEventLabel)
	require.NotEqual(t, updatedStockData.ValueDollars, oldStockData.ValueDollars)
	require.NotEqual(t, updatedStockData.ValueCents, oldStockData.ValueCents)
	require.Equal(t, updatedStockData.ValueDollars, newDollars)
	require.Equal(t, updatedStockData.ValueCents, newCents)
	require.NotEqual(t, updatedStockData.StockID, oldStockData.StockID)
	require.Equal(t, updatedStockData.StockID, newId)
}

func TestDeleteStock(t *testing.T) {
	stock, _ := createRandomStock(t)
	err := testQueries.DeleteStock(context.Background(), stock.Name)
	require.NoError(t, err)
	deletedStock, err := testQueries.GetStock(context.Background(), stock.Name)
	require.NoError(t, err)
	require.Empty(t, deletedStock)
	deletedStockData, err := testQueries.GetStockData(context.Background(), GetStockDataParams{
		StockID: stock.ID,
		Limit:   1,
	})
	require.NoError(t, err)
	require.Empty(t, deletedStockData)
}

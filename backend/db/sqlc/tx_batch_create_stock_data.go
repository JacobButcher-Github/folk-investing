package db

import (
	"context"
	"time"
)

// BatchCreateStockDataParams contains input parameters of batch create stock data transaction
type BatchCreateStockDataParams struct {
	NewStockData []CreateStockDataParams
}

// BatchCreateStockDataResult is result of the batch create stock data transaction
type BatchCreateStockDataResult struct {
	NewStockData []StockDatum
}

// BatchCreateStockDataTx creates a new StockData for each stock given in ID.
func (store *Store) BatchCreateStockDataTx(ctx context.Context, arg BatchCreateStockDataParams) (BatchCreateStockDataResult, error) {
	var result BatchCreateStockDataResult
	err := store.retryableTx(ctx,
		5,
		2*time.Second,
		func(q *Queries) error {
			var err error
			var newStockData StockDatum
			for i := range arg.NewStockData {
				newStockData, err = q.CreateStockData(ctx, CreateStockDataParams{
					StockID:      arg.NewStockData[i].StockID,
					EventLabel:   arg.NewStockData[i].EventLabel,
					ValueDollars: arg.NewStockData[i].ValueDollars,
					ValueCents:   arg.NewStockData[i].ValueCents,
				})
				if err != nil {
					return err
				}
				result.NewStockData = append(result.NewStockData, newStockData)
			}
			return nil
		})
	return result, err
}

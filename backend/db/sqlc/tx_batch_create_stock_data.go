package db

import (
	"context"
	"database/sql"
)

// BatchCreateStockDataParams contains input parameters of batch create stock data transaction
type BatchCreateStockDataParams struct {
	newStockData []CreateStockDataParams
}

// BatchCreateStockDataResult is result of the batch create stock data transaction
type BatchCreateStockDataResult struct {
	newStockData []StockDatum
}

// BatchCreateStockDataTx a new StockData for each stock given in ID.
func (store *Store) BatchCreateStockDataTx(ctx context.Context, arg BatchCreateStockDataParams) (BatchCreateStockDataResult, error) {
	var result BatchCreateStockDataResult
	err := store.execTx(ctx,
		&sql.TxOptions{
			Isolation: sql.LevelSerializable,
		},
		func(q *Queries) error {
			var err error
			var newStockData StockDatum
			for i := range arg.newStockData {
				newStockData, err = q.CreateStockData(ctx, CreateStockDataParams{
					StockID:      arg.newStockData[i].StockID,
					EventLabel:   arg.newStockData[i].EventLabel,
					ValueDollars: arg.newStockData[i].ValueDollars,
					ValueCents:   arg.newStockData[i].ValueCents,
				})
				if err != nil {
					return err
				}
				result.newStockData = append(result.newStockData, newStockData)
			}
			return nil
		})
	return result, err
}

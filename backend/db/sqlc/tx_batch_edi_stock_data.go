package db

import (
	"context"
	"time"
)

type BatchUpdateStockDataParams struct {
	UpdateStockData []UpdateStockDataParams
}

type BatchUpdateStockDataResult struct {
	UpdatedStockData []StockDatum
}

// BatchUpdateStockData updates the StockData for each given StockID + EventLabel combo.
func (store *Store) BatchUpdateStockData(ctx context.Context, arg BatchUpdateStockDataParams) (BatchUpdateStockDataResult, error) {
	var result BatchUpdateStockDataResult
	err := store.retryableTx(ctx, 5, 2*time.Second, func(q *Queries) error {
		var err error
		var updatedStockData StockDatum
		for i := range arg.UpdateStockData {
			updatedStockData, err = q.UpdateStockData(ctx, UpdateStockDataParams{
				NewID:        arg.UpdateStockData[i].NewID,
				NewLabel:     arg.UpdateStockData[i].NewLabel,
				ValueDollars: arg.UpdateStockData[i].ValueDollars,
				ValueCents:   arg.UpdateStockData[i].ValueCents,
				StockID:      arg.UpdateStockData[i].StockID,
				EventLabel:   arg.UpdateStockData[i].EventLabel,
			})
			if err != nil {
				return err
			}
			result.UpdatedStockData = append(result.UpdatedStockData, updatedStockData)
		}
		return nil
	})
	return result, err
}

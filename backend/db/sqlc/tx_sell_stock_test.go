package db

import (
	"context"
	"testing"

	"github.com/JacobButcher-Github/folk-investing/backend/util"
)

func TestSellStockTx(t *testing.T) {
	store := NewStore(testDB)

	randUser := createRandomUser(t)
	randStock, _ := createRandomStock(t)
	userStock, err := testQueries.CreateUserStock(context.Background(), CreateUserStockParams{
		UserID:   randUser.ID,
		StockID:  randStock.ID,
		Quantity: util.RandomInt(5, 100),
	})

	userStartMoney := randUser.Dollars.Int64*100 + randUser.Cents.Int64

	n := 5
	amount := util.RandomInt(1, 100)
	errs := make(chan error)
	results := make(chan BuyStockTxResult)
}

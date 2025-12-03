package api

import (
	//stl
	"database/sql"
	"net/http"

	//go package
	"github.com/gin-gonic/gin"

	//local
	"home/osarukun/repos/tower-investing/backend/token"

	db "github.com/JacobButcher-Github/folk-investing/backend/db/sqlc"
)

type transactionRequest struct {
	UserID  int64 `json:"user_id" binding:"required,min=1"`
	StockID int64 `json:"stock_id" binding:"required,min=1"`
	Amount  int64 `json:"amount" binding:"required,min=1"`
}

func (server *Server) buyTransaction(ctx *gin.Context) {
	var req transactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if !server.validLockout(ctx) {
		return
	}
	if !server.validUser(ctx, req.UserID) {
		return
	}
	if !server.validStock(ctx, req.StockID) {
		return
	}
	if !server.validUserStockBuy(ctx, req.UserID, req.StockID, req.Amount) {
		return
	}

	arg := db.BuyStockTxParams{
		UserID:  authPayload.UserID,
		StockID: req.StockID,
		Amount:  req.Amount,
	}

	buyStockResult, err := server.store.BuyStockTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, buyStockResult)
}

func (server *Server) sellTransaction(ctx *gin.Context) {
	var req transactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if !server.validLockout(ctx) {
		return
	}
	if !server.validUser(ctx, req.UserID) {
		return
	}
	if !server.validStock(ctx, req.StockID) {
		return
	}
	if !server.validUserStockSell(ctx, req.UserID, req.StockID, req.Amount) {
		return
	}

	arg := db.SellStockTxParams{
		UserID:  authPayload.UserID,
		StockID: req.StockID,
		Amount:  req.Amount,
	}

	sellStockResult, err := server.store.SellStockTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, sellStockResult)
}

func (server *Server) validLockout(ctx *gin.Context) bool {
	lockout, err := server.store.GetLockoutStatus(ctx)
	if err != nil {
		return false
	}
	if lockout == 1 {
		return false
	}
	return true
}

func (server *Server) validUser(ctx *gin.Context, userID int64) bool {
	_, err := server.store.GetUserFromId(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	return true
}

func (server *Server) validStock(ctx *gin.Context, stockID int64) bool {
	_, err := server.store.GetStockFromId(ctx, stockID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	return true
}

// validUserStockBuy validates if the user has enough to buy
func (server *Server) validUserStockBuy(ctx *gin.Context, userID int64, stockID int64, amount int64) bool {
	stockData, err := server.store.GetStockData(ctx, db.GetStockDataParams{
		StockID: stockID,
		Limit:   1,
	})
	if err != nil {
		return false
	}
	user, err := server.store.GetUserFromId(ctx, userID)
	if err != nil {
		return false
	}

	totalStockCost := (stockData.ValueDollars*100 + stockData.ValueCents) * amount
	totalUserCurrency := user.Dollars*100 + user.Cents

	if totalUserCurrency-totalStockCost < 0 {
		return false
	}
	return true
}

// validUserStockSell validates if the UserStock exists with the correct amount for selling.
func (server *Server) validUserStockSell(ctx *gin.Context, userID int64, stockID int64, amount int64) bool {
	//check if there's enough stock amount in the UserStock they've chosen
	userStock, err := server.store.GetUserStock(ctx, db.GetUserStockParams{
		UserID:  userID,
		StockID: stockID,
	})
	if err != nil {
		return false
	}
	if userStock.amount < amount {
		return false
	}
	return true
}

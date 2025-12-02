package api

import (
	//stl
	"net/http"

	//go package
	"github.com/gin-gonic/gin"

	//local
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

	arg := db.BuyStockTxParams{
		UserID:  req.UserID,
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

	arg := db.SellStockTxParams{
		UserID:  req.UserID,
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

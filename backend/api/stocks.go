package api

import (
	"database/sql"
	db "home/osarukun/repos/tower-investing/backend/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createStockRequest struct {
	StockName string `json:"stock_name" binding:"required"`
	ImagePath string `json:"image_path"`
}

type stockResponse struct {
	StockName string `json:"stock_name"`
	ImagePath string `json:"image_path"`
}

func newStockResponse(stock db.Stock) stockResponse {
	return stockResponse{
		StockName: stock.Name,
		ImagePath: stock.ImagePath.String,
	}
}

func (server *Server) createStock(ctx *gin.Context) {
	var req createStockRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if req.ImagePath == "" {
		req.ImagePath = "default/image/path"
	}

	arg := db.CreateStockParams{
		Name:      req.StockName,
		ImagePath: sql.NullString{String: req.ImagePath, Valid: true},
	}

	stock, err := server.store.CreateStock(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	rsp := newStockResponse(stock)
	ctx.JSON(http.StatusOK, rsp)
}

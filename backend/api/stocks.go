package api

import (
	"errors"
	"fmt"
	db "home/osarukun/repos/tower-investing/backend/db/sqlc"
	"home/osarukun/repos/tower-investing/backend/token"
	"home/osarukun/repos/tower-investing/backend/util"
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
		ImagePath: stock.ImagePath,
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

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Role != util.AdminRole {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("not an admin account")))
		return
	}

	arg := db.CreateStockParams{
		Name:      req.StockName,
		ImagePath: req.ImagePath,
	}

	stock, err := server.store.CreateStock(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	rsp := newStockResponse(stock)
	ctx.JSON(http.StatusOK, rsp)
}

type createStockDataRequest struct {
	StockIDs     []int64 `json:"stock_ids" binding:"required"`
	EventLabel   string  `json:"event_label" binding:"required"`
	ValueDollars []int64 `json:"value_dollars" binding:"required"`
	ValueCents   []int64 `json:"value_cents" binding:"required"`
}

type createStockDataResponse struct {
	StockData []db.StockDatum `json:"stock_data"`
}

func (server *Server) newStockData(ctx *gin.Context) {
	var req createStockDataRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if len(req.StockIDs) != len(req.ValueDollars) || len(req.ValueDollars) != len(req.ValueCents) {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("length of arrays not equal")))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Role != util.AdminRole {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("not an admin account")))
		return
	}

	var newBatchStockData []db.CreateStockDataParams
	for i := range req.StockIDs {
		newBatchStockData = append(newBatchStockData, db.CreateStockDataParams{
			StockID:      req.StockIDs[i],
			EventLabel:   req.EventLabel,
			ValueDollars: req.ValueDollars[i],
			ValueCents:   req.ValueCents[i],
		})
	}

	arg := db.BatchCreateStockDataParams{
		NewStockData: newBatchStockData,
	}

	result, err := server.store.BatchCreateStockDataTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var rsp createStockDataResponse
	rsp.StockData = result.NewStockData
	ctx.JSON(http.StatusOK, rsp)
}

type getStocksDataResponse struct {
	StockData []db.StockDatum `json:"stock_data"`
}

func (server *Server) stocksData(ctx *gin.Context) {
	stocks, err := server.store.GetAllStocks(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var stockIDs []int64
	for i := range stocks {
		stockIDs = append(stockIDs, stocks[i].ID)
	}

	limit, err := server.store.GetNumberEvents(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.GetStocksDataParams{
		StockIds: stockIDs,
		Limit:    limit,
	}

	res, err := server.store.GetStocksData(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var rsp getStocksDataResponse
	rsp.StockData = res
	ctx.JSON(http.StatusOK, rsp)
}

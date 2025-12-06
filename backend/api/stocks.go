package api

import (
	"errors"
	"fmt"
	db "home/osarukun/repos/tower-investing/backend/db/sqlc"
	"home/osarukun/repos/tower-investing/backend/token"
	"home/osarukun/repos/tower-investing/backend/util"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type createStockRequest struct {
	Image       *multipart.FileHeader `form:"image"`
	information struct {
		StockName string `json:"stock_name" binding:"required"`
		ImagePath string `json:"image_path"`
	} `form:"information" binding:"required"`
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
	if err := ctx.ShouldBindWith(&req, binding.FormMultipart); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if req.information.ImagePath == "" {
		req.information.ImagePath = "../../frontend/public/default_img.webp"
	} else {
		var path strings.Builder
		path.WriteString("../../frontend/public/")
		path.WriteString(req.information.ImagePath)
		req.information.ImagePath = path.String()
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Role != util.AdminRole {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("not an admin account")))
		return
	}

	arg := db.CreateStockParams{
		Name:      req.information.StockName,
		ImagePath: req.information.ImagePath,
	}

	stock, err := server.store.CreateStock(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if req.Image.Filename != "" {
		err = ctx.SaveUploadedFile(req.Image, req.information.ImagePath)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
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

// newStockData creates new stock data according to the list of StockIDs and EventLabels.
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

// stocksData lists data of all stocks up to limit. Specifically for the graph on landing page.
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

type listStockDataRequest struct {
	EventLabel string `json:"event_label" binding:"required"`
}

type listStockDataResponse struct {
	StockData []db.StockDatum `json:"stock_data"`
}

// listStockData takes in an EventLabel and lists  all StockData associated with that EventLabel for all  StockIds
func (server *Server) listStockData(ctx *gin.Context) {

}

type updateStockDataRequest struct {
	EventLabel   string   `json:"event_label" binding:"required"`
	NewLabels    []string `json:"new_label" binding:"required"`
	StockIDs     []int64  `json:"stock_ids" binding:"required"`
	NewIDs       []int64  `json:"new_ids" binding:"required"`
	ValueDollars []int64  `json:"value_dollars" binding:"required"`
	ValueCents   []int64  `json:"value_cents" binding:"required"`
}

type updateStockDataResponse struct {
	UpdatedStockData []db.StockDatum `json:"stock_data"`
}

// updateStockData takes in EventLabel and a list of UpdateStockDataParams to  update those specific  stockdatas
func (server *Server) updateStockData(ctx *gin.Context) {

}

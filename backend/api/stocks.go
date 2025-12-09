package api

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"

	db "github.com/JacobButcher-Github/folk-investing/backend/db/sqlc"
	"github.com/JacobButcher-Github/folk-investing/backend/token"
	"github.com/JacobButcher-Github/folk-investing/backend/util"

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
		req.information.ImagePath = "../../frontend/public/images/default_img.webp"
	} else {
		var path strings.Builder
		path.WriteString("../../frontend/public/images/")
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

type listStocksResponse struct {
	Stocks []db.Stock `json:"stocks"`
}

func (server *Server) listStocks(ctx *gin.Context) {
	stocks, err := server.store.GetAllStocks(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := listStocksResponse{
		Stocks: stocks,
	}
	ctx.JSON(http.StatusOK, rsp)
}

type createStockData struct {
	StockID      int64 `json:"stock_id" binding:"required"`
	ValueDollars int64 `json:"value_dollars" binding:"required"`
	ValueCents   int64 `json:"value_cents" binding:"required"`
}
type createStockDataRequest struct {
	EventLabel    string            `json:"event_label" binding:"required"`
	NewStockDatas []createStockData `json:"new_stock_datas" binding:"required, dive"`
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

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Role != util.AdminRole {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("not an admin account")))
		return
	}

	var newBatchStockData []db.CreateStockDataParams
	for i := range req.NewStockDatas {
		newBatchStockData = append(newBatchStockData, db.CreateStockDataParams{
			StockID:      req.NewStockDatas[i].StockID,
			EventLabel:   req.EventLabel,
			ValueDollars: req.NewStockDatas[i].ValueDollars,
			ValueCents:   req.NewStockDatas[i].ValueCents,
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

// Ordered and mapped by stock_id: []stock_data attributed to it
type getStocksDataResponse struct {
	StockData map[int64][]db.StockDatum `json:"stock_data_by_id"`
}

// stocksData lists data of all stocks up to limit. Specifically for the graph on landing page.
func (server *Server) getStocksData(ctx *gin.Context) {
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
	for _, sd := range res {
		rsp.StockData[sd.StockID] = append(rsp.StockData[sd.StockID], sd)
	}

	ctx.JSON(http.StatusOK, rsp)
}

type listStockDataRequest struct {
	EventLabel string `json:"event_label" binding:"required"`
}

type listStockDataResponse struct {
	StockData []db.StockDatum `json:"stock_data"`
}

// listStockDataByLabel takes in an EventLabel and lists  all StockData associated with that EventLabel for all  StockIds
func (server *Server) listStockDataByLabel(ctx *gin.Context) {
	var req listStockDataRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Role != util.AdminRole {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("not an admin account")))
		return
	}

	res, err := server.store.GetStockDataByLabel(ctx, req.EventLabel)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	rsp := listStockDataResponse{
		StockData: res,
	}
	ctx.JSON(http.StatusOK, rsp)
}

type stockDataUpdate struct {
	StockID      int64   `json:"stock_id" binding:"required"`
	NewID        *int64  `json:"new_id"`
	NewLabel     *string `json:"new_label"`
	ValueDollars *int64  `json:"value_dollars"`
	ValueCents   *int64  `json:"value_cents"`
}

type updateStockDataRequest struct {
	EventLabel string            `json:"event_label" binding:"required"`
	Updates    []stockDataUpdate `json:"updates" binding:"required, dive"`
}

type updateStockDataResponse struct {
	UpdatedStockData []db.StockDatum `json:"stock_data"`
}

// updateStockDataByLabel takes in EventLabel and a list of UpdateStockDataParams to  update those specific  stockdatas
func (server *Server) updateStockDataByLabel(ctx *gin.Context) {
	var req updateStockDataRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Role != util.AdminRole {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("not an admin account")))
		return
	}

	var updateBatchStockData []db.UpdateStockDataParams
	for _, u := range req.Updates {
		updateBatchStockData = append(updateBatchStockData, db.UpdateStockDataParams{
			StockID:      u.StockID,
			EventLabel:   req.EventLabel,
			NewID:        util.NullInt64(u.NewID),
			NewLabel:     util.NullString(u.NewLabel),
			ValueDollars: util.NullInt64(u.ValueDollars),
			ValueCents:   util.NullInt64(u.ValueCents),
		})
	}

	arg := db.BatchUpdateStockDataParams{
		UpdateStockData: updateBatchStockData,
	}

	result, err := server.store.BatchUpdateStockDataTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var rsp updateStockDataResponse
	rsp.UpdatedStockData = result.UpdatedStockData
	ctx.JSON(http.StatusOK, rsp)
}

type deleteStockDataRequest struct {
	EventLabel string `json:"event_label" binding:"required"`
}

// deleteStockDataByLabel takes in an EventLabel and deletes all StockData associated with it.
func (server *Server) deleteStockDataByLabel(ctx *gin.Context) {
	var req deleteStockDataRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Role != util.AdminRole {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("not an admin account")))
		return
	}

	err := server.store.DeleteStockDataByLabel(ctx, req.EventLabel)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

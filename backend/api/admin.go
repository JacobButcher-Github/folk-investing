package api

import (
	"fmt"
	"net/http"
	"time"

	db "github.com/JacobButcher-Github/folk-investing/backend/db/sqlc"
	"github.com/JacobButcher-Github/folk-investing/backend/token"
	"github.com/JacobButcher-Github/folk-investing/backend/util"

	"github.com/gin-gonic/gin"
)

type siteSettingsUpdateRequest struct {
	NumberOfEvents *int64  `json:"number_of_events"`
	ValueSymbol    *string `json:"value_symbol"`
	EventLabel     *string `json:"event_label"`
	Title          *string `json:"title"`
	GiveEachDay    *int64  `json:"give_each_day"`
}

type siteSettingsUpdateResponse struct {
	NumberOfEvents int64  `json:"number_of_events"`
	ValueSymbol    string `json:"value_symbol"`
	EventLabel     string `json:"event_label"`
	Title          string `json:"title"`
	GiveEachDay    int64  `json:"give_each_day"`
}

func newSiteSettingsUpdateResponse(siteSetting db.SiteSetting) siteSettingsUpdateResponse {
	return siteSettingsUpdateResponse{
		NumberOfEvents: siteSetting.NumberOfEventsVisible,
		ValueSymbol:    siteSetting.ValueSymbol,
		EventLabel:     siteSetting.EventLabel,
		Title:          siteSetting.Title,
		GiveEachDay:    siteSetting.GiveEachDay,
	}
}

func (server *Server) siteSettingsUpdate(ctx *gin.Context) {
	var req siteSettingsUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Role != util.AdminRole {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("not an admin account")))
		return
	}

	arg := db.UpdateSettingsParams{
		NumberOfEventsVisible: util.NullInt64(req.NumberOfEvents),
		ValueSymbol:           util.NullString(req.ValueSymbol),
		EventLabel:            util.NullString(req.EventLabel),
		Title:                 util.NullString(req.Title),
		GiveEachDay:           util.NullInt64(req.GiveEachDay),
	}

	updatedSettings, err := server.store.UpdateSettings(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	rsp := newSiteSettingsUpdateResponse(updatedSettings)
	ctx.JSON(http.StatusOK, rsp)
}

type lockoutUpdateRequest struct {
	Lockout     *int64     `json:"lockout"`
	LockoutTime *time.Time `json:"lockout_time"`
}

type lockoutUpdateResponse struct {
	Lockout     int64     `json:"lockout"`
	LockoutTime time.Time `json:"lockout_time"`
}

func newLockoutUpdateResponse(siteSetting db.SiteSetting) lockoutUpdateResponse {
	return lockoutUpdateResponse{
		Lockout:     siteSetting.Lockout,
		LockoutTime: siteSetting.LockoutTimeStart,
	}
}

func (server *Server) adminLockout(ctx *gin.Context) {
	var req lockoutUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Role != util.AdminRole {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("not an admin account")))
		return
	}

	arg := db.UpdateSettingsParams{
		Lockout:          util.NullInt64(req.Lockout),
		LockoutTimeStart: util.NullTime(req.LockoutTime),
	}

	updatedSettings, err := server.store.UpdateSettings(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	rsp := newLockoutUpdateResponse(updatedSettings)
	ctx.JSON(http.StatusOK, rsp)
}

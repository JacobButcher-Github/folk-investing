package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *Server) getSiteSettings(ctx *gin.Context) {
	rsp, err := server.store.GetSiteSettings(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, rsp)
}

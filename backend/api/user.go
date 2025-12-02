package api

import (
	//stl
	"database/sql"
	"net/http"

	//go package
	"github.com/gin-gonic/gin"

	//local
	db "github.com/JacobButcher-Github/folk-investing/backend/db/sqlc"
	"github.com/JacobButcher-Github/folk-investing/backend/util"
)

type createUserRequest struct {
	UserLogin string `json:"user_login" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		UserLogin:      req.UserLogin,
		HashedPassword: hashedPassword,
		Dollars:        100,
		Cents:          0,
	}

	user, err := server.store.CreateUserTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

type getUserRequest struct {
	UserLogin string `uri:"user_login" binding:"required"`
}

func (server *Server) getUser(ctx *gin.Context) {
	var req getUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserFromName(ctx, req.UserLogin)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	//paying for not seperating out user from account here, but I think it's fine.
	user.HashedPassword = ""
	ctx.JSON(http.StatusOK, user)
}

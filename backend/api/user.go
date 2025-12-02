package api

import (
	//stl
	"net/http"

	//go package
	"github.com/gin-gonic/gin"

	//local
	db "github.com/JacobButcher-Github/folk-investing/backend/db/sqlc"
)

type createUserRequest struct {
	UserLogin      string `json:"user_login" binding:"required"`
	HashedPassword string `json:"hashed_password" binding:"required"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.CreateUserParams{
		UserLogin:      req.UserLogin,
		HashedPassword: req.HashedPassword,
		Dollars:        100,
		Cents:          0,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

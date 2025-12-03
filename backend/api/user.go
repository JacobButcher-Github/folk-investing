package api

import (
	//stl
	"database/sql"
	"home/osarukun/repos/tower-investing/backend/token"
	"net/http"
	"time"

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

type userResponse struct {
	UserLogin string `json:"user_login"`
	Dollars   int64  `json:"dollars"`
	Cents     int64  `json:"cents"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		UserLogin: user.UserLogin,
		Dollars:   user.Dollars,
		Cents:     user.Cents,
	}
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

	rsp := newUserResponse(user)
	ctx.JSON(http.StatusOK, rsp)
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

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	user, err := server.store.GetUserFromName(ctx, authPayload.UserLogin)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	//paying for not seperating out user from account here, but I think it's fine.
	rsp := newUserResponse(user)
	ctx.JSON(http.StatusOK, rsp)
}

type loginUserRequest struct {
	UserLogin string `json:"user_login" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

type loginUserResponse struct {
	AccessToken string       `json:"access_token"`
	User        userResponse `json:"user"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.UserLogin)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	duration := 15 * time.Minute
	accessToken, err := server.tokenMaker.CreateToken(
		user.ID,
		user.UserLogin,
		duration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := loginUserResponse{
		AccessToken: accessToken,
		User:        newUserResponse(user),
	}
	ctx.JSON(http.StatusOK, rsp)
}

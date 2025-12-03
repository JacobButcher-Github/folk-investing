package api

import (
	//stl
	"fmt"

	//go package
	"github.com/gin-gonic/gin"

	//local
	"home/osarukun/repos/tower-investing/backend/token"

	db "home/osarukun/repos/tower-investing/backend/db/sqlc"
)

type Server struct {
	store      *db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

// NewServer creates a new http server and setup routing
func NewServer(store *db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker("")
	if err != nil {
		return nil, fmt.Errorf("Cannot create token maker: %w", err)
	}
	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {

	router := gin.Default()
	//user registration and login + token renewal
	router.POST("/users/register", server.createUser)
	router.POST("users/login", server.loginUser)
	router.POST("/tokens/renew_access", server.renewAccessToken)

	//stock data and user information
	router.GET("stocks/stocks_data", server.stocksData)
	router.GET("/users/:user_login", server.getUser)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRoutes.POST("admin/stocks/new_stock", server.createStock)
	authRoutes.POST("admin/stocks/new_stock_data", server.newStockData)
	authRoutes.POST("/transaction/buy_stock", server.buyTransaction)
	authRoutes.POST("/transaction/sell_stock", server.sellTransaction)

	server.router = router
}

// Start runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

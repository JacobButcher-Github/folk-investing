package api

import (
	//stl
	"context"
	"fmt"

	//go package
	"github.com/gin-gonic/gin"

	//local
	"github.com/JacobButcher-Github/folk-investing/backend/token"
	"github.com/JacobButcher-Github/folk-investing/backend/util"

	db "github.com/JacobButcher-Github/folk-investing/backend/db/sqlc"
)

type Server struct {
	config     util.Config
	store      *db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

// NewServer creates a new http server and setup routing
func NewServer(config util.Config, store *db.Store) (*Server, error) {
	pasetoToken := util.RandomString(32)
	tokenMaker, err := token.NewPasetoMaker(pasetoToken)
	if err != nil {
		return nil, fmt.Errorf("Cannot create token maker: %w", err)
	}
	server := &Server{
		config:     config,
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

	//stock data and user information (for leaderboard later probably.)
	router.GET("stocks/get_stocks_data", server.getStocksData)
	router.GET("stocks/list_stocks", server.listStocks)
	router.GET("/users/:user_login", server.getUser)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	//admin group
	authRoutes.POST("admin/stocks/new_stock", server.createStock)
	authRoutes.POST("admin/stocks/new_stock_data", server.newStockData)
	authRoutes.POST("admin/stocks/list_stock_data_by_label", server.listStockDataByLabel)
	authRoutes.POST("admin/stocks/edit_stock_data_by_label", server.updateStockDataByLabel)
	authRoutes.POST("admin/stocks/delete_stock_data_by_label", server.deleteStockDataByLabel)
	authRoutes.POST("admin/user_update", server.adminUserUpdate)
	authRoutes.POST("admin/settings_update", server.siteSettingsUpdate)
	authRoutes.POST("admin/lockout_reset", server.adminLockout)

	//user group
	authRoutes.POST("/users/update_user", server.updateUser)

	//transaction group
	lockoutRoutes := router.Group("/").Use(
		authMiddleware(server.tokenMaker),
		server.LockoutMiddleware(),
	)
	lockoutRoutes.POST("/transaction/buy_stock", server.buyTransaction)
	lockoutRoutes.POST("/transaction/sell_stock", server.sellTransaction)

	server.router = router
}

// Start runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	hashedPassword, err := util.HashPassword(server.config.AdminPassword)
	if err != nil {
		return fmt.Errorf("unable to hash admin password: %w", err)
	}

	arg := db.CreateAdminParams{
		UserLogin:      server.config.AdminUsername,
		Role:           util.AdminRole,
		HashedPassword: hashedPassword,
		Dollars:        100,
		Cents:          0,
	}
	server.store.CreateAdmin(context.Background(), arg)
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

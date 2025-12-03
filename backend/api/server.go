package api

import (
	//stl
	"fmt"

	//go package
	"github.com/gin-gonic/gin"

	//local
	db "github.com/JacobButcher-Github/folk-investing/backend/db/sqlc"
	"github.com/JacobButcher-Github/folk-investing/backend/token"
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

	router.POST("/users", server.createUser)
	router.POST("users/login", server.loginUser)
	router.GET("/users/:user_login", server.getUser)

	router.POST("/transaction/buy_stock", server.buyTransaction)
	router.POST("/transaction/sell_stock", server.sellTransaction)

	server.router = router
}

// Start runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

package api

import (
	//go package
	"github.com/gin-gonic/gin"

	//local
	db "github.com/JacobButcher-Github/folk-investing/backend/db/sqlc"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/users", server.createUser)

	server.router = router
	return server
}

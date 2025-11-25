package main

import (
	//stl
	"log"
	"net/http"
	"os"
	"syscall"

	//import
	"github.com/gin-gonic/gin"
	//local
	// "github.com/JacobButcher-Github/folk-investing/"
)

var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		// Return JSON response
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	if err := r.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

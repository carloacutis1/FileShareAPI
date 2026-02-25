package main

import (
	"gitlab.franciscan.edu/spr-26-backend/basement/internal/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	// create a Gin router
	router := gin.Default()

	// register route: GET /shared-files 
	// will list the shared-files directory in JSON
	router.GET("/shared-files", handlers.ListSharedFiles)

	// start the HTTP server
	router.Run(":8080")
}

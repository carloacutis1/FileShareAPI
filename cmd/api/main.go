package main

import (
	"github.com/gin-gonic/gin"
	"gitlab.franciscan.edu/spr-26-backend/basement/internal/handlers"
)

func main() {
	// create a Gin router
	router := gin.Default()

	// register route: GET /shared-files
	// will list the shared-files directory in JSON
	router.GET("/shared-files", handlers.ListSharedFiles)

	// lists all files
	router.GET("/files", handlers.ListFiles)

	// start the HTTP server
	router.Run(":8080")
}

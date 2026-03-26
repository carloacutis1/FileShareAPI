package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"gitlab.franciscan.edu/spr-26-backend/basement/internal/handlers"
	"gitlab.franciscan.edu/spr-26-backend/basement/internal/routes"
	"gitlab.franciscan.edu/spr-26-backend/basement/internal/database"
)

func main() {
	// create a Gin router
	router := gin.Default()

	// initialize the database
	database.Init()

	// setup routes
	routes.SetupRoutes(router)

	// register route: GET /folders
	// will list the folders table in JSON
	router.GET("/folders", handlers.GetFolders)

	// start the HTTP server
	log.Fatal(router.Run(":8080"))
}

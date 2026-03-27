package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"gitlab.franciscan.edu/spr-26-backend/basement/internal/database"
	"gitlab.franciscan.edu/spr-26-backend/basement/internal/routes"
)

func main() {
	// initialize database connection
	database.Init()

	// create a Gin router
	router := gin.Default()

	// setup routes
	routes.SetupRoutes(router)

	// start the HTTP server
	log.Fatal(router.Run(":8080"))
}

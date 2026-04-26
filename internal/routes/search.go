package routes

import (
	"github.com/gin-gonic/gin"
	"gitlab.franciscan.edu/spr-26-backend/basement/internal/database"
	"gitlab.franciscan.edu/spr-26-backend/basement/internal/handlers"
)

// RegisterSearchRoutes registers search related routes
func SetupSearchRoutes(router *gin.Engine) {
	router.GET("/search", handlers.SearchHandler(database.DB))
}

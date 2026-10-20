package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/carloacutis1/FileShareAPI/internal/database"
	"github.com/carloacutis1/FileShareAPI/internal/handlers"
)

// RegisterSearchRoutes registers search related routes
func SetupSearchRoutes(router *gin.Engine) {
	router.GET("/search", handlers.SearchHandler(database.DB))
}

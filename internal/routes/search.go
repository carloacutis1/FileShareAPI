package routes

import (
	"github.com/carloacutis1/FileShareAPI/internal/database"
	"github.com/carloacutis1/FileShareAPI/internal/handlers"
	"github.com/gin-gonic/gin"
)

// RegisterSearchRoutes registers search related routes
func SetupSearchRoutes(router *gin.Engine) {
	router.GET("/search", handlers.SearchHandler(database.DB))
}

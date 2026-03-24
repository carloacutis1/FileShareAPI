package routes

// Import packages
import (
	"github.com/gin-gonic/gin"
	"gitlab.franciscan.edu/spr-26-backend/basement/internal/handlers"
)

// User model api routes
func SetupUserRoutes(router *gin.Engine) {
	router.GET("/users", handlers.ListUserModels) // Using shared_file.go as a reference - MH
}

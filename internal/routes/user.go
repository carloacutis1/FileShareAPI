package routes

// Import packages
import (
	"github.com/carloacutis1/FileShareAPI/internal/handlers"
	"github.com/gin-gonic/gin"
)

// User model api routes
func SetupUserRoutes(router *gin.Engine) {
	router.GET("/users", handlers.ListUserModels) // Using shared_file.go as a reference - MH
	router.POST("/users", handlers.CreateUserModel)
	router.PUT("/users/:id", handlers.UpdateUserModel)
	router.DELETE("/users/:id", handlers.DeleteUserModel)
}

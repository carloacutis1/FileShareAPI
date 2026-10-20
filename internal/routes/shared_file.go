package routes

// shared_file model api routes
import (
	"github.com/carloacutis1/FileShareAPI/internal/auth"
	"github.com/carloacutis1/FileShareAPI/internal/handlers"
	"github.com/gin-gonic/gin"
)

// write all routes in this function that will be called in routes.go
func SetupSharedFileRoutes(router *gin.Engine) {
	router.GET("/shared-files", handlers.ListSharedFiles)
	router.POST("/shared-files", auth.AuthMiddleware(), handlers.CreateSharedFile)
	router.PUT("/shared-files/:id", auth.AuthMiddleware(), handlers.UpdateSharedFile)
	router.DELETE("/shared-files/:id", auth.AuthMiddleware(), handlers.DeleteSharedFile)
}

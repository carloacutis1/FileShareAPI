package routes

// shared_file model api routes
import (
	"github.com/gin-gonic/gin"
	"gitlab.franciscan.edu/spr-26-backend/basement/internal/handlers"
)

// write all routes in this function that will be called in routes.go
func SetupSharedFileRoutes(router *gin.Engine) {
	router.GET("/shared-files", handlers.ListSharedFiles)
	router.POST("/shared-files", handlers.CreateSharedFile)
	router.DELETE("/shared-files/:id", handlers.DeleteSharedFile)
}

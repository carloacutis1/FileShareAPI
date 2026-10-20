package routes

// folder model api routes
import (
	"github.com/carloacutis1/FileShareAPI/internal/handlers"
	"github.com/gin-gonic/gin"
)

// SetupFolderRoutes registers folder CRUD routes
func SetupFolderRoutes(router *gin.Engine) {
	router.GET("/folders", handlers.ListFolders)
	router.GET("/folders/:id", handlers.GetFolder)
	router.POST("/folders", handlers.CreateFolder)
	router.PUT("/folders/:id", handlers.UpdateFolder)
	router.DELETE("/folders/:id", handlers.DeleteFolder)
}

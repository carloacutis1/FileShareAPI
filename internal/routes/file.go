package routes

// regular file model api routes
import (
	"github.com/gin-gonic/gin"
	"gitlab.franciscan.edu/spr-26-backend/basement/internal/handlers"
)

// write all routes in this function that will be called in routes.go
func SetupFileRoutes(router *gin.Engine) {
	router.GET("/files", handlers.ListFiles)
}

package routes
// shared_file model api routes
import (
	"github.com/gin-gonic/gin"
	"gitlab.franciscan.edu/spr-26-backend/basement/internal/handlers"
)

func SetupSharedFileRoutes(router *gin.Engine)  {
	router.GET("/shared-files", handlers.ListSharedFiles)
}

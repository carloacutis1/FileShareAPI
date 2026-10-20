package routes

// regular file model api routes
import (
	"github.com/gin-gonic/gin"
	"github.com/carloacutis1/FileShareAPI/internal/handlers"
)

// write all routes in this function that will be called in routes.go
func SetupFileRoutes(router *gin.Engine, h *handlers.FileHandler) {
	router.GET("/files", h.ListFiles)
	router.POST("/files", h.CreateFile)
	router.POST("/files/upload", h.CreateFileWithContent)
	router.PUT("/files/:id", h.UpdateFile)
	router.DELETE("/files/:id", h.DeleteFile)
	router.PUT("/files/:id/content", h.UploadFileContent)
	router.GET("/files/:id/content", h.DownloadFileContent)
}

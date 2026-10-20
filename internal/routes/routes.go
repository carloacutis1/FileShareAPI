package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/carloacutis1/FileShareAPI/internal/handlers"
	"github.com/carloacutis1/FileShareAPI/internal/storage"
)

func SetupRoutes(router *gin.Engine, blobStore storage.ObjectStore) {
	fh := handlers.NewFileHandler(blobStore)
	SetupAuthRoutes(router)
	SetupFileRoutes(router, fh)
	SetupFolderRoutes(router)
	SetupSharedFileRoutes(router)
	SetupUserRoutes(router)
	SetupSearchRoutes(router)
}

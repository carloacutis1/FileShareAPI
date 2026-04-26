package routes

import (
	"github.com/gin-gonic/gin"
	"gitlab.franciscan.edu/spr-26-backend/basement/internal/handlers"
	"gitlab.franciscan.edu/spr-26-backend/basement/internal/storage"
)

func SetupRoutes(router *gin.Engine, blobStore storage.ObjectStore) {
	fh := handlers.NewFileHandler(blobStore)
	SetupAuthRoutes(router)
	SetupFileRoutes(router, fh)
	SetupFolderRoutes(router)
	SetupSharedFileRoutes(router)
	SetupUserRoutes(router)
}

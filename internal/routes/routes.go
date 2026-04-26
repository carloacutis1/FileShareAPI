package routes

import "github.com/gin-gonic/gin"

func SetupRoutes(router *gin.Engine) {
	SetupAuthRoutes(router)
	SetupFileRoutes(router)
	SetupFolderRoutes(router)
	SetupSharedFileRoutes(router)
	SetupUserRoutes(router)
}

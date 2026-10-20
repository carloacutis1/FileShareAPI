package routes

import (
	"github.com/carloacutis1/FileShareAPI/internal/auth"
	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.Engine) {
	router.POST("/auth/signup", auth.Signup)
	router.POST("/auth/login", auth.Login)
}

package routes

import (
"github.com/gin-gonic/gin"
"github.com/carloacutis1/FileShareAPI/internal/auth"
)

func SetupAuthRoutes(router *gin.Engine) {
router.POST("/auth/signup", auth.Signup)
router.POST("/auth/login", auth.Login)
}

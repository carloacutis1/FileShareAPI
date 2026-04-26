package routes

import (
"github.com/gin-gonic/gin"
"gitlab.franciscan.edu/spr-26-backend/basement/internal/auth"
)

func SetupAuthRoutes(router *gin.Engine) {
router.POST("/auth/signup", auth.Signup)
router.POST("/auth/login", auth.Login)
}

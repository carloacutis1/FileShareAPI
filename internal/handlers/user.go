package handlers

import (
	"net/http" // Import HTTP Client & Server Functionality - MH

	"github.com/gin-gonic/gin"                                      // Import GIN - MH
	"gitlab.franciscan.edu/spr-26-backend/basement/internal/models" // Connect with models folder containing other user.go - MH
)

// ListUserModels handles GET /users
// Returns all items as JSON
func ListUserModels(c *gin.Context) {
	// Get data from model layer
	items := models.GetAllUserModels()

	// Send JSON response with HTTP 200 OK
	c.JSON(http.StatusOK, items)
}

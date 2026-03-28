package handlers

import (
	"log"
	"net/http" // Import HTTP Client & Server Functionality - MH

	"github.com/gin-gonic/gin"                                                // Import GIN - MH
	"gitlab.franciscan.edu/spr-26-backend/basement/internal/database/queries" // Connect with models folder containing other user.go - MH
)

// ListUserModels handles GET /users
// Returns all items as JSON
func ListUserModels(c *gin.Context) {
	// Get data from model layer
	items, err := queries.GetUsers()
	if err != nil {
		log.Printf("ListUserModels: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	// Send JSON response with HTTP 200 OK
	c.JSON(http.StatusOK, items)
}

package handlers

import (
	"log"
	"net/http" // Import HTTP Client & Server Functionality - MH

	"github.com/gin-gonic/gin"                                                // Import GIN - MH
	"gitlab.franciscan.edu/spr-26-backend/basement/internal/database/queries" // Connect with models folder containing other user.go - MH
	"gitlab.franciscan.edu/spr-26-backend/basement/internal/models"
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

// CreateUserModel handles POST /users
// Reads a JSON body and inserts a new user into the database
func CreateUserModel(c *gin.Context) {
	var input models.UserModel
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if err := queries.CreateUser(input); err != nil {
		log.Printf("CreateUserModel: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	c.JSON(http.StatusCreated, input)
}

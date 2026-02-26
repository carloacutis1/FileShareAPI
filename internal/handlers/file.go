package handlers

// handlers for file domain model

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.franciscan.edu/spr-26-backend/basement/internal/models"
)

// ListSharedFiles handles GET /files
// returns all shared file tuples in the database as a JSON structure
func ListFiles(c *gin.Context) {
	// get the data from the model layer
	Files := models.GetFiles()

	// send JSON response
	c.JSON(http.StatusOK, Files)
}

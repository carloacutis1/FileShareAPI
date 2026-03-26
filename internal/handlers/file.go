package handlers

// handlers for file domain model

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.franciscan.edu/spr-26-backend/basement/internal/models"
)

// ListFiles handles GET /files
// returns all file tuples from the database as a JSON structure
func ListFiles(c *gin.Context) {
	// get the data from the model layer
	files, err := models.GetFiles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// send JSON response
	c.JSON(http.StatusOK, files)
}

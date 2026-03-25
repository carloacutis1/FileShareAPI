package handlers

// handlers for shared_file domain model

import (
	"net/http"
	"gitlab.franciscan.edu/spr-26-backend/basement/internal/database/queries"
	"github.com/gin-gonic/gin"
)

// ListSharedFiles handles GET /shared-files
// returns all shared file tuples in the database as a JSON structure
func ListSharedFiles(c *gin.Context) {
	// get the data from the model layer
	sharedFiles, err := queries.GetSharedFiles()  // now handles 2 values
	if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch shared files"})
			return
	}
	// send JSON response
	c.JSON(http.StatusOK, sharedFiles)
}

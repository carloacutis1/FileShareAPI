package handlers

// handlers for shared_file domain model

import (
	"database/sql"
	"github.com/carloacutis1/FileShareAPI/internal/database/queries"
	"github.com/carloacutis1/FileShareAPI/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ListSharedFiles handles GET /shared-files
// returns all shared file tuples in the database as a JSON structure
func ListSharedFiles(c *gin.Context) {
	// get the data from the model layer
	sharedFiles, err := queries.GetSharedFiles() // now handles 2 values
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch shared files"})
		return
	}
	// send JSON response
	c.JSON(http.StatusOK, sharedFiles)
}

// CreateSharedFile handles POST /shared-files
// creates a shared_file record in the database
func CreateSharedFile(c *gin.Context) {
	var sf models.SharedFile
	if err := c.ShouldBindJSON(&sf); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := queries.CreateSharedFile(sf); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "shared file created"})
}

// UpdateSharedFile handles PUT /shared-files/:id
// updates an existing shared_file record by id
func UpdateSharedFile(c *gin.Context) {
	id := c.Param("id")

	var sf models.SharedFile
	if err := c.ShouldBindJSON(&sf); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := queries.UpdateSharedFile(id, sf); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "shared file not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "shared file updated"})
}

// DeleteSharedFile handles DELETE /shared-files/:id
// deletes a shared_file record by id
func DeleteSharedFile(c *gin.Context) {
	id := c.Param("id")
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	sf, err := queries.GetSharedFile(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "shared file not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check if user is owner or has admin permission
	if sf.OwnerID != userID.(string) && sf.Permissions != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	if err := queries.DeleteSharedFile(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "shared file deleted"})
}

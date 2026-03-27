package handlers

// handlers for file domain model

import (
	"database/sql"
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

// CreateFile handles POST /files
// creates a file record in the database
func CreateFile(c *gin.Context) {
	var file models.File
	if err := c.ShouldBindJSON(&file); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.CreateFile(file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "file created"})
}

// UpdateFile handles PUT /files/:id
// updates an existing file record by id
func UpdateFile(c *gin.Context) {
	id := c.Param("id")

	var file models.File
	if err := c.ShouldBindJSON(&file); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file.ID = id
	if err := models.UpdateFile(file); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "file updated"})
}

// DeleteFile handles DELETE /files/:id
// deletes a file record by id
func DeleteFile(c *gin.Context) {
	id := c.Param("id")

	if err := models.DeleteFile(id); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "file deleted"})
}

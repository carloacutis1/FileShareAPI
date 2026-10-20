package handlers

// HTTP handlers for folder domain model
// This layer handles ONLY HTTP request/response logic

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/carloacutis1/FileShareAPI/internal/database/queries"
	"github.com/carloacutis1/FileShareAPI/internal/models"
)

// ListFolders handles GET /folders
func ListFolders(c *gin.Context) {
	folders, err := queries.GetAllFolders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, folders)
}

// GetFolder handles GET /folders/:id
func GetFolder(c *gin.Context) {
	id := c.Param("id")
	folder, err := queries.GetFolderByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "folder not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, folder)
}

// CreateFolder handles POST /folders
func CreateFolder(c *gin.Context) {
	var folder models.Folder
	if err := c.ShouldBindJSON(&folder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if folder.ID == "" || folder.Name == "" || folder.UserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id, name, and user_id are required"})
		return
	}

	if err := queries.CreateFolder(folder); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "folder created"})
}

// UpdateFolder handles PUT /folders/:id
func UpdateFolder(c *gin.Context) {
	id := c.Param("id")

	var folder models.Folder
	if err := c.ShouldBindJSON(&folder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := queries.UpdateFolder(id, folder); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "folder not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "folder updated"})
}

// DeleteFolder handles DELETE /folders/:id
func DeleteFolder(c *gin.Context) {
	id := c.Param("id")

	if err := queries.DeleteFolder(id); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "folder not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "folder deleted"})
}

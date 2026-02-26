package handlers

import (
    "net/http"
    "gitlab.franciscan.edu/spr-26-backend/basement/internal/models"
    "github.com/gin-gonic/gin"
)

// GetFolders handles GET /folders and returns all folders as JSON.
func GetFolders(c *gin.Context) {
    // Get data from model layer
    folders := models.GetFolders()

    // Send JSON response with HTTP 200 OK
    c.JSON(http.StatusOK, folders)
}
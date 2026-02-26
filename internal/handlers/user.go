package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.franciscan.edu/spr-26-backend/basement/internal/models"
)

// ListUserModels handles GET /your-endpoint

// Returns all items as JSON

func ListUserModels(c *gin.Context) {
	// Get data from model layer
	items := models.GetAllUserModels()

	// Send JSON response with HTTP 200 OK
	c.JSON(http.StatusOK, items)
}

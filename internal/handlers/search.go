package handlers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"gitlab.franciscan.edu/spr-26-backend/basement/internal/database/queries"
	"gitlab.franciscan.edu/spr-26-backend/basement/internal/models"
	"net/http"
)

// SearchHandler returns a gin handler for the unified search endpoint
func SearchHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Query("q")
		userID := c.Query("user_id")

		if query == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Search query 'q' is required"})
			return
		}

		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
			return
		}

		results, err := queries.SearchItems(db, query, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Search failed"})
			return
		}

		if results == nil {
			results = []models.SearchResult{} // Ensure we return empty array, not null
		}

		c.JSON(http.StatusOK, results)
	}
}

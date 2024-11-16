package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/knakul853/accessmesh/internal/store"
)

// GetUsers handles the request to fetch all users
func GetUsers(store *store.MongoStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := store.GetAllUsers()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to fetch users",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"users": users,
		})
	}
}

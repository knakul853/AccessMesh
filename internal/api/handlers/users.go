package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/knakul853/accessmesh/internal/models"
	"github.com/knakul853/accessmesh/internal/store"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetUsers handles the request to fetch all users
func GetUsers(store *store.MongoStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := store.GetAllUsers()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, users)
	}
}

// UpdateUser handles the request to update a user
func UpdateUser(store *store.MongoStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		var userData struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			Role     string `json:"role"`
		}

		if err := c.ShouldBindJSON(&userData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user := &models.User{
			ID:        objID,
			Username:  userData.Username,
			Email:     userData.Email,
			Role:      userData.Role,
			UpdatedAt: time.Now(),
		}

		err = store.UpdateUser(userID, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
	}
}

// DeleteUser handles the request to delete a user
func DeleteUser(store *store.MongoStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("id")
		err := store.DeleteUser(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
	}
}

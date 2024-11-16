package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/knakul853/accessmesh/internal/store"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string            `json:"name" bson:"name"`
	Description string            `json:"description" bson:"description"`
	Permissions []string          `json:"permissions" bson:"permissions"`
}

type RoleHandler struct {
	store *store.MongoStore
}

func NewRoleHandler(store *store.MongoStore) *RoleHandler {
	return &RoleHandler{store: store}
}

// Create handles the creation of a new role
func (h *RoleHandler) Create(c *gin.Context) {
	var role Role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Insert role into database
	result, err := h.store.Roles().InsertOne(c, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create role"})
		return
	}

	role.ID = result.InsertedID.(primitive.ObjectID)
	c.JSON(http.StatusCreated, role)
}

// List returns all roles
func (h *RoleHandler) List(c *gin.Context) {
	var roles []Role
	cursor, err := h.store.Roles().Find(c, primitive.D{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch roles"})
		return
	}
	defer cursor.Close(c)

	if err = cursor.All(c, &roles); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode roles"})
		return
	}

	c.JSON(http.StatusOK, roles)
}

// Get returns a specific role by ID
func (h *RoleHandler) Get(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	var role Role
	if err := h.store.Roles().FindOne(c, primitive.M{"_id": id}).Decode(&role); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	c.JSON(http.StatusOK, role)
}

// Update modifies an existing role
func (h *RoleHandler) Update(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	var role Role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role.ID = id
	result, err := h.store.Roles().ReplaceOne(c, primitive.M{"_id": id}, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update role"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	c.JSON(http.StatusOK, role)
}

// Delete removes a role
func (h *RoleHandler) Delete(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	result, err := h.store.Roles().DeleteOne(c, primitive.M{"_id": id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete role"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role deleted successfully"})
}

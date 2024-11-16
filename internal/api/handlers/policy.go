package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/knakul853/accessmesh/internal/models"
	"github.com/knakul853/accessmesh/internal/store"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PolicyHandler struct {
	store *store.MongoStore
}

func NewPolicyHandler(store *store.MongoStore) *PolicyHandler {
	return &PolicyHandler{store: store}
}

func (h *PolicyHandler) Create(c *gin.Context) {
	var policy models.Policy
	if err := c.ShouldBindJSON(&policy); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.store.Policies().InsertOne(c.Request.Context(), policy)
	if err != nil {
		log.Printf("Error creating policy: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create policy"})
		return
	}

	policy.ID = result.InsertedID.(primitive.ObjectID)
	c.JSON(http.StatusCreated, policy)
}

func (h *PolicyHandler) List(c *gin.Context) {
	log.Println("Listing policies...")
	cur, err := h.store.Policies().Find(c.Request.Context(), nil)
	if err != nil {
		log.Printf("Error listing policies: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list policies"})
		return
	}
	defer cur.Close(c.Request.Context())
	var policies []models.Policy
	for cur.Next(c.Request.Context()) {
		var policy models.Policy
		err := cur.Decode(&policy)
		if err != nil {
			log.Printf("Error decoding policy: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode policy"})
			return
		}
		policies = append(policies, policy)
	}
	if err := cur.Err(); err != nil {
		log.Printf("Error iterating cursor: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to iterate cursor"})
		return
	}
	c.JSON(http.StatusOK, policies)
}

func (h *PolicyHandler) Get(c *gin.Context) {
	log.Println("Getting policy...")
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		log.Printf("Error parsing ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid policy ID"})
		return
	}

	var policy models.Policy
	err = h.store.Policies().FindOne(c.Request.Context(), map[string]interface{}{"_id": id}).Decode(&policy)
	if err != nil {
		log.Printf("Error getting policy: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get policy"})
		return
	}
	c.JSON(http.StatusOK, policy)
}


func (h *PolicyHandler) Update(c *gin.Context) {
	log.Println("Updating policy...")
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		log.Printf("Error parsing ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid policy ID"})
		return
	}

	var policy models.Policy
	if err := c.ShouldBindJSON(&policy); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	policy.ID = id
	_, err = h.store.Policies().UpdateOne(c.Request.Context(), map[string]interface{}{"_id": id}, map[string]interface{}{"$set": policy})
	if err != nil {
		log.Printf("Error updating policy: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update policy"})
		return
	}
	c.JSON(http.StatusOK, policy)
}

func (h *PolicyHandler) Delete(c *gin.Context) {
	log.Println("Deleting policy...")
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		log.Printf("Error parsing ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid policy ID"})
		return
	}

	_, err = h.store.Policies().DeleteOne(c.Request.Context(), map[string]interface{}{"_id": id})
	if err != nil {
		log.Printf("Error deleting policy: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete policy"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "policy deleted"})
}

package api

import (
	"log"
	"github.com/gin-gonic/gin"
	"github.com/knakul853/accessmesh/internal/api/handlers"
	"github.com/knakul853/accessmesh/internal/api/middleware"
	"github.com/knakul853/accessmesh/internal/store"
	"github.com/knakul853/accessmesh/pkg/enforcer"
)

// SetupRoutes sets up the API routes for the application.
// It takes a Gin engine, a store, and an enforcer as parameters.
func SetupRoutes(r *gin.Engine, store *store.MongoStore, enforcer *enforcer.Enforcer) {
	log.Println("Setting up API routes...") // Log that route setup is starting

	policyHandler := handlers.NewPolicyHandler(store)

	api := r.Group("/api/v1")
	api.Use(middleware.AccessControl(enforcer)) // Apply access control middleware to all API routes

	policies := api.Group("/policies") // Group for policy-related routes
	{
		policies.POST("/", policyHandler.Create)     // Create a new policy
		policies.GET("/", policyHandler.List)        // List all policies
		policies.GET("/:id", policyHandler.Get)      // Get a specific policy by ID
		policies.PUT("/:id", policyHandler.Update)   // Update a specific policy by ID
		policies.DELETE("/:id", policyHandler.Delete) // Delete a specific policy by ID
	}
	log.Println("API routes setup complete.") // Log that route setup is finished
}

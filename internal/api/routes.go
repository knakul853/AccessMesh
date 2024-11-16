package api

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/knakul853/accessmesh/internal/api/handlers"
	"github.com/knakul853/accessmesh/internal/api/middleware"
	"github.com/knakul853/accessmesh/internal/services"
	"github.com/knakul853/accessmesh/internal/store"
	"github.com/knakul853/accessmesh/pkg/enforcer"
	"golang.org/x/time/rate"
)

// SetupRoutes sets up the API routes for the application.
// It takes a Gin engine, a store, and an enforcer as parameters.
func SetupRoutes(r *gin.Engine, store *store.MongoStore, enforcer *enforcer.Enforcer) {
	log.Println("Setting up API routes...")

	// Initialize rate limiter
	rateLimiter := middleware.NewIPRateLimiter(rate.Limit(100), 100)
	r.Use(middleware.RateLimiter(rateLimiter))

	// Initialize email service
	emailService := services.NewEmailService(
		"smtp.example.com",       // Replace with your SMTP host
		587,                      // Replace with your SMTP port
		"your-username",          // Replace with your SMTP username
		"your-password",          // Replace with your SMTP password
		"noreply@yourdomain.com", // Replace with your from email
	)

	policyHandler := handlers.NewPolicyHandler(store)
	authHandler := handlers.NewAuthHandler(store, emailService)

	auth := r.Group("/api/v1/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/verify-email", authHandler.VerifyEmail)
		auth.POST("/forgot-password", authHandler.ForgotPassword)
		auth.POST("/reset-password", authHandler.ResetPassword)
	}

	api := r.Group("/api/v1")
	api.Use(middleware.SessionAuth(middleware.SessionConfig{
		JWTSecret: []byte("your-secret-key"), // Replace with your secret key
	}))
	api.Use(middleware.AccessControl(enforcer))

	policies := api.Group("/policies")
	{
		policies.POST("/", policyHandler.Create)
		policies.GET("/", policyHandler.List)
		policies.GET("/:id", policyHandler.Get)
		policies.PUT("/:id", policyHandler.Update)
		policies.DELETE("/:id", policyHandler.Delete)
	}

	log.Println("API routes setup complete.")
}

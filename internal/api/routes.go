package api

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
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

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"} // Next.js dev server
	config.AllowCredentials = true
	config.AllowHeaders = append(config.AllowHeaders, "Authorization")
	r.Use(cors.New(config))

	// Initialize rate limiter
	rateLimiter := middleware.NewIPRateLimiter(rate.Limit(100), 100)
	r.Use(middleware.RateLimiter(rateLimiter))

	// Initialize email service
	emailService := services.NewEmailService(
		os.Getenv("SMTP_HOST"),
		587, // Standard SMTP port
		os.Getenv("SMTP_USERNAME"),
		os.Getenv("SMTP_PASSWORD"),
		os.Getenv("SMTP_FROM_EMAIL"),
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

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
	smtpHost := os.Getenv("SMTP_HOST")
	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	smtpFromEmail := os.Getenv("SMTP_FROM_EMAIL")
	frontendURL := os.Getenv("FRONTEND_URL")

	log.Printf("Initializing email service with config - Host: %s, Username: %s, FromEmail: %s, FrontendURL: %s",
		smtpHost, smtpUsername, smtpFromEmail, frontendURL)

	if smtpHost == "" || smtpUsername == "" || smtpPassword == "" || smtpFromEmail == "" {
		log.Printf("WARNING: Email service configuration is incomplete. Check your environment variables.")
	}

	emailService := services.NewEmailService(
		smtpHost,
		587, // Standard SMTP port
		smtpUsername,
		smtpPassword,
		smtpFromEmail,
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
		JWTSecret: []byte(os.Getenv("JWT_SECRET")),
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

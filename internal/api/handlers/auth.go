package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/knakul853/accessmesh/internal/models"
	"github.com/knakul853/accessmesh/internal/services"
	"github.com/knakul853/accessmesh/internal/store"
	"github.com/knakul853/accessmesh/pkg/auth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthHandler struct {
	store        *store.MongoStore
	emailService *services.EmailService
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type VerifyEmailRequest struct {
	Token string `json:"token" binding:"required"`
}

func NewAuthHandler(store *store.MongoStore, emailService *services.EmailService) *AuthHandler {
	return &AuthHandler{
		store:        store,
		emailService: emailService,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Failed to bind registration request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Processing registration for user: %s, email: %s", req.Username, req.Email)

	var existingUser models.User
	err := h.store.Users().FindOne(c.Request.Context(), bson.M{"username": req.Username}).Decode(&existingUser)
	if err == nil {
		log.Printf("Username already exists: %s", req.Username)
		c.JSON(http.StatusConflict, gin.H{"error": "username already exists"})
		return
	}

	verificationToken, err := services.GenerateToken()
	if err != nil {
		log.Printf("Failed to generate verification token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate verification token"})
		return
	}
	log.Printf("Generated verification token for user: %s", req.Username)

	user := models.User{
		Username:          req.Username,
		Email:             req.Email,
		Password:         req.Password,
		Role:              req.Role,
		EmailVerified:     false,
		VerificationToken: verificationToken,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	if err := user.HashPassword(); err != nil {
		log.Printf("Failed to hash password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	result, err := h.store.Users().InsertOne(c.Request.Context(), user)
	if err != nil {
		log.Printf("Failed to create user in database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	user.ID = result.InsertedID.(primitive.ObjectID)
	log.Printf("Successfully created user in database with ID: %s", user.ID.Hex())

	if err := h.emailService.SendVerificationEmail(user.Email, verificationToken); err != nil {
		log.Printf("Failed to send verification email: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to send verification email: %v", err)})
		return
	}

	user.Password = ""
	log.Printf("Registration successful for user: %s", user.Username)
	c.JSON(http.StatusCreated, user)
}

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	var req VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"email_verified":     true,
			"verification_token": "",
			"updated_at":        time.Now(),
		},
	}

	result, err := h.store.Users().UpdateOne(
		c.Request.Context(),
		bson.M{"verification_token": req.Token},
		update,
	)

	if err != nil || result.ModifiedCount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid verification token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "email verified successfully"})
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	err := h.store.Users().FindOne(c.Request.Context(), bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "if email exists, password reset link will be sent"})
		return
	}

	token, err := services.NewToken(24 * time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate reset token"})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"reset_token":        token.Token,
			"reset_token_expiry": token.ExpiresAt,
			"updated_at":         time.Now(),
		},
	}

	_, err = h.store.Users().UpdateByID(c.Request.Context(), user.ID, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		return
	}

	if err := h.emailService.SendPasswordResetEmail(user.Email, token.Token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send reset email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password reset email sent"})
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	err := h.store.Users().FindOne(c.Request.Context(), bson.M{
		"reset_token": req.Token,
		"reset_token_expiry": bson.M{"$gt": time.Now()},
	}).Decode(&user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or expired reset token"})
		return
	}

	user.Password = req.Password
	if err := user.HashPassword(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"password":          user.Password,
			"reset_token":       "",
			"reset_token_expiry": time.Time{},
			"updated_at":        time.Now(),
		},
	}

	_, err = h.store.Users().UpdateByID(c.Request.Context(), user.ID, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password reset successfully"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	err := h.store.Users().FindOne(c.Request.Context(), bson.M{"username": req.Username}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if err := user.ComparePassword(req.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := auth.GenerateToken(user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	user.Password = "" // Don't send password back

	c.JSON(http.StatusOK, AuthResponse{
		Token: token,
		User:  user,
	})
}

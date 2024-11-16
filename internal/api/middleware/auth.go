package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/knakul853/accessmesh/pkg/auth"
	"github.com/knakul853/accessmesh/pkg/enforcer"
)

func AccessControl(e *enforcer.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		claims, err := auth.ValidateToken(token)
		if err != nil {
			log.Printf("Error validating token: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		allowed, err := e.Enforce(claims.Role, c.Request.URL.Path, c.Request.Method)
		if err != nil {
			log.Printf("Error enforcing policy: %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		if !allowed {
			log.Printf("Access denied for user %s to resource %s with method %s", claims.Role, c.Request.URL.Path, c.Request.Method)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access denied"})
			return
		}

		c.Next()
	}
}

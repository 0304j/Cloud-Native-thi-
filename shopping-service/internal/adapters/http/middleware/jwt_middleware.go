package middleware

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
)

// Context key names
const ContextUserIDKey = "user_id"
const ContextUserRoleKey = "user_role"

// Extract token from Authorization header or cookie
func extractToken(c *gin.Context) (string, bool) {
	// Try Authorization header first (Bearer token)
	auth := c.GetHeader("Authorization")
	if auth != "" {
		parts := strings.SplitN(auth, " ", 2)
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			return parts[1], true
		}
	}

	// Try httpOnly cookie as fallback
	if cookie, err := c.Cookie("jwt_token"); err == nil && cookie != "" {
		return cookie, true
	}

	return "", false
}

// JWTMiddleware returns a Gin middleware that validates JWT and sets user_id and user_role in context
func JWTMiddleware() gin.HandlerFunc {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET not set")
	}
	log.Printf("DEBUG: JWT_SECRET = '%s'", secret)
	return func(c *gin.Context) {
		tokenStr, ok := extractToken(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid authorization token"})
			return
		}
		log.Printf("DEBUG: Received token = '%s'", tokenStr[:50]+"...")
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, jwt.ErrTokenUnverifiable
			}
			return []byte(secret), nil
		}, jwt.WithLeeway(5*time.Second))

		if err != nil || !token.Valid {
			log.Printf("DEBUG: Token validation failed: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token", "details": err.Error()})
			return
		}
		log.Printf("DEBUG: Token validated successfully")
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			return
		}

		// check exp if present (jwt lib already does but double-check)
		if expVal, ok := claims["exp"]; ok {
			switch v := expVal.(type) {
			case float64:
				if int64(v) < time.Now().Unix() {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
					return
				}
			}
		}

		// extract user_id (can be string or number)
		var userID string
		if id, ok := claims["user_id"]; ok {
			switch v := id.(type) {
			case string:
				userID = v
			default:
				userID = strings.TrimSpace(toString(v))
			}
		} else if sub, ok := claims["sub"]; ok {
			userID = toString(sub)
		}

		if userID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user_id claim missing"})
			return
		}
		c.Set(ContextUserIDKey, userID)

		// extract user_role
		if role, ok := claims["role"].(string); ok && role != "" {
			c.Set(ContextUserRoleKey, role)
		}

		c.Next()
	}
}

// Role enforcement middleware
func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		v, ok := c.Get(ContextUserRoleKey)
		if !ok || v != role {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient role"})
			return
		}
		c.Next()
	}
}

// helper to convert interface to string (simple)
func toString(v interface{}) string {
	switch x := v.(type) {
	case string:
		return x
	case float32:
		return strconv.FormatFloat(float64(x), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(x, 'f', -1, 64)
	case int:
		return strconv.Itoa(x)
	case int64:
		return strconv.FormatInt(x, 10)
	default:
		return ""
	}
}

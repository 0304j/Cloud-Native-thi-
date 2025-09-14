package http

import (
	"auth-service/internal/service"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{s}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"` // "admin" oder "user"
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Register(req.Email, req.Password, req.Role); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user registered"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.Authenticate(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role, // "admin" oder "user"
		"exp":     time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// FIX: Verwende os.Getenv() statt h.service.JWTSecret
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "JWT_SECRET not configured"})
		return
	}

	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token error"})
		return
	}

	// Set httpOnly cookie instead of returning token in response
	c.SetCookie(
		"jwt_token",     // cookie name
		signedToken,     // cookie value
		3600,           // max age (1 hour, same as token expiration)
		"/",            // path
		"",             // domain (empty means current domain)
		false,          // secure (set to true in production with HTTPS)
		true,           // httpOnly
	)

	c.JSON(http.StatusOK, gin.H{"message": "login successful"})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// Clear the httpOnly cookie
	c.SetCookie(
		"jwt_token",     // cookie name
		"",             // empty value
		-1,             // negative max age to delete cookie
		"/",            // path
		"",             // domain
		false,          // secure
		true,           // httpOnly
	)

	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

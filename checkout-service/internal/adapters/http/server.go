package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"checkout-service/internal/ports"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router      *gin.Engine
	httpServer  *http.Server
	orderService ports.OrderService
}

func NewServer(orderService ports.OrderService) *Server {
	router := gin.Default()
	
	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		
		c.Next()
	})

	server := &Server{
		router:      router,
		orderService: orderService,
	}

	server.setupRoutes()
	
	return server
}

func (s *Server) setupRoutes() {
	checkoutHandler := NewCheckoutHandler(s.orderService)

	// Health check
	s.router.GET("/health", checkoutHandler.HealthCheck)
	
	// API routes - direct routes for consistency with other services
	s.router.POST("/checkout", checkoutHandler.CreateOrder)

	// API endpoints only - frontend is served by nginx-proxy
}

func (s *Server) Start(port int) error {
	s.httpServer = &http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		Handler:        s.router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.httpServer != nil {
		return s.httpServer.Shutdown(ctx)
	}
	return nil
}
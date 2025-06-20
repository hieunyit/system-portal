package middleware

import (
	"net/http"
	"strings"

	"system-portal/internal/shared/config"

	"github.com/gin-gonic/gin"
)

// CorsMiddleware handles simple CORS headers based on configuration.
type CorsMiddleware struct {
	cfg config.CORSConfig
}

// NewCorsMiddleware creates a new middleware instance.
func NewCorsMiddleware(cfg config.CORSConfig) *CorsMiddleware {
	return &CorsMiddleware{cfg: cfg}
}

// Handler returns the main CORS handler.
func (m *CorsMiddleware) Handler() gin.HandlerFunc {
	allowedOrigins := strings.Join(m.cfg.AllowedOrigins, ",")
	allowedMethods := strings.Join(m.cfg.AllowedMethods, ",")
	allowedHeaders := strings.Join(m.cfg.AllowedHeaders, ",")

	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigins)
		c.Writer.Header().Set("Access-Control-Allow-Methods", allowedMethods)
		c.Writer.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
		if m.cfg.AllowCredentials {
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// SecurityHeaders adds a few basic security headers when enabled.
func (m *CorsMiddleware) SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
		c.Next()
	}
}

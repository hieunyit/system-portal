package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"system-portal/pkg/logger"
)

// AuditMiddleware logs simple request information.
type AuditMiddleware struct{}

func NewAuditMiddleware() *AuditMiddleware { return &AuditMiddleware{} }

func (a *AuditMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		logger.Log.WithFields(map[string]interface{}{
			"method":   c.Request.Method,
			"path":     c.Request.URL.Path,
			"status":   c.Writer.Status(),
			"duration": time.Since(start).String(),
		}).Info("request handled")
	}
}

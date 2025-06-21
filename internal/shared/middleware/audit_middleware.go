package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"system-portal/internal/domains/portal/entities"
	"system-portal/internal/domains/portal/usecases"
	"system-portal/pkg/logger"
)

// AuditMiddleware logs simple request information.
type AuditMiddleware struct {
	uc usecases.AuditUsecase
}

func NewAuditMiddleware(u usecases.AuditUsecase) *AuditMiddleware { return &AuditMiddleware{uc: u} }

func (a *AuditMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		logger.Log.WithFields(map[string]interface{}{
			"method":   c.Request.Method,
			"path":     c.Request.URL.Path,
			"status":   c.Writer.Status(),
			"duration": duration.String(),
		}).Info("request handled")

		if a.uc != nil {
			logEntry := &entities.AuditLog{
				ID:        uuid.New(),
				Action:    c.Request.Method,
				Resource:  c.Request.URL.Path,
				Success:   c.Writer.Status() < 400,
				CreatedAt: time.Now(),
			}
			// Username is available via context but user ID lookup is omitted
			if uname, ok := c.Get("username"); ok {
				if name, ok := uname.(string); ok {
					logEntry.Action = name + " " + logEntry.Action
				}
			}
			if err := a.uc.Add(c.Request.Context(), logEntry); err != nil {
				logger.Log.WithError(err).Warn("failed to add audit log")
			}
		}
	}
}

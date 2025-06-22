package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"system-portal/internal/domains/portal/entities"
	portalrepos "system-portal/internal/domains/portal/repositories"
	"system-portal/internal/domains/portal/usecases"
	"system-portal/pkg/logger"
)

// AuditMiddleware logs simple request information.
type AuditMiddleware struct {
	uc     usecases.AuditUsecase
	users  portalrepos.UserRepository
	groups portalrepos.GroupRepository
}

func NewAuditMiddleware(u usecases.AuditUsecase, userRepo portalrepos.UserRepository, groupRepo portalrepos.GroupRepository) *AuditMiddleware {
	return &AuditMiddleware{uc: u, users: userRepo, groups: groupRepo}
}

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

		if a.uc != nil && (c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut || c.Request.Method == http.MethodDelete) {
			logEntry := &entities.AuditLog{
				ID:           uuid.New(),
				Action:       c.Request.Method,
				Resource:     c.Request.URL.Path,
				Success:      c.Writer.Status() < 400,
				CreatedAt:    time.Now(),
				ResourceName: c.FullPath(),
				IPAddress:    c.ClientIP(),
			}

			if uid, ok := c.Get("userID"); ok {
				if id, ok := uid.(uuid.UUID); ok {
					logEntry.UserID = id
				}
			}
			if uname, ok := c.Get("username"); ok {
				if name, ok := uname.(string); ok {
					logEntry.Username = name
					if logEntry.UserID == uuid.Nil && a.users != nil {
						if usr, err := a.users.GetByUsername(c.Request.Context(), name); err == nil && usr != nil {
							logEntry.UserID = usr.ID
							if a.groups != nil {
								if grp, err := a.groups.GetByID(c.Request.Context(), usr.GroupID); err == nil && grp != nil {
									logEntry.UserGroup = grp.Name
								}
							}
						}
					}
				}
			}
			if group, ok := c.Get("role"); ok && logEntry.UserGroup == "" {
				if g, ok := group.(string); ok {
					logEntry.UserGroup = g
				}
			}

			if err := a.uc.Add(c.Request.Context(), logEntry); err != nil {
				logger.Log.WithError(err).Warn("failed to add audit log")
			}
		}
	}
}

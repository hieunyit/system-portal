package middleware

import (
	nethttp "net/http"

	"github.com/gin-gonic/gin"
	http "system-portal/internal/shared/response"
)

// ValidationMiddleware adds basic request validations.
type ValidationMiddleware struct{}

func NewValidationMiddleware() *ValidationMiddleware { return &ValidationMiddleware{} }

// StrictJSONBinding enforces Content-Type application/json for mutating requests.
func (v *ValidationMiddleware) StrictJSONBinding() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == nethttp.MethodPost || c.Request.Method == nethttp.MethodPut || c.Request.Method == nethttp.MethodPatch {
			if c.ContentType() != "application/json" {
				http.RespondWithBadRequest(c, "invalid content type")
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

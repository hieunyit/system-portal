package middleware

import (
	"strings"

	http "system-portal/internal/shared/response"
	"system-portal/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT access tokens.
type AuthMiddleware struct {
	jwtService *jwt.RSAService
}

// NewAuthMiddleware creates a new middleware instance.
func NewAuthMiddleware(jwtService *jwt.RSAService) *AuthMiddleware {
	return &AuthMiddleware{jwtService: jwtService}
}

// RequireAuth ensures a valid Bearer token is provided.
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			http.RespondWithUnauthorized(c, "missing token")
			c.Abort()
			return
		}

		token := strings.TrimPrefix(header, "Bearer ")
		claims, err := m.jwtService.ValidateAccessToken(token)
		if err != nil {
			http.RespondWithUnauthorized(c, "invalid token")
			c.Abort()
			return
		}

		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}

package handlers

import (
	"system-portal/internal/domains/auth/dto"
	"system-portal/internal/domains/auth/usecases"
	http "system-portal/internal/shared/response"
	"system-portal/pkg/logger"

	"github.com/gin-gonic/gin"
	"strings"
)

// AuthHandler exposes authentication endpoints.
type AuthHandler struct {
	usecase usecases.AuthUsecase
}

// NewAuthHandler creates a new handler instance.
func NewAuthHandler(u usecases.AuthUsecase) *AuthHandler { return &AuthHandler{usecase: u} }

// Login godoc
// @Summary User login
// @Description Authenticate a user and issue JWT tokens
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login credentials"
// @Success 200 {object} dto.TokenResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.WithError(err).Error("failed to bind login request")
		http.RespondWithBadRequest(c, "invalid request")
		return
	}
	access, refresh, userID, role, err := h.usecase.Login(c.Request.Context(), req.Username, req.Password, c.ClientIP())
	if err != nil {
		logger.Log.WithError(err).WithField("username", req.Username).Error("login failed")
		http.RespondWithUnauthorized(c, "login failed")
		return
	}
	logger.Log.WithField("username", req.Username).Info("user logged in")
	c.Set("username", req.Username)
	c.Set("userID", userID)
	c.Set("role", role)
	c.Set("ip", c.ClientIP())
	http.RespondWithSuccess(c, 200, dto.TokenResponse{AccessToken: access, RefreshToken: refresh})
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Issue a new access token using a refresh token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.RefreshRequest true "Refresh token"
// @Success 200 {object} dto.TokenResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.WithError(err).Error("failed to bind refresh token request")
		http.RespondWithBadRequest(c, "invalid request")
		return
	}
	access, refresh, err := h.usecase.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		logger.Log.WithError(err).Error("refresh token failed")
		http.RespondWithUnauthorized(c, "refresh failed")
		return
	}
	logger.Log.Info("refresh token issued")
	http.RespondWithSuccess(c, 200, dto.TokenResponse{AccessToken: access, RefreshToken: refresh})
}

// ValidateToken godoc
// @Summary Validate token
// @Description Validate an access token
// @Tags Authentication
// @Security BearerAuth
// @Produce json
// @Success 200 {string} string "token valid"
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /auth/validate [get]
func (h *AuthHandler) ValidateToken(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		logger.Log.Warn("missing authorization header")
		http.RespondWithBadRequest(c, "missing token")
		return
	}
	if err := h.usecase.Validate(c.Request.Context(), token[7:]); err != nil {
		logger.Log.WithError(err).Warn("invalid token")
		http.RespondWithUnauthorized(c, "invalid token")
		return
	}
	logger.Log.Info("token validated")
	http.RespondWithMessage(c, 200, "token valid")
}

// Logout godoc
// @Summary Logout
// @Description Invalidate the current session
// @Tags Authentication
// @Security BearerAuth
// @Produce json
// @Success 200 {string} string "logged out"
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	token := strings.TrimSpace(strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer "))
	if err := h.usecase.Logout(c.Request.Context(), token); err != nil {
		logger.Log.WithError(err).Warn("logout failed")
	}
	logger.Log.Info("user logged out")
	http.RespondWithMessage(c, 200, "logged out")
}

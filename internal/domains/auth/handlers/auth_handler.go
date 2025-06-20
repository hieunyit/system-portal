package handlers

import (
	"system-portal/internal/domains/auth/dto"
	"system-portal/internal/domains/auth/usecases"
	http "system-portal/internal/shared/response"

	"github.com/gin-gonic/gin"
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
		http.RespondWithBadRequest(c, "invalid request")
		return
	}
	access, refresh, err := h.usecase.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		http.RespondWithUnauthorized(c, "login failed")
		return
	}
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
		http.RespondWithBadRequest(c, "invalid request")
		return
	}
	access, refresh, err := h.usecase.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		http.RespondWithUnauthorized(c, "refresh failed")
		return
	}
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
		http.RespondWithBadRequest(c, "missing token")
		return
	}
	if err := h.usecase.Validate(c.Request.Context(), token[7:]); err != nil {
		http.RespondWithUnauthorized(c, "invalid token")
		return
	}
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
	token := c.GetHeader("Authorization")
	h.usecase.Logout(c.Request.Context(), token)
	http.RespondWithMessage(c, 200, "logged out")
}

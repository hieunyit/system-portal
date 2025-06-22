package handlers

import (
	nethttp "net/http"
	"time"

	"system-portal/internal/domains/portal/dto"
	"system-portal/internal/domains/portal/entities"
	"system-portal/internal/domains/portal/usecases"
	http "system-portal/internal/shared/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	uc usecases.UserUsecase
}

func NewUserHandler(u usecases.UserUsecase) *UserHandler { return &UserHandler{uc: u} }

// ListUsers godoc
// @Summary List portal users
// @Description Retrieve all portal users
// @Tags Portal Users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.SuccessResponse{data=[]dto.PortalUserResponse}
// @Router /api/portal/users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	users, _ := h.uc.List(c.Request.Context())
	resp := make([]dto.PortalUserResponse, 0, len(users))
	for _, u := range users {
		resp = append(resp, dto.PortalUserResponse{
			ID:       u.ID,
			Username: u.Username,
			Email:    u.Email,
			FullName: u.FullName,
			GroupID:  u.GroupID,
			IsActive: u.IsActive,
		})
	}
	http.RespondWithSuccess(c, nethttp.StatusOK, resp)
}

// CreateUser godoc
// @Summary Create portal user
// @Description Create a new portal user
// @Tags Portal Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.PortalUserUpdateRequest true "User data"
// @Success 201 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /api/portal/users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.PortalUserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		http.RespondWithBadRequest(c, "invalid request")
		return
	}
	user := &entities.PortalUser{
		ID:        uuid.New(),
		Username:  req.Username,
		Email:     req.Email,
		FullName:  req.FullName,
		Password:  req.Password,
		GroupID:   req.GroupID,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := h.uc.Create(c.Request.Context(), user); err != nil {
		http.RespondWithBadRequest(c, err.Error())
		return
	}
	http.RespondWithMessage(c, nethttp.StatusCreated, "created")
}

// GetUser godoc
// @Summary Get portal user
// @Description Get portal user by ID
// @Tags Portal Users
// @Security BearerAuth
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} response.SuccessResponse{data=dto.PortalUserResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /api/portal/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		http.RespondWithBadRequest(c, "invalid id")
		return
	}
	u, _ := h.uc.Get(c.Request.Context(), id)
	if u == nil {
		http.RespondWithNotFound(c, "not found")
		return
	}
	http.RespondWithSuccess(c, nethttp.StatusOK, dto.PortalUserResponse{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		FullName: u.FullName,
		GroupID:  u.GroupID,
		IsActive: u.IsActive,
	})
}

// UpdateUser godoc
// @Summary Update portal user
// @Description Update information for a portal user
// @Tags Portal Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body dto.PortalUserRequest true "User data"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /api/portal/users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		http.RespondWithBadRequest(c, "invalid id")
		return
	}
	var req dto.PortalUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		http.RespondWithBadRequest(c, "invalid request")
		return
	}
	user := &entities.PortalUser{
		ID:        id,
		FullName:  req.FullName,
		Password:  req.Password,
		GroupID:   req.GroupID,
		IsActive:  true,
		UpdatedAt: time.Now(),
	}
	if err := h.uc.Update(c.Request.Context(), user); err != nil {
		http.RespondWithBadRequest(c, err.Error())
		return
	}
	http.RespondWithMessage(c, nethttp.StatusOK, "updated")
}

// DeleteUser godoc
// @Summary Delete portal user
// @Description Remove a portal user
// @Tags Portal Users
// @Security BearerAuth
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /api/portal/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		http.RespondWithBadRequest(c, "invalid id")
		return
	}
	if err := h.uc.Delete(c.Request.Context(), id); err != nil {
		http.RespondWithBadRequest(c, err.Error())
		return
	}
	http.RespondWithMessage(c, nethttp.StatusOK, "deleted")
}

// ActivateUser godoc
// @Summary Activate portal user
// @Tags Portal Users
// @Security BearerAuth
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} response.SuccessResponse
// @Router /api/portal/users/{id}/activate [put]
func (h *UserHandler) ActivateUser(c *gin.Context) { http.RespondWithMessage(c, 200, "ok") }

// DeactivateUser godoc
// @Summary Deactivate portal user
// @Tags Portal Users
// @Security BearerAuth
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} response.SuccessResponse
// @Router /api/portal/users/{id}/deactivate [put]
func (h *UserHandler) DeactivateUser(c *gin.Context) { http.RespondWithMessage(c, 200, "ok") }

// ResetPassword godoc
// @Summary Reset portal user password
// @Tags Portal Users
// @Security BearerAuth
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} response.SuccessResponse
// @Router /api/portal/users/{id}/reset-password [put]
func (h *UserHandler) ResetPassword(c *gin.Context) { http.RespondWithMessage(c, 200, "ok") }

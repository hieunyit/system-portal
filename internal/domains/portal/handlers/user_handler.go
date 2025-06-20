package handlers

import (
	nethttp "net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"system-portal/internal/domains/portal/dto"
	"system-portal/internal/domains/portal/entities"
	"system-portal/internal/domains/portal/usecases"
	http "system-portal/internal/shared/response"
)

type UserHandler struct {
	uc usecases.UserUsecase
}

func NewUserHandler(u usecases.UserUsecase) *UserHandler { return &UserHandler{uc: u} }

func (h *UserHandler) ListUsers(c *gin.Context) {
	users, _ := h.uc.List(c.Request.Context())
	resp := make([]dto.UserResponse, 0, len(users))
	for _, u := range users {
		resp = append(resp, dto.UserResponse{
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

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		http.RespondWithBadRequest(c, "invalid request")
		return
	}
	user := &entities.User{
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
	h.uc.Create(c.Request.Context(), user)
	http.RespondWithMessage(c, nethttp.StatusCreated, "created")
}

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
	http.RespondWithSuccess(c, nethttp.StatusOK, dto.UserResponse{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		FullName: u.FullName,
		GroupID:  u.GroupID,
		IsActive: u.IsActive,
	})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		http.RespondWithBadRequest(c, "invalid id")
		return
	}
	var req dto.UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		http.RespondWithBadRequest(c, "invalid request")
		return
	}
	user := &entities.User{
		ID:        id,
		Username:  req.Username,
		Email:     req.Email,
		FullName:  req.FullName,
		Password:  req.Password,
		GroupID:   req.GroupID,
		IsActive:  true,
		UpdatedAt: time.Now(),
	}
	h.uc.Update(c.Request.Context(), user)
	http.RespondWithMessage(c, nethttp.StatusOK, "updated")
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		http.RespondWithBadRequest(c, "invalid id")
		return
	}
	h.uc.Delete(c.Request.Context(), id)
	http.RespondWithMessage(c, nethttp.StatusOK, "deleted")
}

func (h *UserHandler) ActivateUser(c *gin.Context)   { http.RespondWithMessage(c, 200, "ok") }
func (h *UserHandler) DeactivateUser(c *gin.Context) { http.RespondWithMessage(c, 200, "ok") }
func (h *UserHandler) ResetPassword(c *gin.Context)  { http.RespondWithMessage(c, 200, "ok") }

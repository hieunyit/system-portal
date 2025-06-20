package handlers

import (
	nethttp "net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"system-portal/internal/domains/portal/entities"
	"system-portal/internal/domains/portal/usecases"
	http "system-portal/internal/shared/response"
)

type GroupHandler struct{ uc usecases.GroupUsecase }

func NewGroupHandler(u usecases.GroupUsecase) *GroupHandler { return &GroupHandler{uc: u} }

func (h *GroupHandler) ListGroups(c *gin.Context) {
	groups, _ := h.uc.List(c.Request.Context())
	http.RespondWithSuccess(c, nethttp.StatusOK, groups)
}

func (h *GroupHandler) GetGroup(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		http.RespondWithBadRequest(c, "invalid id")
		return
	}
	g, _ := h.uc.Get(c.Request.Context(), id)
	if g == nil {
		http.RespondWithNotFound(c, "not found")
		return
	}
	http.RespondWithSuccess(c, nethttp.StatusOK, g)
}

func (h *GroupHandler) CreateGroup(c *gin.Context) {
	var g entities.Group
	if err := c.ShouldBindJSON(&g); err != nil {
		http.RespondWithBadRequest(c, "invalid request")
		return
	}
	g.ID = uuid.New()
	g.CreatedAt = time.Now()
	g.UpdatedAt = time.Now()
	h.uc.Create(c.Request.Context(), &g)
	http.RespondWithMessage(c, nethttp.StatusCreated, "created")
}

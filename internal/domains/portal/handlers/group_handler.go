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

// ListGroups godoc
// @Summary List portal groups
// @Tags Portal Groups
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.SuccessResponse{data=[]entities.PortalGroup}
// @Router /api/portal/groups [get]
type groupQuery struct {
	Name  string `form:"name"`
	Page  int    `form:"page,default=1"`
	Limit int    `form:"limit,default=20"`
}

func (h *GroupHandler) ListGroups(c *gin.Context) {
	var q groupQuery
	_ = c.ShouldBindQuery(&q)
	filter := &entities.GroupFilter{Name: q.Name, Page: q.Page, Limit: q.Limit}
	groups, total, _ := h.uc.List(c.Request.Context(), filter)
	http.RespondWithSuccess(c, nethttp.StatusOK, gin.H{"groups": groups, "total": total, "page": filter.Page, "limit": filter.Limit})
}

// GetGroup godoc
// @Summary Get portal group
// @Tags Portal Groups
// @Security BearerAuth
// @Produce json
// @Param id path string true "Group ID"
// @Success 200 {object} response.SuccessResponse{data=entities.PortalGroup}
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /api/portal/groups/{id} [get]
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

// CreateGroup godoc
// @Summary Create portal group
// @Tags Portal Groups
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body entities.PortalGroup true "Group data"
// @Success 201 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /api/portal/groups [post]
func (h *GroupHandler) CreateGroup(c *gin.Context) {
	var g entities.PortalGroup
	if err := c.ShouldBindJSON(&g); err != nil {
		http.RespondWithBadRequest(c, "invalid request")
		return
	}
	g.ID = uuid.New()
	g.CreatedAt = time.Now()
	g.UpdatedAt = time.Now()
	if err := h.uc.Create(c.Request.Context(), &g); err != nil {
		http.RespondWithBadRequest(c, err.Error())
		return
	}
	http.RespondWithMessage(c, nethttp.StatusCreated, "created")
}

// UpdateGroup godoc
// @Summary Update portal group
// @Tags Portal Groups
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Group ID"
// @Param request body entities.PortalGroup true "Group data"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /api/portal/groups/{id} [put]
func (h *GroupHandler) UpdateGroup(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		http.RespondWithBadRequest(c, "invalid id")
		return
	}
	var g entities.PortalGroup
	if err := c.ShouldBindJSON(&g); err != nil {
		http.RespondWithBadRequest(c, "invalid request")
		return
	}
	g.ID = id
	g.UpdatedAt = time.Now()
	if err := h.uc.Update(c.Request.Context(), &g); err != nil {
		http.RespondWithBadRequest(c, err.Error())
		return
	}
	http.RespondWithMessage(c, nethttp.StatusOK, "updated")
}

// DeleteGroup godoc
// @Summary Delete portal group
// @Tags Portal Groups
// @Security BearerAuth
// @Produce json
// @Param id path string true "Group ID"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /api/portal/groups/{id} [delete]
func (h *GroupHandler) DeleteGroup(c *gin.Context) {
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

// ListPermissions godoc
// @Summary List permissions
// @Tags Permissions
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.SuccessResponse{data=[]*entities.Permission}
// @Router /api/portal/permissions [get]
func (h *GroupHandler) ListPermissions(c *gin.Context) {
	perms, _ := h.uc.ListPermissions(c.Request.Context())
	http.RespondWithSuccess(c, nethttp.StatusOK, perms)
}

// GetGroupPermissions godoc
// @Summary Get permissions for a group
// @Tags Permissions
// @Security BearerAuth
// @Produce json
// @Param id path string true "Group ID"
// @Success 200 {object} response.SuccessResponse{data=[]*entities.Permission}
// @Router /api/portal/groups/{id}/permissions [get]
func (h *GroupHandler) GetGroupPermissions(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		http.RespondWithBadRequest(c, "invalid id")
		return
	}
	perms, _ := h.uc.GetPermissions(c.Request.Context(), id)
	http.RespondWithSuccess(c, nethttp.StatusOK, perms)
}

type updatePermsRequest struct {
	PermissionIDs []uuid.UUID `json:"permission_ids"`
}

// UpdateGroupPermissions godoc
// @Summary Update group permissions
// @Tags Permissions
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Group ID"
// @Param request body updatePermsRequest true "Permission IDs"
// @Success 200 {object} response.SuccessResponse
// @Router /api/portal/groups/{id}/permissions [put]
func (h *GroupHandler) UpdateGroupPermissions(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		http.RespondWithBadRequest(c, "invalid id")
		return
	}
	var req updatePermsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		http.RespondWithBadRequest(c, "invalid request")
		return
	}
	if err := h.uc.UpdatePermissions(c.Request.Context(), id, req.PermissionIDs); err != nil {
		http.RespondWithBadRequest(c, err.Error())
		return
	}
	http.RespondWithMessage(c, nethttp.StatusOK, "updated")
}

package handlers

import (
	nethttp "net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"system-portal/internal/domains/portal/entities"
	"system-portal/internal/domains/portal/usecases"
	httpresp "system-portal/internal/shared/response"
)

type PermissionHandler struct{ uc usecases.PermissionUsecase }

func NewPermissionHandler(u usecases.PermissionUsecase) *PermissionHandler {
	return &PermissionHandler{uc: u}
}

// ListPermissions godoc
// @Summary List permissions
// @Tags Permissions
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.SuccessResponse{data=[]*entities.Permission}
// @Router /api/portal/permissions [get]
func (h *PermissionHandler) ListPermissions(c *gin.Context) {
	perms, err := h.uc.List(c.Request.Context())
	if err != nil {
		httpresp.RespondWithBadRequest(c, err.Error())
		return
	}
	httpresp.RespondWithSuccess(c, nethttp.StatusOK, perms)
}

// CreatePermission godoc
// @Summary Create permission
// @Tags Permissions
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body entities.Permission true "Permission data"
// @Success 201 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /api/portal/permissions [post]
func (h *PermissionHandler) CreatePermission(c *gin.Context) {
	var p entities.Permission
	if err := c.ShouldBindJSON(&p); err != nil {
		httpresp.RespondWithBadRequest(c, "invalid request")
		return
	}
	p.ID = uuid.New()
	if err := h.uc.Create(c.Request.Context(), &p); err != nil {
		httpresp.RespondWithBadRequest(c, err.Error())
		return
	}
	httpresp.RespondWithMessage(c, nethttp.StatusCreated, "created")
}

// UpdatePermission godoc
// @Summary Update permission
// @Tags Permissions
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Permission ID"
// @Param request body entities.Permission true "Permission data"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /api/portal/permissions/{id} [put]
func (h *PermissionHandler) UpdatePermission(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpresp.RespondWithBadRequest(c, "invalid id")
		return
	}
	var p entities.Permission
	if err := c.ShouldBindJSON(&p); err != nil {
		httpresp.RespondWithBadRequest(c, "invalid request")
		return
	}
	p.ID = id
	if err := h.uc.Update(c.Request.Context(), &p); err != nil {
		httpresp.RespondWithBadRequest(c, err.Error())
		return
	}
	httpresp.RespondWithMessage(c, nethttp.StatusOK, "updated")
}

// DeletePermission godoc
// @Summary Delete permission
// @Tags Permissions
// @Security BearerAuth
// @Produce json
// @Param id path string true "Permission ID"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /api/portal/permissions/{id} [delete]
func (h *PermissionHandler) DeletePermission(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpresp.RespondWithBadRequest(c, "invalid id")
		return
	}
	if err := h.uc.Delete(c.Request.Context(), id); err != nil {
		httpresp.RespondWithBadRequest(c, err.Error())
		return
	}
	httpresp.RespondWithMessage(c, nethttp.StatusOK, "deleted")
}

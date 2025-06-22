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

func (h *PermissionHandler) ListPermissions(c *gin.Context) {
	perms, _ := h.uc.List(c.Request.Context())
	httpresp.RespondWithSuccess(c, nethttp.StatusOK, perms)
}

func (h *PermissionHandler) CreatePermission(c *gin.Context) {
	var p entities.Permission
	if err := c.ShouldBindJSON(&p); err != nil {
		httpresp.RespondWithBadRequest(c, "invalid request")
		return
	}
	p.ID = uuid.New()
	h.uc.Create(c.Request.Context(), &p)
	httpresp.RespondWithMessage(c, nethttp.StatusCreated, "created")
}

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
	h.uc.Update(c.Request.Context(), &p)
	httpresp.RespondWithMessage(c, nethttp.StatusOK, "updated")
}

func (h *PermissionHandler) DeletePermission(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpresp.RespondWithBadRequest(c, "invalid id")
		return
	}
	h.uc.Delete(c.Request.Context(), id)
	httpresp.RespondWithMessage(c, nethttp.StatusOK, "deleted")
}

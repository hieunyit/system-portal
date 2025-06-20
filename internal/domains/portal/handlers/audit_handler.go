package handlers

import (
	"context"
	nethttp "net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"system-portal/internal/domains/portal/entities"
	"system-portal/internal/domains/portal/usecases"
	http "system-portal/internal/shared/response"
)

type AuditHandler struct{ uc usecases.AuditUsecase }

func NewAuditHandler(u usecases.AuditUsecase) *AuditHandler { return &AuditHandler{uc: u} }

// GetAuditLogs godoc
// @Summary List audit logs
// @Tags Audit
// @Security BearerAuth
// @Produce json
// @Success 200 {array} entities.AuditLog
// @Router /api/portal/audit/logs [get]
func (h *AuditHandler) GetAuditLogs(c *gin.Context) {
	logs, _ := h.uc.List(c.Request.Context())
	http.RespondWithSuccess(c, nethttp.StatusOK, logs)
}

// ExportAuditLogs godoc
// @Summary Export audit logs
// @Tags Audit
// @Security BearerAuth
// @Produce json
// @Success 200 {string} string "not implemented"
// @Router /api/portal/audit/logs/export [get]
func (h *AuditHandler) ExportAuditLogs(c *gin.Context) {
	http.RespondWithMessage(c, 200, "not implemented")
}

// GetAuditStats godoc
// @Summary Audit statistics
// @Tags Audit
// @Security BearerAuth
// @Produce json
// @Success 200 {string} string "not implemented"
// @Router /api/portal/audit/stats [get]
func (h *AuditHandler) GetAuditStats(c *gin.Context) {
	http.RespondWithMessage(c, 200, "not implemented")
}

// helper to add log, not used in routes
func (h *AuditHandler) addLog(action, resource string) {
	h.uc.Add(context.Background(), &entities.AuditLog{ID: uuid.New(), Action: action, Resource: resource, Success: true, CreatedAt: time.Now()})
}

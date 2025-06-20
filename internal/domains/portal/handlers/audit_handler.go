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

func (h *AuditHandler) GetAuditLogs(c *gin.Context) {
	logs, _ := h.uc.List(c.Request.Context())
	http.RespondWithSuccess(c, nethttp.StatusOK, logs)
}

func (h *AuditHandler) ExportAuditLogs(c *gin.Context) {
	http.RespondWithMessage(c, 200, "not implemented")
}

func (h *AuditHandler) GetAuditStats(c *gin.Context) {
	http.RespondWithMessage(c, 200, "not implemented")
}

// helper to add log, not used in routes
func (h *AuditHandler) addLog(action, resource string) {
	h.uc.Add(context.Background(), &entities.AuditLog{ID: uuid.New(), Action: action, Resource: resource, Success: true, CreatedAt: time.Now()})
}

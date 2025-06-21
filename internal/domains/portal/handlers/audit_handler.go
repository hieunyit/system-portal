package handlers

import (
	"bytes"
	"context"
	"encoding/csv"
	nethttp "net/http"
	"strconv"
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
	logs, _ := h.uc.List(c.Request.Context())
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	_ = w.Write([]string{"id", "user_id", "action", "resource", "success", "created_at"})
	for _, l := range logs {
		uid := ""
		if l.UserID != nil {
			uid = l.UserID.String()
		}
		_ = w.Write([]string{
			l.ID.String(),
			uid,
			l.Action,
			l.Resource,
			strconv.FormatBool(l.Success),
			l.CreatedAt.Format(time.RFC3339),
		})
	}
	w.Flush()
	c.Header("Content-Disposition", "attachment; filename=audit_logs.csv")
	c.Data(nethttp.StatusOK, "text/csv", buf.Bytes())
}

// GetAuditStats godoc
// @Summary Audit statistics
// @Tags Audit
// @Security BearerAuth
// @Produce json
// @Success 200 {string} string "not implemented"
// @Router /api/portal/audit/stats [get]
func (h *AuditHandler) GetAuditStats(c *gin.Context) {
	logs, _ := h.uc.List(c.Request.Context())
	var total, success, failed int
	for _, l := range logs {
		total++
		if l.Success {
			success++
		} else {
			failed++
		}
	}
	stats := map[string]int{"total": total, "success": success, "failed": failed}
	http.RespondWithSuccess(c, nethttp.StatusOK, stats)
}

// helper to add log, not used in routes
func (h *AuditHandler) addLog(action, resource string) {
	h.uc.Add(context.Background(), &entities.AuditLog{ID: uuid.New(), Action: action, Resource: resource, Success: true, CreatedAt: time.Now()})
}

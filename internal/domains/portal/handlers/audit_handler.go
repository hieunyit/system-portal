package handlers

import (
	"bytes"
	"encoding/csv"
	nethttp "net/http"
	"strconv"
	"time"

	"system-portal/internal/domains/portal/usecases"
	http "system-portal/internal/shared/response"

	"github.com/gin-gonic/gin"
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
	_ = w.Write([]string{"id", "user_id", "username", "user_group", "action", "resource", "ip_address", "user_agent", "duration_ms", "success", "created_at"})
	for _, l := range logs {
		uid := l.UserID.String()
		_ = w.Write([]string{
			l.ID.String(),
			uid,
			l.Username,
			l.UserGroup,
			l.Action,
			l.Resource,
			l.IPAddress,
			l.UserAgent,
			strconv.Itoa(l.DurationMs),
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

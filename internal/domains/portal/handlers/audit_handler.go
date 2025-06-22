package handlers

import (
	"bytes"
	"encoding/csv"
	nethttp "net/http"
	"strconv"
	"time"

	"system-portal/internal/domains/portal/entities"
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
type auditQuery struct {
	Username string     `form:"username"`
	Group    string     `form:"group"`
	IP       string     `form:"ip"`
	From     *time.Time `form:"from" time_format:"2006-01-02"`
	To       *time.Time `form:"to" time_format:"2006-01-02"`
	Page     int        `form:"page,default=1"`
	Limit    int        `form:"limit,default=20"`
}

func (h *AuditHandler) GetAuditLogs(c *gin.Context) {
	var q auditQuery
	_ = c.ShouldBindQuery(&q)
	filter := &entities.AuditFilter{
		Username:  q.Username,
		UserGroup: q.Group,
		IPAddress: q.IP,
		FromTime:  q.From,
		ToTime:    q.To,
		Page:      q.Page,
		Limit:     q.Limit,
	}
	logs, total, _ := h.uc.List(c.Request.Context(), filter)
	http.RespondWithSuccess(c, nethttp.StatusOK, gin.H{"logs": logs, "total": total, "page": filter.Page, "limit": filter.Limit})
}

// ExportAuditLogs godoc
// @Summary Export audit logs
// @Tags Audit
// @Security BearerAuth
// @Produce json
// @Success 200 {string} string "not implemented"
// @Router /api/portal/audit/logs/export [get]
func (h *AuditHandler) ExportAuditLogs(c *gin.Context) {
	var q auditQuery
	_ = c.ShouldBindQuery(&q)
	filter := &entities.AuditFilter{
		Username:  q.Username,
		UserGroup: q.Group,
		IPAddress: q.IP,
		FromTime:  q.From,
		ToTime:    q.To,
	}
	logs, _, _ := h.uc.List(c.Request.Context(), filter)
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	_ = w.Write([]string{"id", "user_id", "username", "user_group", "action", "resource_type", "resource_name", "ip_address", "success", "created_at"})
	for _, l := range logs {
		uid := l.UserID.String()
		_ = w.Write([]string{
			l.ID.String(),
			uid,
			l.Username,
			l.UserGroup,
			l.Action,
			l.ResourceType,
			l.ResourceName,
			l.IPAddress,
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
	var q auditQuery
	_ = c.ShouldBindQuery(&q)
	filter := &entities.AuditFilter{Username: q.Username, UserGroup: q.Group, IPAddress: q.IP, FromTime: q.From, ToTime: q.To}
	logs, _, _ := h.uc.List(c.Request.Context(), filter)
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

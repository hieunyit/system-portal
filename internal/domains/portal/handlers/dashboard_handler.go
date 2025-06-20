package handlers

import (
	"github.com/gin-gonic/gin"
	"system-portal/internal/domains/portal/dto"
	"system-portal/internal/domains/portal/repositories"
	http "system-portal/internal/shared/response"
)

type DashboardHandler struct {
	userRepo  repositories.UserRepository
	auditRepo repositories.AuditRepository
}

func NewDashboardHandler(u repositories.UserRepository, a repositories.AuditRepository) *DashboardHandler {
	return &DashboardHandler{userRepo: u, auditRepo: a}
}

func (h *DashboardHandler) GetDashboardStats(c *gin.Context) {
	users, _ := h.userRepo.List(c.Request.Context())
	http.RespondWithSuccess(c, 200, dto.StatsResponse{Users: len(users)})
}

func (h *DashboardHandler) GetRecentActivities(c *gin.Context) {
	logs, _ := h.auditRepo.List(c.Request.Context())
	resp := make([]dto.AuditResponse, 0, len(logs))
	for _, l := range logs {
		resp = append(resp, dto.AuditResponse{ID: l.ID, UserID: l.UserID, Action: l.Action, Resource: l.Resource, Success: l.Success})
	}
	http.RespondWithSuccess(c, 200, resp)
}

func (h *DashboardHandler) GetUserChartData(c *gin.Context) {
	http.RespondWithSuccess(c, 200, gin.H{})
}

func (h *DashboardHandler) GetActivityChartData(c *gin.Context) {
	http.RespondWithSuccess(c, 200, gin.H{})
}

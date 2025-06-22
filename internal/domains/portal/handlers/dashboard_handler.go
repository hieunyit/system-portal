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

// GetDashboardStats godoc
// @Summary Dashboard statistics
// @Tags Dashboard
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dto.StatsResponse
// @Router /api/portal/dashboard/stats [get]
func (h *DashboardHandler) GetDashboardStats(c *gin.Context) {
	users, _ := h.userRepo.List(c.Request.Context())
	http.RespondWithSuccess(c, 200, dto.StatsResponse{Users: len(users)})
}

// GetRecentActivities godoc
// @Summary Recent activities
// @Tags Dashboard
// @Security BearerAuth
// @Produce json
// @Success 200 {array} dto.AuditResponse
// @Router /api/portal/dashboard/activities [get]
func (h *DashboardHandler) GetRecentActivities(c *gin.Context) {
	logs, _ := h.auditRepo.List(c.Request.Context())
	resp := make([]dto.AuditResponse, 0, len(logs))
	for _, l := range logs {
		resp = append(resp, dto.AuditResponse{
			ID:           l.ID,
			UserID:       l.UserID,
			Action:       l.Action,
			Resource:     l.ResourceType,
			ResourceName: l.ResourceName,
			Success:      l.Success,
		})
	}
	http.RespondWithSuccess(c, 200, resp)
}

// GetUserChartData godoc
// @Summary User chart data
// @Tags Dashboard
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/portal/dashboard/charts/users [get]
func (h *DashboardHandler) GetUserChartData(c *gin.Context) {
	http.RespondWithSuccess(c, 200, gin.H{})
}

// GetActivityChartData godoc
// @Summary Activity chart data
// @Tags Dashboard
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/portal/dashboard/charts/activities [get]
func (h *DashboardHandler) GetActivityChartData(c *gin.Context) {
	http.RespondWithSuccess(c, 200, gin.H{})
}

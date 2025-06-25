package handlers

import (
	"github.com/gin-gonic/gin"
	"system-portal/internal/domains/portal/dto"
	"system-portal/internal/domains/portal/entities"
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
// @Description Overall portal statistics
// @Tags Dashboard
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.SuccessResponse{data=dto.StatsResponse}
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/portal/dashboard/stats [get]
func (h *DashboardHandler) GetDashboardStats(c *gin.Context) {
	users, total, _ := h.userRepo.List(c.Request.Context(), &entities.UserFilter{})
	if total == 0 {
		total = len(users)
	}
	http.RespondWithSuccess(c, 200, dto.StatsResponse{Users: total})
}

// GetRecentActivities godoc
// @Summary Recent activities
// @Description Latest portal activities
// @Tags Dashboard
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.SuccessResponse{data=[]dto.AuditResponse}
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/portal/dashboard/activities [get]
func (h *DashboardHandler) GetRecentActivities(c *gin.Context) {
	filter := &entities.AuditFilter{Page: 1, Limit: 10}
	logs, _, _ := h.auditRepo.List(c.Request.Context(), filter)
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
// @Description Chart statistics of user registrations
// @Tags Dashboard
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.SuccessResponse{data=map[string]interface{}}
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/portal/dashboard/charts/users [get]
func (h *DashboardHandler) GetUserChartData(c *gin.Context) {
	http.RespondWithSuccess(c, 200, gin.H{})
}

// GetActivityChartData godoc
// @Summary Activity chart data
// @Description Chart statistics of portal activities
// @Tags Dashboard
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.SuccessResponse{data=map[string]interface{}}
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/portal/dashboard/charts/activities [get]
func (h *DashboardHandler) GetActivityChartData(c *gin.Context) {
	http.RespondWithSuccess(c, 200, gin.H{})
}

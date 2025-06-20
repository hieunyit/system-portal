// internal/application/handlers/vpn_status_handler.go
package handlers

import (
	nethttp "net/http"
	"system-portal/internal/domains/openvpn/dto"
	"system-portal/internal/domains/openvpn/usecases"
	"system-portal/internal/shared/errors"
	http "system-portal/internal/shared/response"
	"system-portal/pkg/logger"

	"github.com/gin-gonic/gin"
)

type VPNStatusHandler struct {
	vpnStatusUsecase usecases.VPNStatusUsecase
}

func NewVPNStatusHandler(vpnStatusUsecase usecases.VPNStatusUsecase) *VPNStatusHandler {
	return &VPNStatusHandler{
		vpnStatusUsecase: vpnStatusUsecase,
	}
}

// GetVPNStatus godoc
// @Summary Get comprehensive VPN server status
// @Description Get detailed VPN server status including all connected users with their public IPs, connection times, countries, and traffic statistics
// @Tags VPN Status
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.SuccessResponse{data=dto.VPNStatusResponse} "Successful response with VPN status"
// @Failure 401 {object} response.ErrorResponse "Unauthorized - invalid or missing authentication"
// @Failure 500 {object} response.ErrorResponse "Internal server error - failed to retrieve VPN status"
// @Router /api/openvpn/vpn/status [get]
func (h *VPNStatusHandler) GetVPNStatus(c *gin.Context) {
	logger.Log.Info("Getting comprehensive VPN status")

	// Business logic thông qua usecase
	result, err := h.vpnStatusUsecase.GetVPNStatus(c.Request.Context())
	if err != nil {
		logger.Log.WithError(err).Error("Failed to get VPN status from usecase")
		http.RespondWithError(c, errors.InternalServerError("Failed to retrieve VPN status", err))
		return
	}

	// Convert usecase result to DTO
	var connectedUsers []dto.ConnectedUserResponse
	for _, user := range result.ConnectedUsers {
		connectedUsers = append(connectedUsers, dto.ConnectedUserResponse{
			CommonName:         user.CommonName,
			RealAddress:        user.RealAddress,
			VirtualAddress:     user.VirtualAddress,
			VirtualIPv6Address: user.VirtualIPv6Address,
			BytesReceived:      user.BytesReceived,
			BytesSent:          user.BytesSent,
			ConnectedSince:     user.ConnectedSince,
			ConnectedSinceUnix: user.ConnectedSinceUnix,
			Username:           user.Username,
			ClientID:           user.ClientID,
			PeerID:             user.PeerID,
			DataChannelCipher:  user.DataChannelCipher,
			Country:            user.Country,
			ConnectionDuration: user.ConnectionDuration,
		})
	}

	response := dto.VPNStatusResponse{
		TotalConnectedUsers: result.TotalConnectedUsers,
		ConnectedUsers:      connectedUsers,
		Timestamp:           result.Timestamp,
	}

	logger.Log.WithField("total_users", result.TotalConnectedUsers).
		Info("VPN status retrieved successfully")

	http.RespondWithSuccess(c, nethttp.StatusOK, response)
}

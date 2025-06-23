package handlers

import (
	"github.com/gin-gonic/gin"
	nethttp "net/http"
	"system-portal/internal/domains/portal/dto"
	"system-portal/internal/domains/portal/entities"
	"system-portal/internal/domains/portal/usecases"
	httpresp "system-portal/internal/shared/response"
)

type ConfigHandler struct {
	uc     usecases.ConfigUsecase
	reload func()
}

func NewConfigHandler(u usecases.ConfigUsecase, reload func()) *ConfigHandler {
	return &ConfigHandler{uc: u, reload: reload}
}

// GetOpenVPNConfig godoc
// @Summary Get OpenVPN connection
// @Tags Connections
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.SuccessResponse{data=entities.OpenVPNConfig}
// @Failure 404 {object} response.ErrorResponse
// @Router /api/portal/connections/openvpn [get]
func (h *ConfigHandler) GetOpenVPNConfig(c *gin.Context) {
	cfg, err := h.uc.GetOpenVPN(c.Request.Context())
	if err != nil {
		httpresp.RespondWithBadRequest(c, err.Error())
		return
	}
	if cfg == nil {
		httpresp.RespondWithNotFound(c, "not found")
		return
	}
	httpresp.RespondWithSuccess(c, nethttp.StatusOK, cfg)
}

// CreateOpenVPNConfig godoc
// @Summary Create OpenVPN connection
// @Tags Connections
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.OpenVPNConfigRequest true "OpenVPN config"
// @Success 201 {object} response.SuccessResponse
// @Router /api/portal/connections/openvpn [post]
func (h *ConfigHandler) CreateOpenVPNConfig(c *gin.Context) {
	var req dto.OpenVPNConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpresp.RespondWithBadRequest(c, "invalid request")
		return
	}

	// ensure only one configuration exists
	if existing, err := h.uc.GetOpenVPN(c.Request.Context()); err != nil {
		httpresp.RespondWithBadRequest(c, err.Error())
		return
	} else if existing != nil {
		httpresp.RespondWithBadRequest(c, "configuration already exists")
		return
	}

	cfg := &entities.OpenVPNConfig{
		Host:     req.Host,
		Username: req.Username,
		Password: req.Password,
		Port:     req.Port,
	}
	if err := h.uc.SetOpenVPN(c.Request.Context(), cfg); err != nil {
		httpresp.RespondWithBadRequest(c, err.Error())
		return
	}
	if h.reload != nil {
		h.reload()
	}
	httpresp.RespondWithMessage(c, nethttp.StatusCreated, "saved")
}

// UpdateOpenVPNConfig godoc
// @Summary Update OpenVPN connection
// @Tags Connections
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.OpenVPNConfigRequest true "OpenVPN config"
// @Success 200 {object} response.SuccessResponse
// @Router /api/portal/connections/openvpn [put]
func (h *ConfigHandler) UpdateOpenVPNConfig(c *gin.Context) {
	var req dto.OpenVPNConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpresp.RespondWithBadRequest(c, "invalid request")
		return
	}

	// ensure configuration exists before updating
	if existing, err := h.uc.GetOpenVPN(c.Request.Context()); err != nil {
		httpresp.RespondWithBadRequest(c, err.Error())
		return
	} else if existing == nil {
		httpresp.RespondWithNotFound(c, "not found")
		return
	}

	cfg := &entities.OpenVPNConfig{
		Host:     req.Host,
		Username: req.Username,
		Password: req.Password,
		Port:     req.Port,
	}
	if err := h.uc.SetOpenVPN(c.Request.Context(), cfg); err != nil {
		httpresp.RespondWithBadRequest(c, err.Error())
		return
	}
	if h.reload != nil {
		h.reload()
	}
	httpresp.RespondWithMessage(c, nethttp.StatusOK, "saved")
}

// DeleteOpenVPNConfig godoc
// @Summary Delete OpenVPN connection
// @Tags Connections
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.SuccessResponse
// @Router /api/portal/connections/openvpn [delete]
func (h *ConfigHandler) DeleteOpenVPNConfig(c *gin.Context) {
	if err := h.uc.DeleteOpenVPN(c.Request.Context()); err != nil {
		httpresp.RespondWithBadRequest(c, err.Error())
		return
	}
	if h.reload != nil {
		h.reload()
	}
	httpresp.RespondWithMessage(c, nethttp.StatusOK, "deleted")
}

// GetLDAPConfig godoc
// @Summary Get LDAP connection
// @Tags Connections
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.SuccessResponse{data=entities.LDAPConfig}
// @Failure 404 {object} response.ErrorResponse
// @Router /api/portal/connections/ldap [get]
func (h *ConfigHandler) GetLDAPConfig(c *gin.Context) {
	cfg, err := h.uc.GetLDAP(c.Request.Context())
	if err != nil {
		httpresp.RespondWithBadRequest(c, err.Error())
		return
	}
	if cfg == nil {
		httpresp.RespondWithNotFound(c, "not found")
		return
	}
	httpresp.RespondWithSuccess(c, nethttp.StatusOK, cfg)
}

// CreateLDAPConfig godoc
// @Summary Create LDAP connection
// @Tags Connections
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.LDAPConfigRequest true "LDAP config"
// @Success 201 {object} response.SuccessResponse
// @Router /api/portal/connections/ldap [post]
func (h *ConfigHandler) CreateLDAPConfig(c *gin.Context) {
	var req dto.LDAPConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpresp.RespondWithBadRequest(c, "invalid request")
		return
	}

	// ensure only one configuration exists
	if existing, err := h.uc.GetLDAP(c.Request.Context()); err != nil {
		httpresp.RespondWithBadRequest(c, err.Error())
		return
	} else if existing != nil {
		httpresp.RespondWithBadRequest(c, "configuration already exists")
		return
	}

	cfg := &entities.LDAPConfig{
		Host:         req.Host,
		Port:         req.Port,
		BindDN:       req.BindDN,
		BindPassword: req.BindPassword,
		BaseDN:       req.BaseDN,
	}
	if err := h.uc.SetLDAP(c.Request.Context(), cfg); err != nil {
		httpresp.RespondWithBadRequest(c, err.Error())
		return
	}
	if h.reload != nil {
		h.reload()
	}
	httpresp.RespondWithMessage(c, nethttp.StatusCreated, "saved")
}

// UpdateLDAPConfig godoc
// @Summary Update LDAP connection
// @Tags Connections
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.LDAPConfigRequest true "LDAP config"
// @Success 200 {object} response.SuccessResponse
// @Router /api/portal/connections/ldap [put]
func (h *ConfigHandler) UpdateLDAPConfig(c *gin.Context) {
	var req dto.LDAPConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpresp.RespondWithBadRequest(c, "invalid request")
		return
	}

	// ensure configuration exists before updating
	if existing, err := h.uc.GetLDAP(c.Request.Context()); err != nil {
		httpresp.RespondWithBadRequest(c, err.Error())
		return
	} else if existing == nil {
		httpresp.RespondWithNotFound(c, "not found")
		return
	}

	cfg := &entities.LDAPConfig{
		Host:         req.Host,
		Port:         req.Port,
		BindDN:       req.BindDN,
		BindPassword: req.BindPassword,
		BaseDN:       req.BaseDN,
	}
	if err := h.uc.SetLDAP(c.Request.Context(), cfg); err != nil {
		httpresp.RespondWithBadRequest(c, err.Error())
		return
	}
	if h.reload != nil {
		h.reload()
	}
	httpresp.RespondWithMessage(c, nethttp.StatusOK, "saved")
}

// DeleteLDAPConfig godoc
// @Summary Delete LDAP connection
// @Tags Connections
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.SuccessResponse
// @Router /api/portal/connections/ldap [delete]
func (h *ConfigHandler) DeleteLDAPConfig(c *gin.Context) {
	if err := h.uc.DeleteLDAP(c.Request.Context()); err != nil {
		httpresp.RespondWithBadRequest(c, err.Error())
		return
	}
	if h.reload != nil {
		h.reload()
	}
	httpresp.RespondWithMessage(c, nethttp.StatusOK, "deleted")
}

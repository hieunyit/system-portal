package handlers

import (
	"github.com/gin-gonic/gin"
	nethttp "net/http"
	"system-portal/internal/domains/portal/dto"
	"system-portal/internal/domains/portal/entities"
	"system-portal/internal/domains/portal/usecases"
	"system-portal/internal/shared/errors"
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
// @Description Retrieve OpenVPN connection settings
// @Tags Connections
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.SuccessResponse{data=entities.OpenVPNConfig}
// @Failure 404 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/portal/connections/openvpn [get]
func (h *ConfigHandler) GetOpenVPNConfig(c *gin.Context) {
	cfg, err := h.uc.GetOpenVPN(c.Request.Context())
	if err != nil {
		httpresp.RespondWithError(c, errors.BadRequest(err.Error(), err))
		return
	}
	if cfg == nil {
		httpresp.RespondWithError(c, errors.NotFound("not found", nil))
		return
	}
	httpresp.RespondWithSuccess(c, nethttp.StatusOK, cfg)
}

// TestOpenVPN godoc
// @Summary Test OpenVPN connection
// @Description Test connectivity to the OpenVPN server
// @Tags Connections
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.OpenVPNConfigRequest true "OpenVPN config"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/portal/connections/openvpn/test [post]
func (h *ConfigHandler) TestOpenVPN(c *gin.Context) {
	var req dto.OpenVPNConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpresp.RespondWithError(c, errors.BadRequest("invalid request", nil))
		return
	}
	cfg := &entities.OpenVPNConfig{Host: req.Host, Username: req.Username, Password: req.Password, Port: req.Port}
	if err := h.uc.TestOpenVPN(c.Request.Context(), cfg); err != nil {
		httpresp.RespondWithError(c, errors.BadRequest(err.Error(), err))
		return
	}
	httpresp.RespondWithMessage(c, nethttp.StatusOK, "ok")
}

// CreateOpenVPNConfig godoc
// @Summary Create OpenVPN connection
// @Description Create OpenVPN connection settings
// @Tags Connections
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.OpenVPNConfigRequest true "OpenVPN config"
// @Success 201 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/portal/connections/openvpn [post]
func (h *ConfigHandler) CreateOpenVPNConfig(c *gin.Context) {
	var req dto.OpenVPNConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpresp.RespondWithError(c, errors.BadRequest("invalid request", nil))
		return
	}

	// ensure only one configuration exists
	if existing, err := h.uc.GetOpenVPN(c.Request.Context()); err != nil {
		httpresp.RespondWithError(c, errors.BadRequest(err.Error(), err))
		return
	} else if existing != nil {
		httpresp.RespondWithError(c, errors.Conflict("configuration already exists", nil))
		return
	}

	cfg := &entities.OpenVPNConfig{
		Host:     req.Host,
		Username: req.Username,
		Password: req.Password,
		Port:     req.Port,
	}
	if err := h.uc.SetOpenVPN(c.Request.Context(), cfg); err != nil {
		httpresp.RespondWithError(c, errors.BadRequest(err.Error(), err))
		return
	}
	if h.reload != nil {
		h.reload()
	}
	httpresp.RespondWithMessage(c, nethttp.StatusCreated, "saved")
}

// UpdateOpenVPNConfig godoc
// @Summary Update OpenVPN connection
// @Description Update existing OpenVPN connection settings
// @Tags Connections
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.OpenVPNConfigRequest true "OpenVPN config"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/portal/connections/openvpn [put]
func (h *ConfigHandler) UpdateOpenVPNConfig(c *gin.Context) {
	var req dto.OpenVPNConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpresp.RespondWithError(c, errors.BadRequest("invalid request", nil))
		return
	}

	// ensure configuration exists before updating
	if existing, err := h.uc.GetOpenVPN(c.Request.Context()); err != nil {
		httpresp.RespondWithError(c, errors.BadRequest(err.Error(), err))
		return
	} else if existing == nil {
		httpresp.RespondWithError(c, errors.NotFound("not found", nil))
		return
	}

	cfg := &entities.OpenVPNConfig{
		Host:     req.Host,
		Username: req.Username,
		Password: req.Password,
		Port:     req.Port,
	}
	if err := h.uc.SetOpenVPN(c.Request.Context(), cfg); err != nil {
		httpresp.RespondWithError(c, errors.BadRequest(err.Error(), err))
		return
	}
	if h.reload != nil {
		h.reload()
	}
	httpresp.RespondWithMessage(c, nethttp.StatusOK, "saved")
}

// DeleteOpenVPNConfig godoc
// @Summary Delete OpenVPN connection
// @Description Remove OpenVPN connection settings
// @Tags Connections
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.SuccessResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/portal/connections/openvpn [delete]
func (h *ConfigHandler) DeleteOpenVPNConfig(c *gin.Context) {
	if err := h.uc.DeleteOpenVPN(c.Request.Context()); err != nil {
		httpresp.RespondWithError(c, errors.BadRequest(err.Error(), err))
		return
	}
	if h.reload != nil {
		h.reload()
	}
	httpresp.RespondWithMessage(c, nethttp.StatusOK, "deleted")
}

// GetLDAPConfig godoc
// @Summary Get LDAP connection
// @Description Retrieve LDAP connection settings
// @Tags Connections
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.SuccessResponse{data=entities.LDAPConfig}
// @Failure 404 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/portal/connections/ldap [get]
func (h *ConfigHandler) GetLDAPConfig(c *gin.Context) {
	cfg, err := h.uc.GetLDAP(c.Request.Context())
	if err != nil {
		httpresp.RespondWithError(c, errors.BadRequest(err.Error(), err))
		return
	}
	if cfg == nil {
		httpresp.RespondWithError(c, errors.NotFound("not found", nil))
		return
	}
	httpresp.RespondWithSuccess(c, nethttp.StatusOK, cfg)
}

// TestLDAP godoc
// @Summary Test LDAP connection
// @Description Test connectivity to the LDAP server
// @Tags Connections
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.LDAPConfigRequest true "LDAP config"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/portal/connections/ldap/test [post]
func (h *ConfigHandler) TestLDAP(c *gin.Context) {
	var req dto.LDAPConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpresp.RespondWithError(c, errors.BadRequest("invalid request", nil))
		return
	}
	cfg := &entities.LDAPConfig{Host: req.Host, Port: req.Port, BindDN: req.BindDN, BindPassword: req.BindPassword, BaseDN: req.BaseDN}
	if err := h.uc.TestLDAP(c.Request.Context(), cfg); err != nil {
		httpresp.RespondWithError(c, errors.BadRequest(err.Error(), err))
		return
	}
	httpresp.RespondWithMessage(c, nethttp.StatusOK, "ok")
}

// CreateLDAPConfig godoc
// @Summary Create LDAP connection
// @Description Create LDAP connection settings
// @Tags Connections
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.LDAPConfigRequest true "LDAP config"
// @Success 201 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/portal/connections/ldap [post]
func (h *ConfigHandler) CreateLDAPConfig(c *gin.Context) {
	var req dto.LDAPConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpresp.RespondWithError(c, errors.BadRequest("invalid request", nil))
		return
	}

	// ensure only one configuration exists
	if existing, err := h.uc.GetLDAP(c.Request.Context()); err != nil {
		httpresp.RespondWithError(c, errors.BadRequest(err.Error(), err))
		return
	} else if existing != nil {
		httpresp.RespondWithError(c, errors.Conflict("configuration already exists", nil))
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
		httpresp.RespondWithError(c, errors.BadRequest(err.Error(), err))
		return
	}
	if h.reload != nil {
		h.reload()
	}
	httpresp.RespondWithMessage(c, nethttp.StatusCreated, "saved")
}

// UpdateLDAPConfig godoc
// @Summary Update LDAP connection
// @Description Update existing LDAP connection settings
// @Tags Connections
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.LDAPConfigRequest true "LDAP config"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/portal/connections/ldap [put]
func (h *ConfigHandler) UpdateLDAPConfig(c *gin.Context) {
	var req dto.LDAPConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpresp.RespondWithError(c, errors.BadRequest("invalid request", nil))
		return
	}

	// ensure configuration exists before updating
	if existing, err := h.uc.GetLDAP(c.Request.Context()); err != nil {
		httpresp.RespondWithError(c, errors.BadRequest(err.Error(), err))
		return
	} else if existing == nil {
		httpresp.RespondWithError(c, errors.NotFound("not found", nil))
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
		httpresp.RespondWithError(c, errors.BadRequest(err.Error(), err))
		return
	}
	if h.reload != nil {
		h.reload()
	}
	httpresp.RespondWithMessage(c, nethttp.StatusOK, "saved")
}

// DeleteLDAPConfig godoc
// @Summary Delete LDAP connection
// @Description Remove LDAP connection settings
// @Tags Connections
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.SuccessResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/portal/connections/ldap [delete]
func (h *ConfigHandler) DeleteLDAPConfig(c *gin.Context) {
	if err := h.uc.DeleteLDAP(c.Request.Context()); err != nil {
		httpresp.RespondWithError(c, errors.BadRequest(err.Error(), err))
		return
	}
	if h.reload != nil {
		h.reload()
	}
	httpresp.RespondWithMessage(c, nethttp.StatusOK, "deleted")
}

// GetSMTPConfig godoc
// @Summary Get SMTP configuration
// @Description Retrieve SMTP server settings
// @Tags Connections
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.SuccessResponse{data=entities.SMTPConfig}
// @Failure 404 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/portal/connections/smtp [get]
func (h *ConfigHandler) GetSMTPConfig(c *gin.Context) {
	cfg, err := h.uc.GetSMTP(c.Request.Context())
	if err != nil {
		httpresp.RespondWithError(c, errors.BadRequest(err.Error(), err))
		return
	}
	if cfg == nil {
		httpresp.RespondWithError(c, errors.NotFound("not found", nil))
		return
	}
	httpresp.RespondWithSuccess(c, nethttp.StatusOK, cfg)
}

// CreateSMTPConfig godoc
// @Summary Create SMTP configuration
// @Description Create SMTP server configuration in database
// @Tags Connections
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.SMTPConfigRequest true "SMTP config"
// @Success 201 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/portal/connections/smtp [post]
func (h *ConfigHandler) CreateSMTPConfig(c *gin.Context) {
	var req dto.SMTPConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpresp.RespondWithError(c, errors.BadRequest("invalid request", nil))
		return
	}
	if existing, err := h.uc.GetSMTP(c.Request.Context()); err != nil {
		httpresp.RespondWithError(c, errors.BadRequest(err.Error(), err))
		return
	} else if existing != nil {
		httpresp.RespondWithError(c, errors.Conflict("configuration already exists", nil))
		return
	}
	cfg := &entities.SMTPConfig{Host: req.Host, Port: req.Port, Username: req.Username, Password: req.Password, From: req.From, TLS: req.TLS}
	if err := h.uc.SetSMTP(c.Request.Context(), cfg); err != nil {
		httpresp.RespondWithError(c, errors.BadRequest(err.Error(), err))
		return
	}
	httpresp.RespondWithMessage(c, nethttp.StatusCreated, "saved")
}

// UpdateSMTPConfig godoc
// @Summary Update SMTP configuration
// @Description Update SMTP server configuration
// @Tags Connections
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.SMTPConfigRequest true "SMTP config"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/portal/connections/smtp [put]
func (h *ConfigHandler) UpdateSMTPConfig(c *gin.Context) {
	var req dto.SMTPConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpresp.RespondWithError(c, errors.BadRequest("invalid request", nil))
		return
	}
	if existing, err := h.uc.GetSMTP(c.Request.Context()); err != nil {
		httpresp.RespondWithError(c, errors.BadRequest(err.Error(), err))
		return
	} else if existing == nil {
		httpresp.RespondWithError(c, errors.NotFound("not found", nil))
		return
	}
	cfg := &entities.SMTPConfig{Host: req.Host, Port: req.Port, Username: req.Username, Password: req.Password, From: req.From, TLS: req.TLS}
	if err := h.uc.SetSMTP(c.Request.Context(), cfg); err != nil {
		httpresp.RespondWithError(c, errors.BadRequest(err.Error(), err))
		return
	}
	httpresp.RespondWithMessage(c, nethttp.StatusOK, "saved")
}

// DeleteSMTPConfig godoc
// @Summary Delete SMTP configuration
// @Description Remove SMTP configuration from database
// @Tags Connections
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.SuccessResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/portal/connections/smtp [delete]
func (h *ConfigHandler) DeleteSMTPConfig(c *gin.Context) {
	if err := h.uc.DeleteSMTP(c.Request.Context()); err != nil {
		httpresp.RespondWithError(c, errors.BadRequest(err.Error(), err))
		return
	}
	httpresp.RespondWithMessage(c, nethttp.StatusOK, "deleted")
}

// GetEmailTemplate godoc
// @Summary Get email template by action
// @Description Retrieve an email template for the given action
// @Tags Connections
// @Security BearerAuth
// @Produce json
// @Param action path string true "action"
// @Success 200 {object} response.SuccessResponse{data=entities.EmailTemplate}
// @Failure 404 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/portal/connections/templates/{action} [get]
func (h *ConfigHandler) GetEmailTemplate(c *gin.Context) {
	action := c.Param("action")
	tpl, err := h.uc.GetTemplate(c.Request.Context(), action)
	if err != nil {
		httpresp.RespondWithError(c, errors.BadRequest(err.Error(), err))
		return
	}
	if tpl == nil {
		httpresp.RespondWithError(c, errors.NotFound("not found", nil))
		return
	}
	httpresp.RespondWithSuccess(c, nethttp.StatusOK, tpl)
}

// UpdateEmailTemplate godoc
// @Summary Update email template
// @Description Update subject or body of a template
// @Tags Connections
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param action path string true "action"
// @Param request body dto.EmailTemplateRequest true "template"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/portal/connections/templates/{action} [put]
func (h *ConfigHandler) UpdateEmailTemplate(c *gin.Context) {
	action := c.Param("action")
	var req dto.EmailTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpresp.RespondWithError(c, errors.BadRequest("invalid request", nil))
		return
	}
	tpl := &entities.EmailTemplate{Action: action, Subject: req.Subject, Body: req.Body}
	if err := h.uc.SetTemplate(c.Request.Context(), tpl); err != nil {
		httpresp.RespondWithError(c, errors.BadRequest(err.Error(), err))
		return
	}
	httpresp.RespondWithMessage(c, nethttp.StatusOK, "saved")
}

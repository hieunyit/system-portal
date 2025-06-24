package dto

// OpenVPNConfigRequest represents payload to set OpenVPN connection details
type OpenVPNConfigRequest struct {
	Host     string `json:"host" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Port     int    `json:"port" binding:"required"`
}

// LDAPConfigRequest represents payload to set LDAP connection details
type LDAPConfigRequest struct {
	Host         string `json:"host" binding:"required"`
	Port         int    `json:"port" binding:"required"`
	BindDN       string `json:"bindDN" binding:"required"`
	BindPassword string `json:"bindPassword" binding:"required"`
	BaseDN       string `json:"baseDN" binding:"required"`
}

// SMTPConfigRequest represents payload to set SMTP connection details
type SMTPConfigRequest struct {
	Host     string `json:"host" binding:"required"`
	Port     int    `json:"port" binding:"required"`
	Username string `json:"username"`
	Password string `json:"password"`
	From     string `json:"from" binding:"required"`
	TLS      bool   `json:"tls"`
}

// EmailTemplateRequest represents payload to set email template
type EmailTemplateRequest struct {
	Subject string `json:"subject" binding:"required"`
	Body    string `json:"body" binding:"required"`
}

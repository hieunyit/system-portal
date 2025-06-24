package usecases

import (
	"context"
	"system-portal/internal/domains/portal/entities"
)

type ConfigUsecase interface {
	GetOpenVPN(ctx context.Context) (*entities.OpenVPNConfig, error)
	SetOpenVPN(ctx context.Context, cfg *entities.OpenVPNConfig) error
	DeleteOpenVPN(ctx context.Context) error
	TestOpenVPN(ctx context.Context, cfg *entities.OpenVPNConfig) error
	GetLDAP(ctx context.Context) (*entities.LDAPConfig, error)
	SetLDAP(ctx context.Context, cfg *entities.LDAPConfig) error
	DeleteLDAP(ctx context.Context) error
	TestLDAP(ctx context.Context, cfg *entities.LDAPConfig) error

	GetSMTP(ctx context.Context) (*entities.SMTPConfig, error)
	SetSMTP(ctx context.Context, cfg *entities.SMTPConfig) error
	DeleteSMTP(ctx context.Context) error

	GetTemplate(ctx context.Context, action string) (*entities.EmailTemplate, error)
	SetTemplate(ctx context.Context, tpl *entities.EmailTemplate) error
}

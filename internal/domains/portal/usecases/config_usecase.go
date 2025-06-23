package usecases

import (
	"context"
	"system-portal/internal/domains/portal/entities"
)

type ConfigUsecase interface {
	GetOpenVPN(ctx context.Context) (*entities.OpenVPNConfig, error)
	SetOpenVPN(ctx context.Context, cfg *entities.OpenVPNConfig) error
	GetLDAP(ctx context.Context) (*entities.LDAPConfig, error)
	SetLDAP(ctx context.Context, cfg *entities.LDAPConfig) error
}

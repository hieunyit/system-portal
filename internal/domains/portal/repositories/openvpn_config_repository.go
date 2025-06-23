package repositories

import (
	"context"
	"system-portal/internal/domains/portal/entities"
)

type OpenVPNConfigRepository interface {
	Get(ctx context.Context) (*entities.OpenVPNConfig, error)
	Create(ctx context.Context, cfg *entities.OpenVPNConfig) error
	Update(ctx context.Context, cfg *entities.OpenVPNConfig) error
	Delete(ctx context.Context) error
}

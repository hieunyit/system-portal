package repositories

import (
	"context"
	"system-portal/internal/domains/portal/entities"
)

type LDAPConfigRepository interface {
	Get(ctx context.Context) (*entities.LDAPConfig, error)
	Create(ctx context.Context, cfg *entities.LDAPConfig) error
	Update(ctx context.Context, cfg *entities.LDAPConfig) error
}
